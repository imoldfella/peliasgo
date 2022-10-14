package encode

// coded npi is 183mb, so easily fits in ram, does not need a high zoom level

// searching is not necessarily the same resolutions as the map
// 10 might be a good zoom for searching.
// so our "map" should have a very generic idea of what a tile is.
// in the data case, it's just a database shard.

// for npi we just want addresses, probably at zoom 10.
// to do this we need to assign all the points to a tile
/*
# contract lay

should I sort by address for rl compression? the reader can quickly build an npi map if that even matters.

# contract layer


# shared npi layer
  for each tile
	npi_provider(npi,fname,lname,sname,address_id)

alternative
# npi_provider(tile,npi,fname,lname,sname,address_id)
# is this worth an extra level of indirection? won't the addresses just come out in gzip anyway?
# does it help or hurt the client to use this?
# npi_address(tile,address_id,lat,lon,address1,address2,cszip_id)
# tile_cszip(tile, cszip_id, city, state, country)

# plan layer
# plan_provider(tile, contract, npi)

# plan_description (global, not per tile)
# plan_contract(code, contract, price)
*/

// tile sets may provide a symbol dictionary for compression (e.g water, land, dont care)
type Tileset interface {
	Pyramid() Pyramid
	// returns either the tile or a handle to a repeat.
	GetTile(index int, data []byte) ([]byte, int, error)
	Repeats() map[int][]byte
}

type AddressSet struct {
	// this is the coded address block. it has multiple tables in it.
	data  [][]byte
	shill []int
}

// There is only one repeat for the empty set.
func (*AddressSet) Repeats() map[int][]byte {
	return map[int][]byte{
		0: nil,
	}
}

// GetTile implements Tileset
func (a *AddressSet) GetTile(index int, data []byte) ([]byte, int, error) {
	panic("unimplemented")
}

// Pyramid implements Tileset
func (*AddressSet) Pyramid() Pyramid {
	panic("unimplemented")
}

var _ Tileset = (*AddressSet)(nil)

// we could try to compute a zoomlevel based on the addresses presented?
func NewAddressSet(path string, zoomlevel int) {
	// load duckdb here?
}

// this is going to read  a contract,npi csv and build a map layer.
// these can be national in size
// we need to join this with an npi->lat,lon to build the tiles.
// should these be tuples in a full text database? allowable option?
type DataLayer struct {
}

// We can have multiple tables, but only one table needs lat,lon
// we could potentially use things other than lat/lon to tile databases
// for example a contract id might be a way to build a database, but we can't shard this
// like can a long dense pyramid address. so what would that leave us with? a single shard?
// databases that have a shardable primary key (integer defined on 0...N)
// plan_contract(plan, code, contract, price) is not naturally dense. It could be sharded in other ways though. it may not need sharding at all.
// there maybe alternate ways to shard other than hardcore implicit interpolation that works for pyramid ids. For example the top leaf of a b+tree or prolly tree is explicit interpolation.
// this is effectively one shard. each database has an optional interpolation to a shard
// and then an option way to find the tuple (explicit or implicit key)
type Table struct {
	lat, lon []float64
	col      []Column
}
type Column struct {
}
type Schema struct {
	// an array schema has a shardable key
	array bool

	table []Table
}

func NewDataLayer(s *Schema) {

}
