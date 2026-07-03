# Scenario

**Feature**: g++ link

```
# g++ -o app
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"",
		"g++ -o app main.cpp",
	)
	req.Actual = "g++ -o app main.cpp"
	return nil
}
```
