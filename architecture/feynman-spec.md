# Saum: the Feynman spec

Reviewed implementation: `1914e96`. This describes the current product contract,
checks it against the code, and separates confirmed bugs from decisions still needed.
It does not replace the original architecture dossier.

## Explain it without the machinery

Saum is a workshop for building software with AI workers.

- A **factory** is one workshop, with a Git repository and a place on disk.
- A **user story / goal** says what someone should be able to do.
- A **plan** is an ordered list of work items. Each item has acceptance criteria:
  observable statements that explain what would count as doing that item properly.
- The **goal agent** does the work, one item at a time.
- The **factory agent** inspects that work and explains the outcome to the operator.
- The **human operator** decides whether the result is acceptable.
- **Saum**, not either AI, controls the queue and permission to advance.

Think of the goal agent as the builder and the factory agent as the foreman.
The foreman is also AI: its explanation is useful, but it is not a truth oracle.
The human signs off. An AI cannot approve its own item through the exposed MCP functions.

## Follow one goal through the workshop

Suppose the goal is “a user can log in.”

1. A planning agent proposes work items and acceptance criteria.
2. The operator reads the plan. Selected-text comments can be batched into feedback.
   Feedback creates a revised plan; the earlier version stays in history.
3. The operator approves the plan. Saum queues its first work item.
4. The goal agent receives **only that item**, plus the goal and relevant feedback.
5. It implements the item and submits a separate receipt for each criterion,
   in order, using `saum.work_item.submit_acceptance_criteria(...)`.
6. Once every criterion has a receipt, it calls
   `saum.work_item.request_review()` and stops its turn.
7. Only after that turn finishes successfully does Saum queue the factory agent.
8. The factory agent inspects the goal worktree and transcript, then submits an
   executive summary, verification assessment, and risks.
9. After the factory-agent turn finishes, the operator can:
   - **Send Feedback:** resume the goal agent on the **same item**.
   - **Approve:** give the **next item** to the same goal-agent thread.
10. After the final executive summary is approved, every item is approved.
    The current UI calls the goal “Completed.” That does **not** currently mean
    its code has been integrated into `main`, deployed, or independently certified.

The queue must never interpret a review request, an AI summary, or an exit code
as the operator's approval.

## The smallest useful state machine

```text
plan → human plan approval
                  ↓
          item queued → working
                          ├─ blocker → operator reply / explicit resume
                          ├─ failure → explicit retry
                          └─ criterion receipts → review requested
                                                    ↓ successful turn exit
                                             factory-agent summary
                                                    ↓ successful turn exit
                                             waiting for operator
                                               ├─ feedback → SAME item
                                               └─ approval → NEXT item
                                                               ↓ last item
                                                       all items approved
```

There are three different things called “done”:

| Thing | What its completion currently means |
|---|---|
| Agent turn | Codex finished successfully and performed its phase's required MCP submission. |
| Work item | A human approved its latest completed executive summary. |
| Goal | Every work item has that human approval. |
| Delivery into `main` | Not implemented as a workflow step. |
| Deployment | Not implemented as a workflow step. |

An accepted receipt is not automatically a passed test. An approved summary is
not automatically an approval of a particular immutable code version. Those
are important weaknesses below, not details to hide behind the word “done.”

## What owns what

New factories use this layout:

```text
<configured factories directory>/<factory-id>/
├── main/                   main Git checkout; factory-agent working directory
└── <worktree-id>/           goal checkout; goal-agent working directory
```

Planning revisions and later implementation reuse the goal's chosen worktree.
Worktree names currently come from session IDs, not necessarily story IDs.
Existing factories retain their older recorded paths; changing the default
factories directory does not relocate them.

Two IDs must not be confused:

- A **Saum session ID** identifies one queued/executed turn and its recorded output.
- A **Codex thread ID** identifies the conversation resumed across turns.

One approved plan has one goal-thread registry. One factory has one factory-agent
thread registry. Resuming a thread is not the same thing as keeping a process
running forever. Saum starts a new CLI invocation for the next turn.

SQLite stores the plans, receipts, reviews, decisions, queue, and exposed agent
output. Code lives in Git worktrees. Native Codex conversation files also live
in Codex's own storage, outside that SQLite database.

## Rules we should be able to explain and test

1. **Approval controls scheduling.** No implementation turn before plan approval;
   no next-item turn before approval of the current item's executive summary.
2. **Feedback does not advance the queue.** It sends the current item back for changes.
3. **Roles do not grant each other authority.** Goal agents submit criterion receipts
   and request review. Factory agents submit executive summaries. Neither MCP role
   can approve on behalf of the operator.
4. **Old decisions remain inspectable.** Revisions and feedback do not rewrite the
   old plan or executive review.
5. **Receipts can be corrected without erasing history.** Discovering a mistake in
   criterion 1 while working on criterion 2 must have a supported recovery path.
6. **A review identifies the artifact it reviewed.** Approval must not silently
   accept different source bytes from those summarized.
7. **A workspace is actually inside its factory.** A path that merely looks like
   a child directory is insufficient if a symlink redirects it elsewhere.
8. **Thread continuity is checked, not presumed.** Missing or changed thread identity
   cannot silently count as continuing the same conversation.
9. **Failures do not buy unlimited retries.** Interrupted and failed turns require
   an explicit retry. Successful review handoffs are normal workflow, not retries.
10. **Completion says what was delivered.** “All items approved,” “integrated into
    main,” and “deployed” must not be presented as interchangeable claims.

Rules 1–4 are substantially enforced by database transactions, phase-scoped
MCP functions, session ownership checks, and the latest-review gate. The review
found counterexamples to rules 5–8, and missing product decisions around rule 10.

## Confirmed bugs

### B1 · High: approval is not tied to the reviewed source version

**Child-simple version:** the operator can approve a photo of one chair while
Saum accepts a different chair that was swapped in afterward.

`api/work_item_review.go:36` checks the review ID, decision, item state, and summary
turn state. Neither the review nor approval records/checks a commit or content
snapshot of the implementation.

**Reproduced:** create an executive summary, change a source file afterward, then
approve that summary. Approval succeeds and queues the next item.

**Required direction:** freeze or identify the reviewed artifact, bind receipts,
summary, and approval to it, and reject or re-review stale artifacts. A commit ID
alone is insufficient while uncommitted source changes are still allowed.

Probe: `TestFeynmanApprovalMustReferToReviewedBytes`.

### B2 · High: workspace containment is only a textual path check

**Child-simple version:** the workshop checks the address painted on a door,
but does not check where the door actually leads.

`api/revision.go:27` uses `filepath.Rel`, then `os.Stat` at line 31 follows symlinks.
It does not establish that the resolved directory belongs to the factory's Git
worktree set. Repository creation in `api/planner.go` also follows filesystem paths.

**Reproduced:** move an owned worktree aside, replace its path with a symlink to
an outside directory, and prepare an implementation session. Saum accepts the
path even though it resolves outside the factory.

This can grant the next Codex invocation a workspace outside the intended factory.
The reproduction assumes an actor can replace that directory; it is not a claim
that a remote browser alone can perform the replacement.

**Required direction:** check canonical root/repository/worktree ownership,
reject redirected or unregistered worktrees, and account for replacement between
checking a path and starting a process. Git worktrees are not a security boundary
by themselves.

Probe: `TestFeynmanWorkspaceMustResolveInsideFactory`.

### B3 · Medium: a criterion receipt cannot be superseded during its turn

**Child-simple version:** the builder writes down a measurement, discovers it
was wrong, and is forbidden from recording the corrected measurement.

`api/work_item_tools.go:71` rejects a repeated criterion index when its evidence
or summary differs. There is no correction/invalidation function in the exposed
work-item API.

**Reproduced:** submit criterion 0, then submit updated verification for criterion 0.
The server rejects it as already submitted with different evidence.

The agent can preserve a stale receipt and request review, or stop and require a
new turn. Neither is a good normal path for iterative implementation.

**Required direction:** append superseding evidence/invalidation events rather
than overwriting history, and require the latest applicable evidence at review.

Probe: `TestFeynmanCriterionEvidenceMustBeCorrectable`.

### B4 · Medium: native thread continuity does not fail closed

**Child-simple version:** Saum asks to continue a conversation but accepts an
answer from a different conversation without checking the name tag.

`api/planner.go:317` accepts `thread.started` output.
`api/work_item_sessions.go:42` writes that identity into the registry without
comparing it with the thread requested for resumption. Turn completion does not
require a nonempty observed identity.

**Reproduced with a controlled Codex-protocol helper:**

- A resumed invocation omits `thread.started`: the turn still completes.
- It reports a different ID: the turn completes and the registry is overwritten.

Normal native resumption was previously tested successfully. This finding concerns
failure handling when the external protocol stops honoring that assumption,
not evidence that normal Codex routinely changes thread IDs.

**Required direction:** record the requested identity, validate the observed
identity, and require an explicit recovery action before replacing a thread.

Probe: `TestFeynmanNativeThreadIdentityMustFailClosed` (two cases).

## Confirmed design gaps, not an approval-gate bypass

### G1 · Verification receipts are narrative, not independently established facts

`submitAcceptanceCriterion` validates index order, nonempty text, sizes, and
ownership. Its `evidence` entries are strings. `requestWorkItemReview` counts
receipts; it does not establish that their claimed checks occurred or passed.

**Reproduced:** an empty implementation with no executed checks reaches review
using receipts containing “I say it passed.” The factory agent and human still
have to review it; this does not skip human approval.

**Decision needed:** what evidence must Saum itself establish? Prefer references
to recorded command exits, test runs, diffs, screenshots, or explicitly attributed
manual observations. Keep “reported,” “observed,” and “human accepted” distinct.
Do not promise that a second AI can prove arbitrary acceptance criteria true.

Observation: `TestFeynmanObservationNarrativeIsNotProof`.

### G2 · Accepted work does not become the next goal's baseline

`decideWorkItem` records approval and schedules the next item. It does not integrate
the goal branch into `main`. `api/planner.go:164` creates new worktrees from the
main repository's `HEAD`.

**Reproduced:** commit a file in a goal worktree, approve its final summary, observe
that the goal is “Completed,” then create another goal worktree. The approved file
is missing from that new worktree because `main` never received it.

**Decision needed:** define the integration gate. For example, final item approval
could lead to a separate integration review, regression checks against current
main, conflict handling, and approval of an exact resulting commit. Until a policy
is chosen, call the current state “All items approved,” not “delivered into main.”
Do not silently add automatic merges as a side effect of this review.

Observation: `TestFeynmanObservationAcceptedWorkIsNotIntegrated`.

## Other boundaries to make explicit

- **Local-only is not human authentication.** `api/main.go:145` rejects unexpected
  hosts/origins, but a trusted local process can make operator HTTP requests.
  The current model trusts the operator's host. MCP role restrictions and Codex
  sandboxing should not be described as protection against every local process.
- **Planning scope is partly an instruction.** Planning uses a writable worktree
  and allows planning artifacts. Saum hard-gates the implementation queue; it does
  not mechanically distinguish a planning file from application source written
  by the planner. Use a read-only phase if that stronger guarantee is required.
- **A SQLite backup is not a complete continuation backup.** Restoring the database
  without native Codex thread files can leave a registry pointing to a missing
  conversation. The current retry path reuses that ID; there is no operator-facing
  thread-recovery flow. Rehydrating a new thread must be explicit and auditable.
- **Execution is globally serial today.** Separate factory/goal conversations are
  logical identities, not a promise of simultaneous workers. One long invocation
  can delay unrelated queued work and summaries.
- **Budget and retention are incomplete policies.** Per-turn timeouts and explicit
  retries exist, but there is no monetary budget or transcript-retention quota.

## Teach-back questions

- **Can the goal agent unlock the next item by saying “done”?** No; it needs a
  factory-agent summary and a human decision recorded against the latest review.
- **Does submitting all criteria prove the item works?** No; it currently proves
  that the agent submitted a receipt for every criterion.
- **Does approving a summary pin exact code bytes?** Not yet; B1.
- **Does “inside this folder” guarantee filesystem containment?** Not yet; B2.
- **Can the agent revise an earlier criterion receipt?** Not in the same turn; B3.
- **Does Saum always verify that a resumed thread is the same one?** No; B4.
- **Will another goal automatically see the completed goal's code?** No; G2.

## Reproduce without touching live factories or spending model tokens

The opt-in probes use isolated temporary stores, real Git worktrees, and a local
MCP helper process. They do not call the real Codex service or mutate the live DB.

```sh
cd api
go test -tags feynman_review -run '^TestFeynman' -v -count=1 ./...
```

At the reviewed baseline, expect four failing contract test functions (five failed
cases because thread identity has two variants), plus two passing observations
that demonstrate the assurance/delivery gaps. Passing observations do not mean
those gaps are good behavior. Update these opt-in probes as each finding is closed.

The normal suite excludes the `feynman_review` build tag. Its passing result must
not be used to claim that the intentionally exposed contract failures are fixed.

**Fix order:** canonical workspace ownership; artifact-bound review/approval;
correctable, versioned evidence; fail-closed thread identity with explicit recovery;
then agree on integration and independent-verification policies.
