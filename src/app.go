package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	//"strings"
	
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
	// for dev set this to true!
	proxy := len(os.Args)>1 && os.Args[1]=="proxy"
	if proxy {
		// Note: this is just for development - to make it faster
		log.Printf("Proxy static requests")
		log.Printf("Run angular with: ng serve --host 0.0.0.0 --disable-host-check -bh /databox-app-activity-summary/ui/static/")
		router.HandleFunc("/ui", func(w http.ResponseWriter, r *http.Request) {
			log.Printf("proxy %s", r.URL)
			director := func(req *http.Request) {
				req = r
				req.URL.Scheme = "http"
				req.URL.Host = "127.0.0.1:4200"
				req.URL.Path = "/"
			}
			proxy := &httputil.ReverseProxy{Director: director}
			proxy.ServeHTTP(w, r)
		})		
		router.PathPrefix("/ui/static").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("proxy %s", r.URL)
			director := func(req *http.Request) {
				req = r
				req.URL.Scheme = "http"
				req.URL.Host = "127.0.0.1:4200"
				req.URL.Path = "/databox-app-activity-summary/ui/static" + req.URL.Path[len("/ui/static"):]
				//req.Header["Host"] = []string{"127.0.0.1:4200"}
				log.Printf("-> %s", req.URL)
			}
			proxy := &httputil.ReverseProxy{Director: director}
			proxy.ServeHTTP(w, r)
		})
	} else {
		router.HandleFunc("/ui", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "./www/index.html")
		}).Methods("GET")
		
		static := http.StripPrefix("/ui/static", http.FileServer(http.Dir("./www/")))
		router.PathPrefix("/ui/static").Handler(static)
	}
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
