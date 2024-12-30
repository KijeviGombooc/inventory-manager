package dto

import (
	"encoding/json"
	"fmt"
)

type InsertProductsRequest struct {
	WarehouseName string   `json:"warehouseName"`
	Product       any      `json:"product"`
	ParsedProduct IProduct `json:"-"`
	Quantity      int      `json:"quantity"`
}

func (ipr *InsertProductsRequest) ParseProduct() error {
	productJson, ok := ipr.Product.(map[string]interface{})
	if !ok {
		return fmt.Errorf("product is not a map")
	}
	typeString, ok := productJson["type"]
	if !ok {
		return fmt.Errorf("no type field")
	}
	productTypeStr, ok := typeString.(string)
	if !ok {
		return fmt.Errorf("type is not a string")
	}
	switch ProductType(productTypeStr) {
	case Book:
		var bookProduct BookProduct
		productData, _ := json.Marshal(ipr.Product)
		if err := json.Unmarshal(productData, &bookProduct); err != nil {
			return err
		}
		ipr.ParsedProduct = bookProduct
	default:
		return fmt.Errorf("unknown product type")
	}
	return nil
}
