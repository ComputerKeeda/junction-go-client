package components

import (
	"context"
	"fmt"
	"github.com/ComputerKeeda/junction-go-client/types"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosclient"
	"math/rand"
	"time"
)

const (
	JunctionTTCRPC      = "http://34.131.189.98:26657"
	JunctionAPI         = "http://34.131.189.98:1317"
	BatchSize           = 25
	SettlementClientRPC = "http://127.0.0.1:8080"
	DaClientRPC         = "http://127.0.0.1:5050/celestia"
	RpcAUTH             = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBbGxvdyI6WyJwdWJsaWMiLCJyZWFkIiwid3JpdGUiLCJhZG1pbiJdfQ.F0ZCzFpLy0XhURMxkJnSvcGiU3s0vkW0Q0pniqwwJns"
	DaCelRPC            = "http://localhost:26658/"
	DaType              = "mock" // "mock"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randomString(n int, allowedChars string) string {
	s := make([]rune, n)
	for i := range s {
		s[i] = rune(allowedChars[rand.Intn(len(allowedChars))])
	}
	return string(s)
}

// GenerateAddresses Generates n address strings with the "air1" prefix and specified format.
func GenerateAddresses(n int) []string {
	addresses := make([]string, n)
	// Define the characters that can appear in the addresses.
	allowedChars := "abcdefghijklmnopqrstuvwxyz0123456789"
	// Adjust the length according to your specific requirements.
	for i := 0; i < n; i++ {
		// "air1" + 37 characters to match your example length.
		addresses[i] = "air1" + randomString(38, allowedChars)
	}
	return addresses
}

func GenerateRandomWithFavour(lowerBound, upperBound int, favourableSet [2]int, favourableProbability float64) int {
	if lowerBound > upperBound || favourableProbability < 0 || favourableProbability > 1 {
		fmt.Println("Invalid parameters")
		return 0
	}

	// Calculate total range and the favourable range
	totalRange := upperBound - lowerBound + 1
	favourableRange := favourableSet[1] - favourableSet[0] + 1

	if favourableRange <= 0 || favourableRange > totalRange {
		fmt.Println("Invalid favourable set")
		return 0
	}

	// Check if the favourable set is within the total range
	if favourableSet[0] < lowerBound || favourableSet[1] > upperBound || favourableRange <= 0 {
		fmt.Println("Invalid favourable set")
		return 0
	}

	// Calculate the number of favourable outcomes based on the probability
	favourableOutcomes := int(favourableProbability * float64(totalRange))
	if favourableOutcomes < favourableRange {
		favourableOutcomes = favourableRange
	}

	// Generate a random number and adjust for favourable outcomes
	randNum := rand.Intn(totalRange)
	if randNum < favourableOutcomes {
		// Map the first `favourableOutcomes` to the favourable range
		randNum = randNum%favourableRange + favourableSet[0]
	} else {
		// Adjust the random number to exclude the favourable range and map to the rest of the range
		randNum = randNum%favourableOutcomes + lowerBound
		if randNum >= favourableSet[0] && randNum <= favourableSet[1] {
			randNum = favourableSet[1] + 1 + (randNum - favourableSet[0])
		}
	}

	return randNum
}

func GetQueryClient() types.QueryClient {
	ctx := context.Background()
	// getting the account and creating client codes --> End
	client, err := cosmosclient.New(ctx, cosmosclient.WithNodeAddress(JunctionTTCRPC))
	if err != nil {
		Logger.Error("Error creating account client")
	}

	queryClient := types.NewQueryClient(client.Context())

	return queryClient
}
