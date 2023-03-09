// Implements fetch functionality
package dbhandler

import (
	log "github.com/sirupsen/logrus"
)


// Handles Fetch
func (dbh *MongoHandler) FindOneById(id string) (Product, error) {
	log.Debug("Entered Find function .")
	log.Debugf("Searching for element with id '%s' inside database '%s', collection '%s' ...", id, dbh.dbname, dbh.collname)

	// Setup for search
	product := Product{}

	return product, nil
}
