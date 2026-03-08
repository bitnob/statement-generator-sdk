package test

import (
	"testing"
	"time"

	"github.com/statement-generator/sdk/pkg/statements"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStatementGenerator_BasicStatement(t *testing.T) {
	// Arrange
	generator := statements.New()

	account := statements.Account{
		Number:     "****1234",
		HolderName: "John Doe",
		Currency:   "USD",
		Address: &statements.Address{
			Line1:      "123 Main Street",
			Line2:      "Apt 4B",
			City:       "New York",
			State:      "NY",
			PostalCode: "10001",
			Country:    "USA",
		},
	}

	periodStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	periodEnd := time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC)
	openingBalance := decimal.NewFromFloat(1000.00)

	transactions := []statements.Transaction{
		{
			ID:          "TXN001",
			Date:        time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Description: "Direct Deposit - Salary",
			Amount:      decimal.NewFromFloat(5000.00),
			Type:        statements.Credit,
			Reference:   "SAL-2024-01",
		},
		{
			ID:          "TXN002",
			Date:        time.Date(2024, 1, 16, 14, 30, 0, 0, time.UTC),
			Description: "Rent Payment",
			Amount:      decimal.NewFromFloat(-1200.00),
			Type:        statements.Debit,
			Reference:   "RENT-2024-01",
		},
		{
			ID:          "TXN003",
			Date:        time.Date(2024, 1, 20, 9, 15, 0, 0, time.UTC),
			Description: "Grocery Store",
			Amount:      decimal.NewFromFloat(-150.50),
			Type:        statements.Debit,
		},
		{
			ID:          "TXN004",
			Date:        time.Date(2024, 1, 25, 16, 45, 0, 0, time.UTC),
			Description: "Freelance Payment",
			Amount:      decimal.NewFromFloat(800.00),
			Type:        statements.Credit,
			Reference:   "INV-2024-001",
		},
	}

	input := statements.StatementInput{
		Account:        account,
		Transactions:   transactions,
		PeriodStart:    periodStart,
		PeriodEnd:      periodEnd,
		OpeningBalance: openingBalance,
	}

	// Act
	statement, err := generator.Generate(input)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, statement)

	expectedClosingBalance := decimal.NewFromFloat(5449.50)
	expectedTotalCredits := decimal.NewFromFloat(5800.00)
	expectedTotalDebits := decimal.NewFromFloat(1350.50)

	assert.True(t, statement.ClosingBalance.Equal(expectedClosingBalance),
		"Expected closing balance %v, got %v", expectedClosingBalance, statement.ClosingBalance)
	assert.True(t, statement.TotalCredits.Equal(expectedTotalCredits),
		"Expected total credits %v, got %v", expectedTotalCredits, statement.TotalCredits)
	assert.True(t, statement.TotalDebits.Equal(expectedTotalDebits),
		"Expected total debits %v, got %v", expectedTotalDebits, statement.TotalDebits)
	assert.Equal(t, 4, statement.TransactionCount)
}

func TestStatementGenerator_EmptyTransactions(t *testing.T) {
	// Arrange
	generator := statements.New()

	account := statements.Account{
		Number:     "****5678",
		HolderName: "Jane Smith",
		Currency:   "USD",
	}

	openingBalance := decimal.NewFromFloat(1500.00)

	input := statements.StatementInput{
		Account:        account,
		Transactions:   []statements.Transaction{},
		PeriodStart:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		PeriodEnd:      time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC),
		OpeningBalance: openingBalance,
	}

	// Act
	statement, err := generator.Generate(input)

	// Assert
	require.NoError(t, err)
	assert.True(t, statement.ClosingBalance.Equal(openingBalance),
		"Closing balance should equal opening balance with no transactions")
	assert.True(t, statement.TotalCredits.IsZero())
	assert.True(t, statement.TotalDebits.IsZero())
	assert.Equal(t, 0, statement.TransactionCount)
}

func TestStatementGenerator_SingleTransaction(t *testing.T) {
	// Arrange
	generator := statements.New()

	account := statements.Account{
		Number:     "****9999",
		HolderName: "Bob Wilson",
		Currency:   "USD",
	}

	openingBalance := decimal.NewFromFloat(500.00)

	transactions := []statements.Transaction{
		{
			ID:          "TXN001",
			Date:        time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Description: "Payment Received",
			Amount:      decimal.NewFromFloat(250.00),
			Type:        statements.Credit,
		},
	}

	input := statements.StatementInput{
		Account:        account,
		Transactions:   transactions,
		PeriodStart:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		PeriodEnd:      time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC),
		OpeningBalance: openingBalance,
	}

	// Act
	statement, err := generator.Generate(input)

	// Assert
	require.NoError(t, err)

	expectedClosingBalance := decimal.NewFromFloat(750.00)
	assert.True(t, statement.ClosingBalance.Equal(expectedClosingBalance))
	assert.True(t, statement.TotalCredits.Equal(decimal.NewFromFloat(250.00)))
	assert.True(t, statement.TotalDebits.IsZero())
	assert.Equal(t, 1, statement.TransactionCount)
}

func TestStatementGenerator_AllCredits(t *testing.T) {
	// Arrange
	generator := statements.New()

	account := statements.Account{
		Number:     "****1111",
		HolderName: "Alice Cooper",
		Currency:   "USD",
	}

	openingBalance := decimal.NewFromFloat(100.00)

	transactions := []statements.Transaction{
		{
			ID:          "TXN001",
			Date:        time.Date(2024, 1, 5, 10, 0, 0, 0, time.UTC),
			Description: "Deposit 1",
			Amount:      decimal.NewFromFloat(500.00),
			Type:        statements.Credit,
		},
		{
			ID:          "TXN002",
			Date:        time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Description: "Deposit 2",
			Amount:      decimal.NewFromFloat(750.00),
			Type:        statements.Credit,
		},
		{
			ID:          "TXN003",
			Date:        time.Date(2024, 1, 25, 10, 0, 0, 0, time.UTC),
			Description: "Deposit 3",
			Amount:      decimal.NewFromFloat(250.00),
			Type:        statements.Credit,
		},
	}

	input := statements.StatementInput{
		Account:        account,
		Transactions:   transactions,
		PeriodStart:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		PeriodEnd:      time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC),
		OpeningBalance: openingBalance,
	}

	// Act
	statement, err := generator.Generate(input)

	// Assert
	require.NoError(t, err)

	expectedClosingBalance := decimal.NewFromFloat(1600.00)
	expectedTotalCredits := decimal.NewFromFloat(1500.00)

	assert.True(t, statement.ClosingBalance.Equal(expectedClosingBalance))
	assert.True(t, statement.TotalCredits.Equal(expectedTotalCredits))
	assert.True(t, statement.TotalDebits.IsZero())
	assert.Equal(t, 3, statement.TransactionCount)
}

func TestStatementGenerator_AllDebits(t *testing.T) {
	// Arrange
	generator := statements.New()

	account := statements.Account{
		Number:     "****2222",
		HolderName: "Charlie Brown",
		Currency:   "USD",
	}

	openingBalance := decimal.NewFromFloat(2000.00)

	transactions := []statements.Transaction{
		{
			ID:          "TXN001",
			Date:        time.Date(2024, 1, 5, 10, 0, 0, 0, time.UTC),
			Description: "Withdrawal 1",
			Amount:      decimal.NewFromFloat(-300.00),
			Type:        statements.Debit,
		},
		{
			ID:          "TXN002",
			Date:        time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Description: "Withdrawal 2",
			Amount:      decimal.NewFromFloat(-450.00),
			Type:        statements.Debit,
		},
		{
			ID:          "TXN003",
			Date:        time.Date(2024, 1, 25, 10, 0, 0, 0, time.UTC),
			Description: "Withdrawal 3",
			Amount:      decimal.NewFromFloat(-250.00),
			Type:        statements.Debit,
		},
	}

	input := statements.StatementInput{
		Account:        account,
		Transactions:   transactions,
		PeriodStart:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		PeriodEnd:      time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC),
		OpeningBalance: openingBalance,
	}

	// Act
	statement, err := generator.Generate(input)

	// Assert
	require.NoError(t, err)

	expectedClosingBalance := decimal.NewFromFloat(1000.00)
	expectedTotalDebits := decimal.NewFromFloat(1000.00)

	assert.True(t, statement.ClosingBalance.Equal(expectedClosingBalance))
	assert.True(t, statement.TotalCredits.IsZero())
	assert.True(t, statement.TotalDebits.Equal(expectedTotalDebits))
	assert.Equal(t, 3, statement.TransactionCount)
}

func TestStatementGenerator_LargeNumbers(t *testing.T) {
	// Arrange
	generator := statements.New()

	account := statements.Account{
		Number:     "****3333",
		HolderName: "Rich Person",
		Currency:   "USD",
	}

	openingBalance := decimal.NewFromFloat(1_000_000_000.00)

	transactions := []statements.Transaction{
		{
			ID:          "TXN001",
			Date:        time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Description: "Large Deposit",
			Amount:      decimal.NewFromFloat(5_000_000_000.00),
			Type:        statements.Credit,
		},
		{
			ID:          "TXN002",
			Date:        time.Date(2024, 1, 20, 10, 0, 0, 0, time.UTC),
			Description: "Large Withdrawal",
			Amount:      decimal.NewFromFloat(-2_000_000_000.00),
			Type:        statements.Debit,
		},
	}

	input := statements.StatementInput{
		Account:        account,
		Transactions:   transactions,
		PeriodStart:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		PeriodEnd:      time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC),
		OpeningBalance: openingBalance,
	}

	// Act
	statement, err := generator.Generate(input)

	// Assert
	require.NoError(t, err)

	expectedClosingBalance := decimal.NewFromFloat(4_000_000_000.00)
	assert.True(t, statement.ClosingBalance.Equal(expectedClosingBalance))
}

func TestStatementGenerator_MultipleSameDayTransactions(t *testing.T) {
	// Arrange
	generator := statements.New()

	account := statements.Account{
		Number:     "****4444",
		HolderName: "Busy Person",
		Currency:   "USD",
	}

	openingBalance := decimal.NewFromFloat(1000.00)
	sameDate := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	transactions := []statements.Transaction{
		{
			ID:          "TXN001",
			Date:        sameDate,
			Description: "Morning Coffee",
			Amount:      decimal.NewFromFloat(-5.50),
			Type:        statements.Debit,
		},
		{
			ID:          "TXN002",
			Date:        sameDate,
			Description: "Salary",
			Amount:      decimal.NewFromFloat(3000.00),
			Type:        statements.Credit,
		},
		{
			ID:          "TXN003",
			Date:        sameDate,
			Description: "Lunch",
			Amount:      decimal.NewFromFloat(-15.00),
			Type:        statements.Debit,
		},
		{
			ID:          "TXN004",
			Date:        sameDate,
			Description: "Gas",
			Amount:      decimal.NewFromFloat(-50.00),
			Type:        statements.Debit,
		},
	}

	input := statements.StatementInput{
		Account:        account,
		Transactions:   transactions,
		PeriodStart:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		PeriodEnd:      time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC),
		OpeningBalance: openingBalance,
	}

	// Act
	statement, err := generator.Generate(input)

	// Assert
	require.NoError(t, err)

	expectedClosingBalance := decimal.NewFromFloat(3929.50)
	expectedTotalCredits := decimal.NewFromFloat(3000.00)
	expectedTotalDebits := decimal.NewFromFloat(70.50)

	assert.True(t, statement.ClosingBalance.Equal(expectedClosingBalance))
	assert.True(t, statement.TotalCredits.Equal(expectedTotalCredits))
	assert.True(t, statement.TotalDebits.Equal(expectedTotalDebits))
	assert.Equal(t, 4, statement.TransactionCount)
}

func TestStatementGenerator_ZeroAmounts(t *testing.T) {
	// Arrange
	generator := statements.New()

	account := statements.Account{
		Number:     "****5555",
		HolderName: "Zero User",
		Currency:   "USD",
	}

	openingBalance := decimal.NewFromFloat(1000.00)

	transactions := []statements.Transaction{
		{
			ID:          "TXN001",
			Date:        time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Description: "Zero Transaction",
			Amount:      decimal.Zero,
			Type:        statements.Credit,
		},
		{
			ID:          "TXN002",
			Date:        time.Date(2024, 1, 20, 10, 0, 0, 0, time.UTC),
			Description: "Real Transaction",
			Amount:      decimal.NewFromFloat(100.00),
			Type:        statements.Credit,
		},
	}

	input := statements.StatementInput{
		Account:        account,
		Transactions:   transactions,
		PeriodStart:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		PeriodEnd:      time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC),
		OpeningBalance: openingBalance,
	}

	// Act
	statement, err := generator.Generate(input)

	// Assert
	require.NoError(t, err)

	expectedClosingBalance := decimal.NewFromFloat(1100.00)
	assert.True(t, statement.ClosingBalance.Equal(expectedClosingBalance))
	assert.True(t, statement.TotalCredits.Equal(decimal.NewFromFloat(100.00)))
}

func TestStatementGenerator_NegativeOpeningBalance(t *testing.T) {
	// Arrange
	generator := statements.New()

	account := statements.Account{
		Number:     "****6666",
		HolderName: "Overdraft User",
		Currency:   "USD",
	}

	openingBalance := decimal.NewFromFloat(-500.00)

	transactions := []statements.Transaction{
		{
			ID:          "TXN001",
			Date:        time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Description: "Deposit to Cover Overdraft",
			Amount:      decimal.NewFromFloat(1000.00),
			Type:        statements.Credit,
		},
		{
			ID:          "TXN002",
			Date:        time.Date(2024, 1, 20, 10, 0, 0, 0, time.UTC),
			Description: "Small Purchase",
			Amount:      decimal.NewFromFloat(-50.00),
			Type:        statements.Debit,
		},
	}

	input := statements.StatementInput{
		Account:        account,
		Transactions:   transactions,
		PeriodStart:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		PeriodEnd:      time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC),
		OpeningBalance: openingBalance,
	}

	// Act
	statement, err := generator.Generate(input)

	// Assert
	require.NoError(t, err)

	expectedClosingBalance := decimal.NewFromFloat(450.00)
	assert.True(t, statement.ClosingBalance.Equal(expectedClosingBalance))
}

func TestStatementGenerator_WithInstitution(t *testing.T) {
	// Arrange
	institution := &statements.Institution{
		Name: "Test Bank",
		Address: &statements.Address{
			Line1:      "100 Banking Street",
			City:       "Financial District",
			State:      "NY",
			PostalCode: "10004",
			Country:    "USA",
		},
		RegNumber: "REG123456",
		TaxID:     "TAX987654",
	}

	generator := statements.New(statements.WithInstitution(institution))

	account := statements.Account{
		Number:     "****7777",
		HolderName: "Institution Customer",
		Currency:   "USD",
	}

	input := statements.StatementInput{
		Account:        account,
		Transactions:   []statements.Transaction{},
		PeriodStart:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		PeriodEnd:      time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC),
		OpeningBalance: decimal.NewFromFloat(1000.00),
	}

	// Act
	statement, err := generator.Generate(input)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, statement)
}

func TestStatementGenerator_CustomLocale(t *testing.T) {
	// Arrange
	generator := statements.New(statements.WithLocale("en-GB"))

	account := statements.Account{
		Number:     "****8888",
		HolderName: "UK Customer",
		Currency:   "GBP",
	}

	transactions := []statements.Transaction{
		{
			ID:          "TXN001",
			Date:        time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Description: "UK Transaction",
			Amount:      decimal.NewFromFloat(100.50),
			Type:        statements.Credit,
		},
	}

	input := statements.StatementInput{
		Account:        account,
		Transactions:   transactions,
		PeriodStart:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		PeriodEnd:      time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC),
		OpeningBalance: decimal.NewFromFloat(500.00),
	}

	// Act
	statement, err := generator.Generate(input)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, statement)

	// Test that CSV output uses UK date format
	csv := statement.ToCSV()
	assert.Contains(t, csv, "15/01/2024") // UK date format
}

func TestStatementGenerator_QuickPDF(t *testing.T) {
	// Arrange
	transactions := []statements.Transaction{
		{
			ID:          "TXN001",
			Date:        time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Description: "Test Transaction",
			Amount:      decimal.NewFromFloat(100.00),
			Type:        statements.Credit,
		},
	}

	openingBalance := decimal.NewFromFloat(500.00)

	// Act
	pdf, err := statements.QuickPDF(transactions, openingBalance)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, pdf)
	assert.Greater(t, len(pdf), 0, "PDF should not be empty")
}

func TestStatementGenerator_QuickCSV(t *testing.T) {
	// Arrange
	transactions := []statements.Transaction{
		{
			ID:          "TXN001",
			Date:        time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Description: "Test Transaction",
			Amount:      decimal.NewFromFloat(100.00),
			Type:        statements.Credit,
		},
	}

	openingBalance := decimal.NewFromFloat(500.00)

	// Act
	csv, err := statements.QuickCSV(transactions, openingBalance)

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, csv)
	assert.Contains(t, csv, "Test Transaction")
	assert.Contains(t, csv, "100.00")
}