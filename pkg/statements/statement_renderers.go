package statements

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"html/template"
	"strings"
	"time"

	// "github.com/jung-kurt/gofpdf/v2" // Commented out - using SimplePDFRender instead
	"github.com/shopspring/decimal"
)

// renderPDF internal implementation to avoid circular dependency
func (s *Statement) renderPDF() ([]byte, error) {
	// Check config for renderer preference, default to minimalist
	renderer := "minimalist"
	if s.generator != nil && s.generator.config.PDFRenderer != "" {
		renderer = s.generator.config.PDFRenderer
	}

	switch renderer {
	case "enhanced":
		return s.EnhancedPDFRender()
	case "simple":
		return s.SimplePDFRender()
	case "minimalist":
		fallthrough
	default:
		// Use the enhanced minimalist V2 with footer and logo support
		return s.MinimalistPDFRenderV2()
	}
}

/* Commented out - using SimplePDFRender instead
func (s *Statement) renderPDFHeader(pdf *gofpdf.Fpdf, dateFormatter *DateFormatter, addressFormatter *AddressFormatter) {
	if s.Input.Institution != nil {
		pdf.SetFont("Arial", "B", 18)
		pdf.CellFormat(0, 10, s.Input.Institution.Name, "", 1, "C", false, 0, "")

		if s.Input.Institution.Address != nil {
			pdf.SetFont("Arial", "", 10)
			addressLines := addressFormatter.Format(s.Input.Institution.Address)
			for _, line := range addressLines {
				pdf.CellFormat(0, 5, line, "", 1, "C", false, 0, "")
			}
		}
		pdf.Ln(5)
	}

	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(0, 10, "ACCOUNT STATEMENT", "", 1, "C", false, 0, "")
	pdf.SetDrawColor(44, 62, 80)
	pdf.SetLineWidth(0.5)
	pdf.Line(10, pdf.GetY(), 200, pdf.GetY())
	pdf.Ln(5)
}

func (s *Statement) renderPDFAccountInfo(pdf *gofpdf.Fpdf, dateFormatter *DateFormatter, addressFormatter *AddressFormatter) {
	pdf.SetFont("Arial", "B", 12)
	pdf.SetTextColor(44, 62, 80)
	pdf.CellFormat(0, 8, "Account Information", "", 1, "L", false, 0, "")
	pdf.SetTextColor(0, 0, 0)

	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(40, 6, "Account Holder:", "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(0, 6, s.Input.Account.HolderName, "", 1, "L")

	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(40, 6, "Account Number:", "", 0, "L")
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(0, 6, s.Input.Account.Number, "", 1, "L")

	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(40, 6, "Statement Period:", "", 0, "L")
	pdf.SetFont("Arial", "", 10)
	period := dateFormatter.FormatPeriod(s.Input.PeriodStart, s.Input.PeriodEnd)
	pdf.CellFormat(0, 6, period, "", 1, "L")

	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(40, 6, "Currency:", "", 0, "L")
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(0, 6, s.Input.Account.Currency, "", 1, "L")

	if s.Input.Account.Address != nil {
		pdf.SetFont("Arial", "B", 10)
		pdf.CellFormat(40, 6, "Address:", "", 0, "L")
		pdf.SetFont("Arial", "", 10)

		addressLines := addressFormatter.Format(s.Input.Account.Address)
		for i, line := range addressLines {
			if i == 0 {
				pdf.CellFormat(0, 6, line, "", 1, "L")
			} else {
				pdf.CellFormat(40, 6, "", "", 0, "L")
				pdf.CellFormat(0, 6, line, "", 1, "L")
			}
		}
	}
	pdf.Ln(5)
}

func (s *Statement) renderPDFSummary(pdf *gofpdf.Fpdf, currencyFormatter *CurrencyFormatter) {
	pdf.SetFillColor(240, 240, 240)
	pdf.Rect(10, pdf.GetY(), 190, 35, "F")

	pdf.SetFont("Arial", "B", 12)
	pdf.SetTextColor(44, 62, 80)
	pdf.CellFormat(0, 8, "Account Summary", "", 1, "L")
	pdf.SetTextColor(0, 0, 0)

	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(50, 6, "Opening Balance:", "", 0, "L")
	pdf.SetFont("Arial", "", 10)
	openingBalance := currencyFormatter.FormatAmount(s.Input.OpeningBalance, s.Input.Account.Currency)
	pdf.CellFormat(45, 6, openingBalance, "", 0, "R")

	pdf.SetX(105)
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(50, 6, "Total Credits:", "", 0, "L")
	pdf.SetFont("Arial", "", 10)
	totalCredits := currencyFormatter.FormatAmount(s.TotalCredits, s.Input.Account.Currency)
	pdf.CellFormat(45, 6, totalCredits, "", 1, "R")

	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(50, 6, "Total Debits:", "", 0, "L")
	pdf.SetFont("Arial", "", 10)
	totalDebits := currencyFormatter.FormatAmount(s.TotalDebits, s.Input.Account.Currency)
	pdf.CellFormat(45, 6, totalDebits, "", 0, "R")

	pdf.SetX(105)
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(50, 6, "Transaction Count:", "", 0, "L")
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(45, 6, fmt.Sprintf("%d", s.TransactionCount), "", 1, "R")

	pdf.Ln(2)
	pdf.SetFont("Arial", "B", 12)
	pdf.SetTextColor(44, 62, 80)
	pdf.CellFormat(50, 8, "CLOSING BALANCE:", "", 0, "L")
	closingBalance := currencyFormatter.FormatAmount(s.ClosingBalance, s.Input.Account.Currency)
	pdf.CellFormat(45, 8, closingBalance, "", 1, "R")
	pdf.SetTextColor(0, 0, 0)
	pdf.Ln(8)
}

func (s *Statement) renderPDFTransactions(pdf *gofpdf.Fpdf, dateFormatter *DateFormatter, currencyFormatter *CurrencyFormatter) {
	pdf.SetFont("Arial", "B", 12)
	pdf.SetTextColor(44, 62, 80)
	pdf.CellFormat(0, 8, "Transaction History", "", 1, "L")
	pdf.SetTextColor(0, 0, 0)

	// Table header
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

	// Opening balance row
	pdf.SetFont("Arial", "", 9)
	pdf.CellFormat(20, 6, dateFormatter.Format(s.Input.PeriodStart), "B", 0, "L")
	pdf.CellFormat(65, 6, "Opening Balance", "B", 0, "L")
	pdf.CellFormat(25, 6, "", "B", 0, "R")
	pdf.CellFormat(25, 6, "", "B", 0, "R")
	openingBalance := currencyFormatter.FormatAmount(s.Input.OpeningBalance, s.Input.Account.Currency)
	pdf.CellFormat(30, 6, openingBalance, "B", 0, "R")
	pdf.CellFormat(25, 6, "", "B", 1, "L")

	// Transactions
	calculator := NewCalculator()
	sortedTransactions := calculator.SortTransactions(s.Input.Transactions)
	runningBalance := s.Input.OpeningBalance

	for i, txn := range sortedTransactions {
		if pdf.GetY() > 250 {
			pdf.AddPage()
			// Repeat header
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

		runningBalance = runningBalance.Add(txn.Amount)

		if i%2 == 0 {
			pdf.SetFillColor(250, 250, 250)
		} else {
			pdf.SetFillColor(255, 255, 255)
		}

		pdf.SetFont("Arial", "", 9)
		pdf.CellFormat(20, 6, dateFormatter.Format(txn.Date), "B", 0, "L", true, 0, "")

		description := txn.Description
		if len(description) > 40 {
			description = description[:37] + "..."
		}
		pdf.CellFormat(65, 6, description, "B", 0, "L", true, 0, "")

		debitAmount := ""
		if txn.Type == Debit {
			debitAmount = currencyFormatter.FormatAmount(txn.Amount.Abs(), s.Input.Account.Currency)
		}
		pdf.CellFormat(25, 6, debitAmount, "B", 0, "R", true, 0, "")

		creditAmount := ""
		if txn.Type == Credit {
			creditAmount = currencyFormatter.FormatAmount(txn.Amount, s.Input.Account.Currency)
		}
		pdf.CellFormat(25, 6, creditAmount, "B", 0, "R", true, 0, "")

		balance := currencyFormatter.FormatAmount(runningBalance, s.Input.Account.Currency)
		pdf.SetFont("Arial", "B", 9)
		pdf.CellFormat(30, 6, balance, "B", 0, "R", true, 0, "")

		pdf.SetFont("Arial", "", 8)
		reference := txn.Reference
		if len(reference) > 15 {
			reference = reference[:12] + "..."
		}
		pdf.CellFormat(25, 6, reference, "B", 1, "L", true, 0, "")
	}
}

func (s *Statement) renderPDFFooter(pdf *gofpdf.Fpdf, dateFormatter *DateFormatter) {
	pageCount := pdf.PageCount()
	for i := 1; i <= pageCount; i++ {
		pdf.SetPage(i)
		pdf.SetY(-15)
		pdf.SetFont("Arial", "I", 8)
		pdf.SetTextColor(128, 128, 128)
		pageText := fmt.Sprintf("Page %d of %d", i, pageCount)
		pdf.CellFormat(0, 5, pageText, "", 0, "C")
		pdf.SetY(-10)
		generatedText := fmt.Sprintf("Generated: %s", dateFormatter.FormatDateTime(s.GeneratedAt))
		pdf.CellFormat(0, 5, generatedText, "", 0, "C")
	}
}
*/

// renderCSV internal implementation to avoid circular dependency
func (s *Statement) renderCSV() string {
	locale := "en-US"
	if s.generator != nil {
		locale = s.generator.config.Locale
	}

	dateFormatter := NewDateFormatter(locale)
	currencyFormatter := NewCurrencyFormatter(locale)
	addressFormatter := NewAddressFormatter()

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write header
	writer.Write([]string{"Account Statement"})
	writer.Write([]string{"Account Holder:", s.Input.Account.HolderName})
	writer.Write([]string{"Account Number:", s.Input.Account.Number})

	if s.Input.Account.Type != "" {
		writer.Write([]string{"Account Type:", s.Input.Account.Type})
	}

	period := dateFormatter.FormatPeriod(s.Input.PeriodStart, s.Input.PeriodEnd)
	writer.Write([]string{"Period:", period})
	writer.Write([]string{"Currency:", s.Input.Account.Currency})

	openingBalanceStr := currencyFormatter.FormatForCSV(s.Input.OpeningBalance, s.Input.Account.Currency)
	writer.Write([]string{"Opening Balance:", openingBalanceStr})

	if s.Input.Account.Address != nil {
		addressLine := addressFormatter.FormatSingleLine(s.Input.Account.Address)
		writer.Write([]string{"Address:", addressLine})
	}

	writer.Write([]string{})

	// Transaction headers
	headers := []string{"Date", "Description", "Debit", "Credit", "Balance", "Reference"}
	writer.Write(headers)

	// Opening balance row
	openingRow := []string{
		dateFormatter.Format(s.Input.PeriodStart),
		"Opening Balance",
		"",
		"",
		currencyFormatter.FormatForCSV(s.Input.OpeningBalance, s.Input.Account.Currency),
		"",
	}
	writer.Write(openingRow)

	// Transactions
	calculator := NewCalculator()
	sortedTransactions := calculator.SortTransactions(s.Input.Transactions)
	runningBalance := s.Input.OpeningBalance

	for _, txn := range sortedTransactions {
		runningBalance = runningBalance.Add(txn.Amount)

		debitAmount := ""
		creditAmount := ""

		if txn.Type == Debit {
			debitAmount = currencyFormatter.FormatForCSV(txn.Amount.Abs(), s.Input.Account.Currency)
		} else {
			creditAmount = currencyFormatter.FormatForCSV(txn.Amount, s.Input.Account.Currency)
		}

		row := []string{
			dateFormatter.Format(txn.Date),
			txn.Description,
			debitAmount,
			creditAmount,
			currencyFormatter.FormatForCSV(runningBalance, s.Input.Account.Currency),
			txn.Reference,
		}
		writer.Write(row)
	}

	writer.Write([]string{})

	// Summary
	writer.Write([]string{"Summary"})
	totalCredits := currencyFormatter.FormatForCSV(s.TotalCredits, s.Input.Account.Currency)
	writer.Write([]string{"Total Credits:", totalCredits})

	totalDebits := currencyFormatter.FormatForCSV(s.TotalDebits, s.Input.Account.Currency)
	writer.Write([]string{"Total Debits:", totalDebits})

	closingBalance := currencyFormatter.FormatForCSV(s.ClosingBalance, s.Input.Account.Currency)
	writer.Write([]string{"Closing Balance:", closingBalance})

	writer.Write([]string{"Transaction Count:", fmt.Sprintf("%d", s.TransactionCount)})

	generatedAt := dateFormatter.FormatDateTime(s.GeneratedAt)
	writer.Write([]string{"Generated:", generatedAt})

	writer.Flush()
	return buf.String()
}

// renderHTML internal implementation to avoid circular dependency
func (s *Statement) renderHTML() string {
	locale := "en-US"
	htmlTemplate := ""
	if s.generator != nil {
		locale = s.generator.config.Locale
		htmlTemplate = s.generator.config.HTMLTemplate
	}

	dateFormatter := NewDateFormatter(locale)
	currencyFormatter := NewCurrencyFormatter(locale)
	addressFormatter := NewAddressFormatter()

	// Prepare template data
	data := s.prepareHTMLTemplateData(dateFormatter, currencyFormatter, addressFormatter)

	// Use custom template if provided, otherwise use default
	var tmpl *template.Template
	var err error

	if htmlTemplate != "" {
		tmpl, err = template.New("custom").Funcs(s.getHTMLTemplateFuncs(dateFormatter, currencyFormatter)).Parse(htmlTemplate)
	} else {
		tmpl, err = template.New("default").Funcs(s.getHTMLTemplateFuncs(dateFormatter, currencyFormatter)).Parse(getDefaultHTMLTemplate())
	}

	if err != nil {
		return fmt.Sprintf("Template error: %v", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Sprintf("Template execution error: %v", err)
	}

	return buf.String()
}

func (s *Statement) prepareHTMLTemplateData(dateFormatter *DateFormatter, currencyFormatter *CurrencyFormatter, addressFormatter *AddressFormatter) map[string]interface{} {
	calculator := NewCalculator()
	sortedTransactions := calculator.SortTransactions(s.Input.Transactions)

	runningBalance := s.Input.OpeningBalance
	transactionData := make([]map[string]interface{}, 0, len(sortedTransactions))

	for _, txn := range sortedTransactions {
		runningBalance = runningBalance.Add(txn.Amount)

		txnMap := map[string]interface{}{
			"ID":          txn.ID,
			"Date":        dateFormatter.Format(txn.Date),
			"Description": txn.Description,
			"Type":        string(txn.Type),
			"Reference":   txn.Reference,
			"Balance":     currencyFormatter.FormatAmount(runningBalance, s.Input.Account.Currency),
			"Amount":      currencyFormatter.FormatAmount(txn.Amount.Abs(), s.Input.Account.Currency),
		}

		if txn.Type == Debit {
			txnMap["DebitAmount"] = currencyFormatter.FormatAmount(txn.Amount.Abs(), s.Input.Account.Currency)
			txnMap["CreditAmount"] = ""
		} else {
			txnMap["DebitAmount"] = ""
			txnMap["CreditAmount"] = currencyFormatter.FormatAmount(txn.Amount, s.Input.Account.Currency)
		}

		transactionData = append(transactionData, txnMap)
	}

	var addressLines []string
	if s.Input.Account.Address != nil {
		addressLines = addressFormatter.Format(s.Input.Account.Address)
	}

	var institutionAddressLines []string
	if s.Input.Institution != nil && s.Input.Institution.Address != nil {
		institutionAddressLines = addressFormatter.Format(s.Input.Institution.Address)
	}

	data := map[string]interface{}{
		"Account": map[string]interface{}{
			"Number":       s.Input.Account.Number,
			"HolderName":   s.Input.Account.HolderName,
			"Currency":     s.Input.Account.Currency,
			"Type":         s.Input.Account.Type,
			"AddressLines": addressLines,
		},
		"Period": map[string]interface{}{
			"Start":     dateFormatter.Format(s.Input.PeriodStart),
			"End":       dateFormatter.Format(s.Input.PeriodEnd),
			"Formatted": dateFormatter.FormatPeriod(s.Input.PeriodStart, s.Input.PeriodEnd),
		},
		"OpeningBalance":   currencyFormatter.FormatAmount(s.Input.OpeningBalance, s.Input.Account.Currency),
		"ClosingBalance":   currencyFormatter.FormatAmount(s.ClosingBalance, s.Input.Account.Currency),
		"TotalCredits":     currencyFormatter.FormatAmount(s.TotalCredits, s.Input.Account.Currency),
		"TotalDebits":      currencyFormatter.FormatAmount(s.TotalDebits, s.Input.Account.Currency),
		"Transactions":     transactionData,
		"TransactionCount": s.TransactionCount,
		"GeneratedAt":      dateFormatter.FormatDateTime(s.GeneratedAt),
		"Currency":         s.Input.Account.Currency,
	}

	if s.Input.Institution != nil {
		data["Institution"] = map[string]interface{}{
			"Name":         s.Input.Institution.Name,
			"AddressLines": institutionAddressLines,
			"RegNumber":    s.Input.Institution.RegNumber,
			"TaxID":        s.Input.Institution.TaxID,
		}
	}

	return data
}

func (s *Statement) getHTMLTemplateFuncs(dateFormatter *DateFormatter, currencyFormatter *CurrencyFormatter) template.FuncMap {
	return template.FuncMap{
		"formatCurrency": func(amount decimal.Decimal, currency string) string {
			return currencyFormatter.FormatAmount(amount, currency)
		},
		"formatDate": func(date time.Time) string {
			return dateFormatter.Format(date)
		},
		"formatDateTime": func(date time.Time) string {
			return dateFormatter.FormatDateTime(date)
		},
		"join": strings.Join,
		"upper": strings.ToUpper,
		"lower": strings.ToLower,
	}
}

func getDefaultHTMLTemplate() string {
	// Use the minimalist template as default
	return GetMinimalistHTMLTemplate()
}