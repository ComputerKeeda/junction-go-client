package prover

import (
	"encoding/json"
	"github.com/ComputerKeeda/junction-go-client/components"
	"os"
)

// CreateVkPk generates and saves a Proving Key and a Verification Key.
// If either file doesn't exist, it will generate and save new keys. Otherwise, it will print a message stating that both keys already exist.
func CreateVkPk() {
	provingKeyFile := "provingKey.txt"
	verificationKeyFile := "verificationKey.json"

	_, err1 := os.Stat(provingKeyFile)
	_, err2 := os.Stat(verificationKeyFile)

	// If either file doesn't exist, generate and save new keys
	if os.IsNotExist(err1) || os.IsNotExist(err2) {
		provingKey, verificationKey, err := GenerateKeyPair()
		if err != nil {
			components.Logger.Error("Unable to generate key pair" + err.Error())
			return
		}

		// Save Proving Key
		pkFile, err := os.Create(provingKeyFile)
		if err != nil {
			components.Logger.Error("Unable to create Proving Key file" + err.Error())
			return
		}
		_, err = provingKey.WriteTo(pkFile)
		pkFile.Close()
		if err != nil {
			components.Logger.Error("Unable to write Proving Key" + err.Error())
			return
		}

		// Save Verification Key
		file, _ := json.MarshalIndent(verificationKey, "", " ")
		err = os.WriteFile(verificationKeyFile, file, 0644)
		if err != nil {
			components.Logger.Error("Unable to write Verification Key to file" + err.Error())
		}
		components.Logger.Info("Proving key and Verification key generated and saved successfully\n")
	} else {
		components.Logger.Info("Both Proving key and Verification key already exist. No action needed.")
	}
}
