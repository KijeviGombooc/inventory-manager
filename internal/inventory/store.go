package inventory

import (
	"database/sql"

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

func (s *store) BeginTransaction() *Transaction {
	tx, err := s.db.Begin()
	if err != nil {
		panic(err)
	}
	return &Transaction{tx: tx, commited: false}
}

func (s *store) init() error {
	if _, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS warehouses (
			-- id INTEGER PRIMARY KEY AUTOINCREMENT,
			-- name TEXT NOT NULL,
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
			volume INTEGER NOT NULL,
			type TEXT NOT NULL
		)
	`); err != nil {
		return err
	}
	if _, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS warehouse_products (
			warehouse_id INTEGER NOT NULL,
			sku TEXT NOT NULL,
			quantity INTEGER NOT NULL,
			FOREIGN KEY (warehouse_id) REFERENCES warehouses (id),
			FOREIGN KEY (sku) REFERENCES products (sku),
			PRIMARY KEY (warehouse_id, sku)
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

type Transaction struct {
	tx       *sql.Tx
	commited bool
}

func (t *Transaction) GetProductsByWarehouse(name string) ([]any, error) {
	return nil, nil
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

func (t *Transaction) InsertWarehouse(entity WarehouseEntity) error {
	_, err := t.tx.Exec("INSERT INTO warehouses (name, address, capacity) VALUES (?, ?, ?)", entity.Name, entity.Address, entity.Capacity)
	return err
}
