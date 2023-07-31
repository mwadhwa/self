package storage

import (
	"sync"
	"tw/dto"
)

func NewInMemoryStore() ParserStore {
	return &inMemoryStore{
		cache: sync.Map{},
	}
}

type inMemoryStore struct {
	cache sync.Map
}

func (i *inMemoryStore) GetTransaction(address string) ([]*dto.Transaction, error) {
	txns, found := i.cache.Load(address)
	if !found {
		return nil, dto.ErrNoDataFound
	}
	return txns.([]*dto.Transaction), nil
}

// AddTransaction adds the transaction to the address. Map concurrent safe however not thread safe for updating single key concurrently
func (i *inMemoryStore) AddTransaction(address string, add *dto.Transaction) error {
	txns, _ := i.cache.LoadOrStore(address, make([]*dto.Transaction, 0))
	// atomicity not needed as managed by single process
	i.cache.Store(address, append(txns.([]*dto.Transaction), add))
	return nil
}
