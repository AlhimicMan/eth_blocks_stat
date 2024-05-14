package activity_calculator

import (
	"context"
	"eth_blocks_stat/core/adapters"
	"fmt"
	"math/big"
	"sort"
)

const erc20TransferSignature = "0xa9059cbb"

type Calculator struct {
	adapter adapters.GetBlockClientI
}

func NewCalculator(adapter adapters.GetBlockClientI) *Calculator {
	return &Calculator{adapter: adapter}
}

func (c *Calculator) Initialize() error {
	return nil
}

func (c *Calculator) isERC20(txInput string) bool {
	return len(txInput) >= 10 && txInput[:10] == erc20TransferSignature
}

func (c *Calculator) listTransactionsERC20ActivityStat(transactions []adapters.TransactionRec) (map[string]int, error) {
	addStats := map[string]int{}
	for _, tx := range transactions {
		isErc20 := c.isERC20(tx.Input)
		if !isErc20 {
			continue
		}
		aStatFrom, ok := addStats[tx.From]
		if !ok {
			aStatFrom = 0
		}
		addStats[tx.From] = aStatFrom + 1

		aStatTo, ok := addStats[tx.To]
		if !ok {
			aStatTo = 0
		}
		addStats[tx.To] = aStatTo + 1
	}

	return addStats, nil
}

func (c *Calculator) RetrieveTopAddresses(ctx context.Context) (TopActiveAddressesRes, error) {
	lastBlockNumber, err := c.adapter.GetLastBlockNumber(ctx)
	if err != nil {
		return TopActiveAddressesRes{}, fmt.Errorf("failed to retrieve last block number: %w", err)
	}
	resultedStat := make(map[string]int)
	for i := 0; i <= 100; i++ {
		lastBlockNumber.Sub(&lastBlockNumber, big.NewInt(1))
		fmt.Printf("processing block %s from top %d", lastBlockNumber.String(), i)
		block, err := c.adapter.GetBlockRecord(ctx, lastBlockNumber)
		if err != nil {
			return TopActiveAddressesRes{}, fmt.Errorf("failed to retrieve block %s records: %w",
				lastBlockNumber.String(), err)
		}
		blockStat, err := c.listTransactionsERC20ActivityStat(block.Transactions)
		if err != nil {
			return TopActiveAddressesRes{}, fmt.Errorf("failed to retrieve block %s stats: %w",
				block.Number, err)
		}
		for addr, transfers := range blockStat {
			aCount, ok := resultedStat[addr]
			if !ok {
				aCount = 0
			}
			resultedStat[addr] = aCount + transfers
		}
	}
	addrStats := make([]ActiveAddressRes, 0, len(resultedStat))
	for addr, transfers := range resultedStat {
		stRec := ActiveAddressRes{
			Address:   addr,
			Transfers: transfers,
		}
		addrStats = append(addrStats, stRec)
	}
	sort.SliceStable(addrStats, func(i, j int) bool {
		return addrStats[i].Transfers > addrStats[j].Transfers
	})
	resStat := TopActiveAddressesRes{}
	if len(addrStats) > 5 {
		resStat.TopActiveAddresses = addrStats[:5]
	} else {
		resStat.TopActiveAddresses = addrStats
	}
	return resStat, nil
}
