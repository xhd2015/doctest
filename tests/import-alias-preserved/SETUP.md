## Preconditions
- This test verifies that import aliases in Go code blocks are preserved during code generation.
- Inherits Request, Response, and Run from the parent SETUP.md chain.

## Steps
1. The Setup function uses an aliased import (`myos "os"`) and calls `myos.Environ()`.
2. If the alias is correctly preserved in generated code, compilation succeeds.
3. If the alias is lost, the generated code references `myos` which is undefined → compile error.

```go
import (
	"testing"
	myos "os"
)

func Setup(t *testing.T, req *Request) error {
	_ = myos.Environ()
	return nil
}
```
