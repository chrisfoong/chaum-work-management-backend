package work

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"regexp"
	"strings"
	"time"

	"chrisfoong/chaum-work-management-backend/internal/apperr"
	"chrisfoong/chaum-work-management-backend/internal/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const invitationPrefix = "INV:"
const cancelledInvitationPrefix = "CNL:"
const invitationLifetime = 24 * time.Hour

// The existing line_id holds a 47-character token digest until LINE is bound.
// The raw invitation is never stored, logged, or used as an authentication token.
func newInvitation(now time.Time) (string, string, time.Time, error) {
	payload := make([]byte, 40)
	expires := now.Add(invitationLifetime).Truncate(time.Second)
	binary.BigEndian.PutUint64(payload[:8], uint64(expires.Unix()))
	if _, err := rand.Read(payload[8:]); err != nil {
		return "", "", time.Time{}, err
	}
	token := base64.RawURLEncoding.EncodeToString(payload)
	return token, invitationDigest(token), expires, nil
}
func invitationDigest(token string) string {
	digest := sha256.Sum256([]byte(token))
	return invitationPrefix + base64.RawURLEncoding.EncodeToString(digest[:])
}
func invitationExpiry(token string, now time.Time) (time.Time, error) {
	payload, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil || len(payload) != 40 || base64.RawURLEncoding.EncodeToString(payload) != token {
		return time.Time{}, invalid("invitation", "ลิงก์เชิญไม่ถูกต้อง")
	}
	expires := time.Unix(int64(binary.BigEndian.Uint64(payload[:8])), 0)
	if !expires.After(now) || expires.After(now.Add(invitationLifetime+time.Second)) {
		return time.Time{}, conflict("ลิงก์เชิญหมดอายุ กรุณาขอลิงก์ใหม่จากผู้ควบคุมงาน")
	}
	return expires, nil
}

type InvitationUserInput struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Role        string `json:"role"`
	Phone       string `json:"phone_number"`
	BankName    string `json:"bank_name"`
	BankAccount string `json:"bank_account_no"`
	DailyWage   string `json:"daily_wage"`
}
type InvitationTokenInput struct {
	Token string `json:"token"`
}

func invitationPlatform(role string) string {
	if role == "worker" {
		return "worker"
	}
	return "web"
}

func (s *Service) CreateInvitation(ctx context.Context, p auth.Principal, in InvitationUserInput) (any, error) {
	if p.Role != auth.RoleSupervisor || p.UserID == uuid.Nil {
		return nil, apperr.Forbidden("supervisor required")
	}
	in.FirstName = strings.TrimSpace(in.FirstName)
	in.LastName = strings.TrimSpace(in.LastName)
	in.Phone = strings.TrimSpace(in.Phone)
	in.BankName = strings.TrimSpace(in.BankName)
	in.BankAccount = strings.TrimSpace(in.BankAccount)
	for _, field := range []struct {
		value, name string
		max         int
	}{{in.FirstName, "first_name", 100}, {in.LastName, "last_name", 100}, {in.BankName, "bank_name", 100}, {in.BankAccount, "bank_account_no", 50}} {
		if err := text(field.value, field.name, field.max); err != nil {
			return nil, err
		}
	}
	if !regexp.MustCompile(`^[0-9]{10}$`).MatchString(in.Phone) {
		return nil, invalid("phone_number", "10 digits required")
	}
	if in.Role != "worker" && in.Role != "assistant" && in.Role != "supervisor" {
		return nil, invalid("role", "unknown role")
	}
	wage, err := Money(in.DailyWage)
	if err != nil || wage > 9999999999 {
		return nil, invalid("daily_wage", "invalid wage")
	}
	token, digest, expires, err := newInvitation(s.Now())
	if err != nil {
		return nil, err
	}
	var id string
	err = s.Repo.Pool.QueryRow(ctx, `INSERT INTO public."USER"(first_name,last_name,role,phone_number,line_id,bank_name,bank_account_no,daily_wage,is_active)
 VALUES($1,$2,$3,$4,$5,$6,$7,$8::numeric,false) RETURNING user_id::text`, in.FirstName, in.LastName, in.Role, in.Phone, digest, in.BankName, in.BankAccount, Decimal(wage)).Scan(&id)
	if err != nil {
		return nil, err
	}
	return gin.H{"user_id": id, "token": token, "expires_at": expires, "platform": invitationPlatform(in.Role)}, nil
}

func (s *Service) RenewInvitation(ctx context.Context, p auth.Principal, id string) (any, error) {
	if p.Role != auth.RoleSupervisor || p.UserID == uuid.Nil {
		return nil, apperr.Forbidden("supervisor required")
	}
	if err := validID(id, "user_id"); err != nil {
		return nil, err
	}
	token, digest, expires, err := newInvitation(s.Now())
	if err != nil {
		return nil, err
	}
	var role string
	err = s.Repo.Transaction(ctx, func(q Query) error {
		var current string
		var active bool
		err := q.QueryRow(ctx, `SELECT line_id,role::text,is_active FROM public."USER" WHERE user_id=$1 FOR UPDATE`, id).Scan(&current, &role, &active)
		if errors.Is(err, pgx.ErrNoRows) {
			return conflict("ไม่พบบัญชี")
		}
		if err != nil {
			return err
		}
		if active || (!strings.HasPrefix(current, invitationPrefix) && !strings.HasPrefix(current, cancelledInvitationPrefix)) {
			return conflict("บัญชีนี้ผูก LINE แล้วหรือถูกระงับ จึงออกลิงก์เชิญไม่ได้")
		}
		_, err = q.Exec(ctx, `UPDATE public."USER" SET line_id=$2,updated_at=now() WHERE user_id=$1`, id, digest)
		return err
	})
	if err != nil {
		return nil, err
	}
	return gin.H{"user_id": id, "token": token, "expires_at": expires, "platform": invitationPlatform(role)}, nil
}

func (s *Service) OnboardingState(ctx context.Context, subject string) (any, error) {
	var active bool
	var role string
	err := s.Repo.Pool.QueryRow(ctx, `SELECT is_active,role::text FROM public."USER" WHERE line_id=$1`, subject).Scan(&active, &role)
	if errors.Is(err, pgx.ErrNoRows) {
		return gin.H{"status": "unregistered"}, nil
	}
	if err != nil {
		return nil, err
	}
	if !active {
		return gin.H{"status": "inactive"}, nil
	}
	return gin.H{"status": "registered", "platform": invitationPlatform(role)}, nil
}

func (s *Service) PreviewInvitation(ctx context.Context, platform, token string) (any, error) {
	expires, err := invitationExpiry(token, s.Now())
	if err != nil {
		return nil, err
	}
	var first, last, role string
	err = s.Repo.Pool.QueryRow(ctx, `SELECT first_name,last_name,role::text FROM public."USER" WHERE line_id=$1 AND NOT is_active`, invitationDigest(token)).Scan(&first, &last, &role)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, conflict("ลิงก์เชิญถูกใช้แล้วหรือถูกเปลี่ยน กรุณาติดต่อผู้ควบคุมงาน")
	}
	if err != nil {
		return nil, err
	}
	if invitationPlatform(role) != platform {
		return nil, apperr.Forbidden("ลิงก์เชิญนี้ต้องเปิดผ่านช่องทางที่กำหนด")
	}
	return gin.H{"first_name": first, "last_name": last, "role": role, "expires_at": expires}, nil
}

func (s *Service) ClaimInvitation(ctx context.Context, subject, platform, token string) (any, error) {
	if !regexp.MustCompile(`^U[0-9a-f]{32}$`).MatchString(subject) {
		return nil, apperr.Unauthorized("invalid LINE identity")
	}
	if platform != "web" && platform != "worker" {
		return nil, apperr.Forbidden("invalid platform")
	}
	var userID, role string
	err := s.Repo.Transaction(ctx, func(q Query) error {
		if _, err := invitationExpiry(token, s.Now()); err != nil {
			return err
		}
		if err := lock(ctx, q, "line-invitation:"+subject); err != nil {
			return err
		}
		var exists bool
		if err := q.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM public."USER" WHERE line_id=$1)`, subject).Scan(&exists); err != nil {
			return err
		}
		if exists {
			return conflict("บัญชี LINE นี้ผูกกับผู้ใช้ในระบบแล้ว ไม่สามารถผูกซ้ำได้")
		}
		err := q.QueryRow(ctx, `SELECT user_id::text,role::text FROM public."USER" WHERE line_id=$1 AND NOT is_active FOR UPDATE`, invitationDigest(token)).Scan(&userID, &role)
		if errors.Is(err, pgx.ErrNoRows) {
			return conflict("ลิงก์เชิญถูกใช้แล้วหรือถูกเปลี่ยน")
		}
		if err != nil {
			return err
		}
		if invitationPlatform(role) != platform {
			return apperr.Forbidden("บทบาทไม่ตรงกับช่องทางเข้าสู่ระบบ")
		}
		if _, err := invitationExpiry(token, s.Now()); err != nil {
			return err
		}
		if _, err = q.Exec(ctx, `UPDATE public."USER" SET line_id=$2,is_active=true,updated_at=now() WHERE user_id=$1`, userID, subject); err != nil {
			return err
		}
		if role == "worker" {
			if _, err = q.Exec(ctx, query42, userID); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return gin.H{"user_id": userID, "role": role, "status": "registered"}, nil
}

func RegisterOnboarding(engine *gin.Engine, webVerifier, workerVerifier auth.Verifier, s *Service) {
	for _, entry := range []struct {
		prefix, platform string
		verifier         auth.Verifier
	}{{"/api/web", "web", webVerifier}, {"/api/liff", "worker", workerVerifier}} {
		platform := entry.platform
		group := engine.Group(entry.prefix, auth.VerifyIdentity(entry.verifier))
		group.GET("/onboarding", func(c *gin.Context) {
			subject, _ := auth.LineSubject(c)
			out, err := s.OnboardingState(c.Request.Context(), subject)
			respond(c, out, err)
		})
		group.POST("/invitation/preview", input(func(c *gin.Context, _ auth.Principal, in InvitationTokenInput) (any, error) {
			return s.PreviewInvitation(c.Request.Context(), platform, in.Token)
		}))
		group.POST("/invitation/claim", input(func(c *gin.Context, _ auth.Principal, in InvitationTokenInput) (any, error) {
			subject, _ := auth.LineSubject(c)
			return s.ClaimInvitation(c.Request.Context(), subject, platform, in.Token)
		}))
	}
}

// CancelInvitation invalidates the current invitation without deleting the person.
func (s *Service) CancelInvitation(ctx context.Context, p auth.Principal, id string) (any, error) {
	if p.Role != auth.RoleSupervisor || p.UserID == uuid.Nil {
		return nil, apperr.Forbidden("supervisor required")
	}
	if err := validID(id, "user_id"); err != nil {
		return nil, err
	}
	err := s.Repo.Transaction(ctx, func(q Query) error {
		var current string
		var active bool
		err := q.QueryRow(ctx, `SELECT line_id,is_active FROM public."USER" WHERE user_id=$1 FOR UPDATE`, id).Scan(&current, &active)
		if errors.Is(err, pgx.ErrNoRows) {
			return conflict("ไม่พบบัญชี")
		}
		if err != nil {
			return err
		}
		if active || (!strings.HasPrefix(current, invitationPrefix) && !strings.HasPrefix(current, cancelledInvitationPrefix)) {
			return conflict("บัญชีผูก LINE แล้ว ไม่สามารถยกเลิกด้วยลิงก์เชิญได้")
		}
		if strings.HasPrefix(current, cancelledInvitationPrefix) {
			return nil
		}
		marker := cancelledInvitationPrefix + strings.TrimPrefix(current, invitationPrefix)
		_, err = q.Exec(ctx, `UPDATE public."USER" SET line_id=$2,updated_at=now() WHERE user_id=$1`, id, marker)
		return err
	})
	if err != nil {
		return nil, err
	}
	return gin.H{"status": "cancelled", "user_id": id}, nil
}
