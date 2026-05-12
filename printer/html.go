package printer

import (
	"bytes"
	"html/template"
	"strings"

	"github.com/Sylvester-Kapoko/Receipts/domain"
)

type HtmlConfig struct {
	DeveloperName  string
	DeveloperPhone string
	ShowFooter     bool
}

type HtmlFormatter struct {
    fullTmpl    *template.Template
    statusTmpl  *template.Template
    config      HtmlConfig
}


func NewHtmlFormatter(cfg HtmlConfig) *HtmlFormatter {
    return &HtmlFormatter{
        fullTmpl: template.Must(template.New("receipt").Parse(fullTemplate)),
        statusTmpl: template.Must(template.New("status").Parse(statusTemplate)),
        config: cfg,
    }
}

func nonEmpty(s string) bool {
    return strings.TrimSpace(s) != ""
}

type receiptView struct { ... }
type itemView struct { ... }
type footerView struct { ... }


type receiptView struct {
	StoreName     string
	StoreAddr     string
	StorePhone    string
	StoreTaxID    string
	TransactionID string
	Currency      string

	Items []itemView

	Subtotal string
	TaxRate  string
	Tax      string
	Total    string

	PaymentMethod string
	PaymentAmount string
	Change        string

	CreatedAt string

	ShowPhone bool
	ShowTaxID bool

	Footer footerView
}

type itemView struct {
	Name      string
	Qty       int
	UnitPrice string
	LineTotal string
}

type footerView struct {
	Show           bool
	DeveloperName  string
	DeveloperPhone string
	CreatedAt      string
}

func nonEmpty(s string) bool {
	return strings.TrimSpace(s) != ""
}

<div id="payment-status">

<p><b>Payment:</b> {{.PaymentMethod}}</p>
<p><b>Paid:</b> {{.PaymentAmount}}</p>
<p><b>Change:</b> {{.Change}}</p>
<p><b>Total:</b> {{.Total}}</p>

</div>

func (f *HtmlFormatter) RenderStatusFragment(r *domain.Receipt) (string, error)
const fullTemplate = `...`
const statusTemplate = `...`

