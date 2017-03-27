package goleveldown

import (
	"github.com/fiatjaf/levelup"
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

func (l LevelDown) ReadRange(opts *levelup.RangeOpts) levelup.ReadIterator {
	if opts == nil {
		opts = &levelup.RangeOpts{}
	}
	opts.FillDefaults()

	r := util.Range{}
	r.Start = []byte(opts.Start)
	r.Limit = []byte(opts.End)

	iter := l.db.NewIterator(&r, nil)

	// put the cursor in the right place.
	var isvalid bool
	if opts.Reverse {
		isvalid = iter.Last()
	} else {
		isvalid = iter.Next()
	}

	return &ReadIterator{
		iter:    iter,
		opts:    opts,
		count:   1,
		isvalid: isvalid,
	}
}

type ReadIterator struct {
	iter    iterator.Iterator
	opts    *levelup.RangeOpts
	count   int
	isvalid bool
}

func (ri *ReadIterator) Valid() bool {
	if !ri.isvalid {
		return false
	}
	if ri.count > ri.opts.Limit {
		return false
	}
	return true
}

func (ri *ReadIterator) Next() {
	ri.count++
	if ri.opts.Reverse {
		ri.isvalid = ri.iter.Prev()
	} else {
		ri.isvalid = ri.iter.Next()
	}
}

func (ri *ReadIterator) Key() string   { return string(ri.iter.Key()) }
func (ri *ReadIterator) Value() string { return string(ri.iter.Value()) }
func (ri *ReadIterator) Error() error  { return ri.iter.Error() }
func (ri *ReadIterator) Release()      { ri.iter.Release() }
