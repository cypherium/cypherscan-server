package main

import (
	"time"

	"github.com/cypherium/CypherTestNet/go-cypherium/core/types"
)

type ResponseOfGetBlocks struct {
	Total  int64          `json:"total"`
	Blocks []*listTxBlock `json:"blocks"`
}

type listTxBlock struct {
	Number   int64     `json:"number"`
	Time     time.Time `json:"createdAt"`
	Txn      int       `json:"txn"`
	GasUsed  uint64    `json:"gasUsed"`
	GasLimit uint64    `json:"gasLimit"`
}

func transferBlockHeadToListTxBlock(h *types.Header) *listTxBlock {
	return &listTxBlock{
		Number: h.Number.Int64(),
		Time:   time.Unix(h.Time.Int64(), 0),
		// txn
		GasUsed:  h.GasUsed,
		GasLimit: h.GasLimit,
	}
}

type numberDescSorterForListTxBlock []*listTxBlock

func (a numberDescSorterForListTxBlock) Len() int           { return len(a) }
func (a numberDescSorterForListTxBlock) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a numberDescSorterForListTxBlock) Less(i, j int) bool { return a[i].Number > a[j].Number }
