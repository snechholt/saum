# Common Test Baseline

## Purpose

This file defines shared defaults for implementation-facing test scenarios.

Unless a scenario says otherwise, it inherits this baseline. The baseline is intentionally minimal: add a default here when multiple scenarios rely on it, while keeping scenario-specific setup in the scenario itself.

---

## Shared defaults

* the scenario runs in an isolated environment with deterministic seeded state
* automatic deployment retries are disabled with `retry.max_retries = 0`
* the test app defines one `release` image channel with `tag_pattern = "v{version}"`
* the test app uses `image.version_ordering = "integer"`
* test versions therefore use complete OCI tags such as `v1`, `v2`, and `v3`; their extracted ordering values are `1`, `2`, and `3`
* the persisted app catalog is valid for the test app's current image configuration
* catalog entries for these versions use `channel = "release"`
* when a seeded tenant has `current_version` set, `current_image` records a valid logical image identity for that accepted version; when `current_version = null`, `current_image = null`

Tests that need automatic retry capacity explicitly override `retry.max_retries`. Tests that exercise version advancement inherit `v1 < v2 < v3` from the shared channel pattern and integer ordering unless they explicitly configure another policy.

### Seeded deployment context

When a scenario explicitly seeds an existing `deployment_context`, and does not say otherwise:

* `deployment_context.target_image` is a valid logical image identity for the scenario's `deployment_context.target_version`
* `deployment_context.retry_count = 0`
* `deployment_context.retry_after_utc = null`
* when the scenario seeds `deployment_context.attempt`, `deployment_context.attempt.recovery = null`
* when a scenario seeds recovery in progress, start from the [initial recovery state](../appendix/03_persistence_schemas.md#recovery-state) unless the scenario overrides individual recovery fields

No other tenant, runtime, routing, freeze, or deployment-context state is implied by this baseline. Those values remain explicit in each scenario where they matter.

---

### Full garbage-collection loop

When a scenario says to run the full garbage-collection loop:

1. Run the [Garbage Collection Cycle](../architecture/06_runtime_view/garbage_collection_cycle.md).
1. Advance test time past the configured grace period.
1. Run the Garbage Collection Cycle again.

The two cycles model the normal garbage-collection grace period. The first cycle records newly orphaned resources, but does not remove them immediately. After the grace period has elapsed, the second cycle can remove resources that are still eligible for collection.

Tests using this loop make garbage-collection timing deterministic with a controlled clock or equivalent. Garbage-collection-specific tests should spell out individual cycles and timing steps when those boundaries are part of the behavior under test.

---

### Materialization fingerprint fixtures

`traust.materialization_fingerprint` is an opaque value as defined by [Candidate Selection and Suitability](../appendix/13_candidate_selection_and_suitability.md). Test scenarios compare fingerprint identity; they do not prescribe how fingerprints are calculated.

Where a fixture only needs a known match or mismatch between seeded values, the fingerprints may be fictitious opaque values. When an externally seeded candidate must match a fingerprint calculated by the implementation from desired materialized inputs, setup must use the value produced or exposed by that implementation.

## Using the baseline

A scenario should list only the setup that differs from or extends this baseline. When lifecycle state selects the path under test, preconditions should state the relevant persisted values explicitly, such as `status`, `current_version`, `target_version`, and `deployment_context`. Runtime, routing, and catalog facts should likewise be explicit when the scenario depends on them. The harness remains free to choose how it creates that starting state.

Do not assert a baseline field merely because the baseline gives it a value. Assert it only when the field is part of the behavior being tested or when its value provides evidence for the expected outcome.
