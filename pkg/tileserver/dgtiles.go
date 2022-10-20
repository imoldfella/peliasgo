package tileserver

import (
	"github.com/imoldfella/peliasgo/pkg/encode"
)

type DgtileSource struct {
	db       *encode.Database
	mp       *encode.TableReader
	pyr      encode.Pyramid
	metadata string
}

// GetMetadata implements TileSource
func (d *DgtileSource) GetMetadata() string {
	return d.metadata
}

var _ TileSource = (*DgtileSource)(nil)

// get implements TileSource
func (s *DgtileSource) Get(z int, x int, y int) ([]byte, error) {
	id, e := s.pyr.FromXyz(x, y, z)
	if e != nil {
		return nil, e
	}
	return s.mp.Get(id)
}

// where should meta data be stored? Probably 0, since it tells us the size of the pyramid.

func OpenDgtileSource(path string) (*DgtileSource, error) {
	db, e := encode.OpenDatabase(path)
	if e != nil {
		return nil, e
	}

	tbl, e := db.Table("map")
	if e != nil {
		return nil, e
	}

	// we need to parse the metadata in order to create the pyramid encoder
	// but we can return the whole thing anyway.
	metadata, e := tbl.Get(0)
	if e != nil {
		return nil, e
	}

	return &DgtileSource{
		db:       db,
		mp:       tbl,
		pyr:      *encode.NewPyramid(0, 14),
		metadata: string(metadata),
	}, nil
}
func (s *DgtileSource) Close() error {
	return s.db.Close()
}
