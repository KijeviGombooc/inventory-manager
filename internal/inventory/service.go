package inventory

import (
	"fmt"

	"github.com/kijevigombooc/inventory-manager/internal/utils"
)

func NewService(store *store) *service {
	return &service{store: store}
}

type service struct {
	store *store
}

func WarehouseEntityToDto(we WarehouseEntity) WarehouseDto {
	return WarehouseDto(we)
}

func (s *service) GetWarehouses() ([]WarehouseDetailDto, error) {
	trx := s.store.BeginTransaction()
	defer trx.EndTransaction()
	warehouses, err := trx.GetWarehouses()
	if err != nil {
		return nil, err
	}
	result := []WarehouseDetailDto{}
	for _, warehouse := range warehouses {
		productEntities, err := trx.GetProductsByWarehouse(warehouse.Name)
		if err != nil {
			return nil, err
		}
		productDtos, err := utils.MapErrored(productEntities, productEntityWithQuantityToDtoWithQuantity)
		if err != nil {
			return nil, err
		}
		result = append(result, WarehouseDetailDto{
			WarehouseDto: WarehouseEntityToDto(warehouse),
			Products:     productDtos,
		})
	}
	if err := trx.CommitTransaction(); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *service) CreateWarehouse(warehouse WarehouseDto) error {
	trx := s.store.BeginTransaction()
	defer trx.EndTransaction()
	if err := trx.InsertWarehouse(WarehouseEntity(warehouse)); err != nil {
		return err
	}
	if err := trx.CommitTransaction(); err != nil {
		return err
	}
	return nil
}

func (s *service) InsertProducts(warehouse string, product IProduct, quantity int) error {
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
	}
	if remainingQuantity > 0 {
		return fmt.Errorf("not enough capacity in warehouses")
	}
	if err := trx.CommitTransaction(); err != nil {
		return err
	}
	return nil
}

func productDtoToEntity(product IProduct) (IProduct, error) {
	switch product.GetType() {
	case Book:
		bookProduct := product.(BookProductDto)
		return BookProductEntity{
			ProductEntity: ProductEntity(bookProduct.ProductDto),
			Author:        bookProduct.Author,
		}, nil
	default:
		return nil, fmt.Errorf("unknown product type")
	}
}

func productEntityWithQuantityToDtoWithQuantity(productEntity ProductEntityWithQuantity) (ProductDtoWithQuantity, error) {
	switch productEntity.Product.GetType() {
	case Book:
		bookProductEntity := productEntity.Product.(BookProductEntity)
		return ProductDtoWithQuantity{
			IProduct: BookProductDto{
				ProductDto: ProductDto(bookProductEntity.ProductEntity),
				Author:     bookProductEntity.Author,
			},
			Quantity: productEntity.Quantity,
		}, nil
	default:
		return ProductDtoWithQuantity{}, fmt.Errorf("unknown product type")
	}
}

func (s *service) RemoveProducts(sku string, quantity int) error {
	trx := s.store.BeginTransaction()
	defer trx.EndTransaction()
	// TODO: finish this
	return nil
}
