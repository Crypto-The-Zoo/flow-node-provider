package utils

import (
	"InceptionAnimals/app/models"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/onflow/cadence"
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

func MutateScriptAddress(script string) string {
	scriptStr := strings.ReplaceAll(script, "../../contracts/NonFungibleToken.cdc", (os.Getenv("NON_FUNGIBLE_TOKEN_ADDRESS")))
	scriptStr = strings.ReplaceAll(scriptStr, "../../contracts/CryptoZooNFT.cdc", os.Getenv("CRYPTO_ZOO_NFT_ADDRESS"))
	scriptStr = strings.ReplaceAll(scriptStr, "\"", "")

	return scriptStr
}

func ExecuteScript(scriptName string, args []cadence.Value) (cadence.Value, error) {
	ctx := context.Background()

	flowClient, err := ConnectToFlowAccessAPI()
	if err != nil {
		return nil, err
	}

	script, err := ioutil.ReadFile("cadence/scripts/CryptoZoo/" + scriptName + ".cdc")
	if err != nil {
		return nil, err
	}

	scriptStr := MutateScriptAddress(string(script))

	println(scriptStr)

	value, err := flowClient.ExecuteScriptAtLatestBlock(ctx, []byte(scriptStr), args)
	if err != nil {
		return nil, err
	}

	return value, nil
}

func CheckIfTemplateIsMinted(typeID uint64) (cadence.Value, error) {
	scriptName := "check_is_template_minted"
	args := []cadence.Value{cadence.NewUInt64(typeID)}

	scriptRes, err := ExecuteScript(scriptName, args)
	if err != nil {
		return nil, err
	}

	return scriptRes, nil
}
