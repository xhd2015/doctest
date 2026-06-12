# Basic Request Runner Example

This tree demonstrates executable `SETUP.md` and `ASSERT.md` snippets for
`doctest`.

Run it from the repository root:

```sh
cd agents/doctest && go run . test ./tests/testdata/basic-request-runner
```

## Tree

```text
basic-request-runner
├── SETUP.md
├── happy-path
│   ├── SETUP.md
│   └── ASSERT.md
├── expected-error
│   ├── SETUP.md
│   └── ASSERT.md
└── override-run
    ├── SETUP.md
    └── ASSERT.md
```

## Cases

| Path | Purpose |
| --- | --- |
| `happy-path` | Uses inherited root `Run` after leaf setup mutates `Request`. |
| `expected-error` | Uses inherited root `Run` and asserts the returned run error. |
| `override-run` | Defines a leaf `Run`, proving deepest `Run` overrides root `Run`. |
