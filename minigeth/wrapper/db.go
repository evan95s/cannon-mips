package wrapper

import (
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/ethdb/memorydb"
)

type DB struct {
	inner ethdb.Database
}

func NewDB(memoryDB *memorydb.Database) ethdb.Database {
	return &DB{
		inner: rawdb.NewDatabase(memoryDB),
	}
}

func (db *DB) Ancient(kind string, number uint64) ([]byte, error) {
	return db.inner.Ancient(kind, number)
}

func (db *DB) AncientDatadir() (string, error) {
	return db.inner.AncientDatadir()
}

func (db *DB) AncientRange(kind string, start, count, maxBytes uint64) ([][]byte, error) {
	return db.inner.AncientRange(kind, start, count, maxBytes)
}

func (db *DB) AncientSize(kind string) (uint64, error) {
	return db.inner.AncientSize(kind)
}

func (db *DB) Ancients() (uint64, error) {
	return db.inner.Ancients()
}

func (db *DB) Close() error {
	return db.inner.Close()
}

func (db *DB) Compact(start []byte, limit []byte) error {
	return db.inner.Compact(start, limit)
}

func (db *DB) Delete(key []byte) error {
	return db.inner.Delete(key)
}

func (db *DB) Get(key []byte) ([]byte, error) {
	return db.inner.Get(key)
}

func (db *DB) Has(key []byte) (bool, error) {
	return db.inner.Has(key)
}

func (db *DB) HasAncient(kind string, number uint64) (bool, error) {
	return db.inner.HasAncient(kind, number)
}

func (db *DB) MigrateTable(s string, f func([]byte) ([]byte, error)) error {
	return db.inner.MigrateTable(s, f)
}

func (db *DB) ModifyAncients(f func(ethdb.AncientWriteOp) error) (int64, error) {
	return db.inner.ModifyAncients(f)
}

func (db *DB) NewBatch() ethdb.Batch {
	return db.inner.NewBatch()
}

func (db *DB) NewBatchWithSize(size int) ethdb.Batch {
	return db.inner.NewBatchWithSize(size)
}

func (db *DB) NewIterator(prefix []byte, start []byte) ethdb.Iterator {
	return db.inner.NewIterator(prefix, start)
}

func (db *DB) NewSnapshot() (ethdb.Snapshot, error) {
	return db.inner.NewSnapshot()
}

func (db *DB) Put(key []byte, value []byte) error {
	return db.inner.Put(key, value)
}

func (db *DB) ReadAncients(f func(ethdb.AncientReaderOp) error) (err error) {
	return db.inner.ReadAncients(f)
}

func (db *DB) Stat(property string) (string, error) {
	return db.inner.Stat(property)
}

func (db *DB) Sync() error {
	return db.inner.Sync()
}

func (db *DB) Tail() (uint64, error) {
	return db.inner.Tail()
}

func (db *DB) TruncateHead(n uint64) (uint64, error) {
	return db.inner.TruncateHead(n)
}

func (db *DB) TruncateTail(n uint64) (uint64, error) {
	return db.inner.TruncateTail(n)
}
