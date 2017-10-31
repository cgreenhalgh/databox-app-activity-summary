package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	databox "github.com/cgreenhalgh/lib-go-databox"
)

// Note: must match manifest!
const STORE_TYPE = "store-json"

var dataStoreHref = os.Getenv("DATABOX_STORE_ENDPOINT")
var stravaActivityInfo = os.Getenv("DATASOURCE_strava_activity")

// container status entry point
func getStatusEndpoint(w http.ResponseWriter, req *http.Request) {
	if IsStarted() {
		w.Write([]byte("active\n"))
	} else {
		w.Write([]byte("starting\n"))
	}
}

// internal http server.
// Send true on channel when (if) it exits
func server(c chan bool) {
	//
	// Handle Https requests
	//
	router := mux.NewRouter()

	router.HandleFunc("/status", getStatusEndpoint).Methods("GET")
	router.HandleFunc("/ui", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./www/index.html")
	}).Methods("GET")
	
	static := http.StripPrefix("/ui/static", http.FileServer(http.Dir("./www/")))
	router.PathPrefix("/ui/static").Handler(static)

	http.ListenAndServeTLS(":8080", databox.GetHttpsCredentials(), databox.GetHttpsCredentials(), router)
	log.Print("HTTP server exited?!")
	c <- true
}

func main() {
	//log.Printf("Store href %s", dataStoreHref)
	//log.Printf("Strava activity source info %s", stravaActivityInfo)
	
	// embedded web server
	serverdone := make(chan bool)
	go server(serverdone)

	//Wait for my store to become active
	databox.WaitForStoreStatus(dataStoreHref)

	// load state?!
	
	SignalStarted()
	log.Print("App started")
	_ = <-serverdone
}
