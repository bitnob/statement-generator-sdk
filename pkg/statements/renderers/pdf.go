package renderers

import (
	"bytes"
	"fmt"

	"github.com/bitnob/statement-generator-sdk/pkg/statements"
	"github.com/jung-kurt/gofpdf/v2"
)

// PDF implements the PDF renderer
type PDF struct {
	dateFormatter     *statements.DateFormatter
	currencyFormatter *statements.CurrencyFormatter
	addressFormatter  *statements.AddressFormatter
	pageSize          string
}

// NewPDF creates a new PDF renderer
func NewPDF(locale string) *PDF {
	return &PDF{
		dateFormatter:     statements.NewDateFormatter(locale),
		currencyFormatter: statements.NewCurrencyFormatter(locale),
		addressFormatter:  statements.NewAddressFormatter(),
		pageSize:          "A4",
	}
}

// SetPageSize sets the page size (A4, Letter, etc.)
func (p *PDF) SetPageSize(size string) {
	p.pageSize = size
}

// Render renders a statement as PDF
func (p *PDF) Render(statement *statements.Statement) (interface{}, error) {
	return p.RenderPDF(statement)
}

// RenderPDF renders a statement as PDF bytes
func (p *PDF) RenderPDF(statement *statements.Statement) ([]byte, error) {
	// Create new PDF document
	pdf := gofpdf.New("P", "mm", p.pageSize, "")

	// Set document properties
	pdf.SetTitle("Account Statement", false)
	pdf.SetAuthor(statement.Input.Account.HolderName, false)
	pdf.SetSubject("Account Statement", false)
	pdf.SetCreator("Statement Generator SDK", false)

	// Add first page
	pdf.AddPage()

	// Render header
	p.renderHeader(pdf, statement)

	// Render account information
	p.renderAccountInfo(pdf, statement)

	// Render summary
	p.renderSummary(pdf, statement)

	// Render transactions
	p.renderTransactions(pdf, statement)

	// Render footer on all pages
	p.addFooters(pdf, statement)

	// Output PDF
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// renderHeader renders the PDF header
func (p *PDF) renderHeader(pdf *gofpdf.Fpdf, statement *statements.Statement) {
	// Institution header if available
	if statement.Input.Institution != nil {
		// Institution name
		pdf.SetFont("Arial", "B", 18)
		pdf.CellFormat(0, 10, statement.Input.Institution.Name, "", 1, "C", false, 0, "")

		// Institution address if available
		if statement.Input.Institution.Address != nil {
			pdf.SetFont("Arial", "", 10)
			addressLines := p.addressFormatter.Format(statement.Input.Institution.Address)
			for _, line := range addressLines {
				pdf.CellFormat(0, 5, line, "", 1, "C", false, 0, "")
			}
		}

		// Registration number if available
		if statement.Input.Institution.RegNumber != "" {
			pdf.SetFont("Arial", "", 9)
			pdf.CellFormat(0, 5, "Reg No: "+statement.Input.Institution.RegNumber, "", 1, "C", false, 0, "")
		}

		pdf.Ln(5)
	}

	// Statement title
	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(0, 10, "ACCOUNT STATEMENT", "", 1, "C", false, 0, "")

	// Line separator
	pdf.SetDrawColor(44, 62, 80)
	pdf.SetLineWidth(0.5)
	pdf.Line(10, pdf.GetY(), 200, pdf.GetY())
	pdf.Ln(5)
}

// renderAccountInfo renders account information section
func (p *PDF) renderAccountInfo(pdf *gofpdf.Fpdf, statement *statements.Statement) {
	// Section title
	pdf.SetFont("Arial", "B", 12)
	pdf.SetTextColor(44, 62, 80)
	pdf.CellFormat(0, 8, "Account Information", "", 1, "L", false, 0, "")
	pdf.SetTextColor(0, 0, 0)

	pdf.SetFont("Arial", "", 10)

	// Account holder
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(40, 6, "Account Holder:", "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(0, 6, statement.Input.Account.HolderName, "", 1, "L", false, 0, "")

	// Account number
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(40, 6, "Account Number:", "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(0, 6, statement.Input.Account.Number, "", 1, "L", false, 0, "")

	// Account type if available
	if statement.Input.Account.Type != "" {
		pdf.SetFont("Arial", "B", 10)
		pdf.CellFormat(40, 6, "Account Type:", "", 0, "L", false, 0, "")
		pdf.SetFont("Arial", "", 10)
		pdf.CellFormat(0, 6, statement.Input.Account.Type, "", 1, "L", false, 0, "")
	}

	// Statement period
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(40, 6, "Statement Period:", "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "", 10)
	period := p.dateFormatter.FormatPeriod(statement.Input.PeriodStart, statement.Input.PeriodEnd)
	pdf.CellFormat(0, 6, period, "", 1, "L", false, 0, "")

	// Currency
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(40, 6, "Currency:", "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(0, 6, statement.Input.Account.Currency, "", 1, "L", false, 0, "")

	// Address if available (for proof of address)
	if statement.Input.Account.Address != nil {
		pdf.SetFont("Arial", "B", 10)
		pdf.CellFormat(40, 6, "Address:", "", 0, "L", false, 0, "")
		pdf.SetFont("Arial", "", 10)

		addressLines := p.addressFormatter.Format(statement.Input.Account.Address)
		for i, line := range addressLines {
			if i == 0 {
				pdf.CellFormat(0, 6, line, "", 1, "L", false, 0, "")
			} else {
				pdf.CellFormat(40, 6, "", "", 0, "L", false, 0, "") // Empty space for alignment
				pdf.CellFormat(0, 6, line, "", 1, "L", false, 0, "")
			}
		}
	}

	pdf.Ln(5)
}

// renderSummary renders the summary section
func (p *PDF) renderSummary(pdf *gofpdf.Fpdf, statement *statements.Statement) {
	// Background color for summary
	pdf.SetFillColor(240, 240, 240)
	pdf.Rect(10, pdf.GetY(), 190, 35, "F")

	// Section title
	pdf.SetFont("Arial", "B", 12)
	pdf.SetTextColor(44, 62, 80)
	pdf.CellFormat(0, 8, "Account Summary", "", 1, "L", false, 0, "")
	pdf.SetTextColor(0, 0, 0)

	// Create two columns

	// Left column
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(50, 6, "Opening Balance:", "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "", 10)
	openingBalance := p.currencyFormatter.FormatAmount(statement.Input.OpeningBalance, statement.Input.Account.Currency)
	pdf.CellFormat(45, 6, openingBalance, "", 0, "R", false, 0, "")

	// Right column
	pdf.SetX(105)
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(50, 6, "Total Credits:", "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "", 10)
	totalCredits := p.currencyFormatter.FormatAmount(statement.TotalCredits, statement.Input.Account.Currency)
	pdf.CellFormat(45, 6, totalCredits, "", 1, "R", false, 0, "")

	// Second row
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(50, 6, "Total Debits:", "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "", 10)
	totalDebits := p.currencyFormatter.FormatAmount(statement.TotalDebits, statement.Input.Account.Currency)
	pdf.CellFormat(45, 6, totalDebits, "", 0, "R", false, 0, "")

	// Right column
	pdf.SetX(105)
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(50, 6, "Transaction Count:", "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(45, 6, fmt.Sprintf("%d", statement.TransactionCount), "", 1, "R", false, 0, "")

	// Closing balance (prominent)
	pdf.Ln(2)
	pdf.SetFont("Arial", "B", 12)
	pdf.SetTextColor(44, 62, 80)
	pdf.CellFormat(50, 8, "CLOSING BALANCE:", "", 0, "L", false, 0, "")
	closingBalance := p.currencyFormatter.FormatAmount(statement.ClosingBalance, statement.Input.Account.Currency)
	pdf.CellFormat(45, 8, closingBalance, "", 1, "R", false, 0, "")
	pdf.SetTextColor(0, 0, 0)

	pdf.Ln(8)
}

// renderTransactions renders the transaction table
func (p *PDF) renderTransactions(pdf *gofpdf.Fpdf, statement *statements.Statement) {
	// Section title
	pdf.SetFont("Arial", "B", 12)
	pdf.SetTextColor(44, 62, 80)
	pdf.CellFormat(0, 8, "Transaction History", "", 1, "L", false, 0, "")
	pdf.SetTextColor(0, 0, 0)

	// Table header
	p.renderTableHeader(pdf)

	// Opening balance row
	pdf.SetFont("Arial", "", 9)
	pdf.CellFormat(20, 6, p.dateFormatter.Format(statement.Input.PeriodStart), "B", 0, "L", false, 0, "")
	pdf.CellFormat(65, 6, "Opening Balance", "B", 0, "L", false, 0, "")
	pdf.CellFormat(25, 6, "", "B", 0, "R", false, 0, "")
	pdf.CellFormat(25, 6, "", "B", 0, "R", false, 0, "")
	openingBalance := p.currencyFormatter.FormatAmount(statement.Input.OpeningBalance, statement.Input.Account.Currency)
	pdf.CellFormat(30, 6, openingBalance, "B", 0, "R", false, 0, "")
	pdf.CellFormat(25, 6, "", "B", 1, "L", false, 0, "")

	// Transactions
	calculator := statements.NewCalculator()
	sortedTransactions := calculator.SortTransactions(statement.Input.Transactions)
	runningBalance := statement.Input.OpeningBalance

	for i, txn := range sortedTransactions {
		// Check if we need a new page
		if pdf.GetY() > 250 {
			pdf.AddPage()
			p.renderTableHeader(pdf)
		}

		runningBalance = runningBalance.Add(txn.Amount)

		// Alternate row colors
		if i%2 == 0 {
			pdf.SetFillColor(250, 250, 250)
		} else {
			pdf.SetFillColor(255, 255, 255)
		}

		pdf.SetFont("Arial", "", 9)

		// Date
		pdf.CellFormat(20, 6, p.dateFormatter.Format(txn.Date), "B", 0, "L", true, 0, "")

		// Description (truncate if too long)
		description := txn.Description
		if len(description) > 40 {
			description = description[:37] + "..."
		}
		pdf.CellFormat(65, 6, description, "B", 0, "L", true, 0, "")

		// Debit
		debitAmount := ""
		if txn.Type == statements.Debit {
			debitAmount = p.currencyFormatter.FormatAmount(txn.Amount.Abs(), statement.Input.Account.Currency)
			pdf.SetTextColor(220, 53, 69) // Red for debits
		}
		pdf.CellFormat(25, 6, debitAmount, "B", 0, "R", true, 0, "")
		pdf.SetTextColor(0, 0, 0)

		// Credit
		creditAmount := ""
		if txn.Type == statements.Credit {
			creditAmount = p.currencyFormatter.FormatAmount(txn.Amount, statement.Input.Account.Currency)
			pdf.SetTextColor(40, 167, 69) // Green for credits
		}
		pdf.CellFormat(25, 6, creditAmount, "B", 0, "R", true, 0, "")
		pdf.SetTextColor(0, 0, 0)

		// Balance
		balance := p.currencyFormatter.FormatAmount(runningBalance, statement.Input.Account.Currency)
		pdf.SetFont("Arial", "B", 9)
		pdf.CellFormat(30, 6, balance, "B", 0, "R", true, 0, "")

		// Reference
		pdf.SetFont("Arial", "", 8)
		reference := txn.Reference
		if len(reference) > 15 {
			reference = reference[:12] + "..."
		}
		pdf.CellFormat(25, 6, reference, "B", 1, "L", true, 0, "")
	}
}

// renderTableHeader renders the transaction table header
func (p *PDF) renderTableHeader(pdf *gofpdf.Fpdf) {
	pdf.SetFont("Arial", "B", 9)
	pdf.SetFillColor(44, 62, 80)
	pdf.SetTextColor(255, 255, 255)

	pdf.CellFormat(20, 7, "Date", "TB", 0, "L", true, 0, "")
	pdf.CellFormat(65, 7, "Description", "TB", 0, "L", true, 0, "")
	pdf.CellFormat(25, 7, "Debit", "TB", 0, "R", true, 0, "")
	pdf.CellFormat(25, 7, "Credit", "TB", 0, "R", true, 0, "")
	pdf.CellFormat(30, 7, "Balance", "TB", 0, "R", true, 0, "")
	pdf.CellFormat(25, 7, "Reference", "TB", 1, "L", true, 0, "")

	pdf.SetTextColor(0, 0, 0)
	pdf.SetFillColor(255, 255, 255)
}

// addFooters adds footers to all pages
func (p *PDF) addFooters(pdf *gofpdf.Fpdf, statement *statements.Statement) {
	pageCount := pdf.PageCount()

	// Add footer to each page
	for i := 1; i <= pageCount; i++ {
		pdf.SetPage(i)

		// Position at 15mm from bottom
		pdf.SetY(-15)

		// Footer text
		pdf.SetFont("Arial", "I", 8)
		pdf.SetTextColor(128, 128, 128)

		// Page number
		pageText := fmt.Sprintf("Page %d of %d", i, pageCount)
		pdf.CellFormat(0, 5, pageText, "", 0, "C", false, 0, "")

		// Generated timestamp
		pdf.SetY(-10)
		generatedText := fmt.Sprintf("Generated: %s", p.dateFormatter.FormatDateTime(statement.GeneratedAt))
		pdf.CellFormat(0, 5, generatedText, "", 0, "C", false, 0, "")
	}
}

// RenderCompact creates a more compact PDF suitable for large transaction volumes
func (p *PDF) RenderCompact(statement *statements.Statement) ([]byte, error) {
	pdf := gofpdf.New("L", "mm", p.pageSize, "") // Landscape orientation
	pdf.SetFont("Arial", "", 8) // Smaller font

	pdf.AddPage()

	// Compact header
	pdf.SetFont("Arial", "B", 12)
	title := fmt.Sprintf("Statement: %s - %s",
		statement.Input.Account.HolderName,
		statement.Input.Account.Number)
	pdf.CellFormat(0, 8, title, "", 1, "L", false, 0, "")

	period := p.dateFormatter.FormatPeriod(statement.Input.PeriodStart, statement.Input.PeriodEnd)
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(0, 6, period, "", 1, "L", false, 0, "")

	// Compact summary in one line
	pdf.SetFont("Arial", "", 9)
	summaryText := fmt.Sprintf(
		"Opening: %s | Credits: %s | Debits: %s | Closing: %s | Count: %d",
		p.currencyFormatter.FormatAmount(statement.Input.OpeningBalance, statement.Input.Account.Currency),
		p.currencyFormatter.FormatAmount(statement.TotalCredits, statement.Input.Account.Currency),
		p.currencyFormatter.FormatAmount(statement.TotalDebits, statement.Input.Account.Currency),
		p.currencyFormatter.FormatAmount(statement.ClosingBalance, statement.Input.Account.Currency),
		statement.TransactionCount,
	)
	pdf.CellFormat(0, 6, summaryText, "", 1, "L", false, 0, "")
	pdf.Ln(3)

	// Compact transaction table
	p.renderCompactTable(pdf, statement)

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// renderCompactTable renders a compact transaction table
func (p *PDF) renderCompactTable(pdf *gofpdf.Fpdf, statement *statements.Statement) {
	// Table header
	pdf.SetFont("Arial", "B", 8)
	pdf.SetFillColor(200, 200, 200)

	pdf.CellFormat(18, 5, "Date", "1", 0, "L", true, 0, "")
	pdf.CellFormat(80, 5, "Description", "1", 0, "L", true, 0, "")
	pdf.CellFormat(22, 5, "Debit", "1", 0, "R", true, 0, "")
	pdf.CellFormat(22, 5, "Credit", "1", 0, "R", true, 0, "")
	pdf.CellFormat(25, 5, "Balance", "1", 0, "R", true, 0, "")
	pdf.CellFormat(30, 5, "Reference", "1", 1, "L", true, 0, "")

	// Transactions
	pdf.SetFont("Arial", "", 7)
	calculator := statements.NewCalculator()
	sortedTransactions := calculator.SortTransactions(statement.Input.Transactions)
	runningBalance := statement.Input.OpeningBalance

	for _, txn := range sortedTransactions {
		if pdf.GetY() > 180 {
			pdf.AddPage()
			// Repeat header
			pdf.SetFont("Arial", "B", 8)
			pdf.SetFillColor(200, 200, 200)
			pdf.CellFormat(18, 5, "Date", "1", 0, "L", true, 0, "")
			pdf.CellFormat(80, 5, "Description", "1", 0, "L", true, 0, "")
			pdf.CellFormat(22, 5, "Debit", "1", 0, "R", true, 0, "")
			pdf.CellFormat(22, 5, "Credit", "1", 0, "R", true, 0, "")
			pdf.CellFormat(25, 5, "Balance", "1", 0, "R", true, 0, "")
			pdf.CellFormat(30, 5, "Reference", "1", 1, "L", true, 0, "")
			pdf.SetFont("Arial", "", 7)
		}

		runningBalance = runningBalance.Add(txn.Amount)

		pdf.CellFormat(18, 4, p.dateFormatter.Format(txn.Date), "LR", 0, "L", false, 0, "")

		description := txn.Description
		if len(description) > 50 {
			description = description[:47] + "..."
		}
		pdf.CellFormat(80, 4, description, "LR", 0, "L", false, 0, "")

		if txn.Type == statements.Debit {
			pdf.CellFormat(22, 4, txn.Amount.Abs().StringFixed(2), "LR", 0, "R", false, 0, "")
			pdf.CellFormat(22, 4, "", "LR", 0, "R", false, 0, "")
		} else {
			pdf.CellFormat(22, 4, "", "LR", 0, "R", false, 0, "")
			pdf.CellFormat(22, 4, txn.Amount.StringFixed(2), "LR", 0, "R", false, 0, "")
		}

		pdf.CellFormat(25, 4, runningBalance.StringFixed(2), "LR", 0, "R", false, 0, "")
		pdf.CellFormat(30, 4, txn.Reference, "LR", 1, "L", false, 0, "")
	}

	// Bottom border
	pdf.CellFormat(197, 0, "", "T", 0, "L", false, 0, "")
}