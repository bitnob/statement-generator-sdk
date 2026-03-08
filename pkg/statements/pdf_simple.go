package statements

import (
	"bytes"
	"fmt"

	"github.com/jung-kurt/gofpdf/v2"
)

// SimplePDFRender creates a simple PDF rendering implementation
func (s *Statement) SimplePDFRender() ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Header
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Account Statement")
	pdf.Ln(12)

	// Account info
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, fmt.Sprintf("Account Holder: %s", s.Input.Account.HolderName))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Account Number: %s", s.Input.Account.Number))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Currency: %s", s.Input.Account.Currency))
	pdf.Ln(10)

	// Summary
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(40, 10, "Summary")
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, fmt.Sprintf("Opening Balance: %s", s.Input.OpeningBalance.String()))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Total Credits: %s", s.TotalCredits.String()))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Total Debits: %s", s.TotalDebits.String()))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Closing Balance: %s", s.ClosingBalance.String()))
	pdf.Ln(10)

	// Transactions header
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(40, 10, "Transactions")
	pdf.Ln(8)

	// Simple transaction list
	pdf.SetFont("Arial", "", 10)
	calculator := NewCalculator()
	sortedTransactions := calculator.SortTransactions(s.Input.Transactions)
	runningBalance := s.Input.OpeningBalance

	for _, txn := range sortedTransactions {
		runningBalance = runningBalance.Add(txn.Amount)

		txnLine := fmt.Sprintf("%s | %s | %s | Balance: %s",
			txn.Date.Format("2006-01-02"),
			txn.Description,
			txn.Amount.String(),
			runningBalance.String())

		pdf.Cell(40, 10, txnLine)
		pdf.Ln(6)

		// Check for page break
		if pdf.GetY() > 250 {
			pdf.AddPage()
		}
	}

	// Output
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}