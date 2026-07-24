# Scenario

**Feature**: manifest hash miss writes new bytes and updates the entry

```
write sampleGoA -> write sampleGoB
  -> target content is B
  -> manifest entry for rel changes
```

## Steps

1. Inherit write-if-changed Setup.
2. First write sampleGoA; snapshot manifest entry.
3. Run second write with SecondFileContent = sampleGoB.

```go
import (
	"testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Mode = "write-file-second"
	req.SecondFileContent = sampleGoB
	if err := writeGenRelFile(t, req.GenDir, req.RelPath, req.FileContent); err != nil {
		t.Fatalf("first writeGenRelFile: %v", err)
	}
	man := readFileOrEmpty(manifestPath(req.GenDir))
	req.SnapManifestEntryBefore = findManifestLine(man, req.RelPath)
	req.SnapManifestContentBefore = man
	return nil
}
```
