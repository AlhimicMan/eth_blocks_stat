package activity_calculator

import (
	"context"
	"eth_blocks_stat/core/adapters"
	"fmt"
	"math/big"
	"sort"
	"strings"
	"sync"
)

const (
	erc20TransferSignature = "0xa9059cbb"
	blocksToProcess        = 100
)

type Calculator struct {
	adapter adapters.GetBlockClientI
}

func NewCalculator(adapter adapters.GetBlockClientI) *Calculator {
	return &Calculator{adapter: adapter}
}

func (c *Calculator) Initialize() error {
	return nil
}

// isERC20Transfer checks if transfer(address,uint256) token transfer in transaction input
func (c *Calculator) isERC20Transfer(txInput string) bool {
	return len(txInput) == 138 && txInput[:10] == erc20TransferSignature
}

func (c *Calculator) getERC20TransferToAddress(txInput string) string {
	txRunes := []rune(txInput)
	addr := string(txRunes[10:74])
	addr = strings.TrimLeft(addr, "0")
	return "0x" + addr
}

func (c *Calculator) listBlockERC20ActivityStat(block adapters.BlockRecord) (map[string]int, error) {
	addStats := map[string]int{}
	for _, tx := range block.Transactions {
		isErc20 := c.isERC20Transfer(tx.Input)
		if !isErc20 {
			continue
		}
		aStatFrom, ok := addStats[tx.From]
		if !ok {
			aStatFrom = 0
		}
		addStats[tx.From] = aStatFrom + 1

		addrTo := c.getERC20TransferToAddress(tx.Input)
		aStatTo, ok := addStats[addrTo]
		if !ok {
			aStatTo = 0
		}
		addStats[addrTo] = aStatTo + 1

	}

	return addStats, nil
}

func (c *Calculator) retrieveBlockStat(ctx context.Context, blockNum string) (map[string]int, error) {
	block, err := c.adapter.GetBlockRecord(ctx, blockNum)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve block %s records: %w", blockNum, err)
	}
	blockStat, err := c.listBlockERC20ActivityStat(block)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve block %s stats: %w", block.Number, err)
	}
	return blockStat, nil
}

func (c *Calculator) RetrieveTopAddresses(ctx context.Context) (TopActiveAddressesRes, error) {
	lastBlockNumber, err := c.adapter.GetLastBlockNumber(ctx)
	if err != nil {
		return TopActiveAddressesRes{}, fmt.Errorf("failed to retrieve last block number: %w", err)
	}
	blockNum, ok := new(big.Int).SetString(lastBlockNumber, 16)
	if !ok {
		return TopActiveAddressesRes{}, fmt.Errorf("cannot parse block number: %w", err)
	}

	statWg := sync.WaitGroup{}
	blocksProcessed := 0
	var responsesChan = make(chan map[string]int, blocksToProcess)
	for i := 0; i <= blocksToProcess; i++ {
		blockNum.Sub(blockNum, big.NewInt(1))
		blockNumParam := fmt.Sprintf("0x%x", blockNum)
		statWg.Add(1)
		go func(bn string) {
			defer statWg.Done()
			blockStat, err := c.retrieveBlockStat(ctx, bn)
			if err != nil {
				fmt.Printf("failed to retrieve block %s stats: %s", blockNum.String(), err.Error())
			}
			responsesChan <- blockStat
		}(blockNumParam)
	}
	go func() {
		statWg.Wait()
		close(responsesChan)
	}()

	resultedStat := make(map[string]int)
	for blockStat := range responsesChan {
		for addr, transfers := range blockStat {
			aCount, ok := resultedStat[addr]
			if !ok {
				aCount = 0
			}
			resultedStat[addr] = aCount + transfers
		}
		blocksProcessed += 1
	}

	if blocksProcessed < blocksToProcess {
		return TopActiveAddressesRes{}, fmt.Errorf("processed only %d blocks of %d. See logs for details",
			blocksProcessed, blocksToProcess)
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
