
# this is a mess! cleaning it up is todo

# build a map of the world 

copied from planetiler github.

```
wget https://github.com/onthegomap/planetiler/releases/latest/download/planetiler.jar
java -Xmx110g -XX:MaxHeapFreeRatio=40   -jar planetiler.jar \
  --area=planet --bounds=planet --download  --download-threads=10 --download-chunk-size-mb=1000 --fetch-wikidata  --mbtiles=output.mbtiles   --nodemap-type=array --storage=ram 
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
