package dto

type IProduct interface {
	GetBaseProduct() Product
	SetBaseProduct(Product)
	GetType() ProductType
}
