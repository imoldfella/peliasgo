package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type PostgisCoder struct {
	server string
	db     *sql.DB
	stmt   *sql.Stmt
}

func NewPostgis(server string) (Geocoder, error) {

	db, err := sql.Open("postgres", server)
	if err != nil {
		return nil, err
	}
	query := "SELECT ST_X(g.geomout) As lon, ST_Y(g.geomout) As lat FROM tiger.geocode($1,1) As g"
	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}

	return &PostgisCoder{
		server: server,
		db:     db,
		stmt:   stmt,
	}, nil
}

// Geocode implements Geocoder
func (p *PostgisCoder) Geocode(zip string, address string) (float64, float64, error) {
	var lat, lon float64
	err := p.stmt.QueryRow(address+","+zip).Scan(&lat, &lon)
	return lat, lon, err
}

var _ Geocoder = (*PostgisCoder)(nil)
