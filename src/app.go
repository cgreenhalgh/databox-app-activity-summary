package main

import (
	"os"
	"log"
)

var dataStoreHref = os.Getenv("DATABOX_STORE_ENDPOINT")


var stravaActivityInfo = os.Getenv("DATASOURCE_strava_activity")

func main() {
	log.Printf("Store href %s", dataStoreHref)
	log.Printf("Strava activity source info %s", stravaActivityInfo)
}
