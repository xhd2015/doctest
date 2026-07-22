# RunSuite Leaf-Cache Extract — Multi-Prep Identity + Prepare/Record

## Version
0.0.2

**Layer L2 in-process** — nested library tree for extractable multi-prep
leaf-cache helpers (identity, prepare skip list, record passes). **GREEN** —
APIs sealed under `libdoc/leafcache`. No product binary; no `label: heavy`.

Does **not** end-to-end wire workspace `./...` leaf-cache (that **L3** product
path is under parent `workspace/**`). Parent tree owns L2 key/store/partial and
L3 runtime/workspace/cli-plan.

# DSN (Domain Specific Notion)

### Participants

- **RunSuite extract seam** — shared library helpers that single-tree
  `prepareLeafCache` / `recordLeafCachePasses` and future multi-tree RunSuite
  both call after extraction.
- **Leaf identity** — tree-qualified token unique per leaf across trees; keys
  skip lists and fail maps. Not the same as the store hex key.
- **Store key** — `ComputeLeafKey` hex (already mixes abs TreeRoot).
- **Pass plan** — in-memory `keys map[identity]storeKey` plus `skip []identity`
  for GetPass hits when skip is enabled.
- **Record path** — PutPass for non-failed identities after suite JSON accounting.

### Behaviors

1. **FormatLeafIdentity(treeRoot, leafRel)** — stable non-empty token; different
   abs TreeRoots with the same relative path produce different identities.
2. **PreparePassPlan(store, leaves, goVersion, skipEnabled)** — for each LeafRef,
   compute store key under identity; when skipEnabled and GetPass, append identity
   to Skip (sorted). When skipEnabled is false, Skip is empty but Keys still fill.
3. **RecordPasses(store, keys, failed, allPassed)** — PutPass every key when
   allPassed; otherwise PutPass only identities not marked failed; when
   `!allPassed && failed == nil`, store nothing.
4. **Partial suite fail** — failed leaf never PutPass; sibling passes still store
   (product rule).

### Pipeline sketch

```
twins treeA + treeB (same rel "leaf")
  -> FormatLeafIdentity per tree
  -> optional PutPass(storeKey) to warm
  -> PreparePassPlan -> Keys + Skip (tree-qualified)
  -> after suite JSON: RecordPasses(keys, failed, allPassed)
```

## Decision Tree

```
runsuite/
├── identity/
│   ├── same-relpath-two-trees/    FormatLeafIdentity distinct across trees
│   └── stable-roundtrip/          stable + distinguishes leaf rels
└── multi-prep/
    ├── prepare-both-warm-skip/    both warm → 2 tree-qualified skip ids
    ├── prepare-one-warm-only/     only A warm → skip = [idA]
    ├── prepare-skip-disabled/     skipEnabled=false → empty Skip
    └── record-partial-fail/       failed id not PutPass; other is
```

## Test Index

| Leaf | Scenario | Expect |
|------|----------|--------|
| `identity/same-relpath-two-trees` | id(A,leaf) ≠ id(B,leaf) | **GREEN** |
| `identity/stable-roundtrip` | same inputs equal; different rels differ | **GREEN** |
| `multi-prep/prepare-both-warm-skip` | both warm → len(Skip)==2, both ids | **GREEN** |
| `multi-prep/prepare-one-warm-only` | only A warm → Skip=[idA] | **GREEN** |
| `multi-prep/prepare-skip-disabled` | warm + skipEnabled=false → empty Skip | **GREEN** |
| `multi-prep/record-partial-fail` | failed A miss; B GetPass true | **GREEN** |

## How to Run

```sh
doctest vet ./tests/leaf-cache/runsuite/
# L2 — always discovered (no heavy labels)
doctest test ./tests/leaf-cache/runsuite/ -count=1
# Parent L2 mass:  doctest test ./tests/leaf-cache/...
# Parent L3 e2e:   doctest test --label heavy ./tests/leaf-cache/...
```

## Expected public API (sealed)

Package: `github.com/xhd2015/doctest/libdoc/leafcache`

```text
func FormatLeafIdentity(treeRoot, leafRel string) string

type LeafRef struct {
    TreeRoot string
    LeafRel  string
}

type PassPlan struct {
    Keys map[string]string // identity -> store key hex
    Skip []string          // warm identities when skipEnabled
}

func PreparePassPlan(store *Store, leaves []LeafRef, goVersion string, skipEnabled bool) (PassPlan, error)

func RecordPasses(store *Store, keys map[string]string, failed map[string]bool, allPassed bool)
```

Single-tree build helpers should call these after extract (no behavior change
beyond extraction). Workspace wiring is out of scope for this tree.

```go
import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/xhd2015/doctest/libdoc/leafcache"
)

// Request drives one multi-prep surface leaf.
type Request struct {
	Op string // format_identity | format_identity_stable | multi_prep_prepare | multi_prep_record

	WorkDir     string
	ModuleRoot  string
	TreeRoot    string
	LeafDir     string
	ModuleRootB string
	TreeRootB   string
	LeafDirB    string
	GoVersion   string
	StoreRoot   string

	// PrepWarm: "both" | "a" | "none"
	PrepWarm string
	// SkipEnabled for PreparePassPlan
	SkipEnabled    bool
	SkipEnabledSet bool
	// FailedSide: "a" | "b" for multi_prep_record
	FailedSide string
	AllPassed  bool
}

type Response struct {
	Key        string
	Key2       string
	Identity   string
	Identity2  string
	SkipPaths  []string
	Hit        bool
	HitB       bool
	Err        string
}

func Run(t *testing.T, req *Request) (*Response, error) {
	t.Helper()
	resp := &Response{}
	if req.GoVersion == "" {
		req.GoVersion = "go1.25.0"
	}

	switch req.Op {
	case "format_identity":
		relA := twinLeafRel(req.TreeRoot, req.LeafDir)
		relB := twinLeafRel(req.TreeRootB, req.LeafDirB)
		idA := leafcache.FormatLeafIdentity(req.TreeRoot, relA)
		idB := leafcache.FormatLeafIdentity(req.TreeRootB, relB)
		resp.Identity, resp.Identity2 = idA, idB
		resp.Key, resp.Key2 = idA, idB
		return resp, nil

	case "format_identity_stable":
		relLeaf := twinLeafRel(req.TreeRoot, req.LeafDir)
		id1 := leafcache.FormatLeafIdentity(req.TreeRoot, relLeaf)
		id2 := leafcache.FormatLeafIdentity(req.TreeRoot, relLeaf)
		idOther := leafcache.FormatLeafIdentity(req.TreeRoot, "other")
		resp.Key, resp.Key2 = id1, id2
		resp.Identity, resp.Identity2 = id1, idOther
		return resp, nil

	case "multi_prep_prepare":
		st, err := leafcache.NewStore(req.StoreRoot)
		if err != nil {
			resp.Err = err.Error()
			return resp, err
		}
		relA := twinLeafRel(req.TreeRoot, req.LeafDir)
		relB := twinLeafRel(req.TreeRootB, req.LeafDirB)
		idA := leafcache.FormatLeafIdentity(req.TreeRoot, relA)
		idB := leafcache.FormatLeafIdentity(req.TreeRootB, relB)
		resp.Identity, resp.Identity2 = idA, idB
		leaves := []leafcache.LeafRef{
			{TreeRoot: req.TreeRoot, LeafRel: relA},
			{TreeRoot: req.TreeRootB, LeafRel: relB},
		}
		plan, err := leafcache.PreparePassPlan(st, leaves, req.GoVersion, req.SkipEnabled)
		if err != nil {
			resp.Err = err.Error()
			return resp, err
		}
		resp.SkipPaths = plan.Skip
		if plan.Keys != nil {
			resp.Key = plan.Keys[idA]
			resp.Key2 = plan.Keys[idB]
		}
		return resp, nil

	case "multi_prep_record":
		st, err := leafcache.NewStore(req.StoreRoot)
		if err != nil {
			resp.Err = err.Error()
			return resp, err
		}
		relA := twinLeafRel(req.TreeRoot, req.LeafDir)
		relB := twinLeafRel(req.TreeRootB, req.LeafDirB)
		idA := leafcache.FormatLeafIdentity(req.TreeRoot, relA)
		idB := leafcache.FormatLeafIdentity(req.TreeRootB, relB)
		resp.Identity, resp.Identity2 = idA, idB
		kA, err := leafcache.ComputeLeafKey(leafcache.KeyInput{
			ModuleRoot: req.ModuleRoot,
			TreeRoot:   req.TreeRoot,
			LeafDir:    req.LeafDir,
			GoVersion:  req.GoVersion,
		})
		if err != nil {
			resp.Err = err.Error()
			return resp, err
		}
		kB, err := leafcache.ComputeLeafKey(leafcache.KeyInput{
			ModuleRoot: req.ModuleRootB,
			TreeRoot:   req.TreeRootB,
			LeafDir:    req.LeafDirB,
			GoVersion:  req.GoVersion,
		})
		if err != nil {
			resp.Err = err.Error()
			return resp, err
		}
		resp.Key, resp.Key2 = kA, kB
		keys := map[string]string{idA: kA, idB: kB}
		var failed map[string]bool
		if !req.AllPassed {
			failed = map[string]bool{}
			switch req.FailedSide {
			case "a":
				failed[idA] = true
			case "b":
				failed[idB] = true
			}
		}
		leafcache.RecordPasses(st, keys, failed, req.AllPassed)
		hitA, err := st.GetPass(kA)
		if err != nil {
			resp.Err = err.Error()
			return resp, err
		}
		hitB, err := st.GetPass(kB)
		if err != nil {
			resp.Err = err.Error()
			return resp, err
		}
		resp.Hit, resp.HitB = hitA, hitB
		return resp, nil

	default:
		return nil, fmt.Errorf("unknown Op %q", req.Op)
	}
}

// twinLeafRel returns the tree-relative slash path of leafDir under treeRoot.
func twinLeafRel(treeRoot, leafDir string) string {
	if treeRoot == "" || leafDir == "" {
		return "leaf"
	}
	rel, err := filepath.Rel(treeRoot, leafDir)
	if err != nil || rel == "" || rel == "." {
		return "leaf"
	}
	return filepath.ToSlash(rel)
}
```
