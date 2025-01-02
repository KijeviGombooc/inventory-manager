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
		if err := unmarshallProduct(ipr, &BookProduct{}); err != nil {
			return err
		}
	case Consumable:
		if err := unmarshallProduct(ipr, &ConsumableProduct{}); err != nil {
			return err
		}
	case Electronics:
		if err := unmarshallProduct(ipr, &ElectronicsProduct{}); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown product type")
	}
	return nil
}

func unmarshallProduct[T IProduct](ipr *InsertProductsRequest, product T) error {
	productData, _ := json.Marshal(ipr.Product)
	if err := json.Unmarshal(productData, &product); err != nil {
		return err
	}
	ipr.ParsedProduct = product
	return nil
}
