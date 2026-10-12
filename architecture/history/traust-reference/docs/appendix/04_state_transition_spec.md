# Appendix – State Transition Spec

This appendix is the canonical reference for tenant lifecycle transitions and the persisted writes that control deployment lifecycle state.

It defines:

- allowed lifecycle transitions
- transition guards
- canonical persisted deployment operations
- recovery outcome writes

The conceptual state model lives in:

- [Tenant Deployment State](../architecture/08_cross_cutting_concepts/state_and_persistence.md)

Field shapes and meanings are canonical in [Persistence Schemas](03_persistence_schemas.md). Retry eligibility and delay calculation are canonical in [Retry Policy](11_retry_policy.md). Recovery branch selection is canonical in [Deterministic Recovery Decision Order](05_recovery_decision_order.md). Target supersession is defined in [Version Targeting](../architecture/08_cross_cutting_concepts/version_targeting.md).

---

## State Transition Table

| From | To | Guard | Persist timing |
|---|---|---|---|
| `pending` | `deploying` | activation validation and [new-context prerequisites](#new-deployment-context) succeed | Before deployment side effects |
| `running` | `deploying` | `target_version != current_version`; [new-context prerequisites](#new-deployment-context) succeed | Before deployment side effects |
| `deploying` | `failed` | recovery persists a failed outcome | Per recovery outcome |
| `deploying` | `running` | deployment work completed successfully | Before unfreezing the routed target runtime |
| `failed` | `deploying` | deployment retry or superseding-target work is eligible | Before deployment side effects |
| `failed` | `running` | deployer recovery proves successful convergence to the deployment-context target | Before unfreezing the routed target runtime |
| `failed` | `running` | operator _Abandon failed target_; deployment context is terminal and safely settled; `current_version != null`; `target_version == deployment_context.target_version` | As one atomic administrative write after guard validation; no deployment side effects |

Until the `deploying` → `running` transition has been persisted, failures are handled through recovery. Recovery determines whether the tenant can transition to `running`, return to a safe failed state eligible for deployment retry, or remain `failed` with persisted recovery progress.

Recovery can also persist progress while the tenant remains `failed`. Those same-state writes are defined below rather than represented as lifecycle transitions.

### Activation validation

Before a tenant may transition from `pending` to `deploying`, the deployer validates the minimum tenant and app configuration required to begin deployment.

The tenant record must define:

- `tenant`
- `app`
- `hostnames`, with at least one entry
- `update_policy`
- `status`

The app referenced by `tenant.app` must exist and define:

- `app`
- `deployment_mode`
- `image.registry`
- `image.repository`
- `image.version_ordering`

Required fields must satisfy their normal field definitions. These are the minimum activation-validation requirements and are not an exhaustive definition of configuration validity. Additional requirements may be added as deployment requirements evolve.

---

## Persisted Write Model

Unless an operation below says otherwise, tenant-record fields remain unchanged.

### New deployment attempt

A new deployment attempt replaces `deployment_context.attempt` with:

- a new `id`
- `baseline` captured from the currently routed runtime when `current_version` is set, or `null` when `current_version = null`
- `candidate = null`
- `backup_reference = null`
- `recovery = null`

The new attempt is persisted before deployment side effects begin. Recovery never replaces the deployment attempt it is settling.

### New deployment context

Before creating a deployment context:

- the selected `target_version` and its image digest must exist in a catalog that is valid for the app's current image configuration; the canonical catalog-validity rule is defined in [Image Discovery and Targeting - Catalog validity](../architecture/06_runtime_view/image_discovery_and_targeting.md#catalog-validity)
- when `current_version = null`, no tenant route may exist
- when `current_version` is set, the tenant must have a current route to the runtime that will be captured as the rollback baseline
- when `current_version` is set, the [backup capacity policy](06_queue_and_backup_policy.md#backup-artifact-policy) must permit another backup

If a prerequisite is unavailable, the context is not created and deployment side effects do not begin.

Starting work for a new deployment context sets:

- `status = deploying`
- `deployment_context.target_version` from the selected target
- `deployment_context.target_image` from the persisted catalog image source and selected target entry
- `deployment_context.retry_count = 0`
- `deployment_context.retry_after_utc = null`
- `deployment_context.attempt` to a new deployment attempt

The complete context is persisted before deployment side effects begin. The context target version and target image identity remain immutable for the lifetime of that context.

This write is used when deployment first starts from `pending` or `running`, and when a different desired target supersedes a terminal failed context.

### Finalize running

Successful convergence sets:

- `current_version = deployment_context.target_version`
- `current_image = deployment_context.target_image`
- `status = running`
- `deployment_context = null`

The finalized record is persisted before unfreezing the routed target runtime. Clearing `deployment_context` also removes its deployment retry, attempt, backup, and recovery state.

### Canonical persisted operations

| Operation | Required persisted write | Persist timing |
|---|---|---|
| Select target version | After the target-change rules in [Version Targeting](../architecture/08_cross_cutting_concepts/version_targeting.md#target-version-changes) succeed, set `target_version` to the selected version and set `update_policy = manual`. | As one atomic tenant-record write after guard validation. |
| Set update policy to manual | Set `update_policy = manual`. Leave `target_version` unchanged. | As one atomic tenant-record write. |
| Set update policy to auto | After the policy-change rules in [Version Targeting](../architecture/08_cross_cutting_concepts/version_targeting.md#target-version-changes) succeed, set `update_policy = auto`; if automatic resolution selected a newer target, set `target_version` to it in the same write. | As one atomic tenant-record write after guard validation. |
| Create deployment context | Apply **New deployment context** above. | Before deployment side effects. |
| Record selected candidate | Set `deployment_context.attempt.candidate.runtime_instance_id` to the selected runtime, replacing any previous candidate reference. | After candidate selection and before the readiness check or later deployment work uses that candidate. |
| Start automatic deployment retry | Set `status = deploying`. Increment `deployment_context.retry_count`. Clear `deployment_context.retry_after_utc`. Replace `deployment_context.attempt` with a **New deployment attempt**. | Before deployment side effects. |
| Start manual deployment retry | Set `status = deploying`. Clear `deployment_context.retry_after_utc`. Replace `deployment_context.attempt` with a **New deployment attempt**. | Before deployment side effects. |
| Enter migration boundary | Set `deployment_context.attempt.backup_reference` to the canonical path of the completed backup artifact. | After the backup command succeeds and the artifact is non-empty, and before `migrate` is invoked. |
| Initialize recovery | Set `deployment_context.attempt.recovery` to the [initial recovery state](03_persistence_schemas.md#recovery-state). | After recovery determines that the attempt cannot be finalized successfully; before recovery side effects. |
| Start recovery retry | Increment `recovery.retry_count`. Clear `recovery.retry_after_utc`. The deployment attempt and deployment retry fields are unchanged. | Before recovery side effects. |
| Schedule deployment retry | Set `deployment_context.retry_after_utc` according to Retry Policy. | After a safe failed outcome. |
| Schedule recovery retry | Apply the applicable [Recovery Outcome Write](#recovery-outcome-writes) when the failed operation changes the recovery outcome. If `recovery.status` changed during the pass, reset `recovery.retry_count = 0`. Set `recovery.retry_after_utc` to the time calculated by Retry Policy. | As one atomic tenant-record write after a failed recovery operation, before ending the reconciliation pass. |
| Finalize running | Apply **Finalize running** above. | Before unfreezing the routed target runtime. |
| Abandon failed target | Set `status = running`, set `target_version = current_version`, and clear `deployment_context`. If `update_policy = auto`, set it to `manual`. | As one atomic administrative write after guard validation. |

Only starting an automatic deployment retry increments `deployment_context.retry_count`. Scheduling an automatic retry or starting a manual deployment retry leaves it unchanged, so the counter records automatic retries that have actually started.

**Abandon failed target** is a single administrative operation that validates its transition guard and applies its complete persisted write. It must not be reproduced by independently changing the affected tenant fields.

Manual **Retry recovery** uses the same **Start recovery retry** write. It only bypasses the scheduled-time gate for that pass.

---

## Recovery Outcome Writes

The table below defines the lifecycle and recovery fields for each outcome. For a failed recovery operation, **Schedule recovery retry** combines the applicable outcome with its retry metadata in one atomic tenant-record write. A safe baseline write must be durable before any routed-runtime unfreeze.

| Outcome | Required persisted result |
|---|---|
| Safe pre-migration settlement | Set `status = failed`. Keep `backup_reference = null`, `recovery.status = null`, and `recovery.retry_count` unchanged. Clear `recovery.retry_after_utc` and `deployment_context.attempt.candidate`. |
| Failed data recovery | Set `status = failed` and `recovery.status = revert_failed_data`. Keep `backup_reference` unchanged. |
| Failed routing recovery | Set `status = failed` and `recovery.status = revert_failed_routing`. Keep `backup_reference` unchanged. |
| Known-safe application and routing baseline committed | Set `status = failed` and `recovery.status = reverted_safe`. Reset `recovery.retry_count = 0`. Clear `recovery.retry_after_utc`, `deployment_context.attempt.candidate`, and `backup_reference`. |
| Failed unfreeze after safe recovery | Keep `status = failed`. Set `recovery.status = revert_failed_unfreeze`. Keep `backup_reference` clear. |
| Proven successful deployment | Apply **Finalize running**. |

**Schedule recovery retry** also applies to first-deployment `restore --initial` and route-removal failures.

For pre-migration recovery, a routed-runtime unfreeze failure does not introduce another recovery-status value. The tenant remains `failed` with `backup_reference = null` and `recovery.status = null`; only recovery retry metadata changes in the failed-operation write.

First-deployment recovery uses the same `revert_failed_data`, `revert_failed_routing`, and `reverted_safe` outcomes. Its `backup_reference` is already `null`, and it has no baseline runtime for `revert_failed_unfreeze`.

---

## Structural Invariants

- `deployment_context` is present if and only if `status` is `deploying` or `failed`.
- `current_version = null` if and only if `current_image = null`.
- `status = pending` requires `current_version = null`.
- `status = running` requires `current_version` to be set.
- When `deployment_context` is present, `deployment_context.attempt.baseline = null` if and only if `current_version = null`.
- A deployment context's target version and target image identity do not change in place.
- Recovery keeps the same deployment attempt. Deployment retry creates a new attempt only after the previous attempt is safe to retry.
- A new deployment attempt starts with no selected candidate, backup reference, or recovery state.
- A first-deployment attempt (`baseline = null`) has no `backup_reference` and cannot have `recovery.status = revert_failed_unfreeze`.
- `recovery.status = reverted_safe` or `revert_failed_unfreeze` requires `backup_reference = null`.
- For an attempt with a baseline, `recovery.status = revert_failed_data` or `revert_failed_routing` requires `backup_reference` to be set.
- Migration-boundary work must be settled before another deployment attempt begins.

---

## Persistence Mutation Contract

Every persisted state mutation must follow:

- [Appendix – Persistence Mutation Contract](08_persistence_write_contract.md)
