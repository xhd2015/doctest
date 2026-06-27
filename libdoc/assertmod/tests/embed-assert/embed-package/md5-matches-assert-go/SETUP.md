# Scenario

**Feature**: assertmod.ContentMD5 equals MD5 of embedded assert.go file

```
# accessor contract
ContentMD5() -> hex matches file hash
```

## Steps

1. Read `libdoc/assertmod/assert.go` and call `assertmod.ContentMD5()`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.SecondRun = false
	return nil
}
```