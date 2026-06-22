## Preconditions
- The doctest binary is built and available.
- Tests run the CLI directly by calling `cli.Run()` after replacing `os.Stdin`.

## Steps
1. Child SETUP.md files set up stdin source, args, and requirement files.
2. Run replaces `os.Stdin` (pipe, /dev/null, or keep original), calls `cli.Run(req.Args)`, captures stdout/stderr/error.

```go
import (
	"bytes"
	"io"
	"os"
	"testing"
	"github.com/xhd2015/doctest/libdoc/cli"
)

```
