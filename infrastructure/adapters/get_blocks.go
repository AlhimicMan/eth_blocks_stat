package adapters

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"eth_blocks_stat/core/adapters"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"reflect"
)

const (
	rpcVersion = "2.0"
	defaultId  = "getblock.io"
)

type GetBlockClient struct {
	apiKey string
	apiURL string
}

type JSONRpcReq struct {
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	Id      string        `json:"id"`
}

type JSONRpcRes struct {
	Id      string      `json:"id"`
	Jsonrpc string      `json:"jsonrpc"`
	Result  interface{} `json:"result"`
}

func NewGetBlockClient(apiKey string) adapters.GetBlockClientI {
	apiURL := fmt.Sprintf("https://go.getblock.io/%s", apiKey)
	return &GetBlockClient{
		apiKey: apiKey,
		apiURL: apiURL,
	}
}

func (g *GetBlockClient) performRPCCall(ctx context.Context, method string, params []interface{}, result interface{}) error {
	if reflect.TypeOf(result).Kind() != reflect.Ptr {
		return errors.New("result parameter must be pointer to struct where response will be decoded")
	}
	req := JSONRpcReq{
		Jsonrpc: rpcVersion,
		Id:      defaultId,
		Method:  method,
		Params:  params,
	}

	payload, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("cannot prepare request: %w", err)
	}

	res, err := http.Post(g.apiURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("cannot send request: %w", err)
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	callRes := JSONRpcRes{
		Result: result,
	}
	resData, _ := io.ReadAll(res.Body)
	fmt.Printf("resulted body: %s\n", string(resData))
	err = json.NewDecoder(bytes.NewBuffer(resData)).Decode(&callRes)
	if err != nil {
		return fmt.Errorf("cannot parse response: %w", err)
	}
	return nil
}

func (g *GetBlockClient) GetLastBlockNumber(ctx context.Context) (big.Int, error) {
	var blockNumRes string
	err := g.performRPCCall(ctx, "eth_blockNumber", nil, &blockNumRes)
	if err != nil {
		return big.Int{}, fmt.Errorf("cannot perform rpc call eth_blockNumber: %w", err)
	}
	if len(blockNumRes) == 0 {
		return big.Int{}, errors.New("got empty block number")
	}
	blockNum, ok := new(big.Int).SetString(blockNumRes[2:], 16)
	if !ok {
		return big.Int{}, fmt.Errorf("cannot parse block number: %w", err)
	}
	fmt.Printf("last block number: %s, hex: %s\n", blockNum, blockNumRes)

	return *blockNum, nil
}

func (g *GetBlockClient) GetBlockRecord(ctx context.Context, blockNumber big.Int) (adapters.BlockRecord, error) {
	blockRec := adapters.BlockRecord{}
	params := []interface{}{
		fmt.Sprintf("0x%x", &blockNumber),
		true,
	}
	err := g.performRPCCall(ctx, "eth_getBlockByNumber", params, &blockRec)
	if err != nil {
		return adapters.BlockRecord{}, fmt.Errorf("cannot perform rpc call eth_getBlockByNumber: %w", err)
	}
	return blockRec, nil
}
