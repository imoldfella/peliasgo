
let controller = new AbortController();
let signal = controller.signal;

export interface DbTable {
    root: number
    root_length: number
    height: number
}

export interface DbJson {
    chunk_size: number
    chunk_count: number
    table: {
        [key: string]: DbTable
    }
}
export function binarySearch(nums: number[], target: number): number {
    let left: number = 0
    let right: number = nums.length

    while (left < right) {
        const mid: number = Math.floor((left + right) / 2);
        if (nums[mid] === target) return mid;
        if (target < nums[mid]) right = mid;
        else left = mid + 1;
    }
    return left - 1
}
export class DbIndex {
    constructor(public key: number[], public offset: number[]) {

    }
}
  
// function toNum(low: number, high: number): number {
//     return (high >>> 0) * 0x100000000 + (low >>> 0);
// }

export class BytesReader {
    pos = 0
    constructor(public data: Uint8Array) {

    }

  readUvarint() {
    const MSB = 0x80
    , REST = 0x7F
      var res    = 0
      , shift  = 0
      , b
    do {
      b = this.data[this.pos++]
      res += shift < 28
        ? (b & REST) << shift
        : (b & REST) * Math.pow(2, shift)
      shift += 7
    } while (b >= MSB)
    return res
  }
}
export class TableReader {
    constructor(public db: Db, public t: DbTable, public root: DbIndex) {

    }

    async read1(id: number): Promise<Uint8Array> {
        // we need to binary search the sorted values, then take the pivot
        // when we get to the leaf we can return the slice.
        var r = this.root

        let found = 0
        for (let i = 0; i < this.t.height; i++) {
            // find's smallest i such that fn is true; return <= id
            let found = binarySearch(r.key, id)

            // this may find an id that's smaller, but may work because of run length. we may want to keep run length when we decompress?

            const start = r.offset[found]
            const end = r.offset[found + 1]
            r = await this.db.readIndex(start, end - start)
        }
        found = binarySearch(r.key, id)
        const start = r.offset[found]
        const end = r.offset[found + 1]
        return await this.db.readBytes(start, end - start)
    }
    async getJson(id: number): Promise<string> {
        const d = await this.read1(id)
        return new TextDecoder().decode(d)
    }

    async get(id: number): Promise<Uint8Array > {
        const d = await this.read1(id)
        if (d.length <= 5) {
            const id2 = new BytesReader(d).readUvarint()
            return this.read1(id2)
        } else {
            return d
        }
    }
    
}
export interface RangeResponse {
    data: ArrayBuffer;
    etag?: string;
    expires?: string;
    cacheControl?: string;
}

function append(a: Uint8Array, b: Uint8Array) { // a, b TypedArray of same type
    var c = new Uint8Array(a.length + b.length);
    c.set(a, 0);
    c.set(b, a.length);
    return c;
}

export class Db {
    index = new Map<number, DbIndex>()
    table_ = new Map<string, TableReader>()

    constructor(public path: string, public js: DbJson) {

    }

    async table(name: string): Promise<TableReader> {
        const r = this.table_.get(name)
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

    static async open(path: string): Promise<Db> {
        const o = await (await fetch(path + "/index.json")).json()
        return new Db(path, o as DbJson)
    }

    async slice(file: number, from: number, size: number): Promise<Uint8Array> {
        const resp = await fetch(this.path + "/" + file, {
            signal: signal,
            headers: { Range: "bytes=" + from + "-" + (from + size - 1) },
        });
        console.log()
        return new Uint8Array(await resp.arrayBuffer())
    }
    async readBytes(pos: number, size: number): Promise<Uint8Array> {
        const ch = this.js.chunk_size
        const file = Math.floor(pos / ch)
        const offset = pos - file * ch
        const avail = ch - offset

        if (avail >= size) {
            let o = await this.slice(file, offset, size)
            return o
        } else {   // split across two files.
            const b = await this.slice(file, offset, avail)
            const b2 = await this.slice(file + 1, 0, size - avail)
            return append(b, b2)
        }
    }

    async readIndex(pos: number, size: number): Promise<DbIndex> {
        const o = this.index.get(pos)
        if (o) {
            return o
        }

        const b = await this.readBytes(pos, size)
        const rdr = new BytesReader(b)
        let ln = rdr.readUvarint()

        const key: number[] = []
        const offset: number[] = []
        for (let i = 0; i < ln; i++) {
            key[i] = rdr.readUvarint()
        }
        for (let i = 0; i < ln + 1; i++) {
            offset[i] = rdr.readUvarint()
        }
        const idx = new DbIndex(key, offset)
        this.index.set(pos, idx)
        return idx
    }

}