package chain

import (
	"github.com/ComputerKeeda/junction-go-client/components"
	"github.com/ComputerKeeda/junction-go-client/types"
	"strconv"
	"time"
)

type MockDAStruct struct {
	DataBlob    []string
	BatchNumber string
	Commitment  string
}

type DAUploadStruct struct {
	//Proof             ProofStruct `json:"proof"`
	TxnHashes         []string `json:"txnHashes"`
	CurrentStateHash  string   `json:"currentStateHash"`
	PreviousStateHash string   `json:"previousStateHash"`
	MetaData          struct {
		ChainID     string `json:"chainID"`
		BatchNumber int    `json:"batchNumber"`
	} `json:"metaData"`
}

func MockDa(transactions []string, currentStateHash string, batchNumber int) (string, error) {
	daDecode, daDecodeErr := GetOldDaBlob()
	if daDecodeErr != nil {
		return "", daDecodeErr
	}
	stationId, err := GetStationId()
	if err != nil {
		components.Logger.Error(err.Error())
		return "", err
	}
	DaStruct := DAUploadStruct{
		TxnHashes:         transactions,
		CurrentStateHash:  currentStateHash,
		PreviousStateHash: daDecode.PreviousStateHash,
		MetaData: struct {
			ChainID     string `json:"chainID"`
			BatchNumber int    `json:"batchNumber"`
		}{
			ChainID:     stationId,
			BatchNumber: batchNumber,
		},
	}

	_, daResponse := MockDaSubmit(DaStruct)

	da := types.DAStruct{
		DAKey:             daResponse,
		DAClientName:      "mock",
		BatchNumber:       strconv.Itoa(batchNumber),
		PreviousStateHash: daDecode.CurrentStateHash,
		CurrentStateHash:  currentStateHash,
	}
	// ADD DA CLIENT CODES

	status := SetOldDaBlob(da)
	if !status {
		time.Sleep(3 * time.Second)
		_, _ = MockDa(transactions, currentStateHash, batchNumber)
	}
	return "", nil
}

// MockDaSubmit  is a function that mocks the functionality of storing data in a mock database (leveldb). It takes the following parameters:
// - PodData: a byte slice containing the data to be stored
// - podNumber: a string representing the Pod number
//
// The function performs the following steps:
// 1. Computes the SHA256 hash of daData.
// 2. Encodes the hash as a string using hexadecimal encoding.
// 3. Creates a new MockDAStruct instance with the podData, podNumber, and computed hashString.
// 4. Converts the mockData into a byte slice.
// 5. Return True
func MockDaSubmit(PodStruct DAUploadStruct) (bool, string) {
	_ = PodStruct
	return true, "qwertyuiop"
}
