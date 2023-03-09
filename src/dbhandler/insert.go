// Implements insert functionality
package dbhandler

import (
	log "github.com/sirupsen/logrus"
	"context"
	"time"
)


func (dbh *MongoHandler) InsertOne(product *Product) error {
	log.Debug("Entered InsertOne function .")
	log.Debug("Inserting product into database ...")
	
	// Setup for InsertOne
	collection := dbh.client.Database(dbh.dbname).Collection(dbh.collname)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Inserts One
	res, err := collection.InsertOne(ctx, product)
	if err != nil {
		log.Debug("Could not insert product into database .")
		return err
	}
	mid := res.InsertedID
	log.Infof("Inserted element with mongo id '%s', id '%s'", mid, product.Id)

	return nil
}
