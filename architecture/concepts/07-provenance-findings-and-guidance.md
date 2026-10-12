# 07. Provenance, upward findings and learning

**Status:** preserving inspectable evidence is a working direction. Retention, indexing and review-depth policies remain open.

## 1. The information-loss problem

The user challenged the idea that a project should simply return "new state, findings and evidence." Who decides which observations matter upstream? A Test Harness project may know its own immediate needs but not recognize why a failed approach reveals something about Architecture.

A summary such as "the worker spent most of its time implementing a registry and failed" can erase the important fact: its choice followed a plausible interpretation of an unclear requirement. Asking the project to perfectly recognize architectural significance would make the prompt the weakest point in the evidence chain.

The stronger direction is to retain enough execution material for later independent interpretation. Project observations are hints; explicit contradictions are valuable; neither is the only surviving account.

## 2. Evidence should remain reachable

The intended provenance path can include a reconciliation request, authoritative input revisions, plan version, work package or increment, worker/session, branch/commit, artifact or hunk, test/log/observation, reviewer conclusion, human annotation and final settlement.

This does not settle a database schema. The principle is that a compact learning can point back to the concrete source that justified it. Experimental code can be unmerged and still be retained as evidence.

An immutable snapshot matters because a reference to a moving branch can no longer explain a historical conclusion. The original design preferred immutable experimental evidence; exact storage and retention policies remain open.

The evidence can include failed, interrupted and unfinished work. A final uncommitted snapshot is not the same thing as an earlier committed state. A historical first completion claim should not be attributed to a later repaired candidate.

## 3. Retain broadly; inspect proportionally

The user proposed access to all relevant commits, logs and artifacts, while leaving the human or coordinator to decide how deeply to analyze them. Deep review can be delayed until budget is available.

A progressive approach discussed in the main thread was:

- mechanical metadata, tests and identified revisions;
- semantic comparison of the candidate against relevant upstream changes;
- selected provenance inspection around suspicious or consequential areas;
- full comparative trajectory review for foundation or validation campaigns.

These are illustrative review intensities, not a finalized four-level policy. Routine acceptance still needs adequate evidence; progressive disclosure is about economical access, not skipping applicable obligations.

Hints about execution structure may help retrieval. For example, the managing agent's decisions may be more informative than every command subagent's complete transcript. Such indexing should not irreversibly discard the detailed evidence.

## 4. Deferred analysis and current relevance

An expensive campaign can remain available for later analysis. That gives preserved evidence value beyond the current run and allows the human to choose when analysis is worthwhile.

A delayed finding must be checked against newer accepted state before escalation. A loophole may already have been fixed independently. A historical concern may still explain a decision even when it no longer requires action.

Thus the system should distinguish "was true of this earlier evidence" from "requires action now." Exact obsolescence and deduplication rules remain unresolved.

## 5. Upward channels have more than one purpose

A project can report a direct blocker: requirements appear mutually impossible, a required dependency is missing, or a decision exceeds local authority. This may immediately halt a scope of work.

A reviewer can identify a nonblocking upstream issue after local success: valid interpretations diverge, an apparently harmless shortcut violates an unstated assumption, or a workaround introduces a new trust/lifecycle boundary.

Both are valuable. The original vision specifically wanted the second, not only conventional blocker escalation. At the same time, merely inconvenient implementation work should not automatically become an architectural problem.

## 6. Finding classes from the original inventory

**Contradiction:** satisfying one applicable obligation seems to violate another.

**Valid divergence:** independently plausible results interpret the same source differently. This may reflect permitted freedom, insufficient specification, or a consequential choice needing explicit settlement.

**Hidden consequence:** implementing the stated requirements introduces a material property not acknowledged upstream, such as persistence semantics, an operational requirement, ordering, a trust boundary or a latency/availability tradeoff.

**Repeated workaround:** several attempts compensate for the same awkward or missing concept.

**Boundary violation:** apparent progress bypasses a settled product or testing boundary.

**High architectural variance:** alternatives may all be valid, but the choice is important enough not to inherit accidentally from the first worker.

**Suspiciously cheap progress:** a trigger for closer inspection, not evidence of invalidity by itself.

**Repeated inability to progress:** a possible signal of a specification, environment or foundation issue, not proof of any one cause.

These categories guide attention. They are not a demand that every ordinary worker conduct a complete critique of all upstream artifacts.

## 7. The review should explain causality before escalating

The v5 manager work refined generic "report interesting findings" into an authority-derived acceptance basis and salience register. The reviewer identifies consequential categories from the campaign's actual authority, rather than hardcoding Traust-specific interfaces into the reusable prompt.

For a triggered concern it asks what concrete pressure arose, which constraints were authoritative versus locally chosen, whether the proposed mechanism was necessary, what ordinary or sibling alternatives exist, what hidden consequence remains, and who owns the decision.

This matters because a rejected proposal can still reveal a real issue. Terra-high4's separate-root proposal did not prove its own necessity. Yet reviewing the pressure behind it exposed a question about shared-engine cleanup ownership that did not disappear merely by proposing ordinary mounts.

The useful upward output is then a bounded question or bug report with evidence, not a centrally prescribed sequence of implementation steps for the parent project.

## 8. Filter human attention without deleting history

The supplied manager candidates distinguish a detailed evaluation, a causal review trace and vetted human findings. Routine mistakes, unsupported necessity claims, already answered questions and duplicates can be suppressed from the human-facing queue while their dispositions remain recorded.

A required disclosure should still be recorded even when the proposal itself is not endorsed. High salience does not automatically mean "request approval for this mechanism."

Synthesis should preserve those distinctions. A suppressed proposal must not reappear as helpful-looking worker guidance unless the underlying authority and evidence justify it.

This is an important refinement of the product idea: useful transparency does not mean forwarding every failed local choice to the human.

## 9. Learnings are themselves revisable artifacts

A reviewer might infer: "X is a good abstraction; centralize it around bar." The human can agree with the abstraction but object to bar and suggest foo. Re-analysis produces a revised learning.

The original evidence, first interpretation, human comment, revised interpretation and adoption decision should remain connected. Otherwise manager synthesis becomes another opaque authority whose mistakes are difficult to examine.

This supports comments on findings, code, commits, hunks and settlement proposals, not just final approval buttons. A learning can remain a promising reference without becoming a mandatory project rule.

## 10. Four destinations for new knowledge

A general invariant may belong in stable worker/workflow instructions after review. An architectural lesson may belong in a settled project artifact or decision. Temporary directional advice may belong in the next milestone's guidance. A one-off mistake may deserve only rejection and retained evidence.

The user emphasized the third category. Saum should be able to tilt the next attempt toward better choices without continually expanding Ralph's core rules. Foundation-specific advice can be cleared, retired, or embodied in accepted repository structure once its job is done.

A hidden fifth destination is an upstream question: some knowledge cannot honestly be turned into guidance until the appropriate project settles the ambiguity.

## 11. No silent change of authority

A downstream discovery creates a finding. A proposed clarification remains a proposal. A challenge can recommend reopening a decision. Acceptance by the relevant human/policy is what changes authority.

This distinction applies to prompt changes as well as architecture. A manager that quietly makes the next worker obey a new interpretation has effectively changed the rules without settlement, even if no architecture file changed.

Guidance therefore needs scope and status. It may say "prior evidence suggests X" or "Y remains uncertain," not quietly elevate a hypothesis into a universal contract.

## 12. Open implementation choices

How much raw execution must be captured? How should evidence be indexed across heterogeneous tools? How long should large repositories and logs remain available? What does a reference promise if content must later be retired? How are privacy, access permissions and secret-bearing traces handled? The original source names security/trust and retention as architecture topics but does not supply complete answers.

The common requirement is recoverable justification, not unlimited forever-retention at any cost. Details belong in a later architecture decision rather than being invented here.

## Sources

[S01](../history/sources/01-original-manager.md), sections 10, 13-18, 23, 28, 32, 38-39 and 44-47; [S02](../history/sources/02-project-boundary-notes.md), especially the first discussion of raw evidence and delayed review; [manager portability review](../history/evidence/review/MANAGER_PORTABILITY_REVIEW.md), [manager review](../history/evidence/review/MANAGER_REVIEW.md), and [High4 regression](../history/evidence/review/TERRA_HIGH4_REGRESSION.md).
