# Doctest CLI Integration Tests

These doc-style tests specify the command-level contract for the doctest CLI.
They are executable integration tests: each runnable leaf invokes the real
doctest command, captures stdout, stderr, and exit status, and asserts concrete
observable behavior.

The tests intentionally use the public command boundary. Agent-oriented cases
configure fake Codex so no real LLM or network backend is required.

## DSN (Domain Specific Notion)

### Participants
- **`doctest`** — the CLI binary under test; every test builds it fresh, then
  invokes it as a subprocess. It is the single entry point for all behaviors.
- **Test tree** — a directory hierarchy of `.md` files (DOCTEST.md, SETUP.md,
  ASSERT.md) that the CLI reads, interprets, and executes. It models a decision
  tree of scenarios.
- **Fake Codex** — a stand-in for the LLM backend used during agent operations
  (`design`, `implement`). It returns predetermined responses so no real
  network or AI model is involved.
- **The file system** — the doctest binary reads and writes files; tests observe
  side effects (generated code, output files, mapping data).

### Behaviors
- **`build`** — compiles source code from a doc-style test directory into an
  executable test binary.
- **`test`** — runs the compiled test binary, reports pass/fail per leaf, and
  propagates exit codes.
- **`agent`** — orchestrates design and implement sub-agents. The design agent
  writes tests; the implement agent writes implementation code. They
  communicate through session state and progress reports.
- **`skill`** — exposes embedded documentation (doc-spec, code-spec) to users.
- **`help`** — prints usage information at top-level and per subcommand.
- **`vet`** — inspects and validates test tree structure.

