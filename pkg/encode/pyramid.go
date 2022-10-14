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

	return 0, 0, 0, fmt.Errorf("bad parameters")
}

func NewPyramid(min, max int) *Pyramid {
	o := &Pyramid{
		MinZoom: min,
		MaxZoom: max,
		h:       []*hilbert.Hilbert{},
	}
	o.h = make([]*hilbert.Hilbert, o.MaxZoom-o.MinZoom)
	for x := range o.h {
		o.h[x], _ = hilbert.NewHilbert(1 << x)
	}
	for z := o.MinZoom; z < o.MaxZoom; z++ {
		o.Len += (1 << z) * (1 << z)
	}
	return o
}
