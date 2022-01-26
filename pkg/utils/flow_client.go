package utils

import (
	"InceptionAnimals/app/models"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

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

func ServiceAccount(flowClient *client.Client) (flow.Address, *flow.AccountKey, crypto.Signer) {
	privateKey, _ := crypto.DecodePrivateKeyHex(crypto.ECDSA_P256, os.Getenv("MINTER_PRIVATE_KEY"))

	addr := flow.HexToAddress(os.Getenv("MINTER_ADDRESS"))
	acc, _ := flowClient.GetAccount(context.Background(), addr)

	// minterKeyIndex := int(os.Getenv("MINTER_KEY_INDEX"))
	accountKey := acc.Keys[0]
	signer := crypto.NewInMemorySigner(privateKey, accountKey.HashAlgo)

	return addr, accountKey, signer
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

func first(n cadence.Value, _ error) cadence.Value {
	return n
}

func CreateNftTemplate(template models.NFTTemplate) (*flow.TransactionResult, error) {

	templateIsMinted, err := CheckIfTemplateIsMinted(uint64(template.TypeID))
	if err != nil {
		return nil, err
	}
	if templateIsMinted == cadence.NewBool(true) {
		return nil, errors.New("template_is_minted")
	}

	args := []cadence.Value{
		cadence.NewUInt64(uint64(template.TypeID)),
		cadence.NewBool(template.IsPack),
		first(cadence.NewString(template.Name)),
		first(cadence.NewString(template.Description)),
		cadence.NewUInt64(uint64(template.MintLimit)),
		first(cadence.NewUFix64(template.PriceUSD)),
		first(cadence.NewUFix64(template.PriceFlow)),
		cadence.NewDictionary([]cadence.KeyValuePair{
			{Key: first(cadence.NewString("uri")), Value: first(cadence.NewString(template.Metadata.Uri))},
			{Key: first(cadence.NewString("mimetype")), Value: first(cadence.NewString(template.Metadata.Mimetype))},
			{Key: first(cadence.NewString("quality")), Value: first(cadence.NewString(template.Metadata.Quality))},
		}),
		cadence.NewDictionary([]cadence.KeyValuePair{
			{Key: first(cadence.NewString("availableAt")), Value: first(cadence.NewUFix64(fmt.Sprintf("%d", template.Timestamps.AvailableAt.Unix())))},
			{Key: first(cadence.NewString("expiresAt")), Value: first(cadence.NewUFix64(fmt.Sprintf("%d", template.Timestamps.ExpiresAt.Unix())))},
		}),
		cadence.NewBool(template.IsLand),
	}

	txResult, err := sendTransaction("create_nft_template", args)
	if err != nil {
		return nil, err
	}

	return txResult, nil
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

func CheckIfAddressHasCollection(address string) (cadence.Value, error) {
	scriptName := "check_address_has_collection"
	args := []cadence.Value{cadence.NewAddress(flow.HexToAddress(address))}

	scriptRes, err := ExecuteScript(scriptName, args)
	if err != nil {
		return nil, err
	}

	return scriptRes, nil
}

func MintNFT(typeID int, address string) (*flow.TransactionResult, error) {
	addressHasCollection, err := CheckIfAddressHasCollection(address)
	if err != nil {
		return nil, err
	}
	if addressHasCollection == cadence.NewBool(false) {
		return nil, errors.New("missing_collection")
	}

	args := []cadence.Value{
		cadence.NewAddress(flow.HexToAddress(address)),
		cadence.NewUInt64(uint64(typeID)),
	}

	txResult, err := sendTransaction("mint_crypto_zoo_nft", args)
	if err != nil {
		return nil, err
	}

	return txResult, nil
}

// WaitForSeal wait fot the process to seal
func WaitForSeal(ctx context.Context, c *client.Client, id flow.Identifier) (*flow.TransactionResult, error) {
	result, err := c.GetTransactionResult(ctx, id)
	if err != nil {
		return nil, err
	}

	for result.Status != flow.TransactionStatusSealed {
		time.Sleep(time.Second)
		result, err = c.GetTransactionResult(ctx, id)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}
