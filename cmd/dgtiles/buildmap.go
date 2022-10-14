package main

// spread the tiles equally amount N files.
// in most cases it will make sense to load the entire file. but its also possible to read the beginning and then the rail found

// 1. convert all the tiles to hilbert
// 2. split gthe hilbert id's by the number of files.
// 3. build each hilbert key'd file multithreaded. (not lots of thread, should block on disk)

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
)

// we should put the top re-used tiles in its own file most won't be reused
// z, hilbert, tile_or_ref.
//

// track the repeated

type Builder struct {
	Outdir      string
	Mbtile      string
	MinZoom     int
	MaxZoom     int
	Nshards     int
	Nthread     int
	MaxFileSize int

	// created a dictionary file we can use for re-used tiles.
	len          int
	tilesPerfile int
	repeated     map[int]int
	nextFile     int64
	chunkCount   int64 // quick check that we haven't exceeded Nfiles
}

func (o *Builder) Build() error {
	getTile := func(n int) []byte {
		return nil
	}
	if o.Nthread == 0 {
		o.Nthread = runtime.NumCPU()
	}
	cnt := 0
	for z := o.MinZoom; z < o.MaxZoom; z++ {
		cnt += (1 << z) * (1 << z)
	}
	o.len = cnt
	o.tilesPerfile = o.len / o.Nshards
	o.chunkCount = int64(o.Nshards)

	countfile := func() {
		atomic.AddInt64(&o.chunkCount, 1)
	}

	var wg sync.WaitGroup
	wg.Add(o.Nthread)
	for i := 0; i < o.Nthread; i++ {
		go func() {

			for {
				f := (int)(atomic.AddInt64(&o.nextFile, 1))

				w, e := os.Create(fmt.Sprintf("%s/d%d", o.Outdir, f))
				if e != nil {
					panic(e)
				}
				a := f * (o.tilesPerfile)
				if a > o.len {
					wg.Done()
					return
				}
				b := a + o.tilesPerfile
				if b > o.len {
					b = o.len
				}
				start := make([]int, b-a)

				for a != b {
					ppos := 0
					countfile()
					for ; a != b && ppos < o.MaxFileSize; a++ {
						b := getTile(a)
						w.Write(b)
						start[a] = ppos + len(b)
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

func NewBuilder(mbtile string, outdir string) *Builder {
	return &Builder{
		Outdir:       outdir,
		Mbtile:       mbtile,
		MinZoom:      0,
		MaxZoom:      0,
		Nshards:      10000,
		Nthread:      0,
		len:          0,
		tilesPerfile: 0,
		repeated:     map[int]int{},
		nextFile:     0,
	}
}
