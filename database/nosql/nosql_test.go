package nosql

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
)

var (
	dbconn = "postgres://postgres:postgres@127.0.0.1/postgres?sslmode=disable&port=5432"
)

func TestNoSQL(t *testing.T) {
	if db, err := sql.Open("postgres", dbconn); err == nil {
		ex := New(func() (*sql.Tx, error) {
			return db.Begin()
		})
		res, err := ex.Query("public.tutor_arr_length", []int{0, 1, 2})
		t.Log(res, err)
	} else {
		t.Error(err)
	}
}
