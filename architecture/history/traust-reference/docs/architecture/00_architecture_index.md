# Architecture Index

This section is the architecture home for the deployer. It owns the system shape, runtime flows, cross-cutting concepts, tradeoffs, decisions, quality requirements, and risks.

Implementation contracts remain in [Specs](../specs/index.md). Architecture owns canonical architectural rules and their rationale. Exact canonical schemas, transition tables, deterministic reference rules, and low-level policy contracts live in the [Appendix](../appendix/index.md).

When architecture summarizes an exact appendix contract for reader ergonomics, the appendix remains authoritative for those details. Architecture should link to it rather than creating a competing full definition.

## Reading path

1. [Introduction and Goals](01_introduction_and_goals.md)
1. [Constraints](02_constraints.md)
1. [Context and Scope](03_context_and_scope.md)
1. [Solution Strategy](04_solution_strategy.md)
1. [Building Block View](05_building_block_view/index.md)
1. [Runtime View](06_runtime_view/index.md)
1. [Deployment View](07_deployment_view.md)
1. [Cross-Cutting Concepts](08_cross_cutting_concepts/index.md)
1. [Architecture Decisions](09_architecture_decisions.md)
1. [Quality Requirements](10_quality_requirements.md)
1. [Risks and Technical Debt](11_risks_and_debt.md)
1. [Glossary](12_glossary.md)

## Reference areas

- [Implementation Specs](../specs/index.md)
- [Canonical Appendix](../appendix/index.md)
- [Test Scenarios](../tests/0_index.md)

## Documentation ownership

| Area | Owns | Does not own |
|---|---|---|
| Architecture | System shape, runtime flow, concepts, canonical architectural rules and rationale, decisions, quality requirements, risks | Exact schemas, exact command contracts, deterministic reference tables |
| Specs | Implementation-level contracts, endpoint behavior, command invocation, integration mechanics | Architecture rules and rationale or canonical reference contracts |
| Appendix | Exact canonical schemas, transition tables, recovery order, persistence/write contracts, directory rules, and reference contracts | Architecture narrative or runtime walkthroughs |
| Tests | Scenario validation and expected evidence | Architecture rationale or canonical contracts |
