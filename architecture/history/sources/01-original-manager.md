# Deterministic Outer Manager for Agentic Software Development

> **Status:** idea inventory / design handoff, not a settled architecture.
>
> This document captures the concepts, constraints, examples, speculative features, and open questions discussed while evolving the Traust Ralph workflow into a more general software-development system. It intentionally preserves more detail than a normal architecture note so later contexts can recover the reasoning without having to reconstruct it from chat history.

## 1. Core premise

The system is not meant to be “give an agent a vague prompt and let it build an app.” It assumes a comparatively mature, structured body of documentation/backlog that humans primarily own and refine, while agents progressively materialize that specification into derived artifacts such as architecture, test harnesses, detailed designs, interaction models, and implementation.

The strongest examples discussed were:

- **TODO application:** the project would not begin from only “make me a TODO app.” It would have a settled specification comparable in richness to the Traust documentation. A first architectural baseline might be considered established once the app can host its real HTTP API, initialize real state, allow a user to self-sign-up and log in, while TODO listing/managing remains intentionally deferred.
- **Traust test harness:** the backlog is the documented scenario suite. The first architectural baseline can be considered established when Scenario 01 runs through the intended public test-suite / replaceable-SUT boundary with the expected real infrastructure and evidence semantics.
- **Traust implementation:** its backlog could be the product documentation plus obligations such as passing `traust-test-suite` for successive scenarios.

The quality and structure of the backlog matters enormously. A major assumption behind the whole system is:

> Humans spend much of their effort defining, refining, reviewing, and settling the backlog/documentation. Agents turn that settled intent into progressively more concrete artifacts and code.

This is closer to **specification-driven materialization** than unconstrained code generation.

---

## 2. The permanent object is the settled state, not an agent session

The most important conceptual shift is that the permanent thing is **not**:

- a particular agent conversation,
- a single “winning” experiment branch,
- a prompt revision,
- or code that can always be regenerated from scratch.

The permanent object is a **versioned graph of settled artifacts and accepted facts**, with provenance for how those facts were established.

Experiments are temporary branches around that settled graph. They produce:

- evidence,
- candidate implementations,
- architectural alternatives,
- invalid approaches,
- implementation friction,
- possible upstream problems,
- and distilled learnings.

Nothing experimental becomes truth merely because an agent produced it. A review / settlement transition is what advances accepted state.

This matters because downstream artifacts accumulate reviewed nuance. The system should not behave as if everything is a pure function that can always be deleted and regenerated from the latest upstream prose.

---

## 3. Two different graphs

A useful distinction is between two graphs that should not be conflated.

### 3.1 Artifact dependency graph

This describes the durable SDLC/artifact relationships.

Example:

```text
Architecture
    ├──────────────► Test Harness
    │
    ▼
Component / Design Breakdown
    │
    ▼
Interaction Models       (optional layer)
    │
    ▼
Implementation
```

A different project could insert, remove, or split layers.

Possible artifacts/layers include:

- product / system architecture,
- behavioral specification,
- test harness,
- component design,
- API/interface definitions,
- data models,
- interaction models / sequence diagrams / state machines,
- implementation,
- deployment design,
- operational runbooks,
- security model,
- etc.

The artifact graph is how the user chooses **which level of engineering to retain direct control over and which level to leave to agents**.

If the user wants to own direct component interactions, create an explicit interaction-model layer and settle those interactions before implementation.

If the user does not care to specify that level, omit the layer and let downstream agents resolve it within the constraints of settled upstream artifacts.

A useful principle is:

> **Agent discretion is the set of decisions intentionally left unresolved by settled upstream artifacts.**

This makes the system useful even for teams that do not want to “outsource development to agents.” They can choose the resolution at which agents get discretion.

### 3.2 Execution DAG

Separately, there is the workflow DAG that advances one artifact.

Example:

```text
analyze frontier
      │
      ▼
 ┌────┼────┐
 ▼    ▼    ▼
Exp A Exp B Exp C
 └────┼────┘
      ▼
   review
      │
      ▼
 human gate
      │
      ▼
synthesize / implement
      │
      ▼
 verify / settle
```

This DAG is operational. It says which nodes run, which models they use, which commits they start from, which outputs they can see, budgets, stopping conditions, and approval gates.

The UI can expose and manipulate this DAG directly.

The artifact graph answers **what depends on what**.

The execution DAG answers **how we are currently trying to advance one part of it**.

---

## 4. Deterministic orchestrator, not a super-agent

The top-level manager should primarily be a **deterministic application/workflow engine**.

It should not be another unconstrained agent that simply receives a broad prompt such as:

> We hit problem X; launch three workers, evaluate them, choose one, and continue.

That would make the experiment controller capable of silently changing scope, leaking information between experiments, or steering workers in ways that undermine independence.

Instead, the deterministic orchestrator owns things such as:

- DAG structure,
- node types,
- parent commits,
- artifact versions,
- budgets,
- selected models / reasoning levels,
- repository/worktree isolation,
- allowed information edges,
- stopping conditions,
- human approval gates,
- which nodes are allowed to mutate settled state,
- and when dependent artifacts become ready for reconciliation.

Agents operate *inside* nodes.

This gives a useful trust boundary:

> Agents may reason freely within a bounded task. The application controls what tasks exist, what inputs they receive, and what transitions are legal.

---

## 5. The worker primitive

Ralph-like agents are **worker primitives**, not the whole autonomous system.

A worker can:

- load repository state,
- read the relevant documentation/backlog,
- identify the next concrete missing capability,
- implement and test it,
- persist local work state,
- discover blockers,
- propose questions,
- and make substantial coherent progress.

The wider system adds the capabilities that a single worker fundamentally lacks:

- comparison against alternative solutions,
- independent evaluation,
- branch-level validity monitoring,
- synthesis across experiments,
- provenance,
- selective reuse,
- human-controlled settlement,
- and cross-artifact change propagation.

---

## 6. Initial analysis: determine what is enough to settle a baseline

A project should not hardcode “Scenario 01” or “the first feature” as the universal definition of a settled foundation.

Before heavy materialization starts, an analysis/review step should inspect the backlog/specification and propose the **foundation frontier**:

> What is the smallest representative vertical slice that exercises enough of the project’s important assumptions to establish a useful architectural baseline?

The proposal should include:

- the candidate milestone,
- why it is representative,
- which architectural assumptions it will pressure,
- what evidence will demonstrate that it works,
- which behaviors are deliberately deferred,
- and possibly what additional small amount of follow-on pressure is needed before declaring the foundation stable.

Examples:

### TODO app

Possible baseline:

- launch the real HTTP server,
- initialize real persistent state,
- allow self-sign-up,
- allow login,
- demonstrate the relevant persistence/restart behavior,
- explicitly defer TODO CRUD/listing.

This establishes much of the basic system grammar without implementing the full backlog.

### Traust test harness

Possible baseline:

- Scenario 01 executes via the public suite entrypoint,
- a separately replaceable SUT occupies the SUT role,
- real required infrastructure boundaries are exercised,
- evidence arrives through the same semantics available to an independent implementation,
- pass/fail behavior is correct.

Scenario 02 might optionally be used as a small amount of pressure if needed to show the Scenario 01 architecture is not immediately brittle, but this is an analysis decision rather than a universal rule.

The proposed foundation frontier should itself be reviewed/approved.

---

## 7. Backlog/frontier-driven settlement

The outer system should reason in terms of a **settlement frontier** rather than “let this agent run until its token budget expires.”

For Traust, the natural units are scenarios.

```text
settled through S00
       │
       ├── S01 experiment A
       ├── S01 experiment B
       └── S01 experiment C
              │
              ▼
            review
              │
              ▼
          settle S01
              │
              ▼
          begin S02
```

An experimental worker that declares the current frontier complete should generally stop rather than continue into the next backlog unit.

Reason:

- the experiment stays interpretable,
- the branch contains changes attributable to one frontier,
- a good S01 cannot be followed by weird speculative S02/S03 choices that contaminate evaluation,
- the reviewer can compare like with like,
- the synthesis/settlement decision is cleaner.

This does cost context-reconstruction tokens when the next frontier begins. That may be worthwhile.

This does **not** necessarily mean one Ralph iteration per scenario. A worker can still use multiple substantial internal iterations to complete one outer frontier.

Later, when the system is mature, the orchestrator may permit grouping low-risk backlog units. That is an optimization, not the conceptual default.

---

## 8. Experiment independence

Experiment nodes should be genuinely independent unless the DAG explicitly says otherwise.

A typical experiment definition includes:

- common parent commit,
- common settled artifact state,
- common frontier / requirements,
- selected model,
- reasoning level,
- budget,
- tool/environment policy,
- no visibility into sibling experiment outputs.

The purpose is not merely redundancy. Independent experiments perform **architectural search**.

Different workers may:

- discover different valid approaches,
- expose ambiguity,
- find implementation constraints,
- overbuild,
- underbuild,
- or discover that an apparent shortcut violates the intended boundary.

The outer system extracts value from that diversity.

---

## 9. Manager/monitor must not steer live experiments

The outer manager should not tell a running experiment:

> You are choosing the wrong architecture; try a registry proxy instead.

That contaminates the experiment and collapses independent search.

However, it should be able to stop wasting compute when the experiment has crossed an **objective invalidation boundary**.

Example:

- worker marks Scenario 01 covered,
- monitor mechanically evaluates the accepted Scenario 01 validity gate,
- branch does not actually exercise a replaceable external SUT,
- experiment is stopped / marked invalid.

The monitor answers:

> **Is this still a valid experiment under already-settled rules?**

It does not answer:

> **What architecture should the worker choose?**

This distinction should be explicit.

Automatic termination should be conservative and based on settled, objective conditions, not architectural taste.

---

## 10. Explicit situations the manager/reviewer should look for

The review layer should not perform an unbounded “find anything interesting” critique. A useful design is to define explicit classes of signals that warrant escalation, synthesis, or human review.

### 10.1 Blocked by contradiction

A worker reports that requirements appear impossible to satisfy simultaneously.

Example:

> Requirement A appears incompatible with Requirement C because satisfying A requires behavior explicitly forbidden by C.

This is often detected locally by the implementation worker.

### 10.2 Valid divergence / ambiguous interpretation

Two or more independent experiments satisfy the apparent requirements, but they interpret an upstream sentence or concept differently.

This is more interesting because no worker necessarily thinks it has failed.

Example:

- Experiment A interprets “publish state before candidate start” as persisted-state publication.
- Experiment B interprets it as an API acknowledgment.
- both can construct plausible valid systems.

The comparative reviewer detects the divergence and raises an upstream clarification candidate.

This is a particularly valuable use of parallel experiments.

### 10.3 Hidden consequence

Implementation is possible, but doing it introduces a material property the upstream artifact does not acknowledge.

Possible examples:

- an unexpected persistence semantic,
- a security/trust boundary,
- an operational requirement,
- a consistency model,
- a lifecycle assumption,
- an ordering guarantee,
- a latency/availability tradeoff.

This is a broad/systemic question and should **not** become a requirement that every ordinary implementation worker perform a complete architectural critique.

It belongs especially to **validation experiments** or dedicated review modes where broad questions are explicitly in scope.

### 10.4 Repeated workaround

Multiple workers independently compensate for the same awkward or missing upstream concept.

This can indicate:

- missing specification,
- missing abstraction,
- or an upstream rule that is technically satisfiable but poorly aligned with implementation reality.

### 10.5 Boundary violation

A branch appears to progress only by violating an already-settled product or artifact boundary.

This is a good candidate for objective monitoring/termination if the violation is mechanically recognizable.

### 10.6 High architectural variance

Several valid experiments solve the same new problem in materially different ways.

That does not mean one is wrong. It means the choice may be consequential enough to deserve explicit settlement instead of being inherited accidentally from whichever worker happened to run first.

### 10.7 Suspiciously cheap progress

Not a validity rule by itself, but unusually cheap/deep progress compared with peers may be a useful trigger for closer review, especially when prior experiments show that skipped boundaries can create misleading throughput.

### 10.8 Repeated inability to make progress

If several independent workers fail at the same frontier for different reasons, that may indicate a specification gap, environmental dependency, or foundational issue rather than a model-specific failure.

---

## 11. Two very different experiment purposes

A major distinction emerged between **materialization experiments** and **validation experiments**.

### 11.1 Materialization experiments

Goal:

> Find a good way to advance the current artifact frontier.

Outputs may feed the same artifact’s settled state.

Possible outcomes:

- continue a particularly coherent branch,
- transplant a mechanism,
- give a synthesis worker selected references,
- or use learnings only and rebuild coherently from the common base.

### 11.2 Validation experiments

Goal:

> Use attempted materialization as a probe to validate some upstream artifact or documentation.

These branches are deliberately throwaway.

Example:

- documentation currently defines A/B/C,
- before investing in D…N, run implementation/test-harness spikes against A/B/C,
- ask whether A/B/C are implementable, ambiguous, contradictory, or hiding major consequences,
- summarize findings back into the documentation project.

The validation experiments are **not** meant to produce code to merge.

Their result is a report/message to another project/scope.

This should be represented explicitly in the workflow so a validation campaign does not accidentally turn into normal product development.

---

## 12. Specification development can use throwaway derivations as validation

This creates an interesting alternative to both pure waterfall and pure “code first” development.

Suppose the team wants to specify A…N.

Writing all of A…N before implementation is risky because practical limitations may invalidate assumptions.

Instead:

```text
write A / B / C
      │
      ▼
focused textual review
      │
      ▼
throwaway validation campaign
  ├── implementation spike(s)
  ├── test-harness spike(s)
  └── perhaps other derived artifacts
      │
      ▼
findings about ambiguity / contradiction / consequences
      │
      ▼
human revises A / B / C
      │
      ▼
continue with D / E / ...
```

The code may be disposable; the findings are not.

This is especially powerful when **multiple different downstream surfaces derive independently from the same documentation**.

For Traust:

- the actual Traust implementation,
- and the Traust test harness,

can both act as independent validation surfaces.

If they independently interpret the same documentation differently, that is a strong ambiguity signal.

Focused document review and attempted materialization catch different problem classes and can complement each other.

---

## 13. Experiment outputs: knowledge first, code optionally

During early architectural exploration, the primary product of experiments may be **knowledge**, not code.

A reviewer might synthesize:

```text
X:
  Experiment A's Docker-boundary interception is promising.
  Essential principle: ...
  Avoid coupling to ...

Y:
  Experiment B's persisted-state observer is simpler.
  Essential principle: ...

Z:
  Experiment C demonstrated assumption Z does not hold because ...
  Treat this as a constraint, not as an implementation recommendation.

Rejected:
  Experiment D bypassed the public SUT boundary and therefore does not establish valid coverage.
```

A fresh synthesis worker can then materialize one coherent system from the common base plus these learnings.

This avoids choosing a “winning branch” when:

- A solved X well but Y poorly,
- B solved Y well but X poorly,
- C exposed a false assumption,
- no branch is globally best.

It also avoids merging all branches into a Frankenstein architecture.

---

## 14. But do not collapse experiment evidence into prose and discard the source

Even if experimental branches are not merged, keep them as immutable evidence attached to the DAG.

A reviewer statement such as:

> Experiment B found a promising candidate-create interception mechanism.

should retain provenance to:

- experiment node,
- source commit,
- affected files/hunks,
- relevant tests,
- logs/evidence,
- and the reviewer reasoning that extracted the learning.

The normal synthesis context can remain compact while allowing the human/reviewer to drill into the underlying source.

This avoids turning reviewer prose into an information-loss bottleneck.

---

## 15. Selectable policies for using experiment code

There should not be one global rule such as “always start fresh” or “always continue the winner.”

The orchestrator could support several explicit synthesis/input policies.

### Learnings only

```text
experiments
    ↓
review / distilled learnings
    ↓
fresh synthesis from common base
```

Useful for maximizing coherence and testing whether the distilled knowledge is sufficient.

### Reference

The synthesis worker receives distilled learnings plus links/permission to inspect selected commits.

```text
learned principle
    +
selected reference commit
    ↓
fresh coherent implementation
```

Useful when a mechanism is genuinely excellent and re-inventing it would waste effort.

### Transplant

A specific commit/hunk/mechanism is explicitly ported or cherry-picked into the new candidate, possibly by an adaptation node.

This preserves code directly but makes the reuse explicit and reviewable.

### Continue

A particularly coherent experiment branch becomes the base for the next frontier.

This may be appropriate later when the settled architecture is mature and experiments are narrow.

These policies themselves are experimentable. The system can compare whether “learnings only” vs “reference” vs “transplant” produces better coherence/cost.

---

## 16. Human annotations and PR-like review as first-class provenance

A future UI could allow comments attached to:

- commits,
- files,
- hunks/lines,
- reviewer learnings,
- architectural decisions,
- settlement proposals.

Examples:

> Use `variable_name` here instead of `varname`.

> Good principle of X, but the mechanism should be abstracted into package `foo`.

> X is a useful abstraction, but logic should therefore be centralized around `foo`, not `bar`.

These comments can themselves become inputs to structured workflows.

### “Discuss solutions” on selected lines

The user selects code/learning lines and asks:

> X should be in `foo` instead of `bar`.

The UI could offer standard workflows such as:

- **Possible implementations** — branch several implementation alternatives and render code-preview comparisons.
- **Repercussions / impact analysis** — analyze consequences on components, interfaces, tests, dependencies, migration, etc.
- **Revise learning** — ask a reviewer to reformulate the distilled architectural lesson in light of the comment.
- **Free-form discussion** — let a routing/manager agent choose the best available representation/workflow, falling back to text if nothing structured fits.

The important point is that the comment and all derived analyses remain connected by provenance.

---

## 17. Reviewer learnings are also versioned/reviewable

The reviewer is not infallible.

Example:

```text
Experiment evidence
      ↓
Learning v1:
  "X is a good abstraction; centralize around bar"
      ↓
Human annotation:
  "X is useful, but should live in foo"
      ↓
Re-analysis
      ↓
Learning v2:
  "Keep X, centralize around foo"
      ↓
Human approval
```

The system should preserve the original evidence and both learning revisions.

This makes “manager synthesis” itself auditable rather than magical.

---

## 18. Human approval boundaries

The system should distinguish ordinary autonomous progress from changes to the rules of the game.

In particular, **prompt/guideline/specification changes should generally require human approval**.

If several experiments expose a loophole, the reviewer may propose:

> Existing instruction allows interpretation X. Proposed clarification Y.

But the system should stop at a review gate rather than silently updating project instructions.

This is particularly important during project bootstrap, when many loopholes may be found quickly.

A likely pattern for a new project is:

```text
experiment cohort
      ↓
review
      ↓
several instruction/spec loopholes found
      ↓
human approves/edits clarifications
      ↓
new cohort
      ↓
fewer loopholes
      ↓
...
```

The decreasing rate of proposed foundational corrections may itself be a useful maturity signal.

---

## 19. Autonomy should mean widening deterministic auto-approval, not making a super-agent sovereign

A future user might say:

> Auto-approve workflow classes A and C.

That should not mean “let the manager do whatever it thinks is equivalent to A/C.”

It means the deterministic workflow engine allows certain known transition types without pausing for human approval.

Example policy:

```text
auto-approve:
  - ordinary backlog continuation
  - experiment fan-out up to 3 nodes
  - routine reconciliation with no upstream findings

require approval:
  - prompt/project-guideline changes
  - upstream specification reinterpretation
  - new test-only/public interface
  - reopening settled architectural decisions
  - widening experiment scope
  - major re-architecture
```

As trust grows, the user can widen the set of automatically permitted transitions.

The deterministic application remains the authority.

---

## 20. Lifecycle stages: useful defaults, not a rigid global state machine

Three broad modes were identified:

### 20.1 Baseline / foundation discovery

Goal:

> Establish the first representative vertical slice and settle the general grammar of the artifact.

Characteristics:

- heavy parallel exploration,
- strong comparative review,
- likely more human gates,
- willingness to discard branches,
- greater sensitivity to architectural mistakes,
- specification/prompt loopholes expected to surface.

### 20.2 Incremental buildout / settlement

Goal:

> Advance the backlog frontier using the settled baseline.

Characteristics:

- usually one coherent trunk,
- experiments around consequential forks,
- monitor adherence to established principles,
- latent foundational problems may still surface,
- individual backlog units settle progressively.

### 20.3 Maintenance / reconciliation

Goal:

> React to upstream changes while preserving accumulated accepted state.

Characteristics:

- normally incremental,
- strong emphasis on impact/reconciliation,
- avoid unnecessary architectural churn,
- experiments are targeted around genuinely new or risky changes.

However these should be **loose/descriptive modes**, not irreversible states for the entire system.

Each artifact/project can be in its own lifecycle posture.

Architecture might be mature while a new subsystem design is in baseline discovery.

The top-level system can show those statuses, but work is performed within each project/artifact scope.

---

## 21. Maintenance can still launch heavy local exploration

A mature project receiving a major upstream change should not have to “return to foundation discovery” globally.

Example:

- settled downstream decisions X and Y exist,
- architecture revision A13 introduces new requirement Z,
- user/manager decides Z is consequential,
- launch several experiments exploring how to accommodate Z while preserving X/Y.

This is a **scoped exploration inside maintenance**.

A useful explicit operation could be:

> **Challenge current solution based on A13**

Scope:

- current downstream layer only,
- focus on whether existing settled decisions still work under Z,
- do not turn the exercise into a broad critique of upstream architecture,
- but allow concrete upstream findings to surface if necessary.

Possible outputs:

```text
Preserve:
- Y
- Q
- public interface

Reopen proposal:
- X, because A13/Z changes its assumptions

Upstream finding:
- none
```

or:

```text
Preserve:
- X
- Y

Result:
- Z fits without changing settled downstream decisions
```

or occasionally:

```text
Potential upstream issue:
- satisfying Z while preserving required X appears incompatible with A13
```

Nothing is reopened automatically. The output goes through human review.

---

## 22. Scoped reopening instead of wiping settled state

“Go back to foundation discovery” is too destructive as a mental model for mature artifacts.

If a deep issue appears, explicitly identify which decisions are reopened and which remain settled.

```text
CURRENT SETTLED STATE
      │
      │ reopen:
      │   - routing strategy
      │   - deployment observation mechanism
      │
      │ preserve:
      │   - public API
      │   - data model
      │   - unrelated accepted decisions
      ▼
parallel exploration
      │
      ▼
review / proposal
      │
      ▼
human settlement
```

If evidence shows the problem reaches deeper, the human can widen the reopened scope.

Even an extreme foundational rework should start with the previous system’s learnings, provenance, accepted facts, and known failures rather than pretending to be a blank-slate project.

---

## 23. Authority flows downward; evidence can flow upward

Artifact dependencies should be asymmetric.

### Forward direction

```text
upstream settled artifact
        │
        └── constraints/change ──► downstream reconciliation
```

A settled architecture change can automatically mark dependent artifacts stale and schedule reconciliation work.

### Reverse direction

```text
downstream work
      │
      └── findings/proposals ──► upstream review
```

Downstream work must not silently rewrite upstream truth.

It can report:

- ambiguity,
- contradiction,
- infeasibility,
- hidden consequence,
- repeated workaround,
- request for clarification,
- or a proposed change.

A human review/settlement process decides whether the upstream artifact changes.

This preserves hierarchy without assuming humans can predict all implementation realities in advance.

---

## 24. Reconciliation, not regeneration

Suppose:

```text
Architecture A12

Design D8
  based_on: A12

Test Harness T19
  based_on: A12

Implementation I27
  based_on:
    Design D8
    Test Harness T19
```

Architecture advances to A13.

The system should not delete D8/T19/I27 and regenerate everything.

Instead:

```text
Architecture A13       SETTLED
Design D8              STALE vs A13
Test Harness T19       STALE vs A13
Implementation I27     EXISTING SETTLED STATE, NOT YET READY TO RECONCILE
```

Design and test harness may reconcile independently:

```text
A13 ──► Design reconciliation
 │
 └────► Test Harness reconciliation
```

Suppose Design becomes D9 based on A13 first. Implementation still waits because Test Harness remains based on A12.

Once Test Harness becomes T20 based on A13, the implementation becomes ready:

```text
Implementation I27
   inputs now settled:
     Design D9 @ A13
     Test Harness T20 @ A13
   => READY FOR RECONCILIATION
```

This resembles a reactive dependency/build graph, but with **persistent materialized state** rather than disposable render output.

The system preserves reviewed nuance and reconciles deltas.

---

## 25. “Consistent cut” view of system state

Because different artifacts may temporarily be based on different upstream revisions, the UI can distinguish:

- latest settled revision of each artifact,
- stale artifacts,
- artifacts currently reconciling,
- artifacts blocked waiting for sibling dependencies,
- the latest fully consistent set of versions.

This “consistent cut” can matter operationally.

For example, the latest architecture may be A13 while production implementation still corresponds to the last fully settled A12-derived graph. The system should make that explicit rather than implying everything is automatically current.

---

## 26. Cross-project dependency messages, not one global agent context

Each artifact/project is worked on in isolation with its own backlog, state, experiments, and lifecycle.

The top-level system does not directly “do architecture + test harness + implementation” in one context.

Projects communicate through explicit dependency/messages:

```text
upstream constraints/change  ──► downstream project
upstream revision settled    ──► downstream stale/readiness status

downstream findings/proposal ──► upstream review
```

Example:

1. Architecture A13 settles.
2. Test Harness project is marked stale and begins its own reconciliation workflow.
3. Design project is marked stale and begins its own reconciliation workflow.
4. Implementation does not start yet because its upstream inputs are not mutually current.
5. Design D9 settles.
6. Test Harness T20 later settles.
7. Implementation is now ready to reconcile against the complete current dependency set.

The top level can visualize this but should not collapse all work into one mega-agent.

---

## 27. Organic buildup is a feature

The system should not assume “tests are mechanically derived from specs, therefore any downstream artifact is safely reproducible from scratch.”

The desired strategy is an **organic buildup of settled state**.

Over time, actual code/test/design structure embodies many reviewed choices and nuanced consequences.

Those choices may not all be restated verbatim upstream.

Therefore:

- preserve accepted downstream artifacts,
- reconcile them against changes,
- use scoped reopening when necessary,
- retain learnings and provenance,
- do not casually wipe mature state.

This is why later maintenance can be much more stable than initial foundation work: the repository itself becomes a rich constraint that makes old bootstrap mistakes less likely.

Example from the Traust experiments:

- early branches could misunderstand the test-suite/SUT product boundary,
- once a correct public runner/SUT architecture has been settled and accumulated tests rely on it, later agents are much less likely to reproduce the old mistake even if the prompt wording is not infinitely exhaustive.

---

## 28. Avoid fossilizing every model mistake into prompts

As experiments mature, failures become increasingly niche.

Not every bad branch should result in another permanent prompt rule.

The review process should classify findings roughly as:

### General invariant

Repeatedly important across the project. Candidate for durable project/workflow instructions after human approval.

### Architectural lesson

Useful within a specific subsystem/artifact. Better recorded as a settled design/guideline/decision than as a universal agent rule.

### One-off trajectory mistake

Reject the branch. Do not memorialize it.

This keeps prompts from becoming a fossil record of every odd path any model has ever taken.

The outer search/review/settlement loop reduces pressure to make the worker prompt perfect.

---

## 29. Search vs exploitation

A central responsibility of the outer workflow is deciding **when independent search is worth the compute**.

### Early foundation

Search heavily because choosing a bad foundation is expensive.

### Buildout

Exploit the settled architecture by default. Branch selectively when a consequential new fork appears.

### Mature maintenance

Search is exceptional, usually around:

- major upstream changes,
- high-risk new requirements,
- contradictory evidence,
- performance/security boundaries,
- or a challenge/reopen proposal.

This is a more useful framing than “make the manager smarter.”

---

## 30. Potential manager-triggered workflows

Some useful explicit workflow operations discussed or implied:

### Establish baseline

Analyze the backlog and propose the minimum representative vertical slice that can settle a foundation.

### Explore frontier

Launch independent experiments from the same settled state for the current backlog item.

### Review cohort

Compare validity, architecture, complexity, evidence, and learnings across experiments.

### Validate specification

Run intentionally throwaway materializations to find ambiguity, contradiction, hidden consequences, and implementation friction in an upstream draft.

### Challenge current solution based on upstream change

Given A13/Z, test whether settled downstream X/Y remain appropriate.

### Reconcile artifact

Update an existing settled downstream artifact against changed upstream dependencies while preserving unaffected accepted state.

### Reopen decisions

Propose a scoped set of settled facts to reconsider; require review/approval.

### Discuss solutions

From a code/review annotation, branch implementation alternatives.

### Analyze repercussions

From a selected proposed change/comment, analyze system/test/API/dependency consequences without implementing it.

### Revise learning

Re-run synthesis around existing evidence after human feedback on a reviewer conclusion.

### Expand experiment

User adds more parallel nodes/models/budgets to an existing experiment point in the DAG.

---

## 31. The DAG UI as the central control surface

The DAG should not merely visualize hidden orchestration. It can be the user’s programming/control model for the system.

A user could:

- inspect each node,
- see parent revision and inputs,
- see model/reasoning/budget,
- duplicate an experiment with another model,
- add/remove branches before execution,
- change which nodes are reviewers vs workers,
- approve/reject settlement,
- choose synthesis policy,
- inspect provenance,
- comment on commits/hunks,
- add a new analysis branch from any point,
- or configure auto-approval policies.

Example node detail:

```text
Experiment: foundation-07/B
Parent: 3c4f...
Artifact: test-harness
Frontier: Scenario 01
Model: GPT-X
Reasoning: high
Budget: ...
Prompt revision: ...
Inputs:
  - repository @ 3c4f...
  - docs @ ...
Isolation:
  - no sibling experiment outputs

Status:
  STOPPED BY MONITOR

Reason:
  Scenario 01 marked covered
  acceptance gate failed

Outputs:
  - immutable branch snapshot
  - agent state
  - logs
  - monitor report
  - reviewer annotations
```

---

## 32. Information flow should itself be visible

An important possible UI/property is that **context sharing is represented as DAG edges**.

Examples:

- Experiment B does not see Experiment A because there is no edge.
- Reviewer R receives all A/B/C outputs because edges connect them.
- Synthesis S receives the review summary plus only selected reference commit B2 because those are the incoming edges.
- Validation report V sends findings to the Architecture project, not code.

This makes experiment contamination and provenance inspectable.

---

## 33. Model choice is an explicit workflow parameter

The user should be able to control models per node.

Examples:

- 2× Terra + 1× Sol for implementation exploration,
- Astra for meticulous reviewer/reference work,
- cheaper model for routine maintenance,
- expensive model only for a challenge/reopen analysis,
- duplicate the same node across reasoning levels.

Model diversity can be used intentionally as part of architectural search rather than hidden behind one global “manager model.”

---

## 34. Budgets should be explicit

Experiment nodes should have explicit budgets/limits such as:

- token budget,
- maximum worker iterations,
- wall-time cap,
- scenario/frontier completion stop,
- objective invalidation stop.

The exact budget unit remains an open design problem; prior Traust experiments demonstrated that raw iteration count is not a good compute measure.

The important thing is that budgets are an orchestrator-level experimental property, not something workers silently choose.

---

## 35. Comparison should judge the product, not trust worker progress claims

Reviewers/monitors should evaluate actual output against accepted gates.

For Traust, that meant inspecting whether the claimed Scenario 01 coverage actually used the public suite and replaceable SUT boundary, not simply accepting a worker’s coverage metadata.

General principle:

> Agent-reported progress is evidence to inspect, not authoritative truth.

Useful comparison dimensions can include:

- validity against hard requirements,
- backlog depth reached,
- correctness/coverage quality,
- tokens/wall time,
- architectural complexity,
- debt/workarounds,
- introduced interfaces/contracts,
- unresolved blockers,
- coherence with settled upstream state,
- and quality of evidence.

---

## 36. Strong early branches can become references without becoming truth

If one experiment contains an exceptionally good mechanism, later workers may be allowed to inspect it directly.

This is different from declaring the whole branch the winner.

For example:

```text
Reviewer:
  "Experiment A's X is excellent; Experiment B's Y is excellent."

Synthesis inputs:
  - common base
  - learning summary
  - reference to A:X commit/hunk
  - reference to B:Y commit/hunk
```

This allows selective reuse while still asking the synthesis worker to construct a coherent whole.

Whether this produces better systems than pure learnings-only reconstruction is itself an interesting experiment.

---

## 37. Baseline settlement can reduce the need for future prompt hardening

Early bootstrap mistakes matter most because there is little existing structure to constrain the worker.

Once a good baseline is settled, the repository embodies many architectural decisions mechanically.

Thus a loophole like:

- “technically separate SUT process but actually scenario-specific self-test”

may be important to detect during foundation exploration but may never recur after the correct public runner/SUT shape has become established and heavily exercised.

Therefore a system with parallel foundation search + review may need fewer permanent textual guardrails than a single-worker system that must get every bootstrap interpretation right on the first try.

---

## 38. Manager output categories

It may be useful to distinguish what a reviewer/manager node is allowed to emit.

Possible structured categories:

### Validity findings

- branch valid/invalid under settled gate,
- exact evidence,
- reason for rejection.

### Comparative findings

- approach A vs B,
- complexity differences,
- performance/correctness tradeoffs,
- architectural variance.

### Learnings

- reusable principle,
- promising technique,
- rejected technique,
- context/limits.

### Upstream findings

- ambiguity,
- contradiction,
- hidden consequence,
- missing requirement,
- proposed clarification.

### Settlement proposal

- preserve these decisions,
- adopt these new decisions,
- reopen these decisions,
- defer these questions.

### Workflow recommendation

- another experiment cohort is warranted,
- enough evidence exists to synthesize,
- validation mode should run,
- human review required.

These structured outputs could make the deterministic orchestrator much easier to reason about than parsing free-form prose.

---

## 39. Hard distinction: findings vs mutations

An agent can discover or propose nearly anything.

It should not automatically mutate settled authority merely because it found something.

Examples:

- worker discovers contradiction → finding,
- reviewer proposes spec clarification → proposal,
- challenge workflow suggests reopening X → proposal,
- validation experiments reveal ambiguity → upstream message,
- human approves → mutation/settlement.

This separation is one of the main safety/coherence properties of the whole system.

---

## 40. Potential maturity signals

Possible indicators that a project/artifact is becoming mature:

- fewer foundational guideline/spec clarification proposals,
- fewer invalid experiment branches,
- lower variance between valid implementations,
- more frontiers completed by normal single-trunk progression,
- less need for heavy review cohorts,
- changes are mostly reconciliations rather than re-architecture,
- objective monitors rarely stop workers,
- new upstream changes fit established structures without reopening many decisions.

These are not necessarily formal metrics initially, but could eventually inform recommended workflow intensity.

---

## 41. Open question: when should the system automatically suggest more experiments?

The orchestrator might recommend fan-out based on signals such as:

- large upstream change,
- reopened high-impact decision,
- reviewer sees several credible alternatives,
- previous worker is blocked on architectural uncertainty,
- high variance in earlier related work,
- security/performance/operationally critical frontier,
- or user explicitly requests challenge mode.

But automatic fan-out should probably remain an explicit deterministic policy/rule, not an unconstrained manager intuition that silently burns compute.

The user can approve, edit, or auto-approve such recommendations.

---

## 42. Open question: when is a baseline “settled enough”?

The initial baseline analysis should propose this, but there may eventually be reusable criteria:

- first representative end-to-end behavior works,
- key external boundaries are real,
- basic lifecycle/persistence model has been exercised,
- public interfaces are taking shape,
- at least one negative/failure path is understood,
- a small amount of next-step pressure does not immediately invalidate the foundation,
- comparative experiments have converged enough that remaining variance is local rather than foundational.

The exact standard should remain project-specific.

---

## 43. Open question: how much of the review itself should be deterministic?

Possibilities:

- mechanical hard gates first,
- then one or more reviewer agents,
- perhaps a comparative judge that must cite concrete branch evidence,
- possibly multiple independent reviewers for high-impact settlements,
- user can inspect disagreements.

Important principle:

> Deterministic checks should handle objective validity where possible; LLM reviewers should handle comparative/semantic judgment where needed.

---

## 44. Open question: branch retention and storage

Even disposable experiments should likely be retained at least as immutable references because:

- reviewer learnings need provenance,
- users may inspect them later,
- a later agent may need to examine a particularly good mechanism,
- historical failures may explain why a settled rule exists.

Possible long-term storage policies can be handled separately from conceptual provenance.

---

## 45. Open question: how to represent settled decisions

Possible forms:

- repository code/tests themselves,
- architecture/design docs,
- explicit ADR-like records,
- structured decision objects,
- reviewer settlement reports,
- version metadata linking artifact revision to upstream revision.

Likely answer: a combination.

The key is that downstream agents can identify what is actually settled and distinguish it from experimental commentary.

---

## 46. Open question: how to represent messages between projects

Potential message types:

- `upstream_changed`
- `reconciliation_ready`
- `blocked_waiting_for_dependency`
- `ambiguity_found`
- `contradiction_found`
- `hidden_consequence_found`
- `change_proposal`
- `clarification_requested`
- `settled_revision_available`

These could be structured objects with links to exact source revisions/evidence.

---

## 47. Open question: what belongs in prompts vs settled artifact guidance?

The Traust experiments suggest a useful division:

- prompts/workflow rules should contain a small number of high-value generic invariants,
- artifact-specific architecture/design should live in project artifacts,
- one-off mistakes should simply be rejected,
- new prompt/guideline rules should require human approval.

The outer manager reduces the need to encode every possible failure prevention rule in the worker prompt.

---

## 48. Example: Traust test-harness foundation workflow

A possible concrete flow under this system:

```text
Traust documentation/scenario backlog
        │
        ▼
Foundation analysis
  proposed frontier: Scenario 01
        │
        ▼
Human approval
        │
        ▼
 ┌────────────┬────────────┬────────────┐
 │ Terra exp  │ Sol exp    │ Astra exp  │
 │ from C0    │ from C0    │ from C0    │
 └────────────┴────────────┴────────────┘
        │          │          │
        ├── monitor hard Scenario-01 validity gate
        │
        ▼
Comparative review
  - reject invalid boundary interpretations
  - identify good registry strategy
  - identify good SUT launch pattern
  - identify ambiguity/loopholes
  - retain provenance to exact commits
        │
        ├── proposed project-rule change? ─► HUMAN GATE
        │
        ▼
Synthesis policy chosen
  e.g. learnings + selected reference commits
        │
        ▼
Fresh coherent Scenario-01 materialization
        │
        ▼
verification
        │
        ▼
SETTLE S01 BASELINE
        │
        ▼
Begin S02 frontier
```

If a later architecture change introduces Z:

```text
Architecture A13/Z
        │
        ▼
Test Harness becomes stale
        │
        ▼
"Challenge current solution based on A13"
  preserve settled X/Y initially
        │
        ├── exp 1
        ├── exp 2
        └── exp 3
        │
        ▼
review
  - X still works? Y still works?
  - should X be reopened?
  - any concrete upstream issue?
        │
        ▼
human approval
        │
        ▼
reconciliation / scoped reopening
```

---

## 49. Example: TODO application foundation workflow

Given a rich TODO app spec:

```text
Backlog includes:
- service boot
- persistence
- signup/login
- TODO CRUD
- sharing
- reminders
- etc.
```

Foundation analyzer proposes:

```text
Foundation milestone:
- real HTTP API boots
- DB/state initializes
- user self-signup works
- login works
- auth/session survives intended lifecycle

Deferred:
- TODO listing
- TODO CRUD
- sharing
- reminders
```

Several experiments explore the basic system architecture.

Review settles one coherent foundation.

Only then does the backlog frontier move to TODO functionality.

This reduces the risk that each later feature reopens basic choices about HTTP lifecycle, persistence, auth, or state ownership.

---

## 50. Example: spec-validation campaign

Suppose only documentation A/B/C exists.

Goal is not to keep code.

```text
Docs A/B/C
   │
   ▼
Validation campaign
   ├── implementation spike 1
   ├── implementation spike 2
   ├── independent test-harness spike
   └── focused docs reviewer
   │
   ▼
Validation reviewer asks explicitly:
   - Any contradictions?
   - Any valid divergent interpretations?
   - Any repeated workarounds?
   - Any material hidden consequences?
   - Any requirement impossible without introducing undeclared behavior?
   │
   ▼
Report to documentation project
   │
   ▼
Human edits docs
```

No experimental branch is promoted.

This is an intentionally expensive sanity check on the documentation while it is still being built.

---

## 51. Example: optional interaction-model layer

The artifact graph might be:

```text
Architecture
    │
    ▼
Component Design
    │
    ▼
Interaction Models
    │
    ▼
Implementation
```

The interaction-model artifact could contain:

- Mermaid sequence diagrams,
- state transitions,
- component-to-component workflow definitions,
- event ordering,
- request/response contracts.

If the user wants to own this engineering level, downstream implementation agents treat it as settled authority.

If not, remove that layer and let implementation agents resolve interactions themselves.

This illustrates that the system is a **control-resolution framework**, not necessarily a “let agents do all engineering” framework.

---

## 52. Possible relationship to interfaces/header-like artifacts

A team might choose to settle interfaces without implementations.

This could resemble:

- C header files,
- IDLs,
- OpenAPI schemas,
- RPC schemas,
- TypeScript interfaces,
- language-neutral contracts,
- component boundary documents.

The exact representation is secondary.

The important idea is that downstream agents inherit explicit settled boundaries while retaining freedom over internals.

---

## 53. Non-goals / cautions

### Not “one top-level agent does everything”

The manager should not silently invent DAG structure, alter experiment inputs, or merge branches without explicit workflow rules/approval.

### Not perfect prompt engineering

The goal is not to make the Ralph prompt immune to every conceivable model mistake.

Parallel search, objective monitors, comparative review, and settlement provide a higher-level correction mechanism.

### Not automatic upstream rewriting

Downstream evidence can trigger proposals; upstream authority changes only through explicit settlement.

### Not regenerate-everything-on-change

Preserve reviewed state and reconcile deltas.

### Not “pick the winning experiment” by default

Different experiments can contribute different lessons/mechanisms.

### Not merge-everything

Avoid Frankenstein architectures produced by indiscriminately combining all successful experiments.

### Not broad critique in every worker

Ordinary workers should solve the current frontier and report concrete friction. Broad consistency/consequence questions belong to dedicated validation/review workflows.

### Not rigid lifecycle rollback

A mature artifact can run scoped exploration without declaring its entire settled state invalid.

---

## 54. Emerging high-level architecture

A compact representation of the envisaged full system:

```text
                    HUMAN / GUI
                        │
      approve / annotate / alter workflows / set policies
                        │
                        ▼
             DETERMINISTIC ORCHESTRATOR
                        │
       ┌────────────────┼────────────────┐
       │                │                │
       ▼                ▼                ▼
  worker nodes      monitor nodes    reviewer nodes
 (Ralph-like)       (non-steering)   (comparative)
       │                │                │
       └────────────────┼────────────────┘
                        ▼
                proposals / learnings
                        │
                  approval gates
                        │
                        ▼
                SETTLED ARTIFACT GRAPH
                        │
             dependency/change messages
                        │
            ┌───────────┴───────────┐
            ▼                       ▼
      downstream project       upstream review
```

The worker loop remains important, but it is only one component.

---

## 55. Concise mental model

The ideas can be compressed into a few principles:

1. **Humans primarily own the backlog/specification and settlement decisions.**
2. **Agents materialize unresolved space beneath settled artifacts.**
3. **The deterministic orchestrator owns the DAG and information boundaries.**
4. **Parallel experiments are architectural search, not merely redundancy.**
5. **Monitors may stop invalid experiments but should not steer them.**
6. **Review compares actual outputs and can detect divergence no individual worker notices.**
7. **Experiments can produce code, evidence, and learnings; code reuse is optional and explicit.**
8. **Prompt/spec changes are proposals requiring settlement, not automatic manager edits.**
9. **Authority flows downward; evidence/proposals can flow upward.**
10. **Downstream artifacts are reconciled, not blindly regenerated.**
11. **Reopen only scoped settled decisions when new evidence demands it.**
12. **Lifecycle modes guide exploration intensity but do not erase settled state.**
13. **The artifact graph determines what engineering decisions humans retain vs delegate.**
14. **Provenance must survive from experimental code through reviewer learning to final settlement.**
15. **The goal is an organically accumulated, auditable body of accepted system knowledge—not a disposable sequence of agent generations.**

---

## 56. Why this seems promising based on the Traust experiments

The Traust Ralph work exposed a recurring pattern:

- early prompts can be very good yet still contain subtle loopholes,
- different capable models discover different architectures,
- some branches make impressive progress by accidentally solving the wrong product,
- others spend more compute establishing genuinely reusable boundaries,
- comparative review can discover facts that no individual worker knows,
- once a correct structure is established, many bootstrap ambiguities become less likely to recur,
- and the increasingly niche nature of failures suggests diminishing returns from endlessly hardening the inner prompt.

That points toward a broader conclusion:

> The next major leverage point is orchestration, comparison, settlement, and provenance—not making one worker prompt omniscient.

---

## 57. Candidate architecture-doc topics to split out later

This idea inventory could later be decomposed into more formal architecture documents such as:

1. **Concepts and terminology**
   - artifact
   - settled state
   - frontier
   - experiment
   - validation campaign
   - reconciliation
   - reopen proposal
   - project dependency

2. **Artifact graph model**
   - revisions
   - dependency versioning
   - stale/current/blocked/consistent states

3. **Execution DAG model**
   - node types
   - information edges
   - budgets
   - isolation
   - retries/stops

4. **Worker protocol**
   - Ralph-like execution contract
   - durable state
   - frontier scope

5. **Monitor protocol**
   - hard-gate checks
   - termination semantics
   - non-steering guarantee

6. **Review/synthesis model**
   - comparative evaluation
   - structured finding types
   - provenance requirements

7. **Settlement and approvals**
   - human gates
   - auto-approval policies
   - rule/spec mutation rules

8. **Experiment modes**
   - materialization
   - validation
   - challenge
   - scoped re-exploration

9. **Cross-project messaging**
   - upstream/downstream events
   - reconciliation readiness
   - findings/proposals

10. **Provenance and evidence**
    - branch snapshots
    - commit/hunk links
    - review annotations
    - revised learnings

11. **UI/DAG interaction model**
    - node editing
    - branch duplication
    - model selection
    - comments
    - discuss-solutions workflows

12. **Autonomy/policy model**
    - transition permissions
    - auto-approve classes
    - resource limits

13. **Lifecycle / maturity heuristics**
    - baseline discovery
    - incremental settlement
    - maintenance/reconciliation
    - scoped exploration

14. **Specification-validation workflow**
    - throwaway derivations
    - ambiguity detection
    - feedback into docs

15. **Security/trust model**
    - what agents may mutate
    - what the orchestrator enforces
    - experiment isolation
    - approval authority

---

## 58. Things deliberately left unresolved

This discussion produced the system shape, not final implementation choices. Important unresolved questions include:

- exact storage model for artifact revisions and experiment branches,
- exact DAG engine/runtime,
- whether Git commits are sufficient as the primary artifact revision primitive,
- how structured manager/reviewer outputs should be represented,
- how to implement information-isolation guarantees,
- exact budget accounting,
- whether monitors are deterministic programs, LLMs with constrained roles, or both,
- how many reviewer agents are useful for high-impact settlements,
- how to detect “large upstream change” deterministically enough to recommend exploration,
- how to score/compare architectures without overfitting to simplistic metrics,
- when reference-code access helps vs causes imitation/coupling,
- how long to retain expensive experiment repositories/logs,
- how to surface the latest consistent system state in the UI,
- how much of the DAG may be dynamically proposed by agents vs only manually/policy generated,
- how to model partial settlements within very large artifacts,
- and how to migrate/refactor the artifact graph itself when users insert/remove layers.

These are architecture questions for the system to be designed next, rather than assumptions already settled here.

---

## 59. Final framing

A useful final framing is:

> **The system is a deterministic, provenance-preserving workflow around a graph of settled engineering artifacts. Agents are used as bounded search, implementation, review, and synthesis engines. Parallel experiments explore unresolved decisions; monitors prevent invalid work from consuming unlimited compute; reviewers extract comparative learnings; humans settle changes to authoritative state; and downstream artifacts react to upstream changes through explicit reconciliation rather than regeneration.**

The system should support both:

- mature-specification → increasingly autonomous materialization,
- and draft-specification → targeted throwaway materialization → empirical validation → improved specification.

Its purpose is not to remove humans from software engineering. It is to let humans decide **where they want to own decisions**, preserve those decisions once settled, and delegate the unresolved layers to agents under an inspectable, deterministic workflow.
