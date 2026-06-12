## Steps

- Call Run(t, req) from inside Setup.
- The generator emits Run as lowercase `run`, so uppercase `Run`
  in the Setup body will be undefined.

```go
func Setup(t *testing.T, req *Request) error {
    req.Name = "leaf"
    resp, runErr := Run(t, req)
    _ = resp
    _ = runErr
    return nil
}
```
