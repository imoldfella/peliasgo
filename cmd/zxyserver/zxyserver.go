package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var db *sql.DB
var tstmt *sql.Stmt
var meta = map[string]string{}

func main() {
}

func serveDir(path string, port string) error {
	var e error
	db, e = sql.Open("sqlite3", path+"/output.mbtiles")
	if e != nil {
		panic(e)
	}
	tstmt, e = db.Prepare("select tile_data from tiles where zoom_level=? and tile_column=? and tile_row=?")
	if e != nil {
		panic(e)
	}
	r := mux.NewRouter()
	r.Use(func(next http.Handler) http.Handler { return handlers.LoggingHandler(os.Stdout, next) })
	r.HandleFunc("/rlp/{z}/{x}/{y}.pbf", getTile).Methods(http.MethodGet, http.MethodPut, http.MethodPatch, http.MethodOptions)
	r.Use(mux.CORSMethodMiddleware(r))
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(path)))
	http.Handle("/", r)
	//r.PathPrefix("/").Handler(http.FileServer(http.Dir(path)))
	log.Printf("Serving %s on HTTP port: %s\n", path, port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
	return nil
}

func getMetadata(w http.ResponseWriter, r *http.Request) {
	rows, e := db.Query("select name,text from metadata")
	if e != nil {
		panic(e)
	}
	var key, value string
	for rows.Next() {
		rows.Scan(&key, &value)
		meta[key] = value
	}
	rows.Close()

}
func getTile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	x := vars["x"]
	y := vars["y"]
	z := vars["z"]

	log.Printf("z=%v,x=%v,y=%v", z, x, y)
	var data []byte
	e := tstmt.QueryRow(z, x, y).Scan(&data)
	if e == nil {
		w.Header().Add("Content-Type", "application/x-protobuf")
		w.Header().Add("Content-Encoding", "gzip")
		w.Write(data)
	} else {
		w.WriteHeader(204)
	}

}
