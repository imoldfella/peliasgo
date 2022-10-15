package encode

import (
	"database/sql"
)

type MbtileSet struct {
	p         Pyramid
	db        *sql.DB
	getTileId *sql.Stmt
	getTile   *sql.Stmt

	// maps a pyramid index to a value in the Repeats dictionary.
	repeated map[int]int
	// this value is returned, e.g. one entry for water, land, don't care
	Repeats [][]byte
}

var _ Bucketable = (*MbtileSet)(nil)

func NewMbtileSet(path string) (*BlobArray, error) {
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

	reader := &MbtileSet{
		p:         Pyramid{},
		db:        db,
		getTileId: s1,
		getTile:   s2,
	}
	return &BlobArray{
		Len:    reader.p.Len,
		reader: reader,
	}, nil
}

// getChunks implements Bucketable
// repeats at beginning.
func (m *MbtileSet) getChunks(start int, end int) ([][]byte, []int, error) {
	r := make([][]byte, 0, end-start)
	rs := []int{}

	for ; start < len(m.Repeats); start++ {
		r = append(r, m.Repeats[start])
	}

	start -= len(m.Repeats)
	for ; start != end; start++ {
		x, y, z, err := m.p.Xyz(start)
		if err != nil {
			return nil, nil, err
		}
		var tileid int
		err = m.getTileId.QueryRow(x, y, z).Scan(&tileid)
		if err != nil {
			return nil, nil, err
		}

		var b []byte
		sym, ok := m.repeated[tileid]
		if !ok {
			m.getTile.QueryRow(tileid).Scan(&b)
			r = append(r, b)
			rs = append(rs, 0)
		} else {
			rs = append(rs, sym)
		}
	}
	return r, rs, nil
}
