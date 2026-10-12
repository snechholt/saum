# Traust Ralph test-harness context handoff — through v5 experiment

## What we are doing

We are designing an autonomous loop that builds and later maintains the standalone **Traust deployer test harness** from the Traust documentation and its documented test scenarios.

The goal is not merely to generate tests quickly. The loop should be able to work for a long time with limited supervision while preserving the behavior that worked well in the initial experiments:

- read the current scenario and the product documentation it depends on
- try to exercise that scenario as a correct deployer would
- discover the first concrete missing harness capability
- build the smallest useful capability that solves that need
- continue the exercise and let later requirements pressure the architecture naturally
- prove the scenario through realistic integrated boundaries before calling it covered
- then move to the next missing documented requirement

The important design objective is **autonomy without predesigning the harness for Ralph**. We want to constrain milestones and product boundaries, not prescribe the implementation.

The harness itself must remain implementation-independent: a future deployer written in any language should be testable against it.

## Repository split

`ralph/` is deliberately small and generic. It defines the reusable autonomous-worker protocol: how an iteration resumes from durable state, reconciles authoritative source changes, selects work, handles review/blockers, makes progress, and exits.

`test-harness/` is Traust-specific. It contains:

- `PROJECT.md` — project mission and harness-development workflow
- `PROJECT_PREFERENCES.md` — project-specific requirements/preferences
- `ARCHITECTURE.md` — human-authored architectural guidance/candidate techniques
- `specs/` — test-interface guidance and any concrete contracts
- `agent-state/` — mutable Ralph state/memory
- the actual harness implementation, fixtures, scenario coverage, etc. as they emerge

`agent-state/` is intentionally separate from the harness output. Runtime tests must not rely on it.

## How experiments are run

Experiments run in independent **Git worktrees**. Each run/model gets its own worktree/branch.

Once an experiment version starts, its Ralph/project instructions are **frozen** for that experiment. We do not tweak prompts mid-run. New lessons go into the next version. Do not spend time re-explaining this experimental discipline when discussing what to run next; it is already established.

Parallel worktrees share the same host/Docker daemon, so the harness must namespace its Docker resources and avoid fixed host ports where practical.

The current outer driver is SantanderAI/ralph, but the repository protocol is intentionally driver-agnostic. The driver owns scheduling/loop termination; repository instructions express work state rather than driver-specific stop mechanics.

## Core behavior we want from Ralph

### Scenario-driven development

The documentation is authoritative. Ralph discovers the scenario suite/order dynamically from the docs rather than maintaining a duplicate hardcoded roadmap.

When no valid active work is recorded, it compares documented requirements against actual harness coverage and persists the useful result in `agent-state/TODO.md`. The TODO is a cache, not a second source of truth.

On an obviously incomplete/empty harness, Ralph should not deeply analyze every future scenario. It should identify the documented progression, select the earliest actionable missing coverage, and analyze that scenario deeply. Future scenarios are context, not permission to implement their capabilities early.

### Build from concrete pressure

For the active scenario/invariant:

1. read the scenario, common baseline, referenced product docs, and materially relevant surrounding docs
2. identify the behavior that must be controlled and evidence that must be observed
3. attempt the scenario manually/with a stand-in as though it were the deployer
4. identify the first missing harness capability
5. implement the smallest reusable capability needed now
6. continue the exercise; add more only when concrete gaps appear
7. create known-good and representative known-bad exercises
8. map the meaningful requirements to mechanical assertions/coverage
9. persist useful design/capability/question state

Do not prebuild abstractions because later scenarios might need them. Later scenarios are allowed to invalidate earlier abstractions.

### Integrated coverage, not synthetic success

A scenario is not covered merely because validator/collector unit tests pass or normalized/synthetic evidence can be constructed.

Before advancing, its known-good self-exercise must go through the same kinds of real harness boundaries a conforming deployer run would use. If the scenario depends on Docker, routing, registry behavior, application commands, persisted state, etc., the self-exercise should exercise those boundaries rather than substituting in-memory evidence for them.

The stand-in is allowed to perform deployer actions during harness development. In the current v4 baseline the rule is explicit:

> During self-exercise, the stand-in may perform any documented action a deployer would perform, including producing documented persisted state. It must not directly populate collector output or assertion evidence.

The evidence must still arrive through a boundary available with the same semantics to an independent deployer.

### Prefer real/off-the-shelf infrastructure

Use real Docker. For standard external systems such as registries/routing/HTTP services, prefer ordinary off-the-shelf implementations when they already provide what the active scenario needs. Add proxies/wrappers/purpose-built fixtures only when the active requirement creates a concrete need.

Do not reimplement an external protocol merely for self-containment.

### Test-only deployer APIs are exceptional

Ralph has wide freedom inside the harness, but imposing a new required API on every deployer implementation is a substantial product-boundary decision.

Current rule:

- a new deployer test-only interface is a last resort; convenience/determinism alone is not justification
- when Ralph first thinks one is necessary, it must **not add it in that iteration**
- it persists the question as the next active task and exits cleanly
- the next fresh iteration focuses specifically on finding an implementation-independent alternative: documented APIs, persisted/external state, real side effects, logs, Docker/runtime/environment interposition, etc.
- only if that focused attempt still finds no reasonable alternative may it establish the smallest test-only contract
- the contract records the originating requirement, alternatives considered, and why they were insufficient

The contract inventory is also a human review surface. A test-only capability may later be rejected or promoted into the normal documented Traust API.

### Substantial iterations

Context loading is expensive. Small changes/commits are review units, **not iteration boundaries**.

Once useful context is loaded, Ralph should make a meaningful coherent body of related progress before exiting. Finishing one task, review finding, commit, or scenario sub-step is not by itself a reason to stop.

At the same time, do not broaden scope merely to keep busy.

### Persistent state and maintenance mode

This is intended to be both the initial harness builder and the future harness maintainer; there should not be separate init/update agents.

`agent-state/STATE.md` records `last_reconciled_master`. `master` is the authoritative docs source.

When master advances, Ralph compares the recorded cursor to current master, understands the relevant diff, updates affected work/coverage/state, and advances the cursor once the consequences of that source change are represented durably. It does not maintain a duplicate prose changelog of Git history.

This design should eventually transfer naturally from a Ralph loop to longer-lived/background autonomous agents: durable work state belongs in the repository, not in conversational memory or the loop driver.

### Avoid pointless loops

An iteration must make material progress, learn something that changes durable state, or establish that work is currently blocked/up to date. Reconsidering the same unchanged blocker is not progress.

A local blocker does not necessarily block independent work. But Ralph should not manufacture activity merely because the outer loop can keep running.

### Project-specific requirements

Current Traust preferences include:

- substantial reusable executable harness logic is Go
- shell/container setup/etc. can remain glue
- do not introduce Python for substantive harness logic
- test-owned stable names use the `traust-test-` prefix
- independently running worktrees/harness runs must not collide on Docker resources
- Docker is integral; if the Docker engine/daemon itself is unavailable/unhealthy, record the environment blocker and end the iteration immediately rather than searching for alternatives/workarounds
- avoid transient host-native network-listening executables where practical because Windows firewall repeatedly prompts for newly built `.exe` files; prefer Docker network services or a stable host binary path/name

## Experiment history

The numbered versions below are the actual current experiment series.

### v1

**Runs:** 2× Luna, 2× Terra, 2× Sol.

**Baseline:** initial Ralph/test-harness protocol introduced by `0001-Init-ralph-loop-files.patch`.

The initial baseline already had most of the desired philosophy: generic `ralph/` protocol, Traust-specific project/state, dynamic scenario discovery, persisted master reconciliation, substantial iterations, scenario-driven capability growth, real Docker, self-exercise, known-good/known-bad testing, and implementation-independent evidence.

**What v1 taught us:**

- The persistent-state/resumption model was useful. Agents generally made substantial progress rather than rebuilding everything each iteration.
- The milestone for declaring a scenario covered was too loose. Some trajectories built normalized/synthetic validators and controlled collectors, then advanced even though the scenario had not yet been exercised through a sufficiently real integrated environment.
- Some agents read/analyzed too much of the future suite and effectively reconstructed a roadmap/TODO for scenarios they had not reached yet.
- Different agents chose very different infrastructure approaches. One Luna trajectory spent substantial effort implementing registry behavior itself; a Sol trajectory used an off-the-shelf registry and moved faster. This reinforced the preference for standard real dependencies until a scenario proves a special boundary is needed.
- Parallel runs exposed Docker resource collisions because the worktrees shared one Docker daemon.

These findings drove v2: integrated self-exercise became an explicit coverage gate, TODO/future-scenario analysis was constrained, ordinary external infrastructure was preferred, project guidance was made immutable during runs, and per-run Docker isolation became a requirement.

### v2

**Runs:** 2× Luna, 2× Terra, 2× Sol.

**Changes from v1:** `0002-v2.patch`.

Key changes:

- a unit test that directly constructs normalized evidence is explicitly not enough to count as scenario self-exercise
- coverage requires integrated known-good/known-bad exercise through real harness boundaries
- TODO inventory work is proportional to uncertainty; no deep pre-analysis of unopened future scenarios
- independent worktrees must namespace Docker resources / avoid fixed host ports
- ordinary off-the-shelf external systems are preferred over reimplementing protocols
- human-authored project guidance is not modified by ordinary Ralph runs; run-specific decisions go into `agent-state/`

**What v2 taught us:**

The main correction worked. Sol and Terra stopped doing the earlier broad validator-first progression and generally advanced only after completing integrated scenario work.

State persistence and substantial iterations continued to work well.

Scenario 02 was particularly useful because different agents solved the same persistence-before-candidate boundary differently:

- one Sol trajectory correlated persisted-state publication with actual Docker candidate-creation events
- another Sol trajectory introduced a gate/proxy at the Docker API boundary so candidate creation could be physically blocked while state was inspected
- Terra trajectories were more willing to introduce explicit deployer test-only lifecycle checkpoints

This was a good result in one sense: the protocol constrained **what had to be proved** without prescribing **how**. But it also highlighted that Terra could reach for new required deployer test APIs sooner than desired. This led to the hard fresh-context test-interface challenge added in v3.

The Sol/Terra runs generally progressed several scenarios with real integrated exercises rather than the v1-style shallow coverage.

#### Luna follow-up

The Luna trajectories revealed a separate ambiguity in the self-exercise wording. With more runway, the two replicas interpreted the stand-in boundary differently but both incorrectly:

- one stayed blocked on Scenario 01 waiting for a real deployer
- the other treated that as an external blocker and moved on, accumulating partial later scenarios

Both were effectively treating documented persisted-state writes as something the stand-in could not do. That is why v4 now explicitly says the stand-in may perform any documented deployer action, including producing documented persisted state, while still forbidding it from directly populating collector/assertion evidence.

This should not be treated merely as a Luna-specific guard: two independent Luna trajectories exposed the same protocol ambiguity from different directions.

Luna remains interesting because it is much cheaper and its failure mode was conservative rather than nonsensical, but the data so far does not establish that it is a good lead model for this complex autonomous task.

### v3

**Runs:** 4× Terra, 2× Sol, 2× Astra (one high, one medium).

**Changes from v2:** `0003-v3.patch`.

Key changes:

- a new deployer test-only interface is explicitly a last resort
- discovery of a possible new test-only API forces an iteration boundary
- the next fresh iteration must focus on avoiding that API before it may be established
- Docker is explicitly integral; if the Docker engine itself is unavailable/unhealthy, Ralph stops rather than trying to find another path

**Status:** completed and reviewed.

Initial review found:

- the fresh-context test-interface gate worked well: no v3 trajectory established a deployer test-only contract, and at least one Terra trajectory explicitly carried a possible contract question across an iteration boundary before finding a documented/external alternative
- Sol continued to find implementation-independent Docker/runtime techniques, so the gate did not appear to suppress its architectural creativity
- integrated real-boundary exercises persisted several scenarios deep in the strongest trajectories
- Terra-high3 fell into a repeated false blocker around needing a real external deployer; this independently reinforced the stand-in ambiguity previously exposed by Luna v2
- Astra-high remained markedly more meticulous/slower and exercised substantial real Docker/Traefik/registry behavior, negative cases, races/isolation issues, and external process structure
- Astra-medium was affected by Docker instability and obeyed the hard Docker-stop rule rather than searching for workarounds

A later static product-boundary review materially changed the interpretation of v3 throughput. Some apparently fast trajectories, notably Sol-high1 and Terra-high2, embedded deployer-owned behavior inside scenario/test code. Their high scenario counts therefore did not represent the same usable external test-suite product as trajectories that had already established a separately replaceable SUT boundary. This discovery is the main reason revised v4 became a substantial correction rather than a minor clarification.

### v4

**Current baseline:** revised `test-harness-v4-update` on top of v3.

v4 became a materially larger correction than first planned. Static review of the v3 outputs exposed that the existing integrated-coverage gate still allowed a false interpretation of “covered”: some fast trajectories successfully self-exercised realistic infrastructure while embedding deployer-owned behavior directly inside scenario/test code. They therefore proved the fixture/verifier, but did not yet constitute a replaceable test suite for an independently implemented deployer.

The key distinction found during review was:

- Astra-high and Terra-high1 had established a genuine suite↔SUT boundary: the harness launches/connects to a separately running deployer/stand-in and drives it through documented Traust interfaces.
- Sol-high1 and Terra-high2 progressed much further in raw scenario count, but substantial parts of their self-exercise had scenario/test code itself performing deployer behavior. Their reported scenario coverage was therefore too generous under the intended product definition.
- This also explains part of the apparent throughput gap: Astra-high and Terra-high1 paid earlier for a real consumer-facing SUT boundary that some faster branches skipped.

This was judged to be a prompt-level ambiguity, not merely a model quirk. The previous wording strongly constrained evidence collection, but did not clearly require the stand-in itself to sit behind the same replaceable boundary that a future independent deployer will use. The instruction to “attempt the scenario manually/with a stand-in as though it were the deployer” could even encourage embedding deployer behavior in the test.

We explicitly do **not** require a real/reference deployer during harness development. Ralph remains test-first. It may build tiny scripted known-good/known-bad stand-ins containing only the behavior required by the active scenario, but those stand-ins must run separately and occupy the SUT side of the same boundary as a future real deployer.

The revised v4 therefore requires:

- a stable consumer-facing suite entrypoint approximately `traust-test-suite run --scenario <scenario> --sut-config <path> [--work-root <path>]`
- normal pass/fail semantics and a failing exit status when the scenario fails
- a scenario may be marked covered only when that public suite path exercises a **separately running SUT** through documented interfaces/external effects
- scenario/harness code is the **tester, not the deployer**: it may establish preconditions, launch/connect to the SUT, invoke documented interfaces, inject environmental failures, observe effects, and assert results, but must not perform behavior assigned to the deployer
- known-good and deliberately bad stand-ins must use the same normal SUT boundary
- the stand-in may perform any documented deployer action, including persisted-state writes, but only from the SUT side; it must never directly manufacture collector/assertion evidence
- replacing the stand-in with a conforming independent implementation must require only SUT launch/connection configuration, not scenario-code changes
- no private/self-test runner may bypass the public interface eventually used against real implementations
- host-side runtime files created by a run must live beneath one explicit per-run work root
- the earlier v4 firewall/listener guidance remains: avoid transient host-native listeners where practical
- Docker remains integral and the Docker-engine hard-stop remains in force

`ARCHITECTURE.md`, `PROJECT.md`, and `PROJECT_PREFERENCES.md` were all updated because the old wording in each contributed to the bad interpretation.

The first acceptance question for v4 is therefore **not throughput**. It is:

> When a trajectory first marks Scenario 01 covered, can the public test-suite entrypoint exercise a separately supplied implementation by changing only SUT launch/connection configuration?

Only trajectories that satisfy that product boundary should be compared on scenario depth.

The completed v4 result archive was subsequently reviewed with the first outer-manager/evaluator prototype. That review materially changed the interpretation again: only `sol-high1` passed the full Scenario 01 boundary mechanically, while several other runs still relied on private stand-in cooperation, synthetic discovery/state shortcuts, or non-portable observation paths despite using separate processes. The recurring lesson was that “separate process” was still being mistaken for “independently replaceable.”

During this v4 review we also realized that repeated parallel experiments need an **outer manager**: something outside Ralph that evaluates trajectories, compares cross-run evidence, surfaces systemic issues, and governs what happens next. This is a new idea from the v4 iteration, not mature architecture. The current manager ideas and proposed workflow are maintained separately in `manager.md`; read that file for manager design context rather than duplicating it here.

### v5

**Changes from v4:** targeted correction of the remaining SUT-boundary loophole plus a mock of manager-owned milestone scope.

v5 intentionally did not roll every v4 manager finding into project guidance. It changed two main things:

- **No harness-private SUT cooperation.** Process separation is explicitly insufficient. Inputs crossing into the SUT must be documented Traust interfaces/configuration/runtime inputs, ordinary external dependencies available to a real deployment, or an explicitly established test-only contract. Launch configuration may not smuggle private evidence paths, synchronization gates, expected outcomes, synthetic discovery values, duplicate-state destinations, fixture-specific callbacks, or equivalent hidden harness semantics into the SUT.
- **Externally owned milestone/frontier.** `test-harness/MILESTONES.md` is manager/user-owned and read-only to Ralph. For the v5 experiment it names Scenario 01 as the active milestone. Ralph may decompose and work that milestone, but it must not mark the milestone accepted or begin the next one. When it believes no actionable work remains, it records `ready-for-review`; manager evaluation decides whether the milestone is actually accepted, reopened, or advanced.

The current SantanderAI Ralph runner does not understand `ready-for-review`, so a logically finished worker may continue receiving fresh iterations until manually stopped. This is a driver limitation, not permission for Ralph to advance the milestone.

The v5 batch also exposed an important new reason to test the manager rather than immediately patching Ralph again. Multiple trajectories spent substantial effort around authenticated OCI discovery and related setup. Manual inspection suggests at least two candidate systemic issues:

- the Traust documentation uses `registries[].auth.secret_ref` and says registry auth is resolved by the deployer, but does not clearly define implementation-independent secret-resolution semantics for a standalone test environment;
- existing project guidance strongly prefers off-the-shelf external infrastructure and warns against reimplementing protocols, which may be over-rigid when a tiny deterministic protocol-faithful fixture would be much simpler than production-grade infrastructure.

These are **benchmark hypotheses for the manager**, not instructions to hardcode into the next Ralph version. A key purpose of the next context is to see whether analysis of the final v5 trajectories independently surfaces these kinds of architectural/prompt issues, and whether the manager can turn repeated experiment friction into useful human-review findings. Do not jump directly to v6 merely because these manual hypotheses are already visible.

A related open question is whether Scenario 01 is too broad as the very first settled milestone and whether experiment evidence should cause the outer layer to revise/decompose milestones. Do not assume an OCI warmup milestone a priori: that would encode hindsight into onboarding. The interesting question is whether the manager can notice repeated disproportionate friction and decide that a settled implementation, guidance change, upstream clarification, or milestone adjustment is warranted.

## Model observations so far

These remain provisional behavioral observations, not permanent model rankings. The v3 throughput comparison must now be interpreted in light of the corrected SUT-boundary requirement.

### Sol

Sol-high1 appeared to be the strongest v3 throughput trajectory, reaching Scenario 07 in three large iterations. However, static review showed that substantial deployer behavior was embedded in scenario/test code rather than living behind a separately replaceable SUT boundary. Its raw scenario count therefore overstates completed product coverage under the revised definition.

Sol-high2 progressed much less far and used substantially less compute: roughly 389k logged tokens / 49 minutes versus Sol-high1 at roughly 662k / 82 minutes. Thus “same number of iterations” was not equal work. Sol-high2 nevertheless discovered a useful narrow Docker candidate-create interception technique.

Sol remains architecturally creative, but v4 must establish whether that creativity still yields high throughput once the suite↔SUT boundary is mandatory.

### Terra

Terra showed the largest v3 variance. Excluding the blocked trajectory, runs used roughly similar total token budgets yet ended at materially different scenario depths. This reflects real architectural/path dependence rather than iteration count alone.

Terra-high1 was slower partly because it established a real external SUT/process boundary early. Terra-high2 moved much faster but, like Sol-high1, embedded too much deployer behavior in the test path and therefore solved an easier problem than intended. Terra-high4 took a more conservative fixture path and also suffered some Docker/environment friction.

Terra-high3 repeatedly treated the lack of a real external deployer as a blocker despite having a working Scenario 01 self-exercise. This reinforced the stand-in ambiguity already seen in Luna v2 and motivated the earlier v4 clarification that a stand-in may perform documented deployer actions.

Terra is therefore especially informative for v4: the experiment should show whether it can retain the economy of its fast trajectories while respecting a genuine replaceable SUT boundary.

### Luna

Luna exposed the stand-in ambiguity most clearly in v2: one trajectory stayed blocked waiting for a real deployer and another moved on while leaving partial scenarios. The revised v4 wording should directly remove that interpretation. Luna remains interesting as a cheap model, especially for a later reasoning-effort experiment, but should first be judged on whether it now understands the product boundary correctly.

### Astra

Astra-high was much slower in raw scenario progression, but later review showed that some of that apparent slowness was legitimate investment in a separately running SUT boundary that faster trajectories skipped. It remains unusually meticulous/integration-focused and is useful as a quality/reference trajectory. The v4 review should determine how much of its extra work was necessary product structure versus overinvestment.

### Iterations are not a compute unit

Do not compare trajectories primarily by iteration count. v3 demonstrated that an iteration may vary enormously in duration and token use. Future reviews should record at least:

- valid scenario depth under the public SUT-boundary definition
- logged tokens
- wall time
- correctness/coverage quality
- architectural complexity/debt
- test-only interfaces and blockers
- whether an iteration ended after a genuinely substantial coherent body of work

## What the next context should do

1. Read this handoff first. Treat the v1-v4 experiment history above as fixed context rather than reconstructing it from older conversations.
2. Read `manager.md` for the separately maintained outer-manager design and goals. This handoff intentionally does not duplicate that design.
3. Inspect the final v5 prompt/project files and the final v5 experiment archive. If partial snapshots from this context differ from the final archive, the final archive is authoritative.
4. Analyze the v5 trajectories **before** designing v6. Establish what each worker actually achieved, where it stopped, what consumed disproportionate effort, and whether `ready-for-review` claims satisfy the v5 product boundary.
5. Compare runs only after judging them independently. Look for repeated failure/misunderstanding, repeated blockers/workarounds, common successful techniques, suspiciously expensive or cheap paths, genuine architectural variance, and issues in upstream documentation/project guidance rather than treating every odd model choice as a prompt bug.
6. From that evidence, derive **what a good outer manager should have reported**. The focus is now not merely “what tweak should we make?” but whether the manager can replace the manual cross-run analysis we keep performing.
7. Only after the independent analysis, compare the manager-relevant conclusions with the manual benchmark hypotheses recorded in the v5 section above (including `secret_ref` semantics and potentially over-rigid external-infrastructure guidance). Missing those signals is evidence about the manager mechanism.
8. Evaluate the current manager prompts/workflow against that target: would they reliably surface and classify the important v5 findings without being spoon-fed the answers? Improve the manager if needed.
9. Do **not** immediately encode OCI as a new warmup milestone or otherwise bake hindsight into Ralph/onboarding. If milestone decomposition or a settled fixture approach is warranted, it should emerge as a manager conclusion from experiment evidence and then go through the appropriate human/outer-layer decision.
10. Keep worker prompt changes separate from manager findings. The point of this next stage is to test the manager as a stand-in for our own review process, not to skip that process because we already have intuitions about v6.

## High-level success criterion

The desired end state is one autonomous Traust test-harness designer/maintainer protocol that can:

- build the harness from scratch from the documented scenario progression
- later reconcile existing coverage against changes on `master`
- autonomously choose reasonable harness architecture under concrete scenario pressure
- avoid speculative future design
- prove coverage through realistic implementation-independent boundaries
- continue through local problems without unnecessary human stops
- surface genuine product/specification/API decisions cleanly for human review
- avoid infinite/repetitive activity when nothing changed

The point is not to make Ralph perfectly generic or perfectly deterministic. It is to preserve enough freedom that it can discover good solutions while making the few important workflow/product-boundary constraints reliable across different capable models.
