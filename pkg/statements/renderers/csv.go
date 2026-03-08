package renderers

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strings"

	"github.com/statement-generator/sdk/pkg/statements"
)

// CSV implements the CSV renderer
type CSV struct {
	dateFormatter     *statements.DateFormatter
	currencyFormatter *statements.CurrencyFormatter
	addressFormatter  *statements.AddressFormatter
}

// NewCSV creates a new CSV renderer
func NewCSV(locale string) *CSV {
	return &CSV{
		dateFormatter:     statements.NewDateFormatter(locale),
		currencyFormatter: statements.NewCurrencyFormatter(locale),
		addressFormatter:  statements.NewAddressFormatter(),
	}
}

// Render renders a statement as CSV
func (c *CSV) Render(statement *statements.Statement) (interface{}, error) {
	return c.RenderCSV(statement)
}

// RenderCSV renders a statement as CSV string
func (c *CSV) RenderCSV(statement *statements.Statement) (string, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write header section
	if err := c.writeHeader(writer, statement); err != nil {
		return "", err
	}

	// Empty line
	writer.Write([]string{})

	// Write transaction header
	headers := []string{"Date", "Description", "Debit", "Credit", "Balance", "Reference"}
	if err := writer.Write(headers); err != nil {
		return "", err
	}

	// Write opening balance row
	openingRow := []string{
		c.dateFormatter.Format(statement.Input.PeriodStart),
		"Opening Balance",
		"",
		"",
		c.currencyFormatter.FormatForCSV(statement.Input.OpeningBalance, statement.Input.Account.Currency),
		"",
	}
	if err := writer.Write(openingRow); err != nil {
		return "", err
	}

	// Calculate and write transactions
	calculator := statements.NewCalculator()
	sortedTransactions := calculator.SortTransactions(statement.Input.Transactions)
	runningBalance := statement.Input.OpeningBalance

	for _, txn := range sortedTransactions {
		runningBalance = runningBalance.Add(txn.Amount)

		debitAmount := ""
		creditAmount := ""

		if txn.Type == statements.Debit {
			debitAmount = c.currencyFormatter.FormatForCSV(txn.Amount.Abs(), statement.Input.Account.Currency)
		} else {
			creditAmount = c.currencyFormatter.FormatForCSV(txn.Amount, statement.Input.Account.Currency)
		}

		row := []string{
			c.dateFormatter.Format(txn.Date),
			txn.Description,
			debitAmount,
			creditAmount,
			c.currencyFormatter.FormatForCSV(runningBalance, statement.Input.Account.Currency),
			txn.Reference,
		}
		if err := writer.Write(row); err != nil {
			return "", err
		}
	}

	// Empty line before summary
	writer.Write([]string{})

	// Write summary section
	if err := c.writeSummary(writer, statement); err != nil {
		return "", err
	}

	// Flush the writer
	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// writeHeader writes the header section of the CSV
func (c *CSV) writeHeader(writer *csv.Writer, statement *statements.Statement) error {
	// Title
	if err := writer.Write([]string{"Account Statement"}); err != nil {
		return err
	}

	// Account holder info
	if err := writer.Write([]string{"Account Holder:", statement.Input.Account.HolderName}); err != nil {
		return err
	}

	// Account number
	if err := writer.Write([]string{"Account Number:", statement.Input.Account.Number}); err != nil {
		return err
	}

	// Account type if available
	if statement.Input.Account.Type != "" {
		if err := writer.Write([]string{"Account Type:", statement.Input.Account.Type}); err != nil {
			return err
		}
	}

	// Period
	period := c.dateFormatter.FormatPeriod(statement.Input.PeriodStart, statement.Input.PeriodEnd)
	if err := writer.Write([]string{"Period:", period}); err != nil {
		return err
	}

	// Currency
	if err := writer.Write([]string{"Currency:", statement.Input.Account.Currency}); err != nil {
		return err
	}

	// Opening balance
	openingBalance := c.currencyFormatter.FormatForCSV(statement.Input.OpeningBalance, statement.Input.Account.Currency)
	if err := writer.Write([]string{"Opening Balance:", openingBalance}); err != nil {
		return err
	}

	// Address if available (for proof of address)
	if statement.Input.Account.Address != nil {
		addressLine := c.addressFormatter.FormatSingleLine(statement.Input.Account.Address)
		if err := writer.Write([]string{"Address:", addressLine}); err != nil {
			return err
		}
	}

	// Institution info if available
	if statement.Input.Institution != nil {
		if err := writer.Write([]string{"Institution:", statement.Input.Institution.Name}); err != nil {
			return err
		}
		if statement.Input.Institution.Address != nil {
			instAddressLine := c.addressFormatter.FormatSingleLine(statement.Input.Institution.Address)
			if err := writer.Write([]string{"Institution Address:", instAddressLine}); err != nil {
				return err
			}
		}
		if statement.Input.Institution.RegNumber != "" {
			if err := writer.Write([]string{"Registration Number:", statement.Input.Institution.RegNumber}); err != nil {
				return err
			}
		}
	}

	return nil
}

// writeSummary writes the summary section of the CSV
func (c *CSV) writeSummary(writer *csv.Writer, statement *statements.Statement) error {
	// Summary header
	if err := writer.Write([]string{"Summary"}); err != nil {
		return err
	}

	// Total credits
	totalCredits := c.currencyFormatter.FormatForCSV(statement.TotalCredits, statement.Input.Account.Currency)
	if err := writer.Write([]string{"Total Credits:", totalCredits}); err != nil {
		return err
	}

	// Total debits
	totalDebits := c.currencyFormatter.FormatForCSV(statement.TotalDebits, statement.Input.Account.Currency)
	if err := writer.Write([]string{"Total Debits:", totalDebits}); err != nil {
		return err
	}

	// Closing balance
	closingBalance := c.currencyFormatter.FormatForCSV(statement.ClosingBalance, statement.Input.Account.Currency)
	if err := writer.Write([]string{"Closing Balance:", closingBalance}); err != nil {
		return err
	}

	// Transaction count
	if err := writer.Write([]string{"Transaction Count:", fmt.Sprintf("%d", statement.TransactionCount)}); err != nil {
		return err
	}

	// Generated timestamp
	generatedAt := c.dateFormatter.FormatDateTime(statement.GeneratedAt)
	if err := writer.Write([]string{"Generated:", generatedAt}); err != nil {
		return err
	}

	return nil
}

// RenderSimple creates a simpler CSV format without header
func (c *CSV) RenderSimple(statement *statements.Statement) (string, error) {
	var lines []string

	// Header row
	lines = append(lines, "Date,Description,Amount,Balance,Type,Reference")

	// Calculate running balances
	calculator := statements.NewCalculator()
	sortedTransactions := calculator.SortTransactions(statement.Input.Transactions)
	runningBalance := statement.Input.OpeningBalance

	// Opening balance row
	lines = append(lines, fmt.Sprintf("%s,Opening Balance,,%s,,",
		c.dateFormatter.Format(statement.Input.PeriodStart),
		c.currencyFormatter.FormatForCSV(runningBalance, statement.Input.Account.Currency)))

	// Transaction rows
	for _, txn := range sortedTransactions {
		runningBalance = runningBalance.Add(txn.Amount)

		amount := c.currencyFormatter.FormatForCSV(txn.Amount, statement.Input.Account.Currency)
		balance := c.currencyFormatter.FormatForCSV(runningBalance, statement.Input.Account.Currency)

		lines = append(lines, fmt.Sprintf("%s,%s,%s,%s,%s,%s",
			c.dateFormatter.Format(txn.Date),
			escapeCSV(txn.Description),
			amount,
			balance,
			txn.Type,
			txn.Reference))
	}

	return strings.Join(lines, "\n"), nil
}

// escapeCSV escapes special characters in CSV fields
func escapeCSV(s string) string {
	if strings.ContainsAny(s, ",\"\n\r") {
		s = strings.ReplaceAll(s, "\"", "\"\"")
		s = "\"" + s + "\""
	}
	return s
}