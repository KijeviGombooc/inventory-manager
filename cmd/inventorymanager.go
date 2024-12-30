package main

import (
	dbsql "database/sql"
	"log"
	"net/http"

	"github.com/kijevigombooc/inventory-manager/internal/inventory"
	"github.com/kijevigombooc/inventory-manager/internal/inventory/store/sql"
)

func main() {
	db, err := dbsql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	store := sql.NewInventoryStore(db)

	service := inventory.NewService(store)

	mux := http.NewServeMux()
	handler := inventory.NewHandler(service)
	handler.RegisterRoutes(mux)
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
