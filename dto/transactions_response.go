package dto

type TransactionResponse struct {
	Transactions []*Transaction `json:"transactions"`
}
