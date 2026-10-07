package work

import (
	"bytes"
	"chrisfoong/chaum-work-management-backend/internal/auth"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Storage struct {
	Base, Key, Bucket string
	Client            *http.Client
}

func NewStorage(base, key, bucket string) *Storage {
	return &Storage{strings.TrimRight(base, "/"), key, bucket, &http.Client{Timeout: 15 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}}
}
func (st *Storage) request(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	if st.Key == "" || st.Base == "" || st.Bucket == "" {
		return nil, conflict("storage is not configured")
	}
	r, e := http.NewRequestWithContext(ctx, method, st.Base+"/storage/v1"+path, body)
	if e != nil {
		return nil, e
	}
	r.Header.Set("Authorization", "Bearer "+st.Key)
	r.Header.Set("apikey", st.Key)
	return r, nil
}
func owned(p auth.Principal, path string) bool {
	parts := strings.Split(path, "/")
	if len(parts) != 2 || parts[0] != p.UserID.String() {
		return false
	}
	_, e := uuid.Parse(parts[1])
	return e == nil
}
func mediaAllowed(p auth.Principal, mime string) bool {
	return mime == "image/jpeg" || mime == "image/png" || (p.Role != auth.RoleWorker && mime == "application/pdf")
}
func (st *Storage) Verify(ctx context.Context, p auth.Principal, path string) error {
	if !owned(p, path) {
		return invalid("photo_path", "must be an object uploaded by the actor")
	}
	r, e := st.request(ctx, http.MethodGet, "/object/authenticated/"+url.PathEscape(st.Bucket)+"/"+path, nil)
	if e != nil {
		return e
	}
	resp, e := st.Client.Do(r)
	if e != nil {
		return fmt.Errorf("storage verification failed")
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return invalid("photo_path", "object not found")
	}
	body, e := io.ReadAll(io.LimitReader(resp.Body, 5*1024*1024+1))
	if e != nil {
		return fmt.Errorf("storage read failed")
	}
	if len(body) == 0 || len(body) > 5*1024*1024 || !mediaAllowed(p, http.DetectContentType(body)) {
		return invalid("photo_path", "invalid file type or size")
	}
	return nil
}
func (st *Storage) Upload(c *gin.Context) {
	p := actor(c)
	body, e := io.ReadAll(io.LimitReader(c.Request.Body, 5*1024*1024+1))
	if e != nil || len(body) == 0 || len(body) > 5*1024*1024 {
		respond(c, nil, invalid("file", "1 byte to 5MB required"))
		return
	}
	mime := http.DetectContentType(body)
	if !mediaAllowed(p, mime) {
		respond(c, nil, invalid("file", "JPEG/PNG required; staff receipts may be PDF"))
		return
	}
	path := p.UserID.String() + "/" + uuid.NewString()
	r, e := st.request(c.Request.Context(), http.MethodPost, "/object/"+url.PathEscape(st.Bucket)+"/"+path, bytes.NewReader(body))
	if e != nil {
		respond(c, nil, e)
		return
	}
	r.Header.Set("Content-Type", mime)
	r.Header.Set("x-upsert", "false")
	resp, e := st.Client.Do(r)
	if e != nil {
		respond(c, nil, fmt.Errorf("storage upload failed"))
		return
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, io.LimitReader(resp.Body, 4096))
	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		respond(c, nil, conflict("storage upload rejected"))
		return
	}
	respond(c, map[string]any{"object_path": path, "mime_type": mime, "size": len(body), "warning": "upload is separate from the business transaction; unreferenced uploads require cleanup"}, nil)
}

// Read checks ownership or a staff permission to an existing evidence/receipt row.
func (st *Storage) Read(s *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		p := actor(c)
		path := c.Query("path")
		parts := strings.Split(path, "/")
		if len(parts) != 2 {
			respond(c, nil, invalid("path", "invalid object path"))
			return
		}
		if _, e := uuid.Parse(parts[0]); e != nil {
			respond(c, nil, invalid("path", "invalid object path"))
			return
		}
		if _, e := uuid.Parse(parts[1]); e != nil {
			respond(c, nil, invalid("path", "invalid object path"))
			return
		}
		if !owned(p, path) {
			var allowed bool
			query := `SELECT EXISTS(SELECT 1 FROM work_evidence WHERE photo_url=$1)`
			if p.Role == auth.RoleSupervisor {
				query = `SELECT EXISTS(SELECT 1 FROM work_evidence WHERE photo_url=$1 UNION ALL SELECT 1 FROM expense_claim WHERE receipt_photo_url=$1)`
			}
			if p.Role == auth.RoleWorker {
				respond(c, nil, invalid("path", "not your object"))
				return
			}
			if e := s.Repo.Pool.QueryRow(c.Request.Context(), query, path).Scan(&allowed); e != nil || !allowed {
				respond(c, nil, conflict("file not accessible"))
				return
			}
		}
		r, e := st.request(c.Request.Context(), http.MethodGet, "/object/authenticated/"+url.PathEscape(st.Bucket)+"/"+path, nil)
		if e != nil {
			respond(c, nil, e)
			return
		}
		resp, e := st.Client.Do(r)
		if e != nil {
			respond(c, nil, fmt.Errorf("storage read failed"))
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			respond(c, nil, conflict("file not found"))
			return
		}
		body, e := io.ReadAll(io.LimitReader(resp.Body, 5*1024*1024+1))
		if e != nil || len(body) > 5*1024*1024 {
			respond(c, nil, conflict("file unavailable"))
			return
		}
		c.Header("Cache-Control", "no-store")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Data(200, http.DetectContentType(body), body)
	}
}

type LinePush struct {
	Token  string
	Client *http.Client
}

func (l *LinePush) Send(ctx context.Context, user, msg string) error {
	if l.Token == "" {
		return conflict("LINE messaging is not configured")
	}
	body, _ := json.Marshal(map[string]any{"to": user, "messages": []map[string]string{{"type": "text", "text": msg}}})
	r, _ := http.NewRequestWithContext(ctx, "POST", "https://api.line.me/v2/bot/message/push", bytes.NewReader(body))
	r.Header.Set("Authorization", "Bearer "+l.Token)
	r.Header.Set("Content-Type", "application/json")
	resp, e := l.Client.Do(r)
	if e != nil {
		return fmt.Errorf("LINE push failed")
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("LINE push rejected (%d)", resp.StatusCode)
	}
	return nil
}
