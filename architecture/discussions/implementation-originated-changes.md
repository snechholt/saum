# Implementation-Originated Changes and Upward Reconciliation

**Status:** Exploratory design direction. Several mechanisms and policy choices remain unresolved.

## 1. Motivation

Saum's authority graph is directionally hierarchical: Architecture governs Design and Test Harness, which in turn govern Implementation (in the illustrative diamond). Changes originating upstream are reconciled downward.

However, developers will sometimes work directly in Implementation without first initiating an architectural or design change. This may be exploratory prototyping, a targeted fix, or a feature that eventually becomes complete.

Saum must support this workflow without requiring the developer to retroactively initiate the entire top-down process manually, or treating directly modified implementation as an inherently invalid state.

This is primarily an **ergonomics and flexibility requirement**, not an attempt to replace top-down development or optimize for minimum development cost. Top-down reconciliation remains the preferred model for planned architectural changes.

The broader principle is:

**Changes may originate anywhere, but authoritative decisions retain their established direction and precedence.**

## 2. Separate change initiation, assessment, and authority

A change originating in a downstream project may provide evidence or motivate reconsideration of upstream authority. It does not thereby gain authority over upstream projects.

In particular, upward reconciliation must be fundamentally more conservative than downward reconciliation:

- **Downward:** An accepted upstream decision establishes obligations for downstream projects. They must reconcile with that authority or surface a genuine conflict.
- **Upward:** A downstream change is a candidate. Upstream projects evaluate whether it should be accepted, modified, or rejected. Their existing decisions remain authoritative by default.

Saum must not treat upstream reconciliation as finding whatever architectural adjustments are necessary to legitimize code that already exists.

The burden is on the proposed downstream change to justify deviations from accepted authority.

That justification may succeed, but acceptance must be a genuine decision rather than an automatic accommodation.

This asymmetry is fundamental to preserving the value of Saum's architecture-first model.

## 3. Distinguish exploratory state from proposed final state

A major challenge arises if Saum monitors ordinary development activity, such as commits on a feature branch.

A temporary shortcut during prototyping is not equivalent to an intentional proposal to change the final design.

For example, Implementation might bypass a mandated factory pattern because direct construction makes an experiment easier. That does not necessarily mean the developer believes the factory requirement should be removed.

Saum needs a way to distinguish:

- **Exploration:** Temporary changes, experiments, incomplete work, and intentionally deferred conformity.
- **Proposed implementation:** A candidate intended to become the accepted state, subject to reconciliation.
- **Accepted or released implementation:** A further level of commitment, with potentially stricter requirements.

These distinctions may be supplied explicitly by the human, inferred from workflow context, or established through project events. Inferences should remain correctable.

Exploratory changes should not automatically launch expensive upstream reconciliation.

Saum may still detect conflicts and inform the developer of relevant existing constraints. The appropriate response may be a lightweight observation rather than a request to resolve the conflict immediately.

The developer should retain freedom to experiment without making Saum intrusive or encouraging users to mute it.

Early implementations may support only exploration and proposed changes, with richer distinctions between acceptance, merge, and deployment added later. The overall conceptual model should nevertheless accommodate them.

## 4. Proactive monitoring and guidance

Saum could monitor implementation activity and recognize emerging divergence from governing specifications.

For instance, when code bypasses an established factory pattern, Saum could identify the relevant design requirement and surface it while the developer is still working.

This offers some benefits of top-down development even when the work originates locally.

Monitoring should be sensitive to intent, development phase, significance, and risk.

Ordinary or temporary deviations should produce little or no interruption. Significant violations of architectural invariants may justify more prominent guidance.

Continuous monitoring is not necessarily required initially. Commit-based, checkpoint-based, or explicitly invoked assessment may provide a useful starting point.

**Monitoring, warning, requesting reconciliation, and blocking acceptance are distinct operations.** Their triggers and consequences should be governed by policy.

## 5. Configurable strictness, with upstream authority as the default

Saum should support configurable behavior along a spectrum, conceptually ranging from:

- **Adaptation-oriented:** Give considerable weight to the developer's intended implementation; actively evaluate whether upstream specifications should accommodate it.
- **Authority-oriented:** Strongly preserve accepted upstream decisions; require substantial justification before reconsidering them, and reject incompatible downstream changes by default.

This could be represented by an adjustable policy or strictness setting, rather than necessarily a literal numerical slider.

The setting should influence the degree of scrutiny, the burden of justification, escalation thresholds, and the extent of autonomous proposal generation.

It must not make existing authority disappear. Even in adaptation-oriented operation, upstream revisions require proper approval.

Conversely, strict operation should not prohibit legitimate architectural evolution.

A useful distinction may exist between the developer's preferred workflow and non-negotiable project constraints. Critical architectural, security, or integrity requirements may mandate review regardless of the selected strictness.

**The unresolved design challenge is preserving genuine upstream authority while allowing deliberate, justified bottom-up evolution.**

## 6. Upward assessment before downward reconciliation

Two possible strategies were considered:

1. Jump directly from the changed Implementation to Architecture.
2. Propagate through immediate graph neighbors, potentially through multiple intermediate projects.

Immediately involving Architecture ensures the anchor gets visibility, but risks bypassing intermediate responsibilities and requiring upstream projects to understand distant implementation details.

Reconciling each intermediate project immediately during upward propagation has a different problem: a later architectural decision may invalidate the already completed intermediate reconciliations.

The emerging direction is to separate **upward assessment and proposed decisions** from **accepted revisions and downward reconciliation**.

A candidate sequence is:

1. Implementation changes and its intent are identified.
2. Its immediate governing projects assess the implications.
3. Relevant findings and proposed authority changes propagate upward.
4. Higher-level projects review decisions within their authority.
5. Once the necessary upstream decisions are settled, reconciliation proceeds downward through affected projects.
6. Existing implementation is retained wherever it already conforms; otherwise it must be revised.

Intermediate assessments need not immediately modify their accepted specifications.

This avoids repeated premature settlements while preserving the ability to discover architectural implications through intermediate projects.

## 7. Propagation and stopping rules

The exact upward propagation rule remains unresolved.

An initial idea was to stop when an intermediate project determines that a change is entirely within its own authority.

However, this creates a trust problem: **can Design independently determine that Architecture has no interest in a proposed Design revision?**

A stronger candidate rule is:

**When a project proposes changing its own authoritative specification in response to a downstream change, that proposal must be exposed to its governing parent projects before settlement.**

The parent may determine that the change has no implications for its own authority, allowing propagation to stop. Otherwise the issue continues upward.

This gives a tentative pattern:

- No authoritative revision needed at the current project: propagation may stop, provided no higher-level obligation has been implicated.
- Authoritative revision proposed: propagate to the governing parent for assessment.
- Parent identifies further implications: continue upward.
- Parent confirms the proposal is within delegated discretion: stop escalation and proceed toward settlement.

An important unresolved question is how to determine that no higher authority is implicated when no local revision is proposed. Simply allowing a downstream agent to declare a change harmless could miss architectural consequences.

Possible mechanisms include explicit delegation boundaries, parent-defined review policies, or lightweight independent checks.

The important property is that the responsibility for deciding what a parent cares about should ultimately derive from the parent's authority, rather than solely from its child.

## 8. Preserve graph locality without restricting evidence

Upward reconciliation should preferably operate through immediate dependency relationships, preserving each project's limited responsibilities.

However, a downstream change may be relevant to an ancestor several levels away.

Saum can maintain a shared change record containing the original intent, candidate implementation, evidence, identified conflicts, affected scope, proposals, and decisions.

Projects can consult relevant evidence from that record without acquiring normative dependencies on distant descendants.

Thus a Design project reconciling after an architectural decision may consider:

- Its accepted governing architectural requirements.
- Existing Design decisions and project state.
- The downstream candidate implementation that motivated reconsideration.
- Relevant discoveries and validation evidence.

The candidate implementation helps inform what Design should preserve where appropriate. It does not become another authoritative parent.

This distinction is especially important for Test Harness: tests should reflect independently justified requirements, not merely encode the behavior of the current implementation.

The approach should scale to deeper dependency chains and diamonds, with relevant assessments potentially occurring in parallel and findings consolidated at shared ancestors.

It cannot guarantee that a single upward and downward pass will resolve everything. Genuine conflicts or newly discovered constraints may require further negotiation.

## 9. Independent critical review

Retrospective reviews have an inherent risk of rationalizing decisions already embodied in working code.

Saum should counter this through independent assessment, not merely by asking the implementing agent to explain or defend its approach.

One proposed mechanism is:

1. Capture the original problem and intended outcomes independently of the implementation.
2. Start an isolated review context with that intent and the relevant accepted authority, but without the new implementation's solution.
3. Independently examine how the desired outcome should be achieved, including alternatives, risks, architectural implications, and quality requirements.
4. Compare that assessment against the actual implementation and its evidence.
5. Require explicit justification for consequential differences.

The independent assessment is a critical reference, not automatically a replacement authority.

Reviews must have a genuine ability to reject the implemented approach and require changes.

Retrospective ADRs and other documentation should accurately describe decisions made during development versus conclusions reached during subsequent review, rather than inventing a fictional prospective decision process.

## 10. Scope of reconciliation

Not every implementation change should produce an architectural revision.

Routine fixes, refactoring, dependency maintenance, and changes fully covered by existing authority may require only local validation.

The relevant question is whether the change creates or alters obligations that other projects are responsible for governing.

Review obligations should be proportional to semantic impact and risk, while preserving mandatory checks where required.

For consequential changes, review should consider relevant architectural invariants, security, reliability, lifecycle properties, data integrity, interfaces, quality attributes, alternatives, and implications for affected projects.

The goal is comparable rigor for consequential decisions regardless of where development began—not a claim that retrospective review can discover every possible issue.

## 11. Change requests as another initiation mechanism

This discussion also suggests separating **requests for change** from **observed candidate changes**.

Saum might receive an intent directly, such as a request to implement a feature, without any project having changed yet.

Alternatively, Saum might detect an implementation change that has already occurred.

Both can enter a common assessment and reconciliation process, although the first requires determining where work should begin, while the second must evaluate existing candidate work.

Saum coordinates the process but does not independently become the authority deciding what the architecture should be.

This potentially generalizes to changes originating from incidents, test discoveries, external requirements, dependencies, and other sources.

## 12. Open design questions

The main questions to resolve are:

1. **Intent signaling:** How should Saum distinguish temporary experimentation from changes intended for acceptance, without excessive manual annotation?
2. **Strictness policy:** How should adaptation-oriented versus authority-oriented behavior be configured, and which constraints always override it?
3. **Upward review:** Must every proposed authoritative revision be offered to its governing parents? How can parents predefine areas of delegated discretion?
4. **Stopping conditions:** Who can determine that higher authority is unaffected, especially when the intermediate project proposes no change?
5. **Traversal:** How should upward assessment work across deeper DAGs, converging paths, and multiple affected branches without premature settlement?
6. **Settlement:** Which guarantees are required at local completion, project acceptance, merge, and release?
7. **Review independence:** How can isolated planning and architectural challenge remain independent of the candidate solution while still receiving enough context?
8. **Monitoring ergonomics:** How can Saum provide useful guidance without constant interruptions or reconciliation churn?

## 13. Suggested evaluation

A future experiment could compare architecture-first and implementation-first workflows for representative changes, including at least one with substantial architectural implications.

The main interest is developer ergonomics, interruption, quality of architectural review, rework, and successful settlement—not simply which workflow is cheaper or faster.

Implementation-first development is not intended to displace the architectural process. It may nevertheless be particularly appropriate for exploration and some classes of routine change.

## Summary

The emerging model preserves **directional authority but allows change initiation anywhere**.

Downstream modifications become candidate evidence rather than new authority. They may trigger upward assessment and proposals, with upstream projects remaining conservative and retaining the ability to reject them.

Prototyping and proposed final implementations must be distinguished so that ordinary exploratory work does not continuously trigger formal reconciliation.

Upward assessment should avoid premature intermediate settlements. Proposed authority changes should receive appropriate parent review, while shared evidence allows projects to understand relevant downstream work without creating new normative graph edges.

Once governing decisions are accepted, reconciliation runs downward, preserving compatible existing work where possible.

The central design tension is to make direct implementation work ergonomically viable **without weakening the authority of the upper layers or turning retrospective reconciliation into automatic architectural accommodation**.