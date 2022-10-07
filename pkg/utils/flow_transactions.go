package utils

import (
	"context"
	"errors"
	"time"

	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/client"
)

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
