package sql

import (
	"database/sql"

	"github.com/kijevigombooc/inventory-manager/internal/inventory/store"
	"github.com/kijevigombooc/inventory-manager/internal/inventory/store/sql/query"
	_ "github.com/mattn/go-sqlite3"
)

func NewInventoryStore(db *sql.DB) *inventoryStore {
	s := &inventoryStore{db: db}
	if err := s.Init(); err != nil {
		panic(err)
	}
	return s
}

type inventoryStore struct {
	db *sql.DB
}

func (s *inventoryStore) Init() error {
	// TODO: check if required
	if _, err := s.db.Exec("PRAGMA foreign_keys=ON"); err != nil {
		return err
	}
	if _, err := s.db.Exec(query.CreateWarehousesTable); err != nil {
		return err
	}
	if _, err := s.db.Exec(query.CreateBrandsTable); err != nil {
		return err
	}
	if _, err := s.db.Exec(query.CreateProductsTable); err != nil {
		return err
	}
	if _, err := s.db.Exec(query.CreateWarehouseProductsTable); err != nil {
		return err
	}
	if _, err := s.db.Exec(query.CreateBookProductsTable); err != nil {
		return err
	}
	return nil
}

func (s *inventoryStore) BeginTransaction() store.Transaction {
	tx, err := s.db.Begin()
	if err != nil {
		panic(err)
	}
	return &SqlTransaction{tx: tx, commited: false}
}
