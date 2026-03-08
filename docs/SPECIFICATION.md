# Statement Generator SDK Specification v1.0

## Executive Summary

The Statement Generator SDK is a lightweight, embeddable library for generating bank and financial account statements from transaction data. It produces professional statements in PDF, CSV, and HTML formats without requiring external dependencies or services. This specification defines the requirements for implementing the SDK in any programming language.

## Table of Contents

1. [Overview](#overview)
2. [Core Principles](#core-principles)
3. [Data Models](#data-models)
4. [API Specification](#api-specification)
5. [Output Formats](#output-formats)
6. [Formatting Rules](#formatting-rules)
7. [Implementation Requirements](#implementation-requirements)
8. [Test Cases](#test-cases)
9. [Language Implementation Guide](#language-implementation-guide)
10. [Success Criteria](#success-criteria)

## Overview

### Purpose

Financial institutions and fintech companies need to generate account statements for their customers. This SDK provides a standardized, simple solution that can be embedded directly into existing systems without the complexity of deploying and maintaining a separate service.

### Scope

Version 1.0 focuses on core functionality:
- Transaction list processing
- Balance calculation and validation
- Multi-format output (PDF, CSV, HTML)
- Professional formatting and layout
- Template customization for HTML output

### Target Users

- Fintech companies without existing statement generation
- Financial institutions needing standardized statements
- Development teams requiring embedded statement generation
- Organizations wanting to maintain control over statement generation

## Core Principles

1. **SDK-first Architecture**: Distributed as a library/package, not a service
2. **Zero External Dependencies**: No external services, binaries, or API calls required
3. **Pure Functions**: Stateless operations with no side effects
4. **Minimal API Surface**: Simple to use, difficult to misuse
5. **Format Agnostic**: Same input data produces consistent output across all formats
6. **Internationalization Ready**: Built-in support for multiple currencies and locales

## Data Models

### Core Types

#### Transaction

The fundamental unit of account activity.

```typescript
interface Transaction {
  id: string                    // Unique transaction identifier
  date: datetime                 // Transaction date
  description: string            // Transaction description
  amount: decimal                // Amount (positive for credit, negative for debit)
  type: "credit" | "debit"       // Transaction type
  balance?: decimal              // Balance after transaction (optional)
  reference?: string             // External reference (optional)
}
```

#### Account

Account holder information.

```typescript
interface Account {
  number: string                 // Account number (can be masked, e.g., "****1234")
  holder_name: string            // Account holder's name
  currency: string               // ISO 4217 currency code (NGN, USD, EUR, etc.)
}
```

#### Institution

Organization branding information (optional).

```typescript
interface Institution {
  name: string                   // Institution name
  logo?: bytes                   // Logo image data (optional, PNG/JPEG format)
}
```

#### StatementInput

Complete input data for statement generation.

```typescript
interface StatementInput {
  account: Account               // Account details
  transactions: Transaction[]    // List of transactions
  period_start: date            // Statement start date
  period_end: date              // Statement end date
  opening_balance: decimal       // Balance at period start
  institution?: Institution      // Institution details (optional)
}
```

## API Specification

### Primary Interface

#### StatementGenerator Class

```typescript
class StatementGenerator {
  // Constructor
  constructor(config?: {
    institution?: Institution    // Default institution for all statements
    locale?: string             // Locale for formatting (default: "en-US")
  })

  // Core generation method
  generate(input: StatementInput): Statement

  // Static convenience methods
  static quick_pdf(transactions: Transaction[], opening_balance: decimal): bytes
  static quick_csv(transactions: Transaction[], opening_balance: decimal): string

  // Template customization
  set_html_template(template: string): void
}
```

#### Statement Class

```typescript
class Statement {
  // Calculated properties
  closing_balance: decimal       // Final balance after all transactions
  total_credits: decimal         // Sum of all credits
  total_debits: decimal          // Sum of all debits
  transaction_count: integer     // Number of transactions

  // Export methods
  to_pdf(): bytes               // Generate PDF as byte array
  to_csv(): string              // Generate CSV as string
  to_html(): string             // Generate HTML as string
}
```

### HTML Template Variables

When using custom HTML templates, the following variables are available:

```javascript
{
  account: Account,              // Account information
  institution: Institution,      // Institution information
  period: {                     // Statement period
    start: date,
    end: date
  },
  opening_balance: decimal,     // Starting balance
  closing_balance: decimal,     // Ending balance
  total_credits: decimal,       // Sum of credits
  total_debits: decimal,        // Sum of debits
  transactions: Transaction[],  // Transaction list
  generated_at: datetime        // Generation timestamp
}
```

## Output Formats

### PDF Format

The PDF output should be professional and suitable for official records.

**Required Structure:**
1. Header: Institution branding and statement title
2. Account Information: Account holder details and statement period
3. Summary Section: Opening balance, totals, closing balance
4. Transaction Table: Chronological list with running balance
5. Footer: Page numbering and generation timestamp

**Example Layout:**
```
----------------------------------------
[LOGO] INSTITUTION NAME
        ACCOUNT STATEMENT
----------------------------------------
Account Holder: John Doe
Account Number: ****1234
Statement Period: January 1 - January 31, 2024
Currency: USD

SUMMARY
Opening Balance:             $1,000.00
Total Credits:              +$5,000.00
Total Debits:               -$1,200.00
Closing Balance:             $4,800.00
----------------------------------------
TRANSACTIONS

Date       Description         Debit      Credit     Balance
------------------------------------------------------------
2024-01-01 Opening Balance                           $1,000.00
2024-01-15 Salary                        $5,000.00   $6,000.00
2024-01-16 Rent              $1,200.00              $4,800.00
----------------------------------------
Generated: 2024-02-01 10:30:45 UTC
Page 1 of 1
```

### CSV Format

CSV output should be Excel-compatible with clear sections.

**Required Structure:**
```csv
Account Statement
Account Holder: John Doe
Account Number: ****1234
Period: 2024-01-01 to 2024-01-31
Currency: USD
Opening Balance: 1000.00

Date,Description,Debit,Credit,Balance,Reference
2024-01-01,Opening Balance,,,1000.00,
2024-01-15,Salary,,5000.00,6000.00,SAL-001
2024-01-16,Rent,1200.00,,4800.00,RENT-001

Summary
Total Credits: 5000.00
Total Debits: 1200.00
Closing Balance: 4800.00
Transaction Count: 2
Generated: 2024-02-01 10:30:45 UTC
```

### HTML Format

HTML output must be self-contained with embedded CSS.

**Default Template Structure:**
```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Account Statement</title>
    <style>
        /* Embedded CSS for self-contained document */
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Arial, sans-serif;
            max-width: 800px;
            margin: 0 auto;
            padding: 20px;
            color: #333;
        }
        .header {
            border-bottom: 2px solid #000;
            padding-bottom: 20px;
            margin-bottom: 20px;
        }
        .summary {
            background: #f5f5f5;
            padding: 15px;
            border-radius: 5px;
            margin: 20px 0;
        }
        .summary-row {
            display: flex;
            justify-content: space-between;
            margin: 5px 0;
        }
        table {
            width: 100%;
            border-collapse: collapse;
            margin: 20px 0;
        }
        th {
            background: #f0f0f0;
            padding: 10px;
            text-align: left;
            font-weight: 600;
        }
        td {
            padding: 8px 10px;
            border-bottom: 1px solid #e0e0e0;
        }
        .debit { color: #d00; }
        .credit { color: #0a0; }
        .footer {
            margin-top: 30px;
            padding-top: 20px;
            border-top: 1px solid #e0e0e0;
            font-size: 0.9em;
            color: #666;
        }
        @media print {
            body { margin: 0; padding: 10px; }
            .summary { background: none; border: 1px solid #000; }
        }
    </style>
</head>
<body>
    <div class="header">
        <h1>Account Statement</h1>
        <div>Account Holder: {{account.holder_name}}</div>
        <div>Account Number: {{account.number}}</div>
        <div>Period: {{period.start}} to {{period.end}}</div>
    </div>

    <div class="summary">
        <h2>Summary</h2>
        <div class="summary-row">
            <span>Opening Balance:</span>
            <span>{{opening_balance}}</span>
        </div>
        <div class="summary-row">
            <span>Total Credits:</span>
            <span class="credit">{{total_credits}}</span>
        </div>
        <div class="summary-row">
            <span>Total Debits:</span>
            <span class="debit">{{total_debits}}</span>
        </div>
        <div class="summary-row">
            <strong>Closing Balance:</strong>
            <strong>{{closing_balance}}</strong>
        </div>
    </div>

    <table>
        <thead>
            <tr>
                <th>Date</th>
                <th>Description</th>
                <th>Debit</th>
                <th>Credit</th>
                <th>Balance</th>
            </tr>
        </thead>
        <tbody>
            {{#each transactions}}
            <tr>
                <td>{{date}}</td>
                <td>{{description}}</td>
                <td class="debit">{{debit_amount}}</td>
                <td class="credit">{{credit_amount}}</td>
                <td>{{balance}}</td>
            </tr>
            {{/each}}
        </tbody>
    </table>

    <div class="footer">
        <p>Generated: {{generated_at}}</p>
        <p>Total Transactions: {{transaction_count}}</p>
    </div>
</body>
</html>
```

## Formatting Rules

### Currency Formatting

Implementations must support locale-aware currency formatting:

| Currency | Format Example | Negative Format |
|----------|---------------|-----------------|
| USD | $1,234.56 | ($1,234.56) |
| EUR | €1.234,56 | -€1.234,56 |
| GBP | £1,234.56 | -£1,234.56 |
| NGN | ₦1,234.56 | -₦1,234.56 |
| JPY | ¥1,234 | -¥1,234 |
| Default | 1,234.56 | -1,234.56 |

### Date Formatting

Default to ISO 8601 (YYYY-MM-DD) with locale-specific alternatives:

| Locale | Format |
|--------|--------|
| Default | 2024-01-15 |
| en-US | 01/15/2024 |
| en-GB | 15/01/2024 |
| de-DE | 15.01.2024 |

### Number Precision

- Maintain decimal precision as provided in input
- Use banker's rounding when necessary
- Display 2 decimal places for most currencies
- Display 0 decimal places for zero-decimal currencies (JPY, KRW)

### Balance Calculation Rules

```
If transaction includes balance:
    Use provided balance for display
    Validate consistency with calculated balance

If transaction does not include balance:
    balance[n] = balance[n-1] + amount[n]
    where: credit amounts are positive, debit amounts are negative
```

## Implementation Requirements

### Required Features

1. **Input Validation**
   - Validate transaction dates fall within statement period
   - Validate decimal amounts (no invalid numbers)
   - Validate currency codes against ISO 4217
   - Validate required fields are present

2. **Calculation Engine**
   - Calculate running balance when not provided
   - Sum total credits and debits
   - Validate closing balance matches calculated value
   - Handle decimal precision correctly

3. **Error Handling**
   - Return clear, actionable error messages
   - Fail fast on invalid input
   - No silent failures or data corruption
   - Validate data before generation begins

4. **Performance Requirements**
   - Generate statements with 1,000 transactions in < 1 second
   - Generate statements with 10,000 transactions in < 10 seconds
   - Memory usage should not exceed 100MB for 10,000 transactions
   - Support streaming where applicable

5. **PDF Requirements**
   - Embed fonts (no system font dependencies)
   - Support A4 and Letter page sizes
   - Include page numbers for multi-page statements
   - Support vector logos (not just raster images)

## Test Cases

### Basic Test Case

```json
{
  "description": "Basic statement with credits and debits",
  "input": {
    "account": {
      "number": "****1234",
      "holder_name": "John Doe",
      "currency": "USD"
    },
    "period_start": "2024-01-01",
    "period_end": "2024-01-31",
    "opening_balance": 1000.00,
    "transactions": [
      {
        "id": "TXN001",
        "date": "2024-01-15",
        "description": "Direct Deposit - Salary",
        "amount": 5000.00,
        "type": "credit",
        "reference": "SAL-2024-01"
      },
      {
        "id": "TXN002",
        "date": "2024-01-16",
        "description": "Rent Payment",
        "amount": -1200.00,
        "type": "debit",
        "reference": "RENT-2024-01"
      },
      {
        "id": "TXN003",
        "date": "2024-01-20",
        "description": "Grocery Store",
        "amount": -150.50,
        "type": "debit"
      },
      {
        "id": "TXN004",
        "date": "2024-01-25",
        "description": "Freelance Payment",
        "amount": 800.00,
        "type": "credit",
        "reference": "INV-2024-001"
      }
    ]
  },
  "expected": {
    "closing_balance": 5449.50,
    "total_credits": 5800.00,
    "total_debits": 1350.50,
    "transaction_count": 4
  }
}
```

### Edge Cases

Implementations must handle these scenarios correctly:

1. **Empty Transaction List**
   - Input: No transactions
   - Expected: Opening balance equals closing balance

2. **Single Transaction**
   - Input: One transaction only
   - Expected: Correct balance calculation

3. **All Credits**
   - Input: Only credit transactions
   - Expected: Total debits = 0

4. **All Debits**
   - Input: Only debit transactions
   - Expected: Total credits = 0

5. **Large Numbers**
   - Input: Amounts > 1,000,000,000
   - Expected: Maintain precision

6. **Many Transactions**
   - Input: 10,000+ transactions
   - Expected: Performance within limits

7. **Multiple Same-Day Transactions**
   - Input: Several transactions on same date
   - Expected: Maintain chronological order

8. **Zero Amounts**
   - Input: Transaction with 0.00 amount
   - Expected: Handle gracefully

9. **Negative Opening Balance**
   - Input: Opening balance < 0
   - Expected: Correct calculations

## Language Implementation Guide

### Recommended Project Structure

```
statement-generator/
├── README.md                    # Getting started guide
├── LICENSE                      # MIT or Apache 2.0
├── docs/
│   └── API.md                  # API documentation
├── src/
│   ├── generator.{ext}         # Main StatementGenerator class
│   ├── models.{ext}            # Data models
│   ├── pdf_renderer.{ext}      # PDF generation logic
│   ├── csv_renderer.{ext}      # CSV generation logic
│   ├── html_renderer.{ext}     # HTML generation logic
│   ├── formatters.{ext}        # Currency and date formatting
│   └── calculator.{ext}        # Balance calculations
├── templates/
│   └── default.html            # Default HTML template
├── tests/
│   ├── test_generator.{ext}    # Core tests
│   ├── test_calculations.{ext} # Calculation tests
│   ├── test_formats.{ext}      # Output format tests
│   └── fixtures/               # Test data
└── examples/
    ├── basic_usage.{ext}       # Simple example
    ├── custom_template.{ext}   # HTML template customization
    └── multi_currency.{ext}    # Multiple currency example
```

### Language-Specific Recommendations

#### Go
- Use `shopspring/decimal` for decimal arithmetic
- Use `jung-kurt/gofpdf` or `johnfercher/maroto` for PDF generation
- Embed templates with `embed` package

#### Python
- Use `decimal.Decimal` for money calculations
- Use `reportlab` or `weasyprint` for PDF generation
- Use `jinja2` for HTML templating

#### JavaScript/TypeScript
- Use `decimal.js` or `big.js` for precise decimals
- Use `pdfkit` or `jspdf` for PDF generation
- Use template literals or handlebars for HTML

#### Java
- Use `BigDecimal` for monetary calculations
- Use `iText` or `Apache PDFBox` for PDF
- Use Thymeleaf or similar for templates

#### C#/.NET
- Use `decimal` type for money
- Use `iTextSharp` or `QuestPDF` for PDF
- Use Razor or similar for templates

### Implementation Checklist

Essential features for v1.0:

- [ ] Core data models (Transaction, Account, Institution, StatementInput)
- [ ] StatementGenerator class with configuration
- [ ] Input validation with clear error messages
- [ ] Balance calculation engine
- [ ] PDF generation with embedded fonts
- [ ] CSV generation with proper formatting
- [ ] HTML generation with default template
- [ ] HTML template override capability
- [ ] Currency formatting for major currencies (USD, EUR, GBP, NGN, JPY)
- [ ] Date formatting with locale support
- [ ] Comprehensive test suite covering all edge cases
- [ ] API documentation
- [ ] At least 2 working examples
- [ ] README with quick start guide

## Success Criteria

A compliant v1.0 implementation must:

1. **Functionality**
   - Generate valid PDF, CSV, and HTML from test data
   - Calculate balances and totals correctly
   - Handle all specified edge cases
   - Support template customization for HTML

2. **Performance**
   - Generate 1,000-transaction statement in < 1 second
   - Generate 10,000-transaction statement in < 10 seconds
   - Memory usage < 100MB for 10,000 transactions

3. **Quality**
   - No external service dependencies
   - Thread-safe for concurrent use
   - Comprehensive error handling
   - Clean, maintainable code

4. **Documentation**
   - Complete API documentation
   - Working examples
   - Clear README with quick start
   - Test coverage > 80%

## Version History

| Version | Date | Description |
|---------|------|-------------|
| 1.0 | 2024-03-04 | Initial specification |

## License

This specification is released under Creative Commons CC0 1.0 Universal. Implementations should use MIT or Apache 2.0 license for maximum adoption.

## Contributing

To propose changes to this specification:
1. Open an issue describing the proposed change
2. Submit a pull request with specification updates
3. Changes require review from at least 2 implementers

## Reference Implementations

Official reference implementations:
- Go: [github.com/bitnob/statements-go](https://github.com/bitnob/statements-go) (planned)
- Python: [github.com/bitnob/statements-python](https://github.com/bitnob/statements-python) (planned)

---

*This specification enables any developer to create a compatible Statement Generator SDK in their language of choice while maintaining consistency across implementations.*