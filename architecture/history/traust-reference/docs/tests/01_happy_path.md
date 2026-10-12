# Scenario 01 – Happy Path Deployment

## Intent

Validate a full successful upgrade flow where image discovery finds a newer version and this triggers tenant convergence from current version to target version.

This scenario also validates the three-stage readiness/verification sequence:

1. candidate readiness check before migration
1. post-migration internal verification before routing
1. hostname verification after routing

---

## Preconditions

* start from the [Common Test Baseline](00_common_test_baseline.md)
* persisted app catalog contains `v1` for the `release` channel with its immutable image digest
* registry contains `v1` and a newer `v2` tag with distinct immutable image digests; `v2` is not yet present in the persisted app catalog
* tenant has `status = running`, `current_version = v1`, `target_version = v1`, and `deployment_context = null`
* routing points to the existing successful `v1` runtime, which is healthy and unfrozen for the tenant
* tenant update policy is set to `auto`
* backup/migrate/verification and `appctl` command contracts are available and functioning for the test fixture versions

---

## Steps

1. Trigger image discovery (or run one main worker loop that includes discovery + reconciliation evaluation).
1. Wait for deployer processing to complete for the target tenant.
1. Collect observable outputs (tenant record, routing target, and relevant runtime signals).

In this scenario, image discovery is expected to:
* update tenant `target_version` (because policy is `auto`)
* trigger the deployment flow.

---

## Expected Outcomes

* `v2` is discovered as the newer release target
* the persisted catalog records `v2` with `channel = release` and its resolved immutable image digest
* `v2` is selected for update
* the candidate readiness check is performed before entering the migration boundary
* post-migration verification is not performed before migration
* the candidate treats the tenant as frozen without an explicit candidate freeze command
* `appctl freeze` is called for the currently routed runtime before backup
* post-migration internal verification is performed after migration and before routing switch
* hostname verification is performed after routing switch through the real tenant hostname
* deployment state is committed as `running` before `appctl unfreeze` is called for the routed target runtime
* tenant is updated and state converges to `running`
* `current_version` is set to `target_version`
* container for `target_version` exists and is healthy
* `deployment_context` is cleared
* routing points to container for `target_version`
* previously active container is no longer active in routing for this tenant
* the routed target runtime is unfrozen and the inactive previous runtime remains frozen for this tenant

---

## Persisted-State Checks

Verify tenant record after completion:

* `status = running`
* `current_version = target_version`
* `deployment_context = null`

Verify persisted discovery and target-resolution state:

* target version exists in persisted available image/version set
* tenant `target_version` is set to the selected/discovered version

