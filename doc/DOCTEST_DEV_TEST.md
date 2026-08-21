---
name: doctest-dev-test
description: single-agent development followed by significance-first doctest design and verification
metadata:
  version: "0.1.0"
---

--begin of skill doctest-dev-test--

# Dev-test

Develop the requested change first, then add tests and verify it. Work as one
agent: inspect, edit, test, and verify directly without sub-agent delegation.

## Rules

- Understand the requested behavior and inspect relevant code and existing tests.
  For a bug, identify the mechanism before editing.
- Write the implementation before writing new tests. Existing tests may be read
  to learn contracts.
- Tests are mandatory, but may be GREEN on their first run. Do not manufacture
  RED, stage or seal tests, or preserve a wrong test merely because it was added.
- Make the smallest sound change. Keep unrelated refactors out of scope.
- For plan phases, complete develop → test → verify for each phase in order.
- Do not commit, push, open a pull request, deploy, or release unless requested.

## Workflow

1. **Understand** — determine scope, current behavior, and verification surfaces;
   ask only when a decision blocks correct implementation.
2. **Develop** — implement the behavior directly, following project conventions.
3. **Design tests** — after implementation, add coverage using **Test design**.
4. **Test** — run the narrowest relevant tests first and fix genuine code or test
   defects until they pass.
5. **Verify** — run broader regression checks and user-facing scenarios.
6. **Report** — list changed behavior, tests run, and anything unverified.

For doctest-backed changes, adapt paths to the feature:

```sh
doctest vet ./tests/<feature>/
doctest test ./tests/<feature>/
doctest test ./...
```

Also run focused package tests and the project suite where applicable. For web
changes, exercise the flow in a browser, check related routes and edge states,
and check desktop and mobile layouts. If a required check cannot run, state why;
do not claim it passed.

## Test design

**Significance-first:** put the distinction with the greatest user-visible effect
nearest the root. Nest secondary variants beneath the behavior they refine. Do
not organize primarily by files, functions, or incidental setup. Each leaf has
one clear behavioral outcome.

**Mutually exclusive and collectively exhaustive (MECE):** sibling branches do
not overlap and together cover all meaningful outcomes in the declared scope.
Each scenario belongs to one sibling; intentional exclusions are stated. Keep
shared preconditions in parent setup instead of repeating them in every leaf.

| Layer | Target share | Use for |
|-------|-------------:|---------|
| **L1 Go test** | 10–20% | Pure helpers and flat edge-case tables |
| **L2 in-process doctest** | **70–85%** | Main public-behavior and short-CLI scenarios |
| **L3 e2e doctest** | 5–10% | Sparse full integration with a required process boundary |

Default to L2. Use L1 when a flat table is clearer. Use L3 only when a real
binary, install layout, or multi-step integration is essential; every L3 leaf
has `label: e2e`. Help, argument validation, and other short paths stay L2.
Choose the cheapest layer that reliably fails when the behavior is wrong.
Keep in-process tests parallel-safe: inject writers, environment, and paths
instead of mutating process-global stdio, environment, or working directory.

--end of skill doctest-dev-test--
