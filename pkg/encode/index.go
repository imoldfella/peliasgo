package encode

type Number interface {
	uint32 | uint64
}
type Index32 = Index[uint32]

func compress32(d []uint32) []byte {
	// v := make([]uint, len(d))
	// for i := range d{
	// 	v[i]=uint(d[i])
	// }
	// obj,_ := ef.From(v)
	// return
	return []byte{}
}

// we always write the key ascending.
// we alternate between key,run,pos,offset and
type Index[T Number] struct {
	db     *DbEncoder
	key    []uint32
	offset []uint64
	size   []uint32

	// there might be a smarter approach, but for now keep these in ram until we can append them onto the end
	leafPivot []T
	leafData  []byte
}

func NewIndex[T Number](d *DbEncoder) *Index[T] {
	return &Index[T]{
		db: d,
	}
}

const (
	kDense = iota
	kSparse
)

func (x *Index[T]) pack() {
	// a run of values stored sequentially in the log
	dense := []uint64{}

	// skey will hold all the non-sequential keys
	// the first key is non-sequential
	// then
	skey := []uint32{}
	srun := []uint32{}
	soffset := []uint64{}
	ssize := []uint32{}

	state := 1
	for i := 1; i < len(x.key); i++ {
		seq := x.key[i-1]+1 == x.key[i]
		repeat := x.offset[i-1] == x.offset[i]
		seqvalue := x.offset[i-1]+uint64(x.size[i-1]) == x.offset[i]

		if seq && seqvalue {
			dense = append(dense, x.offset[i])
			srun[len(srun)-1]++
		} else if false {

		} else if seq && repeat {
			srun[len(srun)-1]++
		} else {
			srun = append(srun, 1)
			skey = append(skey, x.key[i])
		}
	}
	//v1 := compress32(x.key)

}
func (x *Index[T]) Add(key uint32, offset uint64, size uint32) {
	x.key = append(x.key, key)
	x.offset = append(x.offset, offset)
	x.size = append(x.size, size)

	if len(x.key) == 32*1024 {
		x.pack()
	}
}

func (x *Index[T]) Close() {

}
