package leveldown

import (
	"github.com/fiatjaf/go-levelup"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type LevelDown struct {
	db *leveldb.DB
}

func NewDatabase(path string) levelup.DB {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		panic(err)
	}
	return &LevelDown{db}
}

func (l LevelDown) Put(key, value string) error {
	return l.db.Put([]byte(key), []byte(value), nil)
}

func (l LevelDown) Get(key string) (string, error) {
	data, err := l.db.Get([]byte(key), nil)
	if err != nil {
		if err == errors.ErrNotFound {
			return "", levelup.NotFound
		}
		return "", err
	}
	return string(data), nil
}

func (l LevelDown) Del(key string) error {
	return l.db.Delete([]byte(key), nil)
}

func (l LevelDown) Batch(ops []levelup.Operation) error {
	batch := new(leveldb.Batch)
	for _, op := range ops {
		switch op["type"] {
		case "put":
			batch.Put([]byte(op["key"]), []byte(op["value"]))
		case "del":
			batch.Delete([]byte(op["key"]))
		}
	}
	return l.db.Write(batch, nil)
}

func (l LevelDown) ReadRange(opts levelup.RangeOpts) levelup.ReadIterator {
	r := util.Range{}
	if opts.Start != "" {
		r.Start = []byte(opts.Start)
	}
	if opts.End != "" {
		r.Limit = []byte(opts.End)
	}
	if opts.Limit <= 0 {
		opts.Limit = 9999999
	}

	iter := l.db.NewIterator(&r, nil)
	return &ReadIterator{
		iter:    iter,
		opts:    opts,
		scanned: 0,
	}
}

type ReadIterator struct {
	iter    iterator.Iterator
	opts    levelup.RangeOpts
	scanned int
}

func (ri *ReadIterator) Next() bool {
	has := ri.iter.Next()
	if has {
		ri.scanned++
		if ri.opts.Limit < ri.scanned {
			return false
		}
	}
	return has
}

func (ri *ReadIterator) Key() string   { return string(ri.iter.Key()) }
func (ri *ReadIterator) Value() string { return string(ri.iter.Value()) }
func (ri *ReadIterator) Error() error  { return ri.iter.Error() }
func (ri *ReadIterator) Release()      { ri.iter.Release() }
