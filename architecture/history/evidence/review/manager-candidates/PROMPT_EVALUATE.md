# Manager evaluation pass

Read CAMPAIGN.md to resolve the authorities, experiments, scope, output locations and decision owners. This procedure supplies review operations, not project acceptance rules or a catalogue of exceptional decisions. Derive those from the supplied authority. Do not repair experiments or select their next architecture in this pass.

## Derive the review basis

Resolve source identities, precedence and permitted mutations. Distinguish authoritative requirements and preferences from suggestions, worker interpretations and experimental additions. If the supplied sources leave precedence or a requirement unsettled, record that uncertainty instead of inventing authority.

Derive two separate records:

- **Acceptance basis:** obligations, permitted variation, consumption paths, success criteria, evidence requirements and scope limits established by the sources.
- **Salience register:** decision categories the sources treat as unusually consequential, exceptional, last-resort, restricted, approval-requiring, boundary-changing or human-review-worthy. Use the project's own category names. For each, cite the source, explain its significance, identify the triggering changes and affected assumptions, and record any required process, decision owner and limits on worker authority.

Salience is a reason for deeper review, not an automatic reason to notify a human or require approval. Derive process obligations exactly: a required investigation is not necessarily a required approval. Do not seed categories from previous campaigns or equate an unfamiliar proposal with an exceptional decision. Evidence can reveal that a change touches a source-defined category indirectly; update the register with that source and causal connection.

Inventory snapshots, available history, claims, execution evidence, recorded decisions and gaps. Read shared sources once where practical. Preserve unfinished or uncommitted work without attributing it to an earlier committed result.

## Evaluate each experiment independently

For each experiment:

1. Record its stated status, active obligation, material claims and available historical evidence. Include observable claims of success even when the worker has not declared completion. Assess initial and final validity separately when the result changed.
2. Trace the relevant consumption path from inputs through responsible actions, effects, observations, predicates and reported results. Derive which responsibilities and boundaries matter from the acceptance basis; do not assume every project has the same execution structure.
3. Check each decisive claim against its originating requirement. Examine allowed inputs, identity resolution, evidence availability and lifetime, and whether the observation proves the claimed property. Consider permitted implementation variation where the authority allows it.
4. Check false acceptance and false rejection. Inspect supplied demonstrations for shared assumptions. Where useful, construct a requirement-specific counterexample or permitted variation outside the experiment sources. Distinguish static reasoning, local probes and full execution; do not substitute one evidence strength for another.
5. Where a requirement depends on ordering, distinguish initiation, completion, publication and eventual state. Determine whether the actual observation establishes the required relation or leaves a prohibited sequence possible.
6. Separate observed failures from their diagnoses. Record prerequisites, compatibility, missing inputs, conflicting requirements, implementation defects and interruptions at the level supported by evidence.
7. Map decisions and effects to the salience register, including changes that alter a protected assumption without using the project's category name. Perform the causal review below before classifying any such item for human attention.
8. Record available usage, elapsed execution time, missing measurements, restarts, material increments and stopping reasons. Judge progress from what changed and why work stopped, not from iteration count, duration or code size alone.

Keep the worker's recorded status, the observable result and the evaluator's verdict separate. Use VALID, INVALID, WORKING, BLOCKED or UNCERTAIN with an explanation grounded in the acceptance basis. Invalidating a completion claim does not invalidate every local mechanism or observation.

## Review consequential decisions causally

For every item touching a derived high-salience category, build an evidence-backed causal record:

1. **Pressure:** What concrete requirement or failure led to the decision? Reconstruct the path from the requirement through the chosen approach to the obstacle and proposed response. Separate observed pressure from the worker's explanation.
2. **Assumptions:** Which constraints come from authority, and which come from local implementation, architecture or environment choices? Would the obstacle remain if an optional choice changed?
3. **Necessity:** What establishes that the proposed response is required? Check any source-prescribed investigation or authorization separately from technical necessity. Completing a process does not prove its conclusion.
4. **Alternatives:** Examine ordinary options and relevant sibling-demonstrated mechanisms. Compare them against the same obligations and scope. State what each actually demonstrates, what it does not, and whether an alternative removes the need or merely moves the difficulty. A successful component on an invalid experiment remains usable evidence within its proven scope.
5. **Consequences:** Check the response and alternatives against surrounding authority. Do they change an operating assumption, impose broader obligations, shift ownership, weaken the property being assessed, or conceal another failure? Do not infer necessity from a response that assumes the desired behavior.
6. **Residual value:** If the proposal is unnecessary or defective, determine whether its pressure still exposes an upstream restriction, ambiguity, hidden consequence or simplification opportunity. Recheck that possibility after accounting for ordinary alternatives; a failed proposal is neither automatically a useful finding nor automatically worthless.
7. **Disposition:** Decide what survives. Record evidence, counterevidence, confidence, gaps and the appropriate owner or next action.

Compare independent assessments only after their initial causal records exist. Use cross-experiment evidence to challenge diagnoses and refine those records before producing human-facing findings. Architectural variation alone is not a new requirement or a defect. An efficiency conclusion requires alternatives that meet the same obligations.

## Filter human attention

Suppress raw proposals that reduce to an avoidable local choice, an unsupported necessity claim, routine corrective work, a duplicate, or a question already answered by authority. Retain their disposition and useful evidence in the review trace. Direct ordinary corrections and bounded validation to the workflow allowed by the campaign; lack of validation alone does not require a human decision.

Surface an item only when causal review leaves a material finding or decision worth human attention: for example, an unresolved authoritative restriction or conflict, a consequential change requiring the designated owner's decision, or a supported broader improvement. Explain why local work cannot settle it, or why its broader consequences merit review. State the surviving concern rather than forwarding the worker's preferred mechanism. The selected finding may differ entirely from the original proposal.

Honor any explicit source-required reporting or approval obligation. Keep the required record while distinguishing it from the evaluator's recommendation; the filter does not erase mandatory disclosures. Aggregate related issues and omit discarded mechanisms from the default human-facing account unless they are necessary to explain the surviving concern.

## Outputs

Write to the campaign's selected output location:

- **REVIEW_TRACE.md:** acceptance basis and source-derived salience register; per-item causal records; rejected explanations; suppression and routing decisions; provenance and evidence gaps. This is the audit record, not the default notification payload.
- **EVALUATION.md:** independent experiment verdicts, claim identities, supported achievements and gaps, measurements, comparative findings and evidence limits. Link to the trace where deeper reasoning is needed.
- **HUMAN_FINDINGS.md:** only vetted findings and explicitly required disclosures. Each item identifies the originating authority, concrete pressure, alternatives and counterevidence that matter, surviving consequence, affected owner, requested decision or useful action, and confidence. If nothing survives the filter, say so without manufacturing a question.

Cite consequential claims. Preserve source identities, measurement methods and focused probe output. Summarize routine inspection instead of dumping logs. Do not mutate authority or promote a proposed mechanism through the wording of the report. Settlement and subsequent synthesis are separate steps.
