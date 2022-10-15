package encode

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
)

// blob 0 can hold configuration information for the entire array.
// including zoom levels
// future work: instead of duplicate reference allow a diff program

// we could require the references up front, these are always small, from shallow.
// then we could
type BlobArray interface {
	Len() int
	Copy(o *SnapshotWriter, start, end int) error
}

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
prefix_symbol - array of blobs
prefix_#shard - array of symbols or blobs
prefix_#shared_data - overflow
*/

type EfBytes struct {
	Len   int
	Bytes []byte
}

// each bucket will have a
func writeBuckets(o BlobArray, shards int, prefix string, maxFileSize int, nthread int) error {

	if nthread == 0 {
		nthread = runtime.NumCPU()
	}
	ln := o.Len()
	perFile := ln / shards
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
					chunk, ref, err := o.Read(from, to)
					if err != nil {
						panic(err)
					}

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

	wg.Wait()
	return nil
}
