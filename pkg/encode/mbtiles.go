package encode

import (
	"database/sql"
	"log"
	"math"
	"os"
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

func MbtilesConvert(inpath, outpath string, maxfiles int) error {
	db, e := sql.Open("sqlite3", inpath)
	if e != nil {
		return e
	}
	getTileData, e := db.Prepare("select tile_data from tiles_data where tile_data_id=?")
	if e != nil {
		return e
	}
	pyr := NewPyramid(0, 15)
	//shallow := make([]uint64, pyr.Len)
	// prv[pyr] = tile_data
	// lowpyr[tile_data] = pyr

	// this is sparse, and we only want to visit the pyramid addresses that we care about. It's not clear how run length encoding works then?
	pyrToTdi := make([]uint32, pyr.Len)
	tdiToLowPyr := make([]uint32, pyr.Len)
	tdiToLowPyrLen := make([]uint32, pyr.Len)
	pyrAll := make([]uint32, pyr.Len)
	pyrN := 0
	for i := range tdiToLowPyr {
		tdiToLowPyr[i] = math.MaxUint32
	}

	// this is count, tile. lets us sort by count so we can pack the most reused tiles to make them easy to cache. Count here _should_ be a weight that includes usefullness but future work.
	rs, e := db.Query("select zoom_level , tile_column ,tile_row integer, tile_data_id from tiles_shallow")
	if e != nil {
		return e
	}
	var zoom_level, tile_column, tile_row, tile_data_id uint32
	n := 0
	for rs.Next() {
		rs.Scan(&zoom_level, &tile_column, &tile_row, &tile_data_id)
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

	// low helps because we already know where it is. we could reverse sort?
	// by using three columns here, we can run length encode the tiles.
	// pyr,start,len  (pyr -> fat pointer)
	// note that len is not exactly needed because we can stream from where we want to the end of the file.
	wr := OpenSplitLog(outpath, maxfiles)
	tileStart := make([]uint64, pyr.Len*3)
	var data []byte
	pos := uint64(0)
	previous_id := uint32(math.MaxUint32)
	j := 0

	// we only want to run length encode while the pyr increments by 1.
	// if it jumps by more than one, then we need to insert a don't know even if
	//
	for ii := range pyrAll {
		if ii%10000 == 0 {
			log.Printf("%d", ii)
		}

		i := pyrAll[ii]
		tile_data_id = pyrToTdi[i]
		if tile_data_id == previous_id {
			continue
		}
		tileStart[j+2] = uint64(i)
		if uint32(i) == tdiToLowPyr[tile_data_id] {
			e = getTileData.QueryRow(tile_data_id).Scan(&data)
			if e != nil {
				return e
			}
			wr.Write(data)
			tileStart[j] = pos
			tileStart[j+1] = uint64(len(data))
			tdiToLowPyrLen[tile_data_id] = uint32(len(data))
			pos += uint64(len(data))
		} else {
			pyr := tdiToLowPyr[tile_data_id]
			tileStart[j] = tileStart[pyr*2]
			tileStart[j+1] = tileStart[pyr*2+1]
		}
		j += 3
	}
	// so

	// we need a highPyr->filestart
	// this is sparse, so we need

	// we need a pyr -> highPyr, also sparse, run length encoded.
	// can we have just a pyr->filestart, sparse (encoding runs)
	wr.Close()
	os.WriteFile(outpath+".idx", toBytes(tileStart[:j]), os.ModePerm)

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
