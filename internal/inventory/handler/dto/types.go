package dto

import (
	"encoding/json"
	"fmt"
)

type ProductType string

const (
	Book ProductType = "Book"
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
		var bookProduct BookProductDto
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

type RemoveProductsRequest struct {
	WarehouseName string `json:"warehouseName"`
	Sku           string `json:"sku"`
	Quantity      int    `json:"quantity"`
}

type ProductDto struct {
	SKU   string      `json:"sku"`
	Name  string      `json:"name"`
	Price int         `json:"price"`
	Type  ProductType `json:"type"`
}

type ProductDtoWithQuantity struct {
	IProduct `json:",inline"`
	Quantity int `json:"quantity"`
}

type BookProductDto struct {
	ProductDto `json:",inline"`
	Author     string `json:"author"`
}

type WarehouseDto struct {
	Name     string `json:"name"`
	Address  string `json:"address"`
	Capacity int    `json:"capacity"`
}

type WarehouseDetailDto struct {
	WarehouseDto `json:",inline"`
	Products     []ProductDtoWithQuantity `json:"products"`
}

type IProduct interface {
	GetSKU() string
	GetName() string
	GetPrice() int
	GetType() ProductType
}

func (p ProductDto) GetSKU() string {
	return p.SKU
}

func (p ProductDto) GetName() string {
	return p.Name
}

func (p ProductDto) GetPrice() int {
	return p.Price
}

func (p ProductDto) GetType() ProductType {
	return p.Type
}
