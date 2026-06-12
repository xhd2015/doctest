## Setup

- Redefines `myHelper` with the same name as the root.
- The generator emits both as closures using `:=`,
  causing "no new variables on left side of :=".

```go
func myHelper(s string) string {
    return "child: " + s
}

func Setup(t *testing.T, req *Request) error {
    req.Name = myHelper("name")
    return nil
}
```
