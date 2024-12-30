package main

import (
	dbsql "database/sql"
	"log"
	"net/http"

	"github.com/kijevigombooc/inventory-manager/internal/inventory/handler/rest"
	"github.com/kijevigombooc/inventory-manager/internal/inventory/service"
	"github.com/kijevigombooc/inventory-manager/internal/inventory/store/sql"
)

func main() {
	db, err := dbsql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	store := sql.NewInventoryStore(db)

	service := service.NewInventoryService(store)

	mux := http.NewServeMux()
	handler := rest.NewInventoryHandler(service)
	handler.RegisterRoutes(mux)
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
