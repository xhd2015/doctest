## Preconditions
- `--status` without `--session-id` must error with a clear message.
- Exit code must be 0 per the "always exit 0" requirement.

## Steps
1. Run `doctest agent implement --status` without `--session-id`.

```go
func Setup(t *testing.T, req *Request) error {
    req.Args = []string{"agent", "implement", "--status"}
    return nil
}
```
