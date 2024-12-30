package store

type Store interface {
	Init() error
	BeginTransaction() Transaction
}
