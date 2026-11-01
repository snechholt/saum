# Saum

A **Factory** is Saum's top-level concept. A Factory contains Projects.
For now, Projects are only names in a list; their behavior and internals have not
been built. There are no agents, issues, or architecture workflows.

## Develop

Requires Bun 1.4+ and Go 1.26+.

```sh
bun install
bun dev
```

Open **http://127.0.0.1:5173**. The same command starts the Go API on
`127.0.0.1:8080`; Ctrl-C stops both.

Projects are stored in SQLite at `.saum/factory.sqlite`, which is not committed.
The Go driver is pure Go; no database service or C toolchain is required.

## Layout

- `api/` — Go server.
- `sites/dashboard/` — React + Kumo dashboard.
- `packages/` — reserved for shared JS packages.
- `examples/` — reserved for projects built with Saum.
- `architecture/` — original design dossier, unchanged.

## Checks

```sh
bun run check   # Oxfmt, type-aware Oxlint, strict TypeScript, Go race tests
bun run build
```

The website builds to `sites/dashboard/dist/`. Cloudflare deployment is not wired
up yet. The API remains local.
