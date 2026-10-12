# Saum: idea inventory and evolving conceptual design

**Dossier date:** October 11, 2026  
**Status:** working design dossier; not a settled architecture, implementation specification, or committed roadmap.  
**Purpose:** preserve the ideas, reasoning, examples, disagreements, experimental lessons, and unresolved questions developed in the main discussion and two attached branch summaries. This material can later become input to a Saum Architecture project.

Saum is the provisional name of the proposed system. **Traust is the separate deployer/orchestrator whose test-harness experiments helped motivate it.** The two must not be conflated.

## The current center of gravity

The working idea is a system for managing a graph of independently governed project artifacts. Accepted upstream state constrains downstream work. Projects use their own methodologies to reconcile their existing state against those constraints. Saum coordinates the cross-project process, retains evidence, and controls which candidates become accepted state under human or configured policy authority.

Independent experimentation, comparative evaluation, and synthesis remain important, but their position has evolved: from defining the whole system to potentially being reusable ways of advancing or testing a reconciliation. Exactly how much of that machinery belongs in the core remains open.

The persistent object is accepted state and its provenance, not an agent persona or conversation. The system should preserve useful accumulated work, support challenges to it, and distinguish observations from decisions about what should become authoritative.

## How to read this collection

The topical documents are the working synthesis. The examples preserve concrete pressures that led to the abstractions. The history records shifts in thinking and the sources preserve original wording, including stronger claims and proposals that should not automatically be inherited.

These files are **documentation topics, not proposed Saum project nodes**. The collection does not settle a storage model, API, workflow engine, state machine, permission system, or integration strategy.

| File | What it preserves |
| --- | --- |
| [Glossary](glossary.md) | Terminology, overloaded words, and notation. |
| [01. Purpose and system boundary](concepts/01-purpose-and-system-boundary.md) | Motivation, generic scope, current boundary, and competing interpretations. |
| [02. Projects, graphs, and authority](concepts/02-projects-graphs-and-authority.md) | Nodes, revisions, dependency/context/evidence relationships, consistent state. |
| [03. Reconciliation and planning](concepts/03-reconciliation-and-planning.md) | Change propagation, global/local planning, living plans and planning horizons. |
| [04. Project interface and negotiation](concepts/04-project-interface-and-negotiation.md) | Three contracts, loose coupling, intake questions, blockers and proposals. |
| [05. Validation and settlement](concepts/05-validation-and-settlement.md) | Independent evaluation, evidence soundness, checkpoints, final acceptance. |
| [06. Experiments and synthesis](concepts/06-experiments-and-synthesis.md) | Independence, modes, scopes, monitoring, selective reuse and unresolved ownership. |
| [07. Provenance, findings, and guidance](concepts/07-provenance-findings-and-guidance.md) | Evidence retention, upward learning, attention filtering and temporary worker guidance. |
| [08. Human control, policy, and UI](concepts/08-human-control-policy-and-ui.md) | Approval, bounded autonomy, local/global review, graph interaction and annotations. |
| [09. Lifecycle, cost, and empirical evaluation](concepts/09-lifecycle-cost-and-empirical-evaluation.md) | Bootstrap, foundation, buildout, maintenance, maturity and proof obligations. |
| [10. Legacy reconstruction](concepts/10-legacy-project-reconstruction.md) | Cosmos, characterization, evidence flowing upward, and the normative pivot. |
| [11. Behaviour models and Allium](concepts/11-behaviour-models-allium-and-documentation.md) | arc42, constrained materialization, interpretation risk and contextual access. |
| [12. Existing approaches and integration](concepts/12-existing-approaches-and-integration.md) | What earlier comparisons contributed, without treating product claims as verified facts. |
| [13. Delivery scope and self-hosting](concepts/13-delivery-scope-and-self-hosting.md) | Prototype, candidate MVC/MVP, later possibilities and architecture-project handoff. |
| [Open questions](open-questions.md) | Indexed unresolved issues, including tensions between the discussions. |
| [Traust Scenario 01](examples/01-traust-scenario-01.md) | The motivating evidence, review limitations, and what remains unproved. |
| [Baz change](examples/02-baz-reconciliation.md) | A large change through a realistic small number of projects. |
| [Traust documentation graph](examples/03-traust-artifact-graph.md) | Architecture, Behaviour Model, Harness, Design and Implementation. |
| [Cosmos onboarding](examples/04-cosmos-reconstruction.md) | Walking from legacy code to an accepted authority graph. |
| [School and other governance shapes](examples/05-school-and-other-governance-shapes.md) | Domain independence, TODO foundation, optional interaction/interface layers. |
| [Evolution of the idea](history/evolution.md) | Earlier positions, refinements and unresolved forks. |
| [Sources and coverage](history/sources-and-coverage.md) | Provenance, original-section coverage, branch integration and source limits. |

For a first reading: 01 -> 02 -> 04 -> 03 -> 05, then the examples. For experimentation: 06 -> 07 -> 09 -> Scenario 01. For a later architecture project: include the open questions and evolution rather than supplying only the apparently clean model.

## Status vocabulary

**Working direction** means a repeatedly supported or currently preferred idea, not a ratified rule. **Tentative** means a promising proposal whose alternatives remain live. **Open** means the discussions did not settle the question. **Exploratory** means an interesting possibility rather than a required capability. **Earlier approach** means a historical position that has been refined, displaced, or deprioritized, not necessarily disproved. **Reported observation** refers to what the supplied experiment review records; it is not a new reproduction of the experiment.

A section can have more than one status. Absence of a warning on every sentence does not promote this dossier into authority. The attached branch summaries are brainstorming inputs even where they use words such as "must," "conclusion," or "solution."

Implementation priority is a separate axis. A concept may be architecturally central yet deferred in an initial implementation. Conversely, a convenient prototype mechanism is not necessarily a lasting architecture choice.

## Source and editorial discipline

The original manager inventory, the user's project-boundary notes, and both branch summaries are retained unchanged under [history/sources](history/sources/README.md). Selected experiment reviews/prompts and Traust reference documents are also included. Their SHA-256 fingerprints appear in the source manifest.

The topical documents integrate those sources with the main conversation. They preserve alternatives rather than silently resolving them. Source references use portable Markdown links and original section names/numbers. Main-thread references use thematic IDs defined in the source map; they identify discussion episodes, not invented transcript exports.

No new market research, Allium capability verification, or worker experiment was performed for this dossier. Statements previously made about current third-party products are retained as historical comparison hypotheses or questions to verify. In particular, do not build implementation requirements around earlier claims about Allium commands or Spec Kit capabilities without checking the actual tool/version.

The review archive contains the Traust documentation snapshot used by the v5 review. Selected references here come from that archive; they are not represented as a fresh inspection of an independently mounted copy of the later documentation upload.

## Using this in another context

Treat this collection as evidence and design input. Ask a subsequent reviewer to distinguish the current direction, alternative designs, empirical observations, and missing decisions. Do not ask it to implement every idea described here as though they were one consistent requirement set.

When adding another branch discussion, preserve its source, integrate the relevant detail into the topic files, record any change in interpretation in the history, and update the unresolved questions. Keep one home for each detailed explanation and use cross-links elsewhere.

There is intentionally no combined `SAUM_FULL.md`. The individual Markdown files are the maintained form.
