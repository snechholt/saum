# Example 05. School planning and other governance shapes

**Status:** conceptual examples used to test scope and control resolution, not researched recommendations for operating a school or company.

## 1. School strategy, annual plan and grade plans

The user proposed a school principal governing strategy documents, an annual school plan below them, and grade-level implementation plans further down.

```text
School strategy
       |
       v
Annual school plan
   /      |      \
Grade 1 Grade 2 Grade 3 plans
```

A change in strategy creates a reconciliation question for the annual plan. Changes there can affect grade plans. Each project may use teachers, workshops, documents or domain-specific tools rather than coding agents.

The manager's role follows from the artifact and objective. Nothing requires a universal Principal Agent persona. The same questions apply as in software: which revision governs the plan, what changed, what must be preserved, what evidence supports the candidate, and who can accept it?

## 2. Alternatives and workshops as experimentation

The user observed that producing several annual-plan alternatives and discussing them in a workshop resembles the experiment/review/settlement pattern. Saum could request isolated alternatives, summarize meaningful differences and record the human choice.

That does not imply every school process must use this workflow. A project may already have its own alternatives process. The point is that generic candidate comparison can make sense outside software and therefore may not require knowledge of worktrees or coding agents in the core abstraction.

The example also exposes limits: the meaning of validation is domain-specific. Evidence that a plan follows strategy is different from evidence of educational effectiveness. A generic coordinator cannot simply substitute software-test notions for domain judgment.

## 3. Company and product hierarchy

The main thread also sketched company objectives, product strategy, product architecture, web/data projects and an operating or go-to-market plan. These need not be a fixed tree of roles; some projects may have multiple governing inputs such as strategy and budget.

This tests the difference between authority dependencies and ordinary work dependencies. A plan can be based on accepted strategy and budget revisions even if no task is currently blocked on either.

The generic organizational ceiling is an architectural possibility, not a decision to market the first implementation as a CEO automation platform.

## 4. TODO application foundation

The original inventory used a rich TODO application specification to show that foundation selection should be project-specific. A representative first slice might boot a real HTTP server, initialize persistence, support signup/login and demonstrate the intended restart behavior, while deferring TODO CRUD, sharing and reminders.

The point is not that every app must begin with authentication. It is to identify the smallest slice that pressures enough core assumptions to establish useful structure. A different product could have a completely different frontier.

Experiments can compare approaches; review can accept a coherent baseline; later work should build on that structure rather than rediscover it.

## 5. Optional interaction-model authority

One possible software graph is Architecture -> Component Design -> Interaction Models -> Implementation. The middle artifact could contain sequence diagrams, state machines or component interaction contracts.

It is useful when the human wants to govern those interactions explicitly or several consumers need them. If that level of control is unnecessary, the layer can be omitted and implementation receives more discretion.

The same idea applies to interface/header-like artifacts such as IDLs, API schemas or language-neutral contracts. The representation matters less than the fact that the boundary is intentionally settled and downstream internals remain flexible.

## 6. The shared lesson

Project boundaries should reflect independently governed knowledge, not a universal team chart or every internal task. Different workflows can fit behind the interface, but resulting claims still need appropriate evidence and acceptance.

The examples help expose over-generalization too. A concept expressed only in terms of a Git branch, unit test or named engineering role may be infrastructure rather than a universal semantic primitive.

## Sources

[S01](../history/sources/01-original-manager.md), sections 1, 3, 6, 49, 51-52; [S02](../history/sources/02-project-boundary-notes.md), school alternatives/workshop comparison; main-thread M04-M06.
