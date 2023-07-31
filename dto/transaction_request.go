package dto

type TransactionRequest struct {
	Address string `json:"address" validate:"required"`
}
