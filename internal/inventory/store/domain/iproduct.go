package domain

type IProduct interface {
	GetBaseProduct() Product
	SetBaseProduct(product Product)
	GetType() ProductType
}
