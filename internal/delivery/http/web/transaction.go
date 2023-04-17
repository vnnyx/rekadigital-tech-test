package web

type TransactionCreateReq struct {
	CustomerName string `json:"name"`
	Menu         string `json:"menu"`
	Price        int64  `json:"price"`
	Qty          int64  `json:"qty"`
	Payment      string `json:"payment"`
}

type TransactionDTO struct {
	TransactionID string `json:"id"`
	CustomerID    string `json:"customer_id,omitempty"`
	CustomerName  string `json:"name,omitempty"`
	Menu          string `json:"menu"`
	Price         int64  `json:"price"`
	Qty           int64  `json:"qty"`
	Payment       string `json:"payment"`
	Total         int64  `json:"total"`
}
