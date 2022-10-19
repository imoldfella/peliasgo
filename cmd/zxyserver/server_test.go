package main

import "testing"

func Test_start(t *testing.T) {
	serveDir("../../build/flat", "8081")
}
