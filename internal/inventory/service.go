package inventory

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
		products, err := trx.GetProductsByWarehouse(warehouse.Name)
		if err != nil {
			return nil, err
		}
		result = append(result, WarehouseDetailDto{
			WarehouseDto: WarehouseEntityToDto(warehouse),
			Products:     products,
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

func (s *service) InsertProducts(products []interface{}, quantity int) error {
	trx := s.store.BeginTransaction()
	defer trx.EndTransaction()
	// TODO: finish this
	return nil
}

func (s *service) RemoveProducts(sku string, quantity int) error {
	trx := s.store.BeginTransaction()
	defer trx.EndTransaction()
	// TODO: finish this
	return nil
}
