package service

import (
	"tw/api"
	"tw/dto"
)

type TransactionService struct {
	parser api.Parser
}

func NewTransactionService(parser api.Parser) *TransactionService {

	return &TransactionService{
		parser: parser,
	}
}

func (t *TransactionService) Subscribe(req *dto.TransactionRequest) error {
	t.parser.Subscribe(req.Address)
	return nil
}

func (t *TransactionService) GetCurrentBlock() int {
	return t.parser.GetCurrentBlock()
}

func (t *TransactionService) GetTransactions(req *dto.TransactionRequest) (*dto.TransactionResponse, error) {

	txns := t.parser.GetTransactions(req.Address)
	return &dto.TransactionResponse{
		Transactions: txns,
	}, nil
}
