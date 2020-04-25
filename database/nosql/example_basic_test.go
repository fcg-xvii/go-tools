package nosql_test

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/fcg-xvii/go-tools/database/nosql"
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

func Example_basic() {
	// As a database will use postgres and plv8 (https://plv8.github.io/)
	// On the database side, you must create a function that counts the number of elements in the input array
	// As input argument function accept json object { "input": array } and will return object { "input": array, output: count }
	//
	// CREATE OR REPLACE FUNCTION public.arr_count(data jsonb)
	// RETURNS jsonb AS
	// $BODY$
	//   if(typeof(data.input) != 'object' || !data.input.length) {
	//	   plv8.elog(ERROR, 'Incoming data must be array')
	//   }
	//   data.output = data.input.length
	//   return data
	// $BODY$
	// LANGUAGE plv8 IMMUTABLE STRICT

	// dbConn read from config file before
	// example dbConn string: postgres://postgres:postgres@127.0.0.1/postgres?sslmode=disable&port=5432

	// open the postgres database
	db, err := sql.Open("postgres", dbConn)
	if err != nil {
		fmt.Println(err)
		return
	}

	// setup open new database transaction method
	openTX := func() (*sql.Tx, error) {
		return db.Begin()
	}

	// create api object
	api := nosql.New(openTX)

	// setup request data
	data := map[string]interface{}{
		"input": []int{1, 2, 3},
	}

	// call the function on database side
	result, err := api.Call("public.arr_count", data)
	if err == nil {
		fmt.Println(err)
		return
	}

	// request completed
	fmt.Println(result)
}
