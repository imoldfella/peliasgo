package main

// spread the tiles equally amount N files.
// in most cases it will make sense to load the entire file. but its also possible to read the beginning and then the rail found

// 1. convert all the tiles to hilbert
// 2. split gthe hilbert id's by the number of files.
// 3. build each hilbert key'd file multithreaded. (not lots of thread, should block on disk)

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/google/hilbert"
)

// Create a Hilbert curve for mapping to and from a 16 by 16 space.

// we should put the top re-used tiles in its own file most won't be reused
// z, hilbert, tile_or_ref.
//

// track the repeated

type Pyramid struct {
	Len     int
	MinZoom int
	MaxZoom int
	h       []*hilbert.Hilbert // one fore each zoom
}

func NewPyramid(min, max int) *Pyramid {
	o := &Pyramid{
		MinZoom: min,
		MaxZoom: max,
		h:       []*hilbert.Hilbert{},
	}
	o.h = make([]*hilbert.Hilbert, o.MaxZoom-o.MinZoom)
	for x := range o.h {
		o.h[x], _ = hilbert.NewHilbert(1 << x)
	}
	for z := o.MinZoom; z < o.MaxZoom; z++ {
		o.Len += (1 << z) * (1 << z)
	}
	return o
}

type Map struct {
	p      Pyramid
	Outdir string
	Mbtile string

	Nshards     int
	Nthread     int
	MaxFileSize int

	// created a dictionary file we can use for re-used tiles.

	tilesPerfile int
	repeated     map[int]int
	nextShard    int64
	chunkCount   int64 // quick check that we haven't exceeded Nfiles

}

type SplitLog struct {
	b      []byte
	w      io.Writer
	prefix string
	count  int
}

// Close implements io.WriteCloser
func (*SplitLog) Close() error {
	panic("unimplemented")
}

func (o *SplitLog) Flush() error {
	return os.WriteFile(fmt.Sprintf("%s%d", o.prefix, o.count), o.b, os.ModePerm)
	o.b = o.b[:0]
	o.count++
	return nil
}

// Write implements io.Writer
func (*SplitLog) Write(p []byte) (n int, err error) {
	// write as much as we can,then start the next file
	return 0, nil
}

var _ io.WriteCloser = (*SplitLog)(nil)

func OpenSplitLog(p string, maxfile int) *SplitLog {
	return &SplitLog{
		prefix: p,
		count:  0,
		b:      make([]byte, maxfile),
	}
}

func (o *Map) Build() error {
	ln := o.p.Len
	o.tilesPerfile = ln / o.Nshards
	o.chunkCount = int64(o.Nshards)

	countfile := func() {
		atomic.AddInt64(&o.chunkCount, 1)
	}

	// return either bytes of the file or nil and the dictionary id as negative
	getTile := func(w io.Writer, id int) (int, error) {
		return 0, nil
	}

	var wg sync.WaitGroup
	wg.Add(o.Nthread)
	for i := 0; i < o.Nthread; i++ {
		go func() {
			for {
				f := (int)(atomic.AddInt64(&o.nextShard, 1))
				df := OpenSplitLog(o.Outdir+"/d", o.MaxFileSize)
				a := f * (o.tilesPerfile)
				if a > ln {
					wg.Done()
					return
				}
				b := a + o.tilesPerfile
				if b > ln {
					b = ln
				}

				cl := make([]int, b-a)

				for a != b {
					ppos := 0
					countfile()
					for ; a != b && ppos < o.MaxFileSize; a++ {
						cl[a] = getTile(a)
						if id > 0 {
							w.Write(b)
						}
						cl[a] = len(b)
						ppos += len(b)
					}
				}
				// compress the starting integers for this file
				// estimate
			}
		}()
	}
	wg.Wait()
	return nil
}

func NewBuilder(mbtile string, outdir string) *Map {
	return &Map{
		Outdir: outdir,
		Mbtile: mbtile,

		Nshards: 10000,
		Nthread: runtime.NumCPU(),

		// computed, return values
		tilesPerfile: 0,
		repeated:     map[int]int{},
		nextShard:    0,
	}
}
