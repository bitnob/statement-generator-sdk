package statements

import (
	"fmt"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

// CurrencyFormatter handles currency formatting
type CurrencyFormatter struct {
	locale string
}

// NewCurrencyFormatter creates a new currency formatter
func NewCurrencyFormatter(locale string) *CurrencyFormatter {
	if locale == "" {
		locale = "en-US"
	}
	return &CurrencyFormatter{locale: locale}
}

// FormatAmount formats an amount with currency symbol and proper formatting
func (f *CurrencyFormatter) FormatAmount(amount decimal.Decimal, currency string) string {
	// Determine decimal places
	decimalPlaces := 2
	if f.IsZeroDecimalCurrency(currency) {
		decimalPlaces = 0
		amount = amount.Round(0)
	} else {
		amount = amount.Round(2)
	}

	// Get currency symbol
	symbol := f.GetCurrencySymbol(currency)

	// Format based on locale
	var formatted string
	absAmount := amount.Abs()

	switch f.locale {
	case "de-DE", "fr-FR", "it-IT", "es-ES":
		// European format: 1.234,56
		formatted = f.formatEuropean(absAmount, decimalPlaces)
		if symbol != currency {
			if amount.IsNegative() {
				formatted = "-" + symbol + formatted
			} else {
				formatted = symbol + formatted
			}
		}

	case "en-GB":
		// UK format: 1,234.56
		formatted = f.formatAnglo(absAmount, decimalPlaces)
		if amount.IsNegative() {
			formatted = "-" + symbol + formatted
		} else {
			formatted = symbol + formatted
		}

	case "ja-JP":
		// Japanese format: ¥1,234
		formatted = f.formatAnglo(absAmount, decimalPlaces)
		if amount.IsNegative() {
			formatted = "-" + symbol + formatted
		} else {
			formatted = symbol + formatted
		}

	default: // en-US and others
		// US format: $1,234.56 or ($1,234.56) for negative
		formatted = f.formatAnglo(absAmount, decimalPlaces)
		if currency == "USD" && amount.IsNegative() {
			formatted = "(" + symbol + formatted + ")"
		} else if amount.IsNegative() {
			formatted = "-" + symbol + formatted
		} else {
			formatted = symbol + formatted
		}
	}

	return formatted
}

// FormatForCSV formats an amount for CSV export (no currency symbol)
func (f *CurrencyFormatter) FormatForCSV(amount decimal.Decimal, currency string) string {
	decimalPlaces := 2
	if f.IsZeroDecimalCurrency(currency) {
		decimalPlaces = 0
	}
	return amount.StringFixed(int32(decimalPlaces))
}

// formatAnglo formats with Anglo-Saxon conventions (1,234.56)
func (f *CurrencyFormatter) formatAnglo(amount decimal.Decimal, decimalPlaces int) string {
	str := amount.StringFixed(int32(decimalPlaces))

	// Add thousand separators
	parts := strings.Split(str, ".")
	intPart := parts[0]

	// Add commas
	result := ""
	for i, digit := range intPart {
		if i > 0 && (len(intPart)-i)%3 == 0 {
			result += ","
		}
		result += string(digit)
	}

	if len(parts) > 1 && decimalPlaces > 0 {
		result += "." + parts[1]
	}

	return result
}

// formatEuropean formats with European conventions (1.234,56)
func (f *CurrencyFormatter) formatEuropean(amount decimal.Decimal, decimalPlaces int) string {
	str := amount.StringFixed(int32(decimalPlaces))

	// Replace decimal point with comma
	str = strings.Replace(str, ".", ",", 1)

	// Add thousand separators
	parts := strings.Split(str, ",")
	intPart := parts[0]

	// Add dots for thousands
	result := ""
	for i, digit := range intPart {
		if i > 0 && (len(intPart)-i)%3 == 0 {
			result += "."
		}
		result += string(digit)
	}

	if len(parts) > 1 && decimalPlaces > 0 {
		result += "," + parts[1]
	}

	return result
}

// GetCurrencySymbol returns the symbol for a currency
func (f *CurrencyFormatter) GetCurrencySymbol(currency string) string {
	symbols := map[string]string{
		"USD": "$", "EUR": "€", "GBP": "£", "JPY": "¥", "CNY": "¥",
		"INR": "₹", "KRW": "₩", "NGN": "₦", "BRL": "R$", "CAD": "$",
		"AUD": "$", "CHF": "CHF", "SEK": "kr", "NOK": "kr", "DKK": "kr",
		"PLN": "zł", "RUB": "₽", "TRY": "₺", "MXN": "$", "ZAR": "R",
		"SGD": "$", "HKD": "$", "NZD": "$", "THB": "฿", "PHP": "₱",
		"IDR": "Rp", "MYR": "RM", "VND": "₫", "EGP": "E£", "ILS": "₪",
		"CLP": "$", "PEN": "S/", "COP": "$", "ARS": "$", "UYU": "$",
		"PYG": "₲", "BOB": "Bs", "VES": "Bs", "GTQ": "Q", "HNL": "L",
		"NIO": "C$", "CRC": "₡", "PAB": "B/", "DOP": "RD$", "CUP": "₱",
		"JMD": "J$", "TTD": "TT$", "BBD": "$", "BSD": "$", "BZD": "BZ$",
		"KES": "KSh", "GHS": "₵", "ETB": "Br", "UGX": "USh", "TZS": "TSh",
		"RWF": "FRw", "ZMW": "ZK", "BWP": "P", "MUR": "₨", "MAD": "د.م.",
		"TND": "د.ت", "LYD": "ل.د", "SDG": "ج.س.", "IQD": "ع.د", "JOD": "د.ا",
		"KWD": "د.ك", "SAR": "﷼", "AED": "د.إ", "QAR": "﷼", "OMR": "﷼",
		"YER": "﷼", "IRR": "﷼", "SYP": "£", "LBP": "£", "PKR": "₨",
		"LKR": "₨", "NPR": "₨", "BDT": "৳", "MMK": "K", "KHR": "៛",
		"LAK": "₭", "MNT": "₮", "GEL": "₾", "AMD": "֏", "AZN": "₼",
		"KZT": "₸", "UZS": "лв", "KGS": "лв", "TJS": "ЅМ", "TMT": "m",
		"AFN": "؋", "ALL": "L", "BAM": "KM", "BGN": "лв", "HRK": "kn",
		"CZK": "Kč", "HUF": "Ft", "ISK": "kr", "MDL": "L", "MKD": "ден",
		"RON": "lei", "RSD": "Дин.", "UAH": "₴", "BYN": "Br", "GIP": "£",
		"FKP": "£", "SHP": "£", "XOF": "CFA", "XAF": "FCFA", "XCD": "$",
		"XPF": "₣", "AWG": "ƒ", "ANG": "ƒ", "SRD": "$", "GYD": "$",
		"HTG": "G", "GMD": "D", "GNF": "FG", "SLL": "Le", "LRD": "$",
		"CVE": "$", "STN": "Db", "BIF": "FBu", "DJF": "Fdj", "ERN": "Nfk",
		"KMF": "CF", "CDF": "FC", "MRU": "UM", "SCR": "₨", "MVR": "Rf",
		"BTN": "Nu.", "BND": "$", "TWD": "NT$", "MOP": "MOP$", "KPW": "₩",
		"SOS": "S", "SSP": "£", "LSL": "M", "SZL": "E", "NAD": "$",
		"AOA": "Kz", "MZN": "MT", "MGA": "Ar", "TOP": "T$", "WST": "T",
		"SBD": "$", "VUV": "VT", "PGK": "K", "FJD": "$",
	}

	if symbol, ok := symbols[currency]; ok {
		return symbol
	}
	return currency
}

// IsZeroDecimalCurrency checks if a currency has zero decimal places
func (f *CurrencyFormatter) IsZeroDecimalCurrency(currency string) bool {
	zeroDecimal := map[string]bool{
		"JPY": true, "KRW": true, "VND": true, "CLP": true,
		"PYG": true, "UGX": true, "RWF": true, "GNF": true,
		"BIF": true, "DJF": true, "KMF": true, "XOF": true,
		"XAF": true, "XPF": true, "MGA": true, "VUV": true,
	}
	return zeroDecimal[currency]
}

// DateFormatter handles date formatting
type DateFormatter struct {
	locale   string
	timezone *time.Location
}

// NewDateFormatter creates a new date formatter
func NewDateFormatter(locale string) *DateFormatter {
	if locale == "" {
		locale = "en-US"
	}
	return &DateFormatter{
		locale:   locale,
		timezone: time.UTC,
	}
}

// Format formats a date according to locale
func (f *DateFormatter) Format(date time.Time) string {
	switch f.locale {
	case "en-US":
		return date.Format("01/02/2006")
	case "en-GB", "fr-FR", "es-ES", "it-IT":
		return date.Format("02/01/2006")
	case "de-DE":
		return date.Format("02.01.2006")
	case "ja-JP":
		return date.Format("2006/01/02")
	default:
		return date.Format("2006-01-02") // ISO format
	}
}

// FormatDateTime formats a date with time
func (f *DateFormatter) FormatDateTime(date time.Time) string {
	switch f.locale {
	case "en-US":
		return date.Format("01/02/2006 3:04:05 PM")
	case "en-GB":
		return date.Format("02/01/2006 15:04:05")
	case "de-DE":
		return date.Format("02.01.2006 15:04:05")
	case "fr-FR", "es-ES", "it-IT":
		return date.Format("02/01/2006 15:04:05")
	case "ja-JP":
		return date.Format("2006/01/02 15:04:05")
	default:
		return date.Format("2006-01-02 15:04:05")
	}
}

// FormatPeriod formats a date range
func (f *DateFormatter) FormatPeriod(start, end time.Time) string {
	separator := " to "
	if f.locale == "de-DE" {
		separator = " bis "
	}
	return f.Format(start) + separator + f.Format(end)
}

// FormatMonth formats a month name
func (f *DateFormatter) FormatMonth(date time.Time) string {
	return date.Format("January 2006")
}

// AddressFormatter handles address formatting
type AddressFormatter struct{}

// NewAddressFormatter creates a new address formatter
func NewAddressFormatter() *AddressFormatter {
	return &AddressFormatter{}
}

// Format formats an address as multiple lines
func (f *AddressFormatter) Format(address *Address) []string {
	if address == nil {
		return []string{}
	}

	lines := []string{}

	// Line 1
	if address.Line1 != "" {
		lines = append(lines, address.Line1)
	}

	// Line 2 (if present)
	if address.Line2 != "" {
		lines = append(lines, address.Line2)
	}

	// City, State, Postal Code line
	cityLine := ""
	if address.City != "" {
		cityLine = address.City
	}
	if address.State != "" {
		if cityLine != "" {
			cityLine += ", " + address.State
		} else {
			cityLine = address.State
		}
	}
	if address.PostalCode != "" {
		if cityLine != "" {
			cityLine += " " + address.PostalCode
		} else {
			cityLine = address.PostalCode
		}
	}
	if cityLine != "" {
		lines = append(lines, cityLine)
	}

	// Country
	if address.Country != "" {
		lines = append(lines, address.Country)
	}

	return lines
}

// FormatSingleLine formats an address as a single line
func (f *AddressFormatter) FormatSingleLine(address *Address) string {
	if address == nil {
		return ""
	}

	parts := []string{}

	if address.Line1 != "" {
		parts = append(parts, address.Line1)
	}
	if address.Line2 != "" {
		parts = append(parts, address.Line2)
	}
	if address.City != "" {
		parts = append(parts, address.City)
	}
	if address.State != "" {
		parts = append(parts, address.State)
	}
	if address.PostalCode != "" {
		parts = append(parts, address.PostalCode)
	}
	if address.Country != "" {
		parts = append(parts, address.Country)
	}

	return strings.Join(parts, ", ")
}

// FormatForPDF formats an address for PDF rendering
func (f *AddressFormatter) FormatForPDF(address *Address) string {
	lines := f.Format(address)
	return strings.Join(lines, "\n")
}

// NumberFormatter handles number formatting
type NumberFormatter struct {
	locale string
}

// NewNumberFormatter creates a new number formatter
func NewNumberFormatter(locale string) *NumberFormatter {
	if locale == "" {
		locale = "en-US"
	}
	return &NumberFormatter{locale: locale}
}

// FormatInteger formats an integer with thousand separators
func (f *NumberFormatter) FormatInteger(n int) string {
	str := fmt.Sprintf("%d", n)

	// Add thousand separators based on locale
	if f.locale == "de-DE" || f.locale == "fr-FR" {
		// Use dots for thousands
		result := ""
		for i, digit := range str {
			if i > 0 && (len(str)-i)%3 == 0 {
				result += "."
			}
			result += string(digit)
		}
		return result
	}

	// Use commas for thousands (default)
	result := ""
	for i, digit := range str {
		if i > 0 && (len(str)-i)%3 == 0 {
			result += ","
		}
		result += string(digit)
	}
	return result
}