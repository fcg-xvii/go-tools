//Package nosql makes from SQL database to NoSQL.
//Source database must support json type
package nosql

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

// OpenTXMethod is callback functon to open transaction in NoSQL object
type OpenTXMethod func() (*sql.Tx, error)

// NoSQL object
type NoSQL struct {
	openMethod OpenTXMethod
}

// New is NoSQL object constructor
func New(openMethod OpenTXMethod) *NoSQL {
	return &NoSQL{openMethod}
}

// Query execute request to database with (function(json) json) interface.
// Function should return json value
func (_self *NoSQL) Query(function string, data interface{}) (res interface{}, err error) {
	// convert data to raw json
	var raw []byte
	if raw, err = json.Marshal(data); err != nil {
		return
	}
	// open tx
	var tx *sql.Tx
	if tx, err = _self.openMethod(); err == nil {
		// execute query and scan result
		var rRaw []byte
		row := tx.QueryRow(fmt.Sprintf("select * from %v($1)", function), raw)
		if err = row.Scan(&rRaw); err != nil {
			tx.Rollback()
			return
		}
		// convert raw result data to json and return
		if err = json.Unmarshal(rRaw, &res); err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}
	return
}
