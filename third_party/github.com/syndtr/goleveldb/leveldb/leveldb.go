// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package leveldb

import (
	"context"
	"fmt"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
	errors_ "github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"

	"github.com/searKing/golang/go/container/hashring"
	"github.com/searKing/golang/go/errors"
	sync_ "github.com/searKing/golang/go/sync"
)

type ConsistentDB struct {
	mu sync.Mutex

	dbByPath map[string]*leveldb.DB
	pool     *hashring.StringNodeLocator
	options  *opt.Options

	PathPrefix string
	PoolSize   int

	subject sync_.Subject
}

func NewConsistentDB(prefix string, poolSize int, o *opt.Options) (*ConsistentDB, error) {
	db := &ConsistentDB{
		dbByPath:   make(map[string]*leveldb.DB),
		pool:       hashring.NewStringNodeLocator(),
		options:    o,
		PathPrefix: prefix,
		PoolSize:   poolSize,
	}
	if err := db.Init(db.PathPrefix, db.PoolSize, o); err != nil {
		return nil, err
	}
	return db, nil
}

func (cdb *ConsistentDB) Close() error {
	cdb.mu.Lock()
	defer cdb.mu.Unlock()
	cdb.subject.PublishBroadcast(context.Background(), fmt.Errorf("leveldb closed"))
	var errs []error
	for _, _db := range cdb.dbByPath {
		if _db == nil {
			continue
		}
		if err := _db.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	cdb.dbByPath = make(map[string]*leveldb.DB)
	cdb.pool.RemoveAllNodes()
	return errors.Multi(errs...)
}

func (cdb *ConsistentDB) Init(pathPrefix string, poolSize int, o *opt.Options) (err error) {
	if err := cdb.Close(); err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = cdb.Close()
		}
	}()

	cdb.mu.Lock()
	defer cdb.mu.Unlock()
	if poolSize < 1 {
		poolSize = 1
	}
	cdb.PathPrefix = pathPrefix
	cdb.PoolSize = poolSize
	cdb.options = o
	var paths []string
	for i := 0; i < poolSize; i++ {
		path := fmt.Sprintf("%s-%d", pathPrefix, i)
		// Open the db and recover any potential corruptions
		_db, _err := leveldb.OpenFile(path, o)
		if _, corrupted := _err.(*errors_.ErrCorrupted); corrupted {
			_db, _err = leveldb.RecoverFile(path, nil)
		}
		// (Re)check for errors and abort if opening of the db failed
		if _err != nil {
			err = _err
			return err
		}
		if db, has := cdb.dbByPath[path]; has {
			_ = db.Close()
		}

		cdb.dbByPath[path] = _db
		paths = append(paths, path)
	}
	cdb.pool.AddNodes(paths...)
	return nil
}

func (cdb *ConsistentDB) LevelDB(router string) (path string, db *leveldb.DB) {
	if cdb == nil {
		return "", nil
	}
	cdb.mu.Lock()
	defer cdb.mu.Unlock()
	_dbPath, has := cdb.pool.Get(router)
	if !has {
		return _dbPath, nil
	}

	_db, has := cdb.dbByPath[_dbPath]
	if !has {
		return _dbPath, nil
	}
	return _dbPath, _db
}

func (cdb *ConsistentDB) AllLevelDBNodeByPath() map[string]*leveldb.DB {
	cdb.mu.Lock()
	defer cdb.mu.Unlock()
	var dbBypath = map[string]*leveldb.DB{}
	for k, v := range cdb.dbByPath {
		dbBypath[k] = v
	}
	return dbBypath
}

// This is for convenience

// Subscribe returns a channel that's closed when awoken by PublishSignal or PublishBroadcast in convenience function below.
func (cdb *ConsistentDB) Subscribe() (<-chan any, context.CancelFunc) {
	return cdb.subject.Subscribe()
}

// Stats populates s with database statistics.
func (cdb *ConsistentDB) Stats(router string, s *leveldb.DBStats) error {
	path, db := cdb.LevelDB(router)
	if db == nil {
		return fmt.Errorf("leveldb not found, %s", router)
	}
	cdb.subject.PublishBroadcast(context.Background(),
		fmt.Sprintf("leveldb[%s] Select by %s Stats", path, router))
	return db.Stats(s)
}

// Write apply the given batch to the DB. The batch records will be applied
// sequentially. Write might be used concurrently, when used concurrently and
// batch is small enough, write will try to merge the batches. Set NoWriteMerge
// option to true to disable write merge.
//
// It is safe to modify the contents of the arguments after Write returns but
// not before. Write will not modify content of the batch.
func (cdb *ConsistentDB) Write(router string, batch *leveldb.Batch, wo *opt.WriteOptions) error {
	path, db := cdb.LevelDB(router)
	if db == nil {
		return fmt.Errorf("leveldb not found, %s", router)
	}
	cdb.subject.PublishBroadcast(context.Background(),
		fmt.Sprintf("leveldb[%s] Select by %s Write batch", path, router))
	return db.Write(batch, wo)
}

// Put sets the value for the given key. It overwrites any previous value
// for that key; a DB is not a multi-map. Write merge also applies for Put, see
// Write.
//
// It is safe to modify the contents of the arguments after Put returns but not
// before.
func (cdb *ConsistentDB) Put(router string, key, value []byte, wo *opt.WriteOptions) error {
	path, db := cdb.LevelDB(router)
	if db == nil {
		return fmt.Errorf("leveldb not found, %s", router)
	}
	cdb.subject.PublishBroadcast(context.Background(),
		fmt.Sprintf("leveldb[%s] Select by %s Put %s with %d bytes", path, router, key, len(value)))
	return db.Put(key, value, wo)
}

// Delete deletes the value for the given key. Delete will not returns error if
// key doesn't exist. Write merge also applies for Delete, see Write.
//
// It is safe to modify the contents of the arguments after Delete returns but
// not before.
func (cdb *ConsistentDB) Delete(router string, key []byte, wo *opt.WriteOptions) error {
	path, db := cdb.LevelDB(router)
	if db == nil {
		return nil
	}
	cdb.subject.PublishBroadcast(context.Background(),
		fmt.Sprintf("leveldb[%s] Select by %s Delete %s", path, router, key))
	return db.Delete(key, wo)
}

// Has returns true if the DB does contains the given key.
//
// It is safe to modify the contents of the argument after Has returns.
func (cdb *ConsistentDB) Has(router string, key []byte, ro *opt.ReadOptions) (ret bool, err error) {
	_, db := cdb.LevelDB(router)
	if db == nil {
		return false, nil
	}
	return db.Has(key, ro)
}

// Front returns the first element of list l or nil if the list is empty.
func (cdb *ConsistentDB) Front(path string, ro *opt.ReadOptions) (key, value []byte, err error) {
	cdb.mu.Lock()
	_db := cdb.dbByPath[path]
	cdb.mu.Unlock()
	if _db == nil {
		return nil, nil, fmt.Errorf("leveldb not found, %s", path)
	}
	it := _db.NewIterator(nil, ro)
	defer it.Release()
	if it.First() {
		k := it.Key()
		v := it.Value()
		return k, v, nil
	}
	return nil, nil, leveldb.ErrNotFound
}

// Back returns the last element of list l or nil if the list is empty.
func (cdb *ConsistentDB) Back(path string, ro *opt.ReadOptions) (key, value []byte, err error) {
	cdb.mu.Lock()
	_db := cdb.dbByPath[path]
	cdb.mu.Unlock()
	if _db == nil {
		return nil, nil, fmt.Errorf("leveldb not found, %s", path)
	}
	it := _db.NewIterator(nil, ro)
	defer it.Release()
	if it.Last() {
		k := it.Key()
		v := it.Value()
		return k, v, nil
	}
	return nil, nil, leveldb.ErrNotFound
}

// CompactRange compacts the underlying DB for the given key range.
// In particular, deleted and overwritten versions are discarded,
// and the data is rearranged to reduce the cost of operations
// needed to access the data. This operation should typically only
// be invoked by users who understand the underlying implementation.
//
// A nil Range.Start is treated as a key before all keys in the DB.
// And a nil Range.Limit is treated as a key after all keys in the DB.
// Therefore if both is nil then it will compact entire DB.
func (cdb *ConsistentDB) CompactRange(router string, r util.Range) error {
	_, db := cdb.LevelDB(router)
	if db == nil {
		return fmt.Errorf("leveldb not found, %s", router)
	}
	return db.CompactRange(r)
}

// GetProperty returns value of the given property router.
//
// Property names:
//
//	leveldb.num-files-at-level{n}
//		Returns the number of files at level 'n'.
//	leveldb.stats
//		Returns statistics of the underlying DB.
//	leveldb.iostats
//		Returns statistics of effective disk read and write.
//	leveldb.writedelay
//		Returns cumulative write delay caused by compaction.
//	leveldb.sstables
//		Returns sstables list for each level.
//	leveldb.blockpool
//		Returns block pool stats.
//	leveldb.cachedblock
//		Returns size of cached block.
//	leveldb.openedtables
//		Returns number of opened tables.
//	leveldb.alivesnaps
//		Returns number of alive snapshots.
//	leveldb.aliveiters
//		Returns number of alive iterators.
func (cdb *ConsistentDB) GetProperty(router string, name string) (value string, err error) {
	_, db := cdb.LevelDB(router)
	if db == nil {
		return "", fmt.Errorf("leveldb not found, %s", router)
	}
	return db.GetProperty(name)
}

// GetSnapshot returns a latest snapshot of the underlying DB. A snapshot
// is a frozen snapshot of a DB state at a particular point in time. The
// content of snapshot are guaranteed to be consistent.
//
// The snapshot must be released after use, by calling Release method.
func (cdb *ConsistentDB) GetSnapshot(router string) (*leveldb.Snapshot, error) {
	_, db := cdb.LevelDB(router)
	if db == nil {
		return nil, fmt.Errorf("leveldb not found, %s", router)
	}
	return db.GetSnapshot()
}

// NewIterator returns an iterator for the latest snapshot of the
// underlying DB.
// The returned iterator is not safe for concurrent use, but it is safe to use
// multiple iterators concurrently, with each in a dedicated goroutine.
// It is also safe to use an iterator concurrently with modifying its
// underlying DB. The resultant key/value pairs are guaranteed to be
// consistent.
//
// Slice allows slicing the iterator to only contains keys in the given
// range. A nil Range.Start is treated as a key before all keys in the
// DB. And a nil Range.Limit is treated as a key after all keys in
// the DB.
//
// WARNING: Any slice returned by iterator (e.g. slice returned by calling
// Iterator.Key() or Iterator.Key() methods), its content should not be modified
// unless noted otherwise.
//
// The iterator must be released after use, by calling Release method.
//
// Also read Iterator documentation of the leveldb/iterator package.
func (cdb *ConsistentDB) NewIterator(router string, slice *util.Range, ro *opt.ReadOptions) iterator.Iterator {
	_, db := cdb.LevelDB(router)
	if db == nil {
		return iterator.NewEmptyIterator(fmt.Errorf("leveldb not found, %s", router))
	}
	return db.NewIterator(slice, ro)
}

// SetReadOnly makes DB read-only. It will stay read-only until reopened.
func (cdb *ConsistentDB) SetReadOnly(router string) error {
	_, db := cdb.LevelDB(router)
	if db == nil {
		return fmt.Errorf("leveldb not found, %s", router)
	}
	return db.SetReadOnly()
}

// SizeOf calculates approximate sizes of the given key ranges.
// The length of the returned sizes are equal with the length of the given
// ranges. The returned sizes measure storage space usage, so if the user
// data compresses by a factor of ten, the returned sizes will be one-tenth
// the size of the corresponding user data size.
// The results may not include the sizes of recently written data.
func (cdb *ConsistentDB) SizeOf(router string, ranges []util.Range) (leveldb.Sizes, error) {
	_, db := cdb.LevelDB(router)
	if db == nil {
		return nil, fmt.Errorf("leveldb not found, %s", router)
	}
	return db.SizeOf(ranges)
}

// OpenTransaction opens an atomic DB transaction. Only one transaction can be
// opened at a time. Subsequent call to Write and OpenTransaction will be blocked
// until in-flight transaction is committed or discarded.
// The returned transaction handle is safe for concurrent use.
//
// Transaction is expensive and can overwhelm compaction, especially if
// transaction size is small. Use with caution.
//
// The transaction must be closed once done, either by committing or discarding
// the transaction.
// Closing the DB will discard open transaction.
func (cdb *ConsistentDB) OpenTransaction(router string) (*leveldb.Transaction, error) {
	_, db := cdb.LevelDB(router)
	if db == nil {
		return nil, fmt.Errorf("leveldb not found, %s", router)
	}
	return db.OpenTransaction()
}
