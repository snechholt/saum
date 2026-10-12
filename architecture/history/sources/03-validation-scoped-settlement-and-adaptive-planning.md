# Zonotron — Validation, Scoped Settlement and Adaptive Reconciliation Planning

**Branch summary — October 11, 2026**

**Scope:** This discussion started with the concern that, even if Zonotron delegates reconciliation methodology to individual projects, it still needs to determine whether the resulting artifacts actually conform to their upstream dependencies. It developed into a model for validation authority, dynamically negotiated settlement increments, persistent reconciliation plans, and incremental acceptance.

This is a working architectural direction, not a finalized specification.

## 1. Starting problem: Zonotron cannot treat project outputs as black boxes

The broader Zonotron architecture delegates *how reconciliation is performed* to individual projects. A Test Harness project could use Ralph, Spec Kit, custom agents or some other methodology without Zonotron needing to manage its internal workflow.

However, Zonotron remains responsible for accepting the resulting state into its authoritative artifact graph.

For example:

- Architecture A13 is settled.
- Test Harness produces candidate TH5, claiming conformance with A13.
- Zonotron must determine whether TH5 actually satisfies its obligations before accepting it.

Two concerns emerged:

1. **Zonotron must retain settlement authority.** Accepting a project's claim of completion is insufficient.
2. **Validation requires understanding the resulting artifacts.** Arbitrary formats and heterogeneous project methodologies make universal validation difficult.

The second concern initially suggested that Zonotron would need to perfectly parse and understand every project's artifacts.

### Emerging conclusion

**Zonotron must own settlement authority, but it does not need to implement every validation procedure itself.**

Validation can be performed by separate processes, including deterministic tools, independent agents, domain-specific validators, experiments and human reviewers.

Zonotron controls the acceptance process, evidence requirements and settlement decision.

This distinguishes three responsibilities:

| Responsibility                               | Owner                                                      |
| -------------------------------------------- | ---------------------------------------------------------- |
| Produce candidate artifact                   | Project                                                    |
| Evaluate conformance to upstream obligations | Independent validation process                             |
| Decide whether evidence warrants acceptance  | Zonotron, under configured human/automatic approval policy |

Importantly, the validation process is not infallible either. Independent evaluation reduces certain risks but does not guarantee correctness.

**Settlement therefore means accepted under a recorded validation policy with available evidence, not mathematically proven correct.**

---

## 2. Three contracts at the project boundary

The discussion produced a useful division into three conceptual contracts.

### A. Reconciliation contract

Defines what a project must achieve.

Contains, conceptually:

- authoritative upstream revisions;
- current settled project state;
- reconciliation objective;
- scope of required changes;
- constraints that must be preserved;
- authority for changing those obligations.

For example:

> Reconcile Test Harness against Architecture A13, establishing coverage of its scenario-testing obligations.

This contract is relatively stable during execution. Changing the implementation approach does not change its obligations.

### B. Artifact contract

Defines what a project produces and how it can be inspected.

May include:

- candidate revision;
- artifact location and format;
- execution provenance;
- tests and other evidence;
- project-reported status;
- optional observations;
- claimed coverage.

Projects can use arbitrary internal methodologies, but resulting artifacts need to be inspectable through some suitable mechanism.

A project producing an unfamiliar proprietary format with no available interpreter or meaningful validation evidence cannot expect automatic conformance verification.

Zonotron may require additional tooling, alternative evidence or human review.

### C. Validation contract

Defines what conformance means for a specific dependency relationship and how it should be evaluated.

Different dependency edges imply different obligations.

Examples:

- **Architecture → Allium:** The behavioral model faithfully represents applicable architectural requirements.
- **Architecture → Test Harness:** The harness adequately verifies required behavior through legitimate testing boundaries.
- **Design → Implementation:** Code conforms to settled structural and interaction decisions.
- **Test Harness → Implementation:** Implementation passes the required externally defined tests.

These relationships cannot all be validated in the same way.

A promising concept is therefore **validation contracts associated with dependency edges**, rather than a universal artifact parser or validation algorithm.

Their exact representation and ownership remain open.

---

## 3. Validation confidence should be proportional to risk

The discussion considered two extremes.

### Routine local change

Implementation updates an external library, fixes several breaking API calls, runs tests and reports:

> Library X updated. Three callers adjusted. No intended behavioral changes.

Requiring expensive independent semantic evaluation and human approval for every such change would be wasteful.

The proposal is a configurable approval policy that can initially be conservative and become more permissive as validation reliability is established.

Examples:

- routine dependency update: tests and automated acceptance;
- change to core state or security logic: additional independent review;
- significant behavioral change: reconciliation-focused review;
- foundational architectural change: human settlement and potentially experiments.

**Trust should be calibrated per class of change, not merely per project.**

Repeated success with dependency upgrades does not establish reliability for architectural changes.

Occasional independent audits of automatically accepted work would help establish whether the policies are working.

### Large or foundational change

Large changes need stronger validation, possibly including decomposition, independent review, comparative experiments and multiple human settlement opportunities.

The key question becomes how to divide a large reconciliation into manageable, independently reviewable increments without freezing an unsuitable execution plan.

---

## 4. Separate objective, execution plan and settlement criteria

This was the central conceptual development.

Three things must not be conflated:

1. **Reconciliation objective:** What must ultimately become true?
2. **Execution plan:** How does the project currently intend to achieve it?
3. **Settlement criteria:** What evidence is sufficient to accept the result?

The reconciliation objective and its authoritative obligations remain stable unless explicitly revised.

The execution plan is provisional. It represents the current understanding of how to progress, and may change as work reveals new information.

Settlement criteria must be grounded in upstream obligations, not solely in the project's chosen implementation.

For example, Scenario 01 might initially be decomposed into:

- OCI registry fixture;
- Docker command interception/proxy;
- replaceable test SUT;
- actual scenario test.

That is an execution hypothesis, not a binding requirement to implement those exact mechanisms.

A project may discover better approaches and revise its plan.

Zonotron should validate the resulting scenario coverage, not merely confirm that the four original tasks were completed.

**The objective constrains the plan; the plan does not redefine the objective.**

---

## 5. Persistent, dynamically recalibrated settlement plans

The user's proposal was to make the negotiated reconciliation breakdown a permanent, versioned artifact that evolves after each accepted increment.

This appears to resolve much of the previous uncertainty about how much task decomposition Zonotron must perform upfront.

### Example: creating a Test Harness from nothing

Architecture A13 defines scenarios S01–S20.

The Test Harness project initially has no implementation.

Zonotron requests:

> Produce a Test Harness satisfying the scenario-testing obligations defined by A13.

The project proposes Plan P1:

```text
[ ] TH1: S01 covered
[ ] TH2: S02–S05 covered
[ ] TH3: S06–S20 covered
```

Zonotron records P1, optionally after human review.

The project implements TH1 and submits a candidate.

Zonotron validates the candidate against the obligations of Scenario 01, potentially with human review.

Once TH1 is accepted as an internal checkpoint, Zonotron asks the project to update its remaining plan.

The project reports that S01 uncovered unexpected infrastructure complexity and proposes P2:

```text
[x] TH1: S01 covered
[ ] TH2: S02 covered
[ ] TH3: S03–S05 covered
[ ] TH4: S06–S20 covered
```

Zonotron stores P2 while retaining P1 and the reason for the revision.

This continues until the full reconciliation objective has been met and validated.

### Important properties

**Plans are persistent and versioned.** Previous plans and revisions remain available for provenance.

**Plans are dynamically recalibrated.** Learning during one increment can change the appropriate breakdown for subsequent work.

**A plan need not initially describe the entire journey.** It might contain only TH1: S01, together with an explicit indication that the remaining scope will be planned afterward.

**The overall obligation remains known.** An incomplete planning horizon must not mean that remaining requirements have been forgotten or implicitly dropped.

**Zonotron can expose current progress.** UI/CLI can show the project reconciling toward A13, its current planned increment, accepted checkpoints and outstanding scope.

**Plan changes and scope changes are different operations.** Reorganizing S02–S05 into smaller increments is ordinary recalibration. Deciding S17 is no longer required changes the reconciliation obligation and needs appropriate authority.

The evolution of the plan is potentially valuable evidence in itself: repeated replanning or unexpected difficulty may reveal weaknesses in upstream assumptions.

---

## 6. Settlement increments versus internal implementation tasks

The negotiated increments should represent **meaningful validation/settlement boundaries**, not necessarily every implementation task.

For example:

```text
Settlement increment:
  TH1 — Scenario 01 coverage

Internal project tasks:
  implement registry fixture
  build runtime adapter
  establish replaceable SUT
  implement observers
  write assertions
  fix infrastructure bugs
```

The project can freely reorganize the internal tasks without requesting approval for every change.

Zonotron is primarily interested in the externally meaningful increment and its acceptance criteria.

This also clarifies the role of the original milestone concept.

**A milestone is a candidate settlement boundary, not a prescribed collection of implementation steps.**

Increments may be sequential, parallel or dependent on each other. Their ordering should reflect real constraints rather than an assumption that everything must progress linearly.

---

## 7. Partial acceptance must not automatically become global settlement

A significant caution emerged around letting downstream projects consume partially completed upstream work.

Suppose TH1 successfully validates Scenario 01, but the Test Harness still needs S02–S20.

It would be tempting to let Implementation begin using TH1 immediately.

However, later harness work might uncover a foundational problem and require TH1 to be substantially revised.

Starting downstream work too early could create expensive rework.

### Proposed conservative default

Distinguish:

- **Accepted internal checkpoint:** A particular increment has passed its validation.
- **Globally settled project revision:** The project has completed reconciliation against its required upstream obligations.

TH1 can be internally accepted without making the Test Harness globally settled against A13.

Internally accepted work remains revisable if later evidence invalidates an earlier decision.

Conceptually, this resembles retaining work on an unmerged reconciliation branch until the overall project is ready for settlement.

For the initial MVC, downstream projects would normally wait for full upstream reconciliation.

Future versions might allow partial downstream consumption with explicit dependency scopes and an understanding that such work is provisional.

---

## 8. Three different questions about whether the plan is sufficient

The discussion clarified that "are the negotiated increments sufficient?" covers three separate concerns.

### A. Completeness

Does the plan collectively account for all upstream obligations?

Example:

Every scenario is tested individually, but the complete Test Harness cannot run them together or lacks a consistent replaceable-SUT boundary.

All increments might appear individually complete while the overall objective remains unmet.

Zonotron should assess plan coverage during negotiation and validate the complete artifact before final settlement.

### B. Soundness

Does the evidence actually demonstrate the claimed behavior?

Example:

A fake OCI registry responds in a way that unintentionally bypasses functionality the test harness was supposed to exercise.

The test passes, but it doesn't establish the intended guarantee.

**Soundness is the most difficult of these concerns.**

It requires understanding what the test or artifact genuinely proves, not merely checking completion status.

The Traust experiments demonstrated how seemingly successful artifacts can fail at this level.

### C. Proportionality

Is the chosen approach unnecessarily complex or expensive?

Example:

The Test Harness project builds a fully OCI-compliant registry when a small fixture implementing three endpoints would adequately exercise the intended contract.

This is primarily an efficiency/design concern, not necessarily a conformance failure.

It can be addressed during planning, local project review or comparative evaluation.

Zonotron may flag it without requiring every implementation choice to be centrally approved.

---

## 9. Validation must consider accumulated progress and integration

Validating only the latest increment is insufficient.

When TH2 is completed, Zonotron should consider:

- Did TH2 establish its new obligations?
- Does the resulting candidate still satisfy previously accepted TH1 obligations?

Regression checks may handle much of this mechanically.

Furthermore, even if every increment is accepted individually, the complete project may still violate cross-cutting or integration constraints.

Therefore, the proposed model has two kinds of validation:

**Incremental validation:** Establishes an increment's new claims and checks that earlier accepted claims remain valid.

**Final reconciliation validation:** Establishes that the integrated artifact satisfies the complete authoritative reconciliation objective.

The latter is required before global settlement in the conservative model.

---

## 10. Projects can propose validation evidence, but cannot be their own sole authority

An interesting idea was that the project's decomposition could inform how Zonotron validates its output.

For example:

> I divided Scenario 01 into these four mechanisms. You may want to validate each one independently.

This is useful because the project understands its implementation.

However, allowing the project to define the entire acceptance standard risks reproducing its blind spots.

Validation should therefore combine two sources:

**Upstream-derived obligations**

What the authoritative artifacts require, independently of the project's chosen solution.

**Project-proposed evidence**

The mechanisms, tests, checkpoints and observations the project believes establish its claims.

These can be compared by an independent reviewer.

The project can revise its proposed evidence strategy as implementation evolves, without gaining authority to weaken upstream obligations.

---

## 11. Suggested end-to-end reconciliation loop

The emerging model is:

```text
Settled upstream revision
          |
          v
Zonotron establishes overall
reconciliation objective
          |
          v
Zonotron ↔ Project negotiation
  - meaningful increments
  - initial planning horizon
  - dependencies/order
  - evidence expectations
          |
          v
Store versioned Plan P1
          |
          v
Project executes next increment
using its own methodology
          |
          v
Candidate + claims + evidence
          |
          v
Independent validation
          |
      +---+---+
      |       |
   accepted  rejected/blocked
      |       |
      |       +--> further work / clarification
      v
Record accepted internal checkpoint
          |
          v
Project reassesses remaining work
          |
          v
Store revised Plan P2, if needed
          |
          v
Continue until full scope covered
          |
          v
Final integrated reconciliation review
          |
          v
Zonotron/human settlement
          |
          v
New globally settled project revision
```

The plan may evolve throughout this loop, but the original reconciliation objective remains authoritative unless explicitly changed.

---

# 12. Possible MVC scope

The following is a proposed minimal implementation direction, not something formally settled during the discussion.

An initial MVC could support:

1. **Versioned project dependencies:** A project knows which upstream revisions it must reconcile against.
2. **Overall reconciliation objective:** Stored separately from the execution plan.
3. **Versioned negotiated plan:** Human/agent-assisted planning with meaningful increments, accepted progress and remaining scope.
4. **Adaptive replanning:** Plans can change after accepted increments without losing history or silently changing the objective.
5. **Increment-level candidates:** Projects submit revisions and claims against the current increment.
6. **Simple validation protocol:** Candidate, evidence, validation result, optional human review.
7. **Internal checkpoints:** Accepted increments that do not yet constitute a globally settled project revision.
8. **Final settlement gate:** Human-controlled acceptance of the fully reconciled project.
9. **Basic status visualization:** Current reconciliation objective, plan version, active increment, accepted increments, outstanding scope and blockers.

For an MVC, validation could initially rely heavily on independent LLM review, available tests, Git diffs and explicit human decisions.

The MVC does not need universal artifact parsing, sophisticated automated trust calibration or partial downstream consumption.

The key thing to demonstrate is that a large reconciliation can be progressively planned, executed, evaluated and settled without Zonotron having to manage the project's internal development methodology.

---

# 13. Future features and unresolved architectural questions

### Validation architecture

- How are validation contracts represented and attached to project dependencies?
- Who defines and maintains them?
- Should validation contracts themselves be versioned and independently settled?
- How does Zonotron choose among deterministic tests, semantic reviewers, formal tools, experiments and human review?
- What should happen when available evidence cannot establish an obligation?
- How do we distinguish insufficient evidence from an actual conformance failure?
- How much semantic understanding must Zonotron possess versus delegate to independent validators?

### Settlement and scope

- How granular should an accepted settlement increment be?
- How are accepted claims represented without introducing an excessively detailed requirements database?
- What does it mean to reopen an earlier accepted checkpoint?
- How are earlier validation results invalidated after a later structural change?
- Can different increments within the same project become independently authoritative?
- Under what conditions can downstream projects consume partially reconciled upstream state?
- How do we prevent expensive downstream rework if partial upstream settlement is later reversed?

### Planning and negotiation

- What specific questions should Zonotron ask a project when negotiating its reconciliation plan?
- How does Zonotron identify missing obligations in a proposed plan?
- How much of the decomposition should be generated by Zonotron versus proposed by the project?
- What triggers replanning?
- Must the project provide a revised plan after every increment, or only when needed?
- How are plan revisions accepted and their reasons preserved?
- How are global dependencies considered when negotiating project-local settlement increments?
- How is the final integration obligation established without requiring a perfect decomposition upfront?

### Trust, autonomy and cost

- What change classes qualify for automatic settlement?
- How is reliability measured per project and change category?
- Should auto-accepted changes be randomly sampled for independent review?
- What confidence/evidence standards justify removing human approval gates?
- How does Zonotron avoid expensive review of routine changes while still catching subtle semantic failures?
- Can the project propose a validation strategy that is cheaper but equally convincing?
- When should Zonotron request experiments or independent alternative implementations?

### Provenance and observability

- How much information from each increment must be retained?
- Should plan versions, evidence, reviews and decisions link directly to Git commits and execution trajectories?
- Can accepted findings be reevaluated if new evidence emerges?
- How do we preserve explanations of why the plan changed?
- Can repeated replanning or unexpected implementation complexity automatically trigger upstream analysis?

---

# 14. Main architectural conclusions

The strongest conclusions from this branch are:

**1. Zonotron owns settlement authority, not necessarily validation execution.**

Validation is pluggable and may require domain-specific tools and human judgment.

**2. Projects own their methodology, but their outputs must be inspectable.**

Unknown artifact formats are not inherently prohibited, but meaningful validation requires suitable inspection/evidence mechanisms.

**3. Reconciliation objective, execution plan and settlement criteria are separate concepts.**

The objective constrains the work. The plan can evolve. Evidence must establish the objective rather than merely completion of the plan.

**4. Reconciliation plans should be persistent, versioned and dynamically recalibrated.**

The system can learn the appropriate decomposition through implementation instead of demanding perfect upfront planning.

**5. Planning horizon and reconciliation scope are independent.**

The system can know the complete obligation while planning only the next meaningful increment.

**6. Incremental acceptance and global settlement are different.**

Accepted progress can remain internal and provisional until the project completes reconciliation.

**7. The combined result needs independent validation.**

All increments being accepted does not necessarily mean the integrated artifact satisfies the complete upstream specification.

**8. Completeness, soundness and proportionality are distinct evaluation concerns.**

Completeness concerns coverage; soundness concerns whether evidence truly proves claims; proportionality concerns cost and complexity.

**9. Validation confidence should scale with risk and demonstrated reliability.**

Routine changes may be automatically accepted, while consequential changes require stronger evidence and potentially human settlement.

**10. Universal automatic semantic validation remains an open limitation.**

Zonotron can provide an auditable, configurable process for establishing confidence, but cannot guarantee perfect conformance of arbitrary artifacts.

## Final synthesis

The emerging model is:

> **Zonotron coordinates reconciliation toward stable, authoritative objectives through persistent, dynamically negotiated settlement plans. Projects independently determine how to execute the work and can revise their plans as they learn. Zonotron validates meaningful increments against upstream obligations, preserves accepted progress and provenance, and settles the complete project revision only when sufficient evidence establishes conformance.**

This extends the earlier Zonotron concept without changing its fundamental direction: a graph of independently authoritative artifacts, advanced through controlled reconciliation rather than blind regeneration.

The most important new concept is the **living settlement plan**, which connects high-level reconciliation objectives to project-local execution while keeping planning flexible and validation independent.

The largest unresolved architectural challenge remains **soundness of semantic validation**: establishing that an artifact genuinely satisfies the intended upstream requirements rather than merely satisfying its own tests or implementation plan.
