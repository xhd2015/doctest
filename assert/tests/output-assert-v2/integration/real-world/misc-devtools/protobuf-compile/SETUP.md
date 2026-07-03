# Scenario

**Feature**: protoc

```
# protoc
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"",
		"Writing output to out.pb.go",
	)
	req.Actual = "Writing output to out.pb.go"
	return nil
}
```
