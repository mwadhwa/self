package api

import (
	"tw/dto"
)

type Parser interface {

	// last parsed block
	GetCurrentBlock() int

	// add address to observer
	Subscribe(address string) bool

	// list of inbound or outbound transactions for an address
	GetTransactions(address string) []*dto.Transaction
}
