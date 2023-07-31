package http

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strings"
	"tw/dto"
)

const (
	rpcURL = "https://cloudflare-eth.com"
)

type rpcRequest struct {
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

type rpcResponse struct {
	Jsonrpc string          `json:"jsonrpc"`
	ID      int             `json:"id"`
	Result  json.RawMessage `json:"result"`
	Error   *rpcError       `json:"error"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func GetLatestBlockNumber() (string, error) {
	requestData := &rpcRequest{
		Jsonrpc: "2.0",
		Method:  "eth_blockNumber",
		Params:  []interface{}{""},
		ID:      1,
	}
	resp, err := getRPCResponse(requestData)
	if err != nil {
		return "", err
	}
	var blockNumber string
	if err := json.Unmarshal(resp.Result, &blockNumber); err != nil {
		log.Printf("Failed to parse block header: %v", err)
		return "", err
	}

	return blockNumber, nil
}

func GetBlock(blockNumber *big.Int) (*dto.Block, error) {
	requestData := &rpcRequest{
		Jsonrpc: "2.0",
		Method:  "eth_getBlockByNumber",
		Params:  []interface{}{fmt.Sprintf("0x%X", blockNumber), true},
		ID:      1,
	}

	resp, err := getRPCResponse(requestData)
	if err != nil {
		log.Printf("unable to fetch transaction from block %s", blockNumber)
	}
	// current block has no further transactions
	if resp == nil || resp.Result == nil {
		log.Printf("no transaction found for block %s", blockNumber)
		return nil, dto.ErrNoDataFound
	}

	var block *dto.Block
	if err := json.Unmarshal(resp.Result, &block); err != nil {
		return nil, fmt.Errorf("failed to parse block data: %v", err)
	}
	block.BlockNumber = blockNumber
	return block, nil
}

func getRPCResponse(requestData *rpcRequest) (*rpcResponse, error) {
	requestJSON, err := json.Marshal(requestData)
	if err != nil {
		log.Printf("Failed to create subscribe request: %v", err)
		return nil, err
	}

	resp, err := http.Post(rpcURL, "application/json", strings.NewReader(string(requestJSON)))
	if err != nil {
		log.Printf("Failed to subscribe to new block headers: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	var response *rpcResponse
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&response); err != nil {
		log.Printf("Failed to decode subscribe response: %v", err)
		return nil, err
	}

	if response.Error != nil {
		log.Printf("Subscribe error: %v", response.Error.Message)
		return nil, fmt.Errorf(response.Error.Message)
	}

	if response.Result == nil {
		log.Printf("Empty response received for new block header subscription")
		return nil, dto.ErrNoDataFound
	}
	return response, nil
}
