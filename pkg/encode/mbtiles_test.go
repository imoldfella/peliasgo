package encode

import "testing"

const allpath = "/Users/jim/dev/datagrove/peliasgo/build/world/output.mbtiles"
const smallpath = "/Users/jim/dev/datagrove/peliasgo/build/world/monaco.mbtiles"
const outpath = "/Users/jim/dev/datagrove/peliasgo/build/flat/d"

func Test_o1(t *testing.T) {
	o := NewDbEncoder(outpath, 20*1024*1024)
	e := o.MbtilesConvert("map", smallpath)
	if e != nil {
		panic(e)
	}

}
