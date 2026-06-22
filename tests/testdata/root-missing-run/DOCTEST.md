# Root Missing Run Fixture

## Version
0.0.2

## DSN (Domain Specific Notion)

### Participants
- **fixture** — validation harness missing Run.

### Behaviors
- **discover** — should fail without func Run.

```go
import "testing"

type Request struct {
	Input string
}

type Response struct {
	Output string
}
```