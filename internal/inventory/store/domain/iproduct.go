package domain

type IProduct interface {
	GetSKU() string
	GetName() string
	GetPrice() int
	GetType() ProductType
}
