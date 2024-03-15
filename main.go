package main

import (
	"context"
	"fmt"
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
	client, err := cosmosclient.New(ctx, cosmosclient.WithAddressPrefix(addressPrefix), cosmosclient.WithNodeAddress(components.JunctionTTCRPC), cosmosclient.WithHome(accountPath), cosmosclient.WithGas("auto"), cosmosclient.WithFees("200amf"))
	if err != nil {
		log.Fatal(err)
	}

	accountName := "alice"

	isAccountExists, _ := components.CheckIfAccountExists(accountName, client, addressPrefix, accountPath)
	if !isAccountExists {
		fmt.Println("Creating rollup account")
		components.CreateAccount(accountName, accountPath)
	}

	newTempAddr, _ := core.GetTempAddr(10)
	components.Logger.Warn(fmt.Sprintf("Station Creator Address: %s\n", newTempAddr))
}
