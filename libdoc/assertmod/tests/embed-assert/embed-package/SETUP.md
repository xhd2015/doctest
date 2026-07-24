# Scenario

**Feature**: libdoc/assertmod ContentMD5 matches on-disk assert.go hash

```
# go:embed assert.go
ContentMD5() == md5(libdoc/assertmod/assert.go)
```

## Preconditions

- `libdoc/assertmod/assert.go` is generated and committed.

## Steps

1. Set `runKind = "assertmod-md5"`.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    req.RunKind = "assertmod-md5"
	req.SecondRun = false
	return nil
}
```