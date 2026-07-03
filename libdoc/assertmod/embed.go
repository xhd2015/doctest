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
	entries, err := legacyV1FS.ReadDir("legacy_v1")
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

// LegacyV1File returns embedded legacy_v1 source bytes for the given filename.
func LegacyV1File(name string) ([]byte, error) {
	return fs.ReadFile(legacyV1FS, path.Join("legacy_v1", name))
}