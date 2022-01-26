package utils

import (
	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk"
)

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
