# Scenario

**Feature**: committed raw-source cache key matches embed script output for assert/*.go

```
# go generate / embed-assert -cache-key
RawSourceCacheKeyMD5() == md5(sorted assert/*.go raw bytes)
```

## Preconditions

- `libdoc/assertmod/cache_key.go` is generated and committed.

## Steps

1. Set `runKind = "embed-cache-key"`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	runKind = "embed-cache-key"
	req.SecondRun = false
	return nil
}
```