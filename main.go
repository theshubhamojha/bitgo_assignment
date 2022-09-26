package main

import (
	"context"
	"log"
	"math"
	"sync"

	"github.com/shubhamiitbhu/bitgo.git/algorithm"
	"github.com/shubhamiitbhu/bitgo.git/block"
	"github.com/shubhamiitbhu/bitgo.git/dto"
	"github.com/shubhamiitbhu/bitgo.git/transactions"
)

func main() {

	ctx := context.Background()
	blockHeight := int32(680000)

	blockService := block.NewBlockService()
	blockHash, err := blockService.GetBlockHash(ctx, blockHeight)
	if err != nil {
		log.Print("error getting block hash")
		return
	}

	transactionsService := transactions.NewTransactionService()
	algorithmService := algorithm.NewAlgorithmService(transactionsService)

	allTx, err := transactionsService.GetAllTransaction(ctx, blockHash)
	if err != nil {
		log.Print("error getting all transactions", "error", err.Error())
		return
	}

	var ancestoralCount []dto.AncestoralCount
	var wg sync.WaitGroup

	batchSize := 400

	startIndex := 0
	batchNumber := 0

	for startIndex < len(allTx) {
		maxAvailableIndex := int(math.Min(float64(batchSize), float64(len(allTx)-batchNumber*batchSize-1)))
		for index := batchNumber * batchSize; index < batchNumber*batchSize+maxAvailableIndex; index++ {
			// we will evaluate the results concurrently in definite batch size
			// this will help us to run the calculations quickly
			// we are adding waitgroups to wait for the responses from a batch
			// we are doing the math to pick up minimum of the batchSize or difference of length and current index to evaluate the number of waitgroups
			// this is done to for the last batch, since we cannot gaurantee it's size to be equal to the batch size
			wg.Add(1)
			position := index
			go func() {
				defer wg.Done()
				ancestorSet, err := algorithmService.BFS(context.TODO(), allTx[position], blockHeight)
				if err != nil {
					log.Print(" error getting ancestor set ", " error ", err.Error(), " transaction_id ", allTx[position])
					return
				}

				// append only transactions with > 0 ancestors. This will be used to figure out top 10 transactions with highest ancestor counts
				if len(ancestorSet) > 0 {
					ancestoralCount = append(ancestoralCount, dto.AncestoralCount{
						ID:    allTx[position],
						Count: int32(len(ancestorSet)),
					})
				}
			}()
		}

		wg.Wait()
		startIndex = (batchNumber + 1) * batchSize
		batchNumber += 1
	}

	topTransactions := algorithmService.GetHighestAncestorSetTransactions(ctx, ancestoralCount, 10)

	log.Print("successfully evaluated the transactions with highest ancestoral set")

	for _, transaction := range topTransactions {
		log.Print("id ", transaction.ID, " count ", transaction.Count)
	}
}
