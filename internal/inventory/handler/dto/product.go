package dto

type Product struct {
	SKU   string      `json:"sku"`
	Name  string      `json:"name"`
	Price int         `json:"price"`
	Brand Brand       `json:"brand"`
	Type  ProductType `json:"type"`
}

func (p Product) GetBaseProduct() Product {
	return p
}

func (p *Product) SetBaseProduct(product Product) {
	p.SKU = product.SKU
	p.Name = product.Name
	p.Price = product.Price
	p.Brand = product.Brand
	p.Type = product.Type
}

func (p Product) GetType() ProductType {
	return p.Type
}
