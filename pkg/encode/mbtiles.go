package encode

import (
	"database/sql"
	"log"
	"math"
	"unsafe"

	"github.com/jfcg/sorty/v2"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/paulmach/orb"
)

// cover countries, states, other localities
// https://pkg.go.dev/github.com/paulmach/orb/maptile/tilecover#section-readme

type Mbtiles struct {
	path string
	pyr  *Pyramid
	// pyr -> tile_id
	// we can keep a cache of tile_ids and only write them once.
	shallow     []uint64
	highPyramid []uint64
	fileStart   []uint64
}

func NewMbtiles(inpath string) (*Mbtiles, error) {
	return &Mbtiles{
		path: inpath,
	}, nil
}

func toBytes(d []uint64) []byte {
	return unsafe.Slice((*byte)(unsafe.Pointer(&d[0])), len(d)<<3)
}

type DbEncoder struct {
	wr *SplitLogWriter
}

func NewDbEncoder(path string, maxfiles int) *DbEncoder {
	wr := OpenSplitLog(path, maxfiles)
	return &DbEncoder{
		wr: wr,
	}
}
func (o *DbEncoder) Close() {

	o.wr.Close()
}

func (o *DbEncoder) MbtilesConvert(table, inpath string) error {

	db, e := sql.Open("sqlite3", inpath)
	if e != nil {
		return e
	}
	getTileData, e := db.Prepare("select tile_data from tiles_data where tile_data_id=?")
	if e != nil {
		return e
	}
	rs, e := db.Query("select zoom_level , tile_column ,tile_row integer, tile_data_id from tiles_shallow")
	if e != nil {
		return e
	}

	pyr := NewPyramid(0, 15)
	pyrToTdi := make([]uint32, pyr.Len)
	// low helps because we already know where it is. we could reverse sort?
	tdiToLowPyr := make([]uint32, pyr.Len)
	pyrAll := make([]uint32, pyr.Len)
	pyrN := 0
	for i := range tdiToLowPyr {
		tdiToLowPyr[i] = math.MaxUint32
	}
	var zoom_level, tile_column, tile_row, tile_data_id uint32
	maxTdi := uint32(0)
	n := 0
	for rs.Next() {
		// are these guaranteed to increase tile_data_id? not clear. we could sort by tile_data_id, is that worth it though?
		rs.Scan(&zoom_level, &tile_column, &tile_row, &tile_data_id)
		if tile_data_id+1 > maxTdi {
			maxTdi = tile_data_id + 1
		}
		// generate the hilbert id and
		id, e := pyr.FromXyz(tile_column, tile_row, zoom_level)
		if e != nil {
			return e
		}
		pyrAll[pyrN] = id
		pyrN++
		//shallow[i] = (id << 32) + uint64(tile_data_id)
		pyrToTdi[id] = tile_data_id
		if tdiToLowPyr[tile_data_id] > id {
			tdiToLowPyr[tile_data_id] = id
		}
		n++
		if n%10000 == 0 {
			log.Printf("%d", n)
		}
	}
	pyrAll = pyrAll[0:pyrN]
	sorty.SortSlice(pyrAll)

	tileStart := make([]uint64, maxTdi)
	tileLength := make([]uint32, maxTdi)
	var data []byte
	pos := uint64(0)
	idx := NewIndex[uint32](o)
	defer idx.Close()
	for j := range pyrAll {
		if j%10000 == 0 {
			log.Printf("%d", j)
		}
		i := pyrAll[j]
		tile_data_id = pyrToTdi[i]
		if uint32(i) == tdiToLowPyr[tile_data_id] {
			e = getTileData.QueryRow(tile_data_id).Scan(&data)
			if e != nil {
				return e
			}
			o.wr.Write(data)
			tileStart[tile_data_id] = pos
			tileLength[tile_data_id] = uint32(len(data))
			idx.Add(i, pos, uint32(len(data)))
			pos += uint64(len(data))
		} else {
			prevpos := tileStart[tile_data_id]
			prevlen := uint32(tile_data_id)
			idx.Add(i, prevpos, prevlen)
		}
	}
	return nil
}

func PartitionRange(from, to, partitions int) []int {
	ln := to - from
	splits := make([]int, partitions+1)
	splits[partitions] = ln
	delta := float64(ln) / float64(partitions)
	o := delta
	for i := 1; i < partitions; i++ {
		splits[i] = int(math.Ceil(o))
		o += delta
	}
	return splits
}

func Unpair(x uint64) (uint32, uint32) {
	return uint32(x >> 32), uint32(x & ((1 << 32) - 1))
}

// Why Partitions?
// 1. Average of one file to read a z,x,y. Two reads with ranges; one streaming for header
// 2. Some parallelism
// 3. Split into smallish files rather than pack is more cdn friendly.

const dry_run = true
