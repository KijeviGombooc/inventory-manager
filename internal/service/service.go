package service

import "github.com/kijevigombooc/inventory-manager/internal/inventory/handler/dto"

type Service interface {
	GetWarehouses() ([]dto.WarehouseDetailDto, error)
	CreateWarehouse(warehouse dto.WarehouseDto) error
	InsertProducts(warehouse string, product dto.IProduct, quantity int) error
	RemoveProducts(warehouseName string, sku string, quantity int) error
}
