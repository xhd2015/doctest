# Scenario

**Feature**: gen-side integrity damage forces cache miss even when fingerprint matches

```
first write (fingerprint + bridges OK)
delete gen go.mod OR a listed placeholder
second write with same sources
  -> rebuild; missing artifact restored
```

## Preconditions

- Source inputs unchanged; only gen integrity fails.

## Steps

1. First write establishes warm-eligible fingerprint.
2. Leaf sets DeleteGenGoMod or DeletePlaceholder for Run.

```go
import (
	"testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Mode = "write-second"
	req.ModPath = defaultModPath
	req.HasMod = true
	return nil
}
```
