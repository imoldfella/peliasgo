

running
wget https://github.com/onthegomap/planetiler/releases/latest/download/planetiler.jar
java -Xmx32g -jar planetiler.jar --osm-path=sources/north_america_latest.osm.pbf

#testing

npm install -g tileserver-gl-light
tileserver-gl-light --mbtiles data/output.mbtiles

we need a way to flatten the tiles, maybe combine the tiles.
we might want a way to limit to a single zoom=10?
currently mbtiles is going to have all the zoom levels.

https://github.com/mapbox/mbtiles-spec/blob/master/1.3/spec.md
Note that in the TMS tiling scheme, the Y axis is reversed from the "XYZ" coordinate system commonly used in the URLs to request individual tiles, so the tile commonly referred to as 11/327/791 is inserted as zoom_level 11, tile_column 327, and tile_row 1256, since 1256 is 2^11 - 1 - 791.

https://github.com/mapbox/tile-cover

(main.Xy) {
 x: (int) 2509,
 y: (int) 9337
}
(main.Xy) {
 x: (int) 5147,
 y: (int) 10784
}

(main.Xy) {
 X: (int) 2496,
 Y: (int) 9328
}
(main.Xy) {
 X: (int) 5152,
 Y: (int) 10784
}