package geocoder

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type PeliasCoder struct {
	server string
}

func NewPelias(server string) (Geocoder, error) {
	return &PeliasCoder{
		server: server,
	}, nil
}

// Geocode implements Geocoder
func (p *PeliasCoder) Geocode(zip string, address string) (float64, float64, error) {
	structured := func(zip, address string) (float64, float64, error) {
		url := fmt.Sprintf("%s/v1/search/structured?address=%s&postalcode=%s", p.server, url.QueryEscape(address), url.QueryEscape(zip))
		resp, err := http.Get(url)
		if err != nil {
			return 0, 0, err
		}
		var target Location
		b, e := io.ReadAll(resp.Body)
		resp.Body.Close()
		if e != nil {
			return 0, 0, err
		}
		e = json.Unmarshal(b, &target)
		if e != nil {
			return 0, 0, e
		}
		if len(target.Features) == 0 {
			return 100, 0, nil
		}
		a := target.Features[0].Geometry.Coordinates
		return a[1], a[0], nil
	}

	unstructured := func(zip, address string) (float64, float64, error) {
		url := fmt.Sprintf("%s/v1/search?text=%s", p.server, url.QueryEscape(address+","+" "+zip))
		resp, err := http.Get(url)
		if err != nil {
			return 0, 0, err
		}
		var target Location
		b, e := io.ReadAll(resp.Body)
		resp.Body.Close()
		if e != nil {
			log.Printf("errror %v, %v", b, err)
			return 0, 0, e
		}
		e = json.Unmarshal(b, &target)
		if len(target.Features) == 0 {
			return structured(zip, address)
		}
		if e != nil || len(target.Features) == 0 {
			log.Printf("not found %v", e)
			return 0, 0, e
		}
		a := target.Features[0].Geometry.Coordinates

		return a[1], a[0], nil
	}
	_ = unstructured

	return structured(zip, address)

}

var _ Geocoder = (*PeliasCoder)(nil)
