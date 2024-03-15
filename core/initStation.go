package core

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ComputerKeeda/junction-go-client/components"
	"github.com/ComputerKeeda/junction-go-client/types"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosaccount"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosclient"
	"github.com/joho/godotenv"
	"github.com/tjarratt/babble"
	"os"
)

func generateRandomWord() string {
	babbler := babble.NewBabbler()
	babbler.Count = 1
	randomWord := babbler.Babble()
	return randomWord
}

// createConfigDirectory creates the config directory if it doesn't exist.
func createConfigDirectory() error {
	path := "./config"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.Mkdir(path, 0755)
	}
	return nil
}

// writeChainInfoToJson writes the chain information to a JSON file.
func writeChainInfoToJson(info types.ChainInfoStruct) error {
	filePath := "./config/chainInfo.json"
	file, err := json.MarshalIndent(info, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, file, 0644)
}

func generateRandomChainInfo() (chainInfo types.ChainInfoStruct, err error) {
	err = godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return chainInfo, err
	}

	daWalletAddress := os.Getenv("DaWalletAddress")
	DaWalletKeypair := os.Getenv("DaWalletKeypair")

	// Initialize the ChainInfoStruct variable
	var chainInfoVar types.ChainInfoStruct

	// Set values for each field
	chainInfoVar.ChainInfo.ChainID = generateRandomWord()
	chainInfoVar.ChainInfo.Key = generateRandomWord()
	chainInfoVar.ChainInfo.Moniker = generateRandomWord()

	chainInfoVar.DaInfo.DaSelected = "celestia"
	chainInfoVar.DaInfo.DaWalletAddress = daWalletAddress
	chainInfoVar.DaInfo.DaWalletKeypair = DaWalletKeypair

	chainInfoVar.SequencerInfo.SequencerType = "AIRC"

	return chainInfoVar, nil
}

func InitStation(addr string, client cosmosclient.Client, ctx context.Context, account cosmosaccount.Account) {

	chainInfo, err := generateRandomChainInfo()
	if err != nil {
		components.Logger.Error("error in generating chain info")
		return
	}

	// Create config directory
	err = createConfigDirectory()
	if err != nil {
		errString := fmt.Sprintf("Error creating config directory: %v", err)
		components.Logger.Error(errString)
		return
	}

	// Write chainInfo to JSON
	err = writeChainInfoToJson(chainInfo)
	if err != nil {
		errString := fmt.Sprintf("error in writing chain info to JSON: %s", err)
		components.Logger.Error(errString)
		return
	}

	components.Logger.Info("chainInfo.json has been successfully created.")

	//randomStationId, randomStationInfo := generateRandomString()
	//// Define a message to create a post
	//msg := &types.MsgInitStation{
	//	Creator:         addr,
	//	Tracks:          []string{"air10vnvsez37eukd9hm9yp3969n6m8y93444upax8", "air10vnvsez37eukd9hm9yp3969n6m8y93444upax8", "air10vnvsez37eukd9hm9yp3969n6m8y93444upax8", "air10vnvsez37eukd9hm9yp3969n6m8y93444upax8"},
	//	VerificationKey: []byte("verificationKey"),
	//	StationId:       randomStationId,
	//	StationInfo:     randomStationInfo,
	//}
	//
	//// Broadcast a transaction from account `alice` with the message
	//// to create a post store response in txResp
	//txResp, err := client.BroadcastTx(ctx, account, msg)
	//if err != nil {
	//	fmt.Println("txResp above")
	//	fmt.Println(txResp)
	//	fmt.Println("txResp below")
	//	log.Fatal(err.Error())
	//}
	//
	//// Print response from broadcasting a transaction
	//fmt.Print("MsgCreatePost:\n\n")
	//fmt.Println(txResp)
	//
	//// Instantiate a query client for your `blog` blockchain
	//queryClient := types.NewQueryClient(client.Context())
	//
	//queryResp, err := queryClient.GetTracks(ctx, &types.QueryGetTracksRequest{StationId: randomStationId})
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Print("\n\nAll posts:\n\n")
	//fmt.Println(queryResp)
}

// Function to generate a random string in the format "stationId-xxxx"
//func generateRandomString() (string, string) {
//	// Seed the random number generator
//	rand.New(rand.NewSource(time.Now().UnixNano()))
//
//	// Generate a random four-digit number
//	randomNumber := rand.Intn(9000) + 1000 // Generates a number between 1000 and 9999
//
//	// Concatenate the string parts
//	randomStationId := fmt.Sprintf("stationId-%d", randomNumber)
//	randomStationInfo := fmt.Sprintf("stationInfo-%d", randomNumber)
//
//	return randomStationId, randomStationInfo
//}
