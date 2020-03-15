//Package nosql makes from SQL database to NoSQL
//Database must support json type
package nosql

import "database/sql"

// NoSQL object
type NoSQL struct {
	openMethod func() (sql.Tx, error)
}

// Query execute request to database with command and paramerers map
func (_self *NoSQL) Query(command string, params map[string]interface{}) (interface{}, error) {
	return
}
