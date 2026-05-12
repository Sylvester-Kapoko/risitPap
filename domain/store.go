package domain

type StoreConfig struct {
	StoreName string
	StoreAddr string
	StorePhone string
	StoreTaxID string
	VATRegistered bool // New: true = we charge vat
	Currency string
	TemplateId string
	LogoPath string
}


