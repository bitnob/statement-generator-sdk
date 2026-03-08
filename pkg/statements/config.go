package statements

import "time"

// DefaultConfig returns the default configuration for statement generation
// Using minimalist black and white design as the default
func DefaultConfig() *Config {
	return &Config{
		Locale:          "en-US",
		HTMLTemplate:    "",  // Empty means use minimalist template
		PDFRenderer:     "minimalist",
		TimeZone:        time.UTC,
		EnableColors:    false, // Black and white only
		AlternatingRows: true,  // Grey/white alternating rows
	}
}

// RenderMode defines the PDF rendering style
type RenderMode string

const (
	// RenderModeMinimalist uses the clean, black and white minimalist design (DEFAULT)
	RenderModeMinimalist RenderMode = "minimalist"

	// RenderModeEnhanced uses the colorful enhanced design with colored headers
	RenderModeEnhanced RenderMode = "enhanced"

	// RenderModeSimple uses the basic simple design
	RenderModeSimple RenderMode = "simple"
)

// GetPDFRenderer returns the appropriate PDF rendering function based on config
func (s *Statement) GetPDFRenderer() func() ([]byte, error) {
	// Default to minimalist if not specified
	renderMode := RenderModeMinimalist

	if s.generator != nil && s.generator.config.PDFRenderer != "" {
		renderMode = RenderMode(s.generator.config.PDFRenderer)
	}

	switch renderMode {
	case RenderModeEnhanced:
		return s.EnhancedPDFRender
	case RenderModeSimple:
		return s.SimplePDFRender
	case RenderModeMinimalist:
		fallthrough
	default:
		// Minimalist is the default
		return s.MinimalistPDFRender
	}
}

// GetHTMLTemplate returns the appropriate HTML template based on config
func GetHTMLTemplate(config *Config) string {
	// If custom template provided, use it
	if config != nil && config.HTMLTemplate != "" {
		return config.HTMLTemplate
	}

	// Default to minimalist template
	return GetMinimalistHTMLTemplate()
}

// DefaultRenderOptions returns the default rendering options
func DefaultRenderOptions() *RenderOptions {
	return &RenderOptions{
		EnableColors:    false, // Black and white by default
		AlternatingRows: true,  // Grey/white alternating rows
		BorderStyle:     "thin", // Thin borders for minimalist look
		FontFamily:      "Helvetica",
		FontSize:        10,
	}
}

// RenderOptions contains options for rendering statements
type RenderOptions struct {
	EnableColors    bool   // Enable colored text (red/green for debits/credits)
	AlternatingRows bool   // Enable alternating row colors
	BorderStyle     string // Border style: "none", "thin", "thick"
	FontFamily      string // Font family for PDF
	FontSize        int    // Base font size for PDF
}