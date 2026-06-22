# Logf: Timestamped Logging

## Version
0.0.2


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

```go
import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
	"github.com/xhd2015/agent-pro/agent/subagent"
)

type Request struct {
	Args	[]string
	Env	[]string
}
type Response struct {
	Stdout string
}
func Run(t *testing.T, req *Request) (*Response, error) {
	fmtStr := "default"
	for _, e := range req.Env {
		if strings.HasPrefix(e, "LOGF_FORMAT=") {
			fmtStr = e[len("LOGF_FORMAT="):]
			break
		}
	}

	var args []any
	for _, a := range req.Args {
		args = append(args, a)
	}

	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		return nil, fmt.Errorf("create pipe: %w", err)
	}
	os.Stdout = w

	if len(args) > 0 {
		subagent.Logf("%s", fmt.Sprintf(fmtStr, args...))
	} else {
		subagent.Logf("%s", fmtStr)
	}

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	if _, readErr := buf.ReadFrom(r); readErr != nil {
		return nil, fmt.Errorf("read pipe: %w", readErr)
	}

	return &Response{Stdout: buf.String()}, nil
}
```
