package encode

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"

	"github.com/bvinc/go-sqlite-lite/sqlite3"
	"github.com/jfcg/sorty/v2"
	_ "github.com/paulmach/orb"
)

// needed for cors?
// .Methods(http.MethodGet, http.MethodPut, http.MethodPatch, http.MethodOptions)
type Metadata struct {
	Tilejson     string            `json:"tilejson"`
	Scheme       string            `json:"scheme"`
	Type         string            `json:"type"`
	Format       string            `json:"format"`
	Tiles        []string          `json:"tiles"`
	Bounds       []float64         `json:"bounds"`
	Center       []int             `json:"center"`
	Name         string            `json:"name"`
	Version      string            `json:"version"`
	Description  string            `json:"description"`
	Minzoom      int               `json:"minzoom"`
	Maxzoom      int               `json:"maxzoom"`
	VectorLayers []json.RawMessage `json:"vector_layers"`
}

func MetadataToJson(db *sqlite3.Conn) ([]byte, error) {
	// can we marshal this?
	meta := map[string]string{}
	rows, e := db.Prepare("select name,value from metadata")
	if e != nil {
		return nil, e
	}
	var key, value string
	for {
		hasRow, e := rows.Step()
		if e != nil {
			log.Print(e)
		}
		if !hasRow {
			break
		}
		rows.Scan(&key, &value)
		meta[key] = value
	}
	rows.Close()

	me := `{
		"tilejson": "2.0.0",
		"scheme": "tms",
		"type": "baselayer",
		"format": "pbf",
		"tiles": [
			"http://localhost:8081/rlp/{z}/{x}/{y}.pbf"
		],
		"bounds": [
			%s
		],
		"center": [%s],
		"name": "OpenMapTiles",
		"version": "3.13.1",
		"description": "Tile config based on OpenMapTiles schema",
		"minzoom": 0,
		"maxzoom": 14
	}
	`

	js := fmt.Sprintf(me, meta["bounds"], meta["center"])
	var x Metadata
	json.Unmarshal([]byte(js), &x)

	var v Metadata
	json.Unmarshal([]byte(meta["json"]), &v)
	x.VectorLayers = v.VectorLayers

	b, _ := json.Marshal(&x)
	return b, nil
}

// cover countries, states, other localities
// https://pkg.go.dev/github.com/paulmach/orb/maptile/tilecover#section-readme

func MbtilesConvert(o *DbEncoder, table, inpath string) error {
	pyr := NewPyramid(0, 15)
	// these are oversized to create a trivial perfect hash. for small maps it would be better to use a normal hash table
	pyrAll := make([]uint32, pyr.Len)
	pyrToTdi := make([]uint32, pyr.Len)
	// low because we write from low to high, lower pyr address already written.
	tdiToLowPyr := make([]uint32, pyr.Len)

	maxTdi := uint32(0)

	db, e := sqlite3.Open(inpath)
	if e != nil {
		return e
	}
	getTileData, e := db.Prepare("select tile_data from tiles_data where tile_data_id=?")
	if e != nil {
		return e
	}
	var data []byte
	getTile := func(tdi uint32) ([]byte, error) {
		defer getTileData.Reset()
		data = data[:0]
		e = getTileData.Bind(int(tdi))
		if e == nil {
			hasRow, _ := getTileData.Step()
			if hasRow {
				getTileData.Scan(&data)
			}
		}
		if len(data) == 0 {
			log.Printf("missing %d", tdi)
			data = data[:0]
		}
		return data, nil
	}

	// I can do this in parallel. splitting on zoom_level and tile_column
	// most of the work in is zoom level, I could also test with an early out.
	readShallow := func() error {
		log.Printf("reading shallow")
		rs, e := db.Prepare("select zoom_level , tile_column ,tile_row integer, tile_data_id from tiles_shallow")
		if e != nil {
			return e
		}
		defer rs.Close()

		var zoom_level, tile_column, tile_row, tile_data_id uint32
		var i1, i2, i3, i4 int
		pyrN := 0
		for {
			if zoom_level > 10 {
				break
			}
			hasRow, e := rs.Step()
			if e != nil {
				return e
			}
			if !hasRow {
				break
			}
			e = rs.Scan(&i1, &i2, &i3, &i4)
			zoom_level = uint32(i1)
			tile_column = uint32(i2)
			tile_row = uint32(i3)
			tile_data_id = uint32(i4)
			if e != nil {
				panic(e)
			}
			if tile_data_id > maxTdi {
				maxTdi = tile_data_id
			}
			// note that pyrid cannot be 0, we +1 from the normal id to leave space for null.
			pyrid, e := pyr.FromXyz(tile_column, tile_row, zoom_level)
			if e != nil {
				return e
			}
			pyrAll[pyrN] = pyrid
			pyrN++
			pyrToTdi[pyrid] = tile_data_id
			if tdiToLowPyr[tile_data_id] == 0 || tdiToLowPyr[tile_data_id] > pyrid {
				tdiToLowPyr[tile_data_id] = pyrid
			}
			if pyrN%1000000 == 0 {
				log.Printf("%d", pyrN/1000000)
			}
		}
		maxTdi++
		pyrAll = pyrAll[0:pyrN] // trim off unused pieces
		return nil
	}

	writeTable := func() error {

		log.Printf("writing table")

		var tile_data_id uint32 = 0
		tileStart := make([]uint64, maxTdi)
		tileLength := make([]uint32, maxTdi)
		idx := OpenIndex32(o, table)

		metadata, e := MetadataToJson(db)
		if e != nil {
			return e
		}
		idx.Add(0, metadata)

		j := 0
		for j < len(pyrAll) {
			if j%1000000 == 0 {
				log.Printf("%d", j/1000000)
			}
			i := pyrAll[j]
			tile_data_id = pyrToTdi[i]
			if uint32(i) == tdiToLowPyr[tile_data_id] {
				data, e := getTile(tile_data_id)
				if e != nil {
					return e
				}
				pos, _ := idx.Add(i, data)
				tileStart[tile_data_id] = pos
				tileLength[tile_data_id] = uint32(len(data))
				pos += uint64(len(data))
				j++
			} else {
				// advance j to the end of the run
				bg := j
				for j++; j < len(pyrAll); j++ {
					if pyrToTdi[pyrAll[j]] != tile_data_id {
						break
					}
				}
				// create a varint pointer to previous block
				var b [32]byte
				// run length
				n := binary.PutUvarint(b[:], uint64(j-bg))
				// start and length of copied block
				n += binary.PutUvarint(b[n:], tileStart[tile_data_id])
				n += binary.PutUvarint(b[n:], uint64(tileLength[tile_data_id]))
				idx.Add(i, b[:n])
			}
		}
		beginIndex := o.wr.Length()
		idx.Close()
		log.Printf("data,index=%d,%d", beginIndex, o.wr.Length()-beginIndex)
		return nil
	}

	e = readShallow()
	if e != nil {
		return e
	}
	log.Printf("sorting")
	sorty.SortSlice(pyrAll)

	e = writeTable()
	if e != nil {
		return e
	}

	return nil
}
