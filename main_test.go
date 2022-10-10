package main

import (
	"log"
	"testing"
)

func Test_test1(t *testing.T) {
	// geocode the diff file
	geocode("./build/diff.csv", "./build/diff2.csv")
}
func Test_test3(t *testing.T) {
	geocode("./build/npia.csv", "./build/npia2.csv")
}
func Test_test2(t *testing.T) {
	ok, lat, lon := get1b("19063", "PA", "5 jonathan morris circle")

	log.Printf("%v,%f,%f", ok, lat, lon)
}
