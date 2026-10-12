## Problem

Zonotron normally assumes an authority graph where upstream artifacts constrain downstream ones, for example:

```text
Architecture
   ├──→ Design
   └──→ Test Harness
             \
              → Implementation
```

Authority flows downward. If Architecture changes, Design/Test Harness/Implementation may become stale and reconcile against the new upstream state.

The awkward case is onboarding an existing legacy system that has only code.

Call the legacy software project **Cosmos**:

```text
Implementation I0
```

There is no settled architecture, design, behavioral model, or test harness. Much of the actual system knowledge exists only implicitly in the implementation.

The obvious onboarding process is to work upward:

```text
Implementation
   ↓
infer Design
   ↓
infer Architecture
```

This raises a semantic problem: during reconstruction, Design and Architecture need to conform to the existing implementation, which seems to reverse Zonotron's normal authority relationships.

That would be undesirable. If Design genuinely depended on Implementation, later implementation changes could incorrectly make Design stale. More importantly, existing code should not automatically become normative: legacy implementation contains bugs, accidental behavior, debt, and historical decisions that should not necessarily be preserved.

The same problem applies to tests. Tests generated from existing behavior may correctly characterize what Cosmos currently does without establishing that this behavior is actually required.

## Proposed solution: reconstruction is evidence flow, not reversed authority

Do **not** reverse the normal artifact graph.

Instead, distinguish:

1. **normative authority/dependency**, and
2. **evidence/provenance used to reconstruct an artifact**.

During onboarding:

```text
                    Implementation I0
                     /             \
             evidence               evidence
                ↓                      ↓
       reconstructed Design     characterization tests
                 \                    /
                  \                  /
                   candidate Architecture
                         +
                    human intent
```

The arrows upward mean:

> this candidate artifact was inferred using evidence from the existing implementation.

They do **not** mean:

> this artifact is normatively governed by the implementation.

This lets Zonotron preserve its normal semantics while still climbing upward from code.

## Current behavior versus intended behavior

A useful distinction is between **characterization** and **acceptance**.

A characterization harness answers:

> What does Cosmos demonstrably do today?

For example:

```text
DELETE /user currently returns 500
```

A characterization test may record that fact.

An acceptance harness answers:

> What have we decided Cosmos is required to do?

The 500 response might later be classified as a bug rather than a required behavior.

The same distinction applies to Design.

Zonotron can reconstruct an **as-built design** from the implementation:

```text
AuthService directly manipulates billing tables
```

That is an observation about Cosmos.

Human review may instead settle:

```text
Authentication must not own billing persistence.
```

The reconstructed artifact therefore starts as evidence/candidate state rather than immediately becoming authority.

Useful artifact statuses during onboarding could conceptually be:

```text
observed / reconstructed
candidate
settled
```

## The onboarding pivot

Legacy onboarding can be seen as six stages:

```text
1. IMPORT
   Implementation I0

2. OBSERVE
   Recover behavior, components, interfaces,
   data flows, deployment structure, etc.

3. RECONSTRUCT
   Produce candidate Design, Test/Behaviour,
   Architecture, etc.
   Mark uncertainty, contradictions, accidental behavior.

4. HUMAN SETTLEMENT
   Decide which reconstructed facts actually represent
   the system we intend to continue owning.

5. PIVOT
   The settled higher-level artifacts now become
   normative upstream authority.

6. RECONCILE
   Reconcile the existing Implementation I0
   against the newly settled graph.
```

The key point is that there is **no actual edge reversal**.

Before settlement:

```text
Implementation → evidence → candidate Design/Architecture
```

After settlement:

```text
Architecture/Design/Test → authority → Implementation
```

At the pivot, the newly reconstructed and reviewed artifacts become the basis for future change.

Implementation may immediately become stale against them:

```text
Design D0          SETTLED
Implementation I0  differs from D0
                   → reconciliation required
```

That is desirable: onboarding may intentionally reveal existing architectural debt.

## Architecture cannot be recovered from code alone

Implementation can reveal a lot:

- components;
- actual dependencies;
- APIs/interfaces;
- schemas/data;
- runtime behavior;
- deployment topology;
- implementation constraints.

But it usually cannot reliably reveal:

- business goals;
- architectural rationale;
- quality priorities;
- organizational constraints;
- intentionally supported versus accidental behavior;
- rejected alternatives;
- known debt;
- what should change rather than be preserved.

So Architecture reconstruction should combine:

```text
code-derived evidence
+
reconstructed design/behavior
+
human explanation of Cosmos's purpose
+
which existing behaviors must remain
+
known bugs/debt
+
business/operational constraints
        ↓
candidate Architecture A0
```

Architecture settlement is therefore inherently partly human.

## Zonotron can propose the target artifact graph

A legacy project might initially contain only:

```text
Implementation
```

During onboarding, Zonotron can inspect Cosmos and propose possible governance structures, for example:

```text
Architecture
    ↓
Implementation
```

or:

```text
Architecture
 ├──→ Design
 └──→ Test Harness
          \
           → Implementation
```

or a richer graph with a behavioral model, security model, deployment model, etc.

The criterion should not be "what documentation can we generate?"

It should be:

> Which bodies of system knowledge are important enough that we want them to exist as independently settled authority?

A small legacy service may need only Architecture → Implementation. A large Cosmos-like system may justify several layers.

## Why this fits Zonotron

This reinforces the idea that Zonotron is about **reconciliation and accumulated authoritative state**, not document generation.

Legacy onboarding is not:

> Generate architecture docs from this repository.

It is:

> Recover enough explicit knowledge from an organically evolved downstream artifact that we can establish a governed artifact graph without discarding what the existing system already knows.

Once onboarding finishes:

```text
legacy Cosmos:
Implementation was effectively the only source of truth

             ↓ reconstruction + settlement

governed Cosmos:
Architecture A0
Design D0
Test Harness T0
Implementation I0/I1

             ↓ future change

Architecture A1
→ normal Zonotron reconciliation
```

This suggests adding **bootstrap/reconstruction** as a first-class Zonotron lifecycle posture alongside foundation discovery, buildout, and maintenance.

Its purpose is:

> Establish the first settled authority graph when the project's knowledge currently exists mainly in mature downstream artifacts rather than explicit upstream specifications.
