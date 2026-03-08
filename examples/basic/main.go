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
	// Create account with address for proof-of-address
	account := statements.Account{
		Number:     "****1234",
		HolderName: "John Doe",
		Currency:   "USD",
		Type:       "Checking",
		Address: &statements.Address{
			Line1:      "123 Main Street",
			Line2:      "Apt 4B",
			City:       "New York",
			State:      "NY",
			PostalCode: "10001",
			Country:    "USA",
		},
	}

	// Create institution
	institution := &statements.Institution{
		Name: "Example Bank",
		Address: &statements.Address{
			Line1:   "100 Financial Plaza",
			City:    "New York",
			State:   "NY",
			Country: "USA",
		},
		RegNumber: "REG-123456",
	}

	// Create transactions
	transactions := []statements.Transaction{
		{
			ID:          "TXN001",
			Date:        time.Date(2024, 1, 5, 10, 0, 0, 0, time.UTC),
			Description: "Opening Deposit",
			Amount:      decimal.NewFromFloat(1000.00),
			Type:        statements.Credit,
			Reference:   "DEP-001",
		},
		{
			ID:          "TXN002",
			Date:        time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC),
			Description: "Direct Deposit - Salary",
			Amount:      decimal.NewFromFloat(5000.00),
			Type:        statements.Credit,
			Reference:   "SAL-2024-01",
		},
		{
			ID:          "TXN003",
			Date:        time.Date(2024, 1, 16, 9, 15, 0, 0, time.UTC),
			Description: "Rent Payment",
			Amount:      decimal.NewFromFloat(-1200.00),
			Type:        statements.Debit,
			Reference:   "RENT-2024-01",
		},
		{
			ID:          "TXN004",
			Date:        time.Date(2024, 1, 20, 11, 45, 0, 0, time.UTC),
			Description: "Grocery Store",
			Amount:      decimal.NewFromFloat(-150.50),
			Type:        statements.Debit,
		},
		{
			ID:          "TXN005",
			Date:        time.Date(2024, 1, 22, 16, 20, 0, 0, time.UTC),
			Description: "Online Transfer",
			Amount:      decimal.NewFromFloat(-500.00),
			Type:        statements.Debit,
			Reference:   "TRF-2024-001",
		},
		{
			ID:          "TXN006",
			Date:        time.Date(2024, 1, 25, 10, 0, 0, 0, time.UTC),
			Description: "Freelance Payment",
			Amount:      decimal.NewFromFloat(800.00),
			Type:        statements.Credit,
			Reference:   "INV-2024-001",
		},
		{
			ID:          "TXN007",
			Date:        time.Date(2024, 1, 28, 15, 30, 0, 0, time.UTC),
			Description: "Utilities",
			Amount:      decimal.NewFromFloat(-200.00),
			Type:        statements.Debit,
			Reference:   "UTIL-2024-01",
		},
	}

	// Create statement input
	input := statements.StatementInput{
		Account:        account,
		Transactions:   transactions,
		PeriodStart:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		PeriodEnd:      time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC),
		OpeningBalance: decimal.NewFromFloat(500.00),
		Institution:    institution,
	}

	// Create generator
	generator := statements.New(
		statements.WithLocale("en-US"),
	)

	// Generate statement
	statement, err := generator.Generate(input)
	if err != nil {
		log.Fatal("Failed to generate statement:", err)
	}

	// Output summary
	fmt.Println("=== Statement Summary ===")
	fmt.Printf("Account Holder: %s\n", account.HolderName)
	fmt.Printf("Account Number: %s\n", account.Number)
	fmt.Printf("Period: %s to %s\n",
		input.PeriodStart.Format("2006-01-02"),
		input.PeriodEnd.Format("2006-01-02"))
	fmt.Printf("Opening Balance: %s\n", input.OpeningBalance.StringFixed(2))
	fmt.Printf("Total Credits: %s\n", statement.TotalCredits.StringFixed(2))
	fmt.Printf("Total Debits: %s\n", statement.TotalDebits.StringFixed(2))
	fmt.Printf("Closing Balance: %s\n", statement.ClosingBalance.StringFixed(2))
	fmt.Printf("Transaction Count: %d\n", statement.TransactionCount)

	// Generate PDF
	fmt.Println("\n=== Generating PDF ===")
	pdfBytes, err := statement.ToPDF()
	if err != nil {
		log.Fatal("Failed to generate PDF:", err)
	}

	err = os.WriteFile("statement.pdf", pdfBytes, 0644)
	if err != nil {
		log.Fatal("Failed to save PDF:", err)
	}
	fmt.Println("PDF saved as statement.pdf")

	// Generate CSV
	fmt.Println("\n=== Generating CSV ===")
	csv := statement.ToCSV()
	err = os.WriteFile("statement.csv", []byte(csv), 0644)
	if err != nil {
		log.Fatal("Failed to save CSV:", err)
	}
	fmt.Println("CSV saved as statement.csv")

	// Generate HTML
	fmt.Println("\n=== Generating HTML ===")
	html := statement.ToHTML()
	err = os.WriteFile("statement.html", []byte(html), 0644)
	if err != nil {
		log.Fatal("Failed to save HTML:", err)
	}
	fmt.Println("HTML saved as statement.html")

	fmt.Println("\n✅ Statement generation complete!")
}