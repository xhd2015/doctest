# readStdinIfPresent Error Propagation Tests

Verify that errors from `os.Stdin.Stat()` and `io.ReadAll()` inside
`readStdinIfPresent()` propagate to callers (`runAgentImplement`,
`runAgentDesign`) instead of being silently swallowed.

## Decision Tree

```
tests/read-stdin-error/                          [Request{Args, StdinFile}]
│                                                Run: replaces os.Stdin, calls cli.Run(args)
├── implement/                                   [prepends "agent","implement" to args]
│   │                                            (via runAgentImplement → readStdinIfPresent)
│   ├── stat-error/                              → err != nil (Stat on closed file propagated)
│   └── read-error/                              → err != nil (ReadAll on directory propagated)
└── design/                                      [prepends "agent","design" to args]
    │                                            (via runAgentDesign → readStdinIfPresent)
    ├── stat-error/                              → err != nil (Stat on closed file propagated)
    └── read-error/                              → err != nil (ReadAll on directory propagated)
```

## Test Index

| Leaf | Description |
|------|-------------|
| `implement/stat-error` | Closed file as stdin during `agent implement` — Stat error propagated |
| `implement/read-error` | Directory as stdin during `agent implement` — ReadAll error propagated |
| `design/stat-error` | Closed file as stdin during `agent design` — Stat error propagated |
| `design/read-error` | Directory as stdin during `agent design` — ReadAll error propagated |

## Coverage Notes

- Happy paths (terminal, pipe with data, empty pipe) are covered by `agent-stdin/`.
- This tree fills the gap for **error propagation** — the two ignored errors in `readStdinIfPresent()`.
- Both callers (`runAgentImplement` and `runAgentDesign`) are tested for both error sources.

## How to Run

```sh
doctest test ./libdoc/cli/tests/read-stdin-error/
```
