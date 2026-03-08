package statements

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/jung-kurt/gofpdf/v2"
	"github.com/shopspring/decimal"
)

// MinimalistPDFRender creates a clean, minimalist black and white PDF statement
func (s *Statement) MinimalistPDFRender() ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")

	// Set margins for clean spacing
	pdf.SetMargins(20, 20, 20)
	pdf.SetAutoPageBreak(true, 25)

	pdf.AddPage()

	// Color palette - minimalist black and white
	black := []int{0, 0, 0}
	white := []int{255, 255, 255}
	lightGrey := []int{248, 248, 248}    // Off-white for alternating rows
	darkGrey := []int{100, 100, 100}     // For secondary text
	borderGrey := []int{220, 220, 220}   // Light borders

	// Institution Header - Clean and minimal
	if s.Input.Institution != nil {
		// Institution name - large, bold
		pdf.SetFont("Helvetica", "B", 20)
		pdf.SetTextColor(black[0], black[1], black[2])
		pdf.CellFormat(0, 10, strings.ToUpper(s.Input.Institution.Name), "", 1, "L", false, 0, "")

		// Thin line separator
		pdf.SetDrawColor(borderGrey[0], borderGrey[1], borderGrey[2])
		pdf.SetLineWidth(0.5)
		y := pdf.GetY()
		pdf.Line(20, y, 190, y)
		pdf.Ln(3)

		// Institution details - small, grey
		if s.Input.Institution.Address != nil {
			pdf.SetFont("Helvetica", "", 9)
			pdf.SetTextColor(darkGrey[0], darkGrey[1], darkGrey[2])
			addr := s.Input.Institution.Address

			// Single line address
			addressLine := fmt.Sprintf("%s, %s, %s %s",
				addr.Line1, addr.City, addr.State, addr.PostalCode)
			pdf.CellFormat(0, 4, addressLine, "", 1, "L", false, 0, "")

			if s.Input.Institution.RegNumber != "" {
				pdf.CellFormat(0, 4, s.Input.Institution.RegNumber, "", 1, "L", false, 0, "")
			}
		}
		pdf.Ln(8)
	}

	// Statement Title - Simple and clean
	pdf.SetFont("Helvetica", "", 14)
	pdf.SetTextColor(black[0], black[1], black[2])
	pdf.CellFormat(0, 8, "Statement of Account", "", 1, "L", false, 0, "")
	pdf.Ln(4)

	// Two-column layout for account info and period
	pdf.SetFont("Helvetica", "", 10)

	// Left column - Account holder info
	pdf.SetTextColor(darkGrey[0], darkGrey[1], darkGrey[2])
	pdf.CellFormat(30, 5, "Account Holder", "", 0, "L", false, 0, "")
	pdf.SetTextColor(black[0], black[1], black[2])
	pdf.SetFont("Helvetica", "B", 10)
	pdf.CellFormat(65, 5, s.Input.Account.HolderName, "", 0, "L", false, 0, "")

	// Right column - Period
	pdf.SetFont("Helvetica", "", 10)
	pdf.SetTextColor(darkGrey[0], darkGrey[1], darkGrey[2])
	pdf.CellFormat(30, 5, "Statement Period", "", 0, "L", false, 0, "")
	pdf.SetTextColor(black[0], black[1], black[2])
	pdf.CellFormat(0, 5, fmt.Sprintf("%s - %s",
		s.Input.PeriodStart.Format("Jan 02, 2006"),
		s.Input.PeriodEnd.Format("Jan 02, 2006")), "", 1, "L", false, 0, "")

	// Account number and currency
	pdf.SetFont("Helvetica", "", 10)
	pdf.SetTextColor(darkGrey[0], darkGrey[1], darkGrey[2])
	pdf.CellFormat(30, 5, "Account Number", "", 0, "L", false, 0, "")
	pdf.SetTextColor(black[0], black[1], black[2])
	pdf.CellFormat(65, 5, s.Input.Account.Number, "", 0, "L", false, 0, "")

	pdf.SetTextColor(darkGrey[0], darkGrey[1], darkGrey[2])
	pdf.CellFormat(30, 5, "Currency", "", 0, "L", false, 0, "")
	pdf.SetTextColor(black[0], black[1], black[2])
	pdf.CellFormat(0, 5, s.Input.Account.Currency, "", 1, "L", false, 0, "")

	// Account type
	if s.Input.Account.Type != "" {
		pdf.SetTextColor(darkGrey[0], darkGrey[1], darkGrey[2])
		pdf.CellFormat(30, 5, "Account Type", "", 0, "L", false, 0, "")
		pdf.SetTextColor(black[0], black[1], black[2])
		pdf.CellFormat(65, 5, s.Input.Account.Type, "", 1, "L", false, 0, "")
	}

	// Address Section (if present) - Minimalist box
	if s.Input.Account.Address != nil {
		pdf.Ln(4)
		pdf.SetFont("Helvetica", "", 10)
		pdf.SetTextColor(darkGrey[0], darkGrey[1], darkGrey[2])
		pdf.CellFormat(30, 5, "Mailing Address", "", 0, "L", false, 0, "")

		pdf.SetTextColor(black[0], black[1], black[2])
		addr := s.Input.Account.Address
		addressText := addr.Line1
		if addr.Line2 != "" {
			addressText += ", " + addr.Line2
		}
		pdf.CellFormat(0, 5, addressText, "", 1, "L", false, 0, "")

		pdf.CellFormat(30, 5, "", "", 0, "L", false, 0, "")
		pdf.CellFormat(0, 5, fmt.Sprintf("%s, %s %s, %s",
			addr.City, addr.State, addr.PostalCode, addr.Country), "", 1, "L", false, 0, "")
	}

	pdf.Ln(8)

	// Balance Summary - Clean grid layout
	pdf.SetDrawColor(borderGrey[0], borderGrey[1], borderGrey[2])
	pdf.SetLineWidth(0.3)

	// Summary header
	pdf.SetFont("Helvetica", "", 11)
	pdf.SetTextColor(black[0], black[1], black[2])
	pdf.CellFormat(0, 8, "Account Summary", "", 1, "L", false, 0, "")

	// Draw a subtle line
	y := pdf.GetY()
	pdf.Line(20, y, 190, y)
	pdf.Ln(2)

	formatter := NewCurrencyFormatter("en-US")

	// Summary grid - clean 2x2 layout
	pdf.SetFont("Helvetica", "", 10)

	// Opening Balance
	pdf.SetTextColor(darkGrey[0], darkGrey[1], darkGrey[2])
	pdf.CellFormat(45, 6, "Opening Balance", "", 0, "L", false, 0, "")
	pdf.SetTextColor(black[0], black[1], black[2])
	pdf.CellFormat(50, 6, formatter.FormatAmount(s.Input.OpeningBalance, s.Input.Account.Currency), "", 0, "R", false, 0, "")

	// Total Credits
	pdf.SetTextColor(darkGrey[0], darkGrey[1], darkGrey[2])
	pdf.CellFormat(45, 6, "Total Credits", "", 0, "L", false, 0, "")
	pdf.SetTextColor(black[0], black[1], black[2])
	pdf.CellFormat(0, 6, formatter.FormatAmount(s.TotalCredits, s.Input.Account.Currency), "", 1, "R", false, 0, "")

	// Total Debits
	pdf.SetTextColor(darkGrey[0], darkGrey[1], darkGrey[2])
	pdf.CellFormat(45, 6, "Total Debits", "", 0, "L", false, 0, "")
	pdf.SetTextColor(black[0], black[1], black[2])
	pdf.CellFormat(50, 6, formatter.FormatAmount(s.TotalDebits, s.Input.Account.Currency), "", 0, "R", false, 0, "")

	// Closing Balance - slightly emphasized
	pdf.SetTextColor(darkGrey[0], darkGrey[1], darkGrey[2])
	pdf.CellFormat(45, 6, "Closing Balance", "", 0, "L", false, 0, "")
	pdf.SetFont("Helvetica", "B", 10)
	pdf.SetTextColor(black[0], black[1], black[2])
	pdf.CellFormat(0, 6, formatter.FormatAmount(s.ClosingBalance, s.Input.Account.Currency), "", 1, "R", false, 0, "")
	pdf.SetFont("Helvetica", "", 10)

	// Draw bottom line
	y = pdf.GetY()
	pdf.Line(20, y, 190, y)
	pdf.Ln(10)

	// Transaction Table - Minimalist design
	pdf.SetFont("Helvetica", "", 11)
	pdf.SetTextColor(black[0], black[1], black[2])
	pdf.CellFormat(0, 8, "Transaction Details", "", 1, "L", false, 0, "")

	// Draw line under title
	y = pdf.GetY()
	pdf.Line(20, y, 190, y)
	pdf.Ln(2)

	// Table Headers - Simple, no background
	pdf.SetFont("Helvetica", "", 9)
	pdf.SetTextColor(darkGrey[0], darkGrey[1], darkGrey[2])

	// Define column widths for better proportion
	dateWidth := 22.0
	descWidth := 70.0
	refWidth := 35.0
	debitWidth := 28.0
	creditWidth := 28.0
	balanceWidth := 30.0

	// Headers with bottom border only
	pdf.CellFormat(dateWidth, 6, "Date", "B", 0, "L", false, 0, "")
	pdf.CellFormat(descWidth, 6, "Description", "B", 0, "L", false, 0, "")
	pdf.CellFormat(refWidth, 6, "Reference", "B", 0, "L", false, 0, "")
	pdf.CellFormat(debitWidth, 6, "Debit", "B", 0, "R", false, 0, "")
	pdf.CellFormat(creditWidth, 6, "Credit", "B", 0, "R", false, 0, "")
	pdf.CellFormat(balanceWidth, 6, "Balance", "B", 1, "R", false, 0, "")

	// Transaction rows
	pdf.SetFont("Helvetica", "", 8)
	calculator := NewCalculator()
	sortedTransactions := calculator.SortTransactions(s.Input.Transactions)
	runningBalance := s.Input.OpeningBalance

	// Use alternating row colors for better readability
	useAlternateColor := false

	for i, txn := range sortedTransactions {
		// Check for page break
		if pdf.GetY() > 255 {
			// Add page number before breaking
			pdf.SetY(-15)
			pdf.SetFont("Helvetica", "", 8)
			pdf.SetTextColor(darkGrey[0], darkGrey[1], darkGrey[2])
			pageNo := pdf.PageNo()
			pdf.CellFormat(0, 10, fmt.Sprintf("Page %d", pageNo), "", 0, "C", false, 0, "")

			pdf.AddPage()

			// Repeat minimal header on new page
			pdf.SetFont("Helvetica", "", 11)
			pdf.SetTextColor(black[0], black[1], black[2])
			pdf.CellFormat(0, 8, "Transaction Details (continued)", "", 1, "L", false, 0, "")
			y := pdf.GetY()
			pdf.Line(20, y, 190, y)
			pdf.Ln(2)

			// Repeat column headers
			pdf.SetFont("Helvetica", "", 9)
			pdf.SetTextColor(darkGrey[0], darkGrey[1], darkGrey[2])
			pdf.CellFormat(dateWidth, 6, "Date", "B", 0, "L", false, 0, "")
			pdf.CellFormat(descWidth, 6, "Description", "B", 0, "L", false, 0, "")
			pdf.CellFormat(refWidth, 6, "Reference", "B", 0, "L", false, 0, "")
			pdf.CellFormat(debitWidth, 6, "Debit", "B", 0, "R", false, 0, "")
			pdf.CellFormat(creditWidth, 6, "Credit", "B", 0, "R", false, 0, "")
			pdf.CellFormat(balanceWidth, 6, "Balance", "B", 1, "R", false, 0, "")
			pdf.SetFont("Helvetica", "", 8)
		}

		runningBalance = runningBalance.Add(txn.Amount)

		// Alternate row background
		if useAlternateColor {
			pdf.SetFillColor(lightGrey[0], lightGrey[1], lightGrey[2])
		} else {
			pdf.SetFillColor(white[0], white[1], white[2])
		}
		useAlternateColor = !useAlternateColor

		pdf.SetTextColor(black[0], black[1], black[2])

		// Date
		pdf.CellFormat(dateWidth, 5, txn.Date.Format("Jan 02"), "", 0, "L", true, 0, "")

		// Description (truncate if needed)
		desc := txn.Description
		if len(desc) > 40 {
			desc = desc[:37] + "..."
		}
		pdf.CellFormat(descWidth, 5, desc, "", 0, "L", true, 0, "")

		// Reference (truncate if needed)
		ref := txn.Reference
		if len(ref) > 20 {
			ref = ref[:17] + "..."
		}
		pdf.CellFormat(refWidth, 5, ref, "", 0, "L", true, 0, "")

		// Debit/Credit columns - clean without symbols
		if txn.Amount.IsNegative() {
			// Debit
			pdf.CellFormat(debitWidth, 5, txn.Amount.Abs().StringFixed(2), "", 0, "R", true, 0, "")
			pdf.CellFormat(creditWidth, 5, "", "", 0, "R", true, 0, "")
		} else {
			// Credit
			pdf.CellFormat(debitWidth, 5, "", "", 0, "R", true, 0, "")
			pdf.CellFormat(creditWidth, 5, txn.Amount.StringFixed(2), "", 0, "R", true, 0, "")
		}

		// Balance
		pdf.CellFormat(balanceWidth, 5, runningBalance.StringFixed(2), "", 1, "R", true, 0, "")

		// Add subtle separator line every 5 transactions for easier reading
		if (i+1)%5 == 0 && i < len(sortedTransactions)-1 {
			pdf.SetDrawColor(borderGrey[0], borderGrey[1], borderGrey[2])
			pdf.SetLineWidth(0.1)
			y := pdf.GetY()
			pdf.Line(20, y, 190, y)
		}
	}

	// Final line after transactions
	pdf.SetDrawColor(borderGrey[0], borderGrey[1], borderGrey[2])
	pdf.SetLineWidth(0.3)
	y = pdf.GetY()
	pdf.Line(20, y, 190, y)

	// Footer section
	pdf.Ln(10)

	// Important notice - minimal
	pdf.SetFont("Helvetica", "", 8)
	pdf.SetTextColor(darkGrey[0], darkGrey[1], darkGrey[2])
	pdf.MultiCell(0, 4, "This is an official bank statement. Please review all transactions and report any discrepancies immediately.", "", "L", false)

	// Page numbering and generation info at the bottom of each page
	pageCount := pdf.PageCount()
	for i := 1; i <= pageCount; i++ {
		pdf.SetPage(i)

		// Page number
		pdf.SetY(-20)
		pdf.SetFont("Helvetica", "", 8)
		pdf.SetTextColor(darkGrey[0], darkGrey[1], darkGrey[2])
		pdf.CellFormat(0, 5, fmt.Sprintf("Page %d of %d", i, pageCount), "", 1, "C", false, 0, "")

		// Generation timestamp
		pdf.SetY(-15)
		pdf.SetFont("Helvetica", "", 7)
		pdf.CellFormat(0, 5, fmt.Sprintf("Generated on %s", s.GeneratedAt.Format("January 2, 2006 at 3:04 PM MST")), "", 0, "C", false, 0, "")
	}

	// Output
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Helper function to format currency without symbols for cleaner look
func formatMinimalAmount(amount decimal.Decimal) string {
	absAmount := amount.Abs()
	formatted := absAmount.StringFixed(2)

	// Add thousand separators
	parts := strings.Split(formatted, ".")
	intPart := parts[0]
	decPart := ""
	if len(parts) > 1 {
		decPart = parts[1]
	}

	// Add commas to integer part
	var result []byte
	for i, digit := range []byte(intPart) {
		if i > 0 && (len(intPart)-i)%3 == 0 {
			result = append(result, ',')
		}
		result = append(result, digit)
	}

	formattedAmount := string(result)
	if decPart != "" {
		formattedAmount += "." + decPart
	}

	return formattedAmount
}