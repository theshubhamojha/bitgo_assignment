package transactions

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/shubhamiitbhu/bitgo.git/dto"
)

type service struct {
}

type Service interface {
	GetInputTransactions(ctx context.Context, transactionID string) (inputTransactions []string, err error)
	GetTransactionBlockHeight(ctx context.Context, transactionID string) (blockHeight int32, err error)
	GetAllTransaction(ctx context.Context, blockHash string) (transactionIDs []string, err error)
}

func NewTransactionService() Service {
	return &service{}
}

func (s *service) GetInputTransactions(ctx context.Context, id string) (inputTransactions []string, err error) {
	url := fmt.Sprintf("https://blockstream.info/api/tx/%s", id)

	response, err := http.Get(url)
	if err != nil {
		log.Print("error getting transaction details", "error", err.Error(), "url", url, "id", id)
		return
	}
	if response == nil {
		log.Print("error getting proper response for transaction", "url", url, "id", id)
		return
	}

	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		log.Print("error reading response body transaction", "error", err.Error(), "id", id)
		return
	}

	var transaction dto.TransactionDetails
	err = json.Unmarshal(responseBody, &transaction)
	if err != nil {
		log.Print("error unmarshalling response for transaction", "error", err.Error(), "id", id, "response_body", string(responseBody))
		return
	}

	for _, transaction := range transaction.InputTransactions {
		inputTransactions = append(inputTransactions, transaction.ID)
	}
	return inputTransactions, nil
}

func (s *service) GetAllTransaction(ctx context.Context, blockHash string) (transactionIDs []string, err error) {
	url := fmt.Sprintf("https://blockstream.info/api/block/%s/txids", blockHash)

	response, err := http.Get(url)
	if err != nil {
		log.Print("error getting all transaction details", "error", err.Error(), "url", url, "block_hash", blockHash)
		return
	}
	if response == nil {
		log.Print("error getting proper response for all transaction", "url", url, "block_hash", blockHash)
		return
	}

	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		log.Print("error reading response body for all transaction", "error", err.Error(), "block_hash", blockHash)
		return
	}

	err = json.Unmarshal(responseBody, &transactionIDs)
	if err != nil {
		log.Print("error unmarshalling response body for all transactions", "error", err.Error(), "block_hash", blockHash)
		return
	}

	return transactionIDs, nil
}

func (s *service) GetTransactionBlockHeight(ctx context.Context, transactionID string) (blockHeight int32, err error) {
	url := fmt.Sprintf("https://blockstream.info/api/tx/%s/status", transactionID)

	response, err := http.Get(url)
	if err != nil {
		log.Print("error getting transaction block height details", "error", err.Error(), "url", url, "transaction_id", transactionID)
		return
	}
	if response == nil {
		log.Print("error getting proper response for all transaction", "url", url, "transaction_id", transactionID)
		return
	}

	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)

	var transactionStatus dto.TransactionStatus
	err = json.Unmarshal(responseBody, &transactionStatus)
	if err != nil {
		log.Print("error unmarshalling transaction status details", "error", err.Error(), "transaction_id", transactionID)
		return
	}

	blockHeight = transactionStatus.BlockHeight
	return blockHeight, nil
}
