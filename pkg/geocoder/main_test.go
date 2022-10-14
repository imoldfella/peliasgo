package geocoder

import (
	"log"
	"testing"
)

func pgcode() Geocoder {
	code, e := NewPostgis("user=jim dbname=jim sslmode=disable")
	if e != nil {
		panic(e)
	}
	return code
}
func photonCode() Geocoder {
	code, e := NewPhoton("http://localhost:2322")
	if e != nil {
		panic(e)
	}
	return code
}

func Test_test1(t *testing.T) {
	// geocode the diff file
	Geocode("./build/diff.csv", "./build/diff2.csv", pgcode(), 4)
}
func Test_test3(t *testing.T) {
	Geocode("./build/npia.csv", "./build/npia2.csv", pgcode(), 1)
}
func Test_test2(t *testing.T) {

	ok, lat, lon := pgcode().Geocode("19063", "5 jonathan morris circle")
	log.Printf("%v,%f,%f", ok, lat, lon)
}

func Test_photon1(t *testing.T) {
	ok, lat, lon := photonCode().Geocode("19063", "5 jonathan morris circle")
	log.Printf("%v,%f,%f", ok, lat, lon)
}
func Test_photonDiff(t *testing.T) {
	// geocode the diff file
	Geocode("./build/diff.csv", "./build/photon/diff2.csv", photonCode(), 40)
}
func Test_photonAll(t *testing.T) {
	Geocode("./build/npia.csv", "./build/photon/npia2.csv", photonCode(), 40)
}

func Test_mapbuild(t *testing.T) {

}
