# Short Display Paths — `build.Test` stderr integration

Integration doc tests verifying that `build.Test` stderr uses `DisplayPath` at
the real call sites: `→` gen-root announcement, `doctest:` header, and `cd`
command preview.

## DSN (Domain Specific Notion)

### Participants

- **`build.Test`** — discovers doctest leaves, generates Go tests under a gen
  root, prints progress to stderr.
- **`announceRoots`** — prints `→ <genRoot>` as the first stderr line.
- **Progress header** — prints `doctest: <dir>` and `─── N test cases`.
- **`cd` preview** — prints `cd <runDir> && go test ...` before executing.
- **`DisplayPath`** — display-only formatter applied at each stderr path print.

### Behaviors

- **Auto gen dir** — mapping-gen cache under home; `→` and `cd` lines use `~/...`
  instead of `/Users/...`.
- **Explicit gen dir under cwd** — user `--gen-dir` under project; `→` line uses
  `./_gen` when cwd is the project root.
- **Test dir under cwd** — `doctest:` line uses `./` prefix for the source tree path.

## Decision Tree

```
short-display-paths
└── gen-dir-source              [how gen root is chosen]
    ├── mapping-gen-cache       auto cache dir → ~/.../mapping-gen/...
    └── explicit-gen-dir-under-cwd  --gen-dir ./_gen under project
```

## Test Index

| Leaf | Description |
|------|-------------|
| `gen-dir-source/mapping-gen-cache` | Auto gen dir stderr uses `~/...mapping-gen...`; `doctest:` uses `./`; no raw home absolute |
| `gen-dir-source/explicit-gen-dir-under-cwd` | Explicit `_gen` under project displays as `→ ./_gen` |

## How to Run

```sh
doctest vet ./libdoc/build/tests/short-display-paths/...
doctest test ./libdoc/build/tests/short-display-paths/...
```