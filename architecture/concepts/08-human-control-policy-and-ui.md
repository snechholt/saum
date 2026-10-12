# 08. Human control, autonomy policy and interaction

**Status:** working human/authority principles; UI and automation features are mostly tentative or exploratory.

## 1. Human ownership is selective, not uniformly hands-on

The user imagines defining product/system intent at a CEO/CTO-like level, then wearing more specific review hats for architecture, design, test approaches, harness implementation and code. This is not necessarily a fixed chain of agent personas.

The graph expresses which bodies of knowledge the user wants to govern independently. Internal local decisions may remain delegated. The same person can operate several acceptance boundaries without every decision becoming a separate project.

Zonotron should therefore support deliberate control resolution: explicit interaction models when their ownership matters, omitted layers when downstream discretion is acceptable.

## 2. Different kinds of human decision

Approving a reconciliation objective sets what should become true. Approving a plan says the proposed route and increments are acceptable to attempt. Accepting an internal checkpoint records justified progress. Global settlement determines what becomes an input others may rely on.

A human can also approve an upstream clarification, authorize scoped reopening, select an experiment policy, or revise a learning. These decisions have different consequences and should not collapse into one generic "approved" status.

The exact representation is open. A manual prototype may use prompts and Git history, but the conceptual distinctions should already be visible.

## 3. Approval of rules versus ordinary progress

The original inventory argued that changes to prompts, guidelines, specifications and public/test-only boundaries generally deserve review because they change the rules of the game. Ordinary implementation under accepted rules can have more autonomy.

Examples of candidate gated actions include adopting a new test-only interface, changing an upstream interpretation, widening experimental scope, reopening architecture, or doing major re-architecture. These were illustrative policy classes, not a settled universal permission table.

A project-specific authority may define a different procedure. The v5 manager review specifically warned against imposing one Traust-specific deferred-question process in every reusable workflow. Derive required procedure from the relevant authority.

## 4. Autonomy as explicit transition permission

The original direction was to widen deterministic auto-approval rather than make a manager agent sovereign. A user could permit ordinary continuation, bounded fan-out, or routine reconciliation under specified conditions while retaining gates for consequential changes.

An agent can recommend a transition. The application controls whether it is permitted, requires review, or exceeds the authorized scope. This preserves responsibility for resource use, information access and authority mutation.

The validation branch adds calibration by change class. Trust earned on routine dependency updates should not automatically authorize architectural changes. Occasional independent audits were proposed to check whether auto-acceptance remains justified.

## 5. A useful distinction between control plane and reasoning

The deterministic layer need not author all plans or inspect every artifact. It can delegate those judgments to bounded processes while enforcing which candidate, objective, input set and policy the resulting decision applies to.

Conversely, a sophisticated reviewer does not acquire authority merely because it has a persuasive explanation. A confident recommendation and a legal settlement transition are different objects.

This separation is part of the safety/coherence argument for Zonotron, not a guarantee that configured policies will be wise.

## 6. The original DAG UI

The original vision made the execution DAG an editable control surface, not merely a picture of hidden automation. A user could inspect nodes, input revisions, models, reasoning levels, budgets, output evidence and stop reasons. They could duplicate an experiment, add a reviewer, change a synthesis policy, or approve/reject a candidate.

Information flow would also be visible: no sibling edge means no sibling input; a reviewer receives all candidates; synthesis receives only permitted findings/references. This could make contamination and provenance easier to inspect.

The later project-encapsulation direction changes how much of that graph the top level needs to expose. A compact core may primarily show project state and reconciliation increments, with deeper execution detail available from the project or optional experiment service.

The rich DAG UI is retained as an attractive possibility, not an initial requirement or a settled product identity.

## 7. A cross-project state view

The working state view would distinguish accepted revisions, active reconciliation targets, plans and checkpoints, outstanding scope, blockers, potentially stale dependants and compatible input sets.

For example, Architecture may be A13, Design already reconciled, Harness still working, and Implementation still valid for its older basis but not ready to adopt the new one. The interface should not imply that latest-at-each-node equals a coherent global state.

The validation branch suggests showing a living plan: current increment, accepted internal checkpoints, remaining obligations and reasons for replanning. This gives the human useful oversight without exposing every internal implementation task.

## 8. Annotation-driven workflows

The original inventory proposed PR-like comments attached to commits, files, hunks, reviewer learnings, decisions and settlement proposals.

A comment such as "keep this abstraction, but place responsibility in foo rather than bar" could trigger several bounded operations: compare possible implementations, analyze repercussions, revise the learning, or hold a free-form discussion. Those operations should preserve their relationship to the original comment and evidence.

A future interface could render alternative code previews or impact analyses. These are exploratory interaction ideas. They do not require a custom interface before the semantic workflow can be tested manually.

## 9. Local approval need not become a new graph node

An implementation project can ask a human to review decomposition or a design fork within its methodology. That alone does not mean Design must become a separate project.

Promote the design when it has independent authority, consumers or a useful reconciliation lifecycle. This avoids confusing workflow visibility with artifact governance.

The user may still drill into local work. Encapsulation is not a ban on inspection or intervention; it is a boundary on what the overarching system must understand and control by default.

## 10. Monitor behavior and intervention

A live experiment monitor should enforce objective, accepted stop conditions rather than steer toward a preferred architecture. A human can choose to stop, widen, revise or restart an experiment, but the change should be explicit so the evidence remains interpretable.

Resource controls and stopping conditions were important original responsibilities. The v5 unchanged-ready redispatch problem shows that ordinary lifecycle plumbing can matter independently of prompting quality.

How emergency cancellation, retries, failures and partial output are represented remains open. The discussions identify them as necessary architectural questions, not completed designs.

## 11. Security and information boundaries

The original inventory explicitly listed a security/trust model: what agents may mutate, what the application enforces, experiment isolation and approval authority. It also made allowed information edges part of experimental provenance.

That leaves substantial work for the eventual architecture: credential handling, execution permissions, evidence visibility and isolation of real side effects. A worktree separates files but should not be assumed to isolate a shared Docker engine or external resource. The Traust review exposed precisely that kind of mismatch.

These are unresolved concerns, not a fully specified sandbox model. A future implementation must avoid equating a logical branch with complete environmental isolation.

## Sources

[S01](../history/sources/01-original-manager.md), sections 3-4, 9, 16-19, 30-34, 39 and 57; [S02](../history/sources/02-project-boundary-notes.md), human workshop and shared-engine concerns; [S03](../history/sources/03-validation-scoped-settlement-and-adaptive-planning.md), sections 3, 5, 7 and 12-13. See [evolution](../history/evolution.md) for the shift in UI/core scope.
