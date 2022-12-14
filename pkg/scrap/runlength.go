package scrap

/*
type RlArray struct {
	run    []uint32
	symbol []uint32
	dense  EfBytes
}*/

func RlEncode(v []int) ([]int, []int) {
	r := make([]int, 0, len(v))

	r = append(r, 0)
	last := 0
	c := v[0]
	for i := 1; i < len(v); i++ {
		if v[i] == c {
			r[last]++
		} else {
			r = append(r, 1)
			c = v[i]
		}
	}

	// extract the values
	o := make([]int, 0, len(r))
	pos := 0
	for i := range r {
		o[i] = v[pos]
		pos += r[i]
	}
	return r, o
}
