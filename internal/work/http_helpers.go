package work

import (
	"chrisfoong/chaum-work-management-backend/internal/apperr"
	"chrisfoong/chaum-work-management-backend/internal/auth"
	"chrisfoong/chaum-work-management-backend/internal/httpx"
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"net/http"
	"strconv"
)

type Handler struct{ Service *Service }

func respond(c *gin.Context, out any, e error) {
	if e != nil {
		var pg *pgconn.PgError
		if errors.As(e, &pg) {
			switch pg.Code {
			case "23505":
				e = conflict("duplicate record")
			case "23503":
				e = invalid("reference", "referenced record does not exist")
			case "22P02", "22003", "23514", "23502":
				e = invalid("body", "invalid value for database schema")
			}
		}
		if isMissing(e) {
			e = apperr.NotFound("record not found")
		}
		httpx.WriteError(c, e)
		return
	}
	if out == nil {
		out = gin.H{"success": true}
	}
	c.JSON(http.StatusOK, out)
}
func actor(c *gin.Context) auth.Principal { p, _ := auth.PrincipalFrom(c); return p }
func input[T any](fn func(*gin.Context, auth.Principal, T) (any, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		var in T
		if e := httpx.BindJSON(c, &in); e != nil {
			respond(c, nil, e)
			return
		}
		out, e := fn(c, actor(c), in)
		respond(c, out, e)
	}
}
func paging(c *gin.Context) (int, int, error) {
	limit := 50
	offset := 0
	var e error
	if v := c.Query("limit"); v != "" {
		limit, e = strconv.Atoi(v)
		if e != nil {
			return 0, 0, invalid("limit", "integer required")
		}
	}
	if v := c.Query("offset"); v != "" {
		offset, e = strconv.Atoi(v)
		if e != nil {
			return 0, 0, invalid("offset", "integer required")
		}
	}
	if limit < 1 || limit > 100 || offset < 0 || offset > 100000 {
		return 0, 0, invalid("pagination", "limit 1-100, offset 0-100000 required")
	}
	return limit, offset, nil
}
func list(fn func(context.Context, auth.Principal, int, int) (json.RawMessage, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		l, o, e := paging(c)
		if e != nil {
			respond(c, nil, e)
			return
		}
		data, e := fn(c.Request.Context(), actor(c), l, o)
		respond(c, gin.H{"data": data, "limit": l, "offset": o}, e)
	}
}
