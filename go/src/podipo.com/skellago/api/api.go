package main

import (
	"fmt"
	"github.com/codegangsta/negroni"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "You have reached the API")
	})

	n := negroni.Classic()
	n.UseHandler(mux)
	n.Run(":9000")
}
