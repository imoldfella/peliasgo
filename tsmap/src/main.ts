import './style.css'
import  { Map as MapGl,AttributionControl, NavigationControl, ResponseCallback, RequestParameters, Cancelable}  from 'maplibre-gl'
import maplibre from 'maplibre-gl'
import 'maplibre-gl/dist/maplibre-gl.css';
import { Protocol } from './protocol'
import { mapstyle } from "./mapstyle"

let controller = new AbortController();
let signal = controller.signal;

let protocol = new Protocol();

interface DbTable {
      root: number
      root_length: number
      height: number
    }

interface DbJson  {
  chunk_size: number
  chunk_count: number
  table: {
    [key: string]: DbTable
  }
}
function binarySearch(nums: number[], target: number): number {
  let left: number = 0
  let right: number = nums.length 

  while (left <= right) {
    const mid: number = Math.floor((left + right) / 2);
    if (nums[mid] === target) return mid;
    if (target < nums[mid]) right = mid;
    else left = mid + 1;
  }
  return left-1
}
class DbIndex {
  constructor(public key : number[], public offset : number[]){

  }
}
class BytesReader {
  constructor(public data: Uint8Array, pos: number){

  }
  readUvarint() : number{

    return 0
  }
}
class TableReader {
  constructor(public db: Db, public t: DbTable, public root: DbIndex){

  }

	async read1(id: number) : Promise<Uint8Array> {
		// we need to binary search the sorted values, then take the pivot
		// when we get to the leaf we can return the slice.
		var r =  this.root

    let found = 0
		for (let i = 0; i < this.t.height; i++) {
			// find's smallest i such that fn is true; return <= id
			let found = binarySearch(r.key, id)
			
			// this may find an id that's smaller, but may work because of run length. we may want to keep run length when we decompress?

			const start = r.offset[found]
			const end = r.offset[found+1]
			r = await this.db.readIndex(start, end-start)
		}
		found = binarySearch(r.key, id)
		const start = r.offset[found]
		const end = r.offset[found+1]
		return this.db.readBytes(start, end-start)
	}
  

  async get(id: number) : Promise<Uint8Array|undefined>{
      const d = await this.read1(id)
      if (d.length <= 5) {
        const id2 = readUvarint(bytes.NewReader(d))
        if e != nil {
          return nil, e
        }
        return read1(uint32(id2))
      } else {
        return d
      }
    }
}
  

function append(a: Uint8Array, b: Uint8Array) { // a, b TypedArray of same type
  var c = new Uint8Array(a.length + b.length);
  c.set(a, 0);
  c.set(b, a.length);
  return c;
}

class Db {
  index = new Map<number, DbIndex>()
  table_ = new Map<string, TableReader>()

  constructor(public path: string, public  js: DbJson){

  }

  async table(name: string) : Promise<TableReader> {
    const r =  this.table_.get(name)
    if (r) {
      return r
    } {
      const t = this.js.table[name]
      const root = await this.readIndex(t.root, t.root_length)
      const o = new TableReader(this, t, root)
      this.table_.set(name, o)
      return o
    }
  }
  
  static async  open(path: string):Promise<Db> {
    const o =  await (await fetch(path + "/index.json")).json()
    return new Db(path, o as DbJson)
  }

  async  slice(file: number, from: number, size: number ) : Promise<Uint8Array> {
    const resp = await fetch(this.path + "/"+file, {
      signal: signal,
      headers: { Range: "bytes=" + from + "-" + (from + size - 1) },
    }); 
    return new Uint8Array(await resp.arrayBuffer())
  }
  async readBytes( pos: number, size: number) : Promise<Uint8Array> {
    const ch = this.js.chunk_size
    const file = Math.floor(pos / ch)
    const offset = pos - file*ch
    const avail = ch - offset
  
    if (avail >= size) {
      // split across two files.
      return this.slice(file, offset, size)
    } else {
      const b = await this.slice(file, offset, avail)
      const b2 = await this.slice(file+1, 0, size-avail)
      return append(b, b2)
    }
  }

   async readIndex(pos: number, size: number) : Promise<DbIndex> {
    const o = this.index.get(pos)
    if (o) {
      return o
    }
  
    const b  = await this.readBytes(pos, size)
    const key : number[] = []
    const offset : number[] = []
    const idx = new DbIndex(key, offset)

    this.index.set(pos, idx)
    return idx
  }
  
}

function foo(
  params: RequestParameters, 
  callback: ResponseCallback<any> ): Cancelable  {
  console.log("wants", params)
  
  if (params.type == "json") {
    // what do we return here? not the style, maybe the metadata?
    // there 
    const tilejson = {
      tiles: [params.url + "/{z}/{x}/{y}"],
      minzoom: 0,
      maxzoom: 15,
    };
    callback(null, tilejson, null, null);
  } else {


   // when is this not json?
  }
  //let instance = this.tiles.get(pmtiles_url);
  //if (!instance) {
    // instance = new PMTiles(pmtiles_url);
    // this.tiles.set(pmtiles_url, instance);
  //}

  // instance.getHeader().then((h) => {
  //   const tilejson = {
  //     tiles: [params.url + "/{z}/{x}/{y}"],
  //     minzoom: h.minZoom,
  //     maxzoom: h.maxZoom,
  //   };
  //   callback(null, tilejson, null, null);
  // });

  return {
    cancel: () => {},
  };
}

maplibre.addProtocol("dg",foo)



async function loadMetadata() {
    await Db.open("db")
  var map = new MapGl({
    container: 'map',
    style: mapstyle,
    center: [7.4,43.7372],
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




