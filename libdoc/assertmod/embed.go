package assertmod

import (
	"crypto/md5"
	"embed"
	"encoding/hex"
	"io/fs"
	"path"
	"sort"
)

//go:generate go run ../../script/generate

//go:embed assert.go
var content []byte

//go:embed legacy_v1
var legacyV1FS embed.FS

//go:embed legacy_v2
var legacyV2FS embed.FS

func Content() []byte {
	return content
}

func ContentMD5() string {
	sum := md5.Sum(content)
	return hex.EncodeToString(sum[:])
}

func RawSourceCacheKeyMD5() string {
	return rawSourceCacheKeyMD5
}

// LegacyV1Filenames returns embedded legacy_v1 source filenames in stable order.
func LegacyV1Filenames() ([]string, error) {
	return nestedFilenames(legacyV1FS, "legacy_v1")
}

// LegacyV1File returns embedded legacy_v1 source bytes for the given filename.
func LegacyV1File(name string) ([]byte, error) {
	return fs.ReadFile(legacyV1FS, path.Join("legacy_v1", name))
}

// LegacyV2Filenames returns embedded legacy_v2 source filenames in stable order.
func LegacyV2Filenames() ([]string, error) {
	return nestedFilenames(legacyV2FS, "legacy_v2")
}

// LegacyV2File returns embedded legacy_v2 source bytes for the given filename.
func LegacyV2File(name string) ([]byte, error) {
	return fs.ReadFile(legacyV2FS, path.Join("legacy_v2", name))
}

func nestedFilenames(efs embed.FS, dir string) ([]string, error) {
	entries, err := efs.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		names = append(names, e.Name())
	}
	sort.Strings(names)
	return names, nil
}