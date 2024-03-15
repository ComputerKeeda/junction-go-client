# Junction Go client

This is a Go client for the Junction Chain RPC. It is a work in progress and is not yet ready for production use.


```go
package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func GenerateVerificationKey() (ProvingKey, VerificationKey, error) {
	// Dummy implementation - replace with your actual logic
	return ProvingKey{}, VerificationKey{}, nil
}

type ProvingKey struct{} // Define ProvingKey struct
func (pk ProvingKey) WriteTo(file *os.File) (int64, error) {
	// Dummy implementation - replace with actual logic to write to file
	return 0, nil
}

type VerificationKey struct{} // Define VerificationKey struct

func main() {
	verificationKeyFile := "verificationKey.json"
	provingKeyFile := "provingKey.txt"

	// Check and generate keys if necessary
	generateIfNotExist(provingKeyFile, verificationKeyFile)

	// Log existing files
	logExistingFiles(provingKeyFile, verificationKeyFile)
}

func generateIfNotExist(provingKeyFile, verificationKeyFile string) {
	if _, err := os.Stat(provingKeyFile); os.IsNotExist(err) {
		if _, err := os.Stat(verificationKeyFile); os.IsNotExist(err) {
			// Generate both keys if neither file exists
			provingKey, verificationKey, err := GenerateVerificationKey()
			if err != nil {
				fmt.Println("Error generating keys:", err)
				return
			}
			writeVerificationKeyToFile(verificationKey, verificationKeyFile)
			writeProvingKeyToFile(provingKey, provingKeyFile)
		}
	} else if _, err := os.Stat(verificationKeyFile); os.IsNotExist(err) {
		// Generate verification key if only proving key file exists
		_, verificationKey, err := GenerateVerificationKey()
		if err != nil {
			fmt.Println("Error generating verification key:", err)
			return
		}
		writeVerificationKeyToFile(verificationKey, verificationKeyFile)
	}
}

func writeVerificationKeyToFile(verificationKey VerificationKey, verificationKeyFile string) {
	vkJSON, err := json.Marshal(verificationKey)
	if err != nil {
		fmt.Println("Error marshalling verification key:", err)
		return
	}
	if err := os.WriteFile(verificationKeyFile, vkJSON, 0644); err != nil {
		fmt.Println("Error writing verification key to file:", err)
	}
}

func writeProvingKeyToFile(provingKey ProvingKey, provingKeyFile string) {
	file, err := os.Create(provingKeyFile)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	if _, err = provingKey.WriteTo(file); err != nil {
		fmt.Println("Error writing proving key to file:", err)
	}
}

func logExistingFiles(provingKeyFile, verificationKeyFile string) {
	if _, err := os.Stat(verificationKeyFile); err == nil {
		fmt.Println("Verification key already exists. No action needed.")
	}
	if _, err := os.Stat(provingKeyFile); err == nil {
		fmt.Println("Proving key already exists. No action needed.")
	}
}
```

```
package main

import (
	"encoding/json"
	"os"
)

// Assuming these are defined somewhere in your project
var components struct {
	Logger struct {
		Info  func(msg string)
		Error func(msg string)
	}
}

func GenerateVerificationKey() (ProvingKey, VerificationKey, error) {
	// Dummy implementation - replace with your actual logic
	return ProvingKey{}, VerificationKey{}, nil
}

type ProvingKey struct{} // Define ProvingKey struct
func (pk ProvingKey) WriteTo(file *os.File) (int64, error) {
	// Dummy implementation - replace with actual logic to write to file
	return 0, nil
}

type VerificationKey struct{} // Define VerificationKey struct

func main() {
	verificationKeyFile := "verificationKey.json"
	provingKeyFile := "provingKey.txt"

	// Check and generate keys if necessary
	generateIfNotExist(provingKeyFile, verificationKeyFile)

	// Log existing files
	logExistingFiles(provingKeyFile, verificationKeyFile)
}

func generateIfNotExist(provingKeyFile, verificationKeyFile string) {
	if _, err := os.Stat(provingKeyFile); os.IsNotExist(err) {
		if _, err := os.Stat(verificationKeyFile); os.IsNotExist(err) {
			// Generate both keys if neither file exists
			provingKey, verificationKey, err := GenerateVerificationKey()
			if err != nil {
				components.Logger.Error("Error generating keys: " + err.Error())
				return
			}
			writeVerificationKeyToFile(verificationKey, verificationKeyFile)
			writeProvingKeyToFile(provingKey, provingKeyFile)
		}
	} else if _, err := os.Stat(verificationKeyFile); os.IsNotExist(err) {
		// Generate verification key if only proving key file exists
		_, verificationKey, err := GenerateVerificationKey()
		if err != nil {
			components.Logger.Error("Error generating verification key: " + err.Error())
			return
		}
		writeVerificationKeyToFile(verificationKey, verificationKeyFile)
	}
}

func writeVerificationKeyToFile(verificationKey VerificationKey, verificationKeyFile string) {
	vkJSON, err := json.Marshal(verificationKey)
	if err != nil {
		components.Logger.Error("Error marshalling verification key: " + err.Error())
		return
	}
	if err := os.WriteFile(verificationKeyFile, vkJSON, 0644); err != nil {
		components.Logger.Error("Error writing verification key to file: " + err.Error())
	}
}

func writeProvingKeyToFile(provingKey ProvingKey, provingKeyFile string) {
	file, err := os.Create(provingKeyFile)
	if err != nil {
		components.Logger.Error("Error creating file: " + err.Error())
		return
	}
	defer file.Close()

	if _, err = provingKey.WriteTo(file); err != nil {
		components.Logger.Error("Error writing proving key to file: " + err.Error())
	}
}

func logExistingFiles(provingKeyFile, verificationKeyFile string) {
	if _, err := os.Stat(verificationKeyFile); err == nil {
		components.Logger.Info("Verification key already exists. No action needed.")
	}
	if _, err := os.Stat(provingKeyFile); err == nil {
		components.Logger.Info("Proving key already exists. No action needed.")
	}
}
```