package components

import (
	"strings"
	"testing"
)

// GenerateRandomWithFavour test functions below

func TestGenerateRandomWithFavour_HappyPath(t *testing.T) {
	lowerBound := 1
	upperBound := 10
	favourableSet := [2]int{6, 8}
	favourableProbability := 0.9

	result := GenerateRandomWithFavour(lowerBound, upperBound, favourableSet, favourableProbability)
	t.Logf("Generated number: %v", result)
	if result < lowerBound || result > upperBound {
		t.Errorf("Generated number is out of bounds: got %v, expected between %v and %v", result, lowerBound, upperBound)
	}
}

func TestGenerateRandomWithFavour_InvalidParameters(t *testing.T) {
	lowerBound := 10
	upperBound := 1
	favourableSet := [2]int{3, 5}
	favourableProbability := 0.5

	result := GenerateRandomWithFavour(lowerBound, upperBound, favourableSet, favourableProbability)

	if result != 0 {
		t.Errorf("Expected 0 for invalid parameters, got %v", result)
	}
}

func TestGenerateRandomWithFavour_InvalidFavourableSet(t *testing.T) {
	lowerBound := 1
	upperBound := 10
	favourableSet := [2]int{11, 15}
	favourableProbability := 0.5

	result := GenerateRandomWithFavour(lowerBound, upperBound, favourableSet, favourableProbability)

	if result != 0 {
		t.Errorf("Expected 0 for invalid favourable set, got %v", result)
	}
}

func TestGenerateRandomWithFavour_FavourableProbabilityOutOfRange(t *testing.T) {
	lowerBound := 1
	upperBound := 10
	favourableSet := [2]int{3, 5}
	favourableProbability := 1.5

	result := GenerateRandomWithFavour(lowerBound, upperBound, favourableSet, favourableProbability)

	if result != 0 {
		t.Errorf("Expected 0 for out of range favourable probability, got %v", result)
	}
}

func TestGenerateAddresses_GeneratesCorrectNumberOfAddresses(t *testing.T) {
	n := 5
	addresses := GenerateAddresses(n)

	t.Logf("Generated addresses: %v", addresses)
	if len(addresses) != n {
		t.Errorf("Expected %v addresses, got %v", n, len(addresses))
	}
}

func TestGenerateAddresses_GeneratesAddressesWithCorrectPrefix(t *testing.T) {
	n := 1
	addresses := GenerateAddresses(n)

	t.Logf("Generated address: %v", addresses[0])
	if !strings.HasPrefix(addresses[0], "air1") {
		t.Errorf("Expected address to start with 'air1', got %v", addresses[0])
	}
}

func TestGenerateAddresses_GeneratesAddressesWithCorrectLength(t *testing.T) {
	n := 1
	addresses := GenerateAddresses(n)

	t.Logf("Generated address: %v", addresses[0])
	if len(addresses[0]) != 42 {
		t.Errorf("Expected address length to be 42, got %v", len(addresses[0]))
	}
}

func TestGenerateAddresses_GeneratesUniqueAddresses(t *testing.T) {
	n := 1000
	addresses := GenerateAddresses(n)
	t.Logf("Generated addresses: %v", addresses)
	addressMap := make(map[string]bool)
	for _, address := range addresses {
		if addressMap[address] {
			t.Errorf("Duplicate address found: %v", address)
		}
		addressMap[address] = true
	}
}

func TestGenerateAddresses_GeneratesAddressesWithAllowedCharactersOnly(t *testing.T) {
	n := 1
	allowedChars := "abcdefghijklmnopqrstuvwxyz0123456789"
	addresses := GenerateAddresses(n)

	t.Logf("Generated address: %v", addresses[0])
	for _, char := range addresses[0][4:] {
		if !strings.ContainsRune(allowedChars, char) {
			t.Errorf("Address contains disallowed character: %v", char)
		}
	}
}
