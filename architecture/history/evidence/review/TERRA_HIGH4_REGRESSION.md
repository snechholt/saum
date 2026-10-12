# Terra-high4 manager regression

The revised reusable evaluator derives exceptional decision classes from project authority, then vets the evidence before producing a human-facing finding. Terra-high4 supplies the regression evidence; its vocabulary and expected conclusions are not embedded in the reusable prompts or campaign manifest.

## Vetted human-facing finding

**Clarify the deployer isolation and ownership boundary before settling an additional root-selection requirement.**

V5 requires independent harness runs to share one Docker daemon without collisions and to clean up only their own resources. The product assumes a single active deployer daemon, protected by `/run/traust/daemon.lock`, and requires startup to remove every container carrying `traust.managed_by=traust` and `traust.role=command`. Giving concurrent SUTs separate filesystem roots and locks does not partition those engine-wide labels or that cleanup. A second conforming daemon following the documented startup rule could therefore remove command work belonging to the first. This is a conditional consequence of the rules, not an observed incident in High4.

High4's paired-root mechanism has not established necessity or complete isolation. An ordinary container mount has demonstrated the persisted-state portion in a sibling; the remaining ownership issue survives both root relocation and that alternative when they share an engine. The project and product owners should clarify the intended isolation domain and how parallel conformance exercises preserve daemon exclusivity and cleanup ownership. Do not settle a universal root-selection interface as the answer to that question.

Confidence is high in the mismatch between separate locks and unscoped command cleanup, and moderate in the scope clarification needed: the product's intended host/installation/engine boundary is not fully specified here. No concurrent independent-daemon reproduction or complete alternative isolation exercise has been run in this review. The finding warrants that clarification and bounded validation; it does not establish a preferred replacement architecture.

Authority: [parallel isolation](../v5_20261003_234931/terra-high4/test-harness/PROJECT.md), lines 129–133; [global daemon lock and startup cleanup](../documentation@11c2a1/docs/appendix/01_global_daemon_lock.md), lines 3, 14 and 40–48. Mechanism: [High4 contracts](../v5_20261003_234931/terra-high4/test-harness/specs/DEPLOYER_TEST_INTERFACE_CONTRACTS.md).

## Input boundary

The intended discovery-test inputs are the generic procedures, [CAMPAIGN.md](manager-candidates/CAMPAIGN.md), the immutable project/worker and product sources it references, and experiment evidence. This report, the earlier review, the handoff's hypotheses and expected answers are withheld from such a test.

This regression was performed by the existing reviewer, who has already seen the earlier archive assessment. It is a source-grounded application of the revised mechanism, not a blind evaluator run. Source claims, inference and untested alternatives are distinguished below.

## Authority-derived salience

| Source-derived concern | Why it receives deeper review | Trigger in High4 |
| --- | --- | --- |
| Additional deployer test-interface requirements | [PROJECT.md](../v5_20261003_234931/terra-high4/test-harness/PROJECT.md), lines 178–186, explicitly calls them last-resort, rejects convenience alone, prescribes deferral and a fresh alternative search, and permits later human rejection or product promotion. | Two additions in the mutable established-contracts inventory: paired filesystem roots and a finalization barrier. |
| Private cooperation across the independent implementation boundary | [PROJECT_PREFERENCES.md](../v5_20261003_234931/terra-high4/test-harness/PROJECT_PREFERENCES.md), lines 21–35, restricts allowed inputs and prevents private semantics from being hidden in launch configuration. | New launch variables are supplied to every SUT by the runner; their allowed status and necessity must be assessed rather than inferred from successful configuration parsing. |
| Daemon ownership and effects on other runs | The [global lock specification](../documentation@11c2a1/docs/appendix/01_global_daemon_lock.md) protects a single-daemon assumption and says multiple active instances require a different coordination design. [PROJECT.md](../v5_20261003_234931/terra-high4/test-harness/PROJECT.md), line 131, requires shared-engine isolation and run-owned cleanup. | Separate runtime roots create separate lock locations while leaving the engine and product command labels shared. This trigger is derived from the mechanism's consequence, not from its title. |

The [stable interface guidance](../v5_20261003_234931/terra-high4/test-harness/specs/DEPLOYER_TEST_INTERFACE.md) confirms that the inventory is a human review surface and **does not require a separate approval workflow for the worker to continue**. The derived process therefore includes a challenge and reviewable justification, not invented pre-approval. The generic manager does not know these project categories in advance.

## Paired-root causal review

### Concrete pressure and history

Scenario 01 requires a deterministic global registry profile and persisted v1 baseline before discovery. The [Configuration API](../documentation@11c2a1/docs/specs/5.3_configuration.md) manages app configuration and tenant desired state; High4 identifies no operation that creates the required global registry profile. The [directory specification](../documentation@11c2a1/docs/appendix/07_directory_structure.md) places persistent state at `/var/lib/traust` and runtime coordination at `/run/traust`. Project preferences contain host-side run files beneath a reported work root, while PROJECT requires parallel isolation.

High4 asks how its command-launched SUT will use the seeded per-run state and avoid a shared runtime lock. It queues the question in patch 0004, commit `92d651420d299af1c72f7625354ffa138b7b0eed`, then establishes the paired-root mechanism in patch 0005, commit `8648543c8a2aea09103850fd82f428c93e99baf4`. The prior question and subsequent contract show the recorded deferred process; they do not prove that all reasonable alternatives failed.

Evidence: [queued challenge](../v5_20261003_234931/terra-high4/patches/0004-test-harness-queue-bootstrap-boundary-review.patch), [established mechanism](../v5_20261003_234931/terra-high4/patches/0005-test-harness-establish-isolated-SUT-root-contract.patch).

### Authoritative constraints versus local assumptions

The documented layout, required baseline, containment and collision avoidance are real constraints. The conclusion that isolation must use a new pair of environment variables is a worker design choice. Lack of a documented root-relocation setting does not establish that ordinary deployment isolation is impossible.

The contract dismisses container mounts as insufficient because an independent implementation needs a declared way to use the mounted roots and because a persistent mount alone leaves the runtime lock shared. This mixes two issues. A mount at the documented canonical path can leave the SUT's persistent path unchanged; runtime coordination must be considered separately. Moreover, a chosen command launcher does not require the launched program to run in the host's filesystem namespace.

The [runner](../v5_20261003_234931/terra-high4/test-harness/internal/runner/runner.go), lines 109–112, rejects a connection-only configuration because it must supply the isolated roots. Its `sutEnvironment` unconditionally appends the two suite-selected variables, lines 297–300. The [design notes](../v5_20261003_234931/terra-high4/test-harness/agent-state/DESIGN.md), lines 23–27, describe command launch as necessary because of the paired-root mechanism. That is a consequence of the chosen mechanism, not independent evidence that the mechanism is necessary.

Ordinary implementation-specific configuration is also allowed by PROJECT_PREFERENCES. Where an implementation already provides suitable deployment configuration, a normal launch binding may be sufficient; the archive does not establish such support universally. Neither this possibility nor a container launcher proves full portability by itself.

### Alternatives and their demonstrated scope

| Alternative | Evidence | Limit |
| --- | --- | --- |
| Mount a per-run persisted directory at the SUT's canonical path | Terra-high2's [runner](../v5_20261003_234931/terra-high2/test-harness/internal/suite/runner.go), line 200, mounts `stateRoot:/var/lib/traust`. Its [patch 0005](../v5_20261003_234931/terra-high2/patches/0005-Add-containerized-SUT-state-boundary.patch) records real bidirectional visibility and resolution of the earlier bootstrap question without a test-only interface. | Demonstrates persisted preconditions and reads in a containerized discovery actor. It does not prove runtime-lock isolation, subsequent host-path translation or full deployment conformance. Terra-high2's overall scenario remains incomplete and its public pass is invalid. |
| Give each SUT an ordinary isolated filesystem namespace retaining canonical paths | Ordinary deployment isolation extends the previous alternative to runtime paths; the need to consider it follows from the project's allowed external-environment alternatives. | Not validated here as a complete solution. Sharing the engine still leaves the product-wide cleanup issue; container-local paths must also agree with the engine's bind-mount path semantics. |
| Use a supported implementation's ordinary deployment settings | PROJECT_PREFERENCES explicitly permits normal implementation configuration. | No universal configuration interface or complete alternative implementation is demonstrated. It cannot be treated as a settled universal answer. |
| Supply an isolated engine view through ordinary environment interposition | PROJECT, line 182, permits runtime/environment interposition. A scoped view could retain a shared physical engine while limiting each SUT's visible command ownership. | No such arrangement is demonstrated here. It must preserve product-visible effects, resource identity, lifecycle and isolation rather than manufacturing evidence. Its feasibility is ordinary validation work, not evidence that a new SUT interface is necessary. |
| Give each SUT a fully independent engine, or serialize independent daemon activity | Either can avoid the specific concurrent shared-engine cleanup condition in principle. | Neither has been shown to preserve the project's shared-daemon parallel-run requirement in the intended deployment envelope. An arrangement that relaxes that requirement needs an owner decision; a nested environment that preserves it needs validation. Neither is an established fix here. |

A sibling success challenges one claimed obstacle; it does not prove the sibling architecture solves every obligation. Necessity of the paired-root mechanism remains **unproved**, rather than conclusively disproved for every supported deployment mode.

### Hidden consequence and residual value

The product's startup cleanup assumes that after the daemon acquires its global lock, every matching command runtime belongs to interrupted earlier work. With separate locks and a shared engine, that inference is no longer assured. Per-run container names do not change a cleanup predicate that uses only the shared product labels. The root contract says it does not alter locking rules or engine semantics, but it changes the effective exclusivity domain when concurrent SUTs use distinct lock files.

This consequence also survives the ordinary-filesystem alternative. It is therefore not merely an argument against High4's chosen variables. It exposes an ownership restriction that matters when interpreting the project's shared-engine isolation requirement against the product's exclusivity/cleanup assumption. A fully isolated or appropriately scoped engine environment could remove the collision condition; those options have not been validated. The residual human finding is the documented ownership/isolation boundary and its consequence, not a claim that no implementation can satisfy both sources or that a new interface is required.

## Finalization-barrier causal review

The second mechanism is picked up through the same source-derived exceptional class. Its pressure is the [state transition requirement](../documentation@11c2a1/docs/appendix/04_state_transition_spec.md) to persist Finalize running before initiating target unfreeze. High4 correctly observes that a fixture hold inside an already-started `appctl` command can occur after the invocation; a later state read cannot prove the required earlier write.

Patch 0002 queues the timing challenge; patch 0003, commit `308be26ad5d3f0db32f29fb42cec65c5046d3870`, establishes a callback after the durable write and before unfreeze. The [contract justification](../v5_20261003_234931/terra-high4/test-harness/specs/DEPLOYER_TEST_INTERFACE_CONTRACTS.md) treats event observation and a proxy as too late and assumes a proxy requires an undocumented endpoint. An interposer holding the relevant execution request **before forwarding it** has different timing semantics from an after-start event. The authority expressly permits runtime/environment interposition. That earlier boundary is an alternative to investigate, not a demonstrated complete replacement here.

The callback's promised position also does not independently prove that no earlier unfreeze occurred. A negative exercise must test the original prohibited order, not merely that an actor honoring the callback supplies a correctly ordered record. High4 has no completed public good/bad exercise establishing that proof.

Disposition: retain the valid timing insight; route necessity and negative-proof work as ordinary scoped validation. Do not forward “approve the barrier” or mandate a new observation mechanism. No additional human decision is established by this timing proposal alone. If bounded validation later demonstrates that an authority change is unavoidable, it can return as a newly vetted finding.

## Filtering result

| Candidate | Disposition | Default human-facing result |
| --- | --- | --- |
| Adopt paired root variables for every independent SUT | Suppress as an unsupported necessity claim; retain source, alternatives and local correction/validation evidence. | No API-approval request. |
| Treat the sibling mount as a complete answer | Suppress the overclaim; retain only demonstrated persisted-state provisioning. | No recommendation to adopt its architecture. |
| Approve the finalization callback because the fresh challenge occurred | Suppress; process evidence does not establish necessity or proof. | No callback-approval request. |
| Clarify isolation and cleanup ownership across parallel SUTs | Surface: the material consequence survives both the proposal and an ordinary alternative, and reconciling the conflicting operating assumptions needs the project/product owners. | The single vetted finding at the top of this report. |

No source requires a separate approval request for each established contract. Existing inventories remain the required review records; filtering changes what is recommended for human attention, not what evidence is retained.

## Regression assertions and limits

- The reusable prompts contain no predetermined domain class or root conclusion; the relevant class and process are derived from cited v5 authority.
- High4's decisions trigger causal review through both their declared class and their effects on protected assumptions.
- Recorded deferral is not treated as proof of necessity or as an approval requirement.
- Sibling evidence is usable within its demonstrated scope despite an invalid overall trajectory.
- Avoidable/unsupported mechanism proposals are suppressed; the surviving finding concerns ownership semantics and does not endorse those proposals.
- The generic synthesis preserves the suppression decisions instead of reviving them as informative worker guidance.

This check uses static source reasoning and archived execution evidence. It adds no full execution, concurrency reproduction or independent-deployer validation. Those gaps limit the proposed solutions, while the cited unconditional cleanup predicate supports the conditional consequence. A fresh withheld-answer evaluation remains necessary to test discovery reliability.
