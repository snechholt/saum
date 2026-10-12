# 13. Delivery options, non-goals and self-hosting

**Status:** candidate implementation boundaries and future possibilities. Nothing here is a committed roadmap.

## 1. Separate architectural importance from first-version scope

A concept can matter to the eventual model without needing a sophisticated first implementation. Provenance may be essential while its initial representation is simple files and commit references. Independent validation may be essential while initially relying on a separate prompt, tests and a human. A rich graph UI may be useful without being needed to test the model.

The original prototype already follows this principle: manually run experiments by command, manually invoke manager evaluation and synthesis, review the result, then use selected learning in a later attempt.

The full platform should not be required to discover whether the learning and acceptance semantics actually help.

## 2. Candidate initial reconciliation slice

The validation branch proposed an initial MVC with versioned project dependencies, a reconciliation objective separate from its plan, persistent negotiated increments, adaptive replanning, increment candidates, simple independent validation, internal checkpoints, final settlement and basic status visibility.

A concrete implementation could initially use Git, Markdown, available tests and explicit human decisions. The source does not select exact schemas, libraries or a database. "MVC" is preserved as the branch's label rather than assigned a new expansion here.

The key demonstration would be that a large objective can be progressively planned, executed, evaluated and accepted without Zonotron controlling the project's internal methodology.

## 3. A separate experimental proof slice

The main-thread prototype focused on a different but related slice: independent attempts -> factual evaluation -> vetted learnings -> human review -> fresh synthesis or guidance -> another attempt.

That tests whether comparison produces reusable knowledge and whether that knowledge improves new work. It need not implement the full persistent project graph first.

These are complementary candidate slices, not a silently unified minimum feature set. The team still needs to choose which next result would reduce the most uncertainty.

## 4. Conservative boundaries proposed for an initial version

The validation branch favored waiting for full upstream reconciliation before downstream consumption, avoiding a universal artifact parser, and using straightforward validators with human judgment rather than advanced trust calibration.

The main thread favored a loose project interface and reuse of existing execution methods. A custom adapter for every framework, complete PR management and a generic agent runtime should not be presumed necessary.

Rich multi-project speculative experiments, partial downstream settlement, automatic graph restructuring and elaborate policy learning are potential later capabilities. They should not be implied by the need for simple candidate comparison or a living plan.

## 5. Responsibilities that may remain central

The stable conceptual center includes identified authority and revisions, explicit objectives, reconciliation readiness, accountable candidate acceptance, preservation of evidence, and a route for upstream findings that does not silently mutate authority.

Planning negotiation and validation contracts are strong candidates for the interface, but their exact minimal form is still open. Experiment orchestration is less settled: its value is central to the motivation, while its placement may be optional or delegated.

This distinction should be made deliberately before implementation. Otherwise the original experiment platform and the later generic reconciliation layer can accidentally combine into an oversized first version.

## 6. Non-goals carried forward

Zonotron is not intended to be one unconstrained agent doing architecture, tests and implementation in one conversation. It is not a promise of perfect prompting, universal semantic validation or automatic upstream rewriting.

It should not delete mature accepted artifacts merely because a parent changed. It should not automatically choose one experimental winner or merge every useful fragment. It should not make every ordinary worker perform broad architectural criticism or turn every local review into a top-level project.

It should not assume all domains use code, worktrees, PRs or the same named roles. Nor should genericity prevent a useful software-focused initial implementation.

## 7. Candidate later capabilities

The original inventory and later discussions preserve several possibilities: editable execution/information graphs; multiple reviewers for high-impact changes; automatic review-intensity recommendations; model/resource choices per node; annotation-triggered analyses; selected code reference/transplant policies; delayed comparative analysis; full-project and subtree experimentation; partial downstream consumption; governance-graph restructuring; richer evidence indexing; and reliability-calibrated automatic settlement.

These belong in the design inventory even when excluded from a first version. Their dependencies, benefits and failure modes need separate evaluation.

## 8. Using this dossier to build Zonotron's own architecture

The user wants this material to become input to an Architecture-level project that Zonotron could hypothetically help build itself. The dossier should therefore be treated as a source corpus, not already accepted Architecture A0.

A possible first assignment would be to extract candidate requirements and invariants, identify incompatible proposals, define the smallest useful system boundary, and propose the decisions that need human settlement. The project could then produce a clearer architecture and scenarios without silently ratifying every brainstorming idea.

A derived Behaviour Model could pressure selected semantics such as request, candidate, checkpoint and settlement relationships, but Allium is not a prerequisite for that first architecture effort. A Test Harness could eventually verify the accepted cross-project protocol before full implementation.

This is a potential self-hosting path, not evidence that self-hosting already works. Early work can remain manual while the resulting artifacts become examples for the future system.

## 9. Avoiding circular authority in self-hosting

The conceptual safeguard already present in the design is that generation does not confer authority. A Zonotron-assisted candidate architecture still needs acceptance under an external human/policy basis. Its own output should not authorize the transition that makes it authoritative merely because the system produced it.

The exact bootstrap trust model is not settled. This paragraph is an editorial implication of the existing authority principle, recorded for future design rather than a new agreed mechanism.

## 10. What should be decided next

The highest-leverage unresolved areas remain the meaning of project authority, the minimum project contract, the negotiation and living-plan boundary, the validation/settlement model, and the status of contextual inputs.

Selecting one bounded prototype objective should follow those decisions far enough to make its results interpretable. The choice of runtime, storage and interface can then be evaluated against the required semantics rather than driving them accidentally.

## Sources

[S01](../history/sources/01-original-manager.md), sections 53-59; [S03](../history/sources/03-validation-scoped-settlement-and-adaptive-planning.md), sections 12-13; main-thread prototype discussion M02-M03 and dossier request M10 in [source map](../history/sources-and-coverage.md).
