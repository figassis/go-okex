package helpers


type Config struct {
	APIKey    string `json:"ApiKey"`
	SecretKey string `json:"SecretKey"`
}

type ErrorCode struct {
	ErrorCode int `json:"error_code"`
}

type UserInfo struct {
	Result bool `json:"result"`
	Info   struct {
		Funds struct {
			Borrow, Free, Freezed map[string]interface{}
		} `json:"funds"`
	} `json:"info"`
}

type OrderInfo struct {
	Result bool `json:"result"`
	Orders []struct {
		Amount     float64 `json:"amount"`
		AvgPrice   int     `json:"avg_price"`
		CreateDate int64   `json:"create_date"`
		DealAmount int     `json:"deal_amount"`
		OrderID    int     `json:"order_id"`
		OrdersID   int     `json:"orders_id"`
		Price      int     `json:"price"`
		Status     int     `json:"status"`
		Symbol     string  `json:"symbol"`
		Type       string  `json:"type"`
	} `json:"orders"`
}

type OrderHistory struct {
	Result       bool `json:"result"`
	Total        int  `json:"total"`
	CurrencyPage int  `json:"currency_page"`
	PageLength   int  `json:"page_length"`
	Orders       []struct {
		Amount     int     `json:"amount"`
		AvgPrice   int     `json:"avg_price"`
		CreateDate int64   `json:"create_date"`
		DealAmount int     `json:"deal_amount"`
		OrderID    int     `json:"order_id"`
		OrdersID   int     `json:"orders_id"`
		Price      float64 `json:"price"`
		Status     int     `json:"status"`
		Symbol     string  `json:"symbol"`
		Type       string  `json:"type"`
	} `json:"orders"`
}

type Receipt struct {
	Result  bool `json:"result"`
	OrderID int  `json:"order_id"`
}
