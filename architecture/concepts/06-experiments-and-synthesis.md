# 06. Independent experiments, comparison and synthesis

**Status:** central motivation and working experimental principles; unsettled ownership and implementation scope.

## 1. Why experiments mattered originally

The original idea grew from fan-out around a shared frontier. Several workers could explore different architectures, interpretations and practical mechanisms. The useful result was not necessarily one complete winner.

An invalid branch may contain a useful observation technique. An unfinished branch may expose an upstream ambiguity. A locally successful mechanism can challenge another worker's claim that the environment makes progress impossible. Comparing independent attempts can therefore produce knowledge no individual worker has available.

This motivates experimentation as architectural search and specification pressure, rather than simply redundancy or a way to obtain more code.

## 2. Three purposes that should remain explicit

### Materialization

Find a good way to advance the current artifact. Outputs may contribute to the accepted result through continuation, reference, selective reuse or synthesis. The desired product is a useful candidate plus what was learned.

### Upstream validation

Use attempted downstream materialization to test upstream intent. A documentation project can write A/B/C, commission implementation and harness spikes, receive findings, improve the documents and only then extend them to D/E. The code can be deliberately disposable.

This complements textual review. It tests consequences that appear only when someone tries to construct a concrete system or an external observer. The workflow should make its purpose explicit so a validation probe does not silently become production development.

### Scoped challenge

Given a major change such as A13/Z, ask whether accepted local decisions X/Y still fit. Preserve unrelated accepted state initially, identify specific reopening proposals, and widen scope only when justified through review.

A mature project can perform a strong local challenge without declaring its entire foundation invalid. The previous system's evidence and failed alternatives remain useful inputs.

## 3. What independence means operationally

A comparable cohort should identify a common starting revision, authoritative inputs, scope/frontier, worker instructions, environment/tool policy, resource bounds, and intended differences such as model or approach. Sibling outputs should not be visible unless explicitly permitted.

This is experimental discipline, not a statistical claim that model outputs are independent. Common models, training, prompts and fixtures can still create correlated blind spots.

The original DAG model made information edges inspectable: a reviewer sees all candidates, while a synthesis worker may see only approved learnings and selected references. That remains useful even if a different runtime executes the work.

Changing inputs mid-run weakens interpretability. New lessons normally belong in the next controlled attempt, not covert live coaching of the current one.

## 4. Stop invalid work without steering architecture

A monitor can check whether a run has crossed an already settled, objective invalidation boundary. For example, a claim of complete coverage may fail a required public-boundary gate.

That differs from telling an experiment to abandon its architecture because a sibling found a better proxy. The latter contaminates search. Monitoring should distinguish invalidity from taste, unexplained expense, or an unexpected but permitted approach.

Automatic termination should be conservative. Ambiguous semantic concerns may justify review rather than an immediate hard stop. Suspiciously cheap progress can trigger inspection but is not itself a validity rule.

The v5 experience also exposed a driver problem: workers that reported ready could be dispatched again with unchanged inputs. The outer loop should be able to wait for review instead of encouraging scope expansion merely to use remaining budget.

## 5. Fixed frontier does not mean one worker iteration

A worker may use many internal iterations to complete a bounded outer objective. Stopping at the frontier keeps experiments comparable and prevents subsequent speculative work from contaminating the result.

The original sequential Scenario 01 -> Scenario 02 model was one example. Later work packages and adaptive increments loosen the assumption of linear progression while retaining useful boundaries.

Low-risk work can eventually be grouped. The important property is an interpretable scope and acceptance basis, not a universal one-scenario-per-run rule.

## 6. Where experimentation belongs: preserve the unresolved fork

The main discussion moved through three positions.

First, Zonotron owned the full experiment DAG: workers, monitors, reviewers, synthesis and gates. Second, narrowing Zonotron to cross-project reconciliation suggested putting experiments inside projects as a reusable capability. Third, the user proposed that Zonotron can branch a project and request multiple ordinary candidate reconciliations without knowing its method.

The third position preserves generic orchestration. Each isolated project instance can use Spec Kit, Ralph, plain agents or humans. Zonotron can compare the results against the same objective and retain provenance.

But this still expands its responsibilities: isolation, dispatch, comparison and resource use are workflow concerns. Whether multi-candidate reconciliation is core, optional infrastructure, or a service available to projects remains tentative. The dossier does not choose one silently.

Project-local experiments remain supported regardless. A project may already have a better domain-specific alternatives process; it should not be forced to replace it.

## 7. Scope of fan-out

The smallest useful scope may be one uncertain work package. Running all six consequences of a change three times because one decision is uncertain is potentially wasteful.

A full-project experiment can be justified when interactions between choices are what need testing. A speculative extension would branch an entire downstream subtree and inspect both local results and the integrated outcome. The user explicitly treated this as a high-cost possibility, not a priority.

The generic question is whether these are one operation applied at different scopes or substantially different capabilities with distinct authority and cost risks. That is still open.

## 8. Comparison should be causal, not a leaderboard

A review should assess actual validity first, then preserve qualified mechanisms and interpret differences. Relevant dimensions include required coverage, evidence strength, architectural coherence, complexity, introduced interfaces, unresolved assumptions, resource use and compatibility with upstream state.

It should not infer quality from code volume, apparent depth, iteration count or low token usage. A cheap incomplete branch is not evidence of a cheaper complete solution. An expensive run does not establish that one particular dependency caused the expense.

Comparative review can look for contradictions, valid divergent interpretations, hidden consequences, repeated workarounds, boundary violations, consequential architectural variance, suspicious apparent progress and repeated failure at the same frontier. These are salience signals, not automatic upstream change requests.

## 9. Knowledge first; code reuse is a policy

The original inventory preserved four distinct reuse policies:

| Policy | Meaning | Why it might be useful |
| --- | --- | --- |
| Learnings only | Fresh work from the common base plus distilled, approved knowledge. | Tests whether reusable understanding is sufficient; encourages a coherent result. |
| Reference | Learnings plus permission to inspect selected immutable code or other artifacts. | Avoids expensive reinvention without adopting the whole branch. |
| Transplant | Explicitly port a selected mechanism, commit or hunk, with adaptation/review. | Reuses genuinely valuable work while keeping adoption visible. |
| Continue | Use a coherent candidate as the next base after appropriate acceptance. | May be efficient for narrow changes on a mature foundation. |

No global policy was chosen. Comparing them is itself a possible experiment. The goal is neither always rebuild nor always pick a winner.

Synthesis should avoid an incoherent mixture of every attractive fragment. A fresh worker can be asked to produce one coherent architecture using the learned constraints, with references limited to what policy permits.

## 10. Guidance for the next cohort

The user clarified that reducing permanent Ralph prompt growth does not mean withholding useful opinions. Zonotron should be able to provide milestone-specific guidance: a known dangerous assumption, promising options, previous failed boundaries, or an uncertainty that remains open.

Such guidance is separate from stable worker protocol and durable accepted project decisions. It can expire after foundation settlement or be reflected in the repository. It must not smuggle an unresolved upstream obligation or rejected interface proposal into the next attempt as a mandatory rule.

The distinction also matters experimentally. To test whether manager learning helps, keep the base, worker protocol and authoritative objective fixed while varying only approved guidance.

## 11. Cost-aware exploration without suppressing discovery

The user proposed an initial analysis of a foundation milestone, identifying consequential unknowns before expensive broad attempts. Some questions may be resolved by human discussion, some by a cheap constrained representation, some by a focused probe, and some by open-ended materialization.

At the same time, Traust's difficult Scenario 01 may be valuable precisely because it pressures many assumptions. Over-specifying the fixture or runtime arrangement would suppress variation that exposed upstream weaknesses.

The design goal is not minimal exploration everywhere. It is to choose where unresolved uncertainty justifies its cost. Foundational exploration should eventually reduce repeated rediscovery; otherwise the system has not established its economic thesis.

## 12. Model selection and budgets

The original design made model, reasoning level and resource limits explicit per node. A cohort might deliberately vary models or duplicate a task at a different reasoning level. Cheap routine execution and expensive high-risk review need not use the same model.

Budget dimensions considered included token usage, iterations, wall time, frontier completion and objective invalidation. Iteration counts are poor proxies for equal compute. Usage totals do not automatically equal billed cost or uncached work.

The exact controls and whether the top-level application or delegated runner enforces them remain open. What survives is that workers should not silently widen scope or resource commitment outside the request's policy.

## Sources

[S01](../history/sources/01-original-manager.md), sections 6-16, 20-22, 28-36, 40-44 and 48-50; [S02](../history/sources/02-project-boundary-notes.md), project-level and subtree fan-out proposals; [v5 manager review](../history/evidence/review/MANAGER_REVIEW.md). Main-thread foundation/guidance discussion is M03 in [source map](../history/sources-and-coverage.md).
