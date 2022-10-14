package main

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
)

/*

# dense array format
message ef {
	int len
	bytes data
}

# run length encoded format with symbols
# the 0 symbol is used to indicate a run of references to the dense vector (in order)
message rl_dict {
	repeated uint32 runlength
	repeated uint32 symbol
	ef dense
}

*/

type EfBytes struct {
	Len   int
	Bytes []byte
}

type RlArray struct {
	run    []uint32
	symbol []uint32
	dense  EfBytes
}

type BucketInfo struct {
	Len        int
	BucketSize int
	reader     Bucketable
}

// we can write the bucket index together with the tail chunk, tricky because we need to compress.
type Bucketable interface {
	// getChunks can return an int or a []byte for each position. the int is to reference a symbol and is only used if the []byte is nil
	getChunks(bucket int, start, end int) ([][]byte, []int)
	getSymbols() [][]byte
}

// each bucket will have a
func writeBuckets(o *BucketInfo, prefix string, maxFileSize int, nthread int) error {
	symdata := o.reader.getSymbols()
	if nthread == 0 {
		nthread = runtime.NumCPU()
	}
	ln := o.Len
	perFile := o.BucketSize
	var nextShard int64

	var wg sync.WaitGroup
	wg.Add(nthread)
	for i := 0; i < nthread; i++ {
		go func() {
			for {
				f := (int)(atomic.AddInt64(&nextShard, 1))
				bpfx := fmt.Sprintf("%s_%d", prefix, f)
				df := OpenSplitLog(bpfx, maxFileSize)
				from := f * perFile
				to := from + perFile

				if from > ln {
					wg.Done()
					return
				}
				if to > ln {
					to = ln
				}

				sym := make([]int, 0, to-from)
				run := make([]int, 0, to-from)
				length := make([]int, 0, to-from)
				c := -1 // not matched by any ref.
				pos := 0
				for from != to {
					// get some chunks from the data
					chunk, ref := o.reader.getChunks(f, from, to)

					df.WriteAll(chunk)
					for i, v := range chunk {
						if c != ref[i] {
							c = ref[i]
							sym = append(sym, c)
							run = append(run, 1)
						} else {
							run[len(run)-1]++
						}
						sym[i] = ref[i]
						if ref[i] != 0 {
							pos += len(v)
							length = append(length, pos)
						}
					}
					from += len(chunk)
				}
				tail := df.b

				var header []byte
				// the header is run length sparse vector with symbols
				// and an ef compressed vector with positions

				// write the header + tail together.
				os.WriteFile(bpfx, append(header, tail...), os.ModePerm)
			}
		}()
	}
	// write the symbol table if it exists, it might not.
	// here we just have a
	if len(symdata) != 0 {
		var symbolTable []byte
		//
		os.WriteFile(prefix, symbolTable, os.ModePerm)
	}

	wg.Wait()
	return nil
}
