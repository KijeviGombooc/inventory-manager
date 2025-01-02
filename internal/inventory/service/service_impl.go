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
	// check if product sku already exists with different type
	productType, err := trx.GetProductTypeBySku(product.GetBaseProduct().SKU)
	if err != nil {
		return err
	}
	if productType != domain.None && productType != domain.ProductType(product.GetType()) {
		return fmt.Errorf("product with sku %s already exists with different type", product.GetBaseProduct().SKU)
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
	var result domain.IProduct = nil
	switch product.GetType() {
	case dto.Book:
		bookProductDto := product.(*dto.BookProduct)
		result = &domain.BookProduct{
			Author: bookProductDto.Author,
		}
	case dto.Consumable:
		consumableProductDto := product.(*dto.ConsumableProduct)
		result = &domain.ConsumableProduct{
			ExpirationDate: consumableProductDto.ExpirationDate,
		}
	case dto.Electronics:
		electronicsProductDto := product.(*dto.ElectronicsProduct)
		result = &domain.ElectronicsProduct{
			WarrantyPeriod: electronicsProductDto.WarrantyPeriod,
		}
	default:
		return nil, fmt.Errorf("unknown product type")
	}
	baseProductDto := product.GetBaseProduct()
	result.SetBaseProduct(domain.Product{
		SKU:   baseProductDto.SKU,
		Name:  baseProductDto.Name,
		Price: baseProductDto.Price,
		Brand: domain.Brand(baseProductDto.Brand),      // TODO: check conversion error
		Type:  domain.ProductType(baseProductDto.Type), // TODO: check conversion error
	})
	return result, nil
}

func productWithQuantityEntityToDto(productWithQuantity domain.ProductWithQuantity) (dto.ProductWithQuantity, error) {
	result := dto.ProductWithQuantity{
		Quantity: productWithQuantity.Quantity,
	}
	switch productWithQuantity.Product.GetType() {
	case domain.Book:
		bookProductEntity := productWithQuantity.Product.(*domain.BookProduct)
		result.IProduct = &dto.BookProduct{
			Author: bookProductEntity.Author,
		}
	case domain.Consumable:
		consumableProductEntity := productWithQuantity.Product.(*domain.ConsumableProduct)
		result.IProduct = &dto.ConsumableProduct{
			ExpirationDate: consumableProductEntity.ExpirationDate,
		}
	case domain.Electronics:
		electronicsProductEntity := productWithQuantity.Product.(*domain.ElectronicsProduct)
		result.IProduct = &dto.ElectronicsProduct{
			WarrantyPeriod: electronicsProductEntity.WarrantyPeriod,
		}
	default:
		return dto.ProductWithQuantity{}, fmt.Errorf("unknown product type")
	}
	baseProductEntity := productWithQuantity.Product.GetBaseProduct()
	result.IProduct.SetBaseProduct(
		dto.Product{
			SKU:   baseProductEntity.SKU,
			Name:  baseProductEntity.Name,
			Price: baseProductEntity.Price,
			Brand: dto.Brand(baseProductEntity.Brand),      // TODO: check conversion error
			Type:  dto.ProductType(baseProductEntity.Type), // TODO: check conversion error
		},
	)
	return result, nil
}
