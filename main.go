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
	client, err := cosmosclient.New(ctx, cosmosclient.WithAddressPrefix(addressPrefix), cosmosclient.WithNodeAddress("http://192.168.1.37:26657"))
	if err != nil {
		log.Fatal(err)
	}

	// Account `alice` was initialized during `ignite chain serve`
	accountName := "alice"

	isAccountExists, _ := components.CheckIfAccountExists(accountName, client, addressPrefix, accountPath)
	if !isAccountExists {
		fmt.Println("Creating rollup account")
		components.CreateAccount(accountName, accountPath)
	}

	account, err := client.Account(accountName)
	if err != nil {
		fmt.Println("Error getting account")
	}

	addr, err := account.Address(addressPrefix)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Account Address:", addr)

	newTempAddr, newTempAccount := core.GetTempAddr(10)
	fmt.Println("Temp Address:", newTempAddr)
	newAccountClient, err := cosmosclient.New(ctx, cosmosclient.WithAddressPrefix(addressPrefix), cosmosclient.WithNodeAddress("http://192.168.1.37:26657"), cosmosclient.WithHome("./temp-account"))
	if err != nil {
		log.Fatal(err)
	}
	core.InitStation(newTempAddr, newAccountClient, ctx, newTempAccount)
}
