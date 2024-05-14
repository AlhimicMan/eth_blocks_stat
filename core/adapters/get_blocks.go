package adapters

import (
	"context"
	"math/big"
)

type GetBlockClientI interface {
	GetLastBlockNumber(ctx context.Context) (big.Int, error)
	GetBlockRecord(ctx context.Context, blockNumber big.Int) (BlockRecord, error)
}
