package encode

import (
	"fmt"

	"github.com/google/hilbert"
)

type Pyramid struct {
	Len     int
	MinZoom int
	MaxZoom int
	h       []*hilbert.Hilbert // one fore each zoom
}

func (p *Pyramid) Xyz(id int) (x, y, z int, e error) {
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

func NewPyramid(minZoom, maxZoom int) *Pyramid {
	o := &Pyramid{
		MinZoom: minZoom,
		MaxZoom: maxZoom,
		h:       []*hilbert.Hilbert{},
	}
	o.h = make([]*hilbert.Hilbert, o.MaxZoom-o.MinZoom)
	for x := range o.h {
		o.h[x], _ = hilbert.NewHilbert(1 << (x + o.MinZoom))
	}
	for z := o.MinZoom; z < o.MaxZoom; z++ {
		o.Len += (1 << z) * (1 << z)
	}
	return o
}
