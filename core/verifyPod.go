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

// VerifyPod This Function takes the Proof and the Inputs and Check that the Submitted Pods is correct or not according to  ZKProof and witness submitted and Generated while submitting the Pods Via SubmitPods in Junction
func VerifyPod(proofByte []byte) bool {
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
	gas := components.GenerateRandomWithFavour(279, 500, [2]int{280, 350}, 0.8)
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

	podDetails, err := chain.GetOldDaBlob()
	if err != nil {
		components.Logger.Error(err.Error())
		return false
	}

	if podNumber > 1 {
		fmt.Println("podDetails.PreviousStateHash", podDetails.PreviousStateHash)
		fmt.Println("podDetails.CurrentStateHash", podDetails.CurrentStateHash)
	}

	verifyPodStruct := types.MsgVerifyPod{
		Creator:                newTempAddr,
		StationId:              stationId,
		PodNumber:              podNumber,
		MerkleRootHash:         podDetails.CurrentStateHash,
		PreviousMerkleRootHash: podDetails.PreviousStateHash,
		ZkProof:                proofByte,
	}

	txRes, errTxRes := accountClient.BroadcastTx(ctx, newTempAccount, &verifyPodStruct)
	if errTxRes != nil {
		components.Logger.Error("error in transaction" + errTxRes.Error())
		return false
	}

	components.Logger.Info("Transaction Hash for VerifyPod: " + txRes.TxHash)

	return true

}
