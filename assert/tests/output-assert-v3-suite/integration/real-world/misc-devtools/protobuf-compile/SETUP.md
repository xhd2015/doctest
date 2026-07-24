# Scenario

**Feature**: protoc

```
# protoc
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"",
		"Writing output to out\\.pb\\.go",
	)
	req.Actual = "Writing output to out.pb.go"
	return nil
}
```
