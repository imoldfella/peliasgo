package encode

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
)

type DbReader struct {
	path   string
	dbJson *DbJson
	Table  map[string]*TableReader
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

type TableReader struct {
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

func (d *DbReader) IsValid() bool {
	return false
}

func NewDbReader(path string) (*DbReader, error) {
	b, e := os.ReadFile(path)
	if e != nil {
		return nil, e
	}
	var js DbJson
	json.Unmarshal(b, &js)
	return &DbReader{
		path:   path,
		dbJson: &js,
	}, nil
}

// we need a way to capture map metadata like zoom levels
func Dump(path, out string) {
	x, e := NewDbReader(path)
	if e != nil {
		panic(e)
	}

	for table, v := range x.Table {
		iter := v.Query(0, math.MaxUint32)
		for iter.IsValid() {
			// write each blob as a file
			os.WriteFile(fmt.Sprintf("%s_%s_%d", out, table, iter.Key), iter.Value, os.ModePerm)
			iter.Next()
		}
	}

}
