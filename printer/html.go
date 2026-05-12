package printer

import (
    "fmt"
    "strings"
    "github.com/Sylvester-Kapoko/Receipts/domain"
)

type HtmlFormatter struct{}

func NewHtmlFormatter() *HtmlFormatter {
    return &HtmlFormatter{}
}

func (f *HtmlFormatter) Format(r *domain.Receipt) string {
    var b strings.Builder
    fmt.Fprintf(&b, `<div class="receipt">`)
    fmt.Fprintf(&b, `<h1 class="store-name">%s</h1>`, r.StoreName)
    fmt.Fprintf(&b, `<p class="store-addr">%s</p>`, r.StoreAddr)
    if r.StorePhone != " " {
    	fmt.Fprintf(&b, `<p class="store-phone">Tel: %s </p>`, r.StorePhone)
    }
    if r.StoreTaxID != " " {
    	fmt.Fprintf(&b, `<p class="store-taxID">KRA: %s </p>`, r.StoreTaxID)
    }	
    fmt.Fprintf(&b, `<p class="tx-id">Transaction: %s</p>`, r.TransactionID)
    fmt.Fprintf(&b, `<hr style="border: none; border-top: 2px dotted #000; margin: 10px 0;">`)

    b.WriteString(`<table class="items"><tr><th>Item</th><th>Qty</th><th>Price</th><th>Total</th></tr>`)
    for _, item := range r.Items {
        fmt.Fprintf(&b, `<tr><td>%s</td><td>%d</td><td>%s</td><td>%s</td></tr>`,
            item.Name, item.Qty, domain.DisplayCurrency(item.UnitPrice, r.Currency), domain.DisplayCurrency(item.LineTotal(), r.Currency))
    }
    b.WriteString(`</table>`)

   
    fmt.Fprintf(&b, `<p class="subtotal">Subtotal: %s</p>`, domain.DisplayCurrency(r.Subtotal(), r.Currency))
     if r.HasVAT {
    fmt.Fprintf(&b, `<p class="tax">Tax (%s%%): %s</p>`, r.TaxRatePct(), domain.DisplayCurrency(r.Tax(), r.Currency))
    }
    fmt.Fprintf(&b, `<p class="total">TOTAL: %s</p>`, domain.DisplayCurrency(r.Total(), r.Currency))
    fmt.Fprintf(&b, `<p class="payment">%s: %s</p>`, r.Payment.Method, domain.DisplayCurrency(r.Payment.Amount, r.Currency))
    fmt.Fprintf(&b, `<p class="change">Change: %s</p>`, domain.DisplayCurrency(r.Change(), r.Currency))
    // Updated Footer Section for a standard receipt look
    fmt.Fprintf(&b, `<div class="footer" style="text-align: center; margin-top: 15px; border-top: 1px solid #000; padding-top: 10px; font-family: monospace;">`)
    fmt.Fprintf(&b, `<p>THANK YOU!<br>%s</p>`, r.CreatedAt.Format("2006-01-02 15:04"))
    fmt.Fprintf(&b, `<p class="developer-info" style="font-size: 0.75em; margin-top: 8px; color: #555;">`)
    fmt.Fprintf(&b, `System developed by Sylvester Kapoko<br>0768592677</p>`)
    fmt.Fprintf(&b, `</div>`)
    b.WriteString(`</div>`)
    return b.String()
}
