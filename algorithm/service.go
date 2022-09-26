package algorithm

import (
	"context"
	"errors"
	"log"
	"sort"

	"github.com/shubhamiitbhu/bitgo.git/dto"
	"github.com/shubhamiitbhu/bitgo.git/transactions"
)

type service struct {
	transactionService transactions.Service
}

type Service interface {
	BFS(ctx context.Context, transactionID string, transactionBlockHeight int32) (list []string, err error)
	GetHighestAncestorSetTransactions(ctx context.Context, ancestoralCount []dto.AncestoralCount, top int32) (transactions []dto.AncestoralCount)
}

func NewAlgorithmService(transactionService transactions.Service) Service {
	return &service{
		transactionService: transactionService,
	}
}

// BFS algorithm is used to capture the ancestors
// It will maintain a map of visited transactions in the isVisitedTransaction variable
// if the transaction is already visited we will igonore them
// if the transaction is not visited, we will calculate all the inputs and will put those enqueue those transactions if they belong to the same block_height
// at the end, we will pop the front element from the queue by calling the Dequeue util function
func (s *service) BFS(ctx context.Context, transactionID string, transactionBlockHeight int32) (list []string, err error) {
	var queue []string
	var isVisitedTransaction = make(map[string]bool)
	var ancestorList []string

	queue = Enqueue(queue, transactionID)

	for len(queue) != 0 {
		frontElement := queue[0]
		if isVisitedTransaction[frontElement] {
			continue
		}
		isVisitedTransaction[frontElement] = true
		inputTransactions, err := s.transactionService.GetInputTransactions(ctx, frontElement)
		if err != nil {
			log.Print("error getting all input transactions for running algorithm", "error", err.Error(), "transaction_id", transactionID)
			return list, err
		}

		for _, transaction := range inputTransactions {
			blockHeight, err := s.transactionService.GetTransactionBlockHeight(ctx, transaction)
			if err != nil {
				log.Print("error getting transaction block height", "error", err.Error(), "transaction_id", transaction)
				return ancestorList, errors.New("error getting transaction block height")
			}

			if transactionBlockHeight == blockHeight {
				ancestorList = append(ancestorList, transaction)
				queue = Enqueue(queue, transaction)
			}
		}

		queue = Dequeue(queue)
	}

	list = ancestorList
	return
}

// GetHighestAncestorSetTransactions will sort all the transaction with > 0 ancestors to figure out the top n transactions
func (s *service) GetHighestAncestorSetTransactions(ctx context.Context, ancestoralCount []dto.AncestoralCount, top int32) []dto.AncestoralCount {
	sort.Slice(ancestoralCount, func(i, j int) bool {
		return ancestoralCount[i].Count > ancestoralCount[j].Count
	})

	// returning if size of array is lesser than the `top` parameter provided
	if top > int32(len(ancestoralCount)) {
		return ancestoralCount
	}

	return ancestoralCount[:top]
}
