# go testcache: stale `(cached)` after dependency source change?

Isolated repros (no doctest binary). Two questions:

1. Does Go ignore **observed** semantic changes in a dependency? → **No** (`./run.sh`)
2. Can Go report `(cached)` after an intermediate Setup-like edit when the write is **unread** (doctest-shaped)? → **Yes** (`./run_doctest_shape.sh`)

---

## Key factor (not cwd)

```text
go test result cache
  id1 = link action ID     ← includes mid.a content ID as packagefile
  id2 = binary content ID  ← hash of linked test binary (buildid-neutral)

Edit mid.Setup that only assigns req.WorkDir, leaf never reads it:
  mid.a rebuilds, leaf may re-inline
  optimizer drops the write from the *linked* image (DCE)
  binary content ID STABLE
  id1 MISS  +  id2 HIT  →  ok  (cached)

testlog / cwd:
  only rechecks env + open/stat/chdir from the previous *test process*
  mid.go / mapping-gen never appear there → not the root cause
```

| Factor | Role |
|--------|------|
| **cwd / chdir** | Not the driver of this false `(cached)` |
| **testlog** | Irrelevant for gen/mid sources |
| **Unread Setup writes + DCE/inline** | Keeps **binary content ID** stable |
| **id1 vs id2 dual lookup** | id1 sees mid rebuild; id2 reuses warm result |
| **Suite fingerprint** | **Removed** — doctest relies on Go testcache; unused Setup writes may stay `(cached)` |

---

## A. Observed `Version()` — Go is fine

```bash
./run.sh
```

| Variant | After `return 1`→`2` | Verdict |
|---------|----------------------|---------|
| direct / blank+Fn | re-run **FAIL** | **GO_OK** |

---

## B. Doctest-shaped unread Setup — bug shape

```bash
./run_doctest_shape.sh        # both
./run_doctest_shape.sh dce    # false (cached)
./run_doctest_shape.sh observe
```

Layout (`doctest_shape/{dce,observe}/`):

```text
droot/   Request, RootSetup, Run
mid/     Setup (WorkDir = MID_V1|V2)     ← edited
leaf/    registry + RunTestLeaf
suite/   blank-import leaf, TestDoctestSuite
```

| Mode | Leaf | After mid V1→V2 | Verdict |
|------|------|-----------------|---------|
| **dce** | never reads WorkDir | **`(cached)`**, content-id **stable** | **BUG_REPRODUCED** |
| **observe** | asserts `WorkDir == MID_V1` | re-run **FAIL** | **GO_OK_RERUN_FAIL** |

Recorded on go1.25.10 darwin/arm64:

```text
VERDICT[dce]=BUG_REPRODUCED  a1_cached=1 cid_same=1
VERDICT[observe]=GO_OK_RERUN_FAIL
```

---

## What this means for doctest

- Not “Go is broken for all dep edits.”
- Matches hierarchical SETUP that only mutates request fields never asserted: users see **false `(cached)`** after editing SETUP.md.
- Product policy: no suite fingerprint; cache follows Go. Doctests cover **meaningful** (observed → miss) vs **unused** (dead write → may stay cached).
