package dto

type ConsumableProduct struct {
	Product        `json:",inline"`
	ExpirationDate string `json:"expirationDate"`
}
