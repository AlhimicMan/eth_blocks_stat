package adapters

type TransactionRec struct {
	BlockNumber      string `json:"blockNumber"`
	From             string `json:"from"`
	To               string `json:"to"`
	Input            string `json:"input"`
	TransactionIndex string `json:"transactionIndex"`
	Hash             string `json:"hash"`
}

type BlockRecord struct {
	Number       string           `json:"number"`
	Transactions []TransactionRec `json:"transactions"`
}
