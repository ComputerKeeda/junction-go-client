package core

import (
	"context"
	"fmt"
	"github.com/ComputerKeeda/junction-go-client/chain"
	"github.com/ComputerKeeda/junction-go-client/components"
	"github.com/ComputerKeeda/junction-go-client/types"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosaccount"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosclient"
)

/*

type MsgValidateVrf struct {
	Creator      string json:"creator,omitempty"`
	StationId    string json:"stationId,omitempty"`
	PodNumber    uint64 json:"podNumber,omitempty"`
	SerializedRc []byte json:"serializedRc,omitempty"`
}

*/

func ValidateVRF(serializedRc []byte) bool {
	// getting the account and creating client codes --> Start
	accountName := "temp-account"
	accountPath := "./temp-account"
	addressPrefix := "air"
	registry, err := cosmosaccount.New(cosmosaccount.WithHome(accountPath))
	if err != nil {
		components.Logger.Error(fmt.Sprintf("Error creating account registry: %v", err))
		return false
	}

	newTempAccount, err := registry.GetByName(accountName)
	if err != nil {
		components.Logger.Error(fmt.Sprintf("Error getting account: %v", err))
		return false
	}

	newTempAddr, err := newTempAccount.Address(addressPrefix)
	if err != nil {
		components.Logger.Error(fmt.Sprintf("Error getting address: %v", err))
		return false
	}

	ctx := context.Background()
	gas := components.GenerateRandomWithFavour(510, 1000, [2]int{520, 700}, 0.7)
	gasFees := fmt.Sprintf("%damf", gas)
	components.Logger.Warn(fmt.Sprintf("Gas Fees Used for validate VRF transaction is: %s\n", gasFees))
	accountClient, err := cosmosclient.New(ctx, cosmosclient.WithAddressPrefix(addressPrefix), cosmosclient.WithNodeAddress(components.JunctionTTCRPC), cosmosclient.WithHome(accountPath), cosmosclient.WithGas("auto"), cosmosclient.WithFees(gasFees))
	if err != nil {
		components.Logger.Error("Error creating account client")
		return false
	}
	// getting the account and creating client codes --> End

	stationId, err := chain.GetStationId()
	if err != nil {
		components.Logger.Error(err.Error())
		return false
	}
	podNumber, err := chain.GetPodNumber()
	if err != nil {
		components.Logger.Error(err.Error())
		return false
	}

	msg := types.MsgValidateVrf{
		Creator:      newTempAddr,
		StationId:    stationId,
		PodNumber:    podNumber,
		SerializedRc: serializedRc,
	}

	txRes, errTxRes := accountClient.BroadcastTx(ctx, newTempAccount, &msg)
	if errTxRes != nil {
		components.Logger.Error("error in transaction" + errTxRes.Error())
		return false
	}

	components.Logger.Info("Transaction Hash For VRF Validation: " + txRes.TxHash)

	return true
}