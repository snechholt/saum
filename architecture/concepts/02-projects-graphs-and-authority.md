# 02. Projects, graphs, revisions and authority

**Status:** the two-graph distinction and preservation of accepted state are working directions. Edge semantics, invalidation granularity and graph restructuring remain open.

## 1. The durable graph and the operational graph

The artifact/project graph describes durable relationships between work products. The execution graph describes how a particular candidate is currently being produced or evaluated.

```text
Durable graph:
Architecture -> Design -----------+
             -> Test Harness -----+-> Implementation

One possible execution:
analyze -> alternatives -> compare -> human review -> synthesize -> validate
```

These graphs answer different questions. The first asks what constrains what, and what accepted state depends on which revisions. The second asks what work is running, with which inputs, resources, isolation and decision gates.

The same project may use different execution methods over time without changing its place in the durable graph. A project's internal execution can itself contain subagents, tests, PRs and local approvals without every such object becoming a top-level node.

## 2. When something deserves a node

The working criterion is **independently meaningful authoritative state**. A component design, interface schema or interaction model deserves its own project when it should be governed, revised, depended on, or reconciled independently.

A document does not deserve a node merely because it can be generated. A local human review does not automatically imply a missing graph layer. Conversely, repeated requirements to approve component boundaries separately, or multiple consumers needing one stable interface, may be evidence that a separate artifact would be useful.

The graph is configurable. Architecture -> Implementation can be reasonable. So can Architecture -> Design -> Implementation, with a parallel Test Harness. Interaction Models may be an optional intermediate layer. API schemas, deployment models, operations guidance or security models are possible additional authorities, not required departments in every system.

The realistic Traust example deliberately has one Implementation project. Adding Baz does not automatically create separate Foo, Export, Timeline and API projects. Those can remain concerns inside one repository. Earlier examples that made each consequence a project hid the difficult intake problem.

## 3. Project, repository, artifact and revision are not yet identical

A Git repository is the immediate implementation model because it can preserve candidate and accepted states and support isolated work. It is not settled that every project is exactly one repository, every commit is a meaningful artifact revision, or every file has equal authority.

Traust's documentation illustrates this. A single documentation repository distinguishes architecture rules, implementation contracts, exact appendix contracts and scenarios. The project graph must not flatten those internal authority distinctions merely because the repository is one node.

Possible representations of accepted decisions include code and tests, architecture documents, ADR-like notes, structured records, review reports and revision metadata. The original discussion expected a combination. The exact minimal combination remains open.

## 4. Three relationship meanings that should not be conflated

### Normative dependency

A downstream state claims compatibility with applicable obligations from identified upstream state. A change can require impact analysis and reconciliation. Different edges mean different things: behavioral fidelity, architectural conformance, test coverage, or passing an external harness are not one interchangeable condition.

### Contextual access

A project may need architecture to understand why a constraint exists even when its immediate work is driven by Design and Test Harness. The user proposed similarly broad access to a Behaviour Model. Reading context is valuable without requiring every relationship to impose its own scheduling gate.

### Evidence and provenance

An artifact or decision can be reconstructed or justified using another artifact without being normatively governed by it. Legacy code can be evidence for candidate architecture. An experiment can support a learning. Neither automatically defines what the future system must do.

This is a conceptual distinction, not a finalized typed-edge schema. The same pair of projects could have more than one relevant relationship.

## 5. Authority goes down; findings can go up

A settled upstream change can create obligations downstream. A downstream discovery can create an upstream review question or proposal, but cannot silently rewrite the authority it was supposed to satisfy.

This permits a useful asymmetry. Downstream work can say "these requirements conflict," "this interpretation creates an unexpected lifecycle property," or "this workaround suggests a missing abstraction." The upstream project then determines how to revise its artifact, subject to its own settlement process.

Reverse messages do not inherently reverse the normative graph or imply a normative dependency cycle. This distinction becomes essential during legacy reconstruction.

## 6. The basic revision example

```text
Architecture A12
Design D8       based on A12
Test Harness T19 based on A12
Implementation I27 based on D8 and T19
```

After A13 is accepted, D8 and T19 need reconciliation or a compatibility assessment against the relevant change. I27 remains the existing accepted implementation for its previous input set. It is not automatically deleted or ready to adopt only half the new inputs.

D9 and T20 can be developed in parallel. Under the conservative example, Implementation begins reconciliation once both are accepted against the compatible new basis. I27 then becomes the base for I28 rather than being regenerated without regard to its accepted nuance.

## 7. Stale is not the same as incorrect

Earlier discussion often used "parent revision changed, therefore child is stale" as shorthand. The intended useful fact is that the child's compatibility with the new input has not yet been established. It does not prove the artifact is semantically wrong.

The distinction between revision-level invalidation and semantic impact is unresolved. A wording-only edit, a change to an irrelevant scope, and a changed behavioral obligation may warrant different responses. Likewise, whether an unchanged artifact can be revalidated against a new input without a content change is a question to settle, not an established mechanism.

Keeping this uncertainty explicit prevents a graph implementation from accidentally encoding "every document edit requires every descendant to be rewritten."

## 8. Consistent state and readiness

The UI could distinguish the latest accepted revision of each project, ongoing candidates, stale or potentially affected projects, blocked work, and a compatible set of accepted revisions suitable for use together.

The "latest consistent cut" is useful when Architecture has moved to A13 but the production implementation still corresponds to an A12-derived configuration. The newest revision of every project taken independently may not be a usable combination.

Examples with one shared architecture revision are illustrative. Exact compatibility may involve multiple parents, different scopes or transitive context. The discussions have not defined an algorithm, nor established that matching revision numbers alone proves compatibility.

Accepted internal checkpoints are another dimension. The validation branch proposes retaining TH1's accepted Scenario 01 progress inside a larger reconciliation without advertising the whole harness as globally current. See [validation and settlement](05-validation-and-settlement.md).

## 9. The unsettled dependency-versus-context problem

One proposed graph gives Implementation hard prerequisites only on Design and Test Harness, while Architecture and Behaviour Model remain visible context. This reduces redundant direct readiness relationships and reflects an intended transitive flow.

However, the rationale for omitting an edge cannot simply be "it creates inconvenient staleness." If Implementation is expected to obey an architectural constraint not carried by Design or Harness, that constraint still needs an accountable path into its objective and validation. Reading a newer contextual document must not silently change the basis on which a candidate was accepted.

The discussions did not settle whether contextual inputs are pinned, whether conformance is checked transitively, or whether some direct architecture/behavior dependencies remain necessary. These are explicit open questions, not reasons to silently convert all context into mandatory dependencies either.

Similarly, disagreement between Architecture and a settled Behaviour Model is not resolved just by saying one is higher in the diagram. The intended scope of delegated interpretation needs to be clear.

## 10. Restructuring the graph

Projects may propose inserting a Design layer, splitting a project, or establishing a shared interface authority when local work reveals governance pressure. This should be a proposal rather than an automatic mutation by a worker.

The key question is what new independently governed knowledge the change creates. Turning every internal workflow stage into a project would make the graph unmanageable; refusing all restructuring would fossilize an unsuitable composition.

Migration of revision links, historical provenance, active reconciliations and acceptance responsibility during restructuring remains unresolved. The old graph should remain intelligible as history even if future work uses a different one.

## Sources

[S01](../history/sources/01-original-manager.md), sections 2-3, 23-27, 32, 45-46 and 51-52; [S02](../history/sources/02-project-boundary-notes.md), realistic Architecture/Harness/Design/Implementation example; [S03](../history/sources/03-validation-scoped-settlement-and-adaptive-planning.md), sections 2 and 7; [S04](../history/sources/04-legacy-project-reconstruction.md), authority versus evidence. Main-thread episodes M05-M09 are indexed in [history](../history/sources-and-coverage.md).
