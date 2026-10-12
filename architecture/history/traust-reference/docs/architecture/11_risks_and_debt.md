# 11 Risks and Technical Debt

## Risk register

| ID | Risk | Why it matters | Mitigation | Revisit trigger |
|---|---|---|---|---|
| R-01 | Single-daemon bottleneck | Limits rollout throughput and recovery parallelism. | Accepted single-worker model with deterministic order. | Queue delay or recovery time becomes operationally unacceptable. |
| R-02 | Recovery misclassification | Can cause unsafe deployment retry or unnecessary restore. | Migration boundary, observed reality, canonical recovery order, tests. | Ambiguous timeout/crash cases are not covered by tests or logs. |
| R-03 | Routing correctness dependency | Misroutes affect cutover safety and garbage-collection safety. | Atomic route publication, hostname verification for cutover, and frozen candidates during first-deployment recovery. | Route mismatch or routing convergence becomes operationally significant. |
| R-04 | Backup availability/capacity | Missing backups make recovery unsafe after the migration boundary has been entered. | Require backup completion before migration; enforce retention and capacity limits. | Backup cap blocks deployments or restore reliability is uncertain. |
| R-05 | Contract drift | Divergence between architecture, specs, appendices, tests, and implementation creates unsafe behavior. | Architecture/test traceability, canonical references, strict docs build. | Specs reveal principles absent from architecture, or tests contradict runtime docs. |
| R-06 | Verification phase misuse | Running full app/env sanity before migration can produce false failures. | Separate the candidate readiness check from post-migration verification in the concepts and contracts. | App teams blur verification responsibilities or pre-migration checks depend on migrated tenant state. |
| R-07 | File-backed state corruption | File-backed state may become incomplete or otherwise invalid after an interrupted write or filesystem failure. | Persistence Mutation Contract and fail-closed rejection of missing or unsupported `schema_version` values. | Any state-file corruption or evidence that an interrupted write left invalid state. |
| R-08 | Application freeze contract failure | A runtime that loses or bypasses freeze state can accept tenant mutation during migration. | Persistent fail-closed freeze state, post-restart verification, and deployment containment on contract violation. | Any mutation while frozen, or a runtime restart after which a previously frozen tenant reports unfrozen. |
| R-09 | Unattributed shared-runtime work | Background work without tenant identity cannot honor a tenant-scoped freeze. | Require tenant attribution before mutation, isolated deployment, or one deployable tenant for all affected customers. | An app cannot map a mutation-capable path to a tenant before execution. |

## Technical debt themes

### D-01 - Limited concurrency

Intentional for safety and simplicity. Revisit only when measured queue delay or recovery-time pressure justifies redesign.

### D-02 - Minimal deployment retry taxonomy

Deployment retry uses a deliberately simple bounded retry/backoff model without a general transient/permanent classification. Revisit if repeated failures show a clear need for such classification.

### D-03 - Log-first observability

Observability is log-first. Revisit if diagnostics become too slow or manual.

### D-04 - Single-host deployment

Appropriate for the intended environment. Revisit only if control-plane HA or multi-host orchestration becomes a product requirement.

## Guardrails

- Behavior-critical docs must preserve runtime ordering.
- Architecture owns canonical architectural rules and rationale, specs own implementation contracts and mechanics, and appendices own exact canonical reference contracts.
- Tests should be linked from quality requirements without rewriting them unnecessarily.
