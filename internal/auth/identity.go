package auth

import (
	"chrisfoong/chaum-work-management-backend/internal/apperr"
	"chrisfoong/chaum-work-management-backend/internal/httpx"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"net"
)

const lineSubjectKey = "auth.verified-line-subject"

func LineSubject(c *gin.Context) (string, bool) {
	value, ok := c.Get(lineSubjectKey)
	subject, valid := value.(string)
	return subject, ok && valid && subject != ""
}

// VerifyIdentity is only for onboarding: it never sets an application principal.
func VerifyIdentity(v Verifier) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, ok := bearerToken(c.GetHeader("Authorization"))
		if !ok {
			httpx.WriteError(c, apperr.Unauthorized("missing bearer token"))
			return
		}
		subject, err := v.Verify(c.Request.Context(), token)
		if err != nil {
			var network net.Error
			if errors.Is(err, ErrVerificationUnavailable) || errors.Is(err, context.DeadlineExceeded) || errors.As(err, &network) {
				httpx.WriteError(c, apperr.Unavailable("LINE verification unavailable"))
			} else {
				httpx.WriteError(c, apperr.Unauthorized("invalid or expired token"))
			}
			return
		}
		if subject == "" {
			httpx.WriteError(c, apperr.Unauthorized("invalid LINE identity"))
			return
		}
		c.Set(lineSubjectKey, subject)
		c.Next()
	}
}
