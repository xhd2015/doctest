# Dot Progress Tests

Tests that `build.Test` prints dot progress **incrementally** — one dot per
test package as it completes — rather than batching them all after `go test`
finishes.

## How to Run

```sh
doctest test ./libdoc/build/tests/dot-progress/...
```
