package getbyid_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"task_2_1/internal/testhelpers"
)

var api = testhelpers.NewClient()

// TC-13: Получить объявление по существующему id → 200, данные совпадают.
func TestGetItemByID_ExistingID_Returns200(t *testing.T) {
	created := testhelpers.CreateItemOrFail(t, api, testhelpers.DefaultRequest(testhelpers.RandomSellerID()))

	items, resp, err := api.MustGetItemByID(created.ID)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode, "body: %s", resp.Body)
	require.NotEmpty(t, items, "items array must not be empty")

	assert.Equal(t, created.ID, items[0].ID)
}

// TC-14: Данные объявления консистентны с созданием.
func TestGetItemByID_DataConsistency(t *testing.T) {
	req := testhelpers.DefaultRequest(testhelpers.RandomSellerID())
	req.Name = "Consistency Check"
	req.Price = 42000
	created := testhelpers.CreateItemOrFail(t, api, req)

	items, resp, err := api.MustGetItemByID(created.ID)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode, "body: %s", resp.Body)
	require.NotEmpty(t, items)

	got := items[0]
	assert.Equal(t, req.Name, got.Name)
	assert.Equal(t, req.Price, got.Price)
	assert.Equal(t, req.SellerID, got.SellerID)
}

// TC-15: Несуществующий id → 404.
func TestGetItemByID_NonExistentID_Returns404(t *testing.T) {
	resp, err := api.GetItemByID(testhelpers.NonExistentUUID)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode, "body: %s", resp.Body)
}

// TC-16: Невалидный формат id → 400.
func TestGetItemByID_InvalidIDFormat_Returns400(t *testing.T) {
	resp, err := api.GetItemByID(testhelpers.InvalidID)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "body: %s", resp.Body)
}

// TC-17: Идемпотентность — два GET возвращают одинаковый результат.
func TestGetItemByID_Idempotent(t *testing.T) {
	created := testhelpers.CreateItemOrFail(t, api, testhelpers.DefaultRequest(testhelpers.RandomSellerID()))

	resp1, err := api.GetItemByID(created.ID)
	require.NoError(t, err)
	resp2, err := api.GetItemByID(created.ID)
	require.NoError(t, err)

	assert.Equal(t, string(resp1.Body), string(resp2.Body), "GET must be idempotent")
}
