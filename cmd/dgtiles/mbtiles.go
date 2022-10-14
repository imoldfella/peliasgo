package main

import (
	"database/sql"

	"github.com/imoldfella/peliasgo/pkg/encode"
)

type MbtileSet struct {
	p         encode.Pyramid
	db        *sql.DB
	getTileId *sql.Stmt
	getTile   *sql.Stmt

	// maps a pyramid index to a value in the Repeats dictionary.
	repeats map[int]int

	// this value is returned, e.g. one entry for water, land, don't care
	Repeats map[int]int
}

var _ Tileset = (*MbtileSet)(nil)

func NewMbtileSet(path string) (*MbtileSet, error) {
	// open sqlite and find the repeated tile ids.
	db, e := sql.Open("sqlite", path)
	if e != nil {
		return nil, e
	}
	s1, e := db.Prepare("select")
	if e != nil {
		return nil, e
	}
	s2, e := db.Prepare(" ")
	if e != nil {
		return nil, e
	}
	// find the repeats an intialize those directories.

	return &MbtileSet{
		p:         encode.Pyramid{},
		db:        db,
		getTileId: s1,
		getTile:   s2,
		repeats:   map[int]int{},
		Repeats:   map[int]int{},
	}, nil
}

// GetTile implements Tileset
func (m *MbtileSet) GetTile(index int, data []byte) ([]byte, int, error) {
	if r, ok := m.repeats[index]; ok {
		return nil, r, nil
	}
	// compute xyz from index, then query for the tile_id,
	x, y, z, err := m.p.Xyz(index)

	var tileid int
	m.getTileId.QueryRow().Scan(&tileid)

}

// Pyramid implements Tileset
func (*MbtileSet) Pyramid() encode.Pyramid {
	panic("unimplemented")
}
