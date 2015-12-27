package main

import (
	"flag"
	"log"
	"mime"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Email, Name, University, Location, TShirtSize, LinkedIn, GitHub, PersonalSite string
	TravelReimbursement                                                           int
	NeedsTravelReimbursement, Accepted, Waitlisted                                bool
}

var addr = flag.String("addr", ":8182", "the address to listen on")

func main() {
	flag.Parse()

	// Had to do this because returns svg as text/xml when running on AppEngine: http://goo.gl/hwZSp2
	mime.AddExtensionType(".svg", "image/svg+xml")
	db, err := gorm.Open("sqlite3", "/tmp/gorm.db")
	if err != nil {
		log.Fatal(err)
	}
	db.CreateTable(&User{})
	s := &server{db}

	r := mux.NewRouter()
	sr := r.PathPrefix("/api").Subrouter()
	sr.HandleFunc("/register", s.registerHandler)
	r.HandleFunc("/{rest:.*}", s.staticHandler)
	http.Handle("/", r)
	log.Printf("Listening on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

type server struct {
	db gorm.DB
}

func (s *server) staticHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("path:", r.URL.Path)
	if r.URL.Path == "/" {
		http.ServeFile(w, r, "static/app.html")
	} else {
		http.ServeFile(w, r, "static/"+r.URL.Path)
	}
}

func (s *server) registerHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("path:", r.URL.Path)
}
