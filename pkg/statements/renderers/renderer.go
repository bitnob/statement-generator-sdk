package renderers

import (
	"github.com/statement-generator/sdk/pkg/statements"
)

// Renderer is the interface for all statement renderers
type Renderer interface {
	Render(statement *statements.Statement) (interface{}, error)
}

// PDFRenderer renders statements as PDF
type PDFRenderer interface {
	Renderer
	RenderPDF(statement *statements.Statement) ([]byte, error)
}

// CSVRenderer renders statements as CSV
type CSVRenderer interface {
	Renderer
	RenderCSV(statement *statements.Statement) (string, error)
}

// HTMLRenderer renders statements as HTML
type HTMLRenderer interface {
	Renderer
	RenderHTML(statement *statements.Statement) (string, error)
	SetTemplate(template string) error
}