package service

import (
	dbsql "database/sql"
	"fmt"
	"reflect"
	"testing"

	"github.com/kijevigombooc/inventory-manager/internal/inventory/handler/dto"
	"github.com/kijevigombooc/inventory-manager/internal/inventory/store/sql"
)

var s Service
var db *dbsql.DB
var warehouses []dto.Warehouse
var bookProducts []dto.BookProduct
var consumableProducts []dto.ConsumableProduct
var electronicsProducts []dto.ElectronicsProduct

func TestMain(m *testing.M) {
	BeforeAll()
	m.Run()
}

func BeforeAll() {
	for i := 0; i <= 10; i++ {
		warehouses = append(warehouses, dto.Warehouse{
			Name:     fmt.Sprintf("Warehouse %d", i+1),
			Address:  fmt.Sprintf("Address %d", i+1),
			Capacity: i,
		})
	}
	for i := 0; i < 10; i++ {
		bookProducts = append(bookProducts, dto.BookProduct{
			Product: dto.Product{
				Name:  fmt.Sprintf("Book %c", 'A'+i),
				Price: 100,
				SKU:   fmt.Sprintf("BOOK-%c", 'A'+i),
				Brand: dto.Brand{
					Name:    "Book Brand",
					Quality: 4,
				},
				Type: dto.Book,
			},
			Author: "Author",
		})
	}
	for i := 0; i < 10; i++ {
		consumableProducts = append(consumableProducts, dto.ConsumableProduct{
			Product: dto.Product{
				Name:  fmt.Sprintf("Consumable %c", 'A'+i),
				Price: 100,
				SKU:   fmt.Sprintf("CONS-%c", 'A'+i),
				Brand: dto.Brand{
					Name:    "Consumable Brand",
					Quality: 4,
				},
				Type: dto.Consumable,
			},
			ExpirationDate: "2024.12.12",
		})
	}
	for i := 0; i < 10; i++ {
		electronicsProducts = append(electronicsProducts, dto.ElectronicsProduct{
			Product: dto.Product{
				Name:  fmt.Sprintf("Electronics %c", 'A'+i),
				Price: 100,
				SKU:   fmt.Sprintf("ETRX-%c", 'A'+i),
				Brand: dto.Brand{
					Name:    "Electronics Brand",
					Quality: 4,
				},
				Type: dto.Electronics,
			},
			WarrantyPeriod: "2 Years",
		})
	}
}

func BeforeEach() {
	db, err := dbsql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	store := sql.NewInventoryStore(db)
	s = NewInventoryService(store)
}

func AfterEach() {
	if db != nil {
		db.Close()
	}
	s = nil
}

func TestCreateWarehouseSuccessful(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	warehouseCapacity := 3
	if err := s.CreateWarehouse(warehouses[warehouseCapacity]); err != nil {
		t.Fatalf("Error creating warehouse: %v", err)
	}
}

func TestCreateWarehouseErrorAlreadyExists(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	warehouseCapacity := 3
	if err := s.CreateWarehouse(warehouses[warehouseCapacity]); err != nil {
		t.Fatalf("Error creating warehouse: %v", err)
	}
	if err := s.CreateWarehouse(warehouses[warehouseCapacity]); err == nil {
		t.Fatalf("Should have failed to create warehouse")
	}
}

func TestListWarehousesEmpty(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	warehouses, err := s.GetWarehouses()
	if err != nil {
		t.Fatalf("Error listing warehouses: %v", err)
	}
	if len(warehouses) != 0 {
		t.Fatalf("Warehouses should be empty")
	}
}

func TestListWarehousesNotEmpty(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	warehouseCapacity := 3
	if err := s.CreateWarehouse(warehouses[warehouseCapacity]); err != nil {
		t.Fatalf("Error creating warehouse: %v", err)
	}
	warehouses, err := s.GetWarehouses()
	if err != nil {
		t.Fatalf("Error listing warehouses: %v", err)
	}
	if len(warehouses) != 1 {
		t.Fatalf("Warehouses should not be empty")
	}
}

func TestListWarehousesWWarehouseNotEmpty(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	warehouseCapacity := 3
	toInsertQuantity := 3
	toInsertProduct := bookProducts[0]
	if err := s.CreateWarehouse(warehouses[warehouseCapacity]); err != nil {
		t.Fatalf("Error creating warehouse: %v", err)
	}
	if err := s.InsertProducts(warehouses[warehouseCapacity].Name, &toInsertProduct, toInsertQuantity); err != nil {
		t.Fatalf("Error inserting product: %v", err)
	}
	warehouses, err := s.GetWarehouses()
	if err != nil {
		t.Fatalf("Error listing warehouses: %v", err)
	}
	if len(warehouses) != 1 {
		t.Fatalf("Warehouses should not be empty")
	}
	if len(warehouses[0].Products) != 1 {
		t.Fatalf("Warehouse should not be empty")
	}
	if warehouses[0].Products[0].Quantity != toInsertQuantity {
		t.Fatalf("Product quantity should be %d", toInsertQuantity)
	}
}

func TestInsertAndListSuccessfulByContent(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	warehouseCapacity := 3
	toInsertQuantity := 3
	toInsertProduct := bookProducts[0]
	if err := s.CreateWarehouse(warehouses[warehouseCapacity]); err != nil {
		t.Fatalf("Error creating warehouse: %v", err)
	}
	if err := s.InsertProducts(warehouses[warehouseCapacity].Name, &toInsertProduct, toInsertQuantity); err != nil {
		t.Fatalf("Error inserting product: %v", err)
	}
	warehouses, err := s.GetWarehouses()
	if err != nil {
		t.Fatalf("Error listing warehouses: %v", err)
	}
	if !reflect.DeepEqual(warehouses[0].Products[0].IProduct, &toInsertProduct) {
		t.Fatalf("Products should be the same: %v, %v", warehouses[0].Products[0].IProduct, toInsertProduct)
	}
}

func TestInsertLocalBookSuccessful(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	warehouseCapacity := 3
	toInsertQuantity := 3
	if err := s.CreateWarehouse(warehouses[warehouseCapacity]); err != nil {
		t.Fatalf("Error creating warehouse: %v", err)
	}
	if err := s.InsertProducts(warehouses[warehouseCapacity].Name, &bookProducts[0], toInsertQuantity); err != nil {
		t.Fatalf("Error inserting product: %v", err)
	}
}

func TestInsertLocalConsumableSuccessful(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	warehouseCapacity := 3
	toInsertQuantity := 3
	if err := s.CreateWarehouse(warehouses[warehouseCapacity]); err != nil {
		t.Fatalf("Error creating warehouse: %v", err)
	}
	if err := s.InsertProducts(warehouses[warehouseCapacity].Name, &consumableProducts[0], toInsertQuantity); err != nil {
		t.Fatalf("Error inserting product: %v", err)
	}
}

func TestInsertLocalElectronicsSuccessful(t *testing.T) {
	BeforeEach()
	defer AfterEach()

	warehouseCapacity := 3
	toInsertQuantity := 3
	if err := s.CreateWarehouse(warehouses[warehouseCapacity]); err != nil {
		t.Fatalf("Error creating warehouse: %v", err)
	}
	if err := s.InsertProducts(warehouses[warehouseCapacity].Name, &electronicsProducts[0], toInsertQuantity); err != nil {
		t.Fatalf("Error inserting product: %v", err)
	}
}

func TestInsertLocalBookErrorNotEnoughSpace(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	warehouseCapacity := 2
	toInsertQuantity := 3
	if err := s.CreateWarehouse(warehouses[warehouseCapacity]); err != nil {
		t.Fatalf("Error creating warehouse: %v", err)
	}
	if err := s.InsertProducts(warehouses[warehouseCapacity].Name, &bookProducts[0], toInsertQuantity); err == nil {
		t.Fatalf("Should have failed to insert product")
	}
}

func TestInsertLocalConsumableErrorNotEnoughSpace(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	warehouseCapacity := 2
	toInsertQuantity := 3
	if err := s.CreateWarehouse(warehouses[warehouseCapacity]); err != nil {
		t.Fatalf("Error creating warehouse: %v", err)
	}
	if err := s.InsertProducts(warehouses[warehouseCapacity].Name, &consumableProducts[0], toInsertQuantity); err == nil {
		t.Fatalf("Should have failed to insert product")
	}
}

func TestInsertLocalElectronicsErrorNotEnoughSpace(t *testing.T) {
	BeforeEach()
	defer AfterEach()

	warehouseCapacity := 2
	toInsertQuantity := 3
	if err := s.CreateWarehouse(warehouses[warehouseCapacity]); err != nil {
		t.Fatalf("Error creating warehouse: %v", err)
	}
	if err := s.InsertProducts(warehouses[warehouseCapacity].Name, &electronicsProducts[0], toInsertQuantity); err == nil {
		t.Fatalf("Should have failed to insert product")
	}
}

func TestInsertLocalMultipleTypesSuccessful(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	warehouseCapacity := 10
	toInsertQuantity := 3
	if err := s.CreateWarehouse(warehouses[warehouseCapacity]); err != nil {
		t.Fatalf("Error creating warehouse: %v", err)
	}
	if err := s.InsertProducts(warehouses[warehouseCapacity].Name, &bookProducts[0], toInsertQuantity); err != nil {
		t.Fatalf("Error inserting product: %v", err)
	}
	if err := s.InsertProducts(warehouses[warehouseCapacity].Name, &consumableProducts[0], toInsertQuantity); err != nil {
		t.Fatalf("Error inserting product: %v", err)
	}
	if err := s.InsertProducts(warehouses[warehouseCapacity].Name, &electronicsProducts[0], toInsertQuantity); err != nil {
		t.Fatalf("Error inserting product: %v", err)
	}
}

func TestInsertLocalMultipleTypesErrorNotEnoughSpace(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	warehouseCapacity := 10
	toInsertQuantity := 4
	if err := s.CreateWarehouse(warehouses[warehouseCapacity]); err != nil {
		t.Fatalf("Error creating warehouse: %v", err)
	}
	if err := s.InsertProducts(warehouses[warehouseCapacity].Name, &bookProducts[0], toInsertQuantity); err != nil {
		t.Fatalf("Error inserting product: %v", err)
	}
	if err := s.InsertProducts(warehouses[warehouseCapacity].Name, &consumableProducts[0], toInsertQuantity); err != nil {
		t.Fatalf("Error inserting product: %v", err)
	}
	if err := s.InsertProducts(warehouses[warehouseCapacity].Name, &electronicsProducts[0], toInsertQuantity); err == nil {
		t.Fatalf("Should have failed to insert product")
	}
}

func TestInsertGlobalBookSuccessful(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	toInsertQuantity := 4
	warehouse1Capacity := 2
	warehouse2Capacity := 3
	if err := s.CreateWarehouse(warehouses[warehouse1Capacity]); err != nil {
		t.Fatalf("Error creating warehouse: %v", err)
	}
	if err := s.CreateWarehouse(warehouses[warehouse2Capacity]); err != nil {
		t.Fatalf("Error creating warehouse: %v", err)
	}
	if err := s.InsertProducts(warehouses[warehouse1Capacity].Name, &bookProducts[0], toInsertQuantity); err != nil {
		t.Fatalf("Error inserting product: %v", err)
	}
}

func TestInsertGlobalConsumableSuccessful(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	toInsertQuantity := 4
	warehouse1Capacity := 2
	warehouse2Capacity := 3
	if err := s.CreateWarehouse(warehouses[warehouse1Capacity]); err != nil {
		t.Fatalf("Error creating warehouse: %v", err)
	}
	if err := s.CreateWarehouse(warehouses[warehouse2Capacity]); err != nil {
		t.Fatalf("Error creating warehouse: %v", err)
	}
	if err := s.InsertProducts(warehouses[warehouse1Capacity].Name, &consumableProducts[0], toInsertQuantity); err != nil {
		t.Fatalf("Error inserting product: %v", err)
	}
}

func TestInsertGlobalElectronicsSuccessful(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	toInsertQuantity := 4
	warehouse1Capacity := 2
	warehouse2Capacity := 3
	if err := s.CreateWarehouse(warehouses[warehouse1Capacity]); err != nil {
		t.Fatalf("Error creating warehouse: %v", err)
	}
	if err := s.CreateWarehouse(warehouses[warehouse2Capacity]); err != nil {
		t.Fatalf("Error creating warehouse: %v", err)
	}
	if err := s.InsertProducts(warehouses[warehouse1Capacity].Name, &electronicsProducts[0], toInsertQuantity); err != nil {
		t.Fatalf("Error inserting product: %v", err)
	}
}

func TestInsertGlobalMultipleTypesSuccessful(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	toInsertQuantity := 3
	warehouse1Capacity := 4
	warehouse2Capacity := 5
	if err := s.CreateWarehouse(warehouses[warehouse1Capacity]); err != nil {
		t.Fatalf("Error creating warehouse: %v", err)
	}
	if err := s.CreateWarehouse(warehouses[warehouse2Capacity]); err != nil {
		t.Fatalf("Error creating warehouse: %v", err)
	}
	if err := s.InsertProducts(warehouses[warehouse1Capacity].Name, &bookProducts[0], toInsertQuantity); err != nil {
		t.Fatalf("Error inserting product: %v", err)
	}
	if err := s.InsertProducts(warehouses[warehouse1Capacity].Name, &consumableProducts[0], toInsertQuantity); err != nil {
		t.Fatalf("Error inserting product: %v", err)
	}
	if err := s.InsertProducts(warehouses[warehouse1Capacity].Name, &electronicsProducts[0], toInsertQuantity); err != nil {
		t.Fatalf("Error inserting product: %v", err)
	}
}

func TestInsertGlobalMultipleTypesErrorNotEnoughSpace(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	toInsertQuantity := 4
	warehouse1Capacity := 4
	warehouse2Capacity := 5
	if err := s.CreateWarehouse(warehouses[warehouse1Capacity]); err != nil {
		t.Fatalf("Error creating warehouse: %v", err)
	}
	if err := s.CreateWarehouse(warehouses[warehouse2Capacity]); err != nil {
		t.Fatalf("Error creating warehouse: %v", err)
	}
	if err := s.InsertProducts(warehouses[warehouse1Capacity].Name, &bookProducts[0], toInsertQuantity); err != nil {
		t.Fatalf("Error inserting product: %v", err)
	}
	if err := s.InsertProducts(warehouses[warehouse1Capacity].Name, &consumableProducts[0], toInsertQuantity); err != nil {
		t.Fatalf("Error inserting product: %v", err)
	}
	if err := s.InsertProducts(warehouses[warehouse1Capacity].Name, &electronicsProducts[0], toInsertQuantity); err == nil {
		t.Fatalf("Should have failed to insert product")
	}
}

func TestRemoveGlobalBookSuccessful(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	toInsert1Quantity := 2
	toInsert2Quantity := 4
	toRemoveQuantity := 6
	warehouse1Capacity := 4
	warehouse2Capacity := 5
	if err := s.CreateWarehouse(warehouses[warehouse1Capacity]); err != nil {
		t.Fatalf("Error creating warehouse: %v", err)
	}
	if err := s.CreateWarehouse(warehouses[warehouse2Capacity]); err != nil {
		t.Fatalf("Error creating warehouse: %v", err)
	}
	if err := s.InsertProducts(warehouses[warehouse1Capacity].Name, &bookProducts[0], toInsert1Quantity); err != nil {
		t.Fatalf("Error inserting product: %v", err)
	}
	if err := s.InsertProducts(warehouses[warehouse2Capacity].Name, &bookProducts[0], toInsert2Quantity); err != nil {
		t.Fatalf("Error inserting product: %v", err)
	}
	if err := s.RemoveProducts(warehouses[warehouse1Capacity].Name, bookProducts[0].SKU, toRemoveQuantity); err != nil {
		t.Fatalf("Error removing product: %v", err)
	}
}

func TestRemoveGlobalBookErrorNotEnoughQuantity(t *testing.T) {
	BeforeEach()
	defer AfterEach()
	toInsert1Quantity := 2
	toInsert2Quantity := 4
	toRemoveQuantity := 7
	warehouse1Capacity := 4
	warehouse2Capacity := 5
	if err := s.CreateWarehouse(warehouses[warehouse1Capacity]); err != nil {
		t.Fatalf("Error creating warehouse: %v", err)
	}
	if err := s.CreateWarehouse(warehouses[warehouse2Capacity]); err != nil {
		t.Fatalf("Error creating warehouse: %v", err)
	}
	if err := s.InsertProducts(warehouses[warehouse1Capacity].Name, &bookProducts[0], toInsert1Quantity); err != nil {
		t.Fatalf("Error inserting product: %v", err)
	}
	if err := s.InsertProducts(warehouses[warehouse2Capacity].Name, &bookProducts[0], toInsert2Quantity); err != nil {
		t.Fatalf("Error inserting product: %v", err)
	}
	if err := s.RemoveProducts(warehouses[warehouse1Capacity].Name, bookProducts[0].SKU, toRemoveQuantity); err == nil {
		t.Fatalf("Should have failed to remove product")
	}
}
