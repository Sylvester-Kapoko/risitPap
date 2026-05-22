package printer

import (
	"bytes"
	"embed"
	"html/template"
	"path"
	"strings"

	"github.com/Sylvester-Kapoko/risitPap/domain" // corrected import
)

//go:embed templates/*
var templateFS embed.FS

func parseTemplate(filename string) *template.Template {
	return template.Must(template.New(path.Base(filename)).ParseFS(templateFS, "templates/"+filename))
}

type HtmlConfig struct {
	DeveloperName  string
	DeveloperPhone string
	ShowFooter     bool
}

type HtmlFormatter struct {
	fullTmpl     *template.Template
	fragmentTmpl *template.Template
	statusTmpl   *template.Template
	config       HtmlConfig
}

func NewHtmlFormatter(cfg HtmlConfig) *HtmlFormatter {
	return &HtmlFormatter{
		fullTmpl:     parseTemplate("receipt.html"),
		fragmentTmpl: parseTemplate("receipt_fragment.html"),
		statusTmpl:   parseTemplate("status_fragment.html"),
		config:       cfg,
	}
}

func nonEmpty(s string) bool {
	return strings.TrimSpace(s) != ""
}

// --------------------------------------------------------------------
//  View types
// --------------------------------------------------------------------

type statusView struct {
	PaymentMethod string
	PaymentAmount string
	Change        string
	Total         string
	State         domain.PaymentStatus // using domain type (which is a string)
}

type receiptView struct {
	ID            string
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

// --------------------------------------------------------------------
//  Format methods
// --------------------------------------------------------------------

func (f *HtmlFormatter) Format(r *domain.Receipt) (string, error) {
	items := make([]itemView, 0, len(r.Items))
	for _, it := range r.Items {
		items = append(items, itemView{
			Name:      it.Name,
			Qty:       it.Qty,
			UnitPrice: domain.DisplayCurrency(it.UnitPrice, r.Currency),
			LineTotal: domain.DisplayCurrency(it.LineTotal(), r.Currency),
		})
	}

	view := receiptView{
		ID:            r.ID,
		StoreName:     r.StoreName,
		StoreAddr:     r.StoreAddr,
		StorePhone:    r.StorePhone,
		StoreTaxID:    r.StoreTaxID,
		TransactionID: r.TransactionID,
		Currency:      r.Currency,
		Items:         items,

		Subtotal: domain.DisplayCurrency(r.Subtotal(), r.Currency),
		TaxRate:  r.TaxRatePct(),
		Tax:      domain.DisplayCurrency(r.Tax(), r.Currency),
		Total:    domain.DisplayCurrency(r.Total(), r.Currency),

		PaymentMethod: r.Payment.Method,
		PaymentAmount: domain.DisplayCurrency(r.Payment.Amount, r.Currency),
		Change:        domain.DisplayCurrency(r.Change(), r.Currency),

		CreatedAt: r.CreatedAt.Format("2006-01-02 15:04"),

		ShowPhone: nonEmpty(r.StorePhone),
		ShowTaxID: nonEmpty(r.StoreTaxID),

		Footer: footerView{
			Show:           f.config.ShowFooter,
			DeveloperName:  f.config.DeveloperName,
			DeveloperPhone: f.config.DeveloperPhone,
			CreatedAt:      r.CreatedAt.Format("2006-01-02 15:04"),
		},
	}

	var buf bytes.Buffer
	err := f.fullTmpl.Execute(&buf, view)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (f *HtmlFormatter) RenderFragment(r *domain.Receipt) (string, error) {
	// For v1 we just call Format – you can later make a lighter fragment template.
	return f.Format(r)
}

func (f *HtmlFormatter) RenderStatusFragment(r *domain.Receipt) (string, error) {
	var buf bytes.Buffer

	view := statusView{
		PaymentMethod: r.Payment.Method,
		PaymentAmount: domain.DisplayCurrency(r.Payment.Amount, r.Currency),
		Change:        domain.DisplayCurrency(r.Change(), r.Currency),
		Total:         domain.DisplayCurrency(r.Total(), r.Currency),
		State:         r.PaymentState(), // returns domain.PaymentStatus
	}

	if err := f.statusTmpl.Execute(&buf, view); err != nil {
		return "", err
	}

	return buf.String(), nil
}
