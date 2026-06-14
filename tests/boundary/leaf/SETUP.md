## Preconditions
- The parent root defines Request{} and Response{} with a stub Run that returns an error.

## Steps
1. No additional setup needed; Run is called automatically with an empty Request.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    t.Log("boundary: parent root leaf, stub Run will return error")
    return nil
}
```
