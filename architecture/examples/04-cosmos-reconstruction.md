# Example 04. Onboarding the legacy application Cosmos

**Status:** walkthrough of the attached reconstruction proposal. Cosmos is hypothetical.

## 1. Begin with what exists

Cosmos initially has Implementation I0 and little or no explicit upstream design, behavior model, architecture or test harness. The system already contains knowledge, but some of it is accidental, undesirable or undocumented.

Saum could propose a governance shape rather than impose one. Perhaps Architecture -> Implementation is enough. Perhaps Design and a Test Harness merit separate authority. The question is what knowledge the human wants to govern, not how many documents an agent can produce.

## 2. Observe without declaring requirements

An analysis recovers components, dependencies, interfaces, data flow, schemas, deployment and observable behavior. It records uncertain interpretations and evidence references.

A characterization harness might establish that `DELETE /user` currently returns 500. An as-built design might record that AuthService directly accesses billing tables. Both may be accurate descriptions and still undesirable future requirements.

The inferred documents are observed/reconstructed or candidate state. They do not automatically acquire authority over future work.

## 3. Add intent that code cannot establish

The human explains Cosmos's purpose, required compatibility, quality priorities, business constraints, known bugs and debt. That input helps separate intentional behavior from accident.

For the two examples, the human might identify the 500 response as a bug and reject authentication's ownership of billing persistence. Or a behavior might need temporary preservation for compatibility. The branch does not prescribe the outcome; it requires a legitimate decision rather than inferring intent solely from code.

## 4. Settle a first upstream basis

Candidate Architecture A0 and Design D0 are reviewed. Characterization tests remain evidence, while acceptance obligations reflect the behavior actually chosen.

The accepted higher-level state can now constrain future implementation. The provenance still says it was reconstructed using I0. That historical evidence relationship does not reverse the normative dependency graph.

## 5. The pivot can reveal immediate debt

I0 may differ from the newly accepted design even though that design was reconstructed partly from I0. That is expected when review distinguishes as-built from intended state.

Implementation is asked to reconcile, not to be discarded. Useful existing mechanisms and compatibility knowledge are preserved while changes address the accepted obligations.

Tests derived from old behavior must not mechanically veto every intended correction. Their status as characterization versus acceptance matters to the review.

## 6. Enter ordinary reconciliation

Once the first authority graph is established, a later A0 -> A1 change can flow through Design/Harness/Implementation like any other governed project.

The special bootstrap posture is no longer the main way work proceeds, though further reconstruction may be needed in poorly understood areas. The source branch does not define partial pivots or mixed lifecycle mechanics; those remain questions.

## 7. What this example tests

Can the system distinguish evidence from authority? Can it preserve implementation knowledge without canonizing bugs? Can it invite human intent at the right point? Can accepted reconstruction produce an actionable discrepancy rather than a cosmetically tidy document? Can the graph be established without pretending the legacy system was designed top-down?

Those are the conceptual requirements illustrated by the branch, not a settled reverse-engineering algorithm.

## Source

[S04: legacy-project reconstruction](../history/sources/04-legacy-project-reconstruction.md). See [the reconstruction topic](../concepts/10-legacy-project-reconstruction.md) for the full proposal and open questions.
