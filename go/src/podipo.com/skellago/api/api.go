package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/codegangsta/negroni"
	"github.com/coocood/qbs"
	_ "github.com/lib/pq" // Needed to make the postgres driver available

	"podipo.com/skellago/be"
)

var logger = log.New(os.Stdout, "[api] ", 0)

func registerDB() error {
	// Register the QBS db
	dsn := &qbs.DataSourceName{
		DbName:   os.Getenv("POSTGRES_DB_NAME"),
		Username: os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Host:     os.Getenv("POSTGRES_PORT_5432_TCP_ADDR"),
		Port:     os.Getenv("POSTGRES_PORT_5432_TCP_PORT"),
		Dialect:  qbs.NewPostgres(),
	}
	dsn.Append("sslmode", "disable")
	qbs.RegisterWithDataSourceName(dsn)
	return nil
}

func main() {
	err := registerDB()
	if err != nil {
		logger.Print("DB Registration Error: " + err.Error())
		return
	}
	err = be.MigrateDB()
	if err != nil {
		logger.Print("DB Migration Error: " + err.Error())
		return
	}

	port, err := strconv.ParseInt(os.Getenv("PORT"), 10, 64)
	if err != nil {
		port = 9000
	}
	logger.Print("Port ", port)

	server := negroni.New()

	staticDir := os.Getenv("STATIC_DIR")
	if staticDir == "" {
		staticDir = "static"
	}
	logger.Print("Static dir ", staticDir)
	static := negroni.NewStatic(http.Dir(staticDir))
	static.Prefix = "/api/static"
	server.Use(static)

	api := be.NewAPI("/api")
	server.UseHandler(api.Mux)

	//store := sessions.NewCookieStore([]byte("1234abcd"))
	//server.Use(sessions.Sessions("shoe_session", store))

	server.Run(":" + strconv.FormatInt(port, 10))
}
