package core

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ComputerKeeda/junction-go-client/chain"
	"github.com/ComputerKeeda/junction-go-client/components"
	"github.com/ComputerKeeda/junction-go-client/types"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosaccount"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosclient"
	"time"
)

func SubmitPod(merkleRootHash string, publicWitness any) bool {
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
	gas := components.GenerateRandomWithFavour(100, 300, [2]int{120, 250}, 0.7)
	gasFees := fmt.Sprintf("%damf", gas)
	components.Logger.Warn(fmt.Sprintf("Gas Fees Used for pod submission is: %s\n", gasFees))
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

	unixTime := time.Now().Unix()
	fmt.Println(unixTime)
	currentTime := fmt.Sprintf("%d", unixTime)

	witnessByte, err := json.Marshal(publicWitness)
	if err != nil {
		components.Logger.Error(err.Error())
		return false
	}

	daDecode, daDecodeErr := chain.GetOldDaBlob()
	if daDecodeErr != nil {
		components.Logger.Error(daDecodeErr.Error())
		return false
	}
	var pMrh string
	if podNumber < 2 {
		pMrh = "0"
	} else {
		pMrh = daDecode.PreviousStateHash
	}
	msg := types.MsgSubmitPod{
		Creator:                newTempAddr,
		StationId:              stationId,
		PodNumber:              podNumber,
		MerkleRootHash:         merkleRootHash,
		PreviousMerkleRootHash: pMrh,
		PublicWitness:          witnessByte,
		Timestamp:              currentTime,
	}

	txRes, errTxRes := accountClient.BroadcastTx(ctx, newTempAccount, &msg)
	if errTxRes != nil {
		components.Logger.Error("error in transaction" + errTxRes.Error())
		return false
	}
	components.Logger.Info("Transaction Hash for SubmitPod: " + txRes.TxHash)
	return true
}
