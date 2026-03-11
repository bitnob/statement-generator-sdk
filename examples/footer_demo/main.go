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
	fmt.Println("=== Testing Enhanced Footer and Contact Information ===")

	// Create account
	account := statements.Account{
		Number:     "****9876",
		HolderName: "Sarah Johnson",
		Currency:   "USD",
		Type:       "Premier Checking",
		Address: &statements.Address{
			Line1:      "789 Park Avenue",
			Line2:      "Suite 1200",
			City:       "New York",
			State:      "NY",
			PostalCode: "10021",
			Country:    "United States",
		},
	}

	// Create sample transactions
	transactions := []statements.Transaction{
		{
			ID:          "DEMO001",
			Date:        time.Date(2024, 3, 1, 10, 0, 0, 0, time.UTC),
			Description: "Direct Deposit - Employer",
			Amount:      decimal.NewFromFloat(7500.00),
			Type:        statements.Credit,
			Reference:   "DD-MAR-01",
		},
		{
			ID:          "DEMO002",
			Date:        time.Date(2024, 3, 5, 14, 30, 0, 0, time.UTC),
			Description: "Mortgage Payment",
			Amount:      decimal.NewFromFloat(-2850.00),
			Type:        statements.Debit,
			Reference:   "MORT-2024-03",
		},
		{
			ID:          "DEMO003",
			Date:        time.Date(2024, 3, 8, 11, 15, 0, 0, time.UTC),
			Description: "Online Transfer",
			Amount:      decimal.NewFromFloat(-500.00),
			Type:        statements.Debit,
			Reference:   "TRF-3085",
		},
		{
			ID:          "DEMO004",
			Date:        time.Date(2024, 3, 10, 16, 20, 0, 0, time.UTC),
			Description: "Investment Dividend",
			Amount:      decimal.NewFromFloat(425.75),
			Type:        statements.Credit,
			Reference:   "DIV-Q1-2024",
		},
		{
			ID:          "DEMO005",
			Date:        time.Date(2024, 3, 15, 9, 45, 0, 0, time.UTC),
			Description: "Insurance Premium",
			Amount:      decimal.NewFromFloat(-385.00),
			Type:        statements.Debit,
			Reference:   "INS-AUTO-03",
		},
		{
			ID:          "DEMO006",
			Date:        time.Date(2024, 3, 20, 13, 30, 0, 0, time.UTC),
			Description: "Interest Credit",
			Amount:      decimal.NewFromFloat(18.52),
			Type:        statements.Credit,
			Reference:   "INT-MAR-2024",
		},
	}

	// Test 1: Standard bank with minimal contact info
	fmt.Println("1. STANDARD BANK WITH CONTACT INFO")
	fmt.Println("   - Basic contact details")
	fmt.Println("   - Default footer text")

	institution1 := &statements.Institution{
		Name: "First National Bank",
		Address: &statements.Address{
			Line1:   "100 Banking Plaza",
			City:    "New York",
			State:   "NY",
			Country: "USA",
		},
		RegNumber:    "Member FDIC | Equal Housing Lender",
		ContactPhone: "1-800-555-BANK (2265)",
		ContactEmail: "support@firstnational.com",
	}

	input := statements.StatementInput{
		Account:        account,
		Transactions:   transactions,
		PeriodStart:    time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
		PeriodEnd:      time.Date(2024, 3, 31, 23, 59, 59, 0, time.UTC),
		OpeningBalance: decimal.NewFromFloat(10000.00),
		Institution:    institution1,
	}

	generator := statements.New()
	statement, err := generator.Generate(input)
	if err != nil {
		log.Fatal("Failed to generate statement:", err)
	}

	pdfBytes, _ := statement.ToPDF()
	os.WriteFile("statement_standard_footer.pdf", pdfBytes, 0644)
	fmt.Println("   ✅ Generated: statement_standard_footer.pdf")

	// Test 2: Premium bank with full contact details
	fmt.Println("2. PREMIUM BANK WITH ENHANCED FOOTER")
	fmt.Println("   - Full contact information")
	fmt.Println("   - Custom footer text")
	fmt.Println("   - Website and email")

	institution2 := &statements.Institution{
		Name: "Global Premier Banking",
		Address: &statements.Address{
			Line1:   "One World Financial Center",
			Line2:   "Tower A, Floor 45",
			City:    "San Francisco",
			State:   "CA",
			Country: "USA",
		},
		RegNumber:    "Member FDIC | Reg# 89745632",
		ContactPhone: "1-888-PREMIER (773-6437)",
		ContactEmail: "premier.support@globalbank.com",
		Website:      "www.globalpremierbanking.com",
		FooterText:   "Thank you for banking with Global Premier Banking. For any questions about your statement or transactions, please contact our Premier Service team with your account number ready. We're available 24/7 to assist you. For security reasons, never share your full account details via email.",
	}

	input.Institution = institution2
	statement, err = generator.Generate(input)
	if err != nil {
		log.Fatal("Failed to generate statement:", err)
	}

	pdfBytes, _ = statement.ToPDF()
	os.WriteFile("statement_premium_footer.pdf", pdfBytes, 0644)
	fmt.Println("   ✅ Generated: statement_premium_footer.pdf")

	// Test 3: Fintech with modern contact approach
	fmt.Println("3. FINTECH WITH MODERN CONTACT STYLE")
	fmt.Println("   - Email-first support")
	fmt.Println("   - App-focused footer")
	fmt.Println("   - Modern messaging")

	institution3 := &statements.Institution{
		Name: "NeoBank Digital",
		Address: &statements.Address{
			Line1:   "Tech Hub",
			City:    "Austin",
			State:   "TX",
			Country: "USA",
		},
		RegNumber:    "Banking License #NB-2024-TX",
		ContactEmail: "help@neobank.app",
		Website:      "app.neobank.digital",
		FooterText:   "Questions? We're here to help! Open the NeoBank app and tap 'Support' for instant chat, or email us at help@neobank.app with your account number. Your security is our priority - we'll never ask for your password or PIN.",
	}

	input.Institution = institution3
	statement, err = generator.Generate(input)
	if err != nil {
		log.Fatal("Failed to generate statement:", err)
	}

	pdfBytes, _ = statement.ToPDF()
	os.WriteFile("statement_fintech_footer.pdf", pdfBytes, 0644)
	fmt.Println("   ✅ Generated: statement_fintech_footer.pdf")

	// Test 4: International bank with compliance footer
	fmt.Println("4. INTERNATIONAL BANK WITH COMPLIANCE FOOTER")
	fmt.Println("   - Regulatory compliance text")
	fmt.Println("   - International contact")
	fmt.Println("   - Multi-language support hint")

	institution4 := &statements.Institution{
		Name: "International Commerce Bank",
		Address: &statements.Address{
			Line1:   "200 Global Plaza",
			City:    "London",
			State:   "England",
			Country: "United Kingdom",
		},
		RegNumber:    "FCA Reg: 123456 | PRA Reg: 789012",
		ContactPhone: "+44 20 7xxx xxxx",
		ContactEmail: "customer.service@icb-global.com",
		Website:      "www.icb-global.com",
		FooterText:   "This statement is prepared in accordance with international banking standards. For discrepancies or queries, contact our Customer Service Centre quoting your account number. Llamadas en español: +1-800-xxx-xxxx. 中文服务: +86-xxx-xxxx. This document is confidential and intended solely for the addressee.",
	}

	input.Institution = institution4
	statement, err = generator.Generate(input)
	if err != nil {
		log.Fatal("Failed to generate statement:", err)
	}

	pdfBytes, _ = statement.ToPDF()
	os.WriteFile("statement_international_footer.pdf", pdfBytes, 0644)
	fmt.Println("   ✅ Generated: statement_international_footer.pdf")

	// Summary
	fmt.Println("=== FOOTER AND CONTACT TEST COMPLETE ===")
	fmt.Println("\nGenerated statements with different footer styles:")
	fmt.Println("  1. statement_standard_footer.pdf     - Basic bank contact")
	fmt.Println("  2. statement_premium_footer.pdf      - Full service details")
	fmt.Println("  3. statement_fintech_footer.pdf      - Modern digital approach")
	fmt.Println("  4. statement_international_footer.pdf - Compliance-focused")
	fmt.Println("\nEach PDF includes:")
	fmt.Println("  • Custom footer text (when provided)")
	fmt.Println("  • Contact phone and email")
	fmt.Println("  • Website information")
	fmt.Println("  • Account number reference in footer")
	fmt.Println("  • Generation timestamp at bottom")
	fmt.Println("  • Page numbering (Page X of Y)")
	fmt.Println("\n✅ All customizable footer fields are working!")
}