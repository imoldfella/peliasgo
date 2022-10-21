import './style.css'
import { Map as MapGl, AttributionControl, NavigationControl, ResponseCallback, RequestParameters, Cancelable } from 'maplibre-gl'
import maplibre from 'maplibre-gl'
import 'maplibre-gl/dist/maplibre-gl.css';
import { mapstyle } from "./mapstyle"
import { Db, TableReader } from './db'
import { Pyramid } from './hilbert';
import { decompressSync } from "fflate";

let pyr = new Pyramid(0, 15)

let db: Db
let tr: TableReader
// function typedArrayToBuffer(array: Uint8Array): ArrayBuffer {
//   return array.buffer.slice(array.byteOffset, array.byteLength + array.byteOffset)
// }

function foo(
  params: RequestParameters,
  callback: ResponseCallback<any>): Cancelable {
  console.log("wants", params)

  if (params.type == "json") {

    const tilejson = {
      tiles: [params.url + "/{z}/{x}/{y}"],
      "scheme": "tms",
      minzoom: 0,
      maxzoom: 14,
    };
    callback(null, tilejson, null, null);
  } else {
    const re = new RegExp(/dg:(.+)\/(\d+)\/(\d+)\/(\d+)/);
    const result = params.url.match(re);
    //const pmtiles_url = result[1];
    if (result) {
      const z = parseInt(result[2]);
      const x = parseInt(result[3]);
      const y = parseInt(result[4]);
      console.log("xyz", x, y, z)

      let id = pyr.FromXyz(x, y, z)
      tr.get(id).then((a) => {
        a = decompressSync(new Uint8Array(a))
        console.log("fetch", id, a)
        callback(undefined,a, undefined, undefined)
      })
    }
   }
    return {
      cancel: () => { },
    };
}

maplibre.addProtocol("dg", foo)




async function loadMetadata() {
  db = await Db.open("db")
  tr = await db.table("map")
  console.log(await tr.getJson(0))
  var map = new MapGl({
    container: 'map',
    style: mapstyle,
    center: [7.4, 43.7372],
    zoom: 0,
    attributionControl: false
  });

  map.addControl(new AttributionControl({
    compact: true,
    customAttribution: '<a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors; Date: 06.2022'
  }));

  map.addControl(new NavigationControl({}));
}
loadMetadata()




