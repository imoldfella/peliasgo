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
	repeated         map[int]int
	repeated_tile_id []int
}

var _ BlobArray = (*MbtileSet)(nil)

func (b *MbtileSet) Len() int {
	return b.p.Len
}

const s1 = "select tile_id,count(*) from shallow_tiles group by 1 having count(*)>1 order by 1"

func NewMbtileSet(path string) (*MbtileSet, error) {
	// open sqlite and find the repeated tile ids.
	db, e := sql.Open("sqlite", path)
	if e != nil {
		return nil, e
	}
	s1, e := db.Prepare(s1)
	if e != nil {
		return nil, e
	}
	s2, e := db.Prepare(" ")
	if e != nil {
		return nil, e
	}
	// find the repeats an intialize those directories.

	return &MbtileSet{
		p:         Pyramid{},
		db:        db,
		getTileId: s1,
		getTile:   s2,
	}, nil
}

// start with the most used tiles

// getChunks implements Bucketable
// repeats at beginning.
func (m *MbtileSet) Read(start int, end int) ([][]byte, []int, error) {
	var b []byte
	r := make([][]byte, 0, end-start)
	rs := []int{}

	for ; start < len(m.repeated_tile_id); start++ {
		m.getTile.QueryRow(m.repeated_tile_id[start]).Scan(&b)
		// we might need to copy the block? does scan overwrite?
		r = r.append(b)
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
