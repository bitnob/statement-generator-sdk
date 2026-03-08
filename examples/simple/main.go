package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bitnob/statement-generator-sdk/pkg/statements"
	"github.com/shopspring/decimal"
)

func main() {
	// Create account
	account := statements.Account{
		Number:     "****5678",
		HolderName: "Jane Smith",
		Currency:   "USD",
		Type:       "Checking",
		Address: &statements.Address{
			Line1:      "123 Main Street",
			Line2:      "Apt 4B",
			City:       "San Francisco",
			State:      "CA",
			PostalCode: "94105",
			Country:    "United States",
		},
	}

	// Create sample transactions
	transactions := []statements.Transaction{
		{
			ID:          "TXN001",
			Date:        time.Date(2024, 3, 1, 10, 0, 0, 0, time.UTC),
			Description: "Direct Deposit - Salary",
			Amount:      decimal.NewFromFloat(5000.00),
			Type:        statements.Credit,
			Reference:   "DD-2024-03",
		},
		{
			ID:          "TXN002",
			Date:        time.Date(2024, 3, 5, 14, 30, 0, 0, time.UTC),
			Description: "Rent Payment",
			Amount:      decimal.NewFromFloat(-1500.00),
			Type:        statements.Debit,
			Reference:   "RENT-MAR",
		},
		{
			ID:          "TXN003",
			Date:        time.Date(2024, 3, 10, 9, 15, 0, 0, time.UTC),
			Description: "Grocery Store",
			Amount:      decimal.NewFromFloat(-125.50),
			Type:        statements.Debit,
			Reference:   "POS-8923",
		},
		{
			ID:          "TXN004",
			Date:        time.Date(2024, 3, 15, 11, 0, 0, 0, time.UTC),
			Description: "Online Transfer",
			Amount:      decimal.NewFromFloat(-500.00),
			Type:        statements.Debit,
			Reference:   "TRF-0315",
		},
		{
			ID:          "TXN005",
			Date:        time.Date(2024, 3, 20, 16, 45, 0, 0, time.UTC),
			Description: "Interest Credit",
			Amount:      decimal.NewFromFloat(15.25),
			Type:        statements.Credit,
			Reference:   "INT-MAR",
		},
	}

	// Create institution with contact details
	institution := &statements.Institution{
		Name: "Example Bank",
		Address: &statements.Address{
			Line1:   "100 Financial Plaza",
			City:    "New York",
			State:   "NY",
			Country: "USA",
		},
		RegNumber:    "Member FDIC",
		ContactPhone: "1-800-EXAMPLE",
		ContactEmail: "support@examplebank.com",
		Website:      "www.examplebank.com",
		FooterText:   "Thank you for banking with us. For any questions, please contact our customer service team.",
	}

	// Create statement input
	input := statements.StatementInput{
		Account:        account,
		Transactions:   transactions,
		PeriodStart:    time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
		PeriodEnd:      time.Date(2024, 3, 31, 23, 59, 59, 0, time.UTC),
		OpeningBalance: decimal.NewFromFloat(2500.00),
		Institution:    institution,
	}

	// Generate statement
	generator := statements.New()
	statement, err := generator.Generate(input)
	if err != nil {
		log.Fatal("Failed to generate statement:", err)
	}

	// Output summary
	fmt.Println("Statement Summary")
	fmt.Println("=================")
	fmt.Printf("Account: %s\n", account.Number)
	fmt.Printf("Period: %s to %s\n",
		input.PeriodStart.Format("Jan 02, 2006"),
		input.PeriodEnd.Format("Jan 02, 2006"))
	fmt.Printf("Opening Balance: $%.2f\n", input.OpeningBalance.InexactFloat64())
	fmt.Printf("Closing Balance: $%.2f\n", statement.ClosingBalance.InexactFloat64())
	fmt.Printf("Transactions: %d\n", statement.TransactionCount)

	// Generate PDF
	pdfBytes, err := statement.ToPDF()
	if err != nil {
		log.Fatal("Failed to generate PDF:", err)
	}

	// Save PDF
	filename := "example_statement.pdf"
	err = os.WriteFile(filename, pdfBytes, 0644)
	if err != nil {
		log.Fatal("Failed to save PDF:", err)
	}

	fmt.Printf("\n✅ Statement saved as %s\n", filename)
}