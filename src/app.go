package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"time"
	//"strings"
	
	"github.com/gorilla/mux"
	databox "github.com/cgreenhalgh/lib-go-databox"
)

// Note: must match manifest!
const STORE_TYPE = "store-json"

var dataStoreHref = os.Getenv("DATABOX_STORE_ENDPOINT")

// container status entry point
func getStatusEndpoint(w http.ResponseWriter, req *http.Request) {
	if IsStarted() {
		w.Write([]byte("active\n"))
	} else {
		w.Write([]byte("starting\n"))
	}
}

// Index of available data 
type DataType struct{
	Title string `json:"title"`
	Id string `json:"id"`
	Available bool `json:"available"`
	dsInfo string
	Source databox.TimeSeries_0_2_0
}

var dataTypes = []DataType{
	DataType{
		Title:"Strava Activities",
		// NB must match client name in manifest
		Id:"strava_activity",
	},
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
		log.Printf("Run angular with: `npm bin`/ng serve --host 0.0.0.0 --disable-host-check -bh /databox-app-activity-summary/ui/static/")
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
	router.HandleFunc("/ui/api/dataTypes", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("GET dataTypes (started? %t)", IsStarted())
		WaitUntilStarted()
		data,err := json.Marshal(dataTypes)
		if err != nil {
			log.Printf("Error marshalling dataTypes: %s", err.Error())
			w.WriteHeader(500)
			w.Write([]byte("Error marshalling dataTypes"))
			return
		}
		w.Write(data)
	}).Methods("GET")
	router.HandleFunc("/ui/api/get/{source}/", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		source := vars["source"]
		log.Printf("GET datasource %s (started? %t)", source, IsStarted())
		WaitUntilStarted()
		var dataType *DataType = nil
		for _,dt := range dataTypes {
			log.Printf("- check %s: %t", dt.Id, dt.Available)
			if source == dt.Id {
				dataType = &dt
				break
			}
		}
		if dataType==nil || !dataType.Available {
			log.Printf("datasource %s not known/available", source)
			w.WriteHeader(404)
			w.Write([]byte("Datasource "+source+" not known/available"))
			return
		}
		data,err := dataType.Source.ReadSince(time.Unix(0,0))
		if err != nil {
			log.Printf("Error getting values for %s: %s", source, err.Error())
			w.WriteHeader(500)
			w.Write([]byte("Error getting values for "+source+": "+err.Error()))
			return
		}
		// process??
		w.Write([]byte(data))
	}).Methods("GET")
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
	log.Print("App's store is ready")

	// load state?!

	// datasources
	for i:=0; i<len(dataTypes); i++ {
		dataType := &dataTypes[i]
		var sourceInfo = os.Getenv("DATASOURCE_"+dataType.Id)
		dataType.dsInfo = sourceInfo
		if sourceInfo=="" {
			log.Printf("Warning: Source %s undefined - skipping", dataType.Id)
			dataType.Available = false
			continue
		}
		var err error = nil
		dataType.Source,err = databox.MakeSourceTimeSeries_0_2_0(sourceInfo)
		if  err != nil {
			log.Printf("Error: making TimeSeries for source %s: %s", dataType.Id, err.Error())
			dataType.Available = false
			continue;
		}
		surl,_ := dataType.Source.StoreURL()
		databox.WaitForStoreStatus(surl)
		dataType.Available = true
		log.Printf("Datasource %s is ready", dataType.Id)
	}
	
	SignalStarted()
	log.Print("App started")
	_ = <-serverdone
}
