package inventory

type ProductEntity struct {
	SKU    string
	Name   string
	Price  int
	Volume int
	Type   ProductType
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

type ProductDto struct {
	SKU    string      `json:"sku"`
	Name   string      `json:"name"`
	Price  int         `json:"price"`
	Volume int         `json:"volume"`
	Type   ProductType `json:"type"`
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
	Products     []any `json:"products"`
}
