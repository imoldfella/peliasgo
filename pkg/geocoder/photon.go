package geocoder

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type PhotonCoder struct {
	server string
}

func NewPhoton(server string) (Geocoder, error) {
	return &PhotonCoder{
		server: server,
	}, nil
}

// Geocode implements Geocoder
func (p *PhotonCoder) Geocode(zip string, address string) (float64, float64, error) {

	url := fmt.Sprintf("%s/api?q=%s", p.server, url.QueryEscape(address+","+" "+zip))
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

	if e != nil || len(target.Features) == 0 {
		return 0, 0, fmt.Errorf("not found")
	}
	a := target.Features[0].Geometry.Coordinates

	return a[1], a[0], nil
}

var _ Geocoder = (*PhotonCoder)(nil)
