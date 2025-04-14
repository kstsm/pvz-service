package tests

import (
	"context"
	"fmt"
	"github.com/gookit/slog"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kstsm/pvz-service/internal/handler"
	"github.com/kstsm/pvz-service/internal/repository"
	"github.com/kstsm/pvz-service/internal/service"
	"log"
	"net/http/httptest"
	"testing"
)

func InitTestPostgres(ctx context.Context) *pgxpool.Pool {
	user := "test_user"
	pass := "test_pass"
	host := "localhost"
	port := "5433"
	db := "test_db"

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, pass, host, port, db)

	slog.Info("Подключение к тестовой базе данных", "host", host, "port", port, "db", db)

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}

	slog.Info("Успешное подключение к тестовой базе данных")
	return pool
}

func SetupTestServer(t *testing.T) (*httptest.Server, context.Context, *pgxpool.Pool) {
	ctx := context.Background()
	conn := InitTestPostgres(ctx)
	t.Cleanup(func() {
		conn.Close()
	})

	repo := repository.NewRepository(conn)
	svc := service.NewService(repo)
	router := handler.NewRouterForTests(ctx, svc)
	ts := httptest.NewServer(router)

	return ts, ctx, conn
}
