package dto

type Product struct {
	SKU   string      `json:"sku"`
	Name  string      `json:"name"`
	Price int         `json:"price"`
	Type  ProductType `json:"type"`
}
