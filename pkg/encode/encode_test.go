package encode

import (
	"bytes"
	"log"
	"testing"

	"github.com/bvinc/go-sqlite-lite/sqlite3"
	"golang.org/x/exp/slices"
)

func Test_pyramid(t *testing.T) {
	for i := 1; i <= 15; i++ {
		pyr := NewPyramid(0, i)

		x, y, z, err := pyr.Xyz(pyr.Len - 1)
		_, y2, _, _ := pyr.Xyz(pyr.Len - 2)
		if err != nil {
			panic(err)
		}
		log.Printf("size %d,%d,%d,%d,%d", pyr.Len, x, y, z, y2)
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
func Test_metadata(t *testing.T) {
	o, e := OpenDatabase("/Users/jim/dev/datagrove/peliasgo/build/flat/db")
	check(e)
	defer o.Close()
	tbl, e := o.Table("map")
	check(e)
	b, e := tbl.Get(1)
	check(e)
	log.Printf("%v", b)
}
func compare1(mbtiles, dgtiles string) {
	o, e := OpenDatabase(dgtiles)
	check(e)
	defer o.Close()

	tbl, e := o.Table("map")
	check(e)
	b, e := tbl.Get(0)
	check(e)
	log.Print(string(b))
}

func compare(mbtiles, dgtiles string, zoom int) {
	i, e := OpenMbtileIterator(mbtiles)
	check(e)
	defer i.Close()
	o, e := OpenDatabase(dgtiles)
	check(e)
	defer o.Close()
	tbl, e := o.Table("map")
	check(e)

	for i.Next() && i.Z < zoom {
		data, e := tbl.Get(i.id)
		check(e)
		if !bytes.Equal(data, i.data) {
			log.Printf("%v,%v", i, data)
		}
	}
}

func Test_compare(t *testing.T) {
	const outpath = "/Users/jim/dev/datagrove/peliasgo/build/flat/db"
	const allpath = "/Users/jim/dev/datagrove/peliasgo/build/flat/output.mbtiles"
	compare(allpath, outpath, 11)

	// compare 10 levels

}
func Test_big(t *testing.T) {
	const outpath = "/Users/jim/dev/datagrove/peliasgo/build/flat/db"
	const allpath = "/Users/jim/dev/datagrove/peliasgo/build/flat/output.mbtiles"
	o, e := CreateDb(outpath, 20*1024*1024)
	if e != nil {
		panic(e)
	}
	defer o.Close()

	maxzoom := 15
	e = MbtilesConvert(o, "map", allpath, maxzoom)
	if e != nil {
		panic(e)
	}
	o.Close()

	compare(allpath, outpath, maxzoom)
}

func Test_o1(t *testing.T) {
	const outpath = "/Users/jim/dev/datagrove/peliasgo/build/flat/d"
	const smallpath = "/Users/jim/dev/datagrove/peliasgo/build/flat/monaco.mbtiles"
	o, e := CreateDb(outpath, 20*1024*1024)
	if e != nil {
		panic(e)
	}
	defer o.Close()

	e = MbtilesConvert(o, "map", smallpath, 15)
	if e != nil {
		panic(e)
	}
}

func Test_compress(t *testing.T) {
	a := []uint32{32, 64, 91, 234}
	x := make([]uint32, len(a))
	var b bytes.Buffer
	Compress(a, &b)
	Uncompress(x, bytes.NewReader(b.Bytes()))
	if !slices.Equal(a, x) {
		log.Fatal("ugh")
	}
}

type MbtileIterator struct {
	conn    *sqlite3.Conn
	stmt    *sqlite3.Stmt
	X, Y, Z int
	data    []byte
	e       error
	pyr     *Pyramid
	id      uint32
}

func (m *MbtileIterator) Next() bool {
	hasRow, _ := m.stmt.Step()
	m.e = m.stmt.Scan(&m.Z, &m.X, &m.Y, &m.data)
	m.id, m.e = m.pyr.FromXyz(m.X, m.Y, m.Z)
	return hasRow && m.e == nil
}

func (m *MbtileIterator) Close() {
	m.stmt.Close()
	m.conn.Close()
}

func OpenMbtileIterator(path string) (*MbtileIterator, error) {
	conn, e := sqlite3.Open(path, sqlite3.OPEN_READONLY)
	if e != nil {
		return nil, e
	}
	const s1 = `select zoom_level,
	tile_column,
	tile_row,
	tile_data
	from tiles`
	stmt, e := conn.Prepare(s1)
	if e != nil {
		return nil, e
	}
	return &MbtileIterator{
		conn: conn,
		stmt: stmt,
		X:    0,
		Y:    0,
		Z:    0,
		data: []byte{},
		pyr:  NewPyramid(0, 15),
	}, nil
}
