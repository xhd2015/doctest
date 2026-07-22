# Scenario

**Feature**: sessionmod.ContentMD5 equals MD5 of the embedded source file on disk

```
ContentMD5() -> hex matches file hash of sessionmod embedded blob
```

## Steps

1. Read embedded source file and call ContentMD5().

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    req.RunKind = "sessionmod-md5"
	return nil
}
```
