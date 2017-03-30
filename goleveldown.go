package goleveldown

import (
	"os"

	"github.com/fiatjaf/levelup"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type LevelDown struct {
	db   *leveldb.DB
	path string
}

func NewDatabase(path string) levelup.DB {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		panic(err)
	}
	return &LevelDown{db, path}
}

func (l LevelDown) Close() { l.db.Close() }
func (l LevelDown) Erase() {
	l.Close()
	os.RemoveAll(l.path)
}

func (l LevelDown) Put(key, value []byte) error {
	return l.db.Put(key, value, nil)
}

func (l LevelDown) Get(key []byte) ([]byte, error) {
	data, err := l.db.Get(key, nil)
	if err != nil {
		if err == errors.ErrNotFound {
			return nil, levelup.NotFound
		}
		return nil, err
	}
	return data, nil
}

func (l LevelDown) Del(key []byte) error {
	return l.db.Delete(key, nil)
}

func (l LevelDown) Batch(ops []levelup.Operation) error {
	batch := new(leveldb.Batch)
	for _, op := range ops {
		switch op.Type {
		case "put":
			batch.Put(op.Key, op.Value)
		case "del":
			batch.Delete(op.Key)
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
	r.Start = opts.Start
	r.Limit = opts.End

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

func (ri *ReadIterator) Key() []byte   { return ri.iter.Key() }
func (ri *ReadIterator) Value() []byte { return ri.iter.Value() }
func (ri *ReadIterator) Error() error  { return ri.iter.Error() }
func (ri *ReadIterator) Release()      { ri.iter.Release() }
