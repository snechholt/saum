# 04. Project interface and intake negotiation

**Status:** working problem definition and tentative contract structure. This is not a wire protocol, file schema or required framework integration.

## 1. The boundary the discussions were trying to find

The user's starting mental model was deliberately simple: a project is standalone, plausibly a Git repository, and receives a request to reconcile against another repository. It should not need Zonotron to control every internal agent role, planning document or pull request.

At the same time, a bare request/response black box loses the most valuable Scenario 01 evidence. A producer can locally recover from a bad approach without recognizing that the approach exposed a weakness in upstream intent. It can claim completion using its own inadequate tests.

The interface therefore needs enough structure for accountable objectives, viable intake, inspectable results and independent acceptance, without importing the whole local methodology into Zonotron.

The distinction is **required participation semantics, optional internal method**. Whether these are carried by prompt instructions, repository files, a driver, messages or APIs remains open.

## 2. Three conceptual contracts

The validation branch proposes separating three contracts rather than one universal integration.

### Reconciliation contract: what must become true

This identifies the current accepted local state, relevant authoritative upstream revisions, the objective and scope, constraints to preserve, and who can change those obligations.

A request might be:

> Reconcile the existing Test Harness against A13, establishing its scenario-testing obligations while preserving the accepted public runner and replaceable-SUT boundary.

The request should not silently prescribe every mechanism the project must use. It can be a large objective whose next increments are negotiated.

### Artifact contract: what can be inspected

This identifies the candidate revision, artifact location and format, claimed coverage, local status, evidence and execution provenance, plus optional observations.

Methodological freedom does not mean arbitrary outputs can be automatically accepted. An unfamiliar format with no interpreter, relevant evidence or available reviewer can block meaningful validation. The appropriate response might be tooling, alternative evidence or human judgment, not pretending Zonotron has a universal parser.

### Validation contract: what conformance means

Architecture -> Behaviour Model asks about faithful behavioral interpretation. Architecture -> Test Harness asks about adequate tests and legitimate external observation. Design -> Implementation asks about structural and interaction conformance. Test Harness -> Implementation asks whether the implementation passes required independently defined tests.

Those are different relations. The branch therefore proposes validation contracts associated with dependency edges. Their author, storage, versioning and settlement rules are unresolved. A contract may need to be intelligible to both producer and independent validator without allowing the producer to weaken its own obligations.

## 3. Inputs should identify a basis, not just say "sync"

A project needs to know what it is reconciling from and toward. Relevant inputs include old and new upstream state, current accepted local state, the requested scope, known dependencies and constraints, and review or resource policy where applicable.

Exact revisions matter to provenance. "Use the latest architecture" is not enough to explain later why a candidate made a particular decision. The handling of a newer upstream revision arriving during work remains open, but the original request's basis should remain knowable.

The input can also include temporary guidance and reference evidence. Those should be distinguishable from authority: a promising earlier mechanism is not a requirement merely because it is attached.

## 4. The conversation before execution

The user identified negotiation as a high-leverage architectural question. The following questions collect the ones developed across the main thread and the validation branch. They are candidate questions, not a compulsory long questionnaire for every change.

### Understanding and authority

What obligations do you understand this request to introduce? Which existing decisions constrain your work? Does anything conflict, remain ambiguous, or require an upstream decision before a credible plan can be made? Which apparent restrictions come from authority, and which come from your chosen local approach?

### Consumable scope

Can you consume this request coherently as one unit? If not, what units suit your local state and methodology? Which units can be combined without obscuring acceptance? Which should remain separate because their design or evidence is uncertain? What is the smallest useful next increment?

The point is not a universal token count or mandatory PR size. It is a project-specific description of what can be executed and reviewed coherently.

### Dependencies, atomicity and parallel work

What must precede what? Which units can proceed independently? Which must be evaluated or settled together? Which require another project or unresolved interface first? Does concurrent work risk changing a shared assumption or producing candidates that cannot be combined?

These answers can influence a global plan without exposing every internal task.

### Uncertainty and exploration

Which parts are ordinary work under the existing structure? Where are there materially different approaches? Which uncertainty might affect other projects or reveal an upstream issue? Would discussion, a focused probe, or multiple candidates provide useful discrimination? Can exploration be limited to one work package rather than the entire reconciliation?

A project should be allowed to say that a particular experiment would add little value. Zonotron should also be able to request validation-focused exploration for reasons outside the project's immediate delivery incentives.

### Evidence and acceptance

What would demonstrate each increment's claims? Which obligations are directly observable, and which currently lack convincing evidence? What must be checked to preserve previous accepted progress? What integrated or cross-cutting constraints cannot be established by accepting increments individually?

Project-proposed evidence is input to validation, not the sole acceptance standard.

### Remaining scope and replanning

How much of the objective can sensibly be planned now? What remains explicitly unplanned? What would trigger replanning? What new information from the first increment is likely to change the rest? Does the proposed plan collectively account for the whole obligation even when its horizon is short?

### Cost and human decisions

What resources or human decisions materially affect the next work? Where would a small clarification be cheaper than further experimentation? Which decisions are consequential enough that the project should pause rather than silently choose? How do local priorities interact with the global objective?

Precise forecasts were not required by the discussions. Explanations of uncertainty may be more useful than unsupported numeric estimates.

## 5. Who proposes, who challenges, who accepts

A working pattern is project proposal, Zonotron global review, project pushback or revision, then acceptance under the relevant policy. Multiple rounds are allowed when they add value; one-shot planning is not required.

Zonotron should not claim complete local domain knowledge. The project should not claim that its preferred methodology determines global authority. The combined conversation establishes an executable contract rather than central micromanagement.

A candidate exit condition is mutual viability: the project can execute it, and the coordinator can account for obligations, dependencies and acceptance. How to stop fruitless negotiation, resolve disagreement or escalate to a human is still open.

## 6. What comes back from execution

The strongest current boundary proposal asks for candidate state, local result/status, evidence references, execution provenance and claims. Explicit blockers and optional observations are useful. It does not rely on every project correctly recognizing all architecture-level implications.

An illustrative response can distinguish:

```text
Candidate: identified artifact revision
Requested basis: identified upstream revisions and objective
Claim: increment's obligations now satisfied
Local status: ready for review
Evidence: tests, outputs, observations and their limits
Provenance: sessions, commits, branches and retained artifacts
Blockers/observations: optional local explanations
```

This is illustrative content, not a finalized schema. The actual status vocabulary is unsettled. The review evidence does show why "working," "blocked," "ready," and "publicly reports success" cannot be conflated.

Preserved and changed decisions may be useful claims, but requiring the producer to supply a perfect semantic change inventory would recreate the self-diagnosis problem. Independent evaluation remains necessary where those claims matter.

## 7. Blockers and subtle findings need different paths

A direct contradiction can be raised immediately: "I cannot satisfy X and Y together." The project may need a clarification before work can continue.

A subtler finding might arise only during comparative review: one failed approach exposes a loophole, different valid approaches reveal discretionary space, or a workaround has a hidden consequence. The result may be useful even when local execution is complete.

Both should be able to reach upstream review. Blocking status is not the sole criterion for upstream relevance. Equally, every local difficulty should not become a human ticket. The causal attention filter is described in [provenance and findings](07-provenance-findings-and-guidance.md).

## 8. Analysis-only and multiple-candidate requests

A project may first receive an impact-analysis request rather than execution authority. It can return a proposed breakdown, dependencies and uncertainty without changing accepted state.

The user's experiment proposal goes further: Zonotron could send the same ordinary request to isolated instances of a project. Each instance uses its normal internal method. The coordinator compares candidates afterward. This avoids requiring every methodology to implement the same experiment machinery.

The scope can be one work package, a complete project reconciliation, or, speculatively, a downstream subtree. These are not all equally appropriate for a first implementation.

## 9. Loose coupling versus adapters

Earlier responses casually listed Kiro, Spec Kit and AO adapters. The user correctly challenged the implication of tight integration. The desired boundary may be much looser: ordinary instructions in, artifact references and evidence out, with a small generic driver responsible for the handoff.

A rich adapter could improve discoverability, capture evidence or expose local capabilities, but it should be an option whose value is demonstrated. A supported methodology need not expose all of its internal workflow concepts to Zonotron.

Similarly, an experiment helper can be a reusable procedure plus infrastructure, not a mandatory internal ritual. Projects can also run their own experiments that Zonotron only sees through candidate evidence.

## 10. What the interface should avoid becoming

It should not require every worker to conduct a broad upstream architecture critique. It should not force a global list of software-specific exceptional cases onto school plans or other domains. It should not require the whole local task graph to be approved centrally.

It also cannot be so weak that an opaque success flag counts as conformance. The open design problem is the smallest contract that preserves meaningful control, evidence and interoperability.

## Sources

[S02](../history/sources/02-project-boundary-notes.md), all major concerns; [S03](../history/sources/03-validation-scoped-settlement-and-adaptive-planning.md), sections 1-2, 4-6, 10-13; [S01](../history/sources/01-original-manager.md), sections 23, 26, 38 and 46. Main-thread AFK negotiation questions are M07 in [source map](../history/sources-and-coverage.md).
