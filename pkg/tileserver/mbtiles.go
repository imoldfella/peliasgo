package tileserver

import (
	//"database/sql"

	"log"
	"sync"

	"github.com/bvinc/go-sqlite-lite/sqlite3"
	"github.com/imoldfella/peliasgo/pkg/encode"
	//_ "github.com/mattn/go-sqlite3"
	//_ "github.com/marcboeker/go-duckdb"
)

type MbtileSource struct {
	db    *sqlite3.Conn
	tstmt *sqlite3.Stmt
	meta  string
	mu    sync.Mutex
}

var _ TileSource = (*MbtileSource)(nil)

// GetMetadata implements TileSource
// why does mbtiles put somethings in fields and some things in json?

func (m *MbtileSource) GetMetadata() string {
	return m.meta
}

// get implements TileSource
func (s *MbtileSource) Get(z int, x int, y int) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var data []byte
	defer s.tstmt.Reset()
	e := s.tstmt.Bind(z, x, y)
	if e == nil {
		hasRow, _ := s.tstmt.Step()
		if hasRow {
			e = s.tstmt.Scan(&data)
		}
	}
	if e != nil {
		log.Printf("sql error %v", e)
	}
	return data, e
}

func OpenMbtileSource(path string) (*MbtileSource, error) {
	//db, e := sql.Open("sqlite3", path)
	db, e := sqlite3.Open(path)
	if e != nil {
		panic(e)
	}
	tstmt, e := db.Prepare("select tile_data from tiles where zoom_level=? and tile_column=? and tile_row=?")
	if e != nil {
		panic(e)
	}
	b, e := encode.MetadataToJson(db)
	return &MbtileSource{
		db:    db,
		tstmt: tstmt,
		meta:  string(b),
		mu:    sync.Mutex{},
	}, e
}
func (s *MbtileSource) Close() error {
	s.tstmt.Close()
	return s.db.Close()
}
