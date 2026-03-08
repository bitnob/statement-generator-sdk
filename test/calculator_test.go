package test

import (
	"testing"
	"time"

	"github.com/statement-generator/sdk/pkg/statements"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalculator_CalculateBalances(t *testing.T) {
	tests := []struct {
		name               string
		transactions       []statements.Transaction
		openingBalance     decimal.Decimal
		expectedBalances   []decimal.Decimal
		expectedClosing    decimal.Decimal
		expectedCredits    decimal.Decimal
		expectedDebits     decimal.Decimal
	}{
		{
			name: "Basic calculation",
			transactions: []statements.Transaction{
				{
					ID:          "TXN001",
					Date:        time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
					Description: "Credit",
					Amount:      decimal.NewFromFloat(500),
					Type:        statements.Credit,
				},
				{
					ID:          "TXN002",
					Date:        time.Date(2024, 1, 16, 10, 0, 0, 0, time.UTC),
					Description: "Debit",
					Amount:      decimal.NewFromFloat(-200),
					Type:        statements.Debit,
				},
			},
			openingBalance:   decimal.NewFromFloat(1000),
			expectedBalances: []decimal.Decimal{decimal.NewFromFloat(1500), decimal.NewFromFloat(1300)},
			expectedClosing:  decimal.NewFromFloat(1300),
			expectedCredits:  decimal.NewFromFloat(500),
			expectedDebits:   decimal.NewFromFloat(200),
		},
		{
			name:             "No transactions",
			transactions:     []statements.Transaction{},
			openingBalance:   decimal.NewFromFloat(1000),
			expectedBalances: []decimal.Decimal{},
			expectedClosing:  decimal.NewFromFloat(1000),
			expectedCredits:  decimal.Zero,
			expectedDebits:   decimal.Zero,
		},
		{
			name: "All credits",
			transactions: []statements.Transaction{
				{
					ID:          "TXN001",
					Date:        time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
					Description: "Credit 1",
					Amount:      decimal.NewFromFloat(100),
					Type:        statements.Credit,
				},
				{
					ID:          "TXN002",
					Date:        time.Date(2024, 1, 16, 10, 0, 0, 0, time.UTC),
					Description: "Credit 2",
					Amount:      decimal.NewFromFloat(200),
					Type:        statements.Credit,
				},
				{
					ID:          "TXN003",
					Date:        time.Date(2024, 1, 17, 10, 0, 0, 0, time.UTC),
					Description: "Credit 3",
					Amount:      decimal.NewFromFloat(300),
					Type:        statements.Credit,
				},
			},
			openingBalance: decimal.NewFromFloat(1000),
			expectedBalances: []decimal.Decimal{
				decimal.NewFromFloat(1100),
				decimal.NewFromFloat(1300),
				decimal.NewFromFloat(1600),
			},
			expectedClosing: decimal.NewFromFloat(1600),
			expectedCredits: decimal.NewFromFloat(600),
			expectedDebits:  decimal.Zero,
		},
		{
			name: "All debits",
			transactions: []statements.Transaction{
				{
					ID:          "TXN001",
					Date:        time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
					Description: "Debit 1",
					Amount:      decimal.NewFromFloat(-100),
					Type:        statements.Debit,
				},
				{
					ID:          "TXN002",
					Date:        time.Date(2024, 1, 16, 10, 0, 0, 0, time.UTC),
					Description: "Debit 2",
					Amount:      decimal.NewFromFloat(-200),
					Type:        statements.Debit,
				},
				{
					ID:          "TXN003",
					Date:        time.Date(2024, 1, 17, 10, 0, 0, 0, time.UTC),
					Description: "Debit 3",
					Amount:      decimal.NewFromFloat(-300),
					Type:        statements.Debit,
				},
			},
			openingBalance: decimal.NewFromFloat(1000),
			expectedBalances: []decimal.Decimal{
				decimal.NewFromFloat(900),
				decimal.NewFromFloat(700),
				decimal.NewFromFloat(400),
			},
			expectedClosing: decimal.NewFromFloat(400),
			expectedCredits: decimal.Zero,
			expectedDebits:  decimal.NewFromFloat(600),
		},
		{
			name: "Negative opening balance",
			transactions: []statements.Transaction{
				{
					ID:          "TXN001",
					Date:        time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
					Description: "Deposit to cover overdraft",
					Amount:      decimal.NewFromFloat(1000),
					Type:        statements.Credit,
				},
			},
			openingBalance:   decimal.NewFromFloat(-500),
			expectedBalances: []decimal.Decimal{decimal.NewFromFloat(500)},
			expectedClosing:  decimal.NewFromFloat(500),
			expectedCredits:  decimal.NewFromFloat(1000),
			expectedDebits:   decimal.Zero,
		},
		{
			name: "Zero amount transaction",
			transactions: []statements.Transaction{
				{
					ID:          "TXN001",
					Date:        time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
					Description: "Zero transaction",
					Amount:      decimal.Zero,
					Type:        statements.Credit,
				},
				{
					ID:          "TXN002",
					Date:        time.Date(2024, 1, 16, 10, 0, 0, 0, time.UTC),
					Description: "Real transaction",
					Amount:      decimal.NewFromFloat(100),
					Type:        statements.Credit,
				},
			},
			openingBalance:   decimal.NewFromFloat(1000),
			expectedBalances: []decimal.Decimal{decimal.NewFromFloat(1000), decimal.NewFromFloat(1100)},
			expectedClosing:  decimal.NewFromFloat(1100),
			expectedCredits:  decimal.NewFromFloat(100),
			expectedDebits:   decimal.Zero,
		},
		{
			name: "High precision decimals",
			transactions: []statements.Transaction{
				{
					ID:          "TXN001",
					Date:        time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
					Description: "Precise credit",
					Amount:      decimal.NewFromFloat(100.123456789),
					Type:        statements.Credit,
				},
				{
					ID:          "TXN002",
					Date:        time.Date(2024, 1, 16, 10, 0, 0, 0, time.UTC),
					Description: "Precise debit",
					Amount:      decimal.NewFromFloat(-50.987654321),
					Type:        statements.Debit,
				},
			},
			openingBalance: decimal.NewFromFloat(1000.111111111),
			expectedBalances: []decimal.Decimal{
				decimal.RequireFromString("1100.234567900"),
				decimal.RequireFromString("1049.246913579"),
			},
			expectedClosing: decimal.RequireFromString("1049.246913579"),
			expectedCredits: decimal.NewFromFloat(100.123456789),
			expectedDebits:  decimal.NewFromFloat(50.987654321),
		},
		{
			name: "Very large numbers",
			transactions: []statements.Transaction{
				{
					ID:          "TXN001",
					Date:        time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
					Description: "Large credit",
					Amount:      decimal.NewFromFloat(9999999999.99),
					Type:        statements.Credit,
				},
				{
					ID:          "TXN002",
					Date:        time.Date(2024, 1, 16, 10, 0, 0, 0, time.UTC),
					Description: "Large debit",
					Amount:      decimal.NewFromFloat(-5555555555.55),
					Type:        statements.Debit,
				},
			},
			openingBalance: decimal.NewFromFloat(1111111111.11),
			expectedBalances: []decimal.Decimal{
				decimal.NewFromFloat(11111111111.10),
				decimal.NewFromFloat(5555555555.55),
			},
			expectedClosing: decimal.NewFromFloat(5555555555.55),
			expectedCredits: decimal.NewFromFloat(9999999999.99),
			expectedDebits:  decimal.NewFromFloat(5555555555.55),
		},
	}

	calculator := statements.NewCalculator()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := calculator.CalculateBalances(tt.transactions, tt.openingBalance)
			require.NoError(t, err)

			// Check closing balance
			assert.True(t, result.ClosingBalance.Equal(tt.expectedClosing),
				"Expected closing balance %v, got %v", tt.expectedClosing, result.ClosingBalance)

			// Check total credits
			assert.True(t, result.TotalCredits.Equal(tt.expectedCredits),
				"Expected total credits %v, got %v", tt.expectedCredits, result.TotalCredits)

			// Check total debits
			assert.True(t, result.TotalDebits.Equal(tt.expectedDebits),
				"Expected total debits %v, got %v", tt.expectedDebits, result.TotalDebits)

			// Check individual balances
			assert.Equal(t, len(tt.expectedBalances), len(result.Balances),
				"Expected %d balances, got %d", len(tt.expectedBalances), len(result.Balances))

			for i, expected := range tt.expectedBalances {
				assert.True(t, result.Balances[i].Equal(expected),
					"Balance %d: expected %v, got %v", i, expected, result.Balances[i])
			}
		})
	}
}

func TestCalculator_ValidateProvidedBalances(t *testing.T) {
	tests := []struct {
		name           string
		transactions   []statements.Transaction
		openingBalance decimal.Decimal
		wantErr        bool
		errContains    string
	}{
		{
			name: "All balances correct",
			transactions: []statements.Transaction{
				{
					ID:          "TXN001",
					Date:        time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
					Description: "Credit",
					Amount:      decimal.NewFromFloat(500),
					Type:        statements.Credit,
					Balance:     decimalPtr(1500),
				},
				{
					ID:          "TXN002",
					Date:        time.Date(2024, 1, 16, 10, 0, 0, 0, time.UTC),
					Description: "Debit",
					Amount:      decimal.NewFromFloat(-200),
					Type:        statements.Debit,
					Balance:     decimalPtr(1300),
				},
			},
			openingBalance: decimal.NewFromFloat(1000),
			wantErr:        false,
		},
		{
			name: "First balance incorrect",
			transactions: []statements.Transaction{
				{
					ID:          "TXN001",
					Date:        time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
					Description: "Credit",
					Amount:      decimal.NewFromFloat(500),
					Type:        statements.Credit,
					Balance:     decimalPtr(1600), // Should be 1500
				},
			},
			openingBalance: decimal.NewFromFloat(1000),
			wantErr:        true,
			errContains:    "balance mismatch for transaction TXN001",
		},
		{
			name: "Second balance incorrect",
			transactions: []statements.Transaction{
				{
					ID:          "TXN001",
					Date:        time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
					Description: "Credit",
					Amount:      decimal.NewFromFloat(500),
					Type:        statements.Credit,
					Balance:     decimalPtr(1500),
				},
				{
					ID:          "TXN002",
					Date:        time.Date(2024, 1, 16, 10, 0, 0, 0, time.UTC),
					Description: "Debit",
					Amount:      decimal.NewFromFloat(-200),
					Type:        statements.Debit,
					Balance:     decimalPtr(1400), // Should be 1300
				},
			},
			openingBalance: decimal.NewFromFloat(1000),
			wantErr:        true,
			errContains:    "balance mismatch for transaction TXN002",
		},
		{
			name: "Mixed provided and calculated",
			transactions: []statements.Transaction{
				{
					ID:          "TXN001",
					Date:        time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
					Description: "Credit",
					Amount:      decimal.NewFromFloat(500),
					Type:        statements.Credit,
					Balance:     decimalPtr(1500),
				},
				{
					ID:          "TXN002",
					Date:        time.Date(2024, 1, 16, 10, 0, 0, 0, time.UTC),
					Description: "Debit",
					Amount:      decimal.NewFromFloat(-200),
					Type:        statements.Debit,
					// No balance provided
				},
				{
					ID:          "TXN003",
					Date:        time.Date(2024, 1, 17, 10, 0, 0, 0, time.UTC),
					Description: "Credit",
					Amount:      decimal.NewFromFloat(100),
					Type:        statements.Credit,
					Balance:     decimalPtr(1400),
				},
			},
			openingBalance: decimal.NewFromFloat(1000),
			wantErr:        false,
		},
		{
			name: "Rounding tolerance",
			transactions: []statements.Transaction{
				{
					ID:          "TXN001",
					Date:        time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
					Description: "Credit with rounding",
					Amount:      decimal.NewFromFloat(100.00),
					Type:        statements.Credit,
					Balance:     decimalPtr(1100.0000001), // Slight rounding difference
				},
			},
			openingBalance: decimal.NewFromFloat(1000),
			wantErr:        false, // Should tolerate tiny rounding differences
		},
	}

	calculator := statements.NewCalculator()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := calculator.ValidateProvidedBalances(tt.transactions, tt.openingBalance)
			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestCalculator_SortTransactions(t *testing.T) {
	transactions := []statements.Transaction{
		{
			ID:   "TXN003",
			Date: time.Date(2024, 1, 20, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:   "TXN001",
			Date: time.Date(2024, 1, 10, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:   "TXN004",
			Date: time.Date(2024, 1, 15, 14, 0, 0, 0, time.UTC),
		},
		{
			ID:   "TXN002",
			Date: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		},
	}

	calculator := statements.NewCalculator()
	sorted := calculator.SortTransactions(transactions)

	// Check that transactions are sorted by date
	assert.Equal(t, "TXN001", sorted[0].ID)
	assert.Equal(t, "TXN002", sorted[1].ID)
	assert.Equal(t, "TXN004", sorted[2].ID)
	assert.Equal(t, "TXN003", sorted[3].ID)

	// Ensure original slice is not modified
	assert.Equal(t, "TXN003", transactions[0].ID)
}