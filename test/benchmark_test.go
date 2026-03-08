package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/bitnob/statement-generator-sdk/pkg/statements"
	"github.com/shopspring/decimal"
)

func generateTransactions(count int) []statements.Transaction {
	transactions := make([]statements.Transaction, count)
	baseDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	for i := 0; i < count; i++ {
		txType := statements.Credit
		amount := decimal.NewFromFloat(100.50)
		if i%3 == 0 {
			txType = statements.Debit
			amount = amount.Neg()
		}

		transactions[i] = statements.Transaction{
			ID:          fmt.Sprintf("TXN%06d", i+1),
			Date:        baseDate.AddDate(0, 0, i),
			Description: fmt.Sprintf("Transaction %d", i+1),
			Amount:      amount,
			Type:        txType,
			Reference:   fmt.Sprintf("REF%06d", i+1),
		}
	}

	return transactions
}

func generateStatementInput(transactionCount int) statements.StatementInput {
	return statements.StatementInput{
		Account: statements.Account{
			Number:     "****1234",
			HolderName: "Benchmark User",
			Currency:   "USD",
			Address: &statements.Address{
				Line1:      "123 Benchmark Street",
				City:       "Test City",
				State:      "TC",
				PostalCode: "12345",
				Country:    "USA",
			},
		},
		Transactions:   generateTransactions(transactionCount),
		PeriodStart:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		PeriodEnd:      time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC),
		OpeningBalance: decimal.NewFromFloat(10000),
		Institution: &statements.Institution{
			Name: "Benchmark Bank",
			Address: &statements.Address{
				Line1:   "100 Financial Ave",
				City:    "Banking City",
				Country: "USA",
			},
		},
	}
}

// Benchmark generation with different transaction counts
func BenchmarkStatementGeneration_10(b *testing.B) {
	benchmarkStatementGeneration(b, 10)
}

func BenchmarkStatementGeneration_100(b *testing.B) {
	benchmarkStatementGeneration(b, 100)
}

func BenchmarkStatementGeneration_1000(b *testing.B) {
	benchmarkStatementGeneration(b, 1000)
}

func BenchmarkStatementGeneration_10000(b *testing.B) {
	benchmarkStatementGeneration(b, 10000)
}

func benchmarkStatementGeneration(b *testing.B, transactionCount int) {
	generator := statements.New()
	input := generateStatementInput(transactionCount)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := generator.Generate(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark PDF generation
func BenchmarkPDFGeneration_100(b *testing.B) {
	benchmarkPDFGeneration(b, 100)
}

func BenchmarkPDFGeneration_1000(b *testing.B) {
	benchmarkPDFGeneration(b, 1000)
}

func BenchmarkPDFGeneration_10000(b *testing.B) {
	benchmarkPDFGeneration(b, 10000)
}

func benchmarkPDFGeneration(b *testing.B, transactionCount int) {
	generator := statements.New()
	input := generateStatementInput(transactionCount)
	statement, err := generator.Generate(input)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := statement.ToPDF()
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark CSV generation
func BenchmarkCSVGeneration_100(b *testing.B) {
	benchmarkCSVGeneration(b, 100)
}

func BenchmarkCSVGeneration_1000(b *testing.B) {
	benchmarkCSVGeneration(b, 1000)
}

func BenchmarkCSVGeneration_10000(b *testing.B) {
	benchmarkCSVGeneration(b, 10000)
}

func benchmarkCSVGeneration(b *testing.B, transactionCount int) {
	generator := statements.New()
	input := generateStatementInput(transactionCount)
	statement, err := generator.Generate(input)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = statement.ToCSV()
	}
}

// Benchmark HTML generation
func BenchmarkHTMLGeneration_100(b *testing.B) {
	benchmarkHTMLGeneration(b, 100)
}

func BenchmarkHTMLGeneration_1000(b *testing.B) {
	benchmarkHTMLGeneration(b, 1000)
}

func BenchmarkHTMLGeneration_10000(b *testing.B) {
	benchmarkHTMLGeneration(b, 10000)
}

func benchmarkHTMLGeneration(b *testing.B, transactionCount int) {
	generator := statements.New()
	input := generateStatementInput(transactionCount)
	statement, err := generator.Generate(input)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = statement.ToHTML()
	}
}

// Benchmark balance calculation
func BenchmarkBalanceCalculation_100(b *testing.B) {
	benchmarkBalanceCalculation(b, 100)
}

func BenchmarkBalanceCalculation_1000(b *testing.B) {
	benchmarkBalanceCalculation(b, 1000)
}

func BenchmarkBalanceCalculation_10000(b *testing.B) {
	benchmarkBalanceCalculation(b, 10000)
}

func benchmarkBalanceCalculation(b *testing.B, transactionCount int) {
	calculator := statements.NewCalculator()
	transactions := generateTransactions(transactionCount)
	openingBalance := decimal.NewFromFloat(10000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := calculator.CalculateBalances(transactions, openingBalance)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark validation
func BenchmarkValidation_100(b *testing.B) {
	benchmarkValidation(b, 100)
}

func BenchmarkValidation_1000(b *testing.B) {
	benchmarkValidation(b, 1000)
}

func BenchmarkValidation_10000(b *testing.B) {
	benchmarkValidation(b, 10000)
}

func benchmarkValidation(b *testing.B, transactionCount int) {
	validator := statements.NewValidator()
	input := generateStatementInput(transactionCount)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := validator.ValidateStatementInput(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark currency formatting
func BenchmarkCurrencyFormatting(b *testing.B) {
	formatter := statements.NewCurrencyFormatter("en-US")
	amount := decimal.NewFromFloat(1234567.89)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = formatter.FormatAmount(amount, "USD")
	}
}

// Benchmark date formatting
func BenchmarkDateFormatting(b *testing.B) {
	formatter := statements.NewDateFormatter("en-US")
	date := time.Date(2024, 1, 15, 14, 30, 45, 0, time.UTC)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = formatter.Format(date)
	}
}

// Memory allocation benchmarks
func BenchmarkMemoryAllocation_1000Transactions(b *testing.B) {
	generator := statements.New()
	input := generateStatementInput(1000)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		statement, err := generator.Generate(input)
		if err != nil {
			b.Fatal(err)
		}
		_, _ = statement.ToPDF()
	}
}

func BenchmarkMemoryAllocation_10000Transactions(b *testing.B) {
	generator := statements.New()
	input := generateStatementInput(10000)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		statement, err := generator.Generate(input)
		if err != nil {
			b.Fatal(err)
		}
		_, _ = statement.ToPDF()
	}
}

// Concurrent generation benchmark
func BenchmarkConcurrentGeneration(b *testing.B) {
	generator := statements.New()
	input := generateStatementInput(100)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			statement, err := generator.Generate(input)
			if err != nil {
				b.Fatal(err)
			}
			_, _ = statement.ToPDF()
		}
	})
}

// Quick methods benchmarks
func BenchmarkQuickPDF(b *testing.B) {
	transactions := generateTransactions(100)
	openingBalance := decimal.NewFromFloat(10000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := statements.QuickPDF(transactions, openingBalance)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkQuickCSV(b *testing.B) {
	transactions := generateTransactions(100)
	openingBalance := decimal.NewFromFloat(10000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := statements.QuickCSV(transactions, openingBalance)
		if err != nil {
			b.Fatal(err)
		}
	}
}