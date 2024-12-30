package domain

type IProduct interface {
	GetBaseProduct() Product
	GetType() ProductType
}
