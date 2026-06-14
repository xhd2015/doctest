# Agent Logf: Unified Event Stream Logging

Verify that all event stream output from sub-agents (traceSession, showStatus)
uses `Logf` for consistent timestamped logging `[2006-01-02T15:04:05]`, while
non-event UI framing (borders, headers) continues without timestamps via
`fmt.Fprintf(os.Stdout, ...)`.

## Decision Tree

```
tests/agent-logf/
├── DOCTEST.md                     # This file
├── SETUP.md                       # Root: stub Run (RED until implemented)
│
├── logf/                          # === What function is being tested? Logf ===
│   ├── SETUP.md                   # Capture stdout via pipe, call subagent.Logf
│   ├── without-trailing-newline/  # Message without \n → newline appended
│   ├── with-trailing-newline/     # Message with \n → no double newline
│   ├── empty-message/             # Empty string → just timestamp + \n
│   ├── format-verbs/              # Format string with %s/%d verbs and args
│   └── special-chars/             # Multiline message, special characters
│
├── trace-session/                 # === What function is tested? traceSession ===
│   ├── SETUP.md                   # Create session dirs, run --trace
│   ├── no-events-file/            # No events.jsonl → (no events yet) via Logf
│   └── with-events/              # Events exist, session finished → event+done lines via Logf
│
└── show-status/                   # === What function is tested? showStatus ===
    ├── SETUP.md                   # Create session dirs, run --status
    ├── session-not-found/         # No session → stderr error, no timestamp
    ├── no-events/                 # Session found, no events → "No events yet" via Logf
    └── with-events/              # Session found, events exist → event lines via Logf
```

## Test Index

### Logf (5 leaves)
| Leaf | Description |
|------|-------------|
| `logf/without-trailing-newline` | Message without `\n` gets `\n` appended |
| `logf/with-trailing-newline` | Message with `\n` keeps exactly one `\n` |
| `logf/empty-message` | Empty format string produces timestamp + `\n` |
| `logf/format-verbs` | Format verbs (`%s`, `%d`) resolved correctly |
| `logf/special-chars` | Multiline content, special characters preserved |

### traceSession (2 leaves)
| Leaf | Description |
|------|-------------|
| `trace-session/no-events-file` | No events file: "(no events yet)" + Done via Logf, borders without timestamps |
| `trace-session/with-events` | Events present: event lines via Logf, Done via Logf, borders without timestamps |

### showStatus (3 leaves)
| Leaf | Description |
|------|-------------|
| `show-status/session-not-found` | Session not found: stderr error, no timestamp |
| `show-status/no-events` | No events: header block without timestamps, "No events yet" via Logf |
| `show-status/with-events` | Events present: header block without timestamps, event lines via Logf |

## How to Run

```sh
doctest test -v ./tests/agent-logf/
```
