package core

import (
	"context"
	"fmt"
	"log"

	"github.com/ComputerKeeda/junction/x/junction/types"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosaccount"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosclient"
)

func InitStation(addr string, client cosmosclient.Client, ctx context.Context, account cosmosaccount.Account) {

	// Define a message to create a post
	msg := &types.MsgInitStation{
		Creator:         addr,
		Tracks:          []string{"air10vnvsez37eukd9hm9yp3969n6m8y93444upax8", "air10vnvsez37eukd9hm9yp3969n6m8y93444upax8", "air10vnvsez37eukd9hm9yp3969n6m8y93444upax8", "air10vnvsez37eukd9hm9yp3969n6m8y93444upax8"},
		VerificationKey: []byte("verificationKey"),
		StationId:       "stationId-324789",
		StationInfo:     "stationInfo-324789",
	}

	// Broadcast a transaction from account `alice` with the message
	// to create a post store response in txResp
	txResp, err := client.BroadcastTx(ctx, account, msg)
	if err != nil {
		fmt.Println("txResp above")
		fmt.Println(txResp)
		fmt.Println("txResp below")
		log.Fatal(err.Error())
	}

	// Print response from broadcasting a transaction
	fmt.Print("MsgCreatePost:\n\n")
	fmt.Println(txResp)

	// Instantiate a query client for your `blog` blockchain
	queryClient := types.NewQueryClient(client.Context())

	queryResp, err := queryClient.GetTracks(ctx, &types.QueryGetTracksRequest{StationId: "stationId-324789"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print("\n\nAll posts:\n\n")
	fmt.Println(queryResp)
}
