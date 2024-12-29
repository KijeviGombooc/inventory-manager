package inventory

import (
	"encoding/json"
	"fmt"
)

type ProductEntity struct {
	SKU   string
	Name  string
	Price int
	Type  ProductType
}

type ProductEntityWithQuantity struct {
	Product  IProduct
	Quantity int
}

type WarehouseProductEntity struct {
	Warehouse string
	Sku       string
	Quantity  int
}

type ProductType string

const (
	Book ProductType = "Book"
)

type BookProductEntity struct {
	ProductEntity
	Author string
}

type BrandEntity struct {
	Name    string
	Quality int
}

type WarehouseEntity struct {
	Name     string
	Address  string
	Capacity int
}

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

func (p ProductEntity) GetSKU() string {
	return p.SKU
}

func (p ProductEntity) GetName() string {
	return p.Name
}

func (p ProductEntity) GetPrice() int {
	return p.Price
}

func (p ProductEntity) GetType() ProductType {
	return p.Type
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
