package dto

type RemoveProductsRequest struct {
	WarehouseName string `json:"warehouseName"`
	Sku           string `json:"sku"`
	Quantity      int    `json:"quantity"`
}
