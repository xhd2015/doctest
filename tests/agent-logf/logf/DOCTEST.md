# Logf: Timestamped Logging

Verify that `subagent.Logf` produces `[2006-01-02T15:04:05]` prefixed output
with correct newline handling, format verbs, and special characters.

## Decision Tree

```
tests/agent-logf/logf/
├── DOCTEST.md                     # This file
├── SETUP.md                       # Root: Request/Response, Run calls subagent.Logf
├── without-trailing-newline/      # Message without \n → newline appended
├── with-trailing-newline/         # Message with \n → no double newline
├── empty-message/                 # Empty string → just timestamp + \n
├── format-verbs/                  # Format string with %s/%d verbs and args
└── special-chars/                 # Multiline message, special characters
```

## Test Index

| Leaf | Description |
|------|-------------|
| `without-trailing-newline` | Message without `\n` gets `\n` appended |
| `with-trailing-newline` | Message with `\n` keeps exactly one `\n` |
| `empty-message` | Empty format string produces timestamp + `\n` |
| `format-verbs` | Format verbs (`%s`, `%d`) resolved correctly |
| `special-chars` | Multiline content, special characters preserved |

## How to Run

```sh
doctest test tests/agent-logf/logf/
```
