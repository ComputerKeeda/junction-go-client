package main

import (
	"context"
	"fmt"
	"github.com/ComputerKeeda/junction-go-client/chain"
	"github.com/ComputerKeeda/junction-go-client/components/prover"
	"github.com/ComputerKeeda/junction-go-client/components/vrf"
	"github.com/ComputerKeeda/junction-go-client/core"
	"github.com/ComputerKeeda/junction-go-client/types"
	"log"
	"time"

	// Importing the general purpose Cosmos blockchain client
	"github.com/ignite/cli/v28/ignite/pkg/cosmosclient"

	// Importing the types package of your blog blockchain
	"github.com/ComputerKeeda/junction-go-client/components"
)

func main() {
	//CreateOldDaBlog

	err := chain.CreateOldDaBlog()
	if err != nil {
		// If there was an error, print it
		fmt.Println("Error:", err)
	} else {
		// If the file was successfully created or already existed, print a success message
		fmt.Println("File ensured:", "oldDaBlog.json")
	}

	ctx := context.Background()
	addressPrefix := "air"
	accountPath := "./accounts"

	// Create a Cosmos client instance
	client, err := cosmosclient.New(ctx, cosmosclient.WithAddressPrefix(addressPrefix), cosmosclient.WithNodeAddress(components.JunctionTTCRPC), cosmosclient.WithHome(accountPath), cosmosclient.WithGas("auto"), cosmosclient.WithFees("25amf"))
	if err != nil {
		log.Fatal(err)
	}

	accountName := "alice"

	isAccountExists, _ := components.CheckIfAccountExists(accountName, client, addressPrefix, accountPath)
	if !isAccountExists {
		fmt.Println("Creating rollup account")
		components.CreateAccount(accountName, accountPath)
	}

	// !Important since this will be used to create a new address for the station

	//newTempAddr, _ := core.GetTempAddr(631478)
	//components.Logger.Warn(fmt.Sprintf("Station Creator Address: %s\n", newTempAddr))

	prover.CreateVkPk()
	stationId := core.InitStation()
	if stationId == "" {
		log.Fatal("Error initializing station")
	} else {
		chain.SetChainID(stationId)
	}

	// generating vrf key pair
	vrfPrivateKey, vrfPublicKey := vrf.NewKeyPair()
	vrfPrivateKeyHex := vrfPrivateKey.String()
	vrfPublicKeyHex := vrfPublicKey.String()
	components.Logger.Debug(fmt.Sprintf("VRF Private Key: %s\n", vrfPrivateKeyHex))
	components.Logger.Debug(fmt.Sprintf("VRF Public Key: %s\n", vrfPublicKeyHex))

	if vrfPrivateKeyHex != "" {
		chain.SetVRFPrivKey(vrfPrivateKeyHex)
	} else {
		components.Logger.Error("Error saving VRF private key")
	}

	if vrfPublicKeyHex != "" {
		chain.SetVRFPubKey(vrfPublicKeyHex)
	} else {
		components.Logger.Error("Error saving VRF public key")
	}

	RecursiveFunctions()
}

func RecursiveFunctions() { // client, ctx context.Context

	var funcLevelSerializedRc []byte
	for {
		initVRFResponse, serializedRc := core.InitVRF()
		if !initVRFResponse {
			time.Sleep(5 * time.Second)
		} else {
			funcLevelSerializedRc = serializedRc
			break
		}
	}

	// validate vrf
	validateVRFResponse := core.ValidateVRF(serializedRc)
	if !validateVRFResponse {
		return
	}

	// check if verification success
	vrfRecord := core.GetVRF()
	if vrfRecord == nil {
		return
	}
	if !vrfRecord.IsVerified {
		components.Logger.Error("Verification failed")
		return
	}

	podNumber, err := chain.GetPodNumber()
	if err != nil {
		components.Logger.Error(err.Error())
		return
	}
	intPodNumber := int(podNumber)

	ps := components.GenerateUniqueRandomValues(25) // Generate 25 sets of unique values
	witnessVector, currentStatusHash, proofByte, pkErr := prover.GenerateProof(ps, intPodNumber)
	if pkErr != nil {
		components.Logger.Error(fmt.Sprintf("Error in generating proof : %s", pkErr.Error()))
		return
	}

	// Mock DA call
	success, daKeyHash := chain.MockDa(ps.TransactionHash, currentStatusHash, intPodNumber)

	// submit pod
	podSubmitSuccess := core.SubmitPod(currentStatusHash, witnessVector)
	if !podSubmitSuccess {
		components.Logger.Error("Pod submission failed")
		return
	}

	for {
		err := chain.UpdatePodNumber()
		if err != nil {
			components.Logger.Error(fmt.Sprintf("Error updating pod number: %s", err.Error()))
			time.Sleep(5 * time.Second)
			continue
		} else {
			break
		}
	}

	components.Logger.Debug(fmt.Sprintf("Pod Number Updated: %d", intPodNumber))
	components.Logger.Debug("🚧🚧🚧🚧🚧🚧🚧🚧🚧🚧🚧🚧🚧🚧🚧🚧🚧🚧🚧🚧🚧🚧🚧🚧🚧")

	RecursiveFunctions()
}
