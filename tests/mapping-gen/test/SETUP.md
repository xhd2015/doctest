# Scenario

**Feature**: the doctest binary is built by root Setup

```
# maps doc-style leaves to Go test packages
doctest build/test -> mapping-gen -> leaf <-> Go package

# gen-dir modes
auto gen-dir -> one package per leaf | explicit gen-dir -> user-specified layout
```

## Preconditions
- The doctest binary is built by root Setup.
- Timeout is set to 120s.

## Steps
1. Prefix args with "test" command.

```go
import (
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	req.Args = append([]string{"test"}, req.Args...)
	return nil
}
```
