# Scenario

**Feature**: pg_dump

```
# pg_dump
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template(
		"",
		"-- PostgreSQL database dump",
	)
	req.Actual = "-- PostgreSQL database dump"
	return nil
}
```
