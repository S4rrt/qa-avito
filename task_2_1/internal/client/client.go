// Package client предоставляет типизированный HTTP-клиент для работы с API.
package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"task_2_1/internal/models"
)

const defaultTimeout = 15 * time.Second

// Client — HTTP-клиент сервиса объявлений.
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// New создаёт новый Client.
func New(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}
}

// Response — сырой HTTP-ответ с телом.
type Response struct {
	StatusCode int
	Body       []byte
	Header     http.Header
}

// ─── raw ──────────────────────────────────────────────────────────────────────

func (c *Client) do(method, path string, body interface{}) (*Response, error) {
	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal body: %w", err)
		}
		reqBody = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, c.baseURL+path, reqBody)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}
	return &Response{
		StatusCode: resp.StatusCode,
		Body:       raw,
		Header:     resp.Header,
	}, nil
}

// ─── API v1 ───────────────────────────────────────────────────────────────────

// CreateItem — POST /api/1/item
func (c *Client) CreateItem(req models.CreateItemRequest) (*Response, error) {
	return c.do(http.MethodPost, "/api/1/item", req)
}

// CreateItemRaw позволяет передать произвольное тело (для негативных тестов).
func (c *Client) CreateItemRaw(body interface{}) (*Response, error) {
	return c.do(http.MethodPost, "/api/1/item", body)
}

// GetItemByID — GET /api/1/item/:id
func (c *Client) GetItemByID(id string) (*Response, error) {
	return c.do(http.MethodGet, "/api/1/item/"+id, nil)
}

// GetItemsBySellerID — GET /api/1/:sellerID/item
func (c *Client) GetItemsBySellerID(sellerID int) (*Response, error) {
	return c.do(http.MethodGet, fmt.Sprintf("/api/1/%d/item", sellerID), nil)
}

// GetItemsBySellerIDRaw — GET /api/1/:sellerID/item с произвольным sellerID в пути.
func (c *Client) GetItemsBySellerIDRaw(sellerID string) (*Response, error) {
	return c.do(http.MethodGet, "/api/1/"+sellerID+"/item", nil)
}

// GetStatisticV1 — GET /api/1/statistic/:id
func (c *Client) GetStatisticV1(id string) (*Response, error) {
	return c.do(http.MethodGet, "/api/1/statistic/"+id, nil)
}

// ─── API v2 ───────────────────────────────────────────────────────────────────

// GetStatisticV2 — GET /api/2/statistic/:id
func (c *Client) GetStatisticV2(id string) (*Response, error) {
	return c.do(http.MethodGet, "/api/2/statistic/"+id, nil)
}

// ─── Typed helpers ────────────────────────────────────────────────────────────

// extractIDFromStatus вытаскивает UUID из строки вида
// "Сохранили объявление - <uuid>"
func extractIDFromStatus(status string) string {
	parts := strings.SplitN(status, " - ", 2)
	if len(parts) == 2 {
		return strings.TrimSpace(parts[1])
	}
	return ""
}

// MustCreateItem создаёт объявление, извлекает id из ответа сервера,
// затем получает полный Item через GET /api/1/item/:id.
// Сервер возвращает {"status":"Сохранили объявление - <uuid>"} — не Item напрямую.
func (c *Client) MustCreateItem(req models.CreateItemRequest) (*models.Item, *Response, error) {
	resp, err := c.CreateItem(req)
	if err != nil {
		return nil, nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, resp, nil
	}

	// Парсим фактический ответ сервера
	var createResp models.CreateItemResponse
	if err := json.Unmarshal(resp.Body, &createResp); err != nil {
		return nil, resp, fmt.Errorf("unmarshal create response: %w (body: %s)", err, resp.Body)
	}

	id := extractIDFromStatus(createResp.Status)
	if id == "" {
		return nil, resp, fmt.Errorf("could not extract id from status: %q", createResp.Status)
	}

	// Получаем полный Item по id
	items, _, err := c.MustGetItemByID(id)
	if err != nil {
		return nil, resp, fmt.Errorf("get item after create: %w", err)
	}
	if len(items) == 0 {
		return nil, resp, fmt.Errorf("item %s not found after create", id)
	}
	return &items[0], resp, nil
}

// MustGetItemByID десериализует ответ в []Item.
func (c *Client) MustGetItemByID(id string) ([]models.Item, *Response, error) {
	resp, err := c.GetItemByID(id)
	if err != nil {
		return nil, nil, err
	}
	var items []models.Item
	if err := json.Unmarshal(resp.Body, &items); err != nil {
		return nil, resp, fmt.Errorf("unmarshal items: %w (body: %s)", err, resp.Body)
	}
	return items, resp, nil
}

// MustGetStatisticV1 десериализует ответ в []Statistics.
func (c *Client) MustGetStatisticV1(id string) ([]models.Statistics, *Response, error) {
	resp, err := c.GetStatisticV1(id)
	if err != nil {
		return nil, nil, err
	}
	var stats []models.Statistics
	if err := json.Unmarshal(resp.Body, &stats); err != nil {
		return nil, resp, fmt.Errorf("unmarshal stats: %w (body: %s)", err, resp.Body)
	}
	return stats, resp, nil
}
