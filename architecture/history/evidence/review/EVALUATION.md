# V5 experiment evaluation

None of the four `ready-for-review` claims establishes complete Scenario 01 coverage under the v5 acceptance gate. The decisive defects include assumptions about SUT identifiers and invocation methods, undisclosed application wiring, and assertions that accept prohibited ordering. The other four trajectories remain incomplete; their partial work and recorded friction still provide useful evidence.

The manager-owned frontier held. All eight final trajectories retain the same Scenario 01 milestone and identical frozen instruction files. The campaign supports keeping independent evaluation outside the worker, but does not establish a foundation ready to settle or a reliable model ranking.

## Scope and evidence

The authoritative inputs are the final archive `v5_20261003_234931`, product documentation at `11c2a1efefcd5bf91d6bed3e6eff070e001211a6`, and the frozen v5 worker/project instructions. `manager-ideas.md` supplies design context; the v4 manager prototype supplies the review workflow.

`INVALID` below means a completion claim fails the full milestone gate. It does not mean every component on that branch is unusable. `WORKING` means no completion claim exists and useful obligations remain. `BLOCKED` records an explicit worker stop, with its cause evaluated separately. No trajectory is credited with a fully accepted scenario.

The archive contains final files, ordered commit patches, and worker logs, rather than live Git worktrees. First-claim commits come from patch headers and coverage/state changes. Three ready claims occur at the final patch. Sol-high2 claims coverage before two subsequent patches; those patches do not repair its decisive identifier and route assumptions. Final partial snapshots are evaluated as partial snapshots, without assigning uncommitted work to an earlier commit.

Local review probes ran in copies under `validation/`; archived worker files were not modified. Worker logs support the reported integrated exercises. This review did not independently rerun Docker integration: the review sandbox could not access the Docker pipe. Synthetic assertion probes below establish specific defects, not integrated coverage.

## Campaign measurements

| Trajectory | Worker state | Review outcome | Logged tokens | Finished log minutes | Go lines including tests |
| --- | --- | --- | ---: | ---: | ---: |
| luna-high1 | working | WORKING | at least 1,773,602 | at least 125.78 | 1,920 |
| luna-high2 | working | WORKING | at least 380,143 | at least 28.35 | 562 |
| sol-high1 | ready-for-review | INVALID coverage claim | 574,052 | 80.63 | 3,914 |
| sol-high2 | ready-for-review | INVALID replacement and coverage claim | 278,457 | 44.52 | 1,637 |
| terra-high1 | ready-for-review | INVALID replacement and coverage claim | at least 1,405,578 | at least 96.97 | 2,903 |
| terra-high2 | blocked | BLOCKED partial worker; public false pass | at least 805,218 | at least 62.20 | 2,591 |
| terra-high3 | ready-for-review | INVALID replacement and coverage claim | 1,059,152 | 96.08 | 2,254 |
| terra-high4 | working | WORKING | at least 1,258,276 | at least 80.27 | 4,128 |

Measurements are reproducible with [measure-archive.ps1](measure-archive.ps1); [archive-inventory.json](archive-inventory.json) retains individual log intervals, token counts, models, patch commits, and missing footers. Five trajectories have one unfinished log. Minutes sum completed process intervals and exclude unfinished intervals and gaps between invocations. Driver restarts reuse iteration numbers, so each distinct log is counted. Token totals are reported usage, not prices or an estimate of uncached computation. The logs identify GPT-5.6 variants; labels do not establish current-model performance.

## Experiment assessments

### Luna high1

**Outcome: WORKING. No covered or ready claim.** The final committed patch is `d37f6ca5dc9248ac27eed37932b71707f880ea1c`.

The branch builds a Go runner, isolated networks and run roots, a real registry, distinct fixture images, canonical persisted preconditions, and an application with readiness, verification and command behavior. API helpers for discovery and convergence exist but are not wired into the runner. The public path explicitly fails with “not executable yet”; all scenario requirements remain uncovered. There is no full external stand-in exercise and no test-only contract.

Recorded blockers concern a registry address/protocol usable by both a container SUT and the Docker daemon, and Traefik 3.1 interoperability with the engine. The worker records an actual old-client API rejection, rather than only an unexplained route timeout. These are integration issues with a healthy daemon. Significant completed infrastructure work coexists with repeated partial state and further driver invocations; the token total alone does not attribute that expenditure to a particular subsystem.

Evidence: [runner](../v5_20261003_234931/luna-high1/test-harness/cmd/traust-test-suite/main.go), [questions](../v5_20261003_234931/luna-high1/test-harness/agent-state/OPEN_QUESTIONS.md), [TODO](../v5_20261003_234931/luna-high1/test-harness/agent-state/TODO.md), and patch 0004.

### Luna high2

**Outcome: WORKING. No covered or ready claim.** The final patch is `eefefcd646a88dec38d75484167832226eb2608d`.

The public CLI connects to an independently running SUT through documented APIs and reads host-visible product persistence. It checks preconditions, discovery, convergence, final runtime health, routing and freeze state. It deliberately returns failure after successful final-state checks because operation-order assertions are missing. Isolated fixture provisioning and integrated good/bad exercises remain absent. No contract or specification blocker is recorded.

The Docker command-event investigation identifies a potentially useful ordinary observation boundary, but remains a probe rather than a completed collector. The branch correctly keeps that distinction visible. Its smaller implementation is unfinished, so its lower expenditure cannot be compared as the cost of delivering the same product.

Evidence: [runner](../v5_20261003_234931/luna-high2/test-harness/cmd/traust-test-suite/main.go), especially its final incomplete error; [TODO](../v5_20261003_234931/luna-high2/test-harness/agent-state/TODO.md); patch 0007.

### Sol high1

**Outcome: INVALID full coverage claim. First covered/ready commit: `4bf8b95d2d693f4b8f9e31183fc35014be56b460`, patch 0007.**

This is the strongest ready branch in several structural respects: strict public launch/connection configuration, separate good and bad executables, contained runtime files, real Docker/Traefik, registry-backed digests with an authentication facade, and application observations carried through documented configured runtime inputs. It establishes a reviewable `initialize-current-route` setup contract after a deferred interface question. The stand-in receives a real registry password through its ordinary implementation configuration; that is not itself an undisclosed evidence channel.

Complete acceptance nevertheless fails:

1. The good actor runs `backup` with Docker exec in the baseline application and `migrate` with Docker exec in the candidate. The [product command flow](../documentation@11c2a1/docs/specs/3_backup_and_migration.md) requires disposable command containers. The assertion accepts those application-originated events and does not detect the wrong execution boundary.
2. `assertEventOrder` compares internal verification with hostname verification, without observing the route switch. A trace with routing switched before internal verification is accepted. A review probe reproduced that false acceptance.
3. Finalization is checked through the tenant API after the fixture's unfreeze event. The fixture changes freeze state before posting that event. This does not establish that persisted finalization preceded initiation of unfreeze, especially for an implementation that writes state concurrently with the held command.
4. Default candidate freeze is inferred from absent explicit freeze calls and final state. There is no direct observation of the candidate's required earlier frozen state.

The known-bad actors demonstrate that some errors are detected, but do not repair these omissions. The route setup contract is a useful review surface: route ownership is implementation-defined, while bootstrap through Configuration/activate APIs is an alternative with a different setup dependency. The evidence does not settle which setup policy should become standard.

Evidence: [stand-in](../v5_20261003_234931/sol-high1/test-harness/internal/standin/standin.go), [assertions](../v5_20261003_234931/sol-high1/test-harness/internal/scenario01/evidence.go), [checkpoint runner](../v5_20261003_234931/sol-high1/test-harness/internal/scenario01/run.go), [application](../v5_20261003_234931/sol-high1/test-harness/internal/fixtureapp/fixture.go), [contract](../v5_20261003_234931/sol-high1/test-harness/specs/DEPLOYER_TEST_INTERFACE_CONTRACTS.md), [probe output](validation/sol-high1-probe.log).

### Sol high2

**Outcome: INVALID. First covered commit: `14d180f508ac3f864558a598d1f43e29299a21c9`, patch 0001. Ready commit: `2682671567c918e0dda239c29c7bb55197ce9522`, patch 0002.**

The branch exposes a public CLI and launches separate good/bad processes. It uses real infrastructure, application-originated events and a contained run root. It has no recorded test-only contracts.

The runner fixes the target runtime to `<app>_<tenant>_v2_m24r` and the route filename to `<app>_<tenant>.yaml`. It uses those identities to wait for unfreeze, inspect final state and republish the route. A conforming SUT choosing another permitted runtime suffix or route filename cannot replace the actor by launch configuration alone. Startup also requires an undocumented `/healthz` endpoint.

The good actor executes backup from the target image rather than the persisted baseline image. Patch 0003 corrects materialized environment selection but leaves that image choice unchanged. The verifier also lacks a route-switch observation: the early-route review probe passes here too. A seeded `secrets/registry/local` file supplies an assumed credential convention, while the actor uses hardcoded Basic credentials; no portable binding to a separately supplied implementation is demonstrated.

This is the cheapest ready claim, but it solves a narrower, actor-specific problem. Two unchanged ready iterations consumed another 42,975 logged tokens and 176 seconds after the main work. That is evidence for driver termination, not evidence that the worker should advance the milestone.

Evidence: [runner and assertions](../v5_20261003_234931/sol-high2/test-harness/cmd/traust-test-suite/main.go), [actor](../v5_20261003_234931/sol-high2/test-harness/internal/standin/run.go), [later materialization patch](../v5_20261003_234931/sol-high2/patches/0003-Seed-runtime-materializations-for-happy-path.patch), [probe output](validation/sol-high2-probe.log), and its three log records in the inventory.

### Terra high1

**Outcome: INVALID. First covered/ready commit: `43c0ffec582b43a0e861e936bc2d86dbb59f9542`, patch 0013.**

The suite launches separate container SUTs, bootstraps v1 through Configuration/activate APIs, publishes v2, triggers discovery and reads product persistence. Parts of runtime collection use documented labels rather than fixed identities. Good and bad public exercises exist, with no recorded test-only contract.

Other collectors defeat that portability. The finalization observer fixes the v2 runtime suffix to `_m01r`. Lifecycle assertions require Docker `exec_start` events containing specific readiness and verification HTTP command strings. A deployer using its own HTTP client will produce valid application requests without those Docker exec events; the review probe confirms rejection. Command evidence is collected by inspecting retained containers after deployment, while the documented disposable `--rm` containers should already be gone. The good actor retains them deliberately.

Discovery in the actor reads local Docker RepoDigests. It does not demonstrate the documented registry authentication/tag-list flow. Routes are directed to an additional fixture directory. The separate process and API use therefore do not establish complete independent replacement.

The branch contains useful external route-volume synchronization and an updated Traefik version that works in its recorded environment. Those mechanisms can be considered separately from its invalid completion claim.

Evidence: [finalization observer](../v5_20261003_234931/terra-high1/test-harness/internal/fixture/finalization_observation.go), [lifecycle assertions](../v5_20261003_234931/terra-high1/test-harness/internal/suite/happy_path.go), [command collector](../v5_20261003_234931/terra-high1/test-harness/internal/fixture/happy_path_evidence.go), [actor](../v5_20261003_234931/terra-high1/test-harness/internal/selfexercise/standin/standin.go), [environment](../v5_20261003_234931/terra-high1/test-harness/internal/fixture/environment.go), [probe output](validation/terra-high1-probe.log).

### Terra high2

**Outcome: BLOCKED partial worker, with a public false-pass defect. No completion claim. Final patch: `e5c07ffbfbb987110643671b1ecddcdaeb28e583`.**

The worker finds a canonical container state mount after a fresh-context interface challenge, avoiding a new deployer contract. It builds a real registry, a container discovery actor, API clients, an application and an event recorder configured through normal app inputs. Integrated positive/negative discovery exercises are recorded. Deployment/routing coverage remains partial.

The final stop attributes Traefik Docker-provider failure to an unavailable socket. The fixture uses Traefik 3.1, while a sibling records an old-client API rejection and others succeed with newer Traefik or an API proxy. The manager should classify the blocker as unresolved fixture/engine compatibility, rather than settle the stronger conclusion that container Docker access is unavailable. The engine itself is recorded healthy.

There is a separate consumer defect: `runScenarioHappyPath` returns success after only catalog discovery and automatic target update; the CLI prints `scenario ...: passed`. Its good discovery actor never upgrades the tenant's current version. A manager that checks only covered/ready state would miss this public false pass. The worker's honest partial coverage file does not make the public result correct.

Evidence: [partial scenario function](../v5_20261003_234931/terra-high2/test-harness/internal/suite/scenario_happy_path.go), [public result](../v5_20261003_234931/terra-high2/test-harness/cmd/traust-test-suite/main.go), [discovery actor](../v5_20261003_234931/terra-high2/test-harness/cmd/scenario01-discovery-standin/main.go), [Traefik fixture](../v5_20261003_234931/terra-high2/test-harness/internal/suite/traefik_fixture.go), [question](../v5_20261003_234931/terra-high2/test-harness/agent-state/OPEN_QUESTIONS.md).

### Terra high3

**Outcome: INVALID. First covered/ready commit: `1bacf4242d9927d03937a262ea4574208b1de4e6`, patch 0012.**

The final patch replaces a remaining discovery gap with real registry tags and manifest requests. Separate one-shot SUT actors exercise real Docker/Traefik and several deliberately incorrect orderings. Launching a one-shot main loop can be a valid SUT boundary; lack of a long-running API is not itself the defect.

The required app declares empty runtime environment and mounts, but the actor supplies fixture-specific runtime identity/tenant/hostname values and copies the baseline's shared volume into the candidate. The verifier depends on that candidate writing `fixture-state/events.log`. An independent SUT cannot derive this wiring from the documented app inputs. The observer also fixes the route filename.

Backup and migration execute inside application containers with abbreviated fixture-specific commands. The catalog assertion accepts any correctly sized digest instead of matching the registry's expected immutable digest; a probe confirms that an unrelated digest passes. Timing assertions use polling and brief fixture sleeps, which provide sampling opportunities without establishing an exact ordering boundary. No test-only contract records these dependencies.

Two unchanged ready iterations consumed 22,958 logged tokens and 42 seconds. The worker respected the frontier; the driver continued dispatching it.

Evidence: [seed and assertions](../v5_20261003_234931/terra-high3/test-harness/internal/suite/scenario01.go), [observer](../v5_20261003_234931/terra-high3/test-harness/internal/suite/scenario01_observer.go), [actor wiring](../v5_20261003_234931/terra-high3/test-harness/fixtures/scenario01-standin/standin.go), [baseline setup](../v5_20261003_234931/terra-high3/test-harness/internal/suite/docker.go), [probe output](validation/terra-high3-probe.log).

### Terra high4

**Outcome: WORKING. No completion claim. Final patch: `fe06339b26cea5d4c0bc7c175f0b73f5b7779695`.**

The branch builds substantial baseline infrastructure, public execution logic, persistence and lifecycle observations, routing and freeze checks. It still lacks separate good/bad SUTs and integrated public self-exercises. Its worker state preserves that gap.

Two explicit test-only contracts are established: paired persistent/runtime roots, and a finalization barrier after the durable running write and before target unfreeze. Deferred challenges and justifications make these contracts reviewable. The finalization discussion correctly notices that observing inside an already-started appctl command can be too late to prove invocation order.

Contract necessity remains open. Canonical container mounts are a demonstrated setup alternative; separate container filesystem namespaces can also isolate runtime coordination. Interposing at the ordinary Docker socket before forwarding an exec request is a candidate earlier observation boundary that should be evaluated separately from an app-internal delay. Neither possibility is a settled replacement design. The branch's plain registry also does not implement the configured token-auth boundary, despite seeding an auth reference.

Evidence: [contracts and justification](../v5_20261003_234931/terra-high4/test-harness/specs/DEPLOYER_TEST_INTERFACE_CONTRACTS.md), [runner](../v5_20261003_234931/terra-high4/test-harness/internal/runner/runner.go), [ordinary app logs](../v5_20261003_234931/terra-high4/test-harness/internal/scenario01/observations.go), [TODO](../v5_20261003_234931/terra-high4/test-harness/agent-state/TODO.md), patches 0002–0005.

## Cross experiment findings

### The frontier held while acceptance remained unreliable

The frozen instruction and milestone hashes match across all eight final snapshots. No later scenario implementation was identified. Recorded ready invocations report unchanged state and exit rather than choosing later work. This is a successful control boundary. The runner's repeated dispatch after readiness is an outer-driver issue.

Completion remains a separate problem. Private cooperation can reside in fixed identifiers, assumed HTTP invocation tools, retained command containers and missing configured app inputs, even when launch configuration looks clean. The manager needs to inspect the entire path from setup to SUT action to fixture observation to assertion.

### Oracles require checks for both false rejection and false acceptance

The identity and exec assumptions reject permitted SUT behavior. The missing route observations, arbitrary digest acceptance and partial public success accept prohibited or unfinished behavior. Good/bad stand-ins sharing the same conventions do not expose either class reliably.

The four review probes demonstrate local assertion defects. They should motivate integrated tests using independently constructed actors or changed permitted conventions before settlement. Unit-level probes cannot replace those integrated tests.

### Environment friction is a campaign level finding

Multiple workers encounter registry transport/addressing, Docker API compatibility and host-to-container file notifications. Successful local mechanisms include loopback registry bridging or host networking, OCI upload through real registry endpoints, updated Traefik, and an external route-volume mirror. A full trajectory can fail while a local mechanism still supplies evidence against a sibling's blocker diagnosis.

The campaign does not prove one universal infrastructure layout. It does justify a reusable environment capability record that distinguishes a failed daemon, an old client, a registry inaccessible from one network namespace, and missing file-watch notifications. Those conditions call for different outer-layer responses.

### Authentication and credential provisioning remain unsettled

The [registry specification](../documentation@11c2a1/docs/specs/0_registry_auth_and_tag_discovery.md) describes `auth.secret_ref` or an equivalent secret source, but does not define provisioning/resolution for a standalone consumer. Sol-high1 explicitly uses ordinary implementation-specific launch configuration; Sol-high2 assumes a local secret file and actor credentials. Other partial or ready paths omit auth or substitute local image inspection. A plain Distribution registry does not automatically supply the documented `/token` endpoint.

This supports a human-review question about credential binding and the Scenario 01 authentication obligation. It does not establish that every implementation needs one universal secret-provider API. Ordinary implementation configuration is already allowed. V5 does not show repeated explicit secret-resolution blockers in worker questions; the stronger signal is divergent assumptions and bypasses. The documented protocol and configuration must be compared before interpreting an unauthenticated positive exercise as coverage.

### Infrastructure guidance cost is a hypothesis with incomplete support

There is substantial transport, image-publication, socket and notification machinery, with several independent workarounds. The preference for production infrastructure may contribute to that cost. However, no complete valid branch demonstrates a cheaper protocol-faithful alternative, and subsystem time is not measured. The evidence warrants a targeted fixture comparison, not an automatic conclusion that off-the-shelf guidance is wrong.

### Test interface review needs a necessity and proof audit

Sol-high1 creates a setup contract for route ownership; Terra-high4 creates isolation and timing contracts. Deferral across iterations makes those decisions inspectable, but fresh context alone does not establish necessity. The manager should check ordinary namespace/interposition alternatives and whether the proposed contract actually detects the original prohibited behavior. A callback promising correct order can otherwise move the assumption into the SUT contract.

There is no full valid divergence to rank here. There is meaningful local architectural variance and evidence about underdefined boundaries. Those are review findings, not settled architectures.

Root isolation also requires an ownership audit. [PROJECT.md](../v5_20261003_234931/terra-high4/test-harness/PROJECT.md) requires parallel runs sharing one Docker daemon and run-owned cleanup. The [global lock specification](../documentation@11c2a1/docs/appendix/01_global_daemon_lock.md) assumes one active daemon and requires startup removal of every container with the shared command-role labels. If separate roots or filesystem namespaces permit concurrent SUTs with separate locks but one engine, that cleanup could remove another SUT's live command work. Root separation alone does not solve that conditional consequence. The [causal regression](TERRA_HIGH4_REGRESSION.md) records the alternatives and limits; no concurrent reproduction has been run, and the intended isolation/ownership domain requires clarification.

## Candidate issues for human review

| Issue | Classification | Evidence supported conclusion |
| --- | --- | --- |
| Completion claims and public passes | Evaluation and monitoring gap | Audit public success even without a covered/ready state; partial Scenario 01 cannot report full pass. |
| Independent identifiers and invocation methods | Trajectory defects under existing v5 rules | Reject coupled claims; do not add a permanent worker rule for every observed spelling/tool choice. |
| Ordering proof | Observation design and coverage gap | An eventual record or ordered endpoint subsequence does not prove write-before-invocation or internal-before-route. |
| Registry credential binding | Specification or integration ambiguity | Define how a supplied SUT gets the fixture credential through allowed ordinary configuration, and settle authentication scope. |
| Traefik and registry environment | Fixture compatibility and diagnosis | Record tested namespace/daemon/client capabilities; distinguish compatibility failures from engine failure. |
| Standard infrastructure preference | Guidance hypothesis | Compare a small protocol-faithful fixture against existing machinery before changing the preference. |
| Explicit test-only contracts | Contract review | Review route bootstrap tradeoffs and challenge necessity/proof of isolation/finalization hooks. |
| Parallel deployer ownership | Project/product operating-boundary question | Separate filesystem locks do not scope engine-wide command cleanup; clarify the intended isolation domain before settling a root mechanism. |
| Ready dispatch | Driver behavior | Stop or await external review when the worker is unchanged and ready. |
| Initial milestone size | Workflow proposal | Consider evidence-based decomposition or intermediate review gates; no evidence mandates an OCI warmup milestone. |

No worker/project/product instructions have been changed by this evaluation. A subsequent manager decision should first reject/reopen the concrete claims, validate useful local mechanisms, and resolve the review questions that affect further investment.

## Verification and evidence gaps

- All four focused review probes passed, meaning they reproduced the expected defect. Their sources and logs are retained under `validation/`.
- Sol-high2's unit suite passed in the review copy. Sol-high1's assertion/package tests passed, but its broader suite reached a Docker-dependent test and failed on sandbox pipe access. This is not evidence that the campaign daemon was unhealthy.
- Terra-high1's suite package passed; its fixture containment unit test failed because it compares an unnormalized temporary-root string with a filepath-normalized child. This is a local Windows test-path result, separate from the decisive SUT boundary rejection.
- Terra-high3's selected pure tests passed. Its Docker tests were excluded deliberately; integrated success is retained as a worker-log claim.
- The exact subsystem share of time/tokens, prices, maintenance behavior on a changed master, and performance of a correctly completed foundation are not established.
- The handoff exposes benchmark hypotheses and the prior manager evaluation. This review can assess whether the proposed manager mechanism addresses the evidence, but cannot claim a blind empirical test of independent hypothesis discovery.
