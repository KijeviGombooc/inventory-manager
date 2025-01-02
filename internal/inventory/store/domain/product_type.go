package domain

type ProductType string

const (
	Book        ProductType = "Book"
	Consumable  ProductType = "Consumable"
	Electronics ProductType = "Electronics"
	None        ProductType = "None"
)
