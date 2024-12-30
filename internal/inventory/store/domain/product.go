package domain

type Product struct {
	SKU   string
	Name  string
	Price int
	Type  ProductType
}

func (p Product) GetSKU() string {
	return p.SKU
}

func (p Product) GetName() string {
	return p.Name
}

func (p Product) GetPrice() int {
	return p.Price
}

func (p Product) GetType() ProductType {
	return p.Type
}
