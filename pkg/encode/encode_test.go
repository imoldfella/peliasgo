package encode

import (
	"log"
	"testing"
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
func compare(mbtiles, dgtiles string) {
	o, e := OpenDatabase(dgtiles)
	defer o.Close()
	check(e)
	tbl, e := o.Table("map")
	b, e := tbl.Get(0)
	check(e)
	log.Print(string(b))
}

func Test_compare(t *testing.T) {
	const outpath = "/Users/jim/dev/datagrove/peliasgo/build/flat/db"
	const allpath = "/Users/jim/dev/datagrove/peliasgo/build/flat/output."
	compare(allpath, outpath)
}
func Test_big(t *testing.T) {
	const outpath = "/Users/jim/dev/datagrove/peliasgo/build/flat/db"
	const allpath = "/Users/jim/dev/datagrove/peliasgo/build/flat/output.mbtiles"
	o, e := CreateDb(outpath, 20*1024*1024)
	if e != nil {
		panic(e)
	}
	defer o.Close()

	e = MbtilesConvert(o, "map", allpath)
	if e != nil {
		panic(e)
	}
	o.Close()

	compare(allpath, outpath)
}

func Test_o1(t *testing.T) {
	const outpath = "/Users/jim/dev/datagrove/peliasgo/build/flat/d"
	const smallpath = "/Users/jim/dev/datagrove/peliasgo/build/flat/monaco.mbtiles"
	o, e := CreateDb(outpath, 20*1024*1024)
	if e != nil {
		panic(e)
	}
	defer o.Close()

	e = MbtilesConvert(o, "map", smallpath)
	if e != nil {
		panic(e)
	}
}
