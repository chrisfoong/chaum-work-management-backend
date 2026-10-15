package schema

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

// The embedded exports must stay byte-identical to the docs/schema authority.
// After a new Supabase export, copy both files into docs/schema and here.
func TestEmbeddedExportsMatchDocs(t *testing.T) {
	for _, name := range []string{"columns.json", "constraints.csv"} {
		embedded, err := files.ReadFile(name)
		if err != nil {
			t.Fatalf("read embedded %s: %v", name, err)
		}
		authority, err := os.ReadFile(filepath.Join("..", "..", "docs", "schema", name))
		if err != nil {
			t.Fatalf("read docs/schema/%s: %v", name, err)
		}
		if !bytes.Equal(embedded, authority) {
			t.Errorf("internal/schema/%s differs from docs/schema/%s; copy the docs version over", name, name)
		}
	}
}
