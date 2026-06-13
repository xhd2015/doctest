## Preconditions
- `--list-sessions` and `--session-id` are mutually exclusive.
- When both are given, stderr gets an error and exit code is 0.

## Steps
1. Run `doctest agent implement --list-sessions --session-id some-id`.

```go
func Setup(t *testing.T, req *Request) error {
    req.Args = []string{"agent", "implement", "--list-sessions", "--session-id", "some-id"}
    return nil
}
```
