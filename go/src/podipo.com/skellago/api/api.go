package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/codegangsta/negroni"

	"podipo.com/skellago/be"
)

var logger = log.New(os.Stdout, "[api] ", 0)

func main() {

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
