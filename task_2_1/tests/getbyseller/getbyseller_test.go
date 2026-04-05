package getbyseller_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"task_2_1/internal/testhelpers"
)

var api = testhelpers.NewClient()

// TC-18: Список объявлений существующего продавца → 200.
func TestGetItemsBySellerID_ExistingSeller_Returns200(t *testing.T) {
	sellerID := testhelpers.RandomSellerID()
	testhelpers.CreateItemOrFail(t, api, testhelpers.DefaultRequest(sellerID))

	resp, err := api.GetItemsBySellerID(sellerID)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode, "body: %s", resp.Body)
}

// TC-19: Все объявления в ответе принадлежат запрошенному продавцу.
func TestGetItemsBySellerID_AllItemsBelongToSeller(t *testing.T) {
	sellerID := testhelpers.RandomSellerID()
	testhelpers.CreateItemOrFail(t, api, testhelpers.DefaultRequest(sellerID))
	testhelpers.CreateItemOrFail(t, api, testhelpers.DefaultRequest(sellerID))

	resp, err := api.GetItemsBySellerID(sellerID)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode, "body: %s", resp.Body)

	var items []struct {
		ID       string `json:"id"`
		SellerID int    `json:"sellerId"`
	}
	require.NoError(t, json.Unmarshal(resp.Body, &items))

	for _, it := range items {
		assert.Equal(t, sellerID, it.SellerID, "item %s has wrong sellerId", it.ID)
	}
}

// TC-20: Нечисловой sellerID в пути → 400.
func TestGetItemsBySellerID_NonNumericSellerID_Returns400(t *testing.T) {
	resp, err := api.GetItemsBySellerIDRaw("abc")
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "body: %s", resp.Body)
}

// TC-21: Продавец без объявлений → 200 (пустой массив) или 404.
func TestGetItemsBySellerID_SellerWithNoItems_Returns200OrEmpty(t *testing.T) {
	resp, err := api.GetItemsBySellerID(testhelpers.RandomSellerID())
	require.NoError(t, err)
	assert.Contains(t, []int{http.StatusOK, http.StatusNotFound}, resp.StatusCode,
		"expected 200 or 404, body: %s", resp.Body)
}

// TC-22: Идемпотентность GET списка.
func TestGetItemsBySellerID_Idempotent(t *testing.T) {
	sellerID := testhelpers.RandomSellerID()
	testhelpers.CreateItemOrFail(t, api, testhelpers.DefaultRequest(sellerID))

	resp1, err := api.GetItemsBySellerID(sellerID)
	require.NoError(t, err)
	resp2, err := api.GetItemsBySellerID(sellerID)
	require.NoError(t, err)

	assert.Equal(t, string(resp1.Body), string(resp2.Body), "GET must be idempotent")
}
