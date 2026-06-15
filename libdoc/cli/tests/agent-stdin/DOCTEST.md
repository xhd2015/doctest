# Agent Stdin Input Tests

Verify that `doctest agent implement` and `doctest agent design` accept
prompts from stdin (pipe/heredoc) when no positional arguments are given.

## Decision Tree

```
tests/agent-stdin/                          [Request{Args, Stdin, StdinSource, ReqContent}]
│                                           Run: replaces os.Stdin, calls cli.Run(args)
├── implement/                              [prepends "agent","implement" to args]
│   ├── no-input/                           → err contains "requires <prompt>"
│   ├── stdin-only/                         → err NOT contains "requires <prompt>"
│   ├── positional-only/                    → err NOT contains "requires <prompt>"
│   ├── positional-wins/                    → err NOT contains "requires <prompt>" (pos wins)
│   ├── stdin-heredoc/                      → err NOT contains "requires <prompt>"
│   ├── stdin-empty-pipe/                   → err contains "requires <prompt>" (empty pipe = no prompt)
│   ├── requirement-only/                   → err NOT contains "requires <prompt>"
│   ├── requirement-positional/             → err NOT contains "requires <prompt>"
│   ├── requirement-stdin/                  → err NOT contains "requires <prompt>"
│   └── requirement-missing/                → err contains "read requirement file"
└── design/                                 [prepends "agent","design" to args]
    ├── no-input/                           → err contains "requires <prompt>"
    ├── stdin-only/                         → err NOT contains "requires <prompt>"
    └── positional-only/                    → err NOT contains "requires <prompt>"
```

## Test Index

| Leaf | Description |
|------|-------------|
| `implement/no-input` | No positional args, stdin is terminal → "requires <prompt>" |
| `implement/stdin-only` | Stdin pipe with content, no positional → prompt from stdin |
| `implement/positional-only` | Positional arg present, no stdin → prompt from args |
| `implement/positional-wins` | Both positional and stdin → positional wins, stdin ignored |
| `implement/stdin-heredoc` | Stdin simulates heredoc with multiline content |
| `implement/stdin-empty-pipe` | Stdin is pipe but empty → "requires <prompt>" |
| `implement/requirement-only` | --requirement file, no prompt → requirement as prompt |
| `implement/requirement-positional` | --requirement + positional → combined prompt |
| `implement/requirement-stdin` | --requirement + stdin → combined prompt |
| `implement/requirement-missing` | --requirement to nonexistent file → file read error |
| `design/no-input` | No positional args, stdin is terminal → "requires <prompt>" |
| `design/stdin-only` | Stdin pipe with content → prompt from stdin |
| `design/positional-only` | Positional arg present → prompt from args |

## How to Run

```sh
doctest test ./libdoc/cli/tests/agent-stdin/
```
