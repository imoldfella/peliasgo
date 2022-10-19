package gop

import "unsafe"

func toBytes(d []uint64) []byte {
	return unsafe.Slice((*byte)(unsafe.Pointer(&d[0])), len(d)<<3)
}
