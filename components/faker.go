package components

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/ComputerKeeda/junction-go-client/types"
	"io"
	rand2 "math/rand"
	"strconv"
)

// generateRandomHex generates a random hexadecimal string of a specified length.
func generateRandomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, bytes); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", bytes), nil
}

// generateEthereumAddress generates a random Ethereum wallet address.
func generateEthereumAddress() (string, error) {
	address, err := generateRandomHex(20) // 20 bytes = 40 hex characters
	if err != nil {
		return "", err
	}
	return "0x" + address, nil
}

// generateTransactionHash generates a random Ethereum transaction hash.
func generateTransactionHash() (string, error) {
	txHash, err := generateRandomHex(32) // 32 bytes = 64 hex characters
	if err != nil {
		return "", err
	}
	return "0x" + txHash, nil
}

// generateRandomString generates a random hexadecimal string of specified length.
func generateRandomString(length int) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}

func generateRandomNumber() string {
	tempAmount := GenerateRandomWithFavour(1, 1000000, [2]int{100, 100000}, 0.7)
	stringAmount := fmt.Sprintf("%d", tempAmount)
	return stringAmount
}

func generateRandomSmallNumber() string { // for nonce
	tempValue := GenerateRandomWithFavour(1, 1000, [2]int{1, 100}, 0.7)
	stringAmount := fmt.Sprintf("%d", tempValue)
	return stringAmount
}

func GenerateUniqueRandomValuesNew(count int) types.PodStruct {
	ps := types.PodStruct{
		From:              make([]string, count),
		To:                make([]string, count),
		Amounts:           make([]string, count),
		TransactionHash:   make([]string, count),
		SenderBalances:    make([]string, count),
		ReceiverBalances:  make([]string, count),
		Messages:          make([]string, count),
		TransactionNonces: make([]string, count),
		AccountNonces:     make([]string, count),
	}

	for i := 0; i < count; i++ {
		randomFrom, _ := generateEthereumAddress()
		randomTo, _ := generateEthereumAddress()
		randomTxHash, _ := generateTransactionHash()

		ps.From[i] = randomFrom // Adjust the length as needed
		ps.To[i] = randomTo
		ps.Amounts[i] = generateRandomNumber() // Assuming amounts are shorter strings
		ps.TransactionHash[i] = randomTxHash
		// Generate sender balances first
		ps.SenderBalances[i] = generateRandomNumber()
		senderBalance, _ := strconv.Atoi(ps.SenderBalances[i])
		// Generate a random amount less than the sender balance
		// For simplicity, let's assume senderBalance is always >= 1

		ps.Amounts[i] = strconv.Itoa(rand2.Intn(senderBalance-1) + 1)
		amount, _ := strconv.Atoi(ps.Amounts[i])
		ps.ReceiverBalances[i] = strconv.Itoa(amount + rand2.Intn(amount) + 1)
		ps.Messages[i] = generateRandomString(255)
		ps.TransactionNonces[i] = generateRandomSmallNumber()
		transactionNonce, _ := strconv.Atoi(ps.TransactionNonces[i])

		// Generate an account nonce that is greater than the transaction nonce
		// For simplicity, let's assume accountNonce could be upto twice the transaction nonce
		ps.AccountNonces[i] = strconv.Itoa(transactionNonce + rand2.Intn(transactionNonce) + 1)
	}

	return ps
}

// GenerateUniqueRandomValues generates a set of unique random values for a PodStruct.
// It takes an integer count as an argument which determines the number of unique random values to be generated.
// It returns a PodStruct filled with unique random values.
//
// The function works as follows:
// - It initializes a new PodStruct with slices of strings, each of size count.
// - It then enters a loop that runs count times.
// - In each iteration, it generates random Ethereum addresses for 'From' and 'To' fields,
//   a random Ethereum transaction hash for 'TransactionHash' field,
//   random numbers for 'Amounts', 'SenderBalances', and 'ReceiverBalances' fields,
//   a random string for 'Messages' field,
//   and random small numbers for 'TransactionNonces' and 'AccountNonces' fields.
// - These generated values are then assigned to the respective fields in the PodStruct.
// - Finally, it returns the filled PodStruct.

func GenerateUniqueRandomValues(count int) types.PodStruct {
	ps := types.PodStruct{
		From:              make([]string, count),
		To:                make([]string, count),
		Amounts:           make([]string, count),
		TransactionHash:   make([]string, count),
		SenderBalances:    make([]string, count),
		ReceiverBalances:  make([]string, count),
		Messages:          make([]string, count),
		TransactionNonces: make([]string, count),
		AccountNonces:     make([]string, count),
	}

	for i := 0; i < count; i++ {
		randomFrom, _ := generateEthereumAddress()
		randomTo, _ := generateEthereumAddress()
		randomTxHash, _ := generateTransactionHash()

		ps.From[i] = randomFrom // Adjust the length as needed
		ps.To[i] = randomTo
		ps.Amounts[i] = "10" // Assuming amounts are shorter strings
		ps.TransactionHash[i] = randomTxHash
		ps.SenderBalances[i] = "20"
		ps.ReceiverBalances[i] = "10"
		ps.Messages[i] = generateRandomString(255)
		ps.TransactionNonces[i] = "7"
		ps.AccountNonces[i] = "8"
	}

	return ps
}
