# V5 manager review

The existing two-pass manager prototype has the right separation between factual evaluation and later synthesis. V5 requires a stronger evidence protocol within that structure. A successful manager should reject all four completion claims, retain useful partial mechanisms, identify the public false pass on Terra-high2, and surface environment, credential and ordering questions without automatically rewriting Ralph or advancing the milestone.

The review basis and detailed provenance are in [EVALUATION.md](EVALUATION.md). Candidate prompts are under [manager-candidates](manager-candidates/README.md). They are proposals for review, not adopted worker instructions or a new orchestrator implementation.

## What the manager should report

| Finding | Required manager treatment | Scope of conclusion |
| --- | --- | --- |
| Four ready claims fail full Scenario 01 acceptance | Reopen/reject each with its decisive evidence | Do not credit scenario depth or select the cheapest claim as a winner. |
| Sol-high1 has promising boundaries but incomplete assertions and a nonconforming good actor | Retain qualified mechanisms and record exact coverage gaps | A structurally promising branch is not complete. |
| Terra-high2 returns public success after discovery alone | Report a consumer false pass even though worker state is partial | Monitor observable success in addition to worker coverage declarations. |
| Runtime names, route files, HTTP exec commands and retained command containers are assumed | Audit whether permitted implementation variation causes rejection | Launch configuration cannot repair collector assumptions. |
| Route timing and persisted-finalization timing are unproved | Name the event actually observed and the order it establishes | Distinguish invocation, completion, route publication, and eventual state. |
| Several branches encounter similar Docker Desktop failures | Compare concrete diagnoses and locally successful mechanisms | A valid local mechanism can challenge a blocker even on an invalid full branch. |
| Credential provisioning differs or is bypassed | Propose an upstream/integration review question | Do not assume a universal secret-provider API is required. |
| Infrastructure work is costly and varied | Propose a bounded comparison of fixture alternatives | No complete valid cheaper alternative has yet been established. |
| Explicit route, root and finalization contracts appear | Audit deferral, necessity, alternatives, semantics and actual proof | Fresh context is a procedural gate, not proof of necessity. |
| Ready iterations continue without changed inputs | Recommend driver wait/stop behavior | The worker should not advance scope to consume the remaining budget. |

The manager should also state positive evidence: milestone ownership held, frozen instruction hashes agree, some workers fail closed, and several implementation-independent local techniques are available for further validation. It should not convert an invalid branch into a total loss of experimental knowledge.

## Assessment of the existing prototype

The [evaluation prompt](../v4-manager-prototype/PROMPT_EVALUATE.md) already asks for independent assessments before comparison, historical claim inspection, actual mechanical validity, provenance and repeated workarounds. The [synthesis prompt](../v4-manager-prototype/PROMPT_SYNTHESIZE.md) separates established knowledge from optional worker guidance and forbids changing authority through that guidance. Those principles should remain.

The earlier [v4 evaluation](../v4-manager-prototype/EVALUATION.md) already surfaced credential/authentication ambiguity. The mechanism therefore has some demonstrated ability to produce upstream findings. V5 should strengthen how evidence is gathered and checked, rather than add a list of predetermined OCI conclusions to the prompt.

The concrete weaknesses are:

1. **Campaign resolution assumes live worktrees.** The current campaign points to a v4 path and documentation inside each worktree. V5 is an archive with one shared documentation tree, patches and incomplete logs. An evaluator needs explicit source identities, available historical reconstruction, and final-snapshot provenance before applying the first-claim rule.
2. **Review states omit honest unfinished work.** `VALID/INVALID/BLOCKED/UNCERTAIN` invites labelling every incomplete run as blocked. V5 contains three working trajectories with actionable tasks. Add `WORKING`, and track public execution results separately from worker state.
3. **The trigger is too narrow.** Inspecting the first covered commit is useful, but Terra-high2 demonstrates a full public pass without a coverage claim. Inspect all claims of scenario success, including public output and zero exit status, and distinguish initial acceptance from later repairs.
4. **The boundary checklist needs an evidence trace.** Separate processes and documented API calls are insufficient. For every decisive assertion, trace the originating SUT action, product/runtime configuration, external observation, identity lookup, lifecycle availability and predicate. Check legal variation rather than matching one supplied actor.
5. **Positive and negative self-tests can share a blind spot.** Require a small set of counterexamples tied to consequential requirements, plus permitted-variation checks. Record a unit assertion probe as local evidence and an integrated independent-actor exercise as stronger evidence. Do not substitute one for the other.
6. **Temporal claims need explicit semantics.** Ask what was observed before what, and whether an action had already started before a checkpoint held it. Record sampling/race limitations. A matching final record after a late event cannot be promoted to an earlier-order guarantee.
7. **Invalid trajectories still contain systemic evidence.** Retain their friction, specification interpretations and locally successful boundaries after invalidating depth. Otherwise the evaluator discards precisely the data needed for manager-owned diagnosis.
8. **Environment, cost and stopping evidence need structured capture.** Distinguish engine failure, compatibility, addressing, daemon trust, filesystem notifications, interruptions and unchanged readiness. Record complete versus missing token/time footers and driver restarts. Do not infer efficiency from iteration count or code size.
9. **Consequential decision review needs authority-derived categories.** Derive the project's exceptional decisions and prescribed process from supplied authority. In v5 that identifies additional deployer test interfaces and their deferred challenge; neither concept belongs in the reusable procedure. Review pressure, local assumptions, necessity, ordinary/sibling alternatives, hidden consequences and residual upstream value.
10. **Human attention needs a causal filter.** Suppress avoidable or unsupported mechanism proposals after deeper review while retaining their evidence and disposition. Surface only a material residual concern or an explicit source-required disclosure. A proposal that fails necessity review can reveal a different useful restriction without becoming an approval request.

These changes refine the manager's work, not the worker's architecture. They do not require another long list of Ralph prohibitions.

The [reusable design review](MANAGER_PORTABILITY_REVIEW.md) maps project-specific acceptance assumptions to campaign authority and generic review operations. The proposed procedures derive an acceptance basis and salience register, then retain causal audit records separately from vetted human findings. Synthesis preserves that filter rather than reviving suppressed proposals through optional guidance.

The [Terra-high4 regression](TERRA_HIGH4_REGRESSION.md) demonstrates the resulting distinction. It does not forward the root or finalization mechanism for approval. The root review instead surfaces an isolation/ownership question: separate lock locations do not scope the product's command cleanup on a shared engine. This conditional consequence survives the ordinary mounted-state alternative and warrants clarification of the project/product operating boundary.

## Comparison with the benchmark hypotheses

### Credential resolution

The hypothesis is supported as an unresolved integration/specification question. The product names `secret_ref` or an equivalent source but does not define the provisioning seam for a supplied implementation. V5 paths range from ordinary launch-configured credentials to a presumed file convention, omitted auth and local image inspection.

The evidence does not support saying that multiple v5 workers explicitly stopped on secret resolution: their recorded questions primarily concern transport, Docker API access and timing. Nor does the evidence establish a need for universal provider semantics. The manager should report the actual variants, the relevant product sentence and the portability question, with confidence qualified accordingly.

### External infrastructure preference

The hypothesis is plausible and worth testing. Several workers add transport bridges, OCI publication code, socket compatibility layers or route mirrors around standard infrastructure. A tiny deterministic fixture might preserve the required protocol with less machinery.

A cost conclusion remains unproved. No full valid simple alternative is present, and the archive does not attribute tokens or time to individual subsystems. The manager should propose a controlled comparison of dependency choices while holding the documented protocol and independent SUT boundary constant. It should not replace the infrastructure preference merely because the existing runs were expensive.

### Milestone decomposition

Scenario 01 pressures registry discovery, configuration, persisted state, image commands, runtime readiness, routing, freeze and timing at once. Repeated friction can justify a proposal for intermediate evidence gates or a smaller approved foundation. It does not establish an OCI warmup as the correct universal onboarding milestone.

A useful next outer-layer decision is whether environment compatibility and oracle validity should be reviewed before another full milestone investment. Milestone scope remains externally owned; workers do not decide this from leftover budget.

## Proposed next workflow

1. Review the v5 assessments and correct disputed evidence before synthesis.
2. Decide which concrete completion defects to send back as scoped review findings, and which questions require upstream/contract review. Do not settle a baseline yet.
3. Validate selected local mechanisms without adopting their whole branches: normal configured app observations, label/state-based identities, current Traefik interoperability, external route mirroring, and canonical container state mounts.
4. If useful, authorize a bounded validation exercise for credential provisioning, temporal observation, or fixture cost. Specify the question and success evidence before dispatch; do not steer running experiments toward a sibling architecture.
5. Run synthesis in a fresh context from reviewed evaluation and selected immutable evidence. Keep proposed worker guidance informative and separate from mandatory corrections or new authority.

The manager candidates preserve the existing manual two-pass workflow. A deterministic application can later own manifests, information boundaries, budgets, review transitions and dispatch termination; this review does not implement that application.

## Reliability limits

This is an evidence-based assessment and a proposed improvement, not a blind replay of the existing manager. The input handoff contains benchmark hypotheses, and the reviewer later inspected the previous evaluation. Independent discovery cannot be claimed after those inputs are visible.

A future test can give a fresh evaluator only the candidate campaign, frozen requirements and archived trajectories, then score its factual findings against this reviewed report. That test should measure decisive gate detection, false positives, evidence accuracy, incomplete-state classification and systemic issue discovery. The benchmark conclusions should remain outside its inputs.

No v6 worker prompt, milestone, product specification or archived trajectory was modified. Adoption of the candidate manager protocol and any resulting upstream changes remains a separate settlement decision.
