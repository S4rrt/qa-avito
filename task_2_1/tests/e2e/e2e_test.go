package e2e_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"task_2_1/internal/models"
	"task_2_1/internal/testhelpers"
)

var api = testhelpers.NewClient()

// TC-29: Создать → получить по id → проверить данные → статистика v1 и v2.
func TestE2E_CreateGetStatistic(t *testing.T) {
	req := models.CreateItemRequest{
		SellerID: testhelpers.RandomSellerID(),
		Name:     "E2E Test Item",
		Price:    9999,
		Statistics: models.Statistics{
			Likes: 7, ViewCount: 150, Contacts: 4,
		},
	}

	// Шаг 1: создать
	created := testhelpers.CreateItemOrFail(t, api, req)
	t.Logf("created id=%s", created.ID)

	// Шаг 2: получить по id
	items, resp, err := api.MustGetItemByID(created.ID)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode, "body: %s", resp.Body)
	require.NotEmpty(t, items)
	assert.Equal(t, req.Name, items[0].Name)
	assert.Equal(t, req.Price, items[0].Price)

	// Шаг 3: статистика v1
	stats, resp3, err := api.MustGetStatisticV1(created.ID)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp3.StatusCode, "body: %s", resp3.Body)
	assert.NotEmpty(t, stats)

	// Шаг 4: статистика v2
	resp4, err := api.GetStatisticV2(created.ID)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp4.StatusCode, "body: %s", resp4.Body)
}

// TC-30: Создать несколько объявлений → все присутствуют в списке продавца.
func TestE2E_CreateMultipleItems_AllPresentInSellerList(t *testing.T) {
	sellerID := testhelpers.RandomSellerID()
	const count = 3

	createdIDs := make(map[string]struct{}, count)
	for i := range count {
		req := testhelpers.DefaultRequest(sellerID)
		req.Price = (i + 1) * 1000
		item := testhelpers.CreateItemOrFail(t, api, req)
		createdIDs[item.ID] = struct{}{}
	}

	resp, err := api.GetItemsBySellerID(sellerID)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode, "body: %s", resp.Body)

	var listed []models.Item
	require.NoError(t, json.Unmarshal(resp.Body, &listed))

	assert.GreaterOrEqual(t, len(listed), count, "seller list must contain all created items")

	for _, it := range listed {
		delete(createdIDs, it.ID)
	}
	assert.Empty(t, createdIDs, "these created items were not found in seller list")
}
