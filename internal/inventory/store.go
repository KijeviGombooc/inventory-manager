package inventory

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func NewStore(db *sql.DB) *store {
	s := &store{db: db}
	if err := s.init(); err != nil {
		panic(err)
	}
	return s
}

type store struct {
	db *sql.DB
}

func (s *store) init() error {
	if _, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS warehouses (
			name TEXT PRIMARY KEY,
			address TEXT NOT NULL,
			capacity INTEGER NOT NULL
		)
	`); err != nil {
		return err
	}
	if _, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS products (
			sku TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			price INTEGER NOT NULL,
			type TEXT NOT NULL
		)
	`); err != nil {
		return err
	}
	if _, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS warehouse_products (
			warehouse_name TEXT NOT NULL,
			sku TEXT NOT NULL,
			quantity INTEGER NOT NULL,
			FOREIGN KEY (warehouse_name) REFERENCES warehouses (name),
			FOREIGN KEY (sku) REFERENCES products (sku),
			PRIMARY KEY (warehouse_name, sku)
		)
	`); err != nil {
		return err
	}
	if _, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS book_products (
			sku TEXT PRIMARY KEY,
			author TEXT NOT NULL,
			FOREIGN KEY (sku) REFERENCES products (sku)
		)
	`); err != nil {
		return err
	}
	return nil
}

func (s *store) BeginTransaction() *Transaction {
	tx, err := s.db.Begin()
	if err != nil {
		panic(err)
	}
	return &Transaction{tx: tx, commited: false}
}

type Transaction struct {
	tx       *sql.Tx
	commited bool
}

func (t *Transaction) CommitTransaction() error {
	t.commited = true
	return t.tx.Commit()
}

func (t *Transaction) RollbackTransaction() error {
	return t.tx.Rollback()
}

func (t *Transaction) EndTransaction() {
	if p := recover(); p != nil {
		t.tx.Rollback()
		panic(p)
	}
	if !t.commited {
		t.tx.Rollback()
	}
}

func (t *Transaction) GetWarehouses() ([]WarehouseEntity, error) {
	rows, err := t.tx.Query("SELECT name, address, capacity FROM warehouses")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []WarehouseEntity
	for rows.Next() {
		var we WarehouseEntity
		if err := rows.Scan(&we.Name, &we.Address, &we.Capacity); err != nil {
			return nil, err
		}
		result = append(result, we)
	}
	return result, nil
}

func (t *Transaction) GetWarehousesOrderedFirstWithName(warehouse string) ([]WarehouseEntity, error) {
	rows, err := t.tx.Query(`
		SELECT name, address, capacity
		FROM warehouses
		ORDER BY CASE WHEN name = ? THEN 0 ELSE 1 END, name
	`, warehouse)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []WarehouseEntity
	for rows.Next() {
		var we WarehouseEntity
		if err := rows.Scan(&we.Name, &we.Address, &we.Capacity); err != nil {
			return nil, err
		}
		result = append(result, we)
	}
	return result, nil
}

func (t *Transaction) InsertWarehouse(entity WarehouseEntity) error {
	_, err := t.tx.Exec("INSERT INTO warehouses (name, address, capacity) VALUES (?, ?, ?)", entity.Name, entity.Address, entity.Capacity)
	return err
}

func (t *Transaction) GetProductsByWarehouse(name string) ([]ProductEntityWithQuantity, error) {
	rows, err := t.tx.Query(`
		SELECT p.sku, p.name, p.price, p.type, wp.quantity
		FROM products p
		JOIN warehouse_products wp ON p.sku = wp.sku
		WHERE wp.warehouse_name = ?
	`, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []ProductEntityWithQuantity
	for rows.Next() {
		product, quantity, err := t.mapCurrentRowsToProduct(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, ProductEntityWithQuantity{
			Product:  product,
			Quantity: quantity,
		})
	}

	return products, nil
}

func (t *Transaction) GetUsedCapacity(name string) (int, error) {
	var usedCapacity int
	if err := t.tx.QueryRow(`
		SELECT IFNULL(SUM(quantity), 0)
		FROM warehouse_products
		WHERE warehouse_name = ?
	`, name).Scan(&usedCapacity); err != nil {
		return 0, err
	}
	return usedCapacity, nil
}

func (t *Transaction) InsertProduct(warehouseName string, product IProduct, toInsertQuantity int) error {
	_, err := t.tx.Exec(`
		INSERT OR IGNORE INTO products (sku, name, price, type)
		VALUES (?, ?, ?, ?)
	`, product.GetSKU(), product.GetName(), product.GetPrice(), product.GetType())
	switch product.GetType() {
	case Book:
		_, err = t.tx.Exec(`
			INSERT OR IGNORE INTO book_products (sku, author)
			VALUES (?, ?)
		`, product.GetSKU(), product.(BookProductEntity).Author)
	}
	if err != nil {
		return err
	}
	_, err = t.tx.Exec(`
		INSERT INTO warehouse_products (warehouse_name, sku, quantity)
		VALUES (?, ?, ?)
		ON CONFLICT (warehouse_name, sku)
		DO UPDATE SET quantity = quantity + ?
	`, warehouseName, product.GetSKU(), toInsertQuantity, toInsertQuantity)
	return err
}

func (t *Transaction) mapCurrentRowsToProduct(rows *sql.Rows) (IProduct, int, error) {
	baseProduct := ProductEntity{}
	quantity := 0
	if err := rows.Scan(
		&baseProduct.SKU,
		&baseProduct.Name,
		&baseProduct.Price,
		&baseProduct.Type,
		&quantity,
	); err != nil {
		return nil, quantity, err
	}
	switch baseProduct.Type {
	case Book:
		product := BookProductEntity{
			ProductEntity: baseProduct,
		}
		err := t.tx.QueryRow("SELECT author FROM book_products WHERE sku = ?", product.SKU).Scan(&product.Author)
		if err != nil {
			return nil, quantity, err
		}
		return product, quantity, nil
	default:
		return nil, quantity, fmt.Errorf("unknown product type: %s", baseProduct.Type)
	}
}
