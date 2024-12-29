package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/kijevigombooc/inventory-manager/internal/inventory"
)

func main() {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	store := inventory.NewStore(db)

	service := inventory.NewService(store)

	mux := http.NewServeMux()
	handler := inventory.NewHandler(service)
	handler.RegisterRoutes(mux)
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
