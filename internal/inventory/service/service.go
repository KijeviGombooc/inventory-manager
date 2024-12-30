package service

import "github.com/kijevigombooc/inventory-manager/internal/inventory/handler/dto"

type Service interface {
	GetWarehouses() ([]dto.WarehouseDetail, error)
	CreateWarehouse(warehouse dto.Warehouse) error
	InsertProducts(warehouse string, product dto.IProduct, quantity int) error
	RemoveProducts(warehouseName string, sku string, quantity int) error
}
