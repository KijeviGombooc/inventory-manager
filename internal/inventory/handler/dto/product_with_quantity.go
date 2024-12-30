package dto

type ProductWithQuantity struct {
	IProduct `json:",inline"`
	Quantity int `json:"quantity"`
}
