# Sources, discussion map and coverage review

**Status:** editorial provenance and coverage map for this dossier. It identifies where material was integrated; it does not certify that brainstorming has become a consistent or complete architecture.

## 1. Source inventory

**S01 - Original manager inventory.** [Preserved source](sources/01-original-manager.md). A 59-section historical handoff describing the initial deterministic manager, both graphs, experiments, review, settlement, provenance and open implementation questions. Its original status explicitly says it is not settled architecture.

**S02 - User project-boundary notes.** [Preserved source](sources/02-project-boundary-notes.md). Direct user concerns about evidence loss, opaque projects, methodology-specific intake, negotiation, realistic graph shape and generic experiment fan-out.

**S03 - Validation/scoped-settlement/adaptive-planning branch.** [Preserved source](sources/03-validation-scoped-settlement-and-adaptive-planning.md). The user supplied it as exploratory, even where it uses strong conclusions. It introduces the three contracts, delegated validation, stable objectives versus living plans, checkpoints and conservative downstream release.

**S04 - Legacy-reconstruction branch.** [Preserved source](sources/04-legacy-project-reconstruction.md). The user supplied it as exploratory. It distinguishes observation from intent, evidence from normative edges, and characterization from acceptance.

**E01 - V5 review archive selections.** [Evidence index](evidence/README.md). These are actual archived review reports and proposed prompts, not new experiment results. They include explicit evidence limitations.

**T01 - Traust documentation references.** [Reference index](traust-reference/README.md). These come from the documentation tree inside E01. They ground the examples and internal ownership map, not a claim of newly inspected current product state.

Byte-preservation details are in [the manifest](source-manifest.md).

## 2. Main-thread thematic source index

These IDs identify visible episodes in the main conversation. They are not IDs of independently exported transcript files, and the paragraphs below are not quotations.

| ID | Discussion episode | Integrated detail |
| --- | --- | --- |
| M01 | Initial critique and comparison with existing systems | Settled state as candidate differentiation; alternatives; historic ancestry; economic and governance objections; experimental proof requirements. |
| M02 | Reviewing v5 and asking whether it proves the idea | Worker claims versus accepted truth; partial evidence; outer frontier; manager reliability limits; blind evaluation and guidance-cohort proposal. |
| M03 | Exploration cost, first-milestone breadth and temporary Ralph input | Discussion versus experiment; avoiding overdecomposition; deliberately open-ended foundation pressure; stable protocol plus scoped opinions. |
| M04 | The user's experience ladder from plain Codex to subagents, issue/PR systems and spec workflows | Delegated unit of work; human review hats; AO, BMAD and spec methods as references; what lower-level capabilities might be reused. |
| M05 | Generalization to school strategy and a layer above methods | Authority in artifacts rather than personas; global graph versus local method; broad change impact versus ordinary task decomposition; loose repository model. |
| M06 | Project-boundary challenge and the pasted notes | Observable not opaque projects; evidence retention; local/global knowledge; iterative intake; project-level alternatives; realistic single Implementation project; internal approval versus independent node. |
| M07 | AFK architectural questions | Minimum contract, negotiation questions and termination; prioritization; settlement meaning; governance scope; genericity and explicit non-goals. |
| M08 | Allium discovery and possible roles | Top-level behavior spec versus constrained downstream materialization; cheaper pressure surfaces; faithfulness and independent interpretation; tool assumptions requiring verification. |
| M09 | Applying the idea to the actual Traust arc42 documentation | Internal documentation layers, R-05 drift risk, Harness-before-Implementation, Behaviour Model input, access to the why, dependency/context distinction and unresolved edge policy. |
| M10 | Request for an exhaustive Markdown design dossier and supplied branches | Preserve current directions, alternatives, reasoning, examples and open questions; archive originals; individual files, no combined export; later Architecture-project use. |

The supplied Project context also contains references to related discussions, but this dossier does not invent their missing detailed contents. The two supplied branch summaries are the basis for their respective integrations.

## 3. Original manager inventory coverage

Every numbered section of S01 is retained in the original and mapped below to one or more places in the synthesis. A mapping means the idea is carried forward, possibly as historical, tentative or open; it does not mean adoption.

| S01 section | Main synthesis location |
| --- | --- |
| 1. Core premise | [Topic 01](../concepts/01-purpose-and-system-boundary.md), [Other examples](../examples/05-school-and-other-governance-shapes.md) |
| 2. The permanent object is the settled state, not an agent session | [Topic 02](../concepts/02-projects-graphs-and-authority.md), [Topic 07](../concepts/07-provenance-findings-and-guidance.md) |
| 3. Two different graphs | [Topic 02](../concepts/02-projects-graphs-and-authority.md), [Topic 08](../concepts/08-human-control-policy-and-ui.md) |
| 4. Deterministic orchestrator, not a super-agent | [Topic 01](../concepts/01-purpose-and-system-boundary.md), [Topic 08](../concepts/08-human-control-policy-and-ui.md) |
| 5. The worker primitive | [Topic 01](../concepts/01-purpose-and-system-boundary.md), [Topic 04](../concepts/04-project-interface-and-negotiation.md) |
| 6. Initial analysis: determine what is enough to settle a baseline | [Topic 09](../concepts/09-lifecycle-cost-and-empirical-evaluation.md), [Other examples](../examples/05-school-and-other-governance-shapes.md) |
| 7. Backlog/frontier-driven settlement | [Topic 03](../concepts/03-reconciliation-and-planning.md), [Topic 06](../concepts/06-experiments-and-synthesis.md) |
| 8. Experiment independence | [Topic 06](../concepts/06-experiments-and-synthesis.md) |
| 9. Manager/monitor must not steer live experiments | [Topic 06](../concepts/06-experiments-and-synthesis.md), [Topic 08](../concepts/08-human-control-policy-and-ui.md) |
| 10. Explicit situations the manager/reviewer should look for | [Topic 07](../concepts/07-provenance-findings-and-guidance.md) |
| 11. Two very different experiment purposes | [Topic 06](../concepts/06-experiments-and-synthesis.md) |
| 12. Specification development can use throwaway derivations as validation | [Topic 06](../concepts/06-experiments-and-synthesis.md), [Topic 11](../concepts/11-behaviour-models-allium-and-documentation.md) |
| 13. Experiment outputs: knowledge first, code optionally | [Topic 06](../concepts/06-experiments-and-synthesis.md), [Topic 07](../concepts/07-provenance-findings-and-guidance.md) |
| 14. But do not collapse experiment evidence into prose and discard the source | [Topic 07](../concepts/07-provenance-findings-and-guidance.md) |
| 15. Selectable policies for using experiment code | [Topic 06](../concepts/06-experiments-and-synthesis.md) |
| 16. Human annotations and PR-like review as first-class provenance | [Topic 08](../concepts/08-human-control-policy-and-ui.md), [Topic 07](../concepts/07-provenance-findings-and-guidance.md) |
| 17. Reviewer learnings are also versioned/reviewable | [Topic 07](../concepts/07-provenance-findings-and-guidance.md) |
| 18. Human approval boundaries | [Topic 08](../concepts/08-human-control-policy-and-ui.md), [Topic 05](../concepts/05-validation-and-settlement.md) |
| 19. Autonomy should mean widening deterministic auto-approval, not making a super-agent sovereign | [Topic 08](../concepts/08-human-control-policy-and-ui.md), [Topic 05](../concepts/05-validation-and-settlement.md) |
| 20. Lifecycle stages: useful defaults, not a rigid global state machine | [Topic 09](../concepts/09-lifecycle-cost-and-empirical-evaluation.md) |
| 21. Maintenance can still launch heavy local exploration | [Topic 03](../concepts/03-reconciliation-and-planning.md), [Topic 09](../concepts/09-lifecycle-cost-and-empirical-evaluation.md) |
| 22. Scoped reopening instead of wiping settled state | [Topic 03](../concepts/03-reconciliation-and-planning.md), [Topic 06](../concepts/06-experiments-and-synthesis.md) |
| 23. Authority flows downward; evidence can flow upward | [Topic 02](../concepts/02-projects-graphs-and-authority.md), [Topic 07](../concepts/07-provenance-findings-and-guidance.md), [Topic 10](../concepts/10-legacy-project-reconstruction.md) |
| 24. Reconciliation, not regeneration | [Topic 02](../concepts/02-projects-graphs-and-authority.md), [Topic 03](../concepts/03-reconciliation-and-planning.md) |
| 25. “Consistent cut” view of system state | [Topic 02](../concepts/02-projects-graphs-and-authority.md) |
| 26. Cross-project dependency messages, not one global agent context | [Topic 01](../concepts/01-purpose-and-system-boundary.md), [Topic 04](../concepts/04-project-interface-and-negotiation.md) |
| 27. Organic buildup is a feature | [Topic 03](../concepts/03-reconciliation-and-planning.md), [Topic 09](../concepts/09-lifecycle-cost-and-empirical-evaluation.md) |
| 28. Avoid fossilizing every model mistake into prompts | [Topic 07](../concepts/07-provenance-findings-and-guidance.md) |
| 29. Search vs exploitation | [Topic 06](../concepts/06-experiments-and-synthesis.md), [Topic 09](../concepts/09-lifecycle-cost-and-empirical-evaluation.md) |
| 30. Potential manager-triggered workflows | [Topic 06](../concepts/06-experiments-and-synthesis.md), [Topic 08](../concepts/08-human-control-policy-and-ui.md) |
| 31. The DAG UI as the central control surface | [Topic 08](../concepts/08-human-control-policy-and-ui.md) |
| 32. Information flow should itself be visible | [Topic 02](../concepts/02-projects-graphs-and-authority.md), [Topic 06](../concepts/06-experiments-and-synthesis.md), [Topic 08](../concepts/08-human-control-policy-and-ui.md) |
| 33. Model choice is an explicit workflow parameter | [Topic 06](../concepts/06-experiments-and-synthesis.md) |
| 34. Budgets should be explicit | [Topic 06](../concepts/06-experiments-and-synthesis.md), [Topic 09](../concepts/09-lifecycle-cost-and-empirical-evaluation.md) |
| 35. Comparison should judge the product, not trust worker progress claims | [Topic 05](../concepts/05-validation-and-settlement.md), [Traust experiment](../examples/01-traust-scenario-01.md) |
| 36. Strong early branches can become references without becoming truth | [Topic 06](../concepts/06-experiments-and-synthesis.md) |
| 37. Baseline settlement can reduce the need for future prompt hardening | [Topic 07](../concepts/07-provenance-findings-and-guidance.md), [Topic 09](../concepts/09-lifecycle-cost-and-empirical-evaluation.md) |
| 38. Manager output categories | [Topic 07](../concepts/07-provenance-findings-and-guidance.md), [Topic 05](../concepts/05-validation-and-settlement.md) |
| 39. Hard distinction: findings vs mutations | [Topic 07](../concepts/07-provenance-findings-and-guidance.md), [Topic 08](../concepts/08-human-control-policy-and-ui.md) |
| 40. Potential maturity signals | [Topic 09](../concepts/09-lifecycle-cost-and-empirical-evaluation.md) |
| 41. Open question: when should the system automatically suggest more experiments? | [Topic 06](../concepts/06-experiments-and-synthesis.md), [Open questions](../open-questions.md) |
| 42. Open question: when is a baseline “settled enough”? | [Topic 09](../concepts/09-lifecycle-cost-and-empirical-evaluation.md), [Open questions](../open-questions.md) |
| 43. Open question: how much of the review itself should be deterministic? | [Topic 05](../concepts/05-validation-and-settlement.md), [Open questions](../open-questions.md) |
| 44. Open question: branch retention and storage | [Topic 07](../concepts/07-provenance-findings-and-guidance.md), [Open questions](../open-questions.md) |
| 45. Open question: how to represent settled decisions | [Topic 02](../concepts/02-projects-graphs-and-authority.md), [Open questions](../open-questions.md) |
| 46. Open question: how to represent messages between projects | [Topic 04](../concepts/04-project-interface-and-negotiation.md), [Open questions](../open-questions.md) |
| 47. Open question: what belongs in prompts vs settled artifact guidance? | [Topic 07](../concepts/07-provenance-findings-and-guidance.md), [Open questions](../open-questions.md) |
| 48. Example: Traust test-harness foundation workflow | [Traust experiment](../examples/01-traust-scenario-01.md), [Baz](../examples/02-baz-reconciliation.md) |
| 49. Example: TODO application foundation workflow | [Other examples](../examples/05-school-and-other-governance-shapes.md) |
| 50. Example: spec-validation campaign | [Topic 06](../concepts/06-experiments-and-synthesis.md), [Traust experiment](../examples/01-traust-scenario-01.md) |
| 51. Example: optional interaction-model layer | [Topic 02](../concepts/02-projects-graphs-and-authority.md), [Other examples](../examples/05-school-and-other-governance-shapes.md) |
| 52. Possible relationship to interfaces/header-like artifacts | [Topic 02](../concepts/02-projects-graphs-and-authority.md), [Other examples](../examples/05-school-and-other-governance-shapes.md) |
| 53. Non-goals / cautions | [Topic 13](../concepts/13-delivery-scope-and-self-hosting.md) |
| 54. Emerging high-level architecture | [Topic 01](../concepts/01-purpose-and-system-boundary.md), [Evolution](evolution.md) |
| 55. Concise mental model | [Topic 01](../concepts/01-purpose-and-system-boundary.md), [Glossary](../glossary.md) |
| 56. Why this seems promising based on the Traust experiments | [Traust experiment](../examples/01-traust-scenario-01.md), [Topic 09](../concepts/09-lifecycle-cost-and-empirical-evaluation.md) |
| 57. Candidate architecture-doc topics to split out later | [Reading map](../README.md), [Open questions](../open-questions.md) |
| 58. Things deliberately left unresolved | [Open questions](../open-questions.md) |
| 59. Final framing | [Topic 01](../concepts/01-purpose-and-system-boundary.md), [Topic 13](../concepts/13-delivery-scope-and-self-hosting.md) |


## 4. Branch integration coverage

### S03: validation, scoped settlement and adaptive planning

| Source section | Integration |
| --- | --- |
| 1. Settlement authority versus validation execution | Topics 01, 04 and 05; evolution 16. |
| 2. Three contracts and edge-specific validation | Topic 04 in detail, Topic 05 for conformance implications. |
| 3. Risk-proportional confidence and change-class trust | Topics 05 and 08; questions 65-66. |
| 4. Objective, execution plan and settlement criteria | Topic 03; glossary; example Baz. |
| 5. Persistent dynamic plans, P1/P2, short horizon and remaining scope | Topic 03 in detail; Topic 04 negotiation. |
| 6. Settlement increments versus internal tasks | Topics 03-05; evolution 12. |
| 7. Internal checkpoints versus global settlement | Topics 03 and 05; conservative policy marked tentative. |
| 8. Completeness, soundness, proportionality | Topic 05, with examples and distinct remedies. |
| 9. Accumulated progress, regression and final integration | Topic 05; Baz example. |
| 10. Upstream obligations plus project-proposed evidence | Topics 04-05; questions 38-41. |
| 11. End-to-end loop | Topics 03-05 and Baz walkthrough; no fixed protocol imposed. |
| 12. Possible MVC scope | Topic 13; distinguished from the experiment proof slice. |
| 13. Future questions | Open-question groups D-J, with representation/ownership left unresolved. |
| 14 and final synthesis | Integrated across Topics 03-05; branch status remains exploratory. |

### S04: legacy-project reconstruction

| Source concern | Integration |
| --- | --- |
| Legacy implementation as the only explicit starting artifact | Topic 10 and Cosmos example. |
| Apparent reversal of authority | Topic 02 evidence distinction; Topic 10 explanation. |
| Evidence/provenance versus normative dependency | Topics 02 and 10; evolution 18. |
| Characterization versus acceptance | Topic 10 and Cosmos's 500-response example. |
| As-built design versus intended design | AuthService/billing example retained in Topic 10 and Cosmos. |
| Observed/reconstructed, candidate, settled | Glossary and Topic 10; labels not promoted to final state machine. |
| Import, observe, reconstruct, settle, pivot, reconcile | Topic 10 stages and Cosmos walkthrough. |
| Immediate implementation discrepancy after pivot | Topic 10 and Cosmos example. |
| Code cannot establish all purpose/rationale | Topic 10; required human/context contribution retained. |
| Proposal of alternative artifact graphs | Topic 10 and Topic 02 node criterion. |
| Reconstruction as lifecycle posture | Topic 09 and Topic 10; partial/mixed posture remains open. |

## 5. Cross-check of later main-thread additions

The source mapping specifically checks the ideas not fully present in the original manager inventory: generic non-software projects; loose project interfaces; negotiation instead of static task size; one large Implementation project; optional rather than mandatory local experimentation; coordinator-triggered project candidates; raw provenance rather than mandatory perfect self-diagnosis; delayed analysis and relevance checks; temporary worker opinions; context access distinct from direct authority; Allium as a behavioral projection; shared-model interpretation risk; and conservative partial-acceptance semantics.

The experiment discussion preserves both positive observations and limitations: no accepted Scenario 01, local probes versus integrated execution, incomplete logs, no current-model ranking, no proven cheaper fixture, no universal secret API, and no blind reliability claim for the candidate manager prompts.

The broader original interaction inventory is retained even where not first-version scope: per-node models/budgets, information-flow visibility, annotations, revised learnings, discuss-solutions and repercussions workflows, explicit reuse policies, objective non-steering monitors, maturity signals and scoped challenges.

## 6. Unresolved tensions deliberately not normalized away

- Generic reconciliation core versus broad experiment-DAG ownership.
- Standalone repository handoffs versus framework-specific adapters.
- Project-local knowledge versus global semantic analysis.
- Observable raw evidence versus affordable review and retention.
- Direct normative dependencies versus contextual access and transitive conformance.
- Stable objectives versus adaptive plans and mutable evidence strategies.
- Internal checkpoint acceptance versus independent downstream authority.
- Shared settled models versus independent discovery pressure.
- Preservation of mature state versus legitimate reopening.
- Strong human control versus review throughput and automatic acceptance.

These are explained in the topic files, history and question register rather than resolved by editorial fiat.

## 7. Source limits and nonclaims

The original source documents and selected archive files are preserved, but the full chat, full worker archive, all executable probes and the complete Traust documentation tree are not duplicated in this Markdown dossier. Original links inside copied sources can therefore require the original archives.

No external product claims have been freshly verified. In particular, earlier Allium analysis/CLI descriptions and Spec Kit feature descriptions remain assumptions for future investigation. The dossier also does not verify later outcomes of experiments proposed in the conversation.

No representation in this collection should be mistaken for a settled graph schema or implemented protocol. Diagrams and revision labels are explanatory unless tied explicitly to a preserved source commit.

The coverage review is about preserving the discussed material with its status. It is not a claim that the ideas form a complete, contradiction-free architecture or that the next project should accept them unchanged.
