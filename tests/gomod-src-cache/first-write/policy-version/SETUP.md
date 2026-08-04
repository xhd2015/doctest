# Scenario

**Feature**: fingerprint policy line is gomod-src=1; legacy gomod-fp stays absent

```
cold WriteGoModWithVendorBridges
  -> doctest.gomod-src starts with "version gomod-src=1"
  -> doctest.gomod-fp never written
```

## Steps

1. Inherit first-write cold fixture.
2. Assert policy version and absence of legacy skip file after measured write.

```go
import (
	"testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	if req.Mode != "write-once" {
		t.Fatalf("policy-version expects Mode write-once, got %q", req.Mode)
	}
	return nil
}
```
