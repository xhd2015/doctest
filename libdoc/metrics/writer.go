package metrics

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// SchemaVersion is the metrics event schema version written into events.
const SchemaVersion = 1

// FlushThreshold is the buffered writer size that forces a mid-run flush (128 KiB).
const FlushThreshold = 128 * 1024

// Writer is a mutex-protected buffered JSONL writer.
// Write appends one JSON object plus '\n'. Data is flushed when the buffer
// reaches FlushThreshold or when Close is called.
type Writer struct {
	mu     sync.Mutex
	file   *os.File
	buf    []byte
	closed bool
}

// OpenWriter opens path for append/write. The file must already exist
// (typically created by CreateRunFile) or be creatable.
func OpenWriter(path string) (*Writer, error) {
	// Prefer open existing; create if missing so tests that fall back to RunFilePath work.
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0o644)
	if err != nil {
		return nil, err
	}
	return &Writer{
		file: f,
		buf:  make([]byte, 0, 4096),
	}, nil
}

// Write marshals v as a single JSON object, appends '\n', and buffers the line.
// Flushes when the buffer size is >= FlushThreshold.
func (w *Writer) Write(v any) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.closed {
		return fmt.Errorf("metrics: write on closed writer")
	}
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	w.buf = append(w.buf, data...)
	w.buf = append(w.buf, '\n')
	if len(w.buf) >= FlushThreshold {
		return w.flushLocked()
	}
	return nil
}

// Close flushes any remaining buffer and closes the underlying file.
func (w *Writer) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.closed {
		return nil
	}
	w.closed = true
	if err := w.flushLocked(); err != nil {
		_ = w.file.Close()
		return err
	}
	return w.file.Close()
}

func (w *Writer) flushLocked() error {
	if len(w.buf) == 0 {
		return nil
	}
	_, err := w.file.Write(w.buf)
	if err != nil {
		return err
	}
	w.buf = w.buf[:0]
	return nil
}
