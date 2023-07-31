package dto

import (
	"fmt"
	"math/big"
)

var ErrNoDataFound = fmt.Errorf("no Data Found")

type Transaction struct {
	Hash        string `json:"hash"`
	BlockNumber string `json:"blockNumber"`
	From        string `json:"from"`
	To          string `json:"to"`
	Value       string `json:"value"`
}

type Block struct {
	BlockNumber  *big.Int
	Transactions []*Transaction `json:"transactions"`
}
