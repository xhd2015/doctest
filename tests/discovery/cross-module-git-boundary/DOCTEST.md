# Cross-Module Git-Aware Discovery Tests

Integration tests for `doctest test ./...` when a nested `go.mod` boundary is
crossed. At non-child module paths, discovery depends on whether ancestor and
nested module roots share the same git work tree.

## Git Context Notes

`gitRoot(dir)` resolves via `git.ShowToplevel(dir)` when inside a work tree,
otherwise `""`.

**Collapsed case:** A child directory nested inside the parent's git worktree
without its own `.git` is **same-git context** — `gitRoot(child) ==
gitRoot(parent)`. This is identical to `parent-git-child-same-git-discovers`;
there is no separate warn path for "parent in git, child has no own repo" when
the child lives inside the parent worktree.

**Distinct mismatch cases** require unequal git roots: separate repos
(`parent-git-child-other-git-warns`) or one side in git and the other outside
any work tree (`parent-not-git-child-git-warns`).

## Decision Tree

```
cross-module-git-boundary/
├── module-path/
│   ├── child-path/
│   │   └── child-module-unchanged/          → child path, same git → discover, no warning
│   └── non-child-path/
│       └── git-context/
│           ├── both-null/
│           │   └── both-not-git-discovers/  → neither in git → discover, no warning
│           ├── same-git/
│           │   ├── parent-git-child-same-git-discovers/  → single repo → discover
│           │   └── lifelog-mirror-discovers/             → lifelog layout repro
│           └── different-git/
│               ├── parent-git-child-other-git-warns/     → separate repos → warn + skip
│               └── parent-not-git-child-git-warns/       → parent not git, child git → warn + skip
```

## Test Index

| Leaf | Description |
|------|-------------|
| `module-path/child-path/child-module-unchanged` | Child module path in single git repo: discover nested tests, no warning |
| `module-path/non-child-path/git-context/both-null/both-not-git-discovers` | Neither side in git: discover sibling module tests, no warning |
| `module-path/non-child-path/git-context/same-git/parent-git-child-same-git-discovers` | Single git repo, non-child paths: discover nested module (lifelog case) |
| `module-path/non-child-path/git-context/same-git/lifelog-mirror-discovers` | Lifelog module layout mirror in single git repo: discover `skill-cli` tests |
| `module-path/non-child-path/git-context/different-git/parent-git-child-other-git-warns` | Separate git repos: warn with `different git repository`, skip child |
| `module-path/non-child-path/git-context/different-git/parent-not-git-child-git-warns` | Parent not in git, child in git: warn with `git repository mismatch`, skip child |

## How to Run

```sh
doctest test ./tests/discovery/cross-module-git-boundary/...
doctest test ./tests/test/dotdotdot/cross-go-mod/...
doctest test ./tests/test/dotdotdot/git-boundary/...
```