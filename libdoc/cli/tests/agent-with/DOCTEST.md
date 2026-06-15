# `doctest agent with` Subcommand Tests

Verify the new `doctest agent with --agent-runner=RUNNER [--model=MODEL] <prog> <args...>` subcommand.

This subcommand sets `DOCTEST_SUBAGENT_AGENT_RUNNER` (and optionally `DOCTEST_SUBAGENT_MODEL`) in the child process environment, then executes `<prog>` with inherited stdin/stdout/stderr.

## Decision Tree

```
tests/agent-with/                             [Request{Args, Env}, Response{ExitCode, Stdout, Stderr, Err}]
│                                              Run: replaces os.Stdout/Stderr, calls cli.Run(args)
├── errors/                                    [prepends "agent","with" to Args]
│   ├── missing-agent-runner/                  --agent-runner flag present but no value
│   ├── missing-prog/                          no <prog> positional argument
│   ├── model-without-value/                   --model flag present but no value
│   └── prog-not-found/                        <prog> does not exist in PATH
└── execution/                                 [prepends "agent","with","--agent-runner=opencode" to Args]
    ├── basic/                                 echo hello → stdout "hello\n", exit 0
    ├── with-model/                            --model=gpt-4 → env has DOCTEST_SUBAGENT_MODEL=gpt-4
    ├── with-extra-args/                       extra args forwarded to child
    └── exits-with-code/                       child exits 42 → exit code 42 propagated
```

## Test Index

### Errors — 4 leaves
| Leaf | Description |
|------|-------------|
| `errors/missing-agent-runner` | `--agent-runner` without a value → error `--agent-runner requires a value` |
| `errors/missing-prog` | No `<prog>` argument → error `agent with requires <prog>` |
| `errors/model-without-value` | `--model` without a value → error `--model requires a value` |
| `errors/prog-not-found` | Nonexistent program → error `executable file not found in $PATH` |

### Execution — 4 leaves
| Leaf | Description |
|------|-------------|
| `execution/basic` | Runs `echo hello`, stdout contains `hello`, exit code 0 |
| `execution/with-model` | `--model=gpt-4` set, child sees `DOCTEST_SUBAGENT_MODEL=gpt-4` |
| `execution/with-extra-args` | Extra args after `<prog>` forwarded to child |
| `execution/exits-with-code` | Child exits 42, parent exit code is 42 |

Total: **8 leaves**.

## How to Run

```sh
doctest test ./libdoc/cli/tests/agent-with/
```
