# Doctest CLI Integration Tests

These doc-style tests specify the command-level contract for the doctest CLI.
They are executable integration tests: each runnable leaf invokes the real
doctest command, captures stdout, stderr, and exit status, and asserts concrete
observable behavior.

The tests intentionally use the public command boundary. Agent-oriented cases
configure fake Codex so no real LLM or network backend is required.

