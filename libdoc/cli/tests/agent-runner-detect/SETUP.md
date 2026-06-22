## Preconditions
- The auto-detection logic is implemented in `github.com/xhd2015/agent-pro/agent/subagent`.
- The CLI defaults in `cli.go` are changed to empty `Options{}`.
- Tests call the CLI directly via `cli.Run()` and capture the returned error.

## Steps
1. Child SETUP.md files configure Args and Env for each test scenario.
2. Run sets env vars via os.Setenv, calls `cli.Run(req.Args)`, captures stdout/stderr and error.

```go
import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"
	"github.com/xhd2015/doctest/libdoc/cli"
)

```
