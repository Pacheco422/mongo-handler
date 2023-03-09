// Implements a Database connection with mongo
package dbhandler

import (
	"net/http"
	"encoding/json"
	"go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/mongo/readpref"
	"io"
	"context"
	"time"
	log "github.com/sirupsen/logrus"
)

type DBHandler interface {
	// insert(p Product) (int, error)
	// update(p Product) (error)
	// delete(id int) (error)
	// fetch(id []int) ([]Product, error)
	// connect(connstr string, r int, tconn int, tcdown int)
	// TODO: search(str ) (err error)
	ConnectWithHeartbeat(s http.Server, connstr string, r int, tconn int, tcdown int, dbhb int)
}

type MongoHandler struct {
	client *mongo.Client // mongo client
	dbname string // database name
	collname string // database collection
}

type Product struct {
	Id string
	Name string
	Price string
	Date string
	Tags string
}


func NewMongoHandler(dbname string, collname string) *MongoHandler {
	return &MongoHandler{dbname: dbname, collname: collname}
}

func (p *Product) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

func (p *Product) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

// Connect tries to estabilish a client with mongodb by using connection string
// connstr. It tries to estabilish connection for another r times, if it is not 
// successful on its first attempt. For each try, it has tconn seconds to connect. 
// For each time the connection attempt fails, it sleeps for a cooldown of tcdown 
// seconds .
func (dbh *MongoHandler) Connect(connstr string, r int, tconn int, tcdown int) error {
	var err error
	// Tries (and retries) to connect
	for i := 0; i <= r; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(tconn)*time.Second)
		defer cancel()
		log.Infof("Trying to connect to database (connection has %d seconds to succeed) ...", tconn)
		dbh.client, err = mongo.Connect(ctx, options.Client().ApplyURI(connstr))
		if err != nil {
			log.Errorf("Could not connect to database using connection string '%s' (%d/%d).", connstr, i, r)
			log.Infof("Retrying connection in %d seconds .", tcdown)
			time.Sleep(time.Duration(tcdown)*time.Second)
		} else if err := dbh.client.Ping(ctx, readpref.Primary()); err != nil {
			log.Errorf("Could not connect to database using connection string '%s' (%d/%d).", connstr, i, r)
			log.Infof("Retrying connection in '%d' seconds .", tconn)
			time.Sleep(time.Duration(tcdown)*time.Second)
		} else {
			log.Info("Connected to database .")
			return nil
		}
	}
	log.Error("Reached maximum number of attemps connecting to database .")
	return err
}


// Connects with a db, with a heartbeat and reconnect. Every dbhb secons, a heartbeat is sent
// If the db does not respond, retries connection.
func (dbh *MongoHandler) ConnectWithHeartbeat(s http.Server, connstr string, r int, tconn int, tcdown int, dbhb int) {
	for range time.Tick(time.Second * time.Duration(dbhb)) {
		log.Debug("Sent heartbeat to mongo connection .")
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(tconn)*time.Second)
		defer cancel()
		if err := dbh.client.Ping(ctx, readpref.Primary()); err != nil {
			log.Info("Lost database connection. Connection will be retried.")
			if err := dbh.Connect(connstr, r, tconn, tcdown); err != nil {
				log.Fatal("Application is shutting down ...")
			}
		}
		log.Debug("Application is ready .")
	}
}
