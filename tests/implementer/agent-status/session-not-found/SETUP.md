## Preconditions
- A session ID that does not correspond to any session directory.
- Exit code must be 0 per the "always exit 0" requirement.

## Steps
1. Run `doctest agent implement --status --session-id nonexistent-status-test`.

```go
func Setup(t *testing.T, req *Request) error {
    req.Args = []string{"agent", "implement", "--status", "--session-id", "nonexistent-status-test"}
    return nil
}
```
