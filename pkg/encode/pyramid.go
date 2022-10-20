package encode

import (
	"fmt"

	"github.com/google/hilbert"
)

// normal pyramid encoding except we add 1 so that 0 can be used for "null"

type Pyramid struct {
	Len     int
	MinZoom int
	MaxZoom int
	h       []*hilbert.Hilbert // one fore each zoom
	start   []int
}

func (p *Pyramid) Xyz(id int) (x, y, z int, e error) {
	id-- // make it 0 based
	for z, h := range p.h {
		sz := h.N * h.N
		if id >= sz {
			id -= sz
		} else {
			x, y, _ = h.Map(id)
			return x, y, z + p.MinZoom, nil
		}
	}
	return 0, 0, 0, fmt.Errorf("out of bounds")
}
func (p *Pyramid) FromXyz(x, y, z int) (uint32, error) {
	return p.FromXyz32(uint32(x), uint32(y), uint32(z))
}

func (p *Pyramid) FromXyz32(x, y, z uint32) (uint32, error) {
	id, e := p.h[z].MapInverse(int(x), int(y))
	id += p.start[z]
	return uint32(id) + 1, e
}

func NewPyramid(minZoom, maxZoom int) *Pyramid {
	o := &Pyramid{
		MinZoom: minZoom,
		MaxZoom: maxZoom,
		h:       []*hilbert.Hilbert{},
	}
	o.h = make([]*hilbert.Hilbert, o.MaxZoom-o.MinZoom)
	o.start = make([]int, 1+o.MaxZoom-o.MinZoom)
	cnt := 0
	for x := range o.h {
		o.h[x], _ = hilbert.NewHilbert(1 << (x + o.MinZoom))
		cnt += o.h[x].N * o.h[x].N
		o.start[x+1] = cnt

	}
	o.Len = cnt
	// for z := o.MinZoom; z < o.MaxZoom; z++ {
	// 	o.Len += (1 << z) * (1 << z)
	// }
	return o
}
