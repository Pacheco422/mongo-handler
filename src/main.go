// Application for interacting with a mongo database .
package main

import (
	"fmt"
	"os"
	"context"
	"time"
	"flag"
	"os/signal"
	"net/http"
	"mongo-handler/probes"
	"mongo-handler/api"
	"mongo-handler/dbhandler"
	log "github.com/sirupsen/logrus"
)

// API variables
var host = "0.0.0.0" // application connection host
var port = 9000 // application connection port
var rto = 10 // maximum time to read a request
var wto = 15 // maximum time allowed for write a response
var ito = 360 // maximum time idle
var dbhb = 10 // database heartbeat
var dbname = "testingdb" // database name
var collname = "colldb" // collection name

var connstr = "localhost:27017" // database connection host
var r = 5 // number of retries permitted when connecting to database
var tconn = 10 // time allowed for each connection attempt
var tcdown = 10 // time between connection attempts

var rh *probes.ReadinessHandler
var dbh *dbhandler.MongoHandler

// Application variables
var grace = 30 // grace period for application shutdown measured in seconds


// Starts application flow, by aggregating some of basic functions required to 
// configure it. Currently, it changes the structure of logs and 
// specifies flags/arguments for the application .
func init() {
	// Logging
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
	
	// Arguments
	flag.Usage = func() {	// Redefines usage, which will show when user calls for help .
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
	var help = flag.Bool("help", false, "Show usage .")
	flag.StringVar(&host, "host", host, "Application api host .")
	flag.IntVar(&port, "port", port, "Application api port .")
	flag.IntVar(&rto, "read-timeout", rto, "Application maximum time to read a request .")
	flag.IntVar(&wto, "write-timeout", wto, "Application maximum time allowed for write a response .")
	flag.IntVar(&ito, "idle-timeout", ito, "Application maximum time idle .")
	flag.StringVar(&connstr, "connection-str", connstr, "Database connection string .")
	flag.IntVar(&dbhb, "db-heartbeat", dbhb, "Database heartbeat .")
	flag.StringVar(&dbname, "db-name", dbname, "Database Name .")
	flag.StringVar(&collname, "coll-name", collname, "Collection Name .")
	flag.IntVar(&r, "retries", r, "Number of retries when trying to connect to database .")
	flag.IntVar(&tconn, "time-connect", tconn, "Time (in seconds) allowed when trying to connect to database .")
	flag.IntVar(&tcdown, "time-cooldown", tcdown, "Time (in seconds) between tries when trying to connect to database .")
	flag.IntVar(&grace, "grace-period", grace, "Grace period (in seconds) when shutting down application .")
	flag.Parse() // Parses flags
	if *help { // Shows usage if user calls for help
		flag.Usage()
		os.Exit(0)
	}
}
	

// Application flow. It includes setup, connection to database,
// instantiation of mux with handlers for interacting with db and
// probes (readiness/liveness checks) and application teardown .
func main() {
	log.Info("Started application flow . ")

	rh = probes.NewReadinessHandler()
	dbh = dbhandler.NewMongoHandler(dbname, collname)

	log.Info("Connecting to Server ...")
	dbh.Connect(connstr, r, tconn, tcdown)

	log.Info("Creating Server ...")
	s := startServer()
	log.Info("Server created .")

	log.Infof("Connecting to database (%s) ...", connstr)
	go dbh.ConnectWithHeartbeat(s, connstr, r, tconn, tcdown, dbhb)
	// Alter state of readiness probe
	log.Info("Application is Ready .")
	rh.Ready = true

	// Detects os signals demanding to shut down the application .
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	log.Infof("Detected signal %s. Commencing application Teardown ...", sig)
	teardown(s)
}

func startServer() http.Server {

	// Create ServeMux
	sm := http.NewServeMux()
	// Register Liveness and Readiness probes
	sm.Handle("/alive", probes.NewLivenessHandler())
	sm.Handle("/ready", rh)
	// Register remaining handlers
	sm.Handle("/insert", api.HandleInsert(dbh))
	sm.Handle("/findById", api.HandleFindById(dbh))

	// Create Server
	ba := fmt.Sprintf("%s:%d", host, port)
	s := http.Server{
		Addr:         ba,     			 				// bind address
		Handler:      sm,                				// set the default handler
		ReadTimeout:  time.Duration(rto) * time.Second, // max time to read request from the client
		WriteTimeout: time.Duration(wto) * time.Second, // max time to write response to the client
		IdleTimeout:  time.Duration(ito) * time.Second, // max time for connections using TCP Keep-Alive
	}

	go func() {
		err := s.ListenAndServe()
		if err != nil {
			log.Errorf("Error starting server (%s) .", ba)
			os.Exit(1)
		}
		log.Info("Server stopped .")
	}()

	return s
}

// Teardown tries to clean things up and exit application gracefully .
func teardown(s http.Server) {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(grace)*time.Second)
	s.Shutdown(ctx)
}
