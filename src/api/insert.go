// Implements insert functionality
package api

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"mongo-handler/dbhandler"
)


// Handles Insert. It inserts into database a product specified
// by the json payload inside the POST .
func HandleInsert(dbh *dbhandler.MongoHandler) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		log.Debug("Entered HandleInsert function .")
		if r.Method == "POST" {
			product := &dbhandler.Product{}
			err := product.FromJSON(r.Body)
			if err != nil {
				log.Errorf("Unable to Unmarshal JSON into Product struct: %s", err)
				http.Error(rw, "Unable to Unmarshal JSON into Product struct.", http.StatusBadRequest)
			}
			if err := dbh.InsertOne(product); err != nil {
				http.Error(rw, "Unable to save product into database .", http.StatusInternalServerError)
			}
		} else {
			errmsg := "Endpoint insert only supports POST method ."
			log.Error(errmsg)
			http.Error(rw, errmsg, http.StatusBadRequest)

		}
	}
}
