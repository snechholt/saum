# Traust V5 archived campaign

## Purpose

Evaluate independent attempts to materialize the externally selected Scenario 01 foundation. Determine actual product validity and reusable local evidence, then surface supported workflow, architecture, environment or upstream questions. This is a completed campaign; no evaluator may modify or steer its workers.

## Immutable inputs

Resolve paths relative to this file:

- Experiment archive: `../../v5_20261003_234931`.
- Shared product documentation: `../../documentation@11c2a1/docs`, source revision `11c2a1efefcd5bf91d6bed3e6eff070e001211a6`.
- Frozen worker baseline: each trajectory's `ralph/{PROMPT,RULES}.md` and `test-harness/{PROJECT,PROJECT_PREFERENCES,ARCHITECTURE,MILESTONES}.md`.
- Referenced stable project guidance: each trajectory's `test-harness/specs/DEPLOYER_TEST_INTERFACE.md`. Resolve other authoritative references reached from the product and worker sources as needed.
- Changes that produced the shared baseline: `../../v1-v5-patches`.
- Experiments: `luna-high1`, `luna-high2`, `sol-high1`, `sol-high2`, `terra-high1`, `terra-high2`, `terra-high3`, `terra-high4`.

Shared requirements may be read once. Verify frozen baseline agreement before comparing trajectories. Each directory contains a final snapshot, ordered `patches/` with commit identities, and `.ralph/logs/` with driver/worker output. There is no requirement to treat this archive as a live Git worktree. Unfinished logs and final files not attributable to a patch must be identified explicitly.

Product behavior comes from the shared product documentation; worker process and project acceptance come from the frozen worker/project guidance. Experimental code, mutable decision inventories and worker justifications are evidence to assess against those sources. Their existence does not establish that an added obligation is necessary or settled. Derive consequential decision categories and their review obligations from the authoritative texts, including materially relevant surrounding requirements.

## Work scope and acceptance

The active milestone is Scenario 01 Happy Path Deployment. Workers may claim local readiness; only external evaluation can accept or advance it. All meaningful scenario requirements and materially relevant referenced product boundaries remain authoritative.

Acceptance requires:

- The public suite path launches or reaches an independently supplied SUT and reports correct success/failure with a failing exit status on failure or incompleteness.
- Inputs cross through documented product/runtime configuration, ordinary external dependencies, or explicitly established test-only contracts.
- Harness code acts as tester; the SUT performs deployer-owned behavior.
- Identifiers, routing conventions, client tools and execution strategies permitted by the specification remain independently replaceable.
- Required evidence survives the documented resource lifecycle and comes from the same boundaries available with an independent SUT.
- Meaningful assertions distinguish correct behavior from representative prohibited behavior, including required ordering.
- Known-good and known-bad exercises use that public path; a good actor respects the active documented boundaries.
- Runtime files are contained under the reported per-run root and Docker resources are isolated to the run.
- Test-only contracts follow the deferred challenge process and retain reviewable necessity/proof evidence.

## Historical and final claims

Inspect the first covered/ready claim and any public whole-scenario success, including on partial workers. Judge initial acceptance from available historical evidence, and assess the final snapshot separately if it contains a repair or deterioration. If a historical result cannot be reconstructed adequately, mark that historical judgment uncertain rather than silently using final code.

No rejected initial claim earns valid later depth automatically. Invalid branches still supply local techniques, friction and specification interpretations for qualified comparative analysis.

## Outputs

Write evaluation, causal review trace, vetted human findings and later synthesis to an external output directory selected for that manager run. Do not overwrite this campaign, experiment sources, prior reports or settled authority. Record source paths/commits, evidence type, confidence, measurements and gaps. Rule, contract, product and milestone changes are proposals only. Prior review conclusions and regression expectations are not evaluator inputs for a discovery test.
