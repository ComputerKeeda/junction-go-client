package core

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ComputerKeeda/junction-go-client/chain"
	"github.com/ComputerKeeda/junction-go-client/components"
	"github.com/ComputerKeeda/junction-go-client/components/vrf"
	"github.com/ComputerKeeda/junction-go-client/types"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosaccount"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosclient"
	"go.dedis.ch/kyber/v3/group/edwards25519"
)

func InitVRF() bool {
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
	gas := components.GenerateRandomWithFavour(69, 120, [2]int{70, 90}, 0.7)
	gasFees := fmt.Sprintf("%damf", gas)
	components.Logger.Warn(fmt.Sprintf("Gas Fees Used for init VRF transaction is: %s\n", gasFees))
	accountClient, err := cosmosclient.New(ctx, cosmosclient.WithAddressPrefix(addressPrefix), cosmosclient.WithNodeAddress(components.JunctionTTCRPC), cosmosclient.WithHome(accountPath), cosmosclient.WithGas("auto"), cosmosclient.WithFees(gasFees))
	if err != nil {
		components.Logger.Error("Error creating account client")
		return false
	}
	// getting the account and creating client codes --> End

	// get variables required to generate or call verifiable random number
	suite := edwards25519.NewBlakeSHA256Ed25519()
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
	privateKeyStr, err := chain.GetPrivateKey()
	if err != nil {
		components.Logger.Error(err.Error())
		return false
	}
	privateKey, err := vrf.LoadHexPrivateKey(privateKeyStr)
	if err != nil {
		components.Logger.Error("Error in loading private key: " + err.Error())
		return false
	}
	publicKey, err := chain.GetPubKey()
	if err != nil {
		components.Logger.Error(err.Error())
		return false
	}

	rc := vrf.RequestCommitmentV2Plus{
		BlockNum:         1,
		StationId:        stationId,
		UpperBound:       5,
		RequesterAddress: newTempAddr,
	}

	serializedRC, err := vrf.SerializeRequestCommitmentV2Plus(rc)
	if err != nil {
		components.Logger.Error(err.Error())
		return false
	}
	proof, vrfOutput, err := vrf.GenerateVRFProof(suite, privateKey, serializedRC, int64(rc.BlockNum))
	if err != nil {
		fmt.Printf("Error generating unique proof: %v\n", err)
		return false
	}

	extraArg := types.ExtraArg{
		SerializedRc: serializedRC,
		Proof:        proof,
		VrfOutput:    vrfOutput,
	}
	// marshal
	extraArgsByte, err := json.Marshal(extraArg)
	if err != nil {
		components.Logger.Error(err.Error())
		return false
	}

	var defaultOccupancy uint64
	defaultOccupancy = 5
	msg := types.MsgInitiateVrf{
		Creator:        newTempAddr,
		PodNumber:      podNumber,
		StationId:      stationId,
		Occupancy:      defaultOccupancy,
		CreatorsVrfKey: publicKey,
		ExtraArg:       extraArgsByte,
	}

	txRes, errTxRes := accountClient.BroadcastTx(ctx, newTempAccount, &msg)
	if errTxRes != nil {
		components.Logger.Error("error in transaction" + errTxRes.Error())
		return false
	}

	components.Logger.Info("Transaction Hash: " + txRes.TxHash)

	return true

}
