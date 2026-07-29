# Scenario

**Feature**: gen go.mod gets require + replace for each modules.txt module

```
modules.txt: example.com/dep v1.2.3, example.com/nogo v0.4.0
  -> require example.com/dep v1.2.3
  -> replace example.com/dep => <modRoot>/vendor/example.com/dep
  -> require example.com/nogo v0.4.0
  -> replace example.com/nogo => <modRoot>/vendor/example.com/nogo
```

## Steps

1. Inherit present vendor fixture.
2. Run WriteGoMod; assert requires and replaces for both sample modules.

```go
import (
	"testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	if req.SampleModPath == "" || req.NoGoModPath == "" {
		t.Fatal("present Setup must set SampleModPath and NoGoModPath")
	}
	return nil
}
```
