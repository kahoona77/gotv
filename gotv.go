package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kahoona77/gotv/controller"
	"github.com/kahoona77/gotv/service"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	// command line flags
	port := flag.Int("port", 8080, "port to serve on")
	logFile := flag.String("log", "xtv.log", "log-file")
	flag.Parse()

	// setup log
	f, err := os.OpenFile(*logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(notFound)

	fs := http.Dir("web")
	fileController := http.FileServer(fs)
	r.PathPrefix("/web/").Handler(http.StripPrefix("/web/", fileController))

	// App Context
	ctx := service.CreateContext(*logFile)
	defer ctx.Close()

	//Controller
	dataController := controller.DataController {ctx}
	r.PathPrefix("/data/").HandlerFunc(dataController.HandleRequests)

	packetsController := controller.PacketsController {ctx}
	r.PathPrefix("/packets/").HandlerFunc(packetsController.HandleRequests)

	showsController := controller.ShowsController {ctx}
	r.PathPrefix("/shows/").HandlerFunc(showsController.HandleRequests)

	ircController := controller.IrcController {ctx}
	r.PathPrefix("/irc/").HandlerFunc(ircController.HandleRequests)

	downloadsController := controller.DownloadsController {ctx}
	r.PathPrefix("/downloads/").HandlerFunc(downloadsController.HandleRequests)

	log.Printf("XTV (Go) started port %d\n", *port)
	fmt.Printf("XTV (Go) started port %d\n", *port)
	addr := fmt.Sprintf(":%d", *port)
	// this call blocks -- the progam runs here forever
	http.Handle("/", r)
	err = http.ListenAndServe(addr, nil)
	log.Println(err.Error())
}

func notFound(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path
	if p == "/" || strings.HasPrefix(p, "/home") || strings.HasPrefix(p, "/search") || strings.HasPrefix(p, "/downloads") || strings.HasPrefix(p, "/logFile") || strings.HasPrefix(p, "/shows") {
		body, _ := ioutil.ReadFile("./web/index.html")
		fmt.Fprintf(w, string(body))
		return
	}
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "404 Not found.")
}
