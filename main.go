package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"sync/atomic"
)

type Geocoder interface {
	Geocode(zip, address string) (float64, float64, error)
}

// note: I should change
func main() {
	p, _ := NewPelias("// http://192.168.1.188:4000")
	geocode("/Users/jim/dev/asset/npi/py/diff.csv", "/Users/jim/dev/asset/npi/py/diff2.csv", p, 20)
}

func geocode(in string, out string, geo Geocoder, threads int) {
	var count int64 = 0
	update := func() {
		c := atomic.AddInt64(&count, 1)
		if c%8 == 0 {
			log.Printf("Count %d", count)
		}
	}

	m, e := os.ReadFile(in)
	if e != nil {
		panic(e)
	}
	fout, e := os.Create(out)
	if e != nil {
		panic(e)
	}
	ferr, e := os.Create(out + ".err")
	if e != nil {
		panic(e)
	}
	var mu sync.Mutex

	mn := strings.Split(string(m), "\n")

	var wg sync.WaitGroup
	wg.Add(threads)
	for i := 0; i <= threads; i++ {
		go func(i int) {
			for j := i; j < len(mn); j += threads {
				addr := mn[j]
				update()
				v := strings.Split(addr, "^")
				if len(v) < 2 {
					continue
				}
				lat, lon, err := geo.Geocode(v[0], v[1])
				mu.Lock()
				if lat == 100 || err != nil {
					if err != nil {
						log.Printf("error %v,%s", err, addr)
					}
					ferr.WriteString(addr + "\n")
				} else {
					fout.WriteString(fmt.Sprintf("%s,%f,%f\n", addr, lat, lon))
				}
				mu.Unlock()
			}
			wg.Done()
		}(i)
	}

	wg.Wait()
}

type Location struct {
	Type     string `json:"type"`
	Features []struct {
		Type     string `json:"type"`
		Geometry struct {
			Type        string    `json:"type"`
			Coordinates []float64 `json:"coordinates"`
		} `json:"geometry"`
		Properties struct {
			Layer       string  `json:"layer"`
			Source      string  `json:"source"`
			Name        string  `json:"name"`
			Housenumber string  `json:"housenumber"`
			Street      string  `json:"street"`
			Postalcode  string  `json:"postalcode"`
			Confidence  float64 `json:"confidence"`
			MatchType   string  `json:"match_type"`
			Accuracy    string  `json:"accuracy"`
			Country     string  `json:"country"`
			CountryA    string  `json:"country_a"`
			Region      string  `json:"region"`
			RegionA     string  `json:"region_a"`
			CountyA     string  `json:"county_a"`
			Locality    string  `json:"locality"`
			Label       string  `json:"label"`
		} `json:"properties"`
	} `json:"features"`
}
