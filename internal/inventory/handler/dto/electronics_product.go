package dto

type ElectronicsProduct struct {
	Product        `json:",inline"`
	WarrantyPeriod string `json:"warrantyPeriod"`
}
