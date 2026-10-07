package application

import (
	"chrisfoong/chaum-work-management-backend/internal/auth"
	"chrisfoong/chaum-work-management-backend/internal/config"
	"chrisfoong/chaum-work-management-backend/internal/contract"
	"chrisfoong/chaum-work-management-backend/internal/db"
	"chrisfoong/chaum-work-management-backend/internal/httpx"

	"chrisfoong/chaum-work-management-backend/internal/schema"
	"chrisfoong/chaum-work-management-backend/internal/server"
	"chrisfoong/chaum-work-management-backend/internal/work"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/joho/godotenv"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Main() {
	_ = godotenv.Load()
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	if e := Run(); e != nil {
		slog.Error("server stopped", "error", e)
		os.Exit(1)
	}
}
func Run() error {
	cfg, e := config.LoadRuntime(os.Getenv)
	if e != nil {
		return e
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	connect, cancel := context.WithTimeout(ctx, 15*time.Second)
	pool, e := db.NewPool(connect, cfg.DatabaseURL)
	if e != nil {
		cancel()
		return fmt.Errorf("database connection failed; verify deployment configuration")
	}
	defer pool.Close()
	if e = schema.Check(connect, pool); e != nil {
		cancel()
		return e
	}
	cancel()
	users := auth.NewUserStore(pool)
	httpx.UseJSONFieldNames()
	binding.EnableDecoderDisallowUnknownFields = true
	r := server.New(server.Deps{DB: pool, AllowedOrigins: os.Getenv("WEB_ALLOWED_ORIGINS"), WebAuth: auth.Authenticate(auth.NewLineVerifier(cfg.WebChannel), users.ByLineID, auth.RoleSupervisor, auth.RoleAssistant), WorkerAuth: auth.Authenticate(auth.NewLineVerifier(cfg.WorkerChannel), users.ByLineID, auth.RoleWorker)})
	r.Engine.OPTIONS("/*path", func(c *gin.Context) { c.Status(204) })
	svc := work.New(pool, cfg.QRSecret)
	files := work.NewStorage(cfg.StorageURL, cfg.StorageKey, cfg.StorageBucket)
	svc.Files = files
	if cfg.MessagingToken != "" {
		svc.Notify = &work.LinePush{Token: cfg.MessagingToken, Client: &http.Client{Timeout: 5 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}}
	} else {
		slog.Warn("LINE notifications disabled: messaging token missing")
	}
	contract.RegisterRoutes(r.Web, contract.NewHandler(contract.NewService(contract.TORRepository{}, pool, db.PoolTx(pool), contractNotify{svc})))
	work.Register(r.Web, r.Worker, svc)
	r.Web.POST("/files", files.Upload)
	r.Worker.POST("/files", files.Upload)
	r.Web.GET("/files", files.Read(svc))
	r.Worker.GET("/files", files.Read(svc))
	if cfg.QRSecret == "" {
		slog.Warn("QR issuance/check-in blocked: signing secret missing")
	}
	srv := &http.Server{Addr: ":" + cfg.Port, Handler: r.Engine, ReadHeaderTimeout: 10 * time.Second, ReadTimeout: 30 * time.Second, WriteTimeout: 30 * time.Second, IdleTimeout: 60 * time.Second}
	done := make(chan error, 1)
	go func() { done <- srv.ListenAndServe() }()
	jobDone := make(chan struct{})
	go func() {
		defer close(jobDone)
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				job, c := context.WithTimeout(ctx, 30*time.Second)
				_, e := svc.Finalize(job)
				c()
				if e != nil {
					slog.Error("attendance job failed")
				}
			}
		}
	}()
	select {
	case e := <-done:
		stop()
		<-jobDone
		if errors.Is(e, http.ErrServerClosed) {
			return nil
		}
		return fmt.Errorf("HTTP server failed")
	case <-ctx.Done():
	}
	shutdown, c := context.WithTimeout(context.Background(), 10*time.Second)
	defer c()
	e = srv.Shutdown(shutdown)
	<-jobDone
	return e
}
