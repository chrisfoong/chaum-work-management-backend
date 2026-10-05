// Command server runs the Chaum Resource Management API.
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"chrisfoong/chaum-work-management-backend/internal/auth"
	"chrisfoong/chaum-work-management-backend/internal/config"
	"chrisfoong/chaum-work-management-backend/internal/contract"
	"chrisfoong/chaum-work-management-backend/internal/db"
	"chrisfoong/chaum-work-management-backend/internal/httpx"
	"chrisfoong/chaum-work-management-backend/internal/notify"
	"chrisfoong/chaum-work-management-backend/internal/server"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	if err := run(); err != nil {
		slog.Error("server stopped", "error", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load(os.Getenv)
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := db.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()

	supabase, err := auth.NewSupabaseVerifier(ctx, cfg.SupabaseURL)
	if err != nil {
		return err
	}
	users := auth.NewUserStore(pool)

	httpx.UseJSONFieldNames()
	router := server.New(server.Deps{
		DB:         pool,
		WebAuth:    auth.Authenticate(supabase, users.ByUserID, auth.RoleSupervisor, auth.RoleAssistant),
		WorkerAuth: auth.Authenticate(auth.NewLineVerifier(cfg.LineChannelID), users.ByLineID, auth.RoleWorker),
	})

	// The real LINE notifier comes in P11; until then every message is logged as skipped.
	notifier := notify.Disabled{}
	slog.Warn("LINE notifier not configured: notifications are skipped, not delivered")

	contractSvc := contract.NewService(contract.TORRepository{}, pool, db.PoolTx(pool), notifier)
	contract.RegisterRoutes(router.Web, contract.NewHandler(contractSvc))

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router.Engine,
		ReadHeaderTimeout: 10 * time.Second,
	}
	errCh := make(chan error, 1)
	go func() {
		slog.Info("listening", "addr", srv.Addr)
		errCh <- srv.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-ctx.Done():
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return srv.Shutdown(shutdownCtx)
}
