# Open questions and unresolved architectural tensions

**Status:** question register, not a requirements backlog or a commitment to implement every possibility. Some entries restate a discussion tension as an explicit review prompt; they do not introduce settled answers.

Question IDs are for cross-context reference. The groups are organizational, not proposed project nodes. The most consequential early choices concern project meaning, authority/context, intake contracts, settlement and the experiment boundary.

## A. Purpose and scope

**Q01. What is the smallest coherent Zonotron product?** Is it primarily a governed revision/reconciliation layer, an experiment-and-learning system, or a deliberately staged combination? The current sources support both prototype paths.

**Q02. How much genericity is required initially?** Should the first implementation directly support non-software projects, or merely avoid assumptions that prevent them later?

**Q03. What can existing tools supply adequately?** Which execution, isolation, planning and review responsibilities disappear when composed with existing methods, and which settlement semantics still need new machinery?

**Q04. What outcome would justify the additional governance cost?** Better long-horizon coherence, fewer false successes, less rework and human review effort are candidates; none is established as the sole objective.

## B. Projects and authority

**Q05. What is the precise project/artifact relationship?** Can one project expose multiple independently accepted outputs, or should independently governed outputs become separate nodes?

**Q06. When should a local artifact be promoted?** How should independent consumers, repeated approval pressure and reconciliation needs weigh against graph complexity?

**Q07. Which internal documents own which obligations?** A repository can contain architecture, exact contracts, scenarios and guidance with different authority. How is that map exposed without flattening it?

**Q08. Which graph relationships are normative, contextual or evidentiary?** Are they different edge types, separate metadata, or another representation?

**Q09. How is conflicting upstream authority handled?** Multiple parents can be individually accepted yet disagree. Who identifies and settles the conflict, and what work can continue?

**Q10. What does acceptance of a derived interpretation authorize?** When may Behaviour Model precision stand on its own, and when must it cause an architecture clarification?

## C. Revisions, context and consistency

**Q11. What counts as a revision?** Are Git commits sufficient, or are logical artifact revisions and acceptance records separate from raw commits?

**Q12. What does a parent change invalidate?** All direct dependants, only affected scopes, or a need for compatibility assessment before deciding? How are nonsemantic edits handled?

**Q13. Can unchanged content be revalidated against new authority?** If so, where is the new compatibility claim recorded?

**Q14. How is contextual access versioned?** Can a worker read newer architecture than the one its accepted Design/Harness used? What prevents silent mixed-version authority?

**Q15. When are direct Architecture/Behaviour Model dependencies required for Implementation?** Can immediate parents faithfully carry all obligations, and how is that established?

**Q16. What makes a consistent cut?** Matching revision bases is not by itself proof of semantic compatibility. What checks or claims define a usable set?

**Q17. What happens when upstream state changes during reconciliation?** Finish against the old basis, rebase, cancel, or negotiate a new objective? Preserve the original request either way.

**Q18. How are graph restructurings governed?** Adding, removing or splitting layers must preserve history and account for active work and inherited obligations.

## D. The project contract

**Q19. What is the minimum mandatory interface?** Which inputs, candidate claims, evidence references, statuses and capabilities must be available across all supported methods?

**Q20. How loose can the handoff remain?** Are Git references and structured Markdown sufficient for the first version, and when does a driver or specialized adapter earn its cost?

**Q21. What inspection capability does an artifact need?** What happens for proprietary formats, unavailable tools or evidence that cannot be interpreted automatically?

**Q22. What local observations are required?** Explicit blockers may be straightforward, but how much provenance must replace potentially lossy self-diagnosis?

**Q23. How are public success claims distinguished from worker status?** The v5 false pass shows why a single completion field can be insufficient.

**Q24. How are analysis-only requests and execution permissions separated?** A project should be able to plan or propose graph changes without implicitly receiving mutation authority.

## E. Negotiation and planning

**Q25. Who has final authority over a decomposition?** The project knows local viability; Zonotron knows global constraints. What happens when they disagree?

**Q26. What terminates negotiation?** Define sufficient viability without requiring an optimal plan or endless discussion.

**Q27. What makes a meaningful settlement increment?** It should be reviewable without freezing internal implementation tasks or creating excessive administrative overhead.

**Q28. How are atomic, sequential and parallel increments represented?** Some work can execute independently but still require integrated acceptance.

**Q29. How is unplanned remaining scope kept visible?** A short planning horizon must not hide missing obligations.

**Q30. What triggers replanning?** Every accepted increment, new evidence, cost/complexity changes, a blocked dependency, or explicit request? Must a no-change plan be reissued?

**Q31. Who approves plan revisions and records their reasons?** Distinguish harmless regrouping from changing the objective or weakening evidence.

**Q32. How are global priorities reconciled with local sequencing?** Business value, risk, critical-path progress and cheapest useful acceptance can conflict.

**Q33. How is the final integration obligation defined early enough?** A plan can evolve, but completion cannot be just the sum of checked tasks.

## F. Validation architecture

**Q34. Who defines validation contracts?** Project, upstream owner, coordinator, independent validator or human? Their maintenance and authority need to be clear.

**Q35. Are validation contracts themselves versioned and settled?** How are historical acceptance decisions interpreted after criteria change?

**Q36. How are validators selected?** Deterministic tests, semantic agents, formal tools, experiments and humans establish different kinds of evidence.

**Q37. How much semantic understanding belongs centrally?** Can it be delegated while keeping acceptance responsibility and cross-project interpretation accountable?

**Q38. What happens when evidence is insufficient?** Distinguish unknown conformance, unavailable inspection, actual failure and the choice to accept residual risk.

**Q39. How are shared blind spots challenged?** Passing producer-authored tests or matching one supplied good actor is not enough for every obligation.

**Q40. How are permitted variations tested?** A harness can falsely reject a valid implementation by depending on private names, invocation methods or retained resources.

**Q41. What evidence establishes temporal claims?** Event initiation, completion, route publication, sampling and eventual state must not be conflated.

**Q42. How many independent reviewers are useful?** Independence and disagreement can help, but add cost and do not guarantee correctness.

## G. Checkpoints and settlement

**Q43. How are accepted claims represented?** Avoid an excessively detailed requirements database while retaining enough meaning for regression and final validation.

**Q44. What exactly distinguishes an internal checkpoint from global settlement?** Clarify consumers, revision identity and authority implications.

**Q45. How is an accepted checkpoint reopened?** What evidence invalidates earlier validation, and which later increments must be reconsidered?

**Q46. When may downstream work consume partial upstream state?** The branch favors waiting initially; future scoped consumption needs explicit provisionality and rework semantics.

**Q47. Can increments become independently authoritative?** At what point does that imply a different project graph rather than more checkpoint machinery?

**Q48. What final evidence is required beyond accepted increments?** Cross-cutting, regression and integration claims need a clear home.

## H. Experiments and synthesis

**Q49. Where does multi-candidate reconciliation belong?** Core, optional service, project-local capability or several interchangeable implementations?

**Q50. When should uncertainty be discussed rather than experimentally explored?** Cheap decisions can suppress discovery; broad experiments can waste resources.

**Q51. What scope should fan out?** One package, one project or a subtree? The latter remains explicitly exploratory.

**Q52. What constitutes adequate independence?** Common inputs, sibling isolation and real environmental side effects all matter.

**Q53. Which conditions justify automatic termination?** Keep objective invalidity distinct from taste, high cost or suspicious progress.

**Q54. How should synthesis reuse experimental artifacts?** Learnings-only, references, transplants and continuation have different coherence/cost tradeoffs.

**Q55. When may reference access undermine the purpose of a fresh attempt?** Useful reuse and unwanted imitation need an explicit experimental policy.

**Q56. How are budgets accounted for?** Tokens, cached usage, process time, human review time and actual price are not interchangeable; iteration count is inadequate alone.

## I. Findings, evidence and learning

**Q57. What provenance must always be retained?** Identify the minimum recoverable chain from input through work and validation to acceptance.

**Q58. How is evidence indexed without irreversibly summarizing it away?** Hierarchical sessions and hints may help, but the raw basis for consequential conclusions must remain reachable.

**Q59. How is review depth selected?** Foundation review and routine maintenance should not necessarily consume the same resources.

**Q60. How are delayed findings checked against current state?** Historical relevance, active actionability and duplicates need distinct treatment.

**Q61. How are rejected proposals retained without flooding humans?** Preserve causal traces while surfacing only material residual concerns or required disclosures.

**Q62. Where does a learning belong?** Stable protocol, project architecture, temporary guidance, an upstream question or only historical rejection?

**Q63. When does temporary guidance expire or become durable?** Avoid both forgetting useful knowledge and fossilizing every worker mistake.

**Q64. How are reviewer interpretations challenged and revised?** Keep human annotations and successive learnings tied to the same underlying evidence.

## J. Autonomy and operations

**Q65. Which changes qualify for automatic acceptance?** Calibrate by risk and demonstrated change-class reliability rather than one global trusted/untrusted flag.

**Q66. How should auto-accepted work be audited?** Sampling, confidence standards and response to failures remain unspecified.

**Q67. How are retries, cancellation, timeouts and partial output represented?** The sources identify the need but do not choose a runtime/state machine.

**Q68. How are resource and information isolation enforced?** A repository branch does not isolate a shared engine, external service or credential.

**Q69. What retention/access policy protects evidence and sensitive material?** Original sources name storage and security concerns but do not provide full policies.

**Q70. Which UI surfaces belong in a first version?** Basic objective/plan/status visibility may be enough before editable DAGs and annotation-driven workflows.

## K. Reconstruction and behavioral models

**Q71. When is reconstructed knowledge ready for the authority pivot?** Separate factual confidence from decisions about what should be preserved.

**Q72. How are bugs, debt and compatibility behavior classified?** Characterization alone cannot decide which behavior becomes acceptance criteria.

**Q73. Can reconstruction and normative governance coexist in different scopes?** The branch proposes a lifecycle posture but leaves partial pivots open.

**Q74. What can the selected Allium tool/version actually establish?** Verify language coverage, commands, supported analysis and limits before making them dependencies of the design.

**Q75. How is fidelity of a behavioral projection validated?** Internal model consistency does not establish faithful or unique interpretation of the prose.

**Q76. When should independent derivations be preserved rather than shared?** Discovery benefits from independent interpretation; settled execution benefits from a common basis.

## L. Proof and self-hosting

**Q77. What experiment would most reduce uncertainty about Zonotron now?** Blind evaluator reliability, guidance transfer, a living-plan reconciliation, or preservation across an upstream change are different tests.

**Q78. What would count as a failed thesis?** The sources argue for comparing against simpler workflows, including human attention and review cost; define meaningful success/failure rather than only collecting interesting examples.

**Q79. How should this dossier become Architecture A0?** Identify and settle decisions explicitly instead of promoting polished synthesis into authority.

**Q80. How does self-hosting avoid self-authorizing changes?** The existing distinction between proposals and settlement provides a principle, but the bootstrap trust arrangement remains open.

## Sources and next use

These questions integrate [S01 sections 41-47 and 57-58](history/sources/01-original-manager.md), [S02](history/sources/02-project-boundary-notes.md), [S03 section 13](history/sources/03-validation-scoped-settlement-and-adaptive-planning.md), [S04](history/sources/04-legacy-project-reconstruction.md), and the main-thread tensions indexed in [history](history/sources-and-coverage.md).

They are a menu for subsequent architecture work, not eighty questions the user must answer before making progress. A bounded prototype can proceed under explicit provisional choices while retaining the unresolved remainder.
