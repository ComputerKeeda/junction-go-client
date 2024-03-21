package core

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ComputerKeeda/junction-go-client/components"
	"github.com/ComputerKeeda/junction-go-client/types"
	"github.com/google/uuid"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosaccount"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosclient"
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
	// Initialize the ChainInfoStruct variable
	var chainInfoVar types.ChainInfoStruct

	// Set values for each field
	chainInfoVar.ChainInfo.ChainID = generateRandomWord()
	chainInfoVar.ChainInfo.Key = generateRandomWord()
	chainInfoVar.ChainInfo.Moniker = generateRandomWord()

	chainInfoVar.DaInfo.DaSelected = "celestia"
	chainInfoVar.DaInfo.DaWalletAddress = "celestia1kd94m7dsh87fd452sdlpqz0as64mp6f07ln27l"
	chainInfoVar.DaInfo.DaWalletKeypair = "0068ECE1E6FB5359"

	chainInfoVar.SequencerInfo.SequencerType = "AIRC"

	return chainInfoVar, nil
}

func InitStation() string {
	accountName := "temp-account"
	accountPath := "./temp-account"
	addressPrefix := "air"
	registry, err := cosmosaccount.New(cosmosaccount.WithHome(accountPath))
	if err != nil {
		components.Logger.Error(fmt.Sprintf("Error creating account registry: %v", err))
	}

	newTempAccount, err := registry.GetByName(accountName)
	if err != nil {
		components.Logger.Error(fmt.Sprintf("Error getting account: %v", err))
	}

	newTempAddr, err := newTempAccount.Address(addressPrefix)
	if err != nil {
		components.Logger.Error(fmt.Sprintf("Error getting address: %v", err))
	}

	chainInfo, err := generateRandomChainInfo()
	if err != nil {
		components.Logger.Error("error in generating chain info")
		return ""
	}

	// Create config directory
	err = createConfigDirectory()
	if err != nil {
		errString := fmt.Sprintf("Error creating config directory: %v", err)
		components.Logger.Error(errString)
		return ""
	}

	// Write chainInfo to JSON
	err = writeChainInfoToJson(chainInfo)
	if err != nil {
		errString := fmt.Sprintf("error in writing chain info to JSON: %s", err)
		components.Logger.Error(errString)
		return ""
	}

	components.Logger.Info("chainInfo.json has been successfully created.")

	verificationKeyContents, err := os.ReadFile("verificationKey.json")
	if err != nil {
		components.Logger.Error("Error reading verificationKey.json file")
		return ""
	}

	//occupancy := components.GenerateRandomWithFavour(1, 10, [2]int{2, 5}, 0.5)
	occupancy := 4
	otherTrackMembers := components.GenerateAddresses(occupancy)

	// tracks voting power creator
	var tracksVotingPower []uint64
	power := uint64(100 / 5)
	for i := 0; i < 5; i++ {
		tracksVotingPower = append(tracksVotingPower, power)
	}

	var tracks []string
	tracks = append(tracks, newTempAddr)
	tracks = append(tracks, otherTrackMembers...)
	stationId := uuid.New()
	randomStationInfo, randomStationInfoError := generateRandomChainInfo()
	if randomStationInfoError != nil {
		components.Logger.Error("Error generating random station info")
		return ""
	}
	chainInfoAsString, err := json.Marshal(randomStationInfo.ChainInfo)
	if err != nil {
		components.Logger.Error("Error marshalling chain info")
		return ""
	}

	ctx := context.Background()
	gas := components.GenerateRandomWithFavour(600, 1200, [2]int{611, 1000}, 0.7)
	gasFees := fmt.Sprintf("%damf", gas)
	accountClient, err := cosmosclient.New(ctx, cosmosclient.WithAddressPrefix("air"), cosmosclient.WithNodeAddress(components.JunctionTTCRPC), cosmosclient.WithHome("./temp-account"), cosmosclient.WithGas("auto"), cosmosclient.WithFees(gasFees))
	if err != nil {
		components.Logger.Error("Error creating account client")
		return ""
	}

	extraArg := types.StationArg{
		TrackType: "Airchains sequencer",
		DaType:    "Celestia",
		Prover:    "Airchains",
	}

	extraArgBytes, err := json.Marshal(extraArg)
	if err != nil {
		components.Logger.Error("Error marshalling extra arg")
		return ""
	}

	newStationData := types.MsgInitStation{
		Creator:           newTempAddr,
		Tracks:            tracks,
		VerificationKey:   verificationKeyContents,
		StationId:         stationId.String(),
		StationInfo:       string(chainInfoAsString),
		TracksVotingPower: tracksVotingPower,
		ExtraArg:          extraArgBytes,
	}

	txResp, err := accountClient.BroadcastTx(ctx, newTempAccount, &newStationData)
	if err != nil {
		components.Logger.Error("Error in broadcasting transaction")
		components.Logger.Error(err.Error())
		return ""
	}

	components.Logger.Info("Station created successfully")
	components.Logger.Info(fmt.Sprintf("Transaction hash: %s", txResp.TxHash))

	return stationId.String()
}
