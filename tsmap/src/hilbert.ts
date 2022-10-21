

// Hilbert represents a 2D Hilbert space of order N for mapping to and from.
// Implements SpaceFilling interface.
export class Hilbert {
    N: number = 0


    // NewHilbert returns a Hilbert space which maps integers to and from the curve.
    // n must be a power of two.
    constructor(n: number) {
        this.N = 1 << n
    }

    // Map transforms a one dimension value, t, in the range [0, n^2-1] to coordinates on the Hilbert
    // curve in the two-dimension space, where x and y are within [0,n-1].
    toXy(t: number): [number,number] {
        let x: number = 0, y: number = 0
        for (let i = 1; i < this.N; i = i * 2) {
            let rx = (t & 2) == 2
            let ry = (t & 1) == 1
            if (rx) {
                ry = !ry
            }

            [x, y] = this.rotate(i, x, y, rx, ry)

            if (rx) {
                x = x + i
            }
            if (ry) {
                y = y + i
            }

            t /= 4
        }

        return [x, y]
    }
    rotate(n: number, x: number, y: number, rx: boolean, ry: boolean): number[] {
        if (!ry) {
            if (rx) {
                x = n - 1 - x
                y = n - 1 - y
            }

            [x, y] = [y, x]
        }
        return [x, y]
    }
    // MapInverse transform coordinates on Hilbert curve from (x,y) to t.
    mapInverse(x: number, y: number): number {
        if (x < 0 || x >= this.N || y < 0 || y >= this.N) {
            throw "bad params"
        }

        let t: number = 0
        for (let i = Math.floor(this.N / 2); i > 0;i =  Math.floor(i / 2)) {
            let rx = (x & i) > 0
            let ry = (y & i) > 0

            let o = (rx ? 3 : 0) ^ (ry ? 1 : 0)
            t += i * i * o

            let v = this.rotate(i, x, y, rx, ry)
            x = v[0]
            y = v[1]
        }

        return t
    }

    // rotate rotates and flips the quadrant appropriately.


}


export class Pyramid {
    Len: number = 0
    MinZoom: number = 0
    MaxZoom: number = 0
    h: Hilbert[] = []// one fore each zoom
    start: number[] = []

    constructor(minZoom: number, maxZoom: number) {
        this. MinZoom =  minZoom
        this.MaxZoom = maxZoom
        let nh = maxZoom - minZoom
        let cnt = 0
        this.start[0]=0
        for (let x=0; x< nh ; x++) {
            this.h[x] = new Hilbert(x + minZoom)
            cnt += this.h[x].N * this.h[x].N
            this.start[x + 1] = cnt
        }
        this.Len = cnt
    }

    Xyz(id: number): [number, number, number] {
        id-- // make it 0 based
        for(let z=0; z<this.h.length; z++) {
            let h = this.h[z]
            let sz = h.N * h.N
            if (id >= sz) {
                id -= sz
            } else {
                let  [x, y] = h.toXy(id)
                return [x, y, z + this.MinZoom]
            }
        }
        return [0, 0, 0]
    }

    FromXyz(x: number, y: number, z: number): number {
        let id = this.h[z].mapInverse(x, y)
        id += this.start[z]
        return id + 1
    }


}