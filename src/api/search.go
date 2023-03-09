// Implements search functionality
package api

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"mongo-handler/dbhandler"
)


// Handles Search
func HandleSearch(dbh *dbhandler.MongoHandler) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		log.Debug("Entered HandleInsert function .")
		if r.Method == "GET" {
			product := &dbhandler.Product{}
			err := product.FromJSON(r.Body)
			if err != nil {
				http.Error(rw, "Unable to Unmarshal JSON into Product struct .", http.StatusBadRequest)
			}
			if err := dbh.InsertOne(product); err != nil {
				http.Error(rw, "Unable to save product into database .", http.StatusInternalServerError)
			}
		} else {
			log.Error("Endpoint insert only supports POST method .")
		}
	}
}

// func validate(fl validator.FieldLevel) bool {
// 	// sku is of format abc-absd-dfsdf
// 	re := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)
// 	matches := re.FindAllString(fl.Field().String(), -1)

// 	if len(matches) != 1 {
// 		return false
// 	}

// 	return true
// }
