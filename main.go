package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/go-playground/validator.v9"
	"log"
	"net/http"
	"tw/api"
	handler2 "tw/server/handler"
	"tw/service"
	"tw/storage"
)

var ctx = context.Background()

func handler() *handler2.TransactionHandler {
	store := storage.NewInMemoryStore()
	validater := validator.New()
	parser := api.NewEthereumParser(ctx, store)
	transactionService := service.NewTransactionService(parser)
	return handler2.NewTransactionHandler(transactionService, validater)

}

func main() {
	t := handler()
	r := mux.NewRouter()
	r.HandleFunc("/block", t.GetCurrentBlockHandler).Methods("GET")
	r.HandleFunc("/transactions/addresses/{address}", t.GetTransactionsHandler).Methods("GET")
	r.HandleFunc("/transactions/subscribe", t.SubscribeHandler).Methods("POST")
	fmt.Println("running server at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))

}
