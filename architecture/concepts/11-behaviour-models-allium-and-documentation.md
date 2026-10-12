# 11. Behaviour models, Allium and architectural documentation

**Status:** the user prefers arc42-style documentation at the top. Allium as a derived Behaviour Model is promising but not adopted or tested here. Tool-specific capability claims remain unverified in this dossier.

## 1. Why Allium entered the discussion

The user discovered JUXT's Allium through a podcast and proposed two roles: a top-level specification, or a downstream materialization of broader architecture that would reveal holes when the prose had to be expressed more precisely.

After discussion, the user preferred arc42 at the top. The interesting role for Allium became a Behaviour Model project below the rich documentation, useful both as a pressure test of that documentation and as a more efficient source of intended behavior for downstream work.

This should not be reduced to "add support for another framework." It is an example of a heterogeneous project whose output can make upstream ambiguity more visible.

## 2. Broad architecture and a constrained behavioral projection

The intended broad architecture contains context, structure, deployment, constraints, decisions, quality goals, risks and rationale as well as behavior. A behavioral representation would focus on applicable entities, rules, states, transitions, invariants and interactions.

The division discussed was roughly: Architecture explains the larger system and why constraints exist; Behaviour Model makes relevant behavioral commitments precise; Design proposes the structure that realizes them; Test Harness determines how they can be demonstrated independently; Implementation builds the result.

That division is a working conceptual fit, not a claim that every relevant concern can be perfectly partitioned. Behavioral and structural choices can interact, and disagreements still need a governed resolution path.

## 3. "Compilation" as a useful but limited metaphor

The user's idea of compiling the documentation means materializing it into a more constrained representation, forcing missing information to become visible. It does not imply that an established deterministic compiler can recover the unique intended model from prose.

Earlier responses described Allium constructs and tooling for lifecycle analysis, invariants and test obligations, including a command named `allium plan`. Those were external-product claims made in the earlier conversation. They have not been verified for this dossier and must not be treated as established capabilities, current command syntax, or guarantees of completeness.

The architectural proposal survives that uncertainty: a constrained behavioral model, with whatever actual checks its tools support, could be a cheaper intermediate pressure surface. The concrete Allium version, supported semantics, checker soundness and test-obligation workflow need investigation before implementation relies on them.

## 4. What modeling pressure could reveal

A sentence such as "deployment succeeds once the candidate is running" leaves many questions: what starts the transition, what running means, what can fail before it, when routing changes, what happens to old state, and what invariant must hold throughout.

Traust's documentation already contains a richer behavioral skeleton than that example: tenant state, deployment context, migration boundaries, readiness, internal and hostname verification, freeze/unfreeze behavior, routing and recovery rules. A model may make inconsistencies between those sources harder to ignore.

The value can arise while writing the model, not only from a formal checker. A missing precondition or unclear event identity may require an upstream question even if the tool itself does not detect it.

## 5. Traust's authority is already internally layered

The archived Traust Architecture Index assigns architecture the system shape, canonical architectural rules and rationale; Specs own implementation contracts and mechanics; Appendices own exact schemas, transition tables and reference rules; Tests own scenario validation and expected evidence.

When architecture summarizes an exact appendix contract, the appendix remains authoritative for those details. The risk register explicitly identifies contract drift among architecture, specs, appendices, tests and implementation as R-05.

Thus the proposed Behaviour Model is not just a paraphrase of one arc42 page. It would need to project applicable behavior from the documented authority structure without creating another competing definition.

Whether those areas remain one Documentation/Architecture project or become separate projects is a governance decision. The existing directory structure alone does not decide it.

## 6. The Test Harness dependency

The user's proposal was that the harness still be created before implementation from architecture-defined scenarios, but additionally depend on Allium for intended behavior. That is a natural candidate relationship:

```text
Architecture/docs -> Behaviour Model
Architecture/docs -------------------+
Behaviour Model ---------------------+-> Test Harness
```

The harness could use precise obligations rather than rediscover all ordering requirements from scattered prose. It would still need broader documentation for real environmental constraints, allowed fixtures, independent SUT boundaries and integration mechanics.

A behavioral rule that discovery advances a target does not by itself establish how an external registry boundary should be exercised, or whether a particular observation would be legitimate. That remains a harness-design and validation problem.

## 7. Design and Implementation should understand the rationale

The user wanted essentially all projects to have access to Architecture and, now, the Behaviour Model. Architecture can provide the underlying why; the model can clarify what must remain true. Both can be useful even when immediate tasks are driven by other artifacts.

A proposed Design dependency on both Architecture and Behaviour Model makes sense when component/interactions must realize settled behavioral commitments. The model does not automatically make Design redundant.

Implementation might consume Design and Harness as direct prerequisites while reading Architecture and Behaviour Model transitively. Alternatively, direct obligations on the latter may be needed. The discussions did not settle that graph. The unresolved context-versus-authority issue is detailed in [graphs](02-projects-graphs-and-authority.md).

## 8. Shared precision can hide a shared mistake

A model-producing agent may silently choose one interpretation of ambiguous prose. If both Harness and Implementation then follow that model, they can agree with each other while missing the original ambiguity.

This is not a reason to reject a shared behavioral model. It changes how it should be validated and settled. During discovery, independent derivations or separate harness/implementation probes can pressure the source. After the interpretation has been deliberately accepted, downstream consistency with it becomes useful.

The proposed distinction is therefore between an unreviewed projection and a settled interpretation. A clean model and passing generated tests are not sufficient evidence that the original prose was faithfully or uniquely interpreted.

## 9. Independent model derivations

One proposed experiment is several Allium/behavioral derivations from the same architecture, followed by comparison before accepting a model. This may be cheaper than several full implementations, but that relative cost is a hypothesis until measured.

Material differences can reveal ambiguity, intentional design freedom, or modeling mistakes. Earlier responses described divergence very strongly as ambiguity evidence. The safer design question is whether the alternatives are actually faithful and consequential; difference alone does not settle the diagnosis.

The comparative process should preserve those distinctions rather than turning every model disagreement into an upstream defect.

## 10. Discovery and settled use need not share the same input policy

During validation of a draft architecture, an independent Harness spike need not be constrained by the same unaccepted model used in another branch. That preserves a useful independent interpretation surface.

After a Behaviour Model revision is accepted as a scoped interpretation of the architecture, normal downstream work can use it as a common basis. The relationship between model settlement and any corresponding architecture clarification remains open: some resolved choices may need to be reflected upstream, while other precision may be legitimately delegated.

The system should be able to express the distinction rather than accidentally sharing draft outputs as if they were authority.

## 11. A general pattern beyond Allium

A derived project can translate a broad artifact into a representation that forces precision. Examples discussed included API/IDL contracts, schemas and formal policy models as well as behavior specifications.

The output may be a durable downstream authority, a validation artifact, or both under explicit scope. The important operation is attempted constrained materialization plus review, not a specific language brand.

The proposed cost ladder was textual review, constrained modeling, supported analysis, independent model alternatives, harness/design materialization and finally real implementation pressure. It is not a guarantee that every later rung costs more or catches a strict superset of defects. Different representations expose different classes of problem.

## 12. Limits that preserve the experiment thesis

Some Scenario 01 difficulties involve actual Docker, registry, filesystem and observability behavior. A purely behavioral model cannot be assumed to reveal every environmental constraint or valid evidence path.

Allium therefore complements rather than replaces attempted materialization. Its strongest possible contribution is exposing some uncertainties earlier and giving downstream work a clearer accepted behavioral basis.

The proposed node name is **Behaviour Model**, with Allium as its chosen internal method. That keeps the graph's purpose understandable if the tool changes.

## Sources

Main-thread episodes M08-M09 in [source map](../history/sources-and-coverage.md); [Traust Architecture Index](../history/traust-reference/docs/architecture/00_architecture_index.md), [risk register](../history/traust-reference/docs/architecture/11_risks_and_debt.md), [Scenario 01](../history/traust-reference/docs/tests/01_happy_path.md); [S03](../history/sources/03-validation-scoped-settlement-and-adaptive-planning.md), section 2, for edge-specific validation. Tool capability details were not independently verified for this synthesis.
