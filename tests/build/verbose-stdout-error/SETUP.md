## Preconditions
- Running `doctest test -v` triggers the verbose code path where
  `goTestBuild.Stdout` and `goTestBuild.Stderr` are set to `w`,
  then `CombinedOutput()` tries to set them again → "Stdout already set".

## Steps
1. Run `doctest test -v` on the testdata fixture.

```go
import "path/filepath"

func Setup(t *testing.T, req *Request) error {
    req.Args = []string{"test", "-v", filepath.Join(DOCTEST_ROOT, "build", "testdata", "verbose-stdout-error")}
    return nil
}
```
