package main

import (
	"flag"
	"os"

	"github.com/imoldfella/peliasgo/pkg/geocoder"
)

// note: I should change
func main() {
	var p geocoder.Geocoder
	command := os.Args[1]
	switch command {
	case "pelias":
		url := flag.String("url", "http://192.168.1.188:4000", "# of iterations")
		flag.Parse()
		p, _ = geocoder.NewPelias(*url)
	case "":

	}
	geocoder.Geocode("/Users/jim/dev/asset/npi/py/diff.csv", "/Users/jim/dev/asset/npi/py/diff2.csv", p, 20)
}
