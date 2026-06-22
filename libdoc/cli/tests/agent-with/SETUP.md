## Preconditions
- The `doctest agent with` subcommand is implemented in `libdoc/cli/cli.go`.
- Tests call the CLI directly via `cli.Run()` after replacing `os.Stdout` and `os.Stderr` with pipes.

## Steps
1. Child SETUP.md files set Args and optionally Env.
2. Run replaces os.Stdout/os.Stderr with pipes (to capture child output), sets request env vars via os.Setenv, then calls `cli.Run(req.Args)`.
3. Stdout, Stderr, and ExitCode are parsed from the pipes and error.

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
