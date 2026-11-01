# Saum

A **Factory** is Saum's top-level concept. Create, rename, or delete Factories
from a sortable, filterable table. A Factory contains **User stories**, not Projects.
Each story has a title and a high-level goal: what the user should be able to
achieve, not how to implement it. Open a Factory to create and browse its stories.
Factories containing User stories cannot be deleted.

Stories are planned and implemented by Codex, one work item at a time. A factory agent prepares an executive summary; only human approval advances the next item.

## Shared factory workspace

Each factory owns one Git repository and shared application. Stories use separate
worktrees, reused across retries. Before a story turn, Saum brings in accepted
factory main changes, checkpointing local work on the story branch. Conflicts
are left for the story agent to resolve; unfinished merges cannot be reviewed.
Other stories' unapproved changes never flow into main automatically. Final
human approval queues integration into a clean, unchanged factory main.

Codex sessions and recorded verification commands run with full host permissions,
including networking and localhost listeners. Only run trusted code: the workflow
approval gates are not an OS security boundary against arbitrary local programs.

## Develop

Requires Bun 1.4+ and Go 1.26+.

```sh
bun install
bun dev
```

Open **http://127.0.0.1:5173**. The same command starts the Go API on
`127.0.0.1:8080`; Ctrl-C stops both.

Factories and their User stories are stored in SQLite at `.saum/factory.sqlite`,
which is not committed. Existing entries from the earlier Project prototype are
preserved as User stories, keeping their IDs, titles, and Factory ownership.
Imported entries may have no goal; no goal is invented for them. Any original
`.saum/projects.json` is left untouched.
The Go driver is pure Go; no database service or C toolchain is required.

## Layout

- `api/` — Go server.
- `sites/dashboard/` — React + Kumo dashboard.
- `packages/ui/` — reusable sortable, filterable tables.
- `examples/` — reserved for projects built with Saum.
- `architecture/` — original design dossier, unchanged.

## Checks

```sh
bun run check   # Oxfmt, Oxlint, strict TypeScript, table tests, Go race tests
bun run build
```

The website builds to `sites/dashboard/dist/`. Cloudflare deployment is not wired
up yet. The API remains local.
