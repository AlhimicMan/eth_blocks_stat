package activity_calculator

import (
	"context"
	"eth_blocks_stat/core/adapters"
	"fmt"
	"strconv"
	"strings"
	"testing"
)

type GetBlocksMock struct {
	blocks       map[int]adapters.BlockRecord
	lastBlockNum int
}

func NewGetBlocksMock(blocks map[int]adapters.BlockRecord) *GetBlocksMock {
	var lastBlockNum int
	for bn := range blocks {
		if bn > lastBlockNum {
			lastBlockNum = bn
		}
	}
	mock := &GetBlocksMock{
		blocks:       blocks,
		lastBlockNum: lastBlockNum,
	}
	return mock
}

func (g *GetBlocksMock) GetLastBlockNumber(ctx context.Context) (string, error) {
	return strconv.Itoa(g.lastBlockNum), nil
}

func (g *GetBlocksMock) GetBlockRecord(ctx context.Context, blockNumber string) (adapters.BlockRecord, error) {
	fmt.Printf("retrieve block %s\n", blockNumber)
	bNumRaw := strings.TrimPrefix(blockNumber, "0x")
	bNum, err := strconv.ParseInt(bNumRaw, 16, 64)
	if err != nil {
		return adapters.BlockRecord{}, fmt.Errorf("block number is not a number: %w", err)
	}
	block, ok := g.blocks[int(bNum)]
	if !ok {
		return adapters.BlockRecord{}, fmt.Errorf("block %d not found", bNum)
	}
	return block, nil
}

func TestCalculatorStat(t *testing.T) {
	// Generate blocks for test
	// Here we have:
	// 4 transactions of address 0x75e89d5979e4f6fba9f97c104c2f0afb3f1dcb88 (1 source and 4 target)
	// 2 transactions from 0xc779dc08bb5ef038fc23a6a5ae38d8003adb2c53 (2 source)
	// 1 transaction to 0x43b603d4cdaed3dfa30855c9e354e300094a0a2d
	// 1 transaction from 0x58edf78281334335effa23101bbe3371b6a36a51
	// 1 transaction from 0x6837260d48e75f38b07c32b1cc28bcd866e00287
	mockBlocks := map[int]adapters.BlockRecord{
		1: {
			Number: "0x1",
			Transactions: []adapters.TransactionRec{
				{
					// To address 0x43b603d4cdaed3dfa30855c9e354e300094a0a2d
					BlockNumber:      "0x1",
					From:             "0xc779dc08bb5ef038fc23a6a5ae38d8003adb2c53",
					To:               "0xdac17f958d2ee523a2206206994597c13d831ec7",
					Input:            "0xa9059cbb00000000000000000000000043b603d4cdaed3dfa30855c9e354e300094a0a2d000000000000000000000000000000000000000000000000000000000c845880",
					TransactionIndex: "0x0",
					Hash:             "",
				},
				{
					// From address 0xcb83ca9633ad057bd88a48a5b6e8108d97ad4472
					BlockNumber:      "0x1",
					From:             "0x75e89d5979e4f6fba9f97c104c2f0afb3f1dcb88",
					To:               "0x74232704659ef37c08995e386a2e26cc27a8d7b1",
					Input:            "0xa9059cbb000000000000000000000000cb83ca9633ad057bd88a48a5b6e8108d97ad4472000000000000000000000000000000000000000000000003017f941f72d08000",
					TransactionIndex: "0x0",
					Hash:             "",
				},
			},
		},
		2: {
			Number: "0x2",
			Transactions: []adapters.TransactionRec{
				{
					// To address 0x75e89d5979e4f6fba9f97c104c2f0afb3f1dcb88
					BlockNumber:      "0x2",
					From:             "0xc779dc08bb5ef038fc23a6a5ae38d8003adb2c53",
					To:               "0xdac17f958d2ee523a2206206994597c13d831ec7",
					Input:            "0xa9059cbb00000000000000000000000075e89d5979e4f6fba9f97c104c2f0afb3f1dcb8800000000000000000000000000000000000000000000004be4e7267b6ae00000",
					TransactionIndex: "0x0",
					Hash:             "",
				},
				{
					// To address 0x75e89d5979e4f6fba9f97c104c2f0afb3f1dcb88
					BlockNumber:      "0x2",
					From:             "0x58edf78281334335effa23101bbe3371b6a36a51",
					To:               "0x74232704659ef37c08995e386a2e26cc27a8d7b1",
					Input:            "0xa9059cbb00000000000000000000000075e89d5979e4f6fba9f97c104c2f0afb3f1dcb88000000000000000000000000000000000000000000000003017f941f72d08000",
					TransactionIndex: "0x0",
					Hash:             "0xec500c389ff65a6bef3eedddb9a26fdb9656b5d85bf91ebb2e2ecfa7e46cf625",
				},
				{
					// To address 0x75e89d5979e4f6fba9f97c104c2f0afb3f1dcb88
					BlockNumber:      "0x2",
					From:             "0x6837260d48e75f38b07c32b1cc28bcd866e00287",
					To:               "0x74232704659ef37c08995e386a2e26cc27a8d7b1",
					Input:            "0xa9059cbb00000000000000000000000075e89d5979e4f6fba9f97c104c2f0afb3f1dcb88000000000000000000000000000000000000000000000003017f941f72d08000",
					TransactionIndex: "0x0",
					Hash:             "0xec500c389ff65a6bef3eedddb9a26fdb9656b5d85bf91ebb2e2ecfa7e46cf625",
				},
			},
		},
	}
	gbMock := NewGetBlocksMock(mockBlocks)
	calculator := NewCalculator(gbMock)
	ctx := context.Background()
	resStat, err := calculator.RetrieveTopAddresses(ctx, 2, 5)
	if err != nil {
		t.Errorf("failed to retrieve top addresses: %s", err.Error())
	}

	wantStat := []struct {
		address string
		items   int
	}{
		{
			address: "0x75e89d5979e4f6fba9f97c104c2f0afb3f1dcb88",
			items:   4,
		},
		{
			address: "0xc779dc08bb5ef038fc23a6a5ae38d8003adb2c53",
			items:   2,
		},
		{
			address: "0x43b603d4cdaed3dfa30855c9e354e300094a0a2d",
			items:   1,
		},
		{
			address: "0x58edf78281334335effa23101bbe3371b6a36a51",
			items:   1,
		},
		{
			address: "0x6837260d48e75f38b07c32b1cc28bcd866e00287",
			items:   1,
		},
	}
	if len(resStat.TopActiveAddresses) != len(wantStat) {
		t.Fatalf("top address count is not equal to %d, have %d", len(wantStat), len(resStat.TopActiveAddresses))
	}
	for i, want := range wantStat {
		addrStat := resStat.TopActiveAddresses[i]
		if addrStat.Address != want.address {
			t.Fatalf("expected %d address %s, have %s", i, want.address, addrStat.Address)
		}
		if addrStat.Transfers != want.items {
			t.Fatalf("expected %d address %d transfers, got %d", i, want.items, addrStat.Transfers)
		}
	}
}
