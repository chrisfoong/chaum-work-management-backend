package storage

import (
	"context"
	"strings"
	"testing"
)

func TestFakeUpload(t *testing.T) {
	f := NewFake()
	path, err := f.Upload(context.Background(), "receipts/r1.jpg", "image/jpeg", strings.NewReader("img"))
	if err != nil {
		t.Fatal(err)
	}
	got, ok := f.Object(path)
	if !ok || string(got) != "img" {
		t.Fatalf("object = %q, %v", got, ok)
	}
}
