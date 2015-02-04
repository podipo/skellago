package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/codegangsta/negroni"
	"github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/cookiestore"
	"podipo.com/skellago/be"
)

// VERSION is the API version
var VERSION = "0.1.0"

var logger = log.New(os.Stdout, "[api] ", 0)

func main() {
	// Get the required environment variables
	port, err := strconv.ParseInt(os.Getenv("PORT"), 10, 64)
	if err != nil {
		logger.Panic("No PORT env variable")
		return
	}
	staticDir := os.Getenv("STATIC_DIR")
	if staticDir == "" {
		logger.Panic("No STATIC_DIR env variable")
		return
	}
	fsDir := os.Getenv("FILE_STORAGE_DIR")
	if fsDir == "" {
		logger.Panic("No FILE_STORAGE_DIR env variable")
		return
	}
	sessionSecret := os.Getenv("SESSION_SECRET")
	if sessionSecret == "" {
		logger.Panic("No SESSION_SECRET env variable")
		return
	}
	frontEndDir := os.Getenv("FRONT_END_DIR") // Optional

	logger.Print("PORT:\t\t", port)
	logger.Print("STATIC_DIR:\t", staticDir)
	logger.Print("FRONT_END_DIR:\t", frontEndDir)
	logger.Print("FILE_STORAGE_DIR:\t", fsDir)

	if be.DBHost == "" { // Try CoreOS etcd
		etcdJSON, err := be.EtcdGet(os.Getenv("COREOS_PRIVATE_IPV4"), "/v2/keys/services/skella/postgres")
		if err != nil {
			logger.Panic("Could not parse the etcd data: " + err.Error())
			return
		}
		logger.Print("Received etcd json: " + etcdJSON)
		var postgresData EtcPostgresData
		err = json.NewDecoder(strings.NewReader(etcdJSON)).Decode(&postgresData)
		if err != nil {
			logger.Panic("Could not parse the etcd data: " + etcdJSON)
			return
		}
		be.DBHost = postgresData.Host
		be.DBPort = strconv.Itoa(postgresData.Port)
	}
	logger.Print("DB host: ", be.DBHost, ":", be.DBPort)

	err = be.InitDB()
	if err != nil {
		logger.Panic("DB Registration Error: " + err.Error())
		return
	}
	err = migrateDB()
	if err != nil {
		logger.Panic("DB Migration Error: " + err.Error())
		return
	}

	fs, err := be.NewLocalFileStorage(fsDir)
	if err != nil {
		logger.Panic("Could not open file storage directory: " + fsDir)
		return
	}

	server := negroni.New()
	store := cookiestore.New([]byte(sessionSecret))
	server.Use(sessions.Sessions(be.AuthCookieName, store))

	if frontEndDir != "" {
		feStatic := negroni.NewStatic(http.Dir(frontEndDir))
		feStatic.Prefix = ""
		server.Use(feStatic)
	}

	static := negroni.NewStatic(http.Dir(staticDir))
	static.Prefix = "/api/static"
	server.Use(static)

	api := be.NewAPI("/api/"+VERSION, VERSION, fs)
	api.AddResource(NewEchoResource(), true)

	server.UseHandler(api.Mux)
	server.Run(":" + strconv.FormatInt(port, 10))
}

type EtcPostgresData struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}
