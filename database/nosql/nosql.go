//Package nosql makes from SQL database to NoSQL
//Database must support json type
package nosql

import (
	"database/sql"
	"encoding/json"
)

// NoSQL object
type NoSQL struct {
	openMethod func() (*sql.Tx, error)
}

// Query execute request to database with command and paramerers map
func (_self *NoSQL) Query(command string, data interface{}) (res interface{}, err error) {
	// convert params to json
	var raw []byte
	if raw, err = json.Marshal(data); err != nil {
		return
	}
	// open tx
	var tx *sql.Tx
	if tx, err = _self.openMethod(); err != nil {
		return
	}
}
