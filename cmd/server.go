package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/gookit/slog"
	"github.com/kstsm/pvz-service/config"
	"github.com/kstsm/pvz-service/database"
	"github.com/kstsm/pvz-service/internal/handler"
	"github.com/kstsm/pvz-service/internal/repository"
	"github.com/kstsm/pvz-service/internal/service"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var cfg = config.Config

func Run() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	conn := database.InitPostgres(ctx)
	defer conn.Close()

	repo := repository.NewRepository(conn)
	svc := service.NewService(repo)
	router := handler.NewHandler(ctx, svc)
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler: router.NewRouter(),
	}

	errChan := make(chan error, 1)

	go func() {
		slog.Info("Запуск сервера", "host", cfg.Server.Host, "port", cfg.Server.Port)
		errChan <- srv.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		slog.Info("Завершаем сервер...")
	case err := <-errChan:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Fatal("Ошибка при запуске сервера", "error", err)
		}
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("Ошибка при завершении сервера", "error", err)
	}
}
