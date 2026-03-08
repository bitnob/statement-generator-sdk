package renderers

import (
	"bytes"
	"html/template"
	"strings"
	"time"

	"github.com/statement-generator/sdk/pkg/statements"
	"github.com/shopspring/decimal"
)

// HTML implements the HTML renderer
type HTML struct {
	dateFormatter     *statements.DateFormatter
	currencyFormatter *statements.CurrencyFormatter
	addressFormatter  *statements.AddressFormatter
	template          *template.Template
	customTemplate    string
}

// NewHTML creates a new HTML renderer
func NewHTML(locale string) *HTML {
	h := &HTML{
		dateFormatter:     statements.NewDateFormatter(locale),
		currencyFormatter: statements.NewCurrencyFormatter(locale),
		addressFormatter:  statements.NewAddressFormatter(),
	}
	h.loadDefaultTemplate()
	return h
}

// Render renders a statement as HTML
func (h *HTML) Render(statement *statements.Statement) (interface{}, error) {
	return h.RenderHTML(statement)
}

// RenderHTML renders a statement as HTML string
func (h *HTML) RenderHTML(statement *statements.Statement) (string, error) {
	// Prepare template data
	data := h.prepareTemplateData(statement)

	// Execute template
	var buf bytes.Buffer
	if err := h.template.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// SetTemplate sets a custom HTML template
func (h *HTML) SetTemplate(templateContent string) error {
	h.customTemplate = templateContent

	tmpl, err := template.New("custom").Funcs(h.getTemplateFuncs()).Parse(templateContent)
	if err != nil {
		return err
	}

	h.template = tmpl
	return nil
}

// prepareTemplateData prepares data for template rendering
func (h *HTML) prepareTemplateData(statement *statements.Statement) map[string]interface{} {
	// Calculate running balances
	calculator := statements.NewCalculator()
	sortedTransactions := calculator.SortTransactions(statement.Input.Transactions)

	// Prepare transaction data with running balances
	runningBalance := statement.Input.OpeningBalance
	transactionData := make([]map[string]interface{}, 0, len(sortedTransactions))

	for _, txn := range sortedTransactions {
		runningBalance = runningBalance.Add(txn.Amount)

		txnMap := map[string]interface{}{
			"ID":          txn.ID,
			"Date":        h.dateFormatter.Format(txn.Date),
			"Description": txn.Description,
			"Type":        string(txn.Type),
			"Reference":   txn.Reference,
			"Balance":     h.currencyFormatter.FormatAmount(runningBalance, statement.Input.Account.Currency),
			"Amount":      h.currencyFormatter.FormatAmount(txn.Amount.Abs(), statement.Input.Account.Currency),
		}

		// Set debit/credit amounts
		if txn.Type == statements.Debit {
			txnMap["DebitAmount"] = h.currencyFormatter.FormatAmount(txn.Amount.Abs(), statement.Input.Account.Currency)
			txnMap["CreditAmount"] = ""
		} else {
			txnMap["DebitAmount"] = ""
			txnMap["CreditAmount"] = h.currencyFormatter.FormatAmount(txn.Amount, statement.Input.Account.Currency)
		}

		transactionData = append(transactionData, txnMap)
	}

	// Prepare address data
	var addressLines []string
	if statement.Input.Account.Address != nil {
		addressLines = h.addressFormatter.Format(statement.Input.Account.Address)
	}

	var institutionAddressLines []string
	if statement.Input.Institution != nil && statement.Input.Institution.Address != nil {
		institutionAddressLines = h.addressFormatter.Format(statement.Input.Institution.Address)
	}

	// Build template data
	data := map[string]interface{}{
		"Account": map[string]interface{}{
			"Number":       statement.Input.Account.Number,
			"HolderName":   statement.Input.Account.HolderName,
			"Currency":     statement.Input.Account.Currency,
			"Type":         statement.Input.Account.Type,
			"AddressLines": addressLines,
		},
		"Period": map[string]interface{}{
			"Start":     h.dateFormatter.Format(statement.Input.PeriodStart),
			"End":       h.dateFormatter.Format(statement.Input.PeriodEnd),
			"Formatted": h.dateFormatter.FormatPeriod(statement.Input.PeriodStart, statement.Input.PeriodEnd),
		},
		"OpeningBalance":   h.currencyFormatter.FormatAmount(statement.Input.OpeningBalance, statement.Input.Account.Currency),
		"ClosingBalance":   h.currencyFormatter.FormatAmount(statement.ClosingBalance, statement.Input.Account.Currency),
		"TotalCredits":     h.currencyFormatter.FormatAmount(statement.TotalCredits, statement.Input.Account.Currency),
		"TotalDebits":      h.currencyFormatter.FormatAmount(statement.TotalDebits, statement.Input.Account.Currency),
		"Transactions":     transactionData,
		"TransactionCount": statement.TransactionCount,
		"GeneratedAt":      h.dateFormatter.FormatDateTime(statement.GeneratedAt),
		"Currency":         statement.Input.Account.Currency,
	}

	// Add institution data if available
	if statement.Input.Institution != nil {
		data["Institution"] = map[string]interface{}{
			"Name":         statement.Input.Institution.Name,
			"AddressLines": institutionAddressLines,
			"RegNumber":    statement.Input.Institution.RegNumber,
			"TaxID":        statement.Input.Institution.TaxID,
		}
	}

	return data
}

// getTemplateFuncs returns template functions
func (h *HTML) getTemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"formatCurrency": func(amount decimal.Decimal, currency string) string {
			return h.currencyFormatter.FormatAmount(amount, currency)
		},
		"formatDate": func(date time.Time) string {
			return h.dateFormatter.Format(date)
		},
		"formatDateTime": func(date time.Time) string {
			return h.dateFormatter.FormatDateTime(date)
		},
		"join": strings.Join,
		"upper": strings.ToUpper,
		"lower": strings.ToLower,
	}
}

// loadDefaultTemplate loads the default HTML template
func (h *HTML) loadDefaultTemplate() {
	defaultTemplate := getDefaultHTMLTemplate()
	tmpl, err := template.New("default").Funcs(h.getTemplateFuncs()).Parse(defaultTemplate)
	if err != nil {
		// This should never happen with our default template
		panic("Failed to parse default HTML template: " + err.Error())
	}
	h.template = tmpl
}

// GetDefaultHTMLTemplate returns the default HTML template
func GetDefaultHTMLTemplate() string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Account Statement - {{.Account.HolderName}}</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Helvetica Neue', Arial, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 900px;
            margin: 0 auto;
            padding: 20px;
            background: #fff;
        }
        .header {
            border-bottom: 3px solid #2c3e50;
            padding-bottom: 20px;
            margin-bottom: 30px;
        }
        .header h1 {
            color: #2c3e50;
            margin-bottom: 10px;
            font-size: 28px;
        }
        .institution-info {
            margin-bottom: 20px;
        }
        .institution-name {
            font-size: 20px;
            font-weight: 600;
            color: #2c3e50;
            margin-bottom: 5px;
        }
        .institution-address {
            color: #666;
            font-size: 14px;
        }
        .account-info {
            background: #f8f9fa;
            padding: 20px;
            border-radius: 8px;
            margin-bottom: 30px;
        }
        .account-info h2 {
            color: #2c3e50;
            margin-bottom: 15px;
            font-size: 20px;
        }
        .info-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
            gap: 15px;
        }
        .info-item {
            display: flex;
            flex-direction: column;
        }
        .info-label {
            font-weight: 600;
            color: #666;
            font-size: 12px;
            text-transform: uppercase;
            letter-spacing: 0.5px;
            margin-bottom: 2px;
        }
        .info-value {
            color: #2c3e50;
            font-size: 15px;
        }
        .address {
            margin-top: 10px;
            padding-top: 10px;
            border-top: 1px solid #dee2e6;
        }
        .address-line {
            color: #2c3e50;
            font-size: 14px;
            line-height: 1.4;
        }
        .summary {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 25px;
            border-radius: 8px;
            margin-bottom: 30px;
        }
        .summary h2 {
            margin-bottom: 20px;
            font-size: 20px;
        }
        .summary-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 20px;
        }
        .summary-item {
            display: flex;
            flex-direction: column;
        }
        .summary-label {
            font-size: 12px;
            opacity: 0.9;
            text-transform: uppercase;
            letter-spacing: 0.5px;
            margin-bottom: 5px;
        }
        .summary-value {
            font-size: 24px;
            font-weight: 600;
        }
        .transactions {
            margin-bottom: 30px;
        }
        .transactions h2 {
            color: #2c3e50;
            margin-bottom: 20px;
            font-size: 20px;
        }
        table {
            width: 100%;
            border-collapse: collapse;
            background: white;
            box-shadow: 0 1px 3px rgba(0,0,0,0.1);
            border-radius: 8px;
            overflow: hidden;
        }
        th {
            background: #2c3e50;
            color: white;
            padding: 12px;
            text-align: left;
            font-weight: 600;
            font-size: 14px;
            text-transform: uppercase;
            letter-spacing: 0.5px;
        }
        td {
            padding: 12px;
            border-bottom: 1px solid #e9ecef;
            font-size: 14px;
        }
        tr:last-child td {
            border-bottom: none;
        }
        tr:hover {
            background: #f8f9fa;
        }
        .amount {
            font-weight: 500;
        }
        .debit {
            color: #dc3545;
        }
        .credit {
            color: #28a745;
        }
        .balance {
            font-weight: 600;
            color: #2c3e50;
        }
        .footer {
            margin-top: 40px;
            padding-top: 20px;
            border-top: 2px solid #e9ecef;
            text-align: center;
            color: #6c757d;
            font-size: 12px;
        }
        .footer p {
            margin-bottom: 5px;
        }
        @media print {
            body {
                margin: 0;
                padding: 10px;
            }
            .summary {
                background: #f8f9fa;
                color: #333;
                border: 2px solid #2c3e50;
            }
            .summary-value {
                color: #2c3e50;
            }
            .summary-label {
                color: #666;
            }
            .header {
                page-break-after: avoid;
            }
            .summary {
                page-break-after: avoid;
            }
            table {
                page-break-inside: avoid;
            }
        }
        @media (max-width: 768px) {
            body {
                padding: 10px;
            }
            .info-grid, .summary-grid {
                grid-template-columns: 1fr;
            }
            table {
                font-size: 12px;
            }
            th, td {
                padding: 8px;
            }
        }
    </style>
</head>
<body>
    <div class="header">
        {{if .Institution}}
        <div class="institution-info">
            <div class="institution-name">{{.Institution.Name}}</div>
            {{if .Institution.AddressLines}}
            <div class="institution-address">
                {{range .Institution.AddressLines}}
                <div>{{.}}</div>
                {{end}}
            </div>
            {{end}}
            {{if .Institution.RegNumber}}
            <div class="institution-address">Reg No: {{.Institution.RegNumber}}</div>
            {{end}}
        </div>
        {{end}}
        <h1>Account Statement</h1>
    </div>

    <div class="account-info">
        <h2>Account Information</h2>
        <div class="info-grid">
            <div class="info-item">
                <span class="info-label">Account Holder</span>
                <span class="info-value">{{.Account.HolderName}}</span>
            </div>
            <div class="info-item">
                <span class="info-label">Account Number</span>
                <span class="info-value">{{.Account.Number}}</span>
            </div>
            {{if .Account.Type}}
            <div class="info-item">
                <span class="info-label">Account Type</span>
                <span class="info-value">{{.Account.Type}}</span>
            </div>
            {{end}}
            <div class="info-item">
                <span class="info-label">Statement Period</span>
                <span class="info-value">{{.Period.Formatted}}</span>
            </div>
            <div class="info-item">
                <span class="info-label">Currency</span>
                <span class="info-value">{{.Currency}}</span>
            </div>
        </div>
        {{if .Account.AddressLines}}
        <div class="address">
            <span class="info-label">Address</span>
            {{range .Account.AddressLines}}
            <div class="address-line">{{.}}</div>
            {{end}}
        </div>
        {{end}}
    </div>

    <div class="summary">
        <h2>Account Summary</h2>
        <div class="summary-grid">
            <div class="summary-item">
                <span class="summary-label">Opening Balance</span>
                <span class="summary-value">{{.OpeningBalance}}</span>
            </div>
            <div class="summary-item">
                <span class="summary-label">Total Credits</span>
                <span class="summary-value">{{.TotalCredits}}</span>
            </div>
            <div class="summary-item">
                <span class="summary-label">Total Debits</span>
                <span class="summary-value">{{.TotalDebits}}</span>
            </div>
            <div class="summary-item">
                <span class="summary-label">Closing Balance</span>
                <span class="summary-value">{{.ClosingBalance}}</span>
            </div>
        </div>
    </div>

    <div class="transactions">
        <h2>Transaction History</h2>
        <table>
            <thead>
                <tr>
                    <th>Date</th>
                    <th>Description</th>
                    <th>Debit</th>
                    <th>Credit</th>
                    <th>Balance</th>
                    <th>Reference</th>
                </tr>
            </thead>
            <tbody>
                <tr>
                    <td>{{.Period.Start}}</td>
                    <td>Opening Balance</td>
                    <td></td>
                    <td></td>
                    <td class="balance">{{.OpeningBalance}}</td>
                    <td></td>
                </tr>
                {{range .Transactions}}
                <tr>
                    <td>{{.Date}}</td>
                    <td>{{.Description}}</td>
                    <td class="amount debit">{{.DebitAmount}}</td>
                    <td class="amount credit">{{.CreditAmount}}</td>
                    <td class="balance">{{.Balance}}</td>
                    <td>{{.Reference}}</td>
                </tr>
                {{end}}
            </tbody>
        </table>
    </div>

    <div class="footer">
        <p>Generated on {{.GeneratedAt}}</p>
        <p>Total Transactions: {{.TransactionCount}}</p>
        <p>This is a computer-generated statement and does not require a signature.</p>
        {{if .Institution}}
        <p>&copy; {{.Institution.Name}}. All rights reserved.</p>
        {{end}}
    </div>
</body>
</html>`
}