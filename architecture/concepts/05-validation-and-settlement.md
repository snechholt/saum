# 05. Validation, confidence and settlement

**Status:** working direction refined by an exploratory branch. The distinction between production, evaluation and acceptance is strong; concrete validation and checkpoint mechanics remain open.

## 1. What settlement means

Settlement means that a candidate has been accepted under a recorded policy using available evidence. It is not a declaration of mathematical correctness, nor proof that no future evidence can invalidate the decision.

A useful acceptance claim is relational:

> This candidate is accepted as satisfying the stated reconciliation objective against these upstream inputs, with these preserved constraints and this evidence.

That is more meaningful than "the artifact looks good in isolation." It also makes the policy and its limitations part of what later reviewers can understand.

The validation branch explicitly rejects the idea that a project report alone can settle a candidate. It equally rejects the need for Zonotron to implement every domain-specific validator itself.

## 2. Three responsibilities

The project produces candidate state. An independent validation process evaluates relevant conformance claims. Zonotron, under human or configured automatic approval policy, decides whether that evidence warrants acceptance.

A validator can be a deterministic test suite, specialized tool, independent agent, experiment, human reviewer or combination. The appropriate choice depends on the relationship and change, not merely the artifact extension.

"Independent" should not be mistaken for infallible. A second model can repeat an assumption. A test harness can share the implementation's blind spot. Human reviewers can overlook gaps. The architecture is a process for justified confidence and accountability, not universal semantic proof.

## 3. Edge-specific conformance

The proposed validation contract depends on what the edge means.

Architecture -> Behaviour Model concerns faithful representation of applicable requirements and disclosure of interpretation choices. Architecture -> Test Harness concerns coverage and whether the harness can actually demonstrate required behavior through legitimate boundaries. Design -> Implementation concerns accepted structure and interaction decisions. Test Harness -> Implementation concerns behavior demonstrated by externally supplied tests.

Passing one relation does not automatically establish another. A program can pass a weak harness while violating architecture. An Allium model can be internally coherent while embodying a mistaken interpretation of prose. A design can look structurally compliant while lacking an executable path for a required behavior.

How contracts are created, versioned, approved and maintained is open. They should not become a hidden second specification that nobody governs.

## 4. Completeness, soundness and proportionality

The branch separates three questions that had previously been bundled into "is this good enough?"

**Completeness:** are all relevant obligations accounted for? Individually tested scenarios might still lack a shared replaceable-SUT boundary or an integrated suite that can run them together. Missing cross-cutting requirements can survive a fully checked plan.

**Soundness:** does the evidence establish the claim? A fake registry might bypass the very behavior that was meant to be tested. A final persisted record does not necessarily prove that the record was committed before unfreeze began. A good/bad stand-in pair can share an invalid assumption. This was identified as the hardest semantic problem.

**Proportionality:** is the approach unnecessarily costly or complex? A full protocol implementation may be excessive where a small faithful fixture is sufficient. That may be a design/efficiency concern without being a conformance failure. It should not automatically invalidate an otherwise sound candidate.

The categories can interact, but their remedies differ. Missing scope requires accounting; weak evidence requires better validation; unnecessary complexity invites design or cost review.

## 5. Evaluate actual claims, not only metadata

The v5 review illustrates why worker readiness is not authority. Four ready claims were rejected by the supplied review. A different worker remained partial/blocked while its public CLI could report success prematurely.

Consequently, evaluation should distinguish worker-reported state, actual public behavior, asserted coverage and the strength of supporting evidence. A final branch snapshot and a historical first completion claim may also be different objects of evaluation.

The original manager wanted comparison at the first claimed frontier completion so later speculative work would not obscure that claim. The review refined this: public zero-exit or success behavior can be a completion claim even without a matching worker status. Candidate repairs should not retroactively change what an earlier claim established.

## 6. Trace the evidence path

For a consequential assertion, a reviewer may need to follow the originating action, product/runtime configuration, external observation, identity lookup, lifecycle availability and final predicate. Otherwise a test may appear independent while assuming a private implementation detail.

Examples from the review include fixed runtime suffixes, particular route-file conventions, one HTTP invocation mechanism, and command containers remaining available long enough to inspect. A legitimate independent implementation may vary those choices. Rejecting it for that variation is a false negative, not a strong test.

Conversely, a test may accept an unrelated digest or prohibited timing sequence. Negative probes should target meaningful counterexamples, and permitted-variation checks should challenge false rejection.

This is a review technique inferred from concrete experimental evidence, not a requirement that every routine change undergo full trajectory analysis.

## 7. Evidence strength and limits

Static inspection, a local unit probe, a synthetic counterexample, an integrated execution and an independent implementation exercise establish different things. A local probe can convincingly expose a defect without establishing full end-to-end success.

The archived v5 review did not independently rerun Docker integration because the review environment lacked access to the Docker pipe. It did run focused probes. That distinction belongs in the result rather than being hidden by an overall verdict.

Similar discipline should apply to any domain: an analysis should say what was observed, what was inferred, what was not tested, and what evidence would resolve the remaining claim. Insufficient evidence should not automatically be reported as proven nonconformance.

## 8. Internal checkpoints versus global settlement

The validation branch proposes a conservative default. An increment can be accepted internally while the overall project remains in reconciliation. For a harness required to cover S01-S20, acceptance of S01 does not automatically make the whole harness current against A13.

Later increments may reveal a foundational issue and require earlier accepted work to change. Letting Implementation depend immediately on every checkpoint could create expensive downstream rework.

The proposed initial model resembles holding accepted progress on a reconciliation branch until the whole objective is ready for global settlement. This is an analogy, not a mandated Git branching scheme.

A checkpoint is therefore meaningful progress with recorded evidence, not an irrevocable promise to downstream consumers. Reopening should be explicit; the exact invalidation rules remain open.

## 9. Incremental and final integrated validation

For a new increment, validation asks both whether the new claims hold and whether earlier accepted claims still hold. Regression checks may handle much of the latter, but their adequacy remains part of the evidence question.

After all planned increments are accepted, final reconciliation validation asks whether the integrated artifact satisfies the complete objective. A collection of individually accepted fragments is not automatically a coherent whole.

This final gate is particularly important when the plan has evolved. The objective remains the stable reference; the final reviewer should not merely audit whether the latest checklist has ticks.

Under the conservative branch proposal, downstream projects normally wait for this global settlement. Partial consumption is a future option requiring scoped authority, explicit provisionality and a rework policy. It is not established as forbidden forever or required initially.

## 10. Producer-proposed evidence is valuable but insufficient alone

A project can explain how it decomposed a problem, which mechanisms it used and which tests it believes demonstrate success. That can guide efficient inspection.

But the acceptance basis also needs upstream-derived obligations independent of the chosen implementation. Otherwise the producer can accidentally define away its own blind spots. A changed evidence strategy can be legitimate without permitting a weaker objective.

A validator should compare the two: what must be established and what the offered evidence actually establishes. Uncovered claims can lead to further work, alternative evidence, clarification or human review.

## 11. Approval intensity by risk and change class

The branch proposes calibrating trust per class of change, not simply labeling an entire project trusted. Repeated success on dependency upgrades does not establish competence for changes to security or state-transition logic.

Possible policies range from tests plus automatic acceptance for routine local updates, to independent semantic review for behavioral changes, to human settlement and experiments for foundational changes. These are examples of policy classes, not a settled classification system.

Occasional independent audits of auto-accepted work were suggested as a way to test reliability. Exact sampling, confidence standards and metrics remain open. Auto-approval is a controlled permission for known transition types, not a grant of unrestricted judgment to a manager agent.

## 12. Local human review and global acceptance

Human review can occur inside a project: approve a proposed decomposition, compare design alternatives or review a code change. That does not automatically make each reviewed item an independently governed project.

Zonotron-level settlement has a different consequence: it changes what other projects may treat as accepted authority. Planning approval, internal checkpoint acceptance and global release should be distinguishable even if the same human performs all three in a prototype.

This distinction preserves the user's desire to wear different review hats without forcing an artificial organizational hierarchy into the graph.

## 13. Remaining architectural difficulty

The hardest unresolved issue is sound semantic validation of arbitrary heterogeneous artifacts. More structure, independent review and provenance improve the process but do not eliminate the problem.

Open questions include who defines edge obligations, what an unavailable validator implies, how a later structural change invalidates earlier evidence, how accepted claims are represented without building a burdensome requirements database, and what evidence justifies reducing human gates.

## Sources

[S03](../history/sources/03-validation-scoped-settlement-and-adaptive-planning.md), sections 1-3 and 7-14; [S01](../history/sources/01-original-manager.md), sections 9, 18-19, 35, 38-39 and 43; [v5 evaluation](../history/evidence/review/EVALUATION.md) and [manager review](../history/evidence/review/MANAGER_REVIEW.md).
