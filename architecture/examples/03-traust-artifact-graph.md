# Example 03. Traust Architecture, Behaviour Model, Harness and Design

**Status:** proposed composition, not an implemented Allium integration or a settled dependency policy.

## 1. Existing documentation context

The supplied review archive includes Traust's documentation snapshot at the revision used for v5. Its Architecture Index is explicit about ownership: architecture rules/rationale, specification contracts/mechanics, appendix exact references and test scenarios are distinct even when stored together.

The risk register names contract drift as R-05. Scenario 01 and the common baseline show that substantial behavior is already specified, including lifecycle state, event order, public boundaries and intentionally flexible fixture setup.

This makes the repository a plausible source for a behavioral projection. It also means a projection must respect the existing internal authority map rather than treating every sentence as equally normative.

## 2. User's proposed addition

The intended development order was to create the independent Test Harness from architecture scenarios before implementing the SUT. The user then proposed adding a Behaviour Model, potentially implemented in Allium, as another Harness dependency.

Architecture remains broadly accessible for rationale. The model likewise becomes useful context for Design and Implementation. The goal is more efficient recovery of intended behavior, not removing access to the richer documents.

A concise candidate input map is:

```text
Behaviour Model:
  Architecture/docs

Test Harness:
  Architecture/docs
  Behaviour Model

Design:
  Architecture/docs
  Behaviour Model

Implementation:
  Design
  Test Harness
  Architecture/docs and Behaviour Model also visible
```

The last line intentionally does not settle whether visibility means direct conformance dependencies, pinned transitive context, or another relationship. That was proposed but not resolved in the main thread.

## 3. What each artifact contributes

Architecture/docs define the product context, external technologies and contracts, relevant scenarios, exact reference rules and rationale.

The Behaviour Model makes applicable commitments explicit in a constrained behavioral representation and exposes missing or conflicting decisions during its creation and review.

Design establishes the structural ownership and interactions that realize those commitments.

Test Harness establishes how independent evidence will demonstrate them while respecting a replaceable SUT and valid external boundaries.

Implementation supplies the working deployer, constrained by the accepted inputs.

These are independent responsibilities, not necessarily separate teams or repositories in every system. Their usefulness as separate nodes depends on whether the user wants to govern their state independently.

## 4. Example: image discovery

The documentation states that a newer eligible release is discovered through the relevant registry/authentication process and affects catalog/target state.

The model can state the behavioral conditions and transitions. Design can decide component responsibility and interactions. Harness must create an appropriate external registry surface and observations without replacing the relevant SUT behavior with an internal fiction. Implementation performs discovery and transition correctly.

If Harness cannot determine credential provisioning semantics, that can become an upstream clarification. If running a registry is merely annoying, that is not automatically a specification problem. If independent modeling also cannot express the credential behavior unambiguously, the combined evidence may justify stronger upstream review.

The choice between a full registry and a smaller faithful fixture remains a design/validation question, not a settled preference inferred from this example.

## 5. Example: running state before unfreeze

Scenario 01 requires running state to be committed before unfreeze is called for the target runtime. A model may express the intended order, but Harness still needs an observation strategy that demonstrates it rather than merely reading the final state afterward.

Design must make the relevant sequence implementable. The implementation can satisfy the rule using more than one valid mechanism. The Harness must not accidentally require one private invocation or identity convention.

This shows why precise behavioral obligations are useful but do not solve evidence soundness by themselves.

## 6. Discovery policy versus accepted shared interpretation

During initial validation, separate model, harness and implementation probes can independently interpret the same documentation. Their disagreements may expose ambiguous intent or permitted discretion.

Once a model revision is deliberately accepted, common downstream use reduces accidental variation. Shared input is then part of the governance benefit, not contamination of an experiment that was meant to be independent.

The input policy should make the phase and authority status clear. A draft model should not quietly become everyone's mandatory interpretation just because it is convenient to read.

## 7. Changes and readiness

If A13 advances to A14, the Behaviour Model may need reconciliation. Harness and Design may also be affected by architectural/environmental changes not fully represented by that model.

Once compatible new inputs are accepted, Implementation can reconcile its previous state. Whether every direct revision change immediately marks it stale, or whether effects flow through its immediate parents, is still an open invalidation policy.

The important requirement is an accountable acceptance basis, not a visually minimal graph. Omitting direct edges must not hide applicable obligations; adding all possible edges must not automatically imply wasteful full rebuilding.

## 8. Why this is a useful Saum example

The projects are heterogeneous: broad prose, behavioral model, design, executable harness and code. Each requires different local methods and different conformance evidence. Yet the surrounding questions are shared: what governs this state, what changed, what is the objective, how can it be inspected, what has been accepted, and what findings need to return upstream?

That makes it a concrete test of the proposed generic interface rather than a reason to put Allium-specific logic in Saum's core.

## Sources

[Architecture Index](../history/traust-reference/docs/architecture/00_architecture_index.md), [risk register](../history/traust-reference/docs/architecture/11_risks_and_debt.md), [Scenario 01](../history/traust-reference/docs/tests/01_happy_path.md), [state-transition reference](../history/traust-reference/docs/appendix/04_state_transition_spec.md); main-thread M08-M09. See [Behaviour Model topic](../concepts/11-behaviour-models-allium-and-documentation.md) for limits on tool-specific claims.
