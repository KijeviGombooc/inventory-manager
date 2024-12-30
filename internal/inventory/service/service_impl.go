package service

import (
	"fmt"

	"github.com/kijevigombooc/inventory-manager/internal/inventory/handler/dto"
	"github.com/kijevigombooc/inventory-manager/internal/inventory/store"
	"github.com/kijevigombooc/inventory-manager/internal/inventory/store/domain"
	"github.com/kijevigombooc/inventory-manager/internal/utils"
)

func NewInventoryService(store store.Store) *inventoryService {
	return &inventoryService{store: store}
}

type inventoryService struct {
	store store.Store
}

func (s *inventoryService) GetWarehouses() ([]dto.WarehouseDetail, error) {
	trx := s.store.BeginTransaction()
	defer trx.EndTransaction()
	warehouses, err := trx.GetWarehouses()
	if err != nil {
		return nil, err
	}
	result := []dto.WarehouseDetail{}
	for _, warehouse := range warehouses {
		productEntities, err := trx.GetProductsByWarehouse(warehouse.Name)
		if err != nil {
			return nil, err
		}
		productDtos, err := utils.MapErrored(productEntities, productWithQuantityEntityToDto)
		if err != nil {
			return nil, err
		}
		result = append(result, dto.WarehouseDetail{
			Warehouse: warehouseEntityToDto(warehouse),
			Products:  productDtos,
		})
	}
	if err := trx.CommitTransaction(); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *inventoryService) CreateWarehouse(warehouse dto.Warehouse) error {
	trx := s.store.BeginTransaction()
	defer trx.EndTransaction()
	if err := trx.InsertWarehouse(domain.Warehouse(warehouse)); err != nil {
		return err
	}
	if err := trx.CommitTransaction(); err != nil {
		return err
	}
	return nil
}

func (s *inventoryService) InsertProducts(warehouse string, product dto.IProduct, quantity int) error {
	trx := s.store.BeginTransaction()
	defer trx.EndTransaction()
	warehouses, err := trx.GetWarehousesOrderedFirstWithName(warehouse)
	if err != nil {
		return err
	}
	remainingQuantity := quantity
	for _, warehouse := range warehouses {
		usedCapacity, err := trx.GetUsedCapacity(warehouse.Name)
		if err != nil {
			return err
		}
		availableCapacity := warehouse.Capacity - usedCapacity
		toInsertQuantity := min(availableCapacity, remainingQuantity)
		productEntity, err := productDtoToEntity(product)
		if err != nil {
			return err
		}
		if err := trx.InsertProduct(warehouse.Name, productEntity, toInsertQuantity); err != nil {
			return err
		}
		remainingQuantity -= toInsertQuantity
		if remainingQuantity == 0 {
			break
		}
		if remainingQuantity < 0 {
			return fmt.Errorf("inserted more products than needed")
		}
	}
	if remainingQuantity > 0 {
		return fmt.Errorf("not enough capacity in warehouses")
	}
	if err := trx.CommitTransaction(); err != nil {
		return err
	}
	return nil
}

func (s *inventoryService) RemoveProducts(warehouseName string, sku string, quantity int) error {
	trx := s.store.BeginTransaction()
	defer trx.EndTransaction()

	warehouseProducts, err := trx.GetWarehouseProductsBySkuOrderedFirstWithName(warehouseName, sku)
	if err != nil {
		return err
	}
	remainingQuantity := quantity
	for _, warehouseProduct := range warehouseProducts {
		removedQuantity, err := trx.RemoveProduct(warehouseProduct.WarehouseName, warehouseProduct.Sku, remainingQuantity)
		if err != nil {
			return err
		}
		remainingQuantity -= removedQuantity
		if remainingQuantity == 0 {
			break
		}
		if remainingQuantity < 0 {
			return fmt.Errorf("removed more products than needed")
		}
	}
	if remainingQuantity > 0 {
		return fmt.Errorf("not enough product in warehouses")
	}
	if err := trx.CommitTransaction(); err != nil {
		return err
	}
	return nil
}

func warehouseEntityToDto(we domain.Warehouse) dto.Warehouse {
	return dto.Warehouse(we)
}

func productDtoToEntity(product dto.IProduct) (domain.IProduct, error) {
	switch product.GetType() {
	case dto.Book:
		bookProductDto := product.(dto.BookProduct)
		return domain.BookProduct{
			Product: domain.Product{
				SKU:   bookProductDto.SKU,
				Name:  bookProductDto.Name,
				Price: bookProductDto.Price,
				Brand: domain.Brand(bookProductDto.Brand),      // TODO: check conversion error
				Type:  domain.ProductType(bookProductDto.Type), // TODO: check conversion error
			},
			Author: bookProductDto.Author,
		}, nil
	default:
		return nil, fmt.Errorf("unknown product type")
	}
}

func productWithQuantityEntityToDto(productEntity domain.ProductWithQuantity) (dto.ProductWithQuantity, error) {
	switch productEntity.Product.GetType() {
	case domain.Book:
		bookProductEntity := productEntity.Product.(domain.BookProduct)
		return dto.ProductWithQuantity{
			IProduct: dto.BookProduct{
				Product: dto.Product{
					SKU:   bookProductEntity.SKU,
					Name:  bookProductEntity.Name,
					Price: bookProductEntity.Price,
					Brand: dto.Brand(bookProductEntity.Brand),      // TODO: check conversion error
					Type:  dto.ProductType(bookProductEntity.Type), // TODO: check conversion error
				},
				Author: bookProductEntity.Author,
			},
			Quantity: productEntity.Quantity,
		}, nil
	default:
		return dto.ProductWithQuantity{}, fmt.Errorf("unknown product type")
	}
}
