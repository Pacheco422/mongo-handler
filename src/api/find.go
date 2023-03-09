// Implements fetch functionality
package api

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"mongo-handler/dbhandler"
	"regexp"
)


// Handles FindById. It finds the product with id given on url .
func HandleFindById(dbh *dbhandler.MongoHandler) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var err error
		log.Debug("Entered HandleFindById function .")
		rgx := regexp.MustCompile(`/([0-9]+)`)
		id := rgx.FindAllStringSubmatch(r.URL.Path, -1)[0][1]

		log.Error(id)
		if r.Method == "GET" {
			var product dbhandler.Product
			if product, err = dbh.FindOneById(id); err != nil {
				http.Error(rw, "Unable to find product inside database .", http.StatusInternalServerError)
			}	
			if err = product.ToJSON(rw); err != nil {
				http.Error(rw, "Unable to marshal product that was found inside database .", http.StatusInternalServerError)
			}
		} else {
			errmsg := "Endpoint findById only supports GET method ."
			log.Error(errmsg)
			http.Error(rw, errmsg, http.StatusBadRequest)
		}
	}
}
