package main

import (
	"flag"
	"log"
	"mime"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gorilla/mux"
)

var addr = flag.String("addr", ":8182", "the address to listen on")

func main() {
	flag.Parse()

	// Had to do this because returns svg as text/xml when running on AppEngine: http://goo.gl/hwZSp2
	mime.AddExtensionType(".svg", "image/svg+xml")

	r := mux.NewRouter()
	apiServer, err := url.Parse("http://localhost:8183")
	if err != nil {
		log.Fatal(err)
	}
	r.Handle("/api/{rest:.*}", httputil.NewSingleHostReverseProxy(apiServer))
	r.HandleFunc("/{rest:.*}", staticHandler)
	http.Handle("/", r)
	log.Printf("Listening on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("path:", r.URL.Path)
	if r.URL.Path == "/" {
		http.ServeFile(w, r, "static/app.html")
	} else {
		http.ServeFile(w, r, "static/"+r.URL.Path)
	}
}
