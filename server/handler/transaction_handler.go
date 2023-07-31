package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"gopkg.in/go-playground/validator.v9"
	"log"
	"net/http"
	"tw/dto"
	"tw/service"
)

type TransactionHandler struct {
	transactionService *service.TransactionService
	validater          *validator.Validate
}

func NewTransactionHandler(transactionService *service.TransactionService, validater *validator.Validate) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
		validater:          validater,
	}
}

func (t *TransactionHandler) SubscribeHandler(w http.ResponseWriter, r *http.Request) {

	txnObj, err := t.validateAndDecodeTransactionReq(r)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	err = t.transactionService.Subscribe(txnObj)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (t *TransactionHandler) validateAndDecodeTransactionReq(r *http.Request) (*dto.TransactionRequest, error) {
	var transactionReq *dto.TransactionRequest
	err := json.NewDecoder(r.Body).Decode(&transactionReq)
	if err != nil {
		log.Printf("validation error. %s", err.Error())
		return nil, err
	}

	err = t.validater.Struct(transactionReq)
	if err != nil {
		log.Printf("validation error. %s", err.Error())
		return nil, err
	}
	return transactionReq, nil
}

func (t *TransactionHandler) GetCurrentBlockHandler(w http.ResponseWriter, r *http.Request) {

	blockNumber := t.transactionService.GetCurrentBlock()

	json.NewEncoder(w).Encode(blockNumber)
}

func (t *TransactionHandler) GetTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	address, found := params["address"]
	if !found {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	txnObj := &dto.TransactionRequest{
		Address: address,
	}
	err := t.validater.Struct(txnObj)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	txns, err := t.transactionService.GetTransactions(txnObj)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(txns)
}
