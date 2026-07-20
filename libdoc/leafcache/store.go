package leafcache

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Store is a disk map of leaf keys to explicit pass markers under Root.
// Product default root concept: $CacheHome/doctest/leaf-cache/v1.
type Store struct {
	Root string
}

// NewStore opens a pass store rooted at root. The directory need not exist yet;
// PutPass creates parents as needed. Roots are isolated: writes under A never
// appear under B.
func NewStore(root string) (*Store, error) {
	if strings.TrimSpace(root) == "" {
		return nil, fmt.Errorf("leafcache: empty store root")
	}
	abs, err := filepath.Abs(root)
	if err != nil {
		return nil, fmt.Errorf("leafcache: store root: %w", err)
	}
	return &Store{Root: abs}, nil
}

// GetPass reports whether key was previously recorded with PutPass.
// Missing keys return false, nil.
func (s *Store) GetPass(key string) (bool, error) {
	if s == nil {
		return false, fmt.Errorf("leafcache: nil store")
	}
	if err := validateKey(key); err != nil {
		return false, err
	}
	path := s.entryPath(key)
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// PutPass records an explicit pass for key. Failures never call this.
func (s *Store) PutPass(key string) error {
	if s == nil {
		return fmt.Errorf("leafcache: nil store")
	}
	if err := validateKey(key); err != nil {
		return err
	}
	path := s.entryPath(key)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	// Marker content is intentionally small; presence is the signal.
	return os.WriteFile(path, []byte("pass\n"), 0o644)
}

func validateKey(key string) error {
	if key == "" {
		return fmt.Errorf("leafcache: empty key")
	}
	if strings.Contains(key, string(filepath.Separator)) || strings.Contains(key, "/") || strings.Contains(key, `\`) {
		return fmt.Errorf("leafcache: key must not contain path separators")
	}
	if strings.Contains(key, "..") {
		return fmt.Errorf("leafcache: key must not contain ..")
	}
	return nil
}

// entryPath shards keys under Root as <aa>/<bb>/<key> when key is long enough.
func (s *Store) entryPath(key string) string {
	if len(key) >= 4 {
		return filepath.Join(s.Root, key[:2], key[2:4], key)
	}
	return filepath.Join(s.Root, key)
}
