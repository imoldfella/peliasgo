package encode

import (
	"encoding/json"
	"fmt"
	"os"
)

type DbWriter struct {
	path   string
	wr     *SplitLogWriter
	dbJson DbJson
}

// Creates a directory of files and a manifest file (index.json)
func CreateDb(path string, chunkSize uint64) (*DbWriter, error) {
	e := os.RemoveAll(path)
	if e != nil {
		return nil, e
	}

	e = os.Mkdir(path, os.ModePerm)
	if e != nil {
		return nil, e
	}
	wr := OpenSplitLog(path, chunkSize)
	return &DbWriter{
		path: path,
		wr:   wr,
		dbJson: DbJson{
			ChunkSize: chunkSize,
			Table:     map[string]*DbTable{},
		},
	}, nil
}

func (o *DbWriter) Close() error {
	o.wr.Close()
	o.dbJson.ChunkCount = o.wr.ChunkCount
	b, _ := json.Marshal(&o.dbJson)
	return os.WriteFile(o.path+"/index.json", b, os.ModePerm)
}

func (d *DbWriter) Add(table string, data []byte, height int) error {
	pos, e := d.wr.Write(data)
	if e != nil {
		return e
	}
	d.dbJson.Table[table] = &DbTable{
		Name:       table,
		Root:       pos,
		RootLength: uint64(len(data)),
		Height:     height,
	}
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
	// filename = {table}/{ChunkCount}
	ChunkCount uint64
	//
}

func OpenSplitLog(p string, maxfile uint64) *SplitLogWriter {
	return &SplitLogWriter{
		maxfile:    maxfile,
		fileBuffer: make([]byte, 0, maxfile),
		prefix:     p,
		ChunkCount: 0,
	}
}

func (s *SplitLogWriter) Length() uint64 {
	return s.ChunkCount*s.maxfile + uint64(len(s.fileBuffer))
}

func (s *SplitLogWriter) Close() error {
	err := os.WriteFile(fmt.Sprintf("%s/%d", s.prefix, s.ChunkCount), s.fileBuffer, os.ModePerm)
	s.ChunkCount++
	s.fileBuffer = nil
	return err
}

// returns the offset in the file where value is written
func (s *SplitLogWriter) Write(value []byte) (uint64, error) {
	vl := uint64(len(value))
	bl := uint64(len(s.fileBuffer))
	start := s.Length()

	// write as much as we can,then start the next file
	var err error
	remain := s.maxfile - bl
	if vl < remain {
		s.fileBuffer = append(s.fileBuffer, value...)
	} else {
		s.fileBuffer = append(s.fileBuffer, value[0:remain]...)
		err = os.WriteFile(fmt.Sprintf("%s/%d", s.prefix, s.ChunkCount), s.fileBuffer, os.ModePerm)
		s.ChunkCount++
		s.fileBuffer = value[remain:]
	}
	return start, err
}

type Index32 struct {
	db     *DbWriter
	table  string
	key    []uint32
	offset []uint64

	packed    [][]byte
	packedKey []uint32
	height    int
}

func OpenIndex32(d *DbWriter, table string, height int) *Index32 {
	return &Index32{
		db:        d,
		table:     table,
		key:       []uint32{},
		offset:    []uint64{},
		packed:    [][]byte{},
		packedKey: []uint32{},
		height:    height,
	}
}

// the problem with building the parent here is we don't know where to point yet.
func (x *Index32) pack() {
	x.offset = append(x.offset, x.db.wr.Length())
	x.packedKey = append(x.packedKey, x.key[0])
	b, e := EncodeDbIndex(x.key, x.offset)
	if e != nil {
		panic(e)
	}
	x.packed = append(x.packed, b)
	x.key = x.key[:0]
	x.offset = x.offset[:0]
}

// we need to manage 0 length values (copy previous value)
func (x *Index32) Add(key uint32, data []byte) (uint64, error) {
	if len(x.key) == 32*1024 {
		x.pack()
	}
	pos, _ := x.db.wr.Write(data)
	x.key = append(x.key, key)
	x.offset = append(x.offset, pos)

	return pos, nil
}

func (x *Index32) Close() {
	// when we close the index we need to write the packed blocks and if more than one a directory to those blocks.
	if len(x.key) > 0 {
		x.pack()
	}

	if len(x.packed) > 1 {
		parent := OpenIndex32(x.db, x.table, x.height+1)
		for i := range x.packed {
			parent.Add(x.packedKey[i], x.packed[i])
		}
		parent.Close()
	} else {
		// this is the top of the tree, so insert into the database
		x.db.Add(x.table, x.packed[0], x.height)
	}
}
