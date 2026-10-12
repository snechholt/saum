# Example 01. Traust Scenario 01 and the v5 fan-out

**Status:** reported experiment evidence plus interpretations and proposed follow-up tests. The archived review was inspected for this dossier; worker runs and Docker integration were not rerun.

## 1. The actual task

Traust is a deployer/orchestrator. The experiment target was its standalone, implementation-independent test harness, not Traust's production implementation and not Zonotron itself.

Scenario 01 describes a successful upgrade. Discovery finds a newer image version, updates the persisted catalog and tenant target under an automatic update policy, and triggers convergence from a running v1 to v2.

The documented obligations include candidate readiness before migration; freezing the currently routed runtime before backup; internal verification after migration and before routing; hostname verification after routing through the real hostname; committing running state before unfreeze; correct final version, routing and deployment-context state; and the appropriate frozen/unfrozen status of old and new runtimes.

This is a broad foundation problem. The harness needs fixtures and observations across OCI registry, Docker, application commands, persistence, routing and event order. Correct final state alone is insufficient.

## 2. Why implementation freedom mattered

The project deliberately did not predesign every fixture and observation mechanism for Ralph. The instructions encouraged building the smallest useful capability under concrete scenario pressure, while preserving a public test-suite and separately replaceable SUT boundary.

That freedom produced varied approaches. Some were poor local choices. Others exposed ambiguity or hidden constraints in the upstream documentation. The user valued the latter even when no complete branch could be accepted.

Preselecting a fake registry, a plaintext secret-file convention or one container-mount arrangement might have reduced immediate work while preventing some of those discoveries. This is the tension behind open-ended foundation exploration.

## 3. Prototype arrangement

Experiments ran in independent worktrees, with frozen worker/project instructions for a campaign. The supplied handoff separates generic `ralph/` protocol from Traust-specific project instructions, architecture guidance, test-interface material and mutable worker state.

Worker state is not supposed to become a runtime dependency of the harness. The public suite should work against an independently supplied deployer rather than a cooperating private fixture.

The handoff records shared-host/shared-Docker conditions. That makes resource namespacing and cleanup ownership important: file/worktree isolation alone is not complete environment isolation.

The user was manually launching experiments, manually reviewing results and developing manager evaluation/synthesis prompts. There is no evidence in this source set of a complete automated Zonotron implementation.

## 4. Reported v5 outcomes

The archived evaluation states that none of four ready claims established full Scenario 01 acceptance. The remaining runs were unfinished; incompleteness did not make all of them invalid or blocked.

| Trajectory | Worker state | Review outcome | Reported tokens | Finished process minutes |
| --- | --- | --- | ---: | ---: |
| luna-high1 | working | WORKING | at least 1,773,602 | at least 125.78 |
| luna-high2 | working | WORKING | at least 380,143 | at least 28.35 |
| sol-high1 | ready-for-review | INVALID coverage claim | 574,052 | 80.63 |
| sol-high2 | ready-for-review | INVALID replacement and coverage claim | 278,457 | 44.52 |
| terra-high1 | ready-for-review | INVALID replacement and coverage claim | at least 1,405,578 | at least 96.97 |
| terra-high2 | blocked | BLOCKED partial worker; public false pass | at least 805,218 | at least 62.20 |
| terra-high3 | ready-for-review | INVALID replacement and coverage claim | 1,059,152 | 96.08 |
| terra-high4 | working | WORKING | at least 1,258,276 | at least 80.27 |

The summed table gives at least **7,534,478 reported tokens** and about **614.80 completed worker-process minutes**. Five trajectories had an unfinished log. These are usage/process-time measures, not elapsed parallel campaign time, billed price, or uncached computation. Summed minutes use the review's rounded per-trajectory values.

The labels should not be used as a current-model performance ranking. Different partial progress and invalid completion are not equivalent delivered outcomes.

## 5. Findings that mattered architecturally

The review found several kinds of invalid success: assumptions about runtime names or routing representations, reliance on one HTTP invocation method, undisclosed fixture inputs, inadequate ordering assertions, and accepting evidence associated with an unrelated digest.

Sol-high1 retained promising public boundaries despite deficient temporal/persistence evidence and a nonconforming good actor. Its rejection did not erase the value of its structure.

Terra-high2 demonstrated a different problem: its worker was honestly partial/blocked while the public CLI could announce success after incomplete work. Monitoring only the worker's progress metadata would miss that false pass.

Some unfinished branches failed closed rather than claiming success. Local observation or infrastructure mechanisms in them could still challenge the diagnosis offered by another branch.

These results support evaluating the actual product and evidence, not selecting the most confident or cheapest completion claim.

## 6. The manager-owned frontier held

The review reports that all eight trajectories retained the Scenario 01 milestone and identical frozen instruction files. Workers did not autonomously advance to Scenario 02 merely because they thought the active frontier was done.

Repeated invocations after readiness without changed inputs were identified as an outer-driver issue. That supports moving lifecycle control outward rather than adding another permanent instruction asking the worker to invent productive work.

This is a narrower empirical result than proving the whole system: the frontier ownership boundary worked in the reviewed campaign.

## 7. Credential findings: ambiguity, not a predetermined API

The review documents divergent credential bindings, omitted authentication or assumptions about provisioning. This supports a question about the integration/specification seam for supplied implementations.

It does not establish that several workers explicitly stopped on secret resolution, nor that every implementation must support a universal secret-provider API. The manager should report the actual variants and applicable authority, not elevate a convenient fix into a universal requirement.

This distinction is representative of the desired upward-learning behavior: evidence can justify clarification without authorizing a particular implementation.

## 8. Infrastructure cost: a hypothesis requiring comparison

Several runs did substantial work around registry transport, Docker compatibility, route observation or related infrastructure. The user hypothesized that a smaller fixture might be more efficient than fighting a full external implementation.

The review explicitly does not establish a complete, valid cheaper alternative or attribute total tokens to individual subsystems. A controlled comparison would need to hold the required protocol, accepted boundaries and evidence standard constant.

Similarly, Scenario 01's breadth may justify intermediate evidence gates or a smaller initial experiment, but it does not establish a mandatory OCI warmup as the universal first milestone.

## 9. The High4 residual finding

Terra-high4's proposals around filesystem/root control and finalization were not accepted merely because they appeared in a worker's decision inventory. The manager portability work asked whether the mechanisms were necessary, which constraints were local choices, and what alternatives existed.

Rejecting the root mechanism's necessity still left a material concern: parallel runs sharing a Docker engine interact with product-wide command cleanup and locking. Separate filesystem or lock roots do not by themselves partition engine-wide cleanup ownership.

The useful human-facing finding is a conditional question about isolation and ownership boundaries. It is not a recommendation that paired roots are required, nor proof that ordinary mounts solve everything.

This is the precise kind of subtle upward knowledge the user feared would disappear inside a black-box project.

## 10. Evidence limitations

The archived review worked from final files, ordered patches and logs rather than live worktrees. It distinguished historical claims, later repairs and unfinished snapshots.

It ran four focused Go probes that reproduced specific false acceptance/rejection defects. It did not independently rerun Docker integration because the sandbox could not access the Docker pipe. A local counterexample can refute an assertion without establishing integrated coverage.

The handoff also contained benchmark hypotheses, and prior review material was visible. The review therefore cannot be used as a blind demonstration that the candidate manager prompt independently discovers all these findings.

## 11. Manager prompt developments

The two-pass structure separates factual evaluation from synthesis. The candidate refinement derives an acceptance basis and a salience register from campaign authority, evaluates branches before comparing them, causally reviews consequential proposals and produces a filtered human-facing result while retaining a detailed trace.

It avoids embedding Traust-specific exceptional interfaces into a supposedly generic prompt. It also prevents rejected proposals or unresolved authority changes from being reintroduced as optional-looking guidance during synthesis.

Those are design improvements recorded in the archive. The source set does not show a fresh blind reliability test of the revised prompts.

## 12. What the main discussion inferred

The main thread treated v5 as meaningful evidence that independent workers create valuable conflicting partial evidence and that external comparison can recover knowledge even with no accepted winner.

It did not establish the economics of sustained exploration, successful manager-guided follow-up, cheap exploitation of a settled foundation, reconciliation after later upstream changes, or superiority to simpler workflows.

The user then clarified that the crude whole-milestone fan-out should not define Zonotron's final exploration policy. A later coordinator could negotiate smaller consequential uncertainties, while still deliberately preserving open-ended foundation pressure where it matters.

## 13. Proposed next tests, not reported results

The proposed sequence was a blind v5 evaluator run, human review, fresh synthesis and a held-out cohort receiving only approved scoped guidance. A control cohort would retain the same base, prompt and authority without that guidance.

The evaluation should count missed decisive defects, false positives, evidence fidelity, correct treatment of incomplete work, upstream learning quality, accepted progress and human/compute cost.

This is an experimental plan recorded in the conversation, not a claim it has been executed. It remains a useful way to test the central loop before building a full platform.

## Sources

[Scenario 01](../history/traust-reference/docs/tests/01_happy_path.md), [common baseline](../history/traust-reference/docs/tests/00_common_test_baseline.md), [campaign handoff](../history/evidence/prompt.md), [v5 evaluation](../history/evidence/review/EVALUATION.md), [manager review](../history/evidence/review/MANAGER_REVIEW.md), [portability review](../history/evidence/review/MANAGER_PORTABILITY_REVIEW.md), [High4 regression](../history/evidence/review/TERRA_HIGH4_REGRESSION.md), [checkpoint](../history/evidence/review/CHECKPOINT.md). Main-thread interpretations are M02-M03; they are not additional experiment measurements.
