package statements

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/jung-kurt/gofpdf/v2"
	"github.com/shopspring/decimal"
)

// EnhancedPDFRender creates a professional PDF with proper table formatting
func (s *Statement) EnhancedPDFRender() ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")

	// Set margins
	pdf.SetMargins(15, 15, 15)
	pdf.SetAutoPageBreak(true, 20)

	pdf.AddPage()

	// Colors
	headerBg := []int{52, 73, 94}      // Dark blue-gray
	headerText := []int{255, 255, 255} // White
	tableBorder := []int{189, 195, 199} // Light gray
	creditColor := []int{46, 125, 50}   // Green
	debitColor := []int{198, 40, 40}    // Red

	// Institution Header
	if s.Input.Institution != nil {
		pdf.SetFont("Arial", "B", 18)
		pdf.SetTextColor(headerBg[0], headerBg[1], headerBg[2])
		pdf.CellFormat(0, 10, s.Input.Institution.Name, "", 1, "C", false, 0, "")

		if s.Input.Institution.Address != nil {
			pdf.SetFont("Arial", "", 10)
			pdf.SetTextColor(100, 100, 100)
			addr := s.Input.Institution.Address
			addrLine := fmt.Sprintf("%s, %s, %s %s",
				addr.Line1, addr.City, addr.State, addr.PostalCode)
			pdf.CellFormat(0, 5, addrLine, "", 1, "C", false, 0, "")

			if s.Input.Institution.RegNumber != "" {
				pdf.SetFont("Arial", "I", 9)
				pdf.CellFormat(0, 5, s.Input.Institution.RegNumber, "", 1, "C", false, 0, "")
			}
		}
		pdf.Ln(5)
	}

	// Statement Title with colored background
	pdf.SetFillColor(headerBg[0], headerBg[1], headerBg[2])
	pdf.SetTextColor(headerText[0], headerText[1], headerText[2])
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(0, 10, "ACCOUNT STATEMENT", "", 1, "C", true, 0, "")
	pdf.Ln(5)

	// Reset text color
	pdf.SetTextColor(0, 0, 0)

	// Account Information Section with Box
	pdf.SetDrawColor(tableBorder[0], tableBorder[1], tableBorder[2])
	pdf.SetLineWidth(0.5)

	// Left column - Account Details
	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(95, 8, "Account Information", "LTR", 0, "L", false, 0, "")
	pdf.CellFormat(95, 8, "Statement Period", "LTR", 1, "L", false, 0, "")

	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(95, 6, fmt.Sprintf("Account Holder: %s", s.Input.Account.HolderName), "LR", 0, "L", false, 0, "")
	pdf.CellFormat(95, 6, fmt.Sprintf("From: %s", s.Input.PeriodStart.Format("January 2, 2006")), "LR", 1, "L", false, 0, "")

	pdf.CellFormat(95, 6, fmt.Sprintf("Account Number: %s", s.Input.Account.Number), "LR", 0, "L", false, 0, "")
	pdf.CellFormat(95, 6, fmt.Sprintf("To: %s", s.Input.PeriodEnd.Format("January 2, 2006")), "LR", 1, "L", false, 0, "")

	accountType := s.Input.Account.Type
	if accountType == "" {
		accountType = "Account"
	}
	pdf.CellFormat(95, 6, fmt.Sprintf("Type: %s", accountType), "LR", 0, "L", false, 0, "")
	pdf.CellFormat(95, 6, fmt.Sprintf("Currency: %s", s.Input.Account.Currency), "LR", 1, "L", false, 0, "")

	pdf.CellFormat(95, 6, "", "LBR", 0, "L", false, 0, "")
	pdf.CellFormat(95, 6, "", "LBR", 1, "L", false, 0, "")

	// Address Section (if present)
	if s.Input.Account.Address != nil {
		pdf.Ln(3)
		pdf.SetFont("Arial", "B", 11)
		pdf.CellFormat(0, 8, "Account Holder Address (Proof of Address)", "LTR", 1, "L", false, 0, "")

		pdf.SetFont("Arial", "", 10)
		addr := s.Input.Account.Address
		pdf.CellFormat(0, 6, addr.Line1, "LR", 1, "L", false, 0, "")
		if addr.Line2 != "" {
			pdf.CellFormat(0, 6, addr.Line2, "LR", 1, "L", false, 0, "")
		}
		pdf.CellFormat(0, 6, fmt.Sprintf("%s, %s %s", addr.City, addr.State, addr.PostalCode), "LR", 1, "L", false, 0, "")
		pdf.CellFormat(0, 6, addr.Country, "LBR", 1, "L", false, 0, "")
	}

	pdf.Ln(5)

	// Account Summary Box
	pdf.SetFillColor(245, 245, 245)
	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(0, 8, "Account Summary", "1", 1, "C", true, 0, "")

	// Summary details in a grid
	pdf.SetFont("Arial", "", 10)
	formatter := NewCurrencyFormatter("en-US") // Default to en-US for now

	// Opening and Closing Balance
	pdf.CellFormat(95, 7, fmt.Sprintf("Opening Balance: %s",
		formatter.FormatAmount(s.Input.OpeningBalance, s.Input.Account.Currency)), "LBR", 0, "L", false, 0, "")
	pdf.CellFormat(95, 7, fmt.Sprintf("Closing Balance: %s",
		formatter.FormatAmount(s.ClosingBalance, s.Input.Account.Currency)), "LBR", 1, "L", false, 0, "")

	// Credits and Debits
	pdf.SetTextColor(creditColor[0], creditColor[1], creditColor[2])
	pdf.CellFormat(95, 7, fmt.Sprintf("Total Credits: %s",
		formatter.FormatAmount(s.TotalCredits, s.Input.Account.Currency)), "LBR", 0, "L", false, 0, "")

	pdf.SetTextColor(debitColor[0], debitColor[1], debitColor[2])
	pdf.CellFormat(95, 7, fmt.Sprintf("Total Debits: %s",
		formatter.FormatAmount(s.TotalDebits, s.Input.Account.Currency)), "LBR", 1, "L", false, 0, "")

	pdf.SetTextColor(0, 0, 0)
	pdf.Ln(5)

	// Transaction Table Header
	pdf.SetFont("Arial", "B", 11)
	pdf.SetFillColor(headerBg[0], headerBg[1], headerBg[2])
	pdf.SetTextColor(headerText[0], headerText[1], headerText[2])
	pdf.CellFormat(0, 8, "Transaction Details", "1", 1, "C", true, 0, "")

	// Table Headers
	pdf.SetFillColor(240, 240, 240)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Arial", "B", 9)

	// Define column widths
	dateWidth := 25.0
	descWidth := 65.0
	refWidth := 30.0
	amountWidth := 35.0
	balanceWidth := 35.0

	// Table headers
	pdf.CellFormat(dateWidth, 7, "Date", "1", 0, "C", true, 0, "")
	pdf.CellFormat(descWidth, 7, "Description", "1", 0, "C", true, 0, "")
	pdf.CellFormat(refWidth, 7, "Reference", "1", 0, "C", true, 0, "")
	pdf.CellFormat(amountWidth, 7, "Amount", "1", 0, "C", true, 0, "")
	pdf.CellFormat(balanceWidth, 7, "Balance", "1", 1, "C", true, 0, "")

	// Transaction rows
	pdf.SetFont("Arial", "", 8)
	calculator := NewCalculator()
	sortedTransactions := calculator.SortTransactions(s.Input.Transactions)
	runningBalance := s.Input.OpeningBalance

	// Alternating row colors
	rowColor1 := []int{255, 255, 255}
	rowColor2 := []int{250, 250, 250}
	useAltColor := false

	for _, txn := range sortedTransactions {
		// Check for page break
		if pdf.GetY() > 260 {
			pdf.AddPage()
			// Repeat headers on new page
			pdf.SetFont("Arial", "B", 9)
			pdf.SetFillColor(240, 240, 240)
			pdf.SetTextColor(0, 0, 0)
			pdf.CellFormat(dateWidth, 7, "Date", "1", 0, "C", true, 0, "")
			pdf.CellFormat(descWidth, 7, "Description", "1", 0, "C", true, 0, "")
			pdf.CellFormat(refWidth, 7, "Reference", "1", 0, "C", true, 0, "")
			pdf.CellFormat(amountWidth, 7, "Amount", "1", 0, "C", true, 0, "")
			pdf.CellFormat(balanceWidth, 7, "Balance", "1", 1, "C", true, 0, "")
			pdf.SetFont("Arial", "", 8)
		}

		runningBalance = runningBalance.Add(txn.Amount)

		// Set row background color
		if useAltColor {
			pdf.SetFillColor(rowColor2[0], rowColor2[1], rowColor2[2])
		} else {
			pdf.SetFillColor(rowColor1[0], rowColor1[1], rowColor1[2])
		}
		useAltColor = !useAltColor

		// Date
		pdf.CellFormat(dateWidth, 6, txn.Date.Format("Jan 02, 2006"), "1", 0, "C", true, 0, "")

		// Description (truncate if too long)
		desc := txn.Description
		if len(desc) > 35 {
			desc = desc[:32] + "..."
		}
		pdf.CellFormat(descWidth, 6, desc, "1", 0, "L", true, 0, "")

		// Reference
		ref := txn.Reference
		if len(ref) > 15 {
			ref = ref[:12] + "..."
		}
		pdf.CellFormat(refWidth, 6, ref, "1", 0, "C", true, 0, "")

		// Amount (color based on credit/debit)
		if txn.Amount.IsPositive() {
			pdf.SetTextColor(creditColor[0], creditColor[1], creditColor[2])
		} else {
			pdf.SetTextColor(debitColor[0], debitColor[1], debitColor[2])
		}
		pdf.CellFormat(amountWidth, 6, formatter.FormatAmount(txn.Amount, s.Input.Account.Currency), "1", 0, "R", true, 0, "")

		// Balance
		pdf.SetTextColor(0, 0, 0)
		pdf.CellFormat(balanceWidth, 6, formatter.FormatAmount(runningBalance, s.Input.Account.Currency), "1", 1, "R", true, 0, "")
	}

	// Footer
	pdf.Ln(10)
	pdf.SetFont("Arial", "I", 8)
	pdf.SetTextColor(128, 128, 128)

	// Page numbers
	pageCount := pdf.PageCount()
	for i := 1; i <= pageCount; i++ {
		pdf.SetPage(i)
		pdf.SetY(-15)
		pdf.CellFormat(0, 10, fmt.Sprintf("Page %d of %d", i, pageCount), "", 0, "C", false, 0, "")
	}

	// Add generation timestamp on last page
	pdf.SetPage(pageCount)
	pdf.SetY(-20)
	pdf.SetFont("Arial", "I", 7)
	pdf.CellFormat(0, 5, fmt.Sprintf("Generated on %s", s.GeneratedAt.Format("January 2, 2006 at 3:04 PM")), "", 1, "C", false, 0, "")

	// Add disclaimer if institution present
	if s.Input.Institution != nil {
		pdf.SetY(-25)
		pdf.CellFormat(0, 5, "This is an official bank statement. Please keep it for your records.", "", 1, "C", false, 0, "")
	}

	// Output
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Helper function to format addresses
func formatAddressLine(addr *Address) string {
	parts := []string{}
	if addr.Line1 != "" {
		parts = append(parts, addr.Line1)
	}
	if addr.Line2 != "" {
		parts = append(parts, addr.Line2)
	}
	if addr.City != "" {
		parts = append(parts, addr.City)
	}
	if addr.State != "" && addr.PostalCode != "" {
		parts = append(parts, fmt.Sprintf("%s %s", addr.State, addr.PostalCode))
	} else if addr.State != "" {
		parts = append(parts, addr.State)
	} else if addr.PostalCode != "" {
		parts = append(parts, addr.PostalCode)
	}
	if addr.Country != "" {
		parts = append(parts, addr.Country)
	}
	return strings.Join(parts, ", ")
}

// Helper to get appropriate currency symbol
func getCurrencySymbol(currency string) string {
	symbols := map[string]string{
		"USD": "$",
		"EUR": "€",
		"GBP": "£",
		"JPY": "¥",
		"NGN": "₦",
		"CNY": "¥",
		"INR": "₹",
		"CAD": "C$",
		"AUD": "A$",
	}
	if symbol, ok := symbols[currency]; ok {
		return symbol
	}
	return currency + " "
}

// Helper to format amount with symbol
func formatAmountWithSymbol(amount decimal.Decimal, currency string) string {
	symbol := getCurrencySymbol(currency)
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

	if amount.IsNegative() {
		return "-" + symbol + formattedAmount
	}
	return symbol + formattedAmount
}