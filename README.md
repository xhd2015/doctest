> I once thought every domain of test needs a dedicated DSL, now they just become markdowns with code annotation

# doctest

Doc-style test tool: write test cases as markdown decision trees with embedded
Go code, then build and run them.

# Installation
```sh
curl -fsSL https://raw.githubusercontent.com/xhd2015/doctest/master/install.sh | bash

# or go install
go install github.com/xhd2015/doctest/cmd/doctest@latest
```

**Agent skill** (primary TDD flow):

```sh
doctest skill tdd --install --global
# add the skill to ~/.agents/skills/doctest-tdd/SKILL.md
```

# Quick start

**Agents** (after skill install, in `claude code`, `codex`, `grok`, `opencode` etc.):

```md
/doctest-tdd <your feature here>
```

The agent will propose plan first, after you confirmed, it will spawn a `[designer] sub-agent` for test case design, and confirm they're all RED; then spawn another `[implementer] sub-agent` for implementation.

**Run doctests:**

```sh
doctest test ./...
# discovery skips labeled leaves; full suite:
doctest test --label-all ./...
# label expr: doctest test --help
```

# Author notes

- Suite leaves run in one process under **`t.Parallel()`** — no process
  `Setenv`/`Chdir` for isolation (`doctest skill lint --show`,
  `doctest skill code-spec --show`).
- Prefer **in-process** library / `cli.RunWithWriter` leaves; reserve binary e2e
  for full integration (`doctest skill design-principle --show`).
- Inject context via `d *session.Doctest` fields; process cwd is undetermined
  (`doctest skill migrate --show`).

Tree layout, DSN (domain sketch), MECE, Setup/Assert:
`doctest skill doc-spec --show` and design content via designer / tdd-lite embeds.

# Commands

| Area | Commands |
|------|----------|
| Core | `vet`, `build`, `test`, `edit` |
| Skills | `skill --list` / `--show` / `--install`, `skills update` |
| Agents | `agent generate`, `fill-code`, `design`, `implement`, `with` |
| Metrics | `metrics path\|last\|top\|phases\|summary\|…` |

```sh
doctest --help
doctest test --help
doctest skill --help
```

# Skills

```sh
doctest skill --list
doctest skill tdd --show              # multi-agent TDD orchestrator
doctest skill tdd-lite --show         # single-agent TDD
doctest skill design-principle --show # L1 / L2 / L3
doctest skill review --show           # design review of trees
```

TDD ephemeral requirement paths, modes, and workflow:
`doctest skill tdd --show` (not duplicated here).

# Development

```sh
go test ./libdoc/...
go run ./cmd/doctest test -v ./tests
go build -o doctest ./cmd/doctest
```

Git hooks (optional):

```sh
go run github.com/xhd2015/git-hooks@latest pre-commit add "script-pre-commit" go run ./script/git-hooks/pre-commit
```

Release:

`go run ./script/github/release --dry-run` 

then without `--dry-run` (local `.upload-credentials.json`).
