package encode

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"sync"

	"github.com/edsrzf/mmap-go"
)

type Database struct {
	path   string
	dbJson *DbJson

	table map[string]*TableReader

	index sync.Map

	// map the files and
	mm []mmap.MMap
	f  []*os.File
	mu sync.Mutex
}

func (d *Database) ReadIndex(pos uint64, size uint64) (*DbIndex, error) {
	ox, ok := d.index.Load(pos)
	if ok {
		return ox.(*DbIndex), nil
	}

	b, e := d.ReadBytes(pos, size)
	if e != nil {
		return nil, e
	}
	o, e := DecodeDbIndex(b)
	d.index.Store(pos, o)
	return o, e
}

func (d *Database) Slice(file uint64, from uint64, size uint64) []byte {
	mm := d.mm[file]
	return mm[from : from+size]
}

func OpenDatabase(path string) (*Database, error) {
	b, e := os.ReadFile(path + "/index.json")
	if e != nil {
		return nil, e
	}
	var js DbJson
	json.Unmarshal(b, &js)

	d := &Database{
		path:   path,
		dbJson: &js,
		table:  map[string]*TableReader{},
		index:  sync.Map{},
		mm:     make([]mmap.MMap, js.ChunkCount),
		f:      make([]*os.File, js.ChunkCount),
		mu:     sync.Mutex{},
	}
	for i := range d.f {
		d.f[i], e = os.OpenFile(fmt.Sprintf("%s/%d", d.path, i), os.O_RDWR, 0644)
		if e != nil {
			return nil, e
		}

		d.mm[i], _ = mmap.Map(d.f[i], mmap.RDWR, 0)
	}

	for name, tbl := range d.dbJson.Table {
		// we might as well read the root.
		idx, e := d.ReadIndex(tbl.Root, tbl.RootLength)
		if e != nil {
			return nil, e
		}

		d.table[name] = &TableReader{
			db:   d,
			t:    tbl,
			root: idx,
		}
	}

	return d, nil
}

func FileSlice(path string, from uint64, size uint64) ([]byte, error) {
	f, e := os.Open(path)
	if e != nil {
		return nil, e
	}
	b := make([]byte, size)
	_, e = f.ReadAt(b, int64(from))
	return b, e
}
func (d *Database) ReadBytes(pos uint64, size uint64) ([]byte, error) {
	ch := d.dbJson.ChunkSize
	file := pos / ch
	offset := pos - file*ch
	avail := ch - offset

	if avail >= size {
		// split across two files.
		return d.Slice(file, offset, size), nil
	} else {
		b := d.Slice(file, offset, avail)
		b2 := d.Slice(file+1, 0, size-avail)
		return append(b, b2...), nil
	}
}

// the root of the tree can be a separate json file that points to the roots of the tables.

// by making this binary and including the root of each table we could eliminate
// one round trip
type DbJson struct {
	ChunkSize  uint64              `json:"chunk_size,omitempty"`
	ChunkCount uint64              `json:"chunk_count,omitempty"`
	Table      map[string]*DbTable `json:"table,omitempty"`
}

type DbTable struct {
	Name       string `json:"name,omitempty"`
	Root       uint64 `json:"root,omitempty"`
	RootLength uint64 `json:"root_length,omitempty"`
	Height     int    `json:"height,omitempty"`
}

func (d *TableReader) Query(begin, end uint32) *TableIterator {

	return &TableIterator{
		Begin: begin,
		End:   end,
		Stack: []*DbIndex{},
		Key:   0,
		Value: []byte{},
	}
}

func (d *Database) Table(name string) (*TableReader, error) {
	table, ok := d.table[name]
	if !ok {
		return nil, fmt.Errorf("bad table %s", name)
	} else {
		return table, nil
	}

}

type TableReader struct {
	db   *Database
	t    *DbTable
	root *DbIndex
}

func (r *TableReader) Close() {

}

// func NewTableReader(db *Database) (*TableReader, error) {

// }

// this might take two reads if the first read returns a link

func (t *TableReader) Get(id uint32) ([]byte, error) {

	read1 := func(id uint32) ([]byte, error) {
		// we need to binary search the sorted values, then take the pivot
		// when we get to the leaf we can return the slice.
		r := t.root
		var e error
		for i := 0; i < t.t.Height; i++ {
			// find's smallest i such that fn is true; return <= id
			found := sort.Search(len(r.Key), func(x int) bool {
				return r.Key[x] >= id
			})
			if found == len(r.Key) || r.Key[found] > id {
				found--
			}
			// this may find an id that's smaller, but may work because of run length. we may want to keep run length when we decompress?

			start := r.Offset[found]
			end := r.Offset[found+1]
			r, e = t.db.ReadIndex(start, end-start)
			if e != nil {
				return nil, e
			}
		}
		found := sort.Search(len(r.Key), func(x int) bool {
			return r.Key[x] >= id
		})
		if found == len(r.Key) || r.Key[found] > id {
			found--
		}
		start := r.Offset[found]
		end := r.Offset[found+1]
		return t.db.ReadBytes(start, end-start)
	}
	d, e := read1(id)
	if e == nil && len(d) <= binary.MaxVarintLen32 {
		id2, e := binary.ReadUvarint(bytes.NewReader(d))
		if e != nil {
			return nil, e
		}
		return read1(uint32(id2))
	} else {
		return d, e
	}
}

type TableIterator struct {
	Begin, End uint32
	Stack      []*DbIndex
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

func (d *Database) Close() error {
	return nil
}

// we need a way to capture map metadata like zoom levels
// func (x *Database) Dump(path, out string) {
// 	for table, v := range x.table {
// 		iter := v.Query(0, math.MaxUint32)
// 		for iter.IsValid() {
// 			// write each blob as a file
// 			os.WriteFile(fmt.Sprintf("%s_%s_%d", out, table, iter.Key), iter.Value, os.ModePerm)
// 			iter.Next()
// 		}
// 	}
// }
