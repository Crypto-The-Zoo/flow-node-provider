package utils

import (
	"InceptionAnimals/app/models"
	"context"
	"fmt"
	"os"

	"github.com/onflow/flow-go-sdk/client"
	"google.golang.org/grpc"
)

func ConnectToFlowAccessAPI() (*client.Client, error) {
	flowAccessAPI := os.Getenv("FLOW_ACCESS_NODE")

	flow, err := client.New(flowAccessAPI, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return flow, nil
}

func GetLatestBlock() (*models.Block, error) {
	currentBlock := &models.Block{}

	ctx := context.Background()

	flowClient, err := ConnectToFlowAccessAPI()
	if err != nil {
		return currentBlock, err
	}

	// get the latest sealed block
	isSealed := true
	latestBlock, err := flowClient.GetLatestBlock(ctx, isSealed)
	if err != nil {
		return currentBlock, err
	}

	currentBlock.ID = fmt.Sprintf("ID: %s", latestBlock.ID)
	currentBlock.Height = fmt.Sprintf("%d", latestBlock.Height)
	currentBlock.Timestamp = latestBlock.Timestamp

	return currentBlock, nil
}
