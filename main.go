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
	for {
		validateVRFResponse := core.ValidateVRF(funcLevelSerializedRc)
		if !validateVRFResponse {
			time.Sleep(5 * time.Second)
		} else {
			break
		}
	}

	// check if verification success
	var vrfRecord *types.VrfRecord
	for {
		vrfRecord = core.GetVRF()
		if vrfRecord == nil {
			components.Logger.Error("VRF record is nil")
			time.Sleep(5 * time.Second)
		} else {
			break
		}
		if !vrfRecord.IsVerified {
			components.Logger.Error("Verification failed")
			time.Sleep(5 * time.Second)
		} else {
			break
		}
	}

	var funcLevelPodNumber uint64
	for {
		podNumber, err := chain.GetPodNumber()
		if err != nil {
			components.Logger.Error(err.Error())
			time.Sleep(5 * time.Second)
		} else {
			funcLevelPodNumber = podNumber
			break
		}
	}

	intPodNumber := int(funcLevelPodNumber)

	var funcLevelWitnessVector any
	var funcLevelCurrentStatusHash string
	var funcLevelProofByte []byte
	var funcLevelPsTxHash []string
	for {
		ps := components.GenerateUniqueRandomValues(25) // Generate 25 sets of unique values
		funcLevelPsTxHash = ps.TransactionHash
		witnessVector, currentStatusHash, proofByte, pkErr := prover.GenerateProof(ps, intPodNumber)
		if pkErr != nil {
			components.Logger.Error(fmt.Sprintf("Error in generating proof : %s", pkErr.Error()))
			time.Sleep(5 * time.Second)
		} else {
			funcLevelWitnessVector = witnessVector
			funcLevelCurrentStatusHash = currentStatusHash
			funcLevelProofByte = proofByte
			break
		}
	}

	// Mock DA call
	for {
		success, _ := chain.MockDa(funcLevelPsTxHash, funcLevelCurrentStatusHash, intPodNumber)
		if success {
			var podSubmitSuccess bool
			// submit pod
			for {
				podSubmitSuccess = core.SubmitPod(funcLevelCurrentStatusHash, funcLevelWitnessVector)
				if !podSubmitSuccess {
					components.Logger.Error("Pod submission failed")
					time.Sleep(3 * time.Second)
					continue
				} else {
					break
				}
			}
			// verify pod
			for {
				verifyPodSuccess := core.VerifyPod(funcLevelProofByte)
				if !verifyPodSuccess {
					components.Logger.Error("Pod Verification failed")
					time.Sleep(5 * time.Second)
					continue
				} else {
					break
				}
			}
			break
		} else {
			//	else we have to show the error here
			components.Logger.Error("DA failed retrying... in 3 seconds")
			time.Sleep(3 * time.Second)
			continue
		}
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
