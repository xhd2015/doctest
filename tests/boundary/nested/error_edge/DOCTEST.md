# Deeply Nested Root — DOCTEST.md Boundary

This is a deeply nested self-contained test root. It verifies that
DOCTEST.md boundaries work at any depth — even inside another nested root.

This root defines its own `Request{ID, Data}` and `Response{Status, Message}`
types, which are completely different from any ancestor types.

## How to Run

```sh
doctest test ./tests/boundary/nested/error_edge/
```
