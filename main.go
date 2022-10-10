package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"sync/atomic"
)

func main() {
	geocode("/Users/jim/dev/asset/npi/py/diff.csv", "/Users/jim/dev/asset/npi/py/diff2.csv")
}

func geocode(in string, out string) {
	var count int64 = 0
	update := func() {
		c := atomic.AddInt64(&count, 1)
		if c%100 == 0 {
			log.Printf("Count %d", count)
		}
	}
	threads := 20
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
				if len(v) < 3 {
					continue
				}
				err, lat, lon := get1b(v[0], v[1], v[2])
				mu.Lock()
				if lat == 100 || err != nil {
					if err != nil {
						log.Printf("error %v", err)
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
func get1(zip, state, address string) (error, float64, float64) {
	url := fmt.Sprintf("http://192.168.1.188:4000/v1/search?text=%s", url.QueryEscape(address+","+state+" "+zip))
	resp, err := http.Get(url)
	if err != nil {
		return err, 0, 0
	}
	var target Location
	b, e := io.ReadAll(resp.Body)
	resp.Body.Close()
	if e != nil {
		log.Printf("errror %v, %v", b, err)
		return e, 0, 0
	}
	e = json.Unmarshal(b, &target)
	if len(target.Features) == 0 {
		return get1b(zip, state, address)
	}
	if e != nil || len(target.Features) == 0 {
		log.Printf("not found %v", e)
		return e, 0, 0
	}
	a := target.Features[0].Geometry.Coordinates

	return nil, a[1], a[0]
}
func get1b(zip, state, address string) (error, float64, float64) {
	url := fmt.Sprintf("http://192.168.1.188:4000/v1/search/structured?region=%s&address=%s&postalcode=%s", url.QueryEscape(state), url.QueryEscape(address), url.QueryEscape(zip))
	resp, err := http.Get(url)
	if err != nil {
		return err, 0, 0
	}
	var target Location
	b, e := io.ReadAll(resp.Body)
	resp.Body.Close()
	if e != nil {
		return err, 0, 0
	}
	e = json.Unmarshal(b, &target)
	if e != nil {
		return e, 0, 0
	}
	if len(target.Features) == 0 {
		return nil, 100, 0
	}
	a := target.Features[0].Geometry.Coordinates
	return nil, a[1], a[0]
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
