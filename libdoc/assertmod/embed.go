package assertmod

import (
	"crypto/md5"
	_ "embed"
	"encoding/hex"
)

//go:generate go run ../../script/generate

//go:embed assert.go
var content []byte

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