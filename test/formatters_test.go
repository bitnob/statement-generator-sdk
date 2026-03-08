package test

import (
	"testing"
	"time"

	"github.com/statement-generator/sdk/pkg/statements"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestCurrencyFormatter_FormatAmount(t *testing.T) {
	tests := []struct {
		name     string
		amount   decimal.Decimal
		currency string
		locale   string
		expected string
	}{
		// USD formatting
		{
			name:     "USD positive",
			amount:   decimal.NewFromFloat(1234.56),
			currency: "USD",
			locale:   "en-US",
			expected: "$1,234.56",
		},
		{
			name:     "USD negative",
			amount:   decimal.NewFromFloat(-1234.56),
			currency: "USD",
			locale:   "en-US",
			expected: "($1,234.56)",
		},
		{
			name:     "USD zero",
			amount:   decimal.Zero,
			currency: "USD",
			locale:   "en-US",
			expected: "$0.00",
		},
		{
			name:     "USD large number",
			amount:   decimal.NewFromFloat(1234567890.12),
			currency: "USD",
			locale:   "en-US",
			expected: "$1,234,567,890.12",
		},
		// EUR formatting
		{
			name:     "EUR positive",
			amount:   decimal.NewFromFloat(1234.56),
			currency: "EUR",
			locale:   "de-DE",
			expected: "€1.234,56",
		},
		{
			name:     "EUR negative",
			amount:   decimal.NewFromFloat(-1234.56),
			currency: "EUR",
			locale:   "de-DE",
			expected: "-€1.234,56",
		},
		// GBP formatting
		{
			name:     "GBP positive",
			amount:   decimal.NewFromFloat(1234.56),
			currency: "GBP",
			locale:   "en-GB",
			expected: "£1,234.56",
		},
		{
			name:     "GBP negative",
			amount:   decimal.NewFromFloat(-1234.56),
			currency: "GBP",
			locale:   "en-GB",
			expected: "-£1,234.56",
		},
		// NGN formatting
		{
			name:     "NGN positive",
			amount:   decimal.NewFromFloat(1234.56),
			currency: "NGN",
			locale:   "en-NG",
			expected: "₦1,234.56",
		},
		{
			name:     "NGN negative",
			amount:   decimal.NewFromFloat(-1234.56),
			currency: "NGN",
			locale:   "en-NG",
			expected: "-₦1,234.56",
		},
		// JPY formatting (zero decimal currency)
		{
			name:     "JPY positive",
			amount:   decimal.NewFromFloat(1234),
			currency: "JPY",
			locale:   "ja-JP",
			expected: "¥1,234",
		},
		{
			name:     "JPY negative",
			amount:   decimal.NewFromFloat(-1234),
			currency: "JPY",
			locale:   "ja-JP",
			expected: "-¥1,234",
		},
		// KRW formatting (zero decimal currency)
		{
			name:     "KRW positive",
			amount:   decimal.NewFromFloat(1234567),
			currency: "KRW",
			locale:   "ko-KR",
			expected: "₩1,234,567",
		},
		// Default formatting (unknown currency)
		{
			name:     "Unknown currency",
			amount:   decimal.NewFromFloat(1234.56),
			currency: "XYZ",
			locale:   "en-US",
			expected: "1,234.56",
		},
		// Precision handling
		{
			name:     "High precision USD",
			amount:   decimal.NewFromFloat(1234.567890),
			currency: "USD",
			locale:   "en-US",
			expected: "$1,234.57", // Should round to 2 decimal places
		},
		{
			name:     "High precision JPY",
			amount:   decimal.NewFromFloat(1234.567890),
			currency: "JPY",
			locale:   "ja-JP",
			expected: "¥1,235", // Should round to 0 decimal places
		},
		// Edge cases
		{
			name:     "Very small amount",
			amount:   decimal.NewFromFloat(0.01),
			currency: "USD",
			locale:   "en-US",
			expected: "$0.01",
		},
		{
			name:     "Fraction cents",
			amount:   decimal.NewFromFloat(0.005),
			currency: "USD",
			locale:   "en-US",
			expected: "$0.01", // Should round up
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := statements.NewCurrencyFormatter(tt.locale)
			result := formatter.FormatAmount(tt.amount, tt.currency)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCurrencyFormatter_FormatForCSV(t *testing.T) {
	tests := []struct {
		name     string
		amount   decimal.Decimal
		currency string
		expected string
	}{
		{
			name:     "Positive amount",
			amount:   decimal.NewFromFloat(1234.56),
			currency: "USD",
			expected: "1234.56",
		},
		{
			name:     "Negative amount",
			amount:   decimal.NewFromFloat(-1234.56),
			currency: "USD",
			expected: "-1234.56",
		},
		{
			name:     "Zero",
			amount:   decimal.Zero,
			currency: "USD",
			expected: "0.00",
		},
		{
			name:     "JPY (zero decimal)",
			amount:   decimal.NewFromFloat(1234),
			currency: "JPY",
			expected: "1234",
		},
		{
			name:     "High precision",
			amount:   decimal.NewFromFloat(1234.567890),
			currency: "USD",
			expected: "1234.57",
		},
	}

	formatter := statements.NewCurrencyFormatter("en-US")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatter.FormatForCSV(tt.amount, tt.currency)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDateFormatter_Format(t *testing.T) {
	testDate := time.Date(2024, 1, 15, 14, 30, 45, 0, time.UTC)

	tests := []struct {
		name     string
		locale   string
		date     time.Time
		expected string
	}{
		{
			name:     "Default ISO format",
			locale:   "",
			date:     testDate,
			expected: "2024-01-15",
		},
		{
			name:     "US format",
			locale:   "en-US",
			date:     testDate,
			expected: "01/15/2024",
		},
		{
			name:     "UK format",
			locale:   "en-GB",
			date:     testDate,
			expected: "15/01/2024",
		},
		{
			name:     "German format",
			locale:   "de-DE",
			date:     testDate,
			expected: "15.01.2024",
		},
		{
			name:     "French format",
			locale:   "fr-FR",
			date:     testDate,
			expected: "15/01/2024",
		},
		{
			name:     "Japanese format",
			locale:   "ja-JP",
			date:     testDate,
			expected: "2024/01/15",
		},
		{
			name:     "Unknown locale defaults to ISO",
			locale:   "xx-XX",
			date:     testDate,
			expected: "2024-01-15",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := statements.NewDateFormatter(tt.locale)
			result := formatter.Format(tt.date)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDateFormatter_FormatDateTime(t *testing.T) {
	testDate := time.Date(2024, 1, 15, 14, 30, 45, 0, time.UTC)

	tests := []struct {
		name     string
		locale   string
		date     time.Time
		expected string
	}{
		{
			name:     "Default ISO format with time",
			locale:   "",
			date:     testDate,
			expected: "2024-01-15 14:30:45",
		},
		{
			name:     "US format with time",
			locale:   "en-US",
			date:     testDate,
			expected: "01/15/2024 2:30:45 PM",
		},
		{
			name:     "UK format with time",
			locale:   "en-GB",
			date:     testDate,
			expected: "15/01/2024 14:30:45",
		},
		{
			name:     "German format with time",
			locale:   "de-DE",
			date:     testDate,
			expected: "15.01.2024 14:30:45",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := statements.NewDateFormatter(tt.locale)
			result := formatter.FormatDateTime(tt.date)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDateFormatter_FormatPeriod(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC)

	tests := []struct {
		name     string
		locale   string
		start    time.Time
		end      time.Time
		expected string
	}{
		{
			name:     "Default ISO format",
			locale:   "",
			start:    start,
			end:      end,
			expected: "2024-01-01 to 2024-01-31",
		},
		{
			name:     "US format",
			locale:   "en-US",
			start:    start,
			end:      end,
			expected: "01/01/2024 to 01/31/2024",
		},
		{
			name:     "UK format",
			locale:   "en-GB",
			start:    start,
			end:      end,
			expected: "01/01/2024 to 31/01/2024",
		},
		{
			name:     "German format",
			locale:   "de-DE",
			start:    start,
			end:      end,
			expected: "01.01.2024 bis 31.01.2024",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := statements.NewDateFormatter(tt.locale)
			result := formatter.FormatPeriod(tt.start, tt.end)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCurrencyFormatter_GetCurrencySymbol(t *testing.T) {
	tests := []struct {
		currency string
		expected string
	}{
		{"USD", "$"},
		{"EUR", "€"},
		{"GBP", "£"},
		{"NGN", "₦"},
		{"JPY", "¥"},
		{"CNY", "¥"},
		{"INR", "₹"},
		{"KRW", "₩"},
		{"BRL", "R$"},
		{"CAD", "$"},
		{"AUD", "$"},
		{"CHF", "CHF"},
		{"SEK", "kr"},
		{"NOK", "kr"},
		{"DKK", "kr"},
		{"PLN", "zł"},
		{"RUB", "₽"},
		{"TRY", "₺"},
		{"MXN", "$"},
		{"ZAR", "R"},
		{"SGD", "$"},
		{"HKD", "$"},
		{"NZD", "$"},
		{"THB", "฿"},
		{"PHP", "₱"},
		{"IDR", "Rp"},
		{"MYR", "RM"},
		{"VND", "₫"},
		{"XYZ", "XYZ"}, // Unknown currency
		{"", ""},       // Empty currency
	}

	formatter := statements.NewCurrencyFormatter("en-US")

	for _, tt := range tests {
		t.Run(tt.currency, func(t *testing.T) {
			result := formatter.GetCurrencySymbol(tt.currency)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCurrencyFormatter_IsZeroDecimalCurrency(t *testing.T) {
	tests := []struct {
		currency string
		expected bool
	}{
		{"JPY", true},
		{"KRW", true},
		{"VND", true},
		{"CLP", true},
		{"USD", false},
		{"EUR", false},
		{"GBP", false},
		{"NGN", false},
		{"", false},
		{"XYZ", false},
	}

	formatter := statements.NewCurrencyFormatter("en-US")

	for _, tt := range tests {
		t.Run(tt.currency, func(t *testing.T) {
			result := formatter.IsZeroDecimalCurrency(tt.currency)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAddressFormatter_Format(t *testing.T) {
	tests := []struct {
		name     string
		address  *statements.Address
		expected []string
	}{
		{
			name: "Complete address",
			address: &statements.Address{
				Line1:      "123 Main Street",
				Line2:      "Apartment 4B",
				City:       "New York",
				State:      "NY",
				PostalCode: "10001",
				Country:    "USA",
			},
			expected: []string{
				"123 Main Street",
				"Apartment 4B",
				"New York, NY 10001",
				"USA",
			},
		},
		{
			name: "Minimal address",
			address: &statements.Address{
				Line1:   "456 Oak Avenue",
				City:    "Los Angeles",
				Country: "USA",
			},
			expected: []string{
				"456 Oak Avenue",
				"Los Angeles",
				"USA",
			},
		},
		{
			name: "International address",
			address: &statements.Address{
				Line1:      "10 Downing Street",
				City:       "London",
				PostalCode: "SW1A 2AA",
				Country:    "United Kingdom",
			},
			expected: []string{
				"10 Downing Street",
				"London SW1A 2AA",
				"United Kingdom",
			},
		},
		{
			name:     "Nil address",
			address:  nil,
			expected: []string{},
		},
		{
			name: "Address with all fields",
			address: &statements.Address{
				Line1:      "789 Business Park",
				Line2:      "Building C, Suite 300",
				City:       "San Francisco",
				State:      "CA",
				PostalCode: "94105",
				Country:    "United States",
			},
			expected: []string{
				"789 Business Park",
				"Building C, Suite 300",
				"San Francisco, CA 94105",
				"United States",
			},
		},
	}

	formatter := statements.NewAddressFormatter()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatter.Format(tt.address)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAddressFormatter_FormatSingleLine(t *testing.T) {
	tests := []struct {
		name     string
		address  *statements.Address
		expected string
	}{
		{
			name: "Complete address",
			address: &statements.Address{
				Line1:      "123 Main Street",
				Line2:      "Apt 4B",
				City:       "New York",
				State:      "NY",
				PostalCode: "10001",
				Country:    "USA",
			},
			expected: "123 Main Street, Apt 4B, New York, NY 10001, USA",
		},
		{
			name: "Minimal address",
			address: &statements.Address{
				Line1:   "456 Oak Avenue",
				City:    "Los Angeles",
				Country: "USA",
			},
			expected: "456 Oak Avenue, Los Angeles, USA",
		},
		{
			name:     "Nil address",
			address:  nil,
			expected: "",
		},
	}

	formatter := statements.NewAddressFormatter()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatter.FormatSingleLine(tt.address)
			assert.Equal(t, tt.expected, result)
		})
	}
}