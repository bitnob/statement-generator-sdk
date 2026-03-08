package statements

// GetMinimalistHTMLTemplate returns a clean, minimalist HTML template for statements
func GetMinimalistHTMLTemplate() string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Statement - {{.Account.HolderName}}</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: -apple-system, 'Helvetica Neue', Arial, sans-serif;
            max-width: 900px;
            margin: 0 auto;
            padding: 40px 20px;
            color: #000;
            background: #fff;
            line-height: 1.5;
        }

        /* Institution Header */
        .institution-header {
            margin-bottom: 30px;
        }
        .institution-name {
            font-size: 20px;
            font-weight: bold;
            text-transform: uppercase;
            letter-spacing: 0.5px;
            margin-bottom: 8px;
        }
        .institution-header hr {
            border: none;
            border-bottom: 0.5px solid #dcdcdc;
            margin: 10px 0 8px 0;
        }
        .institution-details {
            font-size: 10px;
            color: #646464;
            margin-bottom: 5px;
        }

        /* Statement Title */
        .statement-title {
            font-size: 14px;
            font-weight: normal;
            margin: 25px 0 20px 0;
        }

        /* Account Info Grid */
        .account-info {
            margin-bottom: 25px;
        }
        .info-row {
            display: flex;
            margin-bottom: 6px;
            font-size: 11px;
        }
        .info-label {
            color: #646464;
            width: 120px;
        }
        .info-value {
            color: #000;
            flex: 1;
        }
        .info-value.bold {
            font-weight: 600;
        }
        .period-info {
            display: flex;
            gap: 60px;
        }

        /* Address Section */
        .address-section {
            margin: 20px 0;
        }

        /* Summary Section */
        .summary {
            border-top: 0.5px solid #dcdcdc;
            border-bottom: 0.5px solid #dcdcdc;
            padding: 20px 0;
            margin: 25px 0;
        }
        .summary-title {
            font-size: 12px;
            margin-bottom: 15px;
            font-weight: normal;
        }
        .summary-grid {
            display: grid;
            grid-template-columns: repeat(2, 1fr);
            gap: 8px 60px;
        }
        .summary-item {
            display: flex;
            justify-content: space-between;
            font-size: 11px;
        }
        .summary-item .label {
            color: #646464;
        }
        .summary-item .value {
            color: #000;
            text-align: right;
            font-family: 'Courier New', monospace;
        }
        .summary-item.closing .value {
            font-weight: bold;
        }

        /* Transactions Section */
        .transactions-section {
            margin-top: 30px;
        }
        .transactions-title {
            font-size: 12px;
            margin-bottom: 8px;
            font-weight: normal;
        }
        .transactions-divider {
            border: none;
            border-bottom: 0.5px solid #dcdcdc;
            margin-bottom: 10px;
        }

        /* Table Styles */
        table {
            width: 100%;
            border-collapse: collapse;
            font-size: 10px;
        }
        thead th {
            text-align: left;
            padding: 6px 4px;
            border-bottom: 0.5px solid #646464;
            color: #646464;
            font-weight: normal;
            font-size: 9px;
            letter-spacing: 0.3px;
        }
        thead th.amount {
            text-align: right;
            padding-right: 8px;
        }

        /* Alternating row colors */
        tbody tr:nth-child(odd) {
            background: #f8f8f8;
        }
        tbody tr:nth-child(even) {
            background: #ffffff;
        }
        tbody td {
            padding: 5px 4px;
            color: #000;
            font-size: 10px;
        }
        tbody td.date {
            white-space: nowrap;
            width: 60px;
        }
        tbody td.description {
            max-width: 250px;
            overflow: hidden;
            text-overflow: ellipsis;
        }
        tbody td.reference {
            color: #646464;
            font-size: 9px;
        }
        tbody td.amount {
            text-align: right;
            font-family: 'Courier New', monospace;
            padding-right: 8px;
            width: 80px;
        }
        tbody td.balance {
            text-align: right;
            font-family: 'Courier New', monospace;
            padding-right: 8px;
            width: 90px;
        }

        /* Subtle separator every 5 rows */
        tbody tr:nth-child(5n) td {
            border-bottom: 0.5px solid #f0f0f0;
        }

        /* Footer */
        .footer {
            margin-top: 40px;
            padding-top: 20px;
            border-top: 0.5px solid #dcdcdc;
            font-size: 9px;
            color: #646464;
            text-align: center;
        }
        .footer p {
            margin: 3px 0;
        }
        .footer .page-info {
            margin-top: 10px;
            font-size: 8px;
        }

        /* Print Styles */
        @media print {
            body {
                padding: 20px;
                font-size: 10px;
            }
            .summary {
                page-break-inside: avoid;
            }
            table {
                page-break-inside: auto;
            }
            tr {
                page-break-inside: avoid;
                page-break-after: auto;
            }
            tbody tr:nth-child(odd) {
                background: #f8f8f8 !important;
                -webkit-print-color-adjust: exact;
                print-color-adjust: exact;
            }
        }
    </style>
</head>
<body>
    {{if .Institution}}
    <div class="institution-header">
        <div class="institution-name">{{upper .Institution.Name}}</div>
        <hr>
        {{if .Institution.AddressLines}}
        <div class="institution-details">
            {{join .Institution.AddressLines ", "}}
        </div>
        {{end}}
        {{if .Institution.RegNumber}}
        <div class="institution-details">{{.Institution.RegNumber}}</div>
        {{end}}
    </div>
    {{end}}

    <div class="statement-title">Statement of Account</div>

    <div class="account-info">
        <div class="info-row">
            <div class="info-label">Account Holder</div>
            <div class="info-value bold">{{.Account.HolderName}}</div>
            <div class="period-info">
                <div class="info-label">Statement Period</div>
                <div class="info-value">{{.Period.Formatted}}</div>
            </div>
        </div>
        <div class="info-row">
            <div class="info-label">Account Number</div>
            <div class="info-value">{{.Account.Number}}</div>
            <div class="period-info">
                <div class="info-label">Currency</div>
                <div class="info-value">{{.Account.Currency}}</div>
            </div>
        </div>
        {{if .Account.Type}}
        <div class="info-row">
            <div class="info-label">Account Type</div>
            <div class="info-value">{{.Account.Type}}</div>
        </div>
        {{end}}
    </div>

    {{if .Account.AddressLines}}
    <div class="address-section">
        <div class="info-row">
            <div class="info-label">Mailing Address</div>
            <div class="info-value">
                {{range $i, $line := .Account.AddressLines}}
                    {{if $i}}, {{end}}{{$line}}
                {{end}}
            </div>
        </div>
    </div>
    {{end}}

    <div class="summary">
        <div class="summary-title">Account Summary</div>
        <div class="summary-grid">
            <div class="summary-item">
                <span class="label">Opening Balance</span>
                <span class="value">{{.OpeningBalance}}</span>
            </div>
            <div class="summary-item">
                <span class="label">Total Credits</span>
                <span class="value">{{.TotalCredits}}</span>
            </div>
            <div class="summary-item">
                <span class="label">Total Debits</span>
                <span class="value">{{.TotalDebits}}</span>
            </div>
            <div class="summary-item closing">
                <span class="label">Closing Balance</span>
                <span class="value">{{.ClosingBalance}}</span>
            </div>
        </div>
    </div>

    <div class="transactions-section">
        <div class="transactions-title">Transaction Details</div>
        <hr class="transactions-divider">
        <table>
            <thead>
                <tr>
                    <th>Date</th>
                    <th>Description</th>
                    <th>Reference</th>
                    <th class="amount">Debit</th>
                    <th class="amount">Credit</th>
                    <th class="amount">Balance</th>
                </tr>
            </thead>
            <tbody>
                <tr>
                    <td class="date">{{.Period.Start}}</td>
                    <td class="description">Opening Balance</td>
                    <td class="reference"></td>
                    <td class="amount"></td>
                    <td class="amount"></td>
                    <td class="balance">{{.OpeningBalance}}</td>
                </tr>
                {{range .Transactions}}
                <tr>
                    <td class="date">{{.Date}}</td>
                    <td class="description">{{.Description}}</td>
                    <td class="reference">{{.Reference}}</td>
                    <td class="amount">{{.DebitAmount}}</td>
                    <td class="amount">{{.CreditAmount}}</td>
                    <td class="balance">{{.Balance}}</td>
                </tr>
                {{end}}
            </tbody>
        </table>
    </div>

    <div class="footer">
        <p>This is an official bank statement. Please review all transactions and report any discrepancies immediately.</p>
        <p class="page-info">Generated on {{.GeneratedAt}} | Total Transactions: {{.TransactionCount}}</p>
        {{if .Institution}}<p>&copy; {{.Institution.Name}}. All rights reserved.</p>{{end}}
    </div>
</body>
</html>`
}