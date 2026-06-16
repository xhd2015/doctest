# Scenario

**Feature**: `PROGRESS_FILE` environment variable is NOT set

```
# sub-agents report progress to a file
sub-agent --writes--> progress file (env var DOCTEST_PROGRESS_FILE)

# multiple entries append
each step -> structured JSON entry -> append to file
```

## Preconditions
- `PROGRESS_FILE` environment variable is NOT set.

## Steps
1. Invoke `report-progress` without `PROGRESS_FILE` env var.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    req.Args = []string{"some progress"}
    return nil
}
```
