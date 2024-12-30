package dto

type WarehouseDetail struct {
	Warehouse `json:",inline"`
	Products  []ProductWithQuantity `json:"products"`
}
