# 01. Purpose and system boundary

**Status:** working direction, with an unresolved boundary around experiment orchestration.  
**Related:** [graphs](02-projects-graphs-and-authority.md), [interface](04-project-interface-and-negotiation.md), [history](../history/evolution.md).

## 1. The problem being addressed

The original problem was not how to ask a coding agent to build an application from a vague idea. The intended starting point was comparatively mature, structured intent: architecture, scenarios, contracts, and a backlog that a human was willing to own. Agents would materialize that intent into increasingly concrete artifacts.

The Traust experiments then exposed a deeper problem. A worker can make substantial progress, demonstrate something plausible, and still solve the wrong problem or use an invalid evidentiary boundary. Other workers can fail to finish yet produce valuable mechanisms or expose weaknesses in the instructions. Reviewing just the apparent winner loses much of the value.

The motivating need became: retain the useful accepted state and the evidence that produced it; distinguish candidate work from authority; and feed discoveries into either the next local attempt or an upstream specification review without letting a worker silently redefine the task.

This is why the permanent object is not a prompt, role, chat, issue, or winning branch. It is the evolving set of accepted artifacts, decisions and compatibility claims, with provenance.

## 2. What the human is trying to retain

The user wants strong specifications at whichever levels matter, rather than a system that automatically chooses all those levels on their behalf. A project graph can make some decisions explicit while deliberately leaving others to downstream workers.

For example, a human may want to govern architecture and externally observable behavior but leave component decomposition internal to implementation. Another human may insist on separately approving component boundaries and interaction diagrams. Both should fit the model.

The original formulation was:

> Agent discretion is the set of decisions intentionally left unresolved by settled upstream artifacts.

That is a useful intent, but it is not an assertion that every real omission is intentional. Experiments may reveal accidental gaps. The system needs a way to distinguish permitted discretion from a missing decision that affects other projects.

This makes Saum potentially useful even when a human does not want to delegate all engineering. The graph selects the resolution at which the human retains control.

## 3. From software-development manager to generic reconciliation layer

The user's later proposal generalized the model. A project need not be a software repository or a named engineering role. A school strategy can constrain an annual plan, which constrains grade-level plans. A business strategy can constrain product and operating plans. Each has work products that can be revised, interpreted, reviewed and accepted.

The role of an agent follows from the work it is asked to perform against those artifacts. It does not need a permanent persona called Architect, CTO or Principal. Named roles may still be useful inside a methodology, but authority comes from accepted artifacts and configured decision rights, not the persona.

This is an architectural generalization, not a decision to launch a universal organizational product. Software is the immediate test domain. The non-software examples pressure the abstraction and reveal when terms such as repository, pull request or coding agent have leaked into its core.

## 4. The proposed split

The working split is between cross-project governance and project-local method.

Saum coordinates which accepted upstream revisions a project must reconcile against, the relevant objective and scope, cross-project ordering, negotiation of viable work, and acceptance into the authoritative graph. It preserves references to candidate state, evidence and decisions. It can route findings to the appropriate upstream owner.

A project chooses how to do its work: plain coding sessions, a Ralph loop, a specification workflow, issue-to-PR automation, domain tools, experiments, human workshops, or a combination. It understands its current internals and can propose appropriate increments and evidence.

Independent validation evaluates conformance. Saum owns the acceptance process but need not implement every validator or possess one context containing all project internals. This is an important refinement from the validation branch, not a complete validation architecture.

The immediate mental model remains deliberately loose: one standalone repository is instructed to reconcile against identified state in another repository. A rich protocol could initially be expressed through instructions, files and returned references. "Adapter" should not silently imply a custom integration with every internal framework.

## 5. Encapsulation is not opacity

The proposed boundary is not "send request, trust result." A project must expose sufficient candidate state and evidence for meaningful evaluation. Its internal methodology may be private to its operation, but its externally relevant results cannot be beyond inspection.

The Scenario 01 concern makes this especially important. A project might dismiss an unsuccessful registry approach as a bad local idea and omit it from a summary. A cross-project reviewer might discover that the approach followed a plausible reading of an ambiguous upstream contract. A mandatory summary alone is therefore an information-loss boundary.

The stronger direction is observable execution and retained provenance, plus optional project observations and explicit blockers. Saum or a validator can inspect more deeply when warranted. That does not mean every transcript is always loaded into one giant manager context.

## 6. Deterministic control and semantic work are different

The original idea called for a deterministic orchestrator rather than a sovereign super-agent. That distinction survives the generalization, but its scope needs care.

The application can control which transitions are legal, which inputs a run receives, who can settle a candidate, whether a scope change is permitted, and whether an experiment may see its siblings. An agent can propose a plan, assess impact, or recommend settlement inside those boundaries.

Deterministic control does not make an agent's semantic judgment deterministic or correct. It makes the judgment's consequences explicit and reviewable. Similarly, delegated validation does not remove the need to decide which evidence is adequate.

The early model assigned detailed DAG nodes, worktrees, models and budgets directly to the orchestrator. The later model questions whether most of that should be delegated. Both are retained; the current common ground is bounded authority and inspectable information flow, not a settled runtime implementation.

## 7. What did not disappear when the boundary narrowed

Several original ideas remain valuable even if the core becomes smaller: independent candidate search; comparison of failures and partial successes; source-linked learnings; human review of interpretation; incremental reconciliation rather than regeneration; and scoped reopening instead of wiping mature state.

The unresolved question is packaging. Those capabilities could be core operations, reusable services available to projects, optional orchestrator workflows, or project-local methods. The user's project-level fan-out proposal shows that experimentation can sometimes be expressed generically without knowing the internal method.

The simplification is therefore not "abandon the original idea." It is "avoid coupling the knowledge and authority model to one way of doing the work."

## 8. Failure modes to keep visible

The system could become a task tracker with dependency arrows but no meaningful conformance semantics. It could become a second, expensive description of the actual system whose links are harder to maintain than the work itself. Or it could recreate the mega-agent by requiring the top level to interpret every action from every project.

There is also a governance risk: accepted early architecture may gain excessive inertia. Preservation is useful only if evidence can challenge it. Human settlement can become the new throughput bottleneck, while permissive auto-approval can conceal semantic failures.

Finally, independent exploration is not automatically economical. Its value must be demonstrated in accepted progress, avoided rework, or better upstream decisions, not merely interesting analysis.

## 9. Working framing

A useful current framing is:

> Saum coordinates the evolution of independently governed artifacts. It delegates production methods, preserves evidence and accepted state, and controls reconciliation and settlement across their authority relationships.

This does not settle whether a project graph must be a strict DAG in every respect, what constitutes a revision, or how much planning and experimentation the core performs. Those are architecture questions for the next stage.

## Sources

[S01](../history/sources/01-original-manager.md), sections 1-5, 23-27, 53-59; [S02](../history/sources/02-project-boundary-notes.md), especially the concerns about evaluator responsibility, loose tickets and realistic project composition; [S03](../history/sources/03-validation-scoped-settlement-and-adaptive-planning.md), section 1. Main-thread episodes M01, M04 and M05 in [sources and coverage](../history/sources-and-coverage.md).
