## Preconditions
- The path_resolve package is importable from this test tree.
- Each group overrides `Run` to call the specific function under test.

## Steps
1. Build the test binary that imports path_resolve.
2. Each leaf sets input fields via `Setup`.
3. Each leaf asserts results via `Assert`.

## Context
- These are direct unit tests, not CLI integration tests.
- The root `Run` is a stub; each group provides its own `Run`.

```go
type Request struct {
	Input    string
	BasePath string
}

type Response struct {
	BoolResult   bool
	StringResult string
	DirsResult   []string
	RootResult   string
	RootOkResult bool
	ErrResult    string
}

func Run(t *testing.T, req *Request) (*Response, error) {
	t.Fatal("not implemented — group must override Run")
	return nil, nil
}
```
