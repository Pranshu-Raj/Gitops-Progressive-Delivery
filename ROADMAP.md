# Learning-By-Building Roadmap

## Learning Profile
- **What I want to learn:** SRE and platform-engineering interview reasoning through secure GitOps and progressive delivery.
- **Why:** Build a credible resume project while preparing for SRE interviews and landing a job.
- **Time available:** 50-60 hours over 4 weeks, alongside 3-5 additional job-search projects.
- **Budget for cloud/tools:** Free tier only. Oracle Cloud free-tier VM: 4 vCPU and 12 GB RAM.
- **Project idea:** GitOps Progressive Delivery

## Target Project

Build a small Kubernetes environment provisioned with Terraform. GitHub Actions builds an application image, scans it with Trivy, generates an SBOM, signs it with cosign, and pushes it to a registry. Argo CD synchronizes desired state declaratively, while Argo Rollouts performs blue/green and canary releases with Prometheus-driven analysis and automated rollback.

This project must produce interview evidence: supply-chain decisions, GitOps drift demonstrations, rollout experiments, rollback timelines, runbooks, and questions answered from firsthand failures. It is a portfolio project, not a claim of production-grade high availability.

## Project Selection

The project is selected. The implementation is intentionally single-node and resource-aware so it fits the available VM and leaves time for interview preparation.

## Milestone Overview

| Milestone | Focus | Est. time |
|---|---|---:|
| 1 | Delivery model and baseline application | 8 hours |
| 2 | Terraform and Kubernetes foundation | 10 hours |
| 3 | CI supply-chain security | 10 hours |
| 4 | Argo CD GitOps operations | 10 hours |
| 5 | Argo Rollouts and progressive delivery | 10 hours |
| 6 | Operations, portfolio packaging, and interviews | 8 hours |
| **Total** |  | **56 hours** |

## Milestone 1: Delivery Model and Baseline Application

### CP-1.1: Project charter and delivery architecture
**Status:** ⏳ Not started
**Goal:** Define the application, environments, delivery flow, trust boundaries, and project limits.
**You'll learn:** GitOps, progressive-delivery vocabulary, architecture decisions, and scope control.
**Done when:**
- [ ] A charter defines the application, environments, promotion path, and rollback goal.
- [ ] A diagram shows Terraform, Kubernetes, GitHub Actions, registry, cosign, Argo CD, Argo Rollouts, and Prometheus.
- [ ] Trust boundaries identify where credentials, signatures, and deployment decisions cross systems.
- [ ] Non-goals include multi-node HA, multi-region delivery, and production-scale traffic.
**Verify it:** Explain how a commit becomes a running version and where each system makes decisions.
**Read first:** [Argo CD Core Concepts](https://argo-cd.readthedocs.io/en/stable/core_concepts/), [Argo Rollouts](https://argo-rollouts.readthedocs.io/en/stable/)
**Est. time:** 2 hours

<details><summary>Hint 1</summary>Draw desired-state flow and runtime flow separately.</details>
<details><summary>Hint 2</summary>Mark who can build, sign, approve, sync, promote, and roll back.</details>
<details><summary>Hint 3</summary>Describe one safe release and one unsafe release from commit to recovery.</details>

**Stretch (optional):** Add a threat model for the registry, signing identity, CI identity, and cluster credentials.

### CP-1.2: Baseline application and release contract
**Status:** ⏳ Not started
**Goal:** Create a small service with visible version, health, and intentional failure behavior.
**You'll learn:** Release contracts, readiness semantics, version traceability, and rollback conditions.
**Done when:**
- [ ] Health and readiness behavior are defined.
- [ ] The running version is visible in a response or diagnostic endpoint.
- [ ] One deterministic latency or error mode exists for rollout testing.
- [ ] Tests cover normal behavior, invalid input, dependency failure, and version reporting.
- [ ] Logs identify version, environment, and request or operation ID.
**Verify it:** Run two versions and prove which version served a request.
**Read first:** [Kubernetes Probes](https://kubernetes.io/docs/concepts/configuration/liveness-readiness-startup-probes/)
**Est. time:** 3 hours

<details><summary>Hint 1</summary>Define observable conditions that make a release acceptable before building the rollout.</details>
<details><summary>Hint 2</summary>Keep the failure mode deterministic and reversible.</details>
<details><summary>Hint 3</summary>Expose version identity in both user-facing behavior and operational evidence.</details>

**Stretch (optional):** Add one idempotent operation and document its retry behavior.

### CP-1.3: Repository and environment contract
**Status:** ⏳ Not started
**Goal:** Separate application source, infrastructure, platform configuration, and environment-specific desired state.
**You'll learn:** Repository boundaries, promotion models, configuration ownership, and drift prevention.
**Done when:**
- [ ] The repository layout distinguishes application, Terraform, Kubernetes base, overlays, and GitOps manifests.
- [ ] Promotion from development to production-like state is documented.
- [ ] Image tags or digests have a clear ownership and update policy.
- [ ] No secrets or private keys are committed.
- [ ] A clean checkout identifies source, generated output, and desired state.
**Verify it:** Classify five sample changes as application, infrastructure, or GitOps changes and justify each.
**Read first:** [Kubernetes Declarative Management](https://kubernetes.io/docs/tasks/manage-kubernetes-objects/declarative-config/)
**Est. time:** 2 hours

<details><summary>Hint 1</summary>For every file category, name its owner and consumer.</details>
<details><summary>Hint 2</summary>Prefer immutable image references for promotion and rollback.</details>
<details><summary>Hint 3</summary>Trace one change through Git history and identify the reproducible commit.</details>

**Stretch (optional):** Write a promotion policy for development, staging, and production-like environments.

### CP-1.4: Milestone 1 review
**Status:** ⏳ Not started
**Goal:** Consolidate the delivery and release model before provisioning infrastructure.
**You'll learn:** Reconciliation, deployment safety, and failure-domain reasoning.
**Done when:**
- [ ] You explain push-based versus pull-based deployment.
- [ ] You distinguish deployment, release, promotion, rollback, and roll-forward.
- [ ] You identify three delivery-path failure modes.
- [ ] You explain why a signed image is not automatically a secure deployment.
**Verify it:** Answer four questions aloud in 15 minutes using your diagram.
**Read first:** [OpenGitOps Principles](https://opengitops.dev/)
**Est. time:** 1 hour

<details><summary>Hint 1</summary>Use your own diagram as the subject of each answer.</details>
<details><summary>Hint 2</summary>Ask what happens if CI, registry, Argo CD, or the cluster is unavailable.</details>
<details><summary>Hint 3</summary>Separate authenticity, authorization, availability, and correctness.</details>

**Stretch (optional):** Write five interviewer follow-ups for your deployment model.

## Milestone 2: Terraform and Kubernetes Foundation

### CP-2.1: Terraform VM and network provisioning
**Status:** ⏳ Not started
**Goal:** Provision the free-tier host and required network access reproducibly with Terraform.
**You'll learn:** Infrastructure as code, state, variables, outputs, provider authentication, and network boundaries.
**Done when:**
- [ ] Terraform provisions or documents the Oracle VM and network rules.
- [ ] Variables and outputs separate environment values from module logic.
- [ ] State handling and locking limitations are documented.
- [ ] Plan, apply, and destroy workflows are documented and safely tested.
- [ ] Administrative access is restricted and secrets are absent from code.
**Verify it:** Run formatting, validation, plan, and a controlled apply; record resource identifiers.
**Read first:** [Terraform CLI](https://developer.hashicorp.com/terraform/cli), [Oracle VCN Overview](https://docs.oracle.com/en-us/iaas/Content/Network/Concepts/overview.htm)
**Est. time:** 3 hours

<details><summary>Hint 1</summary>Start with the smallest resource graph that proves recreation.</details>
<details><summary>Hint 2</summary>Remember that Terraform state may contain sensitive values even when code does not.</details>
<details><summary>Hint 3</summary>Test destroy only after identifying disposable and backed-up data.</details>

**Stretch (optional):** Add a safe Terraform plan artifact to CI.

### CP-2.2: Kubernetes installation and resource budget
**Status:** ⏳ Not started
**Goal:** Install a small Kubernetes distribution and establish budgets for workloads and controllers.
**You'll learn:** Cluster bootstrap, scheduling, requests, limits, and capacity trade-offs.
**Done when:**
- [ ] Kubernetes installation is repeatable on the VM.
- [ ] Namespaces separate application, GitOps, and rollout resources.
- [ ] CPU and memory requests and limits exist for workloads and key controllers.
- [ ] Only required interfaces are exposed.
- [ ] A smoke test confirms the cluster and application namespace are healthy.
**Verify it:** Recreate the cluster from the documented process and capture node, pod, and resource status.
**Read first:** [Kubernetes Resource Management](https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/)
**Est. time:** 3 hours

<details><summary>Hint 1</summary>Reserve headroom for the OS, Kubernetes, Argo CD, Argo Rollouts, and Prometheus.</details>
<details><summary>Hint 2</summary>Requests affect scheduling; limits affect runtime enforcement.</details>
<details><summary>Hint 3</summary>Verify controllers one at a time so resource problems are not mistaken for application problems.</details>

**Stretch (optional):** Produce a capacity table for every always-running component.

### CP-2.3: Kubernetes manifests and configuration
**Status:** ⏳ Not started
**Goal:** Define workload, service, probes, configuration, and immutable version identity declaratively.
**You'll learn:** Workload configuration, health gates, configuration separation, and manifest validation.
**Done when:**
- [ ] Deployment or Rollout-ready manifests define workload, service, probes, and resources.
- [ ] Configuration is externalized and secrets are excluded from Git.
- [ ] The image reference can change without editing application source.
- [ ] Invalid configuration prevents readiness and produces an actionable log.
- [ ] Rendered manifests validate locally.
**Verify it:** Apply the manifests, change only the image reference, and confirm the reported version changes.
**Read first:** [Kubernetes Deployments](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/), [Kubernetes ConfigMaps](https://kubernetes.io/docs/concepts/configuration/configmap/)
**Est. time:** 2 hours

<details><summary>Hint 1</summary>Make image, replicas, probes, resources, config, and service routing explicit.</details>
<details><summary>Hint 2</summary>Readiness protects users from a live but unready new version.</details>
<details><summary>Hint 3</summary>Validate the rendered result, not only a template fragment.</details>

**Stretch (optional):** Add policy checks for privileged containers, missing resources, and mutable tags.

### CP-2.4: Infrastructure validation review
**Status:** ⏳ Not started
**Goal:** Make Terraform and Kubernetes fail early before CI and deployment.
**You'll learn:** Static validation, policy checks, drift awareness, and safe infrastructure changes.
**Done when:**
- [ ] Terraform formatting and validation pass.
- [ ] Kubernetes manifests render and validate.
- [ ] A deliberately invalid resource and manifest are rejected clearly.
- [ ] Drift and state recovery procedures are documented.
- [ ] Blocking and advisory checks are distinguished.
**Verify it:** Introduce one controlled defect in each layer and confirm both are caught.
**Read first:** [Terraform Validate](https://developer.hashicorp.com/terraform/cli/commands/validate), [Kubernetes Declarative Management](https://kubernetes.io/docs/tasks/manage-kubernetes-objects/declarative-config/)
**Est. time:** 2 hours

<details><summary>Hint 1</summary>Choose checks based on defects they prevent.</details>
<details><summary>Hint 2</summary>Validate syntax, then semantics, then policy assumptions.</details>
<details><summary>Hint 3</summary>Document state and cluster recovery before you need it.</details>

**Stretch (optional):** Add a reproducibility check for expected resource categories.

## Milestone 3: CI Supply-Chain Security

### CP-3.1: GitHub Actions build and test pipeline
**Status:** ⏳ Not started
**Goal:** Build a traceable image and reject changes that fail tests or quality checks.
**You'll learn:** CI stages, immutable artifacts, least-privilege permissions, and build provenance.
**Done when:**
- [ ] Pull requests run tests and static checks.
- [ ] A successful workflow builds an image linked to a commit.
- [ ] Workflow permissions are minimal and documented.
- [ ] Failed steps preserve useful diagnostics.
- [ ] The pipeline does not deploy directly to the cluster.
**Verify it:** Break a test and confirm CI blocks the artifact; inspect image metadata for source identity.
**Read first:** [GitHub Actions Security Hardening](https://docs.github.com/en/actions/security-for-github-actions/security-guides/security-hardening-for-github-actions)
**Est. time:** 3 hours

<details><summary>Hint 1</summary>Separate test, build, scan, sign, and publish stages by trust boundary.</details>
<details><summary>Hint 2</summary>Use immutable artifact identifiers and avoid unnecessary registry or deployment permissions.</details>
<details><summary>Hint 3</summary>Ask what an attacker could change if a pull request workflow had deployment credentials.</details>

**Stretch (optional):** Add dependency caching while documenting cache risks.

### CP-3.2: Trivy scanning and SBOM generation
**Status:** ⏳ Not started
**Goal:** Scan the image and generate an inventory of its components before publication.
**You'll learn:** Vulnerability policy, false positives, SBOM purpose, and supply-chain gates.
**Done when:**
- [ ] Trivy scans the image or filesystem in CI.
- [ ] An SBOM is generated and stored as an artifact.
- [ ] Severity thresholds and exceptions have ownership and expiry.
- [ ] A deliberately vulnerable dependency demonstrates a failing gate.
- [ ] Scan and SBOM results reference the same image digest.
**Verify it:** Add a controlled vulnerability, observe failure, then remove or explicitly justify it.
**Read first:** [Trivy Documentation](https://trivy.dev/latest/), [CycloneDX SBOM](https://cyclonedx.org/)
**Est. time:** 3 hours

<details><summary>Hint 1</summary>Decide what blocks delivery before reviewing findings.</details>
<details><summary>Hint 2</summary>Generate the SBOM from the artifact that will be deployed.</details>
<details><summary>Hint 3</summary>Exceptions must explain risk acceptance and an expiration condition.</details>

**Stretch (optional):** Compare image, repository, and IaC scanning coverage.

### CP-3.3: Cosign signing and verification
**Status:** ⏳ Not started
**Goal:** Sign the published image and prove verification rejects unsigned or altered artifacts.
**You'll learn:** Artifact identity, signing trust, verification policy, and key or identity recovery.
**Done when:**
- [ ] The published image is signed after scanning.
- [ ] Verification succeeds for the expected digest.
- [ ] Verification fails for an unsigned image or mismatched identity.
- [ ] Trust assumptions and rotation or recovery are documented.
- [ ] Private signing material is absent from Git and logs.
**Verify it:** Verify the expected digest, then change the digest or identity and capture rejection.
**Read first:** [Sigstore Cosign](https://docs.sigstore.dev/cosign/signing/signing_with_containers/)
**Est. time:** 2 hours

<details><summary>Hint 1</summary>Signing proves an identity made a statement about a digest; it does not prove the code is correct.</details>
<details><summary>Hint 2</summary>Define the trusted identity and issuer before implementing verification.</details>
<details><summary>Hint 3</summary>Test unsigned, wrong repository, wrong identity, and changed-digest paths.</details>

**Stretch (optional):** Attach the SBOM or provenance as a signed OCI artifact.

### CP-3.4: Supply-chain review
**Status:** ⏳ Not started
**Goal:** Demonstrate an auditable chain from source commit to signed image digest.
**You'll learn:** Provenance, artifact promotion, CI threat modeling, and evidence collection.
**Done when:**
- [ ] One run links commit, tests, digest, scan, SBOM, signature, and registry location.
- [ ] Failed tests, scans, and signing block publication.
- [ ] Workflow permissions and third-party actions are reviewed.
- [ ] The supply-chain threat model is recorded in `DECISIONS.md`.
**Verify it:** Reconstruct the artifact story using only workflow logs, artifacts, and registry metadata.
**Read first:** [SLSA Overview](https://slsa.dev/spec/v1.0/)
**Est. time:** 2 hours

<details><summary>Hint 1</summary>Follow the digest, not the tag, across every step.</details>
<details><summary>Hint 2</summary>List what the pipeline proves and cannot prove.</details>
<details><summary>Hint 3</summary>Review every credential by asking why the job needs it and what log exposure would mean.</details>

**Stretch (optional):** Add provenance attestation and a verification report.

## Milestone 4: Argo CD GitOps Operations

### CP-4.1: GitOps repository and desired state
**Status:** ⏳ Not started
**Goal:** Organize declarative application state so Argo CD can reconcile it predictably.
**You'll learn:** Git as source of truth, overlays, immutable promotion, and manifest review.
**Done when:**
- [ ] GitOps manifests are separate from application build logic.
- [ ] Desired image digest and rollout configuration are reviewable in Git.
- [ ] Environment differences are explicit and minimal.
- [ ] A pull request represents a complete desired-state change.
- [ ] Secrets use a documented safe mechanism.
**Verify it:** Change desired state in Git and identify exactly what Argo CD should change.
**Read first:** [Argo CD Declarative Setup](https://argo-cd.readthedocs.io/en/stable/operator-manual/declarative-setup/)
**Est. time:** 2 hours

<details><summary>Hint 1</summary>Keep generated or mutable state out of the source-of-truth path.</details>
<details><summary>Hint 2</summary>Use image digests or an update mechanism that leaves an auditable Git commit.</details>
<details><summary>Hint 3</summary>Review one change as an approver: impact, rollback, and security consequences should be visible.</details>

**Stretch (optional):** Use an application-of-applications pattern only if it reduces real complexity.

### CP-4.2: Argo CD bootstrap and synchronization
**Status:** ⏳ Not started
**Goal:** Install Argo CD and synchronize the application from Git.
**You'll learn:** Reconciliation loops, sync status, health status, and controller debugging.
**Done when:**
- [ ] Argo CD runs within the VM budget.
- [ ] An Argo CD Application points to the GitOps path.
- [ ] Desired state becomes synchronized without imperative deployment commands.
- [ ] Application health and sync status are observable.
- [ ] Bootstrap and teardown are documented.
**Verify it:** Make a Git change, observe OutOfSync and then Synced, and verify the application.
**Read first:** [Argo CD Getting Started](https://argo-cd.readthedocs.io/en/stable/getting_started/)
**Est. time:** 3 hours

<details><summary>Hint 1</summary>Learn the difference between desired state, live state, sync, and health.</details>
<details><summary>Hint 2</summary>When sync fails, inspect controller events, rendered resources, and Kubernetes events separately.</details>
<details><summary>Hint 3</summary>Keep the first Application small; prove reconciliation before rollout complexity.</details>

**Stretch (optional):** Add a read-only operational view for reviewers.

### CP-4.3: Drift, self-heal, and rollback
**Status:** ⏳ Not started
**Goal:** Demonstrate detection of live changes and recovery through Git.
**You'll learn:** Drift detection, self-healing, reconciliation timing, and Git-based rollback.
**Done when:**
- [ ] A safe manual live-state change produces visible drift.
- [ ] The configured policy restores or flags drift as documented.
- [ ] A Git revert returns the application to the previous desired version.
- [ ] Drift, sync, and rollback evidence is captured.
- [ ] Drift is distinguishable from application failure.
**Verify it:** Change one safe live field, observe drift, restore through Git, and verify final state.
**Read first:** [Argo CD Automated Sync](https://argo-cd.readthedocs.io/en/stable/user-guide/auto_sync/)
**Est. time:** 3 hours

<details><summary>Hint 1</summary>Choose a reversible field that cannot break cluster access.</details>
<details><summary>Hint 2</summary>Record whether Argo CD flags, overwrites, or ignores the change.</details>
<details><summary>Hint 3</summary>Rollback the Git commit so the source of truth becomes correct again.</details>

**Stretch (optional):** Compare automatic self-healing with approval-based recovery.

### CP-4.4: Milestone 4 review
**Status:** ⏳ Not started
**Goal:** Explain GitOps behavior under delivery, drift, controller failure, and rollback.
**You'll learn:** Reconciliation reasoning and operational diagnosis.
**Done when:**
- [ ] You explain why Argo CD is not a CI system.
- [ ] You distinguish OutOfSync, Degraded, Failed, and Healthy.
- [ ] You explain what breaks if Git is unavailable during steady state.
- [ ] You describe recovery from a bad desired-state commit.
**Verify it:** Answer four questions using evidence from the drift experiment.
**Read first:** [Argo CD User Guide](https://argo-cd.readthedocs.io/en/stable/user-guide/)
**Est. time:** 2 hours

<details><summary>Hint 1</summary>Start with what the controller knows locally and what it must fetch.</details>
<details><summary>Hint 2</summary>Separate current runtime state from the ability to make future changes.</details>
<details><summary>Hint 3</summary>Name the safe recovery path and evidence that confirms it worked.</details>

**Stretch (optional):** Write an incident scenario for Argo CD being unavailable during a rollout.

## Milestone 5: Argo Rollouts and Progressive Delivery

### CP-5.1: Blue/green rollout
**Status:** ⏳ Not started
**Goal:** Release a second version with preview validation and controlled promotion.
**You'll learn:** Replica sets, active and preview services, promotion gates, and rollback.
**Done when:**
- [ ] A Rollout defines active and preview services.
- [ ] The new version receives preview validation without immediately replacing active traffic.
- [ ] Promotion is explicit and observable.
- [ ] Abort or rollback restores the previous version.
- [ ] The flow is documented with evidence.
**Verify it:** Deploy A, create B, validate preview, promote B, then abort an intentionally bad B.
**Read first:** [Argo Rollouts Blue/Green](https://argo-rollouts.readthedocs.io/en/stable/features/bluegreen/)
**Est. time:** 3 hours

<details><summary>Hint 1</summary>Identify which service selects active traffic and which selects preview traffic.</details>
<details><summary>Hint 2</summary>Define validation before promotion; otherwise preview is only a second copy.</details>
<details><summary>Hint 3</summary>Observe replica sets, service selectors, and rollout status during transitions.</details>

**Stretch (optional):** Add an automated preview smoke test as a promotion gate.

### CP-5.2: Canary rollout and traffic steps
**Status:** ⏳ Not started
**Goal:** Shift traffic gradually and make each canary step measurable.
**You'll learn:** Canary strategy, exposure control, promotion pauses, and progressive risk.
**Done when:**
- [ ] A canary strategy defines at least three steps.
- [ ] Each step has a visible pause, promotion condition, or analysis decision.
- [ ] Canary and stable versions are identifiable in metrics and logs.
- [ ] Manual abort returns traffic to stable.
- [ ] Resource usage remains within the VM budget.
**Verify it:** Promote a healthy canary through all steps and abort another midway.
**Read first:** [Argo Rollouts Canary](https://argo-rollouts.readthedocs.io/en/stable/features/canary/)
**Est. time:** 3 hours

<details><summary>Hint 1</summary>Choose steps large enough to measure but small enough to limit blast radius.</details>
<details><summary>Hint 2</summary>Make version identity a bounded metric label or equivalent signal.</details>
<details><summary>Hint 3</summary>Watch for traffic imbalance: replica percentage is not always request percentage.</details>

**Stretch (optional):** Compare replica-based canary with ingress or service-mesh routing as a design exercise.

### CP-5.3: Prometheus analysis and automated rollback
**Status:** ⏳ Not started
**Goal:** Use objective metrics to pause, promote, or abort a canary automatically.
**You'll learn:** Analysis templates, metric windows, statistical caution, and automated remediation.
**Done when:**
- [ ] Prometheus measures canary success or error behavior.
- [ ] An AnalysisTemplate evaluates the canary against a stated threshold and window.
- [ ] A healthy canary passes and a deliberately bad canary fails.
- [ ] Failed analysis aborts or rolls back automatically.
- [ ] Query, threshold, decision, and limitations are documented.
**Verify it:** Run passing and failing analyses and capture metric-to-rollout timelines.
**Read first:** [Argo Rollouts Analysis](https://argo-rollouts.readthedocs.io/en/stable/features/analysis/), [Prometheus Querying](https://prometheus.io/docs/prometheus/latest/querying/basics/)
**Est. time:** 3 hours

<details><summary>Hint 1</summary>Start with one user-impact metric with a trustworthy denominator.</details>
<details><summary>Hint 2</summary>Choose a window that avoids startup noise but limits exposure.</details>
<details><summary>Hint 3</summary>Test the analysis query independently before embedding it in the rollout.</details>

**Stretch (optional):** Add a second analysis metric and define disagreement behavior.

### CP-5.4: Progressive-delivery failure drill
**Status:** ⏳ Not started
**Goal:** Exercise bad releases, analysis failure, manual abort, and recovery.
**You'll learn:** Failure domains, rollback semantics, blast-radius control, and evidence-based response.
**Done when:**
- [ ] A bad image or controlled failure causes analysis to fail.
- [ ] Automated rollback restores stable.
- [ ] Manual abort is demonstrated separately.
- [ ] A metric-source or controller failure is analyzed as a limitation.
- [ ] Detection, decision, mitigation, and recovery are timed.
**Verify it:** Repeat the failure drill from its runbook and record recovery time.
**Read first:** [Argo Rollouts Troubleshooting](https://argo-rollouts.readthedocs.io/en/stable/troubleshooting/)
**Est. time:** 1 hour

<details><summary>Hint 1</summary>Use a failure whose blast radius is limited to the canary.</details>
<details><summary>Hint 2</summary>Distinguish a bad release from unavailable Prometheus.</details>
<details><summary>Hint 3</summary>Recovery is complete only when traffic, desired state, rollout status, and user behavior agree.</details>

**Stretch (optional):** Add a signal for an unavailable analysis provider.

## Milestone 6: Operations, Portfolio Packaging, and Interviews

### CP-6.1: Security and operational runbooks
**Status:** ⏳ Not started
**Goal:** Make the delivery system safe to operate and hand off.
**You'll learn:** Least privilege, secret handling, rollback runbooks, and operational readiness.
**Done when:**
- [ ] Runbooks cover failed CI, failed scan, signature failure, OutOfSync, stuck rollout, failed analysis, and rollback.
- [ ] CI, registry, Argo CD, and cluster permissions are documented.
- [ ] Secrets and signing identities have rotation or recovery guidance.
- [ ] Runbooks include symptoms, checks, mitigation, escalation, and recovery validation.
- [ ] Credentials are absent from Git, images, logs, and artifacts.
**Verify it:** Execute two runbooks without undocumented personal knowledge.
**Read first:** [Kubernetes Secrets](https://kubernetes.io/docs/concepts/configuration/secret/), [GitHub Actions OIDC](https://docs.github.com/en/actions/security-for-github-actions/security-hardening-your-deployments/about-security-hardening-with-openid-connect)
**Est. time:** 2 hours

<details><summary>Hint 1</summary>Review permissions by job and controller, not only by platform.</details>
<details><summary>Hint 2</summary>Every failure runbook should say whether retrying is safe.</details>
<details><summary>Hint 3</summary>Validate both the control plane and user-facing service after recovery.</details>

**Stretch (optional):** Replace a long-lived CI credential with short-lived identity federation.

### CP-6.2: Capacity and reliability review
**Status:** ⏳ Not started
**Goal:** Measure delivery-platform cost and explain scaling limits.
**You'll learn:** Capacity planning, controller overhead, headroom, and trade-offs.
**Done when:**
- [ ] CPU, memory, disk, pod count, registry traffic, and rollout overhead are measured.
- [ ] Two capacity scenarios are modeled.
- [ ] Headroom assumptions are stated.
- [ ] The first bottleneck is justified with evidence.
- [ ] A scale-out design explains what would move to managed or multi-node infrastructure.
**Verify it:** Produce a one-page capacity report with measurements and recommendation.
**Read first:** [Google SRE: Handling Overload](https://sre.google/sre-book/handling-overload/)
**Est. time:** 2 hours

<details><summary>Hint 1</summary>Measure controllers and monitoring as part of the platform.</details>
<details><summary>Hint 2</summary>Connect the first saturation signal to delivery or user impact.</details>
<details><summary>Hint 3</summary>State what evidence would justify another node or managed controller.</details>

**Stretch (optional):** Design a ten-times-larger delivery architecture without implementing it.

### CP-6.3: Documentation and resume evidence
**Status:** ⏳ Not started
**Goal:** Turn the project into reproducible documentation and defensible resume claims.
**You'll learn:** Technical writing, quantification, trade-off communication, and honest scope framing.
**Done when:**
- [ ] README covers Terraform, bootstrap, CI gates, signing, GitOps sync, rollouts, rollback, runbooks, and limitations.
- [ ] A clean-environment walkthrough fixes ambiguous steps.
- [ ] Three resume bullets use measured outcomes.
- [ ] A two-minute explanation and ten-minute deep dive are written.
- [ ] Every claim is supported by a test, experiment, metric, or screenshot.
**Verify it:** Deliver the two-minute explanation and remove unsupported claims.
**Read first:** [GitHub: About READMEs](https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/customizing-your-repository/about-readmes)
**Est. time:** 2 hours

<details><summary>Hint 1</summary>Lead with the release-safety problem and measured outcome.</details>
<details><summary>Hint 2</summary>Quantify gates, rollback time, failure modes, or resource usage only when measured.</details>
<details><summary>Hint 3</summary>Use the single-node limitation as a trade-off story, not a hidden weakness.</details>

**Stretch (optional):** Record a five-minute demo from commit to automated rollback.

### CP-6.4: SRE mock interview and final demo
**Status:** ⏳ Not started
**Goal:** Demonstrate the platform and answer SRE questions from firsthand evidence.
**You'll learn:** System design, debugging, supply-chain reasoning, incident response, and behavioral storytelling.
**Done when:**
- [ ] You answer questions about Linux, networking, Kubernetes, Terraform state, CI security, GitOps, progressive delivery, and rollback.
- [ ] You design a multi-node or multi-region version and explain the changes.
- [ ] You explain timeout, retry, idempotency, backpressure, and partial-failure trade-offs.
- [ ] You prepare stories about reducing deployment risk, handling a failed rollout, and learning from a mistake.
- [ ] A 30-minute demo goes from environment to CI artifact to Argo CD sync to rollout failure and recovery.
- [ ] At least 80% of prepared questions are answered correctly or marked for revisit.
**Verify it:** Complete a timed technical, system-design, and behavioral mock interview, then update `LEARNING_LOG.md`.
**Read first:** [Google SRE Workbook](https://sre.google/workbook/)
**Est. time:** 2 hours

<details><summary>Hint 1</summary>Answer from your own failure drill first, then generalize.</details>
<details><summary>Hint 2</summary>State assumptions, failure modes, observability, mitigation, and trade-offs.</details>
<details><summary>Hint 3</summary>When asked to scale, preserve deployment safety before adding tools.</details>

**Stretch (optional):** Ask another engineer to review the promotion and trust model.

## Required Quality Coverage

- [ ] Application tests, Terraform validation, Kubernetes validation, and CI failure paths are covered.
- [ ] Trivy scanning, SBOM generation, cosign signing, and verification are demonstrated.
- [ ] Argo CD sync, drift detection, recovery, and Git rollback are verified.
- [ ] Blue/green, canary, Prometheus analysis, automated rollback, and manual abort are verified.
- [ ] Resource usage, disk growth, controller overhead, and rollout capacity are measured.
- [ ] Runbooks and documentation work from a clean environment.

## Deliberately Deferred

Do not make these primary implementation goals in this four-week project:

- Multi-node or multi-region high availability.
- Production-scale traffic, long retention, or public production exposure.
- Enterprise policy engines, service mesh, external secrets platforms, or multi-cluster fleet management.
- Key-management infrastructure beyond a documented safe signing approach.
- Automatic promotion based on statistically mature production baselines.

Keep deferred items as design exercises. Explain why they were deferred and what evidence would justify adding them.
