package core

import (
	"context"
	"fmt"
	"log"

	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	cosmosBankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/ignite/cli/v28/ignite/pkg/cosmosaccount"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosclient"
)

func SendToken(reqAmount int64, toAddress string, ctx context.Context, account cosmosaccount.Account, adminAddress string) (success bool, message string, txhash string, err error) {

	accountPath := "./accounts"
	addressPrefix := "air"

	// Create a Cosmos client instance
	client, err := cosmosclient.New(ctx, cosmosclient.WithAddressPrefix(addressPrefix), cosmosclient.WithNodeAddress("http://192.168.1.37:26657"), cosmosclient.WithHome(accountPath))
	if err != nil {
		log.Fatal(err)
	}
	// check if admin have enough balance to send tokens
	isBalance, amount, _ := CheckBalance(ctx, adminAddress, client)
	fmt.Printf("Chhheekkkkkkkkk   Admin balance: %v\n", amount)
	// no balance in admin wallet
	if !isBalance {
		return false, "No balance in admin wallet", "", err
	}

	// admin have less than 10 tokens (not enough balance)
	if amount < 10 {
		return false, "Admin don't have enough tokens", "", err
	}

	msg := &cosmosBankTypes.MsgSend{
		FromAddress: adminAddress,
		ToAddress:   toAddress,
		Amount:      cosmosTypes.NewCoins(cosmosTypes.NewInt64Coin("stake", reqAmount)),
	}

	txResp, err := client.BroadcastTx(ctx, account, msg)
	if err != nil {
		return false, "error in transaction", "", err
	}

	return true, "Success", txResp.TxHash, nil
}
