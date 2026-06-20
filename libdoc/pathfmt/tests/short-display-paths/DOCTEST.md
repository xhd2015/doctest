# Short Display Paths — `pathfmt.DisplayPath`

Unit-style doc tests for the display-only path formatter used by doctest CLI
output. Internal filesystem operations keep real absolute paths; only user-facing
stderr lines are shortened.

## DSN (Domain Specific Notion)

### Participants

- **`DisplayPath`** — pure formatter: accepts a filesystem path string, returns a
  shorter string for human-readable output only.
- **`cwd`** — process working directory from `os.Getwd()`, absolutized before
  comparison.
- **`home`** — user home directory from `os.UserHomeDir()`, used as a `~` prefix
  when the path is outside the cwd subtree but still under home.
- **Call sites** (`announceRoots`, `doctest:` header, `cd` command preview) —
  consumers that print paths to stderr; not covered in this tree (see
  `libdoc/build/tests/short-display-paths/`).

### Behaviors

- **Normalize** — `filepath.Abs(path)`; empty or Abs error → return input unchanged.
- **Under cwd** — path equals cwd → `"."`; strict child → `"./" + rel` (rel must not
  start with `".."`).
- **Home shorten** — when not under cwd, if path has prefix `home + sep`, replace
  with `"~"`.
- **Fallback** — return the absolute path unchanged (e.g. `/var/folders/...`).

## Decision Tree

```
short-display-paths
├── cwd-relative                 [path is under cwd subtree]
│   ├── at-cwd                   path == cwd → "."
│   ├── child-path               strict child → "./sub"
│   ├── deep-nested-child        nested child → "./a/b/c"
│   ├── relative-input           relative input under cwd → "./rel"
│   └── parent-path              parent of cwd → home or absolute (not "./..")
├── home-shorten                 [path outside cwd, under home]
│   ├── at-home                  path == home → "~"
│   └── under-home               cache-like path → "~/..."
└── edge-inputs                  [normalization edge cases]
    ├── empty-path               "" → "" (unchanged)
    └── outside-home             temp dir outside home → absolute unchanged
```

## Test Index

| Leaf | Description |
|------|-------------|
| `cwd-relative/at-cwd` | Absolute path equal to cwd displays as `"."` |
| `cwd-relative/child-path` | Strict child of cwd displays as `"./child"` |
| `cwd-relative/deep-nested-child` | Deeply nested child displays as `"./a/b/c"` |
| `cwd-relative/relative-input` | Relative input resolved under cwd displays as `"./rel"` |
| `cwd-relative/parent-path` | Parent of cwd is not shortened to `./..`; uses `~` or absolute |
| `home-shorten/at-home` | Home directory displays as `"~"` when cwd is elsewhere |
| `home-shorten/under-home` | Path under home (not cwd) displays with `~` prefix |
| `edge-inputs/empty-path` | Empty input returned unchanged |
| `edge-inputs/outside-home` | Path outside home stays absolute |

## How to Run

```sh
doctest vet ./libdoc/pathfmt/tests/short-display-paths/...
doctest test ./libdoc/pathfmt/tests/short-display-paths/...
```