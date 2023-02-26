package utils

import (
	"InceptionAnimals/app/models"
	"context"
	"math"
	"os"
	"strings"

	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/client"
	"github.com/onflow/flow-go-sdk/crypto"
	"google.golang.org/grpc"
)

// func ConnectToFlowAccessAPI() (*client.Client, error) {
// 	flowAccessAPI := os.Getenv("FLOW_ACCESS_NODE")

// 	// 40 MB
// 	flow, err := client.New(flowAccessAPI, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*1024)), grpc.WithInsecure())

// 	if err != nil {
// 		return nil, err
// 	}

// 	return flow, nil
// }

func ConnectToFlowAccessAPIWithNode(node string) (*client.Client, error) {

	flow, err := client.New(node, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*1024)), grpc.WithInsecure())

	if err != nil {
		return nil, err
	}

	return flow, nil
}

func GetReferenceBlockId(flowClient *client.Client) flow.Identifier {
	block, _ := flowClient.GetLatestBlock(context.Background(), true)

	return block.ID
}

// func GetLatestBlock() (*models.Block, error) {
// 	currentBlock := &models.Block{}

// 	ctx := context.Background()

// 	flowClient, err := ConnectToFlowAccessAPI()
// 	if err != nil {
// 		return currentBlock, err
// 	}

// 	// get the latest sealed block
// 	isSealed := true
// 	latestBlock, err := flowClient.GetLatestBlock(ctx, isSealed)
// 	if err != nil {
// 		return currentBlock, err
// 	}

// 	currentBlock.ID = fmt.Sprintf("ID: %s", latestBlock.ID)
// 	currentBlock.Height = fmt.Sprintf("%d", latestBlock.Height)
// 	currentBlock.Timestamp = latestBlock.Timestamp

// 	return currentBlock, nil
// }

func GetEventsInBlockHeightRangeAutoNode(eventType string, startHeight uint64, endHeight uint64) ([]models.BlockEvents, error) {
	current_env := os.Getenv("ENV")

	if current_env != "prod" {
		return GetEventsInBlockHeightRange("access.devnet.nodes.onflow.org:9000", eventType, startHeight, endHeight)
	}

	nodeMap := map[string]map[string]interface{}{
		"6": {
			"access":           "access-001.mainnet6.nodes.onflow.org:9000",
			"startBlockHeight": uint64(12609237),
			"endBlockHeight":   uint64(13404174 - 1),
		},
		"7": {
			"access":           "access-001.mainnet7.nodes.onflow.org:9000",
			"startBlockHeight": uint64(13404174),
			"endBlockHeight":   uint64(13950742 - 1),
		},
		"8": {
			"access":           "access-001.mainnet8.nodes.onflow.org:9000",
			"startBlockHeight": uint64(13950742),
			"endBlockHeight":   uint64(14892104 - 1),
		},
		"9": {
			"access":           "access-001.mainnet9.nodes.onflow.org:9000",
			"startBlockHeight": uint64(14892104),
			"endBlockHeight":   uint64(15791891 - 1),
		},
		"10": {
			"access":           "access-001.mainnet10.nodes.onflow.org:9000",
			"startBlockHeight": uint64(15791891),
			"endBlockHeight":   uint64(16755602 - 1),
		},
		"11": {
			"access":           "access-001.mainnet11.nodes.onflow.org:9000",
			"startBlockHeight": uint64(16755602),
			"endBlockHeight":   uint64(17544523 - 1),
		},
		"12": {
			"access":           "access-001.mainnet12.nodes.onflow.org:9000",
			"startBlockHeight": uint64(17544523),
			"endBlockHeight":   uint64(18587478 - 1),
		},
		"13": {
			"access":           "access-001.mainnet13.nodes.onflow.org:9000",
			"startBlockHeight": uint64(18587478),
			"endBlockHeight":   uint64(19050753 - 1),
		},
		"14": {
			"access":           "access-001.mainnet14.nodes.onflow.org:9000",
			"startBlockHeight": uint64(19050753),
			"endBlockHeight":   uint64(21291692 - 1),
		},
		"15": {
			"access":           "access-001.mainnet15.nodes.onflow.org:9000",
			"startBlockHeight": uint64(21291692),
			"endBlockHeight":   uint64(23830813 - 1),
		},
		"16": {
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
			"access":           "access-001.mainnet19.nodes.onflow.org:9000",
			"startBlockHeight": uint64(35858811),
			"endBlockHeight":   uint64(40171634 - 1),
		},
		"20": {
			"access":           "access-001.mainnet20.nodes.onflow.org:9000",
			"startBlockHeight": uint64(40171634),
			"endBlockHeight":   uint64(44950207 - 1),
		},
		"21": {
			"access":           "access-001.mainnet21.nodes.onflow.org:9000",
			"startBlockHeight": uint64(44950207),
			"endBlockHeight":   uint64(47169687 - 1),
		},
		"22": {
			"access":           "access.mainnet.nodes.onflow.org:9000",
			"startBlockHeight": uint64(47169687),
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
	if height >= 47169687 {
		return "22"
	}
	if height >= 44950207 {
		return "21"
	}
	if height >= 40171634 {
		return "20"
	}
	if height >= 35858811 {
		return "19"
	}
	if height >= 31735955 {
		return "18"
	}
	if height >= 27341470 {
		return "17"
	}
	if height >= 23830813 {
		return "16"
	}
	if height >= 21291692 {
		return "15"
	}
	if height >= 19050753 {
		return "14"
	}
	if height >= 18587478 {
		return "13"
	}
	if height >= 17544523 {
		return "12"
	}
	if height >= 16755602 {
		return "11"
	}
	if height >= 15791891 {
		return "10"
	}
	if height >= 14892104 {
		return "9"
	}
	if height >= 13950742 {
		return "8"
	}
	if height >= 13404174 {
		return "7"
	}
	if height >= 12609237 {
		return "6"
	}
	return "UNKNOWN"
}

func GetEventsInBlockHeightRangeRaw(node string, eventType string, startHeight uint64, endHeight uint64) ([]client.BlockEvents, error) {
	blocks := []client.BlockEvents{}

	ctx := context.Background()

	flowClient, err := ConnectToFlowAccessAPIWithNode(node)
	if err != nil {
		return blocks, err
	}

	blocks, err = flowClient.GetEventsForHeightRange(ctx, client.EventRangeQuery{
		Type:        eventType,
		StartHeight: startHeight,
		EndHeight:   endHeight,
	})

	if err != nil {
		panic("failed to query events")
	}

	return blocks, nil
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
		panic(err)
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
	signer, err := crypto.NewInMemorySigner(privateKey, accountKey.HashAlgo)
	if err != nil {
		panic("NewInMemorySigner Error!")
	}

	return addr, accountKey, signer
}
