package adapters

import (
	"context"
)

type GetBlockClientI interface {
	GetLastBlockNumber(ctx context.Context) (string, error)
	GetBlockRecord(ctx context.Context, blockNumber string) (BlockRecord, error)
}
