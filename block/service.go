package block

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
)

type service struct {
}

type Service interface {
	GetBlockHash(ctx context.Context, blockHeight int32) (blockHash string, err error)
}

func NewBlockService() Service {
	return &service{}
}

func (s *service) GetBlockHash(ctx context.Context, blockHeight int32) (blockHash string, err error) {
	url := fmt.Sprintf("https://blockstream.info/api/block-height/%d", blockHeight)

	response, err := http.Get(url)
	if err != nil {
		log.Print("error getting block hash from api", "error", err.Error(), "block_height", blockHeight)
		return
	}
	if response == nil {
		log.Print("error getting proper response for block hash", "url", url, "block_height", blockHeight)
		return
	}

	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		log.Print("error reading response body block height", "error", err.Error(), "block_height", blockHeight)
		return
	}

	blockHash = string(responseBody)
	return
}
