# 12. Existing approaches, integration and competitive questions

**Status:** historical comparison context and hypotheses to investigate. This is not a current product survey, verified feature matrix or proof of uniqueness.

## 1. Why the comparisons matter

The user had mainly used plain Codex interactively: development tasks, review, clarification and iteration. They had not yet explored subagents, issue-driven autonomy or specification workflows in depth. Saum described a hoped-for end state rather than an informed replacement for every existing tool.

The earlier comparison therefore served two purposes: explain what progressively higher-level delegation might feel like, and identify which proposed capabilities could be borrowed rather than built.

The durable question remains: which parts of Saum disappear when existing tools are composed, and which concepts still have no satisfactory home?

## 2. Levels of delegation discussed

**Interactive coding and subagents:** the human still frames the engineering task, such as linking Z to Y on page X. Subagents can investigate patterns, APIs, tests or isolated changes. The added parallelism need not change the unit of authority; the architecture can remain context the human is responsible for maintaining.

**Issue-to-PR work:** the human delegates bounded outcomes and reviews proposed changes. The operational world is tasks, branches, PRs, tests and feedback. The human begins acting more like a lead, but cross-artifact conformance is not implied by the issue lifecycle alone.

**Project orchestration:** larger outcomes can be decomposed into multiple work items, coordinated and reviewed. The relevant question is whether the product only orchestrates execution or also gives durable semantics to authority and reconciliation across independently governed artifacts.

**Specification-oriented work:** requirements, design and task artifacts become part of the production method. This is closer to the desired control over intent, but its scope may still be one feature lifecycle rather than continuous evolution of a configurable project graph.

These distinctions describe the conceptual experiences discussed, not a claim that every product falls neatly into one exclusive category.

## 3. Named systems and the role each played in the discussion

| Name/category | What it was used to think about | What needs checking before relying on it |
| --- | --- | --- |
| Codex, Claude Code, Copilot CLI and local parallel sessions | Worker execution, context isolation, subagents and worktrees. | Actual GUI/CLI support and boundaries in the chosen version. |
| AO / Agent Orchestrator | Coordinating coding harnesses, tasks, branches and PR/CI follow-up. | Exact project identity, installation model and how much planning/state it truly manages. |
| GitHub coding agents, Devin, OpenHands, Cursor background agents, Jules | Delegation of bounded tasks into reviewable code changes. | Current task scope, persistence, isolation and review behavior. |
| Factory / Missions | Larger autonomous outcomes and coordinated parallel work. | Actual architecture/state semantics versus broad product positioning. |
| BMAD | Downward progression through product, architecture, planning, development and review roles. | Which parts are procedures/prompts and which are enforced workflow/state. |
| Kiro | Requirements/design/tasks and possible design-first input. | Whether high-level architecture is merely input context or independently governed state. |
| Spec Kit | Structured specification-to-implementation flow, consistency and decomposition. | Current scope, extension behavior and the earlier claims about "spec of specs." |
| Tessl / OpenSpec-related approaches | Specification fidelity, drift and the spec-as-source philosophy. | Which actual workflows preserve code versus regenerate it; avoid treating an ecosystem as one position. |
| LangGraph, CrewAI, AutoGen | Infrastructure for agent/workflow execution. | Which guarantees are supplied and which domain semantics remain application work. |
| Temporal, Prefect | Durable orchestration, waits and operational workflow concerns. | Suitability for the implementation; no runtime was chosen. |
| Dagster | Analogy of durable assets and dependencies versus jobs that update them. | The analogy does not establish semantic settlement support. |
| Allium | A constrained behavioral artifact and possible upstream validation surface. | Actual language/tool analysis capabilities and version-specific commands. |

The earlier assistant made more specific claims about several products, including current workflow engines, commands and feature breadth. This dossier deliberately does not promote those claims into established facts. The user asked to preserve the design discussions, not perform a new market verification exercise.

## 4. AO as a practical stepping stone

The user remembered an AO/Agent Orchestrator tool as an intermediate step toward autonomy. The comparison used it to separate a coding harness from a coordinator that manages workers, branches and review loops.

That remains a useful evaluation question even if the exact tool/version differs: how much repetitive operational work can be delegated while keeping the underlying coding agent constant?

If such tooling handles isolation, status, CI and PR follow-up well, Saum need not own that machinery just to obtain a candidate. It still needs its own meaningful objective, acceptance and provenance boundary if those are not supplied.

## 5. BMAD and roles

The user saw BMAD as partly overlapping because it progresses through levels of product and engineering work, but different in philosophy. Their desired system starts from strong specifications at chosen levels and lets the graph determine what is governed. A named role workflow is not the same thing.

The remark that BMAD felt aimed at less technical users was the user's impression, not an established demographic fact. Its useful role here is as a reference for downward progression and for asking how much value comes from carefully structured procedures versus role-persona prompts.

Saum can allow BMAD-like work inside a project without making its roles the universal ontology.

## 6. Kiro and broad architecture inputs

Earlier responses described Kiro as supporting requirements-first and design-first feature workflows, potentially accepting architecture documents. The user's question was whether that is the same kind of top-level product architecture as arc42 or merely a feature design.

The key distinction is not whether a tool can read arc42. It is whether it manages separately accepted revisions, dependency obligations, upstream findings and ongoing reconciliation across projects. A tool can be useful inside one node even if it accepts a large architecture document as input.

The actual degree of overlap should be verified rather than assumed from labels such as design-first.

## 7. Spec Kit and architectural decomposition

The conversation treated Spec Kit as the closest challenge to the proposal. Earlier responses said its "spec of specs" decomposition was a high-overhead fallback; the user contrasted that with Saum, where propagation through an existing project graph is central.

That external characterization remains unverified here. The conceptual distinction survives: decomposing one oversized feature is not identical to determining how a new authoritative state affects several already accepted artifacts, including what to preserve, reopen, inform or defer.

Ordinary feature-to-task decomposition may be widely supplied. Architecture-to-reconciliation-scope analysis across heterogeneous governed artifacts was the more specific capability under discussion. Calling all decomposition commodity obscures that difference.

## 8. Reconciliation versus strong spec-as-source

The early critique used spec-as-source/regeneration as a philosophical contrast. One possible approach treats implementation as cheaply reproducible from sufficiently precise specs. Saum instead values accumulated reviewed structure and normally reconciles it incrementally.

Neither position should be turned into a blanket claim about all tools under one brand. The architectural bet is that meaningful downstream knowledge will often remain expensive or unreliable to reconstruct from upstream prose alone.

If regeneration becomes adequate for a particular artifact, a project could use it internally. Saum's preservation principle is about accepted knowledge and governed change, not prohibiting every local regeneration technique.

## 9. Historical methodological ancestry

The discussion connected Saum with model-driven engineering, CASE, traceability, V-model thinking, risk-driven spikes, prototypes and explicit architectural decisions. The lesson was not that the idea is new or dead, but that its ancestors expose recurring costs: maintaining links, synchronizing representations, rigid process, and a second system describing the first.

Claims about historical empirical studies were not reverified here. The design concern stands independently: the graph must earn the cost of maintaining it.

The potentially changed economic assumption is that agent-produced throwaway materializations may be cheap enough to use as tests of upstream intent. That is a hypothesis to measure, not a demonstrated general advantage.

## 10. The strongest arguments against building a new system

The earlier critique identified several serious objections. Many teams may obtain sufficient value from existing specs, agents, Git and review with far less ceremony. Graph bookkeeping can dominate useful work. Accepted decisions can fossilize accidental architecture. Human gates can become bottlenecks. Independent experiments can duplicate discovery and add little benefit on tightly sequential tasks.

Semantic reviewers remain fallible, and deterministic policy can accumulate its own complexity. Better future agents could reduce the need for some intermediate artifacts. A composition of existing tools plus modest dependency metadata might achieve most of the same value.

These are not reasons to abandon the idea. They define the burden of proof: does explicit settlement, evidence retention and cross-project reconciliation provide materially better long-horizon coherence for the intended projects?

## 11. What not to claim as differentiation

Parallel agents, DAG execution, issue-to-PR automation, named roles and generic spec-driven development should not be treated as the unique thesis merely because the prototype uses them.

The more specific candidate differentiation is accepted state with provenance; a durable artifact graph distinct from execution; evidence flowing upward without automatic authority mutation; incremental reconciliation; and experimental pressure that can improve upstream intent even when no candidate wins.

No claim was established that no existing product implements those ideas, or that competitors could not add them. The historical research is a starting point for verification, not a market moat conclusion.

## 12. Practical exploration path retained

The suggested learning progression was to try local parallel/subagent work on familiar tasks, issue/project orchestration with the same underlying coding harness where possible, then one nontrivial specification-oriented workflow, with a role-based method as an optional contrasting experience.

The user accepted that direction. The purpose is to learn what is already solved before designing custom machinery. It is not a dependency on adopting a particular product or a claim that the user has completed those trials.

## Sources

Main-thread comparison and tool-experience episodes M01 and M04-M05 in [sources and coverage](../history/sources-and-coverage.md); [S01](../history/sources/01-original-manager.md), core design distinctions. This file records the implications of the discussion; it does not supply fresh external verification.
