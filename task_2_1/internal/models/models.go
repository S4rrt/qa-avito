package models

// Statistics — статистика объявления.
type Statistics struct {
	Likes     int `json:"likes"`
	ViewCount int `json:"viewCount"`
	Contacts  int `json:"contacts"`
}

// Item — объявление, возвращаемое сервером.
type Item struct {
	ID         string     `json:"id"`
	SellerID   int        `json:"sellerId"`
	Name       string     `json:"name"`
	Price      int        `json:"price"`
	Statistics Statistics `json:"statistics"`
	CreatedAt  string     `json:"createdAt"`
}

// CreateItemRequest — тело запроса создания объявления.
type CreateItemRequest struct {
	SellerID   int        `json:"sellerID"`
	Name       string     `json:"name"`
	Price      int        `json:"price"`
	Statistics Statistics `json:"statistics"`
}

// CreateItemResponse — фактический ответ сервера на создание:
// {"status":"Сохранили объявление - <uuid>"}
type CreateItemResponse struct {
	Status string `json:"status"`
}

// ErrorResponse — тело ответа при ошибке.
type ErrorResponse struct {
	Result interface{} `json:"result"`
	Status string      `json:"status"`
}
