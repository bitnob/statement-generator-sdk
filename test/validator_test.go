package test

import (
	"testing"
	"time"

	"github.com/bitnob/statement-generator-sdk/pkg/statements"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidator_ValidateStatementInput(t *testing.T) {
	tests := []struct {
		name        string
		input       statements.StatementInput
		wantErr     bool
		errContains string
	}{
		{
			name: "Valid input",
			input: statements.StatementInput{
				Account: statements.Account{
					Number:     "****1234",
					HolderName: "John Doe",
					Currency:   "USD",
				},
				Transactions: []statements.Transaction{
					{
						ID:          "TXN001",
						Date:        time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
						Description: "Test",
						Amount:      decimal.NewFromFloat(100),
						Type:        statements.Credit,
					},
				},
				PeriodStart:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				PeriodEnd:      time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC),
				OpeningBalance: decimal.NewFromFloat(1000),
			},
			wantErr: false,
		},
		{
			name: "Missing account number",
			input: statements.StatementInput{
				Account: statements.Account{
					HolderName: "John Doe",
					Currency:   "USD",
				},
				Transactions:   []statements.Transaction{},
				PeriodStart:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				PeriodEnd:      time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC),
				OpeningBalance: decimal.NewFromFloat(1000),
			},
			wantErr:     true,
			errContains: "account number is required",
		},
		{
			name: "Missing holder name",
			input: statements.StatementInput{
				Account: statements.Account{
					Number:   "****1234",
					Currency: "USD",
				},
				Transactions:   []statements.Transaction{},
				PeriodStart:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				PeriodEnd:      time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC),
				OpeningBalance: decimal.NewFromFloat(1000),
			},
			wantErr:     true,
			errContains: "holder name is required",
		},
		{
			name: "Missing currency",
			input: statements.StatementInput{
				Account: statements.Account{
					Number:     "****1234",
					HolderName: "John Doe",
				},
				Transactions:   []statements.Transaction{},
				PeriodStart:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				PeriodEnd:      time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC),
				OpeningBalance: decimal.NewFromFloat(1000),
			},
			wantErr:     true,
			errContains: "currency is required",
		},
		{
			name: "Invalid currency code",
			input: statements.StatementInput{
				Account: statements.Account{
					Number:     "****1234",
					HolderName: "John Doe",
					Currency:   "INVALID",
				},
				Transactions:   []statements.Transaction{},
				PeriodStart:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				PeriodEnd:      time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC),
				OpeningBalance: decimal.NewFromFloat(1000),
			},
			wantErr:     true,
			errContains: "invalid currency code",
		},
		{
			name: "Period end before period start",
			input: statements.StatementInput{
				Account: statements.Account{
					Number:     "****1234",
					HolderName: "John Doe",
					Currency:   "USD",
				},
				Transactions:   []statements.Transaction{},
				PeriodStart:    time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
				PeriodEnd:      time.Date(2024, 1, 1, 23, 59, 59, 0, time.UTC),
				OpeningBalance: decimal.NewFromFloat(1000),
			},
			wantErr:     true,
			errContains: "period end must be after period start",
		},
		{
			name: "Transaction outside period",
			input: statements.StatementInput{
				Account: statements.Account{
					Number:     "****1234",
					HolderName: "John Doe",
					Currency:   "USD",
				},
				Transactions: []statements.Transaction{
					{
						ID:          "TXN001",
						Date:        time.Date(2024, 2, 15, 10, 0, 0, 0, time.UTC), // Outside period
						Description: "Test",
						Amount:      decimal.NewFromFloat(100),
						Type:        statements.Credit,
					},
				},
				PeriodStart:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				PeriodEnd:      time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC),
				OpeningBalance: decimal.NewFromFloat(1000),
			},
			wantErr:     true,
			errContains: "transaction date outside statement period",
		},
		{
			name: "Missing transaction ID",
			input: statements.StatementInput{
				Account: statements.Account{
					Number:     "****1234",
					HolderName: "John Doe",
					Currency:   "USD",
				},
				Transactions: []statements.Transaction{
					{
						Date:        time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
						Description: "Test",
						Amount:      decimal.NewFromFloat(100),
						Type:        statements.Credit,
					},
				},
				PeriodStart:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				PeriodEnd:      time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC),
				OpeningBalance: decimal.NewFromFloat(1000),
			},
			wantErr:     true,
			errContains: "transaction ID is required",
		},
		{
			name: "Missing transaction description",
			input: statements.StatementInput{
				Account: statements.Account{
					Number:     "****1234",
					HolderName: "John Doe",
					Currency:   "USD",
				},
				Transactions: []statements.Transaction{
					{
						ID:     "TXN001",
						Date:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
						Amount: decimal.NewFromFloat(100),
						Type:   statements.Credit,
					},
				},
				PeriodStart:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				PeriodEnd:      time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC),
				OpeningBalance: decimal.NewFromFloat(1000),
			},
			wantErr:     true,
			errContains: "transaction description is required",
		},
		{
			name: "Mismatched transaction type and amount",
			input: statements.StatementInput{
				Account: statements.Account{
					Number:     "****1234",
					HolderName: "John Doe",
					Currency:   "USD",
				},
				Transactions: []statements.Transaction{
					{
						ID:          "TXN001",
						Date:        time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
						Description: "Test",
						Amount:      decimal.NewFromFloat(-100), // Negative amount
						Type:        statements.Credit,           // But marked as credit
					},
				},
				PeriodStart:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				PeriodEnd:      time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC),
				OpeningBalance: decimal.NewFromFloat(1000),
			},
			wantErr:     true,
			errContains: "credit transaction cannot have negative amount",
		},
		{
			name: "Invalid transaction type",
			input: statements.StatementInput{
				Account: statements.Account{
					Number:     "****1234",
					HolderName: "John Doe",
					Currency:   "USD",
				},
				Transactions: []statements.Transaction{
					{
						ID:          "TXN001",
						Date:        time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
						Description: "Test",
						Amount:      decimal.NewFromFloat(100),
						Type:        "invalid_type",
					},
				},
				PeriodStart:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				PeriodEnd:      time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC),
				OpeningBalance: decimal.NewFromFloat(1000),
			},
			wantErr:     true,
			errContains: "invalid transaction type",
		},
		{
			name: "Duplicate transaction IDs",
			input: statements.StatementInput{
				Account: statements.Account{
					Number:     "****1234",
					HolderName: "John Doe",
					Currency:   "USD",
				},
				Transactions: []statements.Transaction{
					{
						ID:          "TXN001",
						Date:        time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
						Description: "Test 1",
						Amount:      decimal.NewFromFloat(100),
						Type:        statements.Credit,
					},
					{
						ID:          "TXN001", // Duplicate ID
						Date:        time.Date(2024, 1, 16, 10, 0, 0, 0, time.UTC),
						Description: "Test 2",
						Amount:      decimal.NewFromFloat(200),
						Type:        statements.Credit,
					},
				},
				PeriodStart:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				PeriodEnd:      time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC),
				OpeningBalance: decimal.NewFromFloat(1000),
			},
			wantErr:     true,
			errContains: "duplicate transaction ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := statements.NewValidator()
			err := validator.ValidateStatementInput(tt.input)

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

func TestValidator_ValidateCurrencyCode(t *testing.T) {
	tests := []struct {
		code    string
		valid   bool
	}{
		{"USD", true},
		{"EUR", true},
		{"GBP", true},
		{"NGN", true},
		{"JPY", true},
		{"CHF", true},
		{"CAD", true},
		{"AUD", true},
		{"CNY", true},
		{"INR", true},
		{"BRL", true},
		{"ZAR", true},
		{"KRW", true},
		{"SGD", true},
		{"HKD", true},
		{"", false},
		{"US", false},
		{"USDD", false},
		{"123", false},
		{"usd", false}, // Must be uppercase
		{"XYZ", false}, // Not a valid ISO code
	}

	validator := statements.NewValidator()

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			err := validator.ValidateCurrencyCode(tt.code)
			if tt.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestValidator_ValidateAddress(t *testing.T) {
	tests := []struct {
		name    string
		address *statements.Address
		wantErr bool
	}{
		{
			name: "Valid complete address",
			address: &statements.Address{
				Line1:      "123 Main Street",
				Line2:      "Apt 4B",
				City:       "New York",
				State:      "NY",
				PostalCode: "10001",
				Country:    "USA",
			},
			wantErr: false,
		},
		{
			name: "Valid minimal address",
			address: &statements.Address{
				Line1:   "123 Main Street",
				City:    "New York",
				Country: "USA",
			},
			wantErr: false,
		},
		{
			name:    "Nil address is valid",
			address: nil,
			wantErr: false,
		},
		{
			name: "Missing line1",
			address: &statements.Address{
				City:    "New York",
				Country: "USA",
			},
			wantErr: true,
		},
		{
			name: "Missing city",
			address: &statements.Address{
				Line1:   "123 Main Street",
				Country: "USA",
			},
			wantErr: true,
		},
		{
			name: "Missing country",
			address: &statements.Address{
				Line1: "123 Main Street",
				City:  "New York",
			},
			wantErr: true,
		},
		{
			name: "Line1 too long",
			address: &statements.Address{
				Line1:   string(make([]byte, 201)), // 201 characters
				City:    "New York",
				Country: "USA",
			},
			wantErr: true,
		},
		{
			name: "Invalid postal code format",
			address: &statements.Address{
				Line1:      "123 Main Street",
				City:       "New York",
				PostalCode: "INVALID-POSTAL-CODE-123456789",
				Country:    "USA",
			},
			wantErr: true,
		},
	}

	validator := statements.NewValidator()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateAddress(tt.address)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidator_ValidateBalanceConsistency(t *testing.T) {
	tests := []struct {
		name           string
		transactions   []statements.Transaction
		openingBalance decimal.Decimal
		wantErr        bool
		errContains    string
	}{
		{
			name: "Consistent balances",
			transactions: []statements.Transaction{
				{
					ID:          "TXN001",
					Date:        time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
					Description: "Deposit",
					Amount:      decimal.NewFromFloat(500),
					Type:        statements.Credit,
					Balance:     decimalPtr(1500),
				},
				{
					ID:          "TXN002",
					Date:        time.Date(2024, 1, 16, 10, 0, 0, 0, time.UTC),
					Description: "Withdrawal",
					Amount:      decimal.NewFromFloat(-200),
					Type:        statements.Debit,
					Balance:     decimalPtr(1300),
				},
			},
			openingBalance: decimal.NewFromFloat(1000),
			wantErr:        false,
		},
		{
			name: "Inconsistent balance",
			transactions: []statements.Transaction{
				{
					ID:          "TXN001",
					Date:        time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
					Description: "Deposit",
					Amount:      decimal.NewFromFloat(500),
					Type:        statements.Credit,
					Balance:     decimalPtr(1600), // Should be 1500
				},
			},
			openingBalance: decimal.NewFromFloat(1000),
			wantErr:        true,
			errContains:    "balance mismatch",
		},
		{
			name: "No provided balances (should not error)",
			transactions: []statements.Transaction{
				{
					ID:          "TXN001",
					Date:        time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
					Description: "Deposit",
					Amount:      decimal.NewFromFloat(500),
					Type:        statements.Credit,
				},
				{
					ID:          "TXN002",
					Date:        time.Date(2024, 1, 16, 10, 0, 0, 0, time.UTC),
					Description: "Withdrawal",
					Amount:      decimal.NewFromFloat(-200),
					Type:        statements.Debit,
				},
			},
			openingBalance: decimal.NewFromFloat(1000),
			wantErr:        false,
		},
		{
			name: "Mixed provided and calculated balances",
			transactions: []statements.Transaction{
				{
					ID:          "TXN001",
					Date:        time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
					Description: "Deposit",
					Amount:      decimal.NewFromFloat(500),
					Type:        statements.Credit,
					Balance:     decimalPtr(1500),
				},
				{
					ID:          "TXN002",
					Date:        time.Date(2024, 1, 16, 10, 0, 0, 0, time.UTC),
					Description: "Withdrawal",
					Amount:      decimal.NewFromFloat(-200),
					Type:        statements.Debit,
					// No balance provided - should be calculated
				},
			},
			openingBalance: decimal.NewFromFloat(1000),
			wantErr:        false,
		},
	}

	validator := statements.NewValidator()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateBalanceConsistency(tt.transactions, tt.openingBalance)
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

// Helper function to create decimal pointer
func decimalPtr(value float64) *decimal.Decimal {
	d := decimal.NewFromFloat(value)
	return &d
}