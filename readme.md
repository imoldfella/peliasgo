
# this is a mess! cleaning it up is todo

# build a map of the world 

copied from planetiler github.

```
wget https://github.com/onthegomap/planetiler/releases/latest/download/planetiler.jar
java -Xmx110g -XX:MaxHeapFreeRatio=40   -jar planetiler.jar \
  --area=planet --bounds=planet --download  --download-threads=10 --download-chunk-size-mb=1000 --fetch-wikidata  --mbtiles=output.mbtiles   --nodemap-type=array --storage=ram 

java -Xmx110g -XX:MaxHeapFreeRatio=40   -jar planetiler.jar \
  --osm-path=data/sources/monaco.osm.pbf --mbtiles=monaco.mbtiles   --nodemap-type=array --storage=ram 

```

## test the build 
npm install -g tileserver-gl-light
tileserver-gl-light --mbtiles data/output.mbtiles


# download npi data, todo




# openaddresses (not currently used)
https://batch.openaddresses.io/data

# links to standards
https://github.com/mapbox/mbtiles-spec/blob/master/1.3/spec.md


# potentially interesting libraries (not used)
https://github.com/murphy214/vector-tile-go
Uses custom protobuf code for speed

https://github.com/kjhsoftware/us-state-polygons/blob/master/LPZStatePolygons.m
state polygons

https://github.com/mapbox/tile-cover
cover a polygon with tiles.

https://github.com/omniscale/imposm3
golang for importing pbf to postgres

https://daylightmap.org/
facebook, microsoft enhanced maps.


# COG future work
cogger - cloud optimized geotiff. great for satellite images if we can get them

lansat geotiff images
https://d9-wret.s3.us-west-2.amazonaws.com/assets/palladium/production/s3fs-public/atoms/files/LSDS-1388-Landsat-Cloud-Optimized-GeoTIFF_DFCB-v2.0.pdf

cog standard
https://portal.ogc.org/files/102116

https://openlayers.org/


random merged notes

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

# whyu bvinc/go-sqlite-lite 
https://turriate.com/articles/making-sqlite-faster-in-go