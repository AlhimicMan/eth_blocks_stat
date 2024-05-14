package adapters

//type TransactionRec struct {
//}

type TransactionRec struct {
	BlockNumber      string `json:"blockNumber"`
	From             string `json:"from"`
	To               string `json:"to"`
	Input            string `json:"input"`
	TransactionIndex string `json:"transactionIndex"`
}

type BlockRecord struct {
	Number       string           `json:"number"`
	Transactions []TransactionRec `json:"transactions"`
}
