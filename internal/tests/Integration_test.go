package tests

import (
	"github.com/kstsm/pvz-service/internal/repository"
	"github.com/kstsm/pvz-service/internal/service"
	"testing"
)

func TestReceptionIntegration(t *testing.T) {
	ts, ctx, pool := SetupTestServer(t)
	defer ts.Close()
	t.Cleanup(func() {
		pool.Close()
	})
	t.Log("Сервер и пул подключений успешно инициализированы")

	repo := repository.NewRepository(pool)
	svc := service.NewService(repo)
	t.Log("Репозиторий и сервис инициализированы")

	pvz, err := svc.CreatePVZ(ctx, "Москва")
	if err != nil {
		t.Fatalf("Ошибка при создании ПВЗ: %v", err)
	}
	t.Logf("ПВЗ создан: ID=%d", pvz.ID)

	reception, err := svc.CreateReception(ctx, pvz.ID)
	if err != nil {
		t.Fatalf("Ошибка при создании приёмки: %v", err)
	}
	t.Logf("Приёмка создана: ID=%d", reception.ID)

	for i := 0; i < 50; i++ {
		productType := "электроника"
		product, err := svc.AddProductToActiveReception(ctx, productType, pvz.ID)
		if err != nil {
			t.Fatalf("Ошибка при добавлении товара #%d: %v", i+1, err)
		}
		if product.ReceptionID != reception.ID {
			t.Fatalf("Неверный ReceptionID у товара #%d: ожидался %v, получен %v", i+1, reception.ID, product.ReceptionID)
		}
	}
	t.Logf("50 товаров добавлены к приёмке ID=%d", reception.ID)

	closedReception, err := svc.CloseLastReception(ctx, pvz.ID)
	if err != nil {
		t.Fatalf("Ошибка при закрытии приёмки: %v", err)
	}
	t.Logf("Приёмка закрыта: ID=%d, Статус=%s", closedReception.ID, closedReception.Status)

	expectedStatus := "close"
	if closedReception.Status != expectedStatus {
		t.Fatalf("Некорректный статус приёмки: ожидался %q, получен %q", expectedStatus, closedReception.Status)
	}
	t.Log("Статус приёмки успешно проверен")
}
