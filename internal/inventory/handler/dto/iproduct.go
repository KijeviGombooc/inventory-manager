package dto

type IProduct interface {
	GetBaseProduct() Product
	GetType() ProductType
}
