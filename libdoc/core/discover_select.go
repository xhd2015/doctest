package core

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/xhd2015/doctest/libdoc/rules"
)

// DiscoverMode controls how much work Discover does before label filtering.
type DiscoverMode int

const (
	// DiscoverFull parses every leaf SETUP chain and ASSERT body (vet / build).
	DiscoverFull DiscoverMode = iota
	// DiscoverSelectThenDeep does a light scan (paths + labels), then the
	// caller filters, then HydrateTreeCases deep-parses only the run set.
	DiscoverSelectThenDeep
)

// DiscoverTreeCasesLight finds all ASSERT.md leaves and reads ASSERT frontmatter
// only (labels/explanation). It does not parse SETUP Go blocks or validate
// intermediate SETUP directories. Root DOCTEST.md is still required and checked
// for Request/Response/Run so selected leaves can hydrate later.
//
// Used by doctest test default discovery so labeled (e.g. heavy) leaves are not
// fully parsed when they will be skipped.
func DiscoverTreeCasesLight(root string) ([]TreeCase, error) {
	info, err := os.Stat(root)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("not a directory: %s", root)
	}

	var verrs []ValidationError
	doctestPath := filepath.Join(root, "DOCTEST.md")
	doctestContent, err := os.ReadFile(doctestPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	doctestDoc, err := ParseDOCTESTDocument(doctestPath, string(doctestContent))
	if err != nil {
		verrs = append(verrs, ValidationError{Path: "DOCTEST.md", Msg: err.Error()})
	} else if doctestDoc.GoBlock == nil {
		verrs = append(verrs, ValidationError{Path: "DOCTEST.md", Msg: "must have a Go code block"})
	} else {
		if v := rules.CheckRootHasRequestResponse(doctestDoc.GoBlock.Types, "DOCTEST.md"); v != nil {
			verrs = append(verrs, ValidationError{Path: v.Path, Msg: v.Msg})
		}
		if v := rules.CheckRootHasRun(doctestDoc.GoBlock.Run != nil, "DOCTEST.md"); v != nil {
			verrs = append(verrs, ValidationError{Path: v.Path, Msg: v.Msg})
		}
	}
	if len(verrs) > 0 {
		return nil, JoinValidationErrors(verrs)
	}

	var cases []TreeCase
	err = filepath.WalkDir(root, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			if d.Name() == "testdata" {
				return filepath.SkipDir
			}
			if path == root {
				return nil
			}
			if _, err := os.Stat(filepath.Join(path, "DOCTEST.md")); err == nil {
				return filepath.SkipDir
			}
			return nil
		}
		if d.Name() != "ASSERT.md" {
			return nil
		}
		leafDir := filepath.Dir(path)
		relLeaf, err := filepath.Rel(root, leafDir)
		if err != nil {
			return err
		}
		if relLeaf == "." {
			relLeaf = ""
		}
		content, err := os.ReadFile(path)
		if err != nil {
			relAssert, _ := filepath.Rel(root, path)
			verrs = append(verrs, ValidationError{Path: relAssert, Msg: err.Error()})
			return nil
		}
		relAssert, _ := filepath.Rel(root, path)
		fm, _, fmErr := ParseAssertFrontmatter(relAssert, string(content))
		if fmErr != nil {
			verrs = append(verrs, ValidationError{Path: relAssert, Msg: fmErr.Error()})
			return nil
		}
		cases = append(cases, TreeCase{
			Name:        CaseName(relLeaf),
			Path:        relLeaf,
			Labels:      append([]string(nil), fm.Labels...),
			Explanation: fm.Explanation,
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	if len(verrs) > 0 {
		return nil, JoinValidationErrors(verrs)
	}
	sort.Slice(cases, func(i, j int) bool { return cases[i].Path < cases[j].Path })
	return cases, nil
}

// HydrateTreeCases deep-parses SETUP chains and ASSERT Go bodies for light cases
// (or any cases missing SetupFiles). Root DOCTEST is re-read once; SETUP parses
// are memoized for the batch.
func HydrateTreeCases(root string, cases []TreeCase) ([]TreeCase, error) {
	if len(cases) == 0 {
		return cases, nil
	}
	doctestPath := filepath.Join(root, "DOCTEST.md")
	doctestContent, err := os.ReadFile(doctestPath)
	if err != nil {
		return nil, err
	}
	doctestDoc, err := ParseDOCTESTDocument(doctestPath, string(doctestContent))
	if err != nil {
		return nil, fmt.Errorf("DOCTEST.md: %w", err)
	}
	setupCache := make(map[string]setupCacheEntry)
	var verrs []ValidationError
	out := make([]TreeCase, 0, len(cases))
	for _, tc := range cases {
		leafDir := root
		if tc.Path != "" {
			leafDir = filepath.Join(root, filepath.FromSlash(tc.Path))
		}
		setupDocs, chainErr := setupChainCached(root, leafDir, doctestDoc, setupCache)
		if chainErr != nil {
			relAssert := tc.Path
			if relAssert != "" {
				relAssert = filepath.ToSlash(filepath.Join(relAssert, "ASSERT.md"))
			} else {
				relAssert = "ASSERT.md"
			}
			verrs = append(verrs, ValidationError{Path: relAssert, Msg: chainErr.Error()})
			continue
		}
		// Mirror full-discover WalkDir rules for SETUP on the leaf's ancestor
		// path: root SETUP.md is optional and need not define Setup; missing
		// intermediate SETUP.md is OK; when a non-root SETUP.md exists it must
		// have a Go block with func Setup.
		for _, doc := range setupDocs {
			if doc.Path == "DOCTEST.md" || doc.Path == "SETUP.md" {
				continue
			}
			if !strings.HasSuffix(doc.Path, "SETUP.md") {
				continue
			}
			setupAbs := filepath.Join(root, filepath.FromSlash(doc.Path))
			if _, statErr := os.Stat(setupAbs); os.IsNotExist(statErr) {
				continue
			}
			if doc.GoBlock == nil {
				verrs = append(verrs, ValidationError{Path: doc.Path, Msg: "must have a Go code block"})
			} else if doc.GoBlock.Setup == nil {
				verrs = append(verrs, ValidationError{Path: doc.Path, Msg: "must have func Setup"})
			}
		}
		assertPath := filepath.Join(leafDir, "ASSERT.md")
		assertContent, err := os.ReadFile(assertPath)
		if err != nil {
			relAssert, _ := filepath.Rel(root, assertPath)
			verrs = append(verrs, ValidationError{Path: relAssert, Msg: err.Error()})
			continue
		}
		relAssert, _ := filepath.Rel(root, assertPath)
		assertDoc, err := ParseAssertDocument(relAssert, string(assertContent))
		if err != nil {
			verrs = append(verrs, ValidationError{Path: relAssert, Msg: err.Error()})
			continue
		}
		if v := rules.CheckChainHasRun(runSource(setupDocs), relAssert); v != nil {
			verrs = append(verrs, ValidationError{Path: v.Path, Msg: v.Msg})
		}
		tc.SetupFiles = setupDocs
		tc.AssertFile = assertDoc
		tc.Labels = append([]string(nil), assertDoc.Frontmatter.Labels...)
		tc.Explanation = assertDoc.Frontmatter.Explanation
		out = append(out, tc)
	}
	if len(verrs) > 0 {
		return nil, JoinValidationErrors(verrs)
	}
	return out, nil
}
