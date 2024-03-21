package chain

import (
	"fmt"
	"github.com/ComputerKeeda/junction-go-client/components"
	"io"
	"os"
	"strconv"
)

func SetChainID(stationId string) {
	// Create or open the file chainId.txt
	file, err := os.Create("data/stationId.txt")
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

func SetVRFPubKey(pubKey string) {
	// Create or open the file chainId.txt
	file, err := os.Create("data/vrfPubKey.txt")
	if err != nil {
		// Handle the error if the file cannot be created
		components.Logger.Error(fmt.Sprintf("error creating vrfPubKey.txt: %v", err))
		return
	}
	defer file.Close()

	// Write the stationId to the file
	_, err = file.WriteString(pubKey)
	if err != nil {
		// Handle the error if the file cannot be written to
		components.Logger.Error(fmt.Sprintf("error writing to vrfPubKey.txt: %v", err))
		return
	}

	// Save the file
	err = file.Sync()
	if err != nil {
		// Handle the error if the file cannot be saved
		components.Logger.Error(fmt.Sprintf("error saving vrfPubKey.txt: %v", err))
		return
	}

	// Print the stationId
	components.Logger.Info(fmt.Sprintf("vrfPubKey ID: %s", pubKey))
}

func SetVRFPrivKey(privateKey string) {
	// Create or open the file chainId.txt
	file, err := os.Create("data/vrfPrivKey.txt")
	if err != nil {
		// Handle the error if the file cannot be created
		components.Logger.Error(fmt.Sprintf("error creating vrfPrivKey.txt: %v", err))
		return
	}
	defer file.Close()

	// Write the stationId to the file
	_, err = file.WriteString(privateKey)
	if err != nil {
		// Handle the error if the file cannot be written to
		components.Logger.Error(fmt.Sprintf("error writing to vrfPrivKey.txt: %v", err))
		return
	}

	// Save the file
	err = file.Sync()
	if err != nil {
		// Handle the error if the file cannot be saved
		components.Logger.Error(fmt.Sprintf("error saving vrfPrivKey.txt: %v", err))
		return
	}

	// Print the stationId
	components.Logger.Info(fmt.Sprintf("vrfPrivKey ID: %s", privateKey))
}

func GetStationId() (stationId string, err error) {
	// get station id
	file, err := os.Open("data/stationId.txt")
	if err != nil {
		return "", err
	}
	defer file.Close()
	buf := make([]byte, 1024) // Buffer size of 1024 bytes
	for {
		n, err := file.Read(buf)
		if err == io.EOF {
			break // End of file reached
		}
		if err != nil {
			return "", err
		}
		stationId = string(buf[:n])
	}

	return stationId, nil
}

func GetPodNumber() (podNumber uint64, err error) {

	var podNumberStr string
	// get station id
	file, err := os.Open("data/podNumber.txt")
	if err != nil {
		return 0, err
	}
	defer file.Close()
	buf := make([]byte, 1024) // Buffer size of 1024 bytes
	for {
		n, err := file.Read(buf)
		if err == io.EOF {
			break // End of file reached
		}
		if err != nil {
			return 0, err
		}
		podNumberStr = string(buf[:n])
	}

	// string to uint64
	num, err := strconv.ParseUint(podNumberStr, 10, 64)
	if err != nil {
		return 0, err

	}

	return num, nil
}

func GetPrivateKey() (privateKey string, err error) {
	// get private Key
	file, err := os.Open("data/vrfPrivKey.txt")
	if err != nil {
		return "", err
	}
	defer file.Close()
	buf := make([]byte, 1024) // Buffer size of 1024 bytes
	for {
		n, err := file.Read(buf)
		if err == io.EOF {
			break // End of file reached
		}
		if err != nil {
			return "", err
		}
		privateKey = string(buf[:n])
	}

	return privateKey, nil
}

func GetPubKey() (pubKey string, err error) {
	// get private Key
	file, err := os.Open("data/vrfPubKey.txt")
	if err != nil {
		return "", err
	}
	defer file.Close()
	buf := make([]byte, 1024) // Buffer size of 1024 bytes
	for {
		n, err := file.Read(buf)
		if err == io.EOF {
			break // End of file reached
		}
		if err != nil {
			return "", err
		}
		pubKey = string(buf[:n])
	}

	return pubKey, nil
}
