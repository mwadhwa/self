package storage

import "tw/dto"

type ParserStore interface {
	GetTransaction(key string) ([]*dto.Transaction, error)
	AddTransaction(key string, add *dto.Transaction) error
}
