# Scenario

**Feature**: ffmpeg -version

```
# ffmpeg
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template(
		"__VER__: 'type=string, example=6.0'\n",
		"ffmpeg version 6\\.0 Copyright",
	)
	req.Actual = "ffmpeg version 6.0 Copyright"
	return nil
}
```
