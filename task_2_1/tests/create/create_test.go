package create_test

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

// TC-01: Создание объявления с валидными данными → 200, поля заполнены.
func TestCreateItem_ValidData_Returns200(t *testing.T) {
	sellerID := testhelpers.RandomSellerID()
	req := testhelpers.DefaultRequest(sellerID)

	item, resp, err := api.MustCreateItem(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode, "body: %s", resp.Body)

	assert.NotEmpty(t, item.ID, "id must not be empty")
	assert.NotEmpty(t, item.CreatedAt, "createdAt must not be empty")
	assert.Equal(t, sellerID, item.SellerID)
	assert.Equal(t, req.Name, item.Name)
	assert.Equal(t, req.Price, item.Price)
}

// TC-02: Два одинаковых запроса → разные id (уникальность).
func TestCreateItem_DuplicateRequest_UniqueIDs(t *testing.T) {
	req := testhelpers.DefaultRequest(testhelpers.RandomSellerID())

	item1 := testhelpers.CreateItemOrFail(t, api, req)
	item2 := testhelpers.CreateItemOrFail(t, api, req)

	assert.NotEqual(t, item1.ID, item2.ID, "ids must be unique")
}

// TC-03: Статистика сохраняется корректно.
func TestCreateItem_StatisticsSaved(t *testing.T) {
	req := testhelpers.DefaultRequest(testhelpers.RandomSellerID())
	req.Statistics = models.Statistics{Likes: 10, ViewCount: 200, Contacts: 7}

	item := testhelpers.CreateItemOrFail(t, api, req)

	assert.Equal(t, 10, item.Statistics.Likes)
	assert.Equal(t, 200, item.Statistics.ViewCount)
	assert.Equal(t, 7, item.Statistics.Contacts)
}

// TC-04: Нулевая цена → сервер отклоняет с 400, хотя 0 валидно.
func TestCreateItem_ZeroPrice_ServerRejects(t *testing.T) {
	req := testhelpers.DefaultRequest(testhelpers.RandomSellerID())
	req.Price = 0

	resp, err := api.CreateItemRaw(req)
	require.NoError(t, err)
	// BUG: ожидается 200, фактически сервер возвращает 400
	assert.Equal(t, http.StatusOK, resp.StatusCode,
		"BUG: price=0 should be valid, but server rejects it: %s", resp.Body)
}

// TC-05: Нулевые значения статистики → сервер отклоняет с 400.
func TestCreateItem_ZeroStatistics_ServerRejects(t *testing.T) {
	req := testhelpers.DefaultRequest(testhelpers.RandomSellerID())
	req.Statistics = models.Statistics{}

	resp, err := api.CreateItemRaw(req)
	require.NoError(t, err)
	// BUG: ожидается 200, фактически сервер возвращает 400
	assert.Equal(t, http.StatusOK, resp.StatusCode,
		"BUG: zero statistics should be valid, but server rejects it: %s", resp.Body)
}

// TC-06: Отсутствует поле name → 400.
func TestCreateItem_MissingName_Returns400(t *testing.T) {
	body := map[string]interface{}{
		"sellerID": testhelpers.RandomSellerID(),
		"price":    500,
		"statistics": map[string]int{
			"likes": 1, "viewCount": 1, "contacts": 1,
		},
	}
	resp, err := api.CreateItemRaw(body)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "body: %s", resp.Body)
}

// TC-07: Отсутствует sellerID → 400.
func TestCreateItem_MissingSellerID_Returns400(t *testing.T) {
	body := map[string]interface{}{
		"name":  "No Seller",
		"price": 500,
		"statistics": map[string]int{
			"likes": 1, "viewCount": 1, "contacts": 1,
		},
	}
	resp, err := api.CreateItemRaw(body)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "body: %s", resp.Body)
}

// TC-08: Отрицательная цена → сервер принимает с 200, хотя должен отклонять.
func TestCreateItem_NegativePrice_ServerAccepts(t *testing.T) {
	req := testhelpers.DefaultRequest(testhelpers.RandomSellerID())
	req.Price = -1

	resp, err := api.CreateItemRaw(req)
	require.NoError(t, err)
	// BUG: ожидается 400, фактически сервер возвращает 200
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode,
		"BUG: negative price should be rejected, but server accepts it: %s", resp.Body)
}

// TC-09: Пустое тело → 400.
func TestCreateItem_EmptyBody_Returns400(t *testing.T) {
	resp, err := api.CreateItemRaw(map[string]interface{}{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "body: %s", resp.Body)
}

// TC-10: sellerID = 0 → 400.
func TestCreateItem_SellerIDZero_Returns400(t *testing.T) {
	req := testhelpers.DefaultRequest(0)
	resp, err := api.CreateItemRaw(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "body: %s", resp.Body)
}

// TC-11: Ответ содержит заголовок Content-Type.
func TestCreateItem_ResponseHasContentType(t *testing.T) {
	req := testhelpers.DefaultRequest(testhelpers.RandomSellerID())
	_, resp, err := api.MustCreateItem(req)
	require.NoError(t, err)
	assert.NotEmpty(t, resp.Header.Get("Content-Type"), "Content-Type header missing")
}

// TC-12: Ответ на создание содержит статус с id объявления.
func TestCreateItem_ResponseContainsID(t *testing.T) {
	req := testhelpers.DefaultRequest(testhelpers.RandomSellerID())
	resp, err := api.CreateItem(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode, "body: %s", resp.Body)

	var body struct {
		Status string `json:"status"`
	}
	require.NoError(t, json.Unmarshal(resp.Body, &body))
	assert.NotEmpty(t, body.Status, "status field must not be empty")
}
