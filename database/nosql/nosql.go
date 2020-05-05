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

// CallJSON accepts raw json bytes and returns result raw json bytes
func (_self *NoSQL) CallJSON(function string, rawJSON []byte) (resRawJSON []byte, err error) {
	// open tx
	var tx *sql.Tx
	if tx, err = _self.openMethod(); err == nil {
		// execute query and scan result
		row := tx.QueryRow(fmt.Sprintf("select * from %v($1)", function), rawJSON)
		if err = row.Scan(&resRawJSON); err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}
	return
}

// CallObjParam accepts interface{} object and returns result raw json bytes
func (_self *NoSQL) CallObjParam(function string, data interface{}) (resRawJSON []byte, err error) {
	// convert incoming object to raw json
	var raw []byte
	if raw, err = json.Marshal(data); err != nil {
		return
	}
	return _self.CallJSON(function, raw)
}

// Call accepts interface{} object and returns result interface{}
func (_self *NoSQL) Call(function string, data interface{}) (res interface{}, err error) {
	var resRawJSON []byte
	if resRawJSON, err = _self.CallObjParam(function, data); err == nil {
		// convert raw result data to obj
		err = json.Unmarshal(resRawJSON, &res)
	}
	return
}
