// Package storage defines the file upload port (Supabase Storage). The real
// implementation is added with the first upload slice (2A receipt photo).
package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"sync"
)

// Storage stores uploaded files. Callers validate type and size first and
// save only the returned path in the database.
type Storage interface {
	Upload(ctx context.Context, path, contentType string, body io.Reader) (storedPath string, err error)
}

// Fake is an in-memory Storage for tests.
type Fake struct {
	mu      sync.Mutex
	objects map[string][]byte
}

// NewFake returns an empty Fake.
func NewFake() *Fake {
	return &Fake{objects: map[string][]byte{}}
}

// Upload stores body under path.
func (f *Fake) Upload(_ context.Context, path, _ string, body io.Reader) (string, error) {
	data, err := io.ReadAll(body)
	if err != nil {
		return "", fmt.Errorf("read upload: %w", err)
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.objects[path] = data
	return path, nil
}

// Object returns the stored bytes for path.
func (f *Fake) Object(path string) ([]byte, bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	data, ok := f.objects[path]
	return bytes.Clone(data), ok
}
