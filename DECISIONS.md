# Design Decisions

Design decisions will be recorded as the project is scoped.

## Initial Project Decisions

- **Date:** 2026-10-01
- **Decision:** Treat GitOps Progressive Delivery as a four-week flagship project with a 56-hour target.
- **Why:** It demonstrates infrastructure as code, secure artifact delivery, declarative operations, progressive release safety, and SRE reasoning in one coherent portfolio story.
- **Alternatives considered:** Building separate CI, Terraform, and Kubernetes demos; expanding directly into a multi-cluster platform.
- **Trade-offs:** The project has a strong end-to-end narrative, but it will not prove enterprise-scale fleet management or production traffic behavior.

- **Date:** 2026-10-01
- **Decision:** Use Terraform, GitHub Actions, Trivy, SBOM generation, cosign, Argo CD, Argo Rollouts, and Prometheus as the core stack.
- **Why:** Together they cover provisioning, build security, artifact identity, reconciliation, progressive delivery, and automated release analysis.
- **Alternatives considered:** Direct `kubectl` deployment; a CI system that pushes directly to the cluster; service-mesh-only traffic management.
- **Trade-offs:** More controllers and concepts must fit on one VM, but the explicit boundaries create better interview evidence.

- **Date:** 2026-10-01
- **Decision:** Run the implementation on one Oracle Cloud free-tier VM with 4 vCPU and 12 GB RAM.
- **Why:** This matches the available budget and makes controller overhead, resource requests, and capacity planning part of the learning.
- **Alternatives considered:** Paid managed Kubernetes; multiple cloud nodes; local-only development.
- **Trade-offs:** Low cost and reproducibility, but the environment is a single point of failure and cannot validate highly available delivery control planes.

- **Date:** 2026-10-01
- **Decision:** Keep GitHub Actions responsible for build and security gates, and keep Argo CD responsible for cluster reconciliation.
- **Why:** Separating CI from deployment makes the GitOps trust model explicit and avoids a pipeline needing direct cluster-admin deployment access.
- **Alternatives considered:** Letting GitHub Actions run `kubectl apply`; using Argo CD only as a dashboard over imperative deployments.
- **Trade-offs:** Promotion requires a desired-state change and reconciliation, but the audit trail and access boundaries are clearer.

- **Date:** 2026-10-01
- **Decision:** Use immutable image digests and signature verification as release inputs, with short-lived or tightly scoped credentials where supported.
- **Why:** Tags can move; a digest, SBOM, scan result, and signature provide a traceable artifact identity.
- **Alternatives considered:** Mutable tags; unsigned images; long-lived registry and cluster credentials.
- **Trade-offs:** The workflow is more deliberate and requires key or identity management, but rollback and provenance are easier to defend.

- **Date:** 2026-10-01
- **Decision:** Include blue/green, canary, automated rollback, failure drills, runbooks, and interview preparation as explicit deliverables.
- **Why:** The project must demonstrate operational judgment and release safety, not only a list of installed tools.
- **Alternatives considered:** Implementing only a successful deployment path; treating interview preparation as separate theory study.
- **Trade-offs:** Some enterprise features are deferred so the completed project contains repeatable evidence of failure and recovery.

## Decision Template
- **Date:**
- **Decision:**
- **Why:**
- **Alternatives considered:**
- **Trade-offs:**
