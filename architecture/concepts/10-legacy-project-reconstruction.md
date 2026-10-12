# 10. Legacy-project reconstruction and the authority pivot

**Status:** exploratory branch proposal with a coherent working direction; not a ratified onboarding protocol.

## 1. The problem

Normal Saum examples start with explicit upstream intent: Architecture constrains Design and Test Harness, which constrain Implementation. But an existing messy application may have only code. The source example calls that application Cosmos.

Much of Cosmos's actual system knowledge is implicit in implementation. To onboard it, the natural work direction is upward: inspect code, recover design and behavior, and reconstruct architecture.

This appears to reverse authority. If Design depended normatively on Implementation, future accidental code changes might redefine the design. Bugs, debt and historical shortcuts could become mandatory simply because they existed first.

The branch's proposed solution is to distinguish evidence flow from normative dependency rather than reversing the authority graph.

## 2. Evidence can flow upward without governing upward

An implementation snapshot can be the source of a reconstructed Design, characterization tests and a candidate Behaviour Model. Those can inform candidate Architecture alongside human explanation.

```text
Implementation I0
  -> evidence for as-built Design
  -> evidence for characterization tests / behavior
                + human intent and constraints
                -> candidate Architecture
```

The arrows here mean "inferred using this evidence," not "required to conform forever to this source." A reconstructed document begins as a candidate account, not accepted future authority.

This is the same asymmetry already valued in experiments: downstream work can teach upstream projects without silently rewriting their normative rules.

## 3. Characterization versus acceptance

A characterization harness asks what Cosmos demonstrably does now. An acceptance harness asks what Cosmos is required to do after the relevant decisions have been settled.

If `DELETE /user` currently returns 500, a characterization test can accurately record that behavior. It does not establish that returning 500 is a requirement. Human review might classify it as a bug, a temporarily tolerated behavior, or something needing further investigation.

The same applies to structure. An as-built design can observe that AuthService directly manipulates billing tables. A settled design may instead state that authentication must not own billing persistence. The discrepancy is valuable onboarding output, not evidence that reconstruction failed.

The branch suggested conceptual statuses such as observed/reconstructed, candidate and settled. Those are useful distinctions, not a finalized universal lifecycle enumeration.

## 4. The six-stage proposal

**Import:** identify the existing implementation snapshot I0.

**Observe:** recover components, actual dependencies, interfaces, schemas, data flows, runtime behavior, deployment structure and constraints that can be supported by evidence.

**Reconstruct:** produce candidate design, tests/behavior and architecture, with uncertainty, contradictions and potentially accidental behavior marked.

**Human settlement:** decide which observations represent the system the organization intends to continue owning, which are debt or bugs, and which require more knowledge.

**Pivot:** the accepted higher-level artifacts become normative inputs for future work.

**Reconcile:** update the existing implementation against the newly accepted authority graph while preserving useful accumulated knowledge.

These stages describe semantics, not necessarily a rigid one-pass workflow. The branch did not settle how much iteration happens between observation and human review.

## 5. No actual reversal at the pivot

Before settlement, Implementation I0 is evidence for candidate higher-level artifacts. After settlement, Architecture/Design/Test constrain future implementation. The original evidence relationships can remain in provenance without becoming normative dependencies.

It is legitimate for I0 to require reconciliation immediately after the pivot. The human may have accepted an intended design that deliberately rejects some existing coupling or behavior.

This avoids two bad outcomes: declaring all legacy behavior required, or throwing away the implementation's accumulated knowledge and pretending the project is greenfield.

## 6. Code alone cannot establish all intended architecture

Implementation can reveal actual structure, interfaces, behavior and constraints. It generally cannot reliably establish business goals, rejected alternatives, quality priorities, organizational constraints, historical rationale, or which accidental behaviors should be preserved.

The branch therefore combines code-derived evidence and reconstructed design/behavior with human purpose, known bugs and debt, preservation requirements, and business/operational constraints.

This is not a claim that every reconstructed fact needs a lengthy human workshop. It means the transition from observed reality to intended authority requires a legitimate source of intent rather than code alone.

## 7. Proposing a governance graph

Saum could inspect a legacy project and propose alternative artifact structures. A small service might need only Architecture -> Implementation. A more complex Cosmos could justify Design and Test Harness, or a Behaviour Model, security model or deployment model.

The criterion is not "how many documents can be generated?" It is "which bodies of knowledge should become independently settled authority?"

Graph proposals remain subject to review. The branch does not establish that all legacy projects need every possible layer or that structure should be inferred automatically from repository folders.

## 8. Relationship to preservation and reconciliation

Legacy onboarding is a particularly strong example of why downstream artifacts are not disposable. Existing behavior may encode undocumented compatibility requirements and hard-won practical knowledge, even where some choices are undesirable.

The intended outcome is to make enough of that knowledge explicit to govern future change. Reconstruction should preserve uncertainty and the evidence for observations, not clean up the story so thoroughly that historical reality is lost.

Once the authority graph is accepted, a later A0 -> A1 change can use normal reconciliation semantics. The unusual part is establishing the first authority basis, not permanently running the graph backward.

## 9. Unresolved questions

How much reconstruction is enough for a first useful authority graph? How should confidence in inferred behavior be recorded? What happens when observations conflict or code paths cannot be exercised? How are known bugs distinguished from temporarily accepted compatibility behavior? Can only part of the legacy system pivot while the rest remains characterized?

The source branch does not settle those mechanics. They should be treated as next design questions rather than silently filled with a universal reverse-engineering method.

A concrete walkthrough is in [Cosmos onboarding](../examples/04-cosmos-reconstruction.md).

## Sources

[S04: original legacy reconstruction branch](../history/sources/04-legacy-project-reconstruction.md), all sections. Related original principles: [S01](../history/sources/01-original-manager.md), sections 23-27. The branch is preserved verbatim; this topic integrates its proposal without promoting it to a final rule.
