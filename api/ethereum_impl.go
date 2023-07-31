package api

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"
	"tw/dto"
	"tw/external/http"
	"tw/storage"
)

func NewEthereumParser(ctx context.Context, storage storage.ParserStore) Parser {

	e := &ethereumImpl{
		subscribedAddress: make([]string, 0),
		store:             storage, // todo change to DI
	}
	blkChan := make(chan string)
	go e.listenToBlocks(ctx, blkChan)
	go e.processBlocks(blkChan)
	return e
}

type ethereumImpl struct {
	subscribedAddress []string
	currBlock         string
	store             storage.ParserStore
}

func (e *ethereumImpl) GetCurrentBlock() int {
	num, err := hexToBlockNumber(e.currBlock)
	if err != nil {
		log.Printf("error while parsing block number %s", e.currBlock)
		return -1
	}

	return int(num.Int64())
}

func (e *ethereumImpl) Subscribe(address string) bool {
	for _, s := range e.subscribedAddress {
		if s == address {
			return false
		}
	}
	e.subscribedAddress = append(e.subscribedAddress, address)
	return true
}

func (e *ethereumImpl) GetTransactions(address string) []*dto.Transaction {
	txns, err := e.store.GetTransaction(address)
	if err != nil {
		if err == dto.ErrNoDataFound {
			println("No record found")
		}
		return []*dto.Transaction{}
	}
	return txns

}

func hexToBlockNumber(hexString string) (*big.Int, error) {
	blockNumber, success := new(big.Int).SetString(strings.TrimPrefix(hexString, "0x"), 16)
	if !success {
		return nil, fmt.Errorf("failed to parse block number")
	}
	return blockNumber, nil
}

// processBlocks process the transactions from a block given block number
func (e *ethereumImpl) processBlocks(blkChan chan string) {

	for blockHex := range blkChan {
		blockNumber, err := hexToBlockNumber(blockHex)
		if err != nil {
			log.Printf("unable to convert block hex to number. %s", blockHex)
			continue
		}
		block, err := http.GetBlock(blockNumber)
		if err != nil {
			if err == dto.ErrNoDataFound {
				log.Printf("Block has no further transactions %s", blockNumber)
			} else {
				log.Printf("Failed to retrieve block %v: %v", blockNumber, err)
			}

			continue
		}

		// Process each transaction in the block
		for _, tx := range block.Transactions {
			if contains(e.subscribedAddress, tx.To) {
				_ = e.store.AddTransaction(tx.To, tx)
			}
		}
	}
}

func contains(addresses []string, address string) bool {
	for _, a := range addresses {
		if a == address {
			return true
		}
	}
	return false
}

// listenToBlocks get the latest blocks to process
func (e *ethereumImpl) listenToBlocks(ctx context.Context, blkChan chan<- string) {

	for {
		select {
		case <-ctx.Done():
			log.Println("context is closed. closing parser")
			close(blkChan)
			log.Printf("parser is closed")

		default:
			blkNumber, err := http.GetLatestBlockNumber()
			if err != nil {
				log.Printf("error in subscription feed. %s", err.Error())
				continue
			}
			if e.currBlock != blkNumber {
				blkChan <- blkNumber
				e.currBlock = blkNumber
			}

		}
	}
}
