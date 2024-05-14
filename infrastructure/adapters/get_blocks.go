package adapters

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"eth_blocks_stat/core/adapters"
	"fmt"
	"net/http"
	"reflect"
	"time"
)

const (
	rpcVersion = "2.0"
	defaultId  = "getblock.io"
	rateLimit  = time.Second / 50 // 50 RPS
)

type GetBlockClient struct {
	apiKey   string
	apiURL   string
	throttle <-chan time.Time
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
	throttle := time.Tick(rateLimit)
	return &GetBlockClient{
		apiKey:   apiKey,
		apiURL:   apiURL,
		throttle: throttle,
	}
}

func (g *GetBlockClient) performRPCCall(ctx context.Context, method string, params []interface{}, result interface{}) error {
	<-g.throttle
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
	err = json.NewDecoder(res.Body).Decode(&callRes)
	if err != nil {
		return fmt.Errorf("cannot parse response: %w", err)
	}
	return nil
}

func (g *GetBlockClient) GetLastBlockNumber(ctx context.Context) (string, error) {
	var blockNumRes string
	err := g.performRPCCall(ctx, "eth_blockNumber", nil, &blockNumRes)
	if err != nil {
		return "", fmt.Errorf("cannot perform rpc call eth_blockNumber: %w", err)
	}
	if len(blockNumRes) < 3 {
		return "", errors.New("got empty block number")
	}

	return blockNumRes[2:], nil
}

func (g *GetBlockClient) GetBlockRecord(ctx context.Context, blockNumber string) (adapters.BlockRecord, error) {
	blockRec := adapters.BlockRecord{}
	params := []interface{}{
		blockNumber,
		true,
	}
	err := g.performRPCCall(ctx, "eth_getBlockByNumber", params, &blockRec)
	if err != nil {
		return adapters.BlockRecord{}, fmt.Errorf("cannot perform rpc call eth_getBlockByNumber: %w", err)
	}
	return blockRec, nil
}
