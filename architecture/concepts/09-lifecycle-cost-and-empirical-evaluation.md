# 09. Lifecycle, economics and empirical evaluation

**Status:** useful lifecycle distinctions and a research agenda, not a proven efficiency model.

## 1. Lifecycle postures are local and reversible in scope

The original inventory described foundation discovery, incremental buildout and maintenance/reconciliation. The reconstruction branch adds a bootstrap posture for projects whose knowledge exists mainly in downstream artifacts.

These should not become one irreversible global state machine. Architecture can be mature while a new component is in foundation exploration. A mature project can run a scoped challenge without erasing all previous accepted knowledge.

The lifecycle posture influences how much exploration and review may be useful; it does not determine authority by itself.

## 2. Reconstruction/bootstrap

For a legacy system, the first task is not implementation from a clean spec. It is recovering enough explicit knowledge to establish a governed graph. Existing code supplies evidence, humans contribute intent, and reconstructed artifacts become normative only after review.

That differs from foundation discovery, where the question is how best to materialize already stated intent. The two can overlap in practice but should not silently share the same assumption about what is authoritative.

See [legacy reconstruction](10-legacy-project-reconstruction.md).

## 3. Foundation discovery

A foundation frontier is proposed by analyzing the actual backlog and important assumptions. It should be the smallest representative scope that can establish useful system grammar, not automatically the first numbered scenario or a predetermined login flow.

For the TODO example, a baseline might include a real HTTP server, persistence initialization, signup, login and relevant restart behavior while deferring TODO CRUD. For Traust, Scenario 01 through a legitimate public harness/SUT boundary was a candidate. A small amount of next-scenario pressure could test whether the baseline is immediately brittle.

The proposal should state why it is representative, what it pressures, what demonstrates success and what is deliberately deferred. It can be reviewed before large compute expenditure.

Foundation work may justify diverse attempts, stronger comparative review and more human gates. It is also where upstream loopholes and missing assumptions are most likely to surface.

## 4. Buildout and exploitation

Once useful structure is accepted, later work should normally exploit it. One coherent line of work can advance the backlog while experiments concentrate on consequential forks.

The repository itself may embody constraints more effectively than an ever-growing prompt. Correct runner/SUT boundaries, accepted test patterns and established lifecycle structure reduce the need to rediscover bootstrap decisions.

This is a hypothesis supported by the conceptual model, not a measured long-horizon result from the supplied v5 campaign.

## 5. Maintenance and scoped reopening

Upstream change should normally cause incremental reconciliation. Useful accepted state is preserved, and local experiments can challenge particular decisions whose assumptions changed.

The system should be able to say "reopen routing strategy; preserve public API and unrelated state model" rather than "reset the whole project to foundation discovery." Even a major rework starts with accumulated evidence, not amnesia.

What constitutes enough evidence to reopen a checkpoint or globally settled decision remains a policy question.

## 6. The economic thesis

The hoped-for payoff is not that parallel agents are intrinsically cheap. It is that paying for targeted exploration and comparative learning can avoid later rework and reduce repeated rediscovery.

The user's concern was that starting from nothing makes broad foundation work especially expensive. Initial frontier analysis, discussion-versus-experiment choices, scoped fan-out and cheaper behavioral projections might reduce that expense.

However, a faster locally obvious choice can hide specification weaknesses. The economic objective includes knowledge gained about upstream intent, not merely the immediate cost of producing code.

The system still needs evidence that the additional knowledge is worth its production and review cost. "Every branch taught us something" is not a sufficient cost-effectiveness argument by itself.

## 7. What v5 establishes, and what it does not

The supplied review records eight trajectories, four rejected ready claims, useful partial mechanisms, one public false success despite partial worker state, and preserved manager ownership of the Scenario 01 frontier. It supports the need for independent evaluation and the information value of comparing partial attempts.

It does not establish an accepted Scenario 01 foundation, a reliable model ranking, a blind successful execution of the revised manager prompts, or a cost-effective complete Zonotron loop.

Reported usage and finished worker process intervals were lower bounds where logs were incomplete. They were not monetary cost, elapsed campaign wall time or subsystem-level attribution. The review could not conclude that a small registry fixture would provide equivalent accepted coverage more cheaply.

The detailed evidence and numerical table are in [the Scenario 01 example](../examples/01-traust-scenario-01.md).

## 8. The prototype learning loop under test

The actual prototype workflow was manual experiment commands, manual review, manager evaluation/synthesis prompts, then potentially feeding reviewed learnings into another worker prompt. That is enough to test important semantics before building the full runtime or UI.

The proposed next proof was a blind manager evaluation on the frozen v5 archive. The evaluator should receive campaign authority and trajectories but not previous assessments, expected conclusions, regression reports or benchmark hypotheses.

Then reviewed synthesis would produce scoped guidance. A fresh worker cohort could be compared against an otherwise equivalent cohort without that guidance. The archive, original prompt and input authority should be held constant where feasible.

The main thread recommended not changing Ralph at the same time, because changing several variables would make attribution harder. No result of that held-out guidance comparison is provided in the source material.

## 9. What to assess in a blind manager run

The review should reject invalid success claims for the right reasons; distinguish honest incomplete work from blockage or invalidity; notice public false passes; preserve useful mechanisms from rejected branches; detect implementation-specific test coupling and temporal evidence gaps; and avoid turning hypotheses into recommendations with false authority.

It should also accurately cite evidence, compare sibling counterevidence, identify material residual upstream concerns and filter routine mistakes. False positives and attention cost matter alongside finding recall.

The manager candidates already aim to derive acceptance and salience from authority, perform causal review, retain a detailed trace, and produce a filtered human-facing result. Their design alignment is not the same as demonstrated reliability.

## 10. What to assess in the next worker cohort

Useful outcomes include fewer repeated invalid assumptions, fewer false completion claims, better portable observation, quicker arrival at useful structure, and lower cost per accepted advancement. Guidance should not merely force imitation of one previous branch or introduce new unjustified constraints.

No single improvement proves the entire thesis. A stronger long-horizon test would show that an accepted foundation supports cheaper later increments and that an upstream change can be reconciled without unnecessary rewriting of reviewed state.

## 11. Comparisons that would test the broader idea

The early critique proposed comparing a good single-agent workflow, a specification-oriented workflow and the settlement model. Important outcomes include architectural regressions after later changes, assumptions that survive review, rework from ambiguity, human review time, resource use per accepted frontier and the ability to explain decisions after several revisions.

Such comparisons should use equivalent objectives and count manager/reviewer work, not only worker tokens. The discussions do not provide results from them.

The strongest competing approach may be composition of existing tools plus modest metadata, rather than one named vendor. The question is which semantics remain genuinely difficult without Zonotron.

## 12. Possible maturity signals

The original inventory proposed decreasing foundational clarification requests, fewer invalid attempts, less consequential variance, more ordinary single-line progress, fewer reopened assumptions, and upstream changes fitting existing structures.

These are candidate signals, not automatic thresholds. Lower variance could mean a settled foundation, but it could also reflect shared guidance suppressing alternatives. Fewer reported findings could mean clarity, or a weaker reviewer. Metrics need interpretation rather than serving as unquestioned authority.

## 13. Bounded ambition

A successful prototype need not show universal autonomy. It can demonstrate a narrower loop: accepted objective -> negotiated increment -> candidate -> independent evidence review -> recorded acceptance -> revised plan or upstream finding.

That would test the cross-project interface while leaving execution methods replaceable. Rich scheduling, policy calibration and UI can follow evidence rather than precede it.

## Sources

[S01](../history/sources/01-original-manager.md), sections 6-7, 20-22, 27-29, 34-35, 40-44 and 48-50; [S03](../history/sources/03-validation-scoped-settlement-and-adaptive-planning.md), sections 3 and 12-13; [S04](../history/sources/04-legacy-project-reconstruction.md), lifecycle proposal; [v5 checkpoint](../history/evidence/review/CHECKPOINT.md) and [evaluation](../history/evidence/review/EVALUATION.md). Main-thread prototype and guidance discussions are M02-M03.
