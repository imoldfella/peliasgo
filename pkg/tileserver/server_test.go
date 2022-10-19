package tileserver

import (
	"log"
	"testing"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
func Test_mbtiles(t *testing.T) {
	src, e := OpenMbtileSource("../../build/flat/output.mbtiles")
	check(e)

	ServeTiles(src, "8081", "/Users/jim/dev/datagrove/peliasgo/build/flat")
}
