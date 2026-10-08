package work

import (
	"bytes"
	"chrisfoong/chaum-work-management-backend/internal/auth"
	"context"
	"github.com/google/uuid"
	"image"
	"image/png"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestContractPNGVerificationChecksBytesOwnerAndSize(t *testing.T) {
	var buffer bytes.Buffer
	if e := png.Encode(&buffer, image.NewRGBA(image.Rect(0, 0, 1, 1))); e != nil {
		t.Fatal(e)
	}
	body := buffer.Bytes()
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls++; w.Write(body) }))
	defer server.Close()
	st := NewStorage(server.URL, "test-key", "private")
	p := auth.Principal{UserID: uuid.New(), Role: auth.RoleSupervisor}
	path := p.UserID.String() + "/" + uuid.NewString()
	if e := st.VerifyPNG(context.Background(), p, path); e != nil {
		t.Fatal("valid PNG", e)
	}
	body = []byte("%PDF-1.7\nnot-a-contract-png")
	if e := st.VerifyPNG(context.Background(), p, path); e == nil {
		t.Fatal("PDF accepted as PNG")
	}
	body = bytes.Repeat([]byte("x"), 5*1024*1024+1)
	if e := st.VerifyPNG(context.Background(), p, path); e == nil {
		t.Fatal("oversized file accepted")
	}
	count := calls
	if e := st.VerifyPNG(context.Background(), p, uuid.NewString()+"/"+uuid.NewString()); e == nil || calls != count {
		t.Fatal("another owner's file fetched")
	}
}
