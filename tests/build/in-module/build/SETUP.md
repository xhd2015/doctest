# Scenario

**Feature**: doctest build command dumps generated files when internal imports detected

```
# internal import + build: temp compile, copy dump to --gen-dir
doctest build <tree> --gen-dir <module>/_gen -> dump *_test.go (no go.mod)
```

## Preconditions

- Operation mode is `doctest build`.

## Steps

1. Prefix doctest args with the `build` subcommand.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Args = append([]string{"build"}, req.Args...)
	return nil
}
```