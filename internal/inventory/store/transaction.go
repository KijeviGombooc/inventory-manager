package store

import (
	"github.com/kijevigombooc/inventory-manager/internal/inventory/store/domain"
)

type Transaction interface {
	CommitTransaction() error
	RollbackTransaction() error
	EndTransaction()
	GetWarehouses() ([]domain.Warehouse, error)
	GetWarehousesOrderedFirstWithName(warehouse string) ([]domain.Warehouse, error)
	InsertWarehouse(entity domain.Warehouse) error
	GetProductsByWarehouse(name string) ([]domain.ProductWithQuantity, error)
	GetUsedCapacity(warehouseName string) (int, error)
	InsertProduct(warehouseName string, product domain.IProduct, toInsertQuantity int) error
	GetWarehouseProductsBySkuOrderedFirstWithName(warehouseName string, sku string) ([]domain.WarehouseProduct, error)
	RemoveProduct(warehouseName string, sku string, toRemoveQuantity int) (int, error)
}
