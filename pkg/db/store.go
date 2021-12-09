package db

var valueStore []uint64

// Store is a mock for storing the block number values in a database
type Store struct {
}

// NewStore is a constructor of Store
func NewStore() *Store {
	return &Store{}
}

// Save stores the block number value in a local in-memory slice
func (s *Store) Save(value uint64) error {
	valueStore = append(valueStore, value)
	return nil
}
