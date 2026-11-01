# UI conventions

Read [.agents/skills/kumo-design/SKILL.md](.agents/skills/kumo-design/SKILL.md)
before changing the dashboard or shared UI. This is the official
[Cloudflare Kumo design skill](https://github.com/cloudflare/kumo/tree/main/skills/kumo-design).

- Use Kumo components and semantic color tokens. Custom CSS is for layout, not a replacement design system.
- Content and controls use 14px text. Density comes from spacing, not tiny type.
- Do not override root font size, letter spacing, or Kumo control typography.
- Keep action buttons icon-only with accessible names and tooltips.
- Keep the interface minimal. Do not add explanatory copy or project features without asking.
- Run checks, build, and inspect the actual browser UI before calling a visual change done.
