package utils

import (
	"InceptionAnimals/app/models"
	"context"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strings"

	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/client"
	"github.com/onflow/flow-go-sdk/crypto"
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

func ConnectToFlowAccessAPIWithNode(node string) (*client.Client, error) {
	flow, err := client.New(node, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return flow, nil
}

func GetReferenceBlockId(flowClient *client.Client) flow.Identifier {
	block, _ := flowClient.GetLatestBlock(context.Background(), true)

	return block.ID
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

func GetEventsInBlockHeightRangeAutoNode(eventType string, startHeight uint64, endHeight uint64) ([]models.BlockEvents, error) {
	nodeMap := map[string]map[string]interface{}{"16": {
		"access":           "access-001.mainnet16.nodes.onflow.org:9000",
		"startBlockHeight": uint64(23830813),
		"endBlockHeight":   uint64(27341470 - 1),
	},
		"17": {
			"access":           "access-001.mainnet17.nodes.onflow.org:9000",
			"startBlockHeight": uint64(27341470),
			"endBlockHeight":   uint64(31735955 - 1),
		},
		"18": {
			"access":           "access-001.mainnet18.nodes.onflow.org:9000",
			"startBlockHeight": uint64(31735955),
			"endBlockHeight":   uint64(35858811 - 1),
		},
		"19": {
			"access":           "access.mainnet.nodes.onflow.org:9000",
			"startBlockHeight": uint64(35858811),
			"endBlockHeight":   math.Inf(1),
		},
	}

	startNode := nodeBelongsTo(startHeight)
	endNode := nodeBelongsTo(endHeight)

	if startNode == endNode {
		return GetEventsInBlockHeightRange(nodeMap[startNode]["access"].(string), eventType, startHeight, endHeight)
	} else {
		firstPart, _ := GetEventsInBlockHeightRange(nodeMap[startNode]["access"].(string), eventType, startHeight, nodeMap[startNode]["endBlockHeight"].(uint64))
		secondPart, _ := GetEventsInBlockHeightRange(nodeMap[endNode]["access"].(string), eventType, nodeMap[endNode]["startBlockHeight"].(uint64), endHeight)

		firstPart = append(firstPart, secondPart...)

		return firstPart, nil
	}
}

func nodeBelongsTo(height uint64) string {
	if height >= 23830813 && height <= 27341470-1 {
		return "16"
	} else if height >= 27341470 && height <= 31735955-1 {
		return "17"
	} else if height >= 31735955 && height <= 35858811-1 {
		return "18"
	}

	return "19"
}

func GetEventsInBlockHeightRange(node string, eventType string, startHeight uint64, endHeight uint64) ([]models.BlockEvents, error) {
	blocks := []client.BlockEvents{}
	blockEvents := []models.BlockEvents{}

	ctx := context.Background()

	flowClient, err := ConnectToFlowAccessAPIWithNode(node)
	if err != nil {
		return blockEvents, err
	}

	blocks, err = flowClient.GetEventsForHeightRange(ctx, client.EventRangeQuery{
		Type:        eventType,
		StartHeight: startHeight,
		EndHeight:   endHeight,
	})

	if err != nil {
		panic("failed to query events")
	}

	// TODO: serialize the payload
	for _, s := range blocks {
		for _, event := range s.Events {

			// Prepare blockEventData
			blockEventData := make(map[string]string)
			eventFields := event.Value.EventType.Fields
			eventValues := event.Value.Fields
			for i, field := range eventFields {
				blockEventData[field.Identifier] = eventValues[i].String()
			}

			blockEvents = append(blockEvents, models.BlockEvents{
				ID:                s.BlockID.Hex(),
				FlowTransactionID: event.TransactionID.Hex(),
				BlockEventData:    blockEventData,
				EventDate:         s.BlockTimestamp,

				BlockID:          s.BlockID.Hex(),
				BlockHeight:      s.Height,
				BlockTimestamp:   s.BlockTimestamp,
				Type:             event.Type,
				TransactionID:    event.TransactionID.Hex(),
				TransactionIndex: event.TransactionIndex,
				EventIndex:       event.EventIndex,
				Data:             blockEventData,
			})
		}
	}

	return blockEvents, nil
}

func MutateScriptAddress(script string) string {
	scriptStr := strings.ReplaceAll(script, "../../contracts/NonFungibleToken.cdc", (os.Getenv("NON_FUNGIBLE_TOKEN_ADDRESS")))
	scriptStr = strings.ReplaceAll(scriptStr, "../../contracts/CryptoZooNFT.cdc", os.Getenv("CRYPTO_ZOO_NFT_ADDRESS"))
	scriptStr = strings.ReplaceAll(scriptStr, "\"", "")

	return scriptStr
}

func First(n cadence.Value, _ error) cadence.Value {
	return n
}

func ServiceAccount(flowClient *client.Client) (flow.Address, *flow.AccountKey, crypto.Signer) {
	privateKey, _ := crypto.DecodePrivateKeyHex(crypto.ECDSA_P256, os.Getenv("MINTER_PRIVATE_KEY"))

	addr := flow.HexToAddress(os.Getenv("MINTER_ADDRESS"))
	acc, _ := flowClient.GetAccount(context.Background(), addr)

	// minterKeyIndex := int(os.Getenv("MINTER_KEY_INDEX"))
	accountKey := acc.Keys[0]
	signer := crypto.NewInMemorySigner(privateKey, accountKey.HashAlgo)

	return addr, accountKey, signer
}

func ExecuteScript(scriptName string, args []cadence.Value) (cadence.Value, error) {
	ctx := context.Background()

	flowClient, err := ConnectToFlowAccessAPI()
	if err != nil {
		return nil, err
	}
	defer func() {
		err := flowClient.Close()
		if err != nil {
			panic(err)
		}
	}()

	script, err := ioutil.ReadFile("cadence/scripts/CryptoZoo/" + scriptName + ".cdc")
	if err != nil {
		return nil, err
	}

	scriptStr := MutateScriptAddress(string(script))

	value, err := flowClient.ExecuteScriptAtLatestBlock(ctx, []byte(scriptStr), args)
	if err != nil {
		return nil, err
	}

	return value, nil
}

func sendTransaction(scriptName string, args []cadence.Value) (*flow.TransactionResult, error) {
	ctx := context.Background()

	// Initialize Flow Client
	flowClient, err := ConnectToFlowAccessAPI()
	if err != nil {
		return nil, err
	}
	defer func() {
		err := flowClient.Close()
		if err != nil {
			panic(err)
		}
	}()

	// Parse transaction script
	script, err := ioutil.ReadFile("cadence/transactions/CryptoZoo/" + scriptName + ".cdc")
	if err != nil {
		return nil, err
	}
	scriptStr := MutateScriptAddress(string(script))

	// Get service account
	serviceAcctAddr, serviceAcctKey, serviceSigner := ServiceAccount(flowClient)

	// Build and sign transaction
	tx := flow.NewTransaction().
		SetPayer(serviceAcctAddr).
		SetProposalKey(serviceAcctAddr, serviceAcctKey.Index, serviceAcctKey.SequenceNumber).
		SetScript([]byte(scriptStr)).
		SetReferenceBlockID(GetReferenceBlockId(flowClient)).
		AddAuthorizer(serviceAcctAddr)

	for _, argument := range args {
		tx.AddArgument(argument)
	}

	err = tx.SignEnvelope(serviceAcctAddr, serviceAcctKey.Index, serviceSigner)
	if err != nil {
		return nil, err
	}

	// Send transaction
	err = flowClient.SendTransaction(ctx, *tx)
	if err != nil {
		return nil, err
	}

	fmt.Printf("--Transaction ID: %s", tx.ID())

	txResult, err := WaitForSeal(ctx, flowClient, tx.ID())
	if err != nil {
		return nil, err
	}

	return txResult, nil
}
