# Example 02. A large Baz change through a small project graph

**Status:** illustrative scenario used to pressure the conceptual boundary. Revision names and work scopes are hypothetical.

## 1. The realistic starting graph

```text
                    Architecture A12
                     /            \
                    v              v
             Test Harness TH4    Design D8
                     \            /
                      v          v
                    Implementation I27
```

There is one Implementation project, not a separate project for every feature, API or component. This was an important user correction to earlier examples.

Architecture A13 adds a full module Baz. It describes connections with Foo and Bar. Foo must also link into Baz; timeline entries must represent it; exports must include its data; APIs must support it; and tenant-specific templates must make it quick to configure.

The source example is intentionally broad enough to exceed one comfortable worker task or feature spec.

## 2. Global impact first, without pretending to know all internals

Zonotron knows that Design and Harness were accepted against A12 and must be considered against A13. It can identify the broad changed obligations and request local impact analysis or intake negotiation.

Design understands its component and interaction model. Harness understands scenario coverage and externally observable evidence. Each can explain what changes locally, which assumptions are uncertain and what dependencies need decisions first.

The global plan should be informed by both reports. Zonotron should not simply send the entire architecture and require every project to reinvent the same cross-project analysis, nor predesign all local tasks itself.

## 3. Parallel upstream reconciliations

Design may propose increments for Baz responsibilities, Foo/Bar interactions and interface effects. Harness may propose coverage for lifecycle, API behavior, template behavior and export/timeline consequences.

These are candidate settlement boundaries, not dictated implementation tasks. Either project may need further experimentation or clarification. Their internal methods can differ.

If Harness discovers that the required observation is ambiguous, it can block on that scope or return a finding. If an apparently valid local approach reveals a hidden architectural consequence, a reviewer may raise it even without a blocker.

## 4. Negotiation rather than one-shot task sizing

For either project, Zonotron can ask which units can be combined, which must be sequential, what evidence establishes them, which obligations remain unplanned and which decisions affect the other project.

A plan can initially cover only the next meaningful increment while retaining the full objective. Later findings may split or combine remaining work. The earlier plan and reasons for change remain evidence.

There is no assumption that "Spec Kit can consume this" or "one PR is the right unit" determines the architecture. The project's method influences viable intake; it does not redefine the objective.

## 5. Internal progress need not release downstream work

Suppose Design has accepted one component increment and Harness has accepted one scenario. Under the conservative branch proposal, that does not automatically make Implementation ready to consume them as globally settled inputs.

The projects may still revise foundations as the rest of Baz is addressed. Waiting for completed, integrated reconciliation reduces speculative downstream rework.

Partial scoped consumption remains a possible later policy. This example does not settle it, and a checked plan item should not imply it automatically.

## 6. Implementation eventually receives a large batch

Once D9 and TH5 are accepted against the compatible A13 basis, I27 can reconcile. It may still face a substantial combination of persistence, API, UI, export, timeline and template work.

Zonotron and Implementation negotiate a plan suited to the actual local codebase and methodology. The project may use a specification workflow per work package, a task orchestrator, plain coding sessions or another method.

The problem being solved is not just feature decomposition. It is satisfying several governing inputs while preserving accepted structure and accounting for the entire upstream change.

## 7. Experiments can be scoped to the uncertain part

Perhaps most changes fit the existing design, but how tenant templates interact with Baz admits several consequential approaches. Zonotron or the project could request alternatives only for that package.

A full-project fan-out would be justified only if interactions among the changes are the uncertainty. It should not be the default merely because the total objective is large.

After comparison, the result may be an accepted local decision, a guidance update, a new experiment, or an upstream question. No candidate automatically authorizes its own architectural interpretation.

## 8. Final validation and preserved knowledge

Implementation's evidence must establish the new obligations, preserve earlier accepted behavior and satisfy integrated constraints. Passing its own task checklist is not sufficient. Passing the Harness is important but may not alone establish conformance to Design or broader architecture.

I28 is accepted under the relevant validation and approval policy. Its basis and provenance should be inspectable. The old accepted configuration remains historical evidence rather than being overwritten conceptually.

## 9. Questions the example leaves open

What happens if A14 arrives halfway through the work? How are cross-project plan conflicts resolved? Which context revisions does Implementation read? Can a no-change artifact be revalidated? What if Design and Harness accept incompatible interpretations despite sharing A13? How should the global plan represent that conflict without making the coordinator omniscient?

These are deliberate architecture questions, not missing implementation steps filled in by assumption.

## Sources

[S02](../history/sources/02-project-boundary-notes.md), particularly the realistic project graph and large batched implementation request; [S03](../history/sources/03-validation-scoped-settlement-and-adaptive-planning.md), adaptive planning and conservative release distinction; main-thread M05-M07.
