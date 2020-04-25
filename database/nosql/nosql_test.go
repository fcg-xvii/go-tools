package nosql

import (
	"database/sql"
	"io/ioutil"
	"log"
	"testing"

	_ "github.com/lib/pq"
)

var (
	dbConn string
)

func init() {
	// read postgres connection string from file z_data.config
	connSource, _ := ioutil.ReadFile("z_data.config")
	dbConn = string(connSource)
	log.Println("DB connection string", dbConn)
}

func TestNoSQL(t *testing.T) {
	if db, err := sql.Open("postgres", dbConn); err == nil {
		ex := New(func() (*sql.Tx, error) {
			return db.Begin()
		})
		res, err := ex.Call("public.arr_count", map[string]interface{}{
			"input": []int{0, 1, 2},
		})
		t.Log(res, err)
	} else {
		t.Error(err)
	}
}
