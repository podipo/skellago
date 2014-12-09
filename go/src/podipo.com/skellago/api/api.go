package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/codegangsta/negroni"
	"github.com/goincremental/negroni-sessions"
	"podipo.com/skellago/be"
)

var logger = log.New(os.Stdout, "[api] ", 0)

func main() {
	err := be.InitDB()
	if err != nil {
		logger.Print("DB Registration Error: " + err.Error())
		return
	}

	port, err := strconv.ParseInt(os.Getenv("PORT"), 10, 64)
	if err != nil {
		port = 9000
	}
	logger.Print("Port ", port)

	server := negroni.New()

	store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))
	server.Use(sessions.Sessions("api_session", store))

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

	server.Run(":" + strconv.FormatInt(port, 10))
}
