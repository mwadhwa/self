package api

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"tw/storage"
)

func TestNewEthereumParser(t *testing.T) {
	ctx, cancelFun := context.WithCancel(context.Background())
	defer cancelFun()
	parser := NewEthereumParser(ctx, storage.NewInMemoryStore())
	parser.Subscribe("0x95ad61b0a150d79219dcf64e1e6cc01f0b64c4ce")
	time.Sleep(10 * time.Second)
	txns := parser.GetTransactions("0x95ad61b0a150d79219dcf64e1e6cc01f0b64c4ce")
	for _, tx := range txns {
		println(tx)
	}
	blk := parser.GetCurrentBlock()
	assert.NotEqual(t, -1, blk)

}
