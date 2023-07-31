package storage

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"tw/dto"
)

func TestGetTransaction(t *testing.T) {

	store := NewInMemoryStore()
	txns, err := store.GetTransaction("xyz")
	assert.Empty(t, txns)
	assert.Error(t, err)
	tx := &dto.Transaction{To: "t1"}

	err = store.AddTransaction("xyz", tx)
	assert.Nil(t, err)
	txns, err = store.GetTransaction("xyz")
	assert.Nil(t, err)
	assert.Contains(t, txns, tx)

}

func TestAddTransaction(t *testing.T) {
	store := NewInMemoryStore()

	err := store.AddTransaction("xyz", &dto.Transaction{To: "t1"})
	assert.Nil(t, err)
}
