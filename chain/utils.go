package chain

import (
	"fmt"
	"github.com/ComputerKeeda/junction-go-client/components"
	"os"
)

func SetChainID(stationId string) {
	// Create or open the file chainId.txt
	file, err := os.Create("data/chainId.txt")
	if err != nil {
		// Handle the error if the file cannot be created
		components.Logger.Error(fmt.Sprintf("error creating chainId.txt: %v", err))
		return
	}
	defer file.Close()

	// Write the stationId to the file
	_, err = file.WriteString(stationId)
	if err != nil {
		// Handle the error if the file cannot be written to
		components.Logger.Error(fmt.Sprintf("error writing to chainId.txt: %v", err))
		return
	}

	// Save the file
	err = file.Sync()
	if err != nil {
		// Handle the error if the file cannot be saved
		components.Logger.Error(fmt.Sprintf("error saving chainId.txt: %v", err))
		return
	}

	// Print the stationId
	components.Logger.Info(fmt.Sprintf("Station ID: %s", stationId))
}
