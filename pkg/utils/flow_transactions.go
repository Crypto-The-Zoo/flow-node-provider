package utils

import (
	"InceptionAnimals/app/models"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/client"
)

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
		First(cadence.NewString(template.Name)),
		First(cadence.NewString(template.Description)),
		cadence.NewUInt64(uint64(template.MintLimit)),
		First(cadence.NewUFix64(template.PriceUSD)),
		First(cadence.NewUFix64(template.PriceFlow)),
		cadence.NewDictionary([]cadence.KeyValuePair{
			{Key: First(cadence.NewString("uri")), Value: First(cadence.NewString(template.Metadata.Uri))},
			{Key: First(cadence.NewString("mimetype")), Value: First(cadence.NewString(template.Metadata.Mimetype))},
			{Key: First(cadence.NewString("quality")), Value: First(cadence.NewString(template.Metadata.Quality))},
		}),
		cadence.NewDictionary([]cadence.KeyValuePair{
			{Key: First(cadence.NewString("availableAt")), Value: First(cadence.NewUFix64(fmt.Sprintf("%d", template.Timestamps.AvailableAt.Unix())))},
			{Key: First(cadence.NewString("expiresAt")), Value: First(cadence.NewUFix64(fmt.Sprintf("%d", template.Timestamps.ExpiresAt.Unix())))},
		}),
		cadence.NewBool(template.IsLand),
	}

	txResult, err := sendTransaction("create_nft_template", args)
	if err != nil {
		return nil, err
	}

	return txResult, nil
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
