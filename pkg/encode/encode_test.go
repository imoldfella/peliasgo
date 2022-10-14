package encode

import (
	"log"
	"testing"
)

func Test_one(t *testing.T) {
	for i := 1; i <= 15; i++ {
		pyr := NewPyramid(0, i)

		x, y, z, err := pyr.Xyz(pyr.Len - 1)
		_, y2, _, _ := pyr.Xyz(pyr.Len - 2)
		if err != nil {
			panic(err)
		}
		log.Printf("size %d,%d,%d,%d,%d", pyr.Len, x, y, z, y2)
	}
}
