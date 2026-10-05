# gitops-progressive-delivery

A small delivery platform on one Oracle free-tier VM. Terraform builds the VM,
cloud-init turns it into a single-node k3s cluster, GitHub Actions builds and
signs the app image, Argo CD keeps the cluster matching this repo, and Argo
Rollouts moves traffic to new versions in steps with Prometheus deciding whether
to keep going.

I built it to have something concrete to point at when talking about release
safety and supply chain. It is deliberately one node. The interesting parts are
the boundaries between the pieces and what happens when a release is bad, not
the scale.

## How a change gets out

Say I change a handler in `app/`.

1. The PR runs `go vet`, the tests, renders every kustomize overlay, and
   validates the Terraform. Nothing is built or pushed from a PR.
2. On merge to `main` the image job builds
   `ghcr.io/pranshu-raj/gitops-progressive-delivery:<sha>`, pushes it, scans
   that digest with Trivy (CRITICAL and HIGH fail the job), writes a CycloneDX
   SBOM of the same digest, then signs the digest with cosign using the
   workflow's OIDC identity and attaches the SBOM as an attestation. A failed
   scan leaves an unsigned image in the registry, which is as good as no image
   once verification is in the path.
3. Nothing deploys yet. CI has no cluster credentials at all.
4. To ship it I open a second PR that changes the image reference in
   `deploy/overlays/prod/kustomization.yaml`. That diff is the whole release.
   Whoever reviews it sees exactly which digest is going where.
5. Argo CD notices the merge, renders the overlay, and applies it. The Rollout
   spec changed, so Argo Rollouts creates a new ReplicaSet and starts the canary.
6. The canary goes 25%, 50%, 75% with two-minute pauses between steps. From the
   first step on, an AnalysisRun asks Prometheus every 30s for the canary's
   non-5xx ratio. One failed measurement aborts the rollout and traffic goes
   back to the stable ReplicaSet on its own. The last step is an indefinite
   pause, so the final promote is a person running `kubectl argo rollouts promote`.
7. Rolling back is `git revert` of the overlay change. Argo CD syncs it,
   Rollouts shifts back. The cluster never has a version that isn't in git.

Dev is the same up to step 5 but uses blue/green instead of canary. The new
version comes up behind `demo-preview`, I poke at it, and promote by hand.
Dev also self-heals: if someone edits a live object, Argo CD puts it back.
Prod only flags the drift.

## What's in the repo

    app/         the service, Go, standard library only
    terraform/   VCN, subnet, security list, the VM, cloud-init that installs k3s
    deploy/      kustomize base + dev/prod overlays. This is the desired state.
    platform/    Argo CD, Argo Rollouts, Prometheus, and the two Argo CD Applications
    scripts/     bootstrap.sh and a curl loop for generating load
    docs/        runbooks, filled in as failure drills get done

Ownership is the point of the split. `app/` produces artifacts and never
touches the cluster. `terraform/` is applied by hand from a laptop and knows
nothing about what runs on the node. `deploy/overlays/<env>` is the only thing
Argo CD reads, and the only way anything changes in the cluster is a merge
there. `platform/` is applied once by `bootstrap.sh` and then left alone.

## The app

`GET /` returns version, environment and hostname as JSON. `/healthz` is
liveness and always answers 200. `/readyz` is readiness and flips to 503 the
moment the process gets SIGTERM, a few seconds before it stops listening, so the
endpoint slice drains before connections are cut. `/metrics` is Prometheus text,
written by hand because pulling in client_golang for two metrics felt like
overkill and it keeps the SBOM down to the base image.

Two environment variables exist purely to make rollouts fail on purpose:

    FAIL_EVERY=5    every fifth request to / returns 500
    SLOW_MS=800     every request to / takes an extra 800ms

Both come from a ConfigMap that kustomize generates with a content hash in its
name. Changing a value in an overlay changes the ConfigMap name, which changes
the pod template, which starts a rollout. That is how the failure drills work:
bump `FAIL_EVERY` in prod, watch the analysis fail, watch it abort.

Every log line carries `version`, `env` and a request id (`X-Request-Id` if the
client sent one, random otherwise). Every metric carries a `version` label, and
Prometheus copies the `rollouts-pod-template-hash` pod label on scrape, so
canary and stable are separable in queries without the app knowing anything
about rollouts.

## Running it

You need an OCI tenancy with a compartment, the OCI CLI config at
`~/.oci/config`, terraform, kubectl, and the Argo Rollouts kubectl plugin.

    cp terraform/terraform.tfvars.example terraform/terraform.tfvars
    terraform -chdir=terraform init
    terraform -chdir=terraform apply

Fill in the region, compartment OCID, your SSH public key, and the /32 you are
sitting on. That CIDR is the only thing allowed to reach port 22 and 6443. The
app itself is never exposed publicly.

Give cloud-init a couple of minutes, then pull the kubeconfig:

    ip=$(terraform -chdir=terraform output -raw public_ip)
    ssh ubuntu@$ip sudo cat /etc/rancher/k3s/k3s.yaml | sed "s/127.0.0.1/$ip/" > ~/.kube/config
    kubectl get nodes

Then the platform:

    bash scripts/bootstrap.sh

That applies Argo CD, Argo Rollouts and Prometheus from `platform/`, waits for
them to come up, applies the two Application objects, and prints the Argo CD
admin password. From then on the cluster follows `main`.

Argo CD UI: `kubectl -n argocd port-forward svc/argocd-server 8080:443`, then
https://localhost:8080.

Terraform state is local and gitignored. For one person and one VM that is
fine. `terraform destroy` takes everything down; there is nothing on the node
worth keeping.

## Day to day

    make test       vet + tests
    make image      local docker build tagged with the short sha
    make render     kustomize build both overlays, to see what Argo CD will apply
    make tf-check   terraform fmt + validate

Watching a rollout:

    kubectl argo rollouts get rollout demo -n demo-prod --watch

The analysis needs traffic to measure. Easiest is a port-forward and the curl
loop:

    kubectl -n demo-prod port-forward svc/demo 8081:80 &
    scripts/load.sh http://localhost:8081/ 0.1

Promote a paused rollout with `kubectl argo rollouts promote demo -n demo-prod`.
Abort one with `kubectl argo rollouts abort demo -n demo-prod`. Abort moves
traffic back to stable but leaves the Rollout Degraded until the spec changes,
so the real rollback is still a revert in git.

## Checking an image

    cosign verify ghcr.io/pranshu-raj/gitops-progressive-delivery@sha256:... \
      --certificate-identity-regexp '^https://github.com/Pranshu-Raj/Gitops-Progressive-Delivery/' \
      --certificate-oidc-issuer https://token.actions.githubusercontent.com

Swap `verify` for `verify-attestation --type cyclonedx` to get the SBOM back.
The identity regexp is the trust decision: only this repo's workflows, signed
through GitHub's OIDC issuer, count.

CI also pushes a moving `main` tag. The overlays start on it only because there
is no digest until the first build runs. Nothing should stay on it. Argo CD
does not resync when a tag moves, because the manifest did not change, which is
exactly why digests are the thing you pin.

## Things to know before you change stuff

OCI's Ubuntu images ship with iptables rules that drop everything except SSH,
separately from the security list in Terraform. cloud-init opens 6443 in
iptables before installing k3s. If you add a NodePort later you have to open it
in both places.

kustomize does not know that a Rollout has a pod template inside it.
`deploy/base/rollout-transform.yaml` tells it where ConfigMap and Service names
live in a Rollout so the generated ConfigMap's hash suffix gets rewritten.
Without that file the pod's `envFrom` points at a ConfigMap that does not exist
and the rollout never becomes Ready.

There is no traffic router, so canary weight is replica count. 25% of four
replicas is one pod, and the Service round-robins across whatever is Ready, so
the real traffic share is whatever the endpoint slice says it is, not the number
in the step.

The analysis query divides canary successes by canary requests. With no traffic
that is NaN, and the `isNaN(result[0]) ||` in the success condition stops an
idle canary from being killed for lack of evidence. That is a choice: no data is
treated as fine. Making sure there is load during a rollout is on me.

Prometheus is one pod with an emptyDir and two days of retention. It restarts,
it forgets. Fine here, not fine anywhere real.

## What this is not

One VM. No HA for the API server, Argo CD, Prometheus, or anything else. No
ingress, no TLS in front of the app, no service mesh, no external secrets store.
Production traffic has never touched it and the numbers in the analysis
template are picked for a demo, not tuned against a baseline. The deferred
list in `ROADMAP.md` says what would change first if any of that mattered.

## Where the rest lives

`ROADMAP.md` is the plan and how far along it is. `DECISIONS.md` has the why
behind each tool choice. `LEARNING_LOG.md` is what I got wrong on the way.
`docs/runbooks/` fills in as the failure drills get run.
