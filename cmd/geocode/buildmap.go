package main

// 1. convert all the tiles to hilbert
// 2. split the hilbert id's by the number of files.
// 3. build each hilbert key'd file multithreaded. (not lots of thread, should block on disk)

type Options struct {
	Nfiles int
}

func DefaultOptions() *Options {
	return &Options{
		Nfiles: 10000,
	}
}
func build(mbtile string, outdir string, opt *Options) {

}
