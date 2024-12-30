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

func (p Product) GetType() ProductType {
	return p.Type
}
