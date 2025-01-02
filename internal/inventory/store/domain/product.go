package domain

type Product struct {
	SKU   string
	Name  string
	Price int
	Brand Brand
	Type  ProductType
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
