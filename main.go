package main

import (
	"context"
	"fmt"
	"github.com/ComputerKeeda/junction-go-client/chain"
	"github.com/ComputerKeeda/junction-go-client/components/prover"
	"github.com/ComputerKeeda/junction-go-client/components/vrf"
	"github.com/ComputerKeeda/junction-go-client/core"
	"log"

	// Importing the general purpose Cosmos blockchain client
	"github.com/ignite/cli/v28/ignite/pkg/cosmosclient"

	// Importing the types package of your blog blockchain
	"github.com/ComputerKeeda/junction-go-client/components"
)

func main() {
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

	//RecursiveFunctions(stationId, client, ctx )
	core.InitVRF()
	core.GetVRF()

	//RecursiveFunctions()
}

func RecursiveFunctions() { // client, ctx context.Context

	/*
		Pod submit
		Pod verify

		Vrf init
		Vrf verify
	*/

	RecursiveFunctions()
}
