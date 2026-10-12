# Glossary and notation

**Status:** editorial normalization of vocabulary already used in the discussions. Several terms remain design choices, not schema names.

## Systems and actors

**Saum:** the provisional name of the overarching reconciliation/governance system. It is not Traust.

**Traust:** the deployer/orchestrator product used in the motivating test-harness experiments. Its Scenario 01 exercises a full upgrade flow, including image discovery, persistence, runtime actions, routing, and ordering guarantees.

**Cosmos:** the hypothetical legacy application used to discuss reconstructing an authority graph from existing code. It is an example, not a proposed Saum subsystem.

**Project:** a unit with its own work products, local methodology and potentially independently settled state. A Git repository is the immediate mental model, not an established universal requirement. A project is neither automatically one microservice nor one task.

**Artifact:** an inspectable work product: documents, code, tests, a model, schema, plan, or other result. Project and artifact are sometimes used loosely in the source material. Whether one project has one or several independently governed artifact revisions remains open.

**Worker / Ralph:** a bounded producer inside a project or workflow. Ralph is the prototype worker protocol. It is not a required implementation for every project.

**Manager:** overloaded. It can mean the deterministic coordination application, an agent doing planning, a comparative reviewer, or the manual prompts being prototyped. The dossier uses the more precise name where possible; semantic work may be delegated while legal transitions remain controlled.

**Validator / reviewer:** a process evaluating evidence or conformance. It may be a program, model, human, or combination. Independent evaluation does not imply infallibility or require a distinct vendor/model in every case.

## State and authority

**Settled state:** a revision or set of claims accepted under a recorded decision/validation policy. It means accepted with available evidence, not universally true or mathematically proven.

**Candidate:** proposed state not yet accepted at the relevant authority boundary. A worker can report it ready without making it settled.

**Authority:** the right of an accepted artifact or decision to constrain downstream obligations. Authority is not the same as useful context, provenance, or evidence of current behavior.

**Normative:** prescribing what is required. **Descriptive / as-built:** recording what currently exists or happens. Reconstruction must not silently turn the latter into the former.

**Artifact dependency graph:** durable relationships among governed project states. This is distinct from the execution graph used to advance them.

**Execution DAG:** the operational arrangement of analysis, workers, experiments, review, approval and synthesis. The early model put this under Saum broadly; its later scope remains contested.

**Context relationship:** permission or encouragement to consult an artifact for rationale or background. Whether this can remain non-governing when it contains applicable requirements is an unresolved boundary issue.

**Evidence / provenance relationship:** records that a claim, reconstruction, learning, or decision was derived using a particular source. It does not automatically create normative authority.

**`based_on`:** shorthand used in examples for the upstream revisions a candidate or accepted state was reconciled against. It is not a final field specification or sufficient proof of compatibility.

**Stale / potentially affected:** an accepted state has not yet been reconciled or revalidated against a relevant changed input. The term does not necessarily mean the existing artifact is wrong. Exact invalidation granularity remains open.

**Ready for reconciliation:** the necessary inputs for the requested work are available and compatible under the chosen policy. Readiness is different from whether the existing artifact is accepted, current, or production-deployed.

**Consistent cut:** shorthand for a compatible selection of accepted revisions across relevant dependencies. Examples use one common architecture revision; the general compatibility rules are not specified.

## Work and acceptance

**Reconciliation:** bring existing accepted state into conformance with changed authoritative obligations while preserving unaffected useful decisions. It is not synonymous with deletion and regeneration.

**Reconciliation objective:** what must ultimately become true. It is separate from the plan for reaching it and from evidence offered to establish it.

**Reconciliation contract:** inputs, scope, constraints and authority for changing the objective. **Artifact contract:** what gets returned and how it can be inspected. **Validation contract:** what conformance means on a dependency relationship and how to evaluate it.

**Work package / reconciliation unit:** a bounded piece of delegated work, potentially suitable for separate review or experimentation. Earlier discussion used these to avoid assuming sequential milestones. Not every work package is an acceptance boundary.

**Settlement increment:** an externally meaningful increment with claims that can be evaluated. Its internal tasks can change without redefining its objective.

**Milestone:** originally a sequential backlog frontier, later refined into a candidate settlement/checkpoint boundary. The term should not imply all execution is sequential.

**Foundation frontier:** the smallest representative scope proposed to establish a useful baseline. Scenario 01 is one example, not a universal starting rule.

**Living settlement plan:** a persistent, versioned, adaptively recalibrated plan linking the objective to proposed increments and remaining obligations. This term comes from the validation branch.

**Planning horizon:** how far ahead the next work is planned. It can be shorter than the total reconciliation obligation without silently dropping the remainder.

**Accepted internal checkpoint:** work accepted for progress within an ongoing reconciliation. Under the conservative branch proposal it is not automatically a globally settled upstream release.

**Global settlement:** acceptance of a project revision for use under the cross-project dependency policy. The conservative initial direction waits for complete reconciliation and integrated validation.

**Scoped reopening:** explicitly reconsider selected accepted decisions while preserving the rest initially. Evidence may justify widening the scope through the relevant authority process.

## Experiments and evidence

**Materialization experiment:** search for a way to advance a target artifact. Code or other outputs may be reused under an explicit policy.

**Validation experiment:** attempted materialization used to pressure an upstream artifact. The intended output is knowledge; implementation may be deliberately disposable.

**Challenge experiment:** test whether accepted local decisions remain suitable under a consequential change. It need not reset the entire project to foundation discovery.

**Independent candidates:** alternatives started from controlled common inputs without unintended sibling visibility. Independence is an experimental condition, not a claim that their reasoning is statistically independent.

**Validity monitor:** checks whether a run has crossed an already established invalidation/stop boundary. It should not covertly coach an experiment toward a preferred architecture.

**Learning:** an interpretation drawn from evidence. Learnings can be wrong, revised and reviewed; they are not themselves automatically authority.

**Synthesis:** construct a coherent result from selected learnings and permitted references. It is not necessarily choosing one winner or merging all branches.

**Guidance:** scoped, informative input intended to steer a later worker within existing authority. It is separate from a mandatory new requirement or permanent worker-protocol change.

**Completeness:** whether obligations are covered. **Soundness:** whether the evidence establishes what is claimed. **Proportionality:** whether the approach's cost/complexity is justified. These are distinct review questions.

**Characterization test:** demonstrates current behavior, including possibly undesirable behavior. **Acceptance test:** demonstrates behavior the relevant authority has decided is required.

## Example notation

`A12 -> A13`, `D8 -> D9`, `TH4 -> TH5`, `B7 -> B8`, and `I27 -> I28` indicate illustrative revisions of Architecture, Design, Test Harness, Behaviour Model and Implementation. They are not actual repository identifiers unless explicitly tied to an archived experiment commit.

The validation branch uses `TH1`, `TH2`, etc. as increment labels as well as harness-shaped names. Topic documents distinguish a checkpoint label from a globally released project revision; the original source remains unchanged.

**MVC / MVP:** both appear in the discussions. The branch says MVC for a minimal initial implementation but does not expand the acronym. The dossier preserves that ambiguity rather than defining a new product milestone.

**Allium:** the candidate language/tool used inside a proposed Behaviour Model project. The project responsibility should remain understandable without requiring Allium to be part of Saum's core.

## Sources

[S01: original manager inventory](history/sources/01-original-manager.md), especially sections 2-4, 6-8, 21-26 and 38; [S03: validation branch](history/sources/03-validation-scoped-settlement-and-adaptive-planning.md), sections 1-10; [S04: reconstruction branch](history/sources/04-legacy-project-reconstruction.md). Main-thread terminology changes are recorded in [evolution](history/evolution.md).
