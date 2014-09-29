package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kahoona77/gotv/domain"
	"github.com/kahoona77/gotv/handler"
	"io/ioutil"
	"labix.org/v2/mgo"
	"log"
	"net/http"
	"strings"
)

func main() {
	// command line flags
	port := flag.Int("port", 8080, "port to serve on")
	flag.Parse()

	//creating db
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(notFound)

	fs := http.Dir("web")
	fileHandler := http.FileServer(fs)
	r.PathPrefix("/web/").Handler(http.StripPrefix("/web/", fileHandler))

	//Repositories
	serverRepo := domain.NewRepository(session, "servers")
	settingsRepo := domain.NewRepository(session, "settings")
	packetsRepo := domain.NewRepository(session, "packets")

	//Handlers
	dataHandler := handler.NewDataHandler(serverRepo, settingsRepo)
	r.PathPrefix("/data/").HandlerFunc(dataHandler.HandleRequests)

	packetsHandler := handler.NewPacketsHandler(packetsRepo)
	r.PathPrefix("/packets/").HandlerFunc(packetsHandler.HandleRequests)


	log.Printf("Running on port %d\n", *port)
	addr := fmt.Sprintf("127.0.0.1:%d", *port)
	// this call blocks -- the progam runs here forever
	http.Handle("/", r)
	err = http.ListenAndServe(addr, nil)
	fmt.Println(err.Error())
}

func notFound(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path
	if p == "/" || strings.HasPrefix(p, "/home") || strings.HasPrefix(p, "/search") || strings.HasPrefix(p, "/downloads") || strings.HasPrefix(p, "/logFile") {
		body, _ := ioutil.ReadFile("./web/index.html")
		fmt.Fprintf(w, string(body))
		return
	}
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "404 Not found.")
}
