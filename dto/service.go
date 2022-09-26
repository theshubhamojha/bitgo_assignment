package dto

type TransactionDetails struct {
	ID                string                    `json:"txid"`
	Size              int32                     `json:"size"`
	InputTransactions []InputTransactionsHolder `json:"vin"`
	Status            StatusHolder              `json:"status"`
}

type InputTransactionsHolder struct {
	ID string `json:"txid"`
}

type StatusHolder struct {
	BlockHeight int32 `json:"block_height"`
}

type TransactionStatus struct {
	BlockHeight int32 `json:"block_height"`
}

type AncestoralCount struct {
	ID    string
	Count int32
}
