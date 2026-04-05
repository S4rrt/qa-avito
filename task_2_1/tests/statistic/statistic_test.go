package statistic_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"task_2_1/internal/models"
	"task_2_1/internal/testhelpers"
)

var api = testhelpers.NewClient()

// TC-23: Статистика по существующему объявлению (v1) → 200.
func TestGetStatisticV1_ExistingID_Returns200(t *testing.T) {
	created := testhelpers.CreateItemOrFail(t, api, testhelpers.DefaultRequest(testhelpers.RandomSellerID()))

	stats, resp, err := api.MustGetStatisticV1(created.ID)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode, "body: %s", resp.Body)
	assert.NotEmpty(t, stats, "statistics array must not be empty")
}

// TC-24: Значения статистики совпадают с переданными при создании (v1).
func TestGetStatisticV1_ValuesMatchCreation(t *testing.T) {
	req := testhelpers.DefaultRequest(testhelpers.RandomSellerID())
	req.Statistics = models.Statistics{Likes: 15, ViewCount: 300, Contacts: 9}
	created := testhelpers.CreateItemOrFail(t, api, req)

	stats, resp, err := api.MustGetStatisticV1(created.ID)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode, "body: %s", resp.Body)
	require.NotEmpty(t, stats)

	assert.Equal(t, 15, stats[0].Likes)
	assert.Equal(t, 300, stats[0].ViewCount)
	assert.Equal(t, 9, stats[0].Contacts)
}

// TC-25: Несуществующий id (v1) → 404.
func TestGetStatisticV1_NonExistentID_Returns404(t *testing.T) {
	resp, err := api.GetStatisticV1(testhelpers.NonExistentUUID)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode, "body: %s", resp.Body)
}

// TC-26: Невалидный формат id (v1) → 400.
func TestGetStatisticV1_InvalidID_Returns400(t *testing.T) {
	resp, err := api.GetStatisticV1(testhelpers.InvalidID)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "body: %s", resp.Body)
}

// TC-27: Статистика по существующему объявлению (v2) → 200.
func TestGetStatisticV2_ExistingID_Returns200(t *testing.T) {
	created := testhelpers.CreateItemOrFail(t, api, testhelpers.DefaultRequest(testhelpers.RandomSellerID()))

	resp, err := api.GetStatisticV2(created.ID)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode, "body: %s", resp.Body)
}

// TC-28: Несуществующий id (v2) → 404.
func TestGetStatisticV2_NonExistentID_Returns404(t *testing.T) {
	resp, err := api.GetStatisticV2(testhelpers.NonExistentUUID)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode, "body: %s", resp.Body)
}
