package encode

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
)

type Database struct {
	path   string
	dbJson *DbJson
	table  map[string]*TableReader
}

func (d *TableReader) Query(begin, end uint32) *TableIterator {
	return &TableIterator{
		Begin: begin,
		End:   end,
		Stack: []*Block{},
		Key:   0,
		Value: []byte{},
	}
}

func (d *Database) Table(name string) (*TableReader, error) {
	return &TableReader{
		db: d,
	}, nil
}

type TableReader struct {
	db *Database
}

func (t *TableReader) Get(id uint32) ([]byte, error) {
	return nil, nil
}

type Block struct {
	key    []uint32
	offset []uint64
}

type TableIterator struct {
	Begin, End uint32
	Stack      []*Block
	Key        uint32
	Value      []byte
}

func (it *TableIterator) IsValid() bool {
	return it.Key != it.End
}
func (it *TableIterator) Next() bool {
	return it.Key != it.End
}

func (d *Database) IsValid() bool {
	return false
}

func OpenDatabase(path string) (*Database, error) {
	b, e := os.ReadFile(path)
	if e != nil {
		return nil, e
	}
	var js DbJson
	json.Unmarshal(b, &js)
	return &Database{
		path:   path,
		dbJson: &js,
	}, nil
}

func (d *Database) Close() error {
	return nil
}

// we need a way to capture map metadata like zoom levels
func (x *Database) Dump(path, out string) {
	for table, v := range x.table {
		iter := v.Query(0, math.MaxUint32)
		for iter.IsValid() {
			// write each blob as a file
			os.WriteFile(fmt.Sprintf("%s_%s_%d", out, table, iter.Key), iter.Value, os.ModePerm)
			iter.Next()
		}
	}
}
