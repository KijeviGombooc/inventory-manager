package query

const CreateWarehousesTable = `
	CREATE TABLE IF NOT EXISTS warehouses (
		name TEXT PRIMARY KEY,
		address TEXT NOT NULL,
		capacity INTEGER NOT NULL
	)
`
const CreateProductsTable = `
	CREATE TABLE IF NOT EXISTS products (
		sku TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		price INTEGER NOT NULL,
		type TEXT NOT NULL
	)
`
const CreateWarehouseProductsTable = `
	CREATE TABLE IF NOT EXISTS warehouse_products (
		warehouse_name TEXT NOT NULL,
		sku TEXT NOT NULL,
		quantity INTEGER NOT NULL,
		FOREIGN KEY (warehouse_name) REFERENCES warehouses (name),
		FOREIGN KEY (sku) REFERENCES products (sku),
		PRIMARY KEY (warehouse_name, sku)
	)
`
const CreateBookProductsTable = `
	CREATE TABLE IF NOT EXISTS book_products (
		sku TEXT PRIMARY KEY,
		author TEXT NOT NULL,
		FOREIGN KEY (sku) REFERENCES products (sku) ON DELETE CASCADE
	)
`
const SelectWarehouses = "SELECT name, address, capacity FROM warehouses"
const InsertIntoWarehouses = "INSERT INTO warehouses (name, address, capacity) VALUES (?, ?, ?)"
const SelectFromBookProducts = "SELECT author FROM book_products WHERE sku = ?"

const SelectWarehousesOrderedFirstWithName = `
			SELECT name, address, capacity
			FROM warehouses
			ORDER BY CASE WHEN name = ? THEN 0 ELSE 1 END, name
		`
const SelectProductsByWarehouse = `
		SELECT p.sku, p.name, p.price, p.type, wp.quantity
		FROM products p
		JOIN warehouse_products wp ON p.sku = wp.sku
		WHERE wp.warehouse_name = ? AND wp.quantity > 0
	`
const SelectUsedCapacitiyByWarehouse = `
	SELECT IFNULL(SUM(quantity), 0)
	FROM warehouse_products
	WHERE warehouse_name = ?
`
const SelectWarehouseProductBySkuOrderedFirstWithName = `
			SELECT wp.warehouse_name, wp.sku, wp.quantity
			FROM warehouse_products wp
			WHERE wp.sku = ?
			ORDER BY CASE WHEN wp.warehouse_name = ? THEN 0 ELSE 1 END, wp.warehouse_name
		`
const SelectWarehouseProductQuantity = `
		SELECT quantity FROM warehouse_products
		WHERE warehouse_name = ? AND sku = ?
	`
const InsertOrIgnoreIntoProducts = `
	INSERT OR IGNORE INTO products (sku, name, price, type)
	VALUES (?, ?, ?, ?)
`
const InsertOrIgnoreIntoBookProducts = `
					INSERT OR IGNORE INTO book_products (sku, author)
					VALUES (?, ?)
				`
const InsertOrUpdateIntoWarehouseProducts = `
				INSERT INTO warehouse_products (warehouse_name, sku, quantity)
				VALUES (?, ?, ?)
				ON CONFLICT (warehouse_name, sku)
				DO UPDATE SET quantity = quantity + ?
			`

const UpdateWarehouseProductQuantity = `
	UPDATE warehouse_products
	SET quantity = CASE
		WHEN quantity - ? < 0 THEN 0
		ELSE quantity - ?
	END
	WHERE warehouse_name = ? AND sku = ?
	RETURNING quantity AS new_quantity
`