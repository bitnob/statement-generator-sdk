package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/statement-generator/sdk/pkg/statements"
	"github.com/shopspring/decimal"
)

func main() {
	// Create sample data
	account := statements.Account{
		Number:     "****5678",
		HolderName: "Template Test User",
		Currency:   "USD",
		Type:       "Checking",
		Address: &statements.Address{
			Line1:      "456 Design Street",
			City:       "San Francisco",
			State:      "CA",
			PostalCode: "94105",
			Country:    "USA",
		},
	}

	transactions := []statements.Transaction{
		{
			ID:          "TEST001",
			Date:        time.Date(2024, 3, 1, 10, 0, 0, 0, time.UTC),
			Description: "Initial Deposit",
			Amount:      decimal.NewFromFloat(5000),
			Type:        statements.Credit,
			Reference:   "INIT-001",
		},
		{
			ID:          "TEST002",
			Date:        time.Date(2024, 3, 5, 14, 30, 0, 0, time.UTC),
			Description: "Payment Received",
			Amount:      decimal.NewFromFloat(1500),
			Type:        statements.Credit,
			Reference:   "PAY-002",
		},
		{
			ID:          "TEST003",
			Date:        time.Date(2024, 3, 10, 9, 15, 0, 0, time.UTC),
			Description: "Office Supplies",
			Amount:      decimal.NewFromFloat(-250),
			Type:        statements.Debit,
			Reference:   "EXP-003",
		},
		{
			ID:          "TEST004",
			Date:        time.Date(2024, 3, 15, 11, 0, 0, 0, time.UTC),
			Description: "Software Subscription",
			Amount:      decimal.NewFromFloat(-99),
			Type:        statements.Debit,
			Reference:   "SUB-004",
		},
		{
			ID:          "TEST005",
			Date:        time.Date(2024, 3, 20, 16, 45, 0, 0, time.UTC),
			Description: "Client Payment",
			Amount:      decimal.NewFromFloat(3000),
			Type:        statements.Credit,
			Reference:   "INV-005",
		},
	}

	input := statements.StatementInput{
		Account:        account,
		Transactions:   transactions,
		PeriodStart:    time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
		PeriodEnd:      time.Date(2024, 3, 31, 23, 59, 59, 0, time.UTC),
		OpeningBalance: decimal.NewFromFloat(1000),
		Institution: &statements.Institution{
			Name: "Demo Bank",
			Address: &statements.Address{
				Line1:   "100 Finance Ave",
				City:    "New York",
				State:   "NY",
				Country: "USA",
			},
			RegNumber: "FDIC #12345",
		},
	}

	fmt.Println("=== Testing Statement Template Styles ===\n")

	// Test 1: Default (Minimalist) Template
	fmt.Println("1. GENERATING WITH DEFAULT (MINIMALIST) TEMPLATE")
	fmt.Println("   - Black and white design")
	fmt.Println("   - Alternating grey/white rows")
	fmt.Println("   - Clean, professional look")

	generator := statements.New()
	statement, err := generator.Generate(input)
	if err != nil {
		log.Fatal("Failed to generate statement:", err)
	}

	pdfBytes, _ := statement.ToPDF()
	os.WriteFile("statement_minimalist.pdf", pdfBytes, 0644)
	fmt.Println("   ✅ Generated: statement_minimalist.pdf\n")

	// Test 2: Enhanced Template
	fmt.Println("2. GENERATING WITH ENHANCED TEMPLATE")
	fmt.Println("   - Colored headers and text")
	fmt.Println("   - Blue headers, red/green for debits/credits")
	fmt.Println("   - More visual styling")

	generator = statements.New(
		statements.WithEnhancedDesign(),
	)
	statement, err = generator.Generate(input)
	if err != nil {
		log.Fatal("Failed to generate statement:", err)
	}

	pdfBytes, _ = statement.ToPDF()
	os.WriteFile("statement_enhanced.pdf", pdfBytes, 0644)
	fmt.Println("   ✅ Generated: statement_enhanced.pdf\n")

	// Test 3: Simple Template
	fmt.Println("3. GENERATING WITH SIMPLE TEMPLATE")
	fmt.Println("   - Basic text-based layout")
	fmt.Println("   - Minimal formatting")
	fmt.Println("   - Fastest generation")

	generator = statements.New(
		statements.WithSimpleDesign(),
	)
	statement, err = generator.Generate(input)
	if err != nil {
		log.Fatal("Failed to generate statement:", err)
	}

	pdfBytes, _ = statement.ToPDF()
	os.WriteFile("statement_simple.pdf", pdfBytes, 0644)
	fmt.Println("   ✅ Generated: statement_simple.pdf\n")

	// Test 4: Custom Configuration
	fmt.Println("4. GENERATING WITH CUSTOM CONFIGURATION")
	fmt.Println("   - Using minimalist with custom locale")

	generator = statements.New(
		statements.WithMinimalistDesign(),
		statements.WithLocale("fr-FR"),
	)
	statement, err = generator.Generate(input)
	if err != nil {
		log.Fatal("Failed to generate statement:", err)
	}

	pdfBytes, _ = statement.ToPDF()
	os.WriteFile("statement_custom.pdf", pdfBytes, 0644)
	fmt.Println("   ✅ Generated: statement_custom.pdf\n")

	// Summary
	fmt.Println("=== TEMPLATE TEST COMPLETE ===")
	fmt.Println("\nDefault template is: MINIMALIST (black & white)")
	fmt.Println("\nGenerated files:")
	fmt.Println("  • statement_minimalist.pdf - Default clean design")
	fmt.Println("  • statement_enhanced.pdf   - Colorful enhanced design")
	fmt.Println("  • statement_simple.pdf     - Basic simple design")
	fmt.Println("  • statement_custom.pdf     - Custom configuration")
	fmt.Println("\nTo change default, use WithEnhancedDesign() or WithSimpleDesign()")
	fmt.Println("Or use WithPDFRenderer('minimalist'|'enhanced'|'simple')")
}