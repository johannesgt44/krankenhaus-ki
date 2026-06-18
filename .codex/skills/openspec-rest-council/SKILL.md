---
name: openspec-rest-council
description: Use when the user wants to prepare, plan, or review a REST server implementation with OpenSpec using a council/panel of perspectives, especially for seminar, exam, workshop, or assignment situations where the language, scripting language, platform, or framework may be announced shortly before implementation. Use GitHub reference implementations for Python/FastAPI and TypeScript/Hono when relevant.
---

# OpenSpec REST Council

Use this skill to prepare REST server work with a small council of perspectives before creating or applying an OpenSpec change.

Keep the workflow lightweight. The council should improve the OpenSpec proposal, requirements, tasks, and review; it must not replace OpenSpec.

## Anti-Drift Rule

For every substantial response while this skill is active, start with:

`Marc, Modus: <mode>.`

Use one of these modes:

- `Council`
- `OpenSpec`
- `Implementation`
- `Review`
- `Coaching`

If this marker disappears, the user may interrupt with:

`Skill drift: return to openspec-rest-council.`

## Reference Projects

Use GitHub as the canonical source so the skill works across machines:

- TypeScript/Hono reference: <https://github.com/johannesgt44/krankenhaus>
- Python/FastAPI reference: <https://github.com/marculm/krankenhaus>

Read `references/rest-reference-map.md` before relying on the reference projects.

If network access is unavailable, ask the user whether local clones exist and where they are. Do not assume local paths are portable.

## Council Roles

Use these perspectives when the user asks for a council, panel, debate, or multiple views:

- **REST Architect**: routes, resources, HTTP methods, status codes, JSON shapes, Problem Details.
- **Framework Adapter**: map the REST blueprint to the announced language/framework.
- **OpenSpec Writer**: turn decisions into proposal, design, specs, and tasks.
- **Testing Critic**: define minimal integration checks and edge cases.
- **Seminar Coach**: keep the plan explainable, small, and practical under time pressure.

If sub-agents are available and the user explicitly asks for agents or a council, use separate agents for independent perspectives. Otherwise, simulate the perspectives clearly in one response.

## OpenSpec Integration

Use the council before creating an OpenSpec proposal or before applying a risky change.

Council output should become:

- requirements
- design decisions
- implementation tasks
- test checklist

After the council summary, continue with the normal OpenSpec workflow when the user asks to propose or implement.

## Guardrails

- Do not implement code unless the user explicitly asks for implementation.
- Do not over-expand the council; keep it short when time matters.
- Prefer the existing OpenSpec skills for proposal, apply, sync, and archive steps.
- Prefer GitHub references over local paths.
