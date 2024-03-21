package components

import (
	"testing"
)

// TestGenerateUniqueRandomValues is a table driven test for testing GenerateUniqueRandomValues function.
func TestGenerateUniqueRandomValues(t *testing.T) {
	tests := []struct {
		name  string
		count int
	}{
		{name: "Zero", count: 0},
		{name: "Three", count: 3},
		{name: "Hundred", count: 100},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := GenerateUniqueRandomValues(tt.count)
			t.Logf("Generated struct: %v", ps)

			if len(ps.From) != tt.count || len(ps.To) != tt.count || len(ps.Amounts) != tt.count ||
				len(ps.TransactionHash) != tt.count || len(ps.SenderBalances) != tt.count || len(ps.ReceiverBalances) != tt.count ||
				len(ps.Messages) != tt.count || len(ps.TransactionNonces) != tt.count || len(ps.AccountNonces) != tt.count {
				t.Errorf("Unexpected length of generated fields in struct, got %d, want %d", len(ps.From), tt.count)
			}
		})
	}
}
