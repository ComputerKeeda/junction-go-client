package main

import (
	"context"
	"fmt"
	"github.com/ComputerKeeda/junction-go-client/chain"
	"github.com/ComputerKeeda/junction-go-client/components/prover"
	"log"

	// Importing the general purpose Cosmos blockchain client
	"github.com/ignite/cli/v28/ignite/pkg/cosmosclient"

	// Importing the types package of your blog blockchain
	"github.com/ComputerKeeda/junction-go-client/components"
	"github.com/ComputerKeeda/junction-go-client/core"
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
	//RecursiveFunctions(stationId)
}

func RecursiveFunctions(stationId string) {

	/*
		Vrf init
		Vrf verify
		Pod submit
		Pod verify
	*/

	RecursiveFunctions(stationId)
}
