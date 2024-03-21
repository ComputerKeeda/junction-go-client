package core

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ComputerKeeda/junction-go-client/components"
	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	cosmosBankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosaccount"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosclient"
)

func GetTempAddr(amount int64) (newTempAddr string, newTempAccount cosmosaccount.Account) {

	accountPath := "./temp-account"
	// Remove the directory and any subdirectories it contains
	err := os.RemoveAll(accountPath)
	if err != nil {
		fmt.Println("Error deleting directory:", err)
	}

	ctx := context.Background()
	addressPrefix := "air"

	// Create a Cosmos newAccountClient instance
	newAccountClient, err := cosmosclient.New(ctx, cosmosclient.WithAddressPrefix(addressPrefix), cosmosclient.WithNodeAddress(components.JunctionTTCRPC), cosmosclient.WithHome(accountPath), cosmosclient.WithGas("auto"), cosmosclient.WithFees("25amf"))
	if err != nil {
		log.Fatal(err)
	}

	adminAccountClient, err := cosmosclient.New(ctx, cosmosclient.WithAddressPrefix(addressPrefix), cosmosclient.WithNodeAddress(components.JunctionTTCRPC), cosmosclient.WithHome("./accounts"), cosmosclient.WithGas("auto"), cosmosclient.WithFees("25amf"))
	if err != nil {
		log.Fatal(err)
	}

	adminAccountName := "alice"
	adminAccount, err := adminAccountClient.Account(adminAccountName)
	if err != nil {
		fmt.Println("Error getting account")
	}

	adminAddress, err := adminAccount.Address(addressPrefix)
	if err != nil {
		log.Fatal(err)
	}

	components.Logger.Warn(fmt.Sprintf("Admin alice address : %s\n", adminAddress))

	accountName := "temp-account"
	components.CreateAccount(accountName, accountPath)

	newTempAccount, err = newAccountClient.Account(accountName)
	if err != nil {
		fmt.Println("Error getting account")
	}

	newTempAddr, err = newTempAccount.Address(addressPrefix)
	if err != nil {
		log.Fatal(err)
	}

	msg := &cosmosBankTypes.MsgSend{
		FromAddress: adminAddress,
		ToAddress:   newTempAddr,
		Amount:      cosmosTypes.NewCoins(cosmosTypes.NewInt64Coin("amf", amount)),
	}

	txResp, err := adminAccountClient.BroadcastTx(ctx, adminAccount, msg)
	if err != nil {
		fmt.Println("Error in sending tokens to temp account")
		fmt.Println("error in transaction", err)
	}

	fmt.Printf("Tx Hash: %v\n", txResp.TxHash)

	return newTempAddr, newTempAccount
}
