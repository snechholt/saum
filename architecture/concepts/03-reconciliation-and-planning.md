# 03. Reconciliation, impact analysis and living plans

**Status:** reconciliation rather than regeneration is a working direction. Persistent adaptive plans are a promising branch proposal, not a final planning protocol.

## 1. Preserve accepted state while changing its basis

Reconciliation starts from an existing accepted artifact and changed authoritative inputs. The goal is to satisfy the new obligations while preserving unaffected decisions and useful accumulated structure.

This differs from treating code, tests or design as disposable renders of prose. Downstream work contains reviewed choices that may not be restated completely upstream. Wiping that work would discard knowledge and can recreate previously resolved bootstrap mistakes.

Preservation is not absolute. A material change may justify reopening an accepted routing strategy, observation mechanism or component boundary. The intended operation names what is reopened and what initially remains preserved; evidence can justify asking to widen that scope.

## 2. Global and local impact analysis

Zonotron knows the governed project relationships, accepted input revisions, cross-project dependencies and overall priorities. A project knows its actual internal state, local domain, coupling and chosen method. Neither should be assumed to know the other's world completely.

For Architecture A12 -> A13 introducing Baz, Zonotron can identify potentially affected projects and request local impact analysis. Design and Test Harness may respond separately. Implementation may not yet be ready because its future input set depends on both.

A realistic implementation project may later receive changes involving Baz itself, Foo/Bar interactions, timeline entries, exports, APIs and tenant templates. Those consequences need not be different projects. They create a substantial intake problem inside one project.

The working response is negotiation. Zonotron and the project shape a viable decomposition together rather than assuming either one can produce the final plan in isolation.

## 3. Three separate objects of reasoning

The validation branch makes a central distinction:

**Objective:** what must ultimately become true, grounded in upstream authority.

**Execution plan:** the current hypothesis about how to achieve it.

**Settlement criteria:** the evidence required to accept the resulting claims.

For Scenario 01, a tentative plan might include a registry fixture, a Docker interception mechanism, a replaceable stand-in and the scenario assertions. Those are possible implementation steps, not automatically four permanent requirements. A better approach may remove or replace them while still meeting the objective.

Conversely, completing every planned task does not establish that the objective has been met. The plan cannot redefine acceptance around whatever the implementation happened to accomplish.

## 4. Persistent, adaptively recalibrated plans

The user's branch proposal is to retain the negotiated breakdown permanently and allow it to change as work teaches the project more. The plan is a versioned artifact of the reconciliation, not an ephemeral TODO and not a new independent upstream authority by default.

Illustrative objective:

> Produce a Test Harness satisfying the scenario-testing obligations of A13.

Illustrative first plan:

```text
P1
[ ] Increment 1: establish S01 coverage
[ ] Increment 2: establish S02-S05 coverage
[ ] Increment 3: establish S06-S20 coverage
```

After the first accepted checkpoint exposes unexpected infrastructure complexity:

```text
P2
[x] Increment 1: establish S01 coverage
[ ] Increment 2: establish S02 coverage
[ ] Increment 3: establish S03-S05 coverage
[ ] Increment 4: establish S06-S20 coverage
```

P1 remains available, along with why P2 replaced it. The object of acceptance is not "the plan was right" but the justified progress and subsequent revised plan. The evolution itself may reveal overly optimistic assumptions or recurring upstream friction.

The source branch uses TH1/TH2 as increment labels. Here they are called increments to avoid confusing them with globally released harness revisions.

## 5. Planning horizon is not obligation scope

A plan may describe only the next meaningful increment. This can be appropriate when foundation work is uncertain and speculative detail would be wasteful.

However, the unresolved remainder must remain explicit. Knowing that S01-S20 are required while planning only S01 is different from silently forgetting S02-S20. Planning can be incomplete without making the objective incomplete.

Likewise, splitting S02-S05 into smaller increments is plan recalibration. Removing S17 from the obligation is a scope/authority change. These should not be hidden behind the same operation.

The discussions leave open whether a project must reassess after every accepted increment or only when something changes materially. Reassessment can be a useful default; a mandatory expensive replan after trivial progress has not been established.

## 6. Settlement increments and internal tasks

The negotiated unit should usually be meaningful to review: Scenario 01 coverage, a settled interface, or another externally legible result. Beneath it, the project may freely reorganize fixtures, runtime adapters, observers, assertions and fixes.

The original milestone concept bundled useful stopping boundaries with sequential execution. The main thread suggested work packages or reconciliation units to remove that assumption. The later branch restores the useful meaning of milestone as a candidate settlement boundary, not as a compulsory linear schedule.

These distinctions can coexist. A settlement increment may contain several tasks. A work package can be a bounded experiment without producing an accepted revision. Some increments can be parallel, some dependent, and some need joint validation. The final naming and representation remain open.

## 7. Negotiation should find a workable plan, not an optimal one

Potential planning concerns include local coherence, cross-project ordering, evidence quality, reviewability, resource cost, risk and independence. A single static "optimal task size" cannot capture all of them.

A project can propose that one request is consumable as a feature spec, another needs prerequisites, and another should be split because the evidence would otherwise be uninterpretable. Zonotron can ask whether units can be combined or whether an apparently local decision affects another dependency.

A useful candidate stopping condition is: the receiving project considers the work executable, and Zonotron considers the plan compatible with the global objective and acceptance constraints. That is a good-enough contract, not proof that no better plan exists.

The termination policy, disagreement handling and approval rights remain open. Endless planning would recreate the cost problem the system is meant to solve.

## 8. Choosing how to resolve uncertainty

Before expensive fan-out, the user proposed distinguishing questions that merit discussion from those that merit experiments. Some can be settled cheaply using an accepted convention. Some require clarification by the human. Others are consequential enough that independent attempts may reveal hidden assumptions.

The objective is not to decompose every mechanism or predesign the project. It is to identify uncertainty with architectural leverage. Scenario 01 may be the right representative milestone and still be too broad for the first undifferentiated experiment.

There is a tension to preserve. Cheap early decisions can suppress valuable discovery. Deciding to fake a registry, choose a secret-file convention, or use only container mounts may make progress easier while hiding the ambiguity that alternative implementations would expose. Conversely, experimenting on trivial plumbing can waste substantial resources.

No universal rule was settled. "A small OCI fixture is cheaper/better" remained a hypothesis in the evidence, not an approved infrastructure policy.

## 9. Prioritization across projects

The user asked whether Zonotron or individual projects should prioritize work. The working distinction is global intent and dependencies versus local execution realities.

Zonotron can care about business importance, readiness, bottlenecks and risk reduction across the graph. A project can explain what must be established first locally and which work benefits from being grouped. A negotiated plan should reconcile those views rather than imposing an arbitrary ordering from either side.

The discussions did not choose a scoring function, resource-allocation policy, or a single definition of priority. Critical-path progress, cheapest useful settlement and reduction of consequential uncertainty may point to different next actions.

## 10. Accepted progress and downstream release

The conservative branch direction allows internally accepted increments while withholding a globally settled project revision until its entire reconciliation objective and integrated constraints have been validated.

This protects downstream projects from depending on a foundation that later local work will substantially revise. It does not make an accepted checkpoint immutable: new evidence can reopen it through a recorded process.

Partial downstream consumption is a possible later capability. It would need explicit scope, dependency and provisionality semantics. The branch does not establish it as an initial requirement, and it should not be implied by showing a checked box in a plan.

## 11. Planning failures worth distinguishing

A missing obligation is a completeness problem. A test that claims the wrong guarantee is a soundness problem. An unnecessarily elaborate mechanism is a proportionality problem. A temporarily unusable environment is not automatically a specification contradiction. A changed implementation plan is not automatically a failed reconciliation.

Keeping these distinct should improve both recovery and upward learning. The right response may be ordinary local work, independent validation, a new evidence strategy, an upstream clarification, or a scoped replan.

## Sources

[S01](../history/sources/01-original-manager.md), sections 6-7, 21-26 and 29; [S02](../history/sources/02-project-boundary-notes.md), discussion of intake size, project-local knowledge and parallel work; [S03](../history/sources/03-validation-scoped-settlement-and-adaptive-planning.md), sections 4-11. See also [Baz](../examples/02-baz-reconciliation.md).
