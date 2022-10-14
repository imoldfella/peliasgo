package main

import (
	"bytes"
	"database/sql"
	"encoding/binary"
	"fmt"
	"math"
	"os"
)

// read the zoom level 14 tiles

func main() {

	db, err := sql.Open("sqlite3", "file")
	if err != nil {
		panic(err)
	}
	_ = db
}

type Xy struct {
	X int
	Y int
}

func Tms(lat, lon float64, zoom int) Xy {
	z := 1 << zoom
	n := float64(z)

	x := int(math.Floor((lon + 180.0) / 360.0 * n))
	if float64(x) >= n {
		x = int(n - 1)
	}
	y := int(math.Floor((1.0 - math.Log(math.Tan(lat*math.Pi/180.0)+1.0/math.Cos(lat*math.Pi/180.0))/math.Pi) / 2.0 * n))
	// flip
	y = (1 << zoom) - y - 1
	return Xy{x, y}
}

func RoundBlock(c1, c2 Xy) (Xy, Xy) {
	// we want to extract 256 tiles at a time, 16 on a side
	x1 := (c1.X) & ^15
	x2 := (c2.X + 15) & ^15
	y1 := (c1.Y) & ^15
	y2 := (c2.Y + 15) & ^15
	return Xy{x1, y1}, Xy{x2, y2}
}

func WriteUsa(db *sql.DB, root string) error {
	// write a file for each block
	c1, c2 := Usbbox()
	d1, d2 := RoundBlock(c1, c2)
	stmt, err := db.Prepare("select tile_data from tiles where zoom_level=14 and tile_column =? and tile_row =?")
	if err != nil {
		return err
	}

	for x := d1.X; x < d2.X; x += 16 {
		for y := d1.Y; y < d2.Y; y += 16 {
			var b bytes.Buffer
			sz := make([]byte, 1024)
			for x1 := 0; x1 < 4; x1++ {
				for y1 := 0; y1 < 4; y1++ {
					xv := x + x1
					yv := y + y1
					var d []byte
					_ = stmt.QueryRow(&xv, &yv).Scan(&d)
					index := 4*x1 + y1
					binary.LittleEndian.PutUint32(sz[4*index:], (uint32)(len(d)))
					b.Write(d)
				}
			}
			fn := fmt.Sprintf("%s/t%d_%d", root, x, y)
			os.WriteFile(fn, append(sz, b.Bytes()...), os.ModePerm)
		}
	}
	return nil
}

func Usbbox() (Xy, Xy) {
	return Tms(24.396308, -124.848974, 14), Tms(49.384358, -66.885444, 14)
}

//https://wiki.openstreetmap.org/wiki/Slippy_map_tilenames
// flipping xyz to tms.
// y = (1<<zoom) - y - 1

/*
https://github.com/mapbox/mbtiles-spec/blob/master/1.3/spec.md
Note that in the TMS tiling scheme, the Y axis is reversed from the "XYZ" coordinate system commonly used in the URLs to request individual tiles, so the tile commonly referred to as 11/327/791 is inserted as zoom_level 11, tile_column 327, and tile_row 1256, since 1256 is 2^11 - 1 - 791.

Given Tile numbers to longitude/latitude :

n = 2 ^ zoom
lon_deg = xtile / n * 360.0 - 180.0
lat_rad = arctan(sinh(π * (1 - 2 * ytile / n)))
lat_deg = lat_rad * 180.0 / π

n = 2 ^ zoom
xtile = n * ((lon_deg + 180) / 360)
ytile = n * (1 - (log(tan(lat_rad) + sec(lat_rad)) / π)) / 2
*/
