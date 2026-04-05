// Package testhelpers содержит общие утилиты для всех тест-пакетов.
package testhelpers

import (
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"task_2_1/internal/client"
	"task_2_1/internal/models"
)

const (
	BaseURL         = "https://qa-internship.avito.com"
	NonExistentUUID = "00000000-0000-0000-0000-000000000000"
	InvalidID       = "not-a-valid-uuid!!!"
)

// NewClient возвращает клиент, настроенный на тестовый стенд.
func NewClient() *client.Client {
	return client.New(BaseURL)
}

// RandomSellerID возвращает случайный sellerID из диапазона [111111, 999999].
func RandomSellerID() int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return 111111 + r.Intn(888889)
}

// DefaultRequest возвращает валидный запрос создания объявления.
func DefaultRequest(sellerID int) models.CreateItemRequest {
	return models.CreateItemRequest{
		SellerID: sellerID,
		Name:     "Test Item",
		Price:    1000,
		Statistics: models.Statistics{
			Likes:     5,
			ViewCount: 100,
			Contacts:  3,
		},
	}
}

// CreateItemOrFail создаёт объявление и завершает тест при любой ошибке.
func CreateItemOrFail(t *testing.T, c *client.Client, req models.CreateItemRequest) *models.Item {
	t.Helper()
	item, resp, err := c.MustCreateItem(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode, "CreateItemOrFail: body: %s", resp.Body)
	require.NotEmpty(t, item.ID, "CreateItemOrFail: got empty id")
	return item
}
