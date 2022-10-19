package encode

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type DbEncoder struct {
	path   string
	wr     *SplitLogWriter
	dbJson DbJson
}

// the root of the tree can be a separate json file that points to the roots of the tables.
type DbJson struct {
	Table []*DbTable `json:"table,omitempty"`
}
type DbTable struct {
	Name       string `json:"name,omitempty"`
	Root       uint64 `json:"root,omitempty"`
	RootLength int    `json:"root_length,omitempty"`
}

// Creates a directory of files and a manifest file (index.json)
func CreateDb(path string, maxfiles uint64) (*DbEncoder, error) {
	e := os.RemoveAll(path)
	if e != nil {
		return nil, e
	}

	e = os.Mkdir(path, os.ModePerm)
	if e != nil {
		return nil, e
	}
	wr := OpenSplitLog(path, maxfiles)
	return &DbEncoder{
		path:   path,
		wr:     wr,
		dbJson: DbJson{},
	}, nil
}

func (o *DbEncoder) Close() error {
	o.wr.Close()
	b, _ := json.Marshal(&o.dbJson)
	return os.WriteFile(o.path+"/index.json", b, os.ModePerm)
}

func (d *DbEncoder) Add(table string, data []byte) error {
	pos, e := d.wr.Write(data)
	if e != nil {
		return e
	}
	d.dbJson.Table = append(d.dbJson.Table, &DbTable{
		Name:       table,
		Root:       pos,
		RootLength: len(data),
	})
	return nil
}

// Create a logical file split over files of at most maxfile bytes.
type SplitLogWriter struct {
	// byte size of target file size in bytes, eg 25M
	maxfile uint64
	// buffer the file an write it at once.
	fileBuffer []byte
	// prefix is database directory
	prefix string
	// filename = {table}/{offset}
	offset uint64
	//
}

func OpenSplitLog(p string, maxfile uint64) *SplitLogWriter {
	return &SplitLogWriter{
		maxfile:    maxfile,
		fileBuffer: make([]byte, 0, maxfile),
		prefix:     p,
		offset:     0,
	}
}

func (s *SplitLogWriter) Length() uint64 {
	return s.offset*s.maxfile + uint64(len(s.fileBuffer))
}

func (s *SplitLogWriter) Close() error {
	err := os.WriteFile(fmt.Sprintf("%s/%d", s.prefix, s.offset), s.fileBuffer, os.ModePerm)
	s.fileBuffer = nil
	return err
}

// returns the offset in the file where value is written
func (s *SplitLogWriter) Write(value []byte) (uint64, error) {
	vl := uint64(len(value))
	bl := uint64(len(s.fileBuffer))

	// write as much as we can,then start the next file
	var err error
	remain := s.maxfile - bl
	if vl < remain {
		s.fileBuffer = append(s.fileBuffer, value...)
	} else {
		s.fileBuffer = append(s.fileBuffer, value[0:remain]...)
		err = os.WriteFile(fmt.Sprintf("%s/%d", s.prefix, s.offset), s.fileBuffer, os.ModePerm)
		s.offset += 1
		s.fileBuffer = value[remain:]
	}
	return s.offset*s.maxfile + bl, err
}

type Number interface {
	uint32 | uint64
}

// todo: partition ranges, try bic
func compress[T uint64 | uint32](d []T, w io.Writer) {
	var b [binary.MaxVarintLen64]byte
	for _, v := range d {
		n := binary.PutUvarint(b[:], uint64(v))
		w.Write(b[0:n])
	}
	w.Write([]byte{0})
}

type Index32 struct {
	db     *DbEncoder
	table  string
	pos    uint64
	key    []uint32
	offset []uint64 // offset is +1 if patch

	packed    [][]byte
	packedKey []uint32
}

func OpenIndex32(d *DbEncoder, table string) *Index32 {
	return &Index32{
		db:        d,
		table:     table,
		pos:       0,
		key:       []uint32{},
		offset:    []uint64{0},
		packed:    [][]byte{},
		packedKey: []uint32{},
	}
}

// the problem with building the parent here is we don't know where to point yet.
func (x *Index32) pack() {
	x.packedKey = append(x.packedKey, x.key[0])
	var b bytes.Buffer
	compress(x.key, &b)
	compress(x.offset, &b) //bufio.NewWriter(&b))
	x.packed = append(x.packed, b.Bytes())
}

// we need to manage 0 length values (copy previous value)
func (x *Index32) Add(key uint32, data []byte) (uint64, error) {
	pos, _ := x.db.wr.Write(data)
	x.pos += uint64(len(data))
	x.key = append(x.key, key)
	x.offset = append(x.offset, x.pos)
	if len(x.key) == 32*1024 {
		x.pack()
	}
	return pos, nil
}

func (x *Index32) Close() {
	// when we close the index we need to write the packed blocks and if more than one a directory to those blocks.
	if len(x.key) > 0 {
		x.pack()
	}

	if len(x.packed) > 1 {
		parent := OpenIndex32(x.db, x.table)
		for i := range x.packed {
			parent.Add(x.packedKey[i], x.packed[i])
		}
		parent.Close()
	} else {
		// this is the top of the tree, so insert into the database
		x.db.Add(x.table, x.packed[0])
	}
}
