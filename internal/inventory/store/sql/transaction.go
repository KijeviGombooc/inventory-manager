package sql

import (
	"database/sql"
	"fmt"

	"github.com/kijevigombooc/inventory-manager/internal/inventory/store/domain"
	"github.com/kijevigombooc/inventory-manager/internal/inventory/store/sql/query"
)

type SqlTransaction struct {
	tx       *sql.Tx
	commited bool
}

func (t *SqlTransaction) CommitTransaction() error {
	t.commited = true
	return t.tx.Commit()
}

func (t *SqlTransaction) RollbackTransaction() error {
	return t.tx.Rollback()
}

func (t *SqlTransaction) EndTransaction() {
	if p := recover(); p != nil {
		t.tx.Rollback()
		panic(p)
	}
	if !t.commited {
		t.tx.Rollback()
	}
}

func (t *SqlTransaction) GetWarehouses() ([]domain.Warehouse, error) {
	rows, err := t.tx.Query(query.SelectWarehouses)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []domain.Warehouse
	for rows.Next() {
		var we domain.Warehouse
		if err := rows.Scan(&we.Name, &we.Address, &we.Capacity); err != nil {
			return nil, err
		}
		result = append(result, we)
	}
	return result, nil
}

func (t *SqlTransaction) GetWarehousesOrderedFirstWithName(warehouse string) ([]domain.Warehouse, error) {
	rows, err := t.tx.Query(query.SelectWarehousesOrderedFirstWithName, warehouse)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []domain.Warehouse
	for rows.Next() {
		var we domain.Warehouse
		if err := rows.Scan(&we.Name, &we.Address, &we.Capacity); err != nil {
			return nil, err
		}
		result = append(result, we)
	}
	return result, nil
}

func (t *SqlTransaction) InsertWarehouse(entity domain.Warehouse) error {
	_, err := t.tx.Exec(query.InsertIntoWarehouses, entity.Name, entity.Address, entity.Capacity)
	return err
}

func (t *SqlTransaction) GetProductsByWarehouse(name string) ([]domain.ProductWithQuantity, error) {
	rows, err := t.tx.Query(query.SelectProductsByWarehouse, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []domain.ProductWithQuantity
	for rows.Next() {
		product, quantity, err := t.mapCurrentRowsToProduct(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, domain.ProductWithQuantity{
			Product:  product,
			Quantity: quantity,
		})
	}
	return products, nil
}

func (t *SqlTransaction) GetUsedCapacity(warehouseName string) (int, error) {
	var usedCapacity int
	if err := t.tx.QueryRow(query.SelectUsedCapacitiyByWarehouse, warehouseName).Scan(&usedCapacity); err != nil {
		return 0, err
	}
	return usedCapacity, nil
}
func (t *SqlTransaction) InsertProduct(warehouseName string, product domain.IProduct, toInsertQuantity int) error {
	baseProduct := product.GetBaseProduct()
	if _, err := t.tx.Exec(
		query.InsertOrIgnoreIntoBrands,
		baseProduct.Brand.Name,
		baseProduct.Brand.Quality,
	); err != nil {
		return err
	}
	if _, err := t.tx.Exec(
		query.InsertOrIgnoreIntoProducts,
		baseProduct.SKU,
		baseProduct.Name,
		baseProduct.Price,
		baseProduct.Brand.Name,
		baseProduct.Type,
	); err != nil {
		return err
	}
	switch product.GetType() {
	case domain.Book:
		if _, err := t.tx.Exec(
			query.InsertOrIgnoreIntoBookProducts,
			baseProduct.SKU,
			product.(domain.BookProduct).Author,
		); err != nil {
			return err
		}
	}
	if _, err := t.tx.Exec(
		query.InsertOrUpdateIntoWarehouseProducts,
		warehouseName,
		baseProduct.SKU,
		toInsertQuantity,
		toInsertQuantity,
	); err != nil {
		return err
	}
	return nil
}

func (t *SqlTransaction) GetWarehouseProductsBySkuOrderedFirstWithName(warehouseName string, sku string) ([]domain.WarehouseProduct, error) {
	rows, err := t.tx.Query(query.SelectWarehouseProductBySkuOrderedFirstWithName, sku, warehouseName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []domain.WarehouseProduct
	for rows.Next() {
		var wpe domain.WarehouseProduct
		if err := rows.Scan(&wpe.WarehouseName, &wpe.Sku, &wpe.Quantity); err != nil {
			return nil, err
		}
		result = append(result, wpe)
	}
	return result, nil
}

func (t *SqlTransaction) RemoveProduct(warehouseName string, sku string, toRemoveQuantity int) (int, error) {
	queryResult := t.tx.QueryRow(query.SelectWarehouseProductQuantity, warehouseName, sku)
	if queryResult.Err() == sql.ErrNoRows {
		return 0, nil
	}
	originalQuantity := 0
	if err := queryResult.Scan(&originalQuantity); err != nil {
		return 0, err
	}
	updateResult := t.tx.QueryRow(query.UpdateWarehouseProductQuantity, toRemoveQuantity, toRemoveQuantity, warehouseName, sku)
	if updateResult.Err() == sql.ErrNoRows {
		return 0, nil
	}
	newQuantity := 0
	if err := updateResult.Scan(&newQuantity); err != nil {
		return 0, err
	}
	// TODO: clear row if not needed (if 0 quantity)
	removedQuantity := originalQuantity - newQuantity
	return removedQuantity, nil
}

func (t *SqlTransaction) mapCurrentRowsToProduct(rows *sql.Rows) (domain.IProduct, int, error) {
	baseProduct := domain.Product{}
	quantity := 0
	if err := rows.Scan(
		&baseProduct.SKU,
		&baseProduct.Name,
		&baseProduct.Price,
		&baseProduct.Brand.Name,
		&baseProduct.Type,
		&quantity,
	); err != nil {
		return nil, quantity, err
	}
	if err := t.tx.QueryRow(query.SelectBrandQuality, baseProduct.Brand.Name).Scan(&baseProduct.Brand.Quality); err != nil {
		return nil, 0, err
	}
	switch baseProduct.Type {
	case domain.Book:
		product := domain.BookProduct{
			Product: baseProduct,
		}
		err := t.tx.QueryRow(query.SelectFromBookProducts, product.SKU).Scan(&product.Author)
		if err != nil {
			return nil, quantity, err
		}
		return product, quantity, nil
	default:
		return nil, quantity, fmt.Errorf("unknown product type: %s", baseProduct.Type)
	}
}
