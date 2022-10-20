package encode

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"
)

type Number interface {
	uint32 | uint64
}

// todo: partition ranges, try bic
func Compress[T uint64 | uint32](d []T, w io.Writer) {
	var b [binary.MaxVarintLen64]byte
	for _, v := range d {
		n := binary.PutUvarint(b[:], uint64(v))
		w.Write(b[0:n])
	}
}
func Uncompress[T uint64 | uint32](d []T, r io.ByteReader) error {
	for i := range d {
		v, e := binary.ReadUvarint(r)
		if e != nil {
			return e
		}
		d[i] = T(v)
	}
	return nil
}

type DbIndex struct {
	Key    []uint32
	Offset []uint64 // offset is +1 if patch
}

func checkMonotonic[T Number](v []T) {
	for i := 1; i < len(v); i++ {
		if v[i-1] >= v[i] {
			log.Fatalf("not monotonic")
		}
	}
}
func EncodeDbIndex(key []uint32, offset []uint64) ([]byte, error) {
	if len(offset) != len(key)+1 {
		log.Fatalf("bad parameters")
	}
	checkMonotonic(key)
	checkMonotonic(offset)

	var b bytes.Buffer
	var bx [binary.MaxVarintLen64]byte
	n := binary.PutUvarint(bx[:], uint64(len(key)))
	b.Write(bx[0:n])
	Compress(key, &b)
	Compress(offset, &b) //bufio.NewWriter(&b))
	return b.Bytes(), nil
}

func DecodeDbIndex(data []byte) (*DbIndex, error) {
	r := bytes.NewReader(data)
	ln, e := binary.ReadUvarint(r)
	if e != nil {
		return nil, e
	}
	key := make([]uint32, ln)
	e = Uncompress(key, r)
	if e != nil {
		return nil, e
	}
	val := make([]uint64, ln+1)
	e = Uncompress(val, r)
	if e != nil {
		return nil, e
	}
	return &DbIndex{
		Key:    key,
		Offset: val,
	}, nil
}
