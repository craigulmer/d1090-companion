
// To build for 32b:  GOARCH=386 go build
package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

var enableStartup bool
var enableCapture bool
var enableStdout bool

var d1090Link string // = "http://192.168.1.1:8080"
var d1090Raw string // = "192.168.1.1:30003"
var googleApiKey string

var rootDir string  // = "/path/to/root"
var dataDir string //= rootDir+"/data"
var htmlDir string //= rootDir+"/html"

var myPort string

func getTodaysDate() string {
	t := time.Now().UTC()
	sdate := t.Format("060102")
	return sdate
}
func cleanDayName(day_name string) string {

	//Get last item in path, strip off all suffixes 
	tokens := strings.Split(day_name,"/")
	day_name = tokens[len(tokens)-1]
	spot := strings.Index(day_name,".")
	if(spot > 0){
		day_name = day_name[0:spot]
	}
	return day_name
}
func AddSafeHeaders(w http.ResponseWriter) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("X-Frame-Options", "SAMEORIGIN")
	w.Header().Set("Strict-Transport-Security", "max-age=2592000; includeSubDomains")
}

func AddSafeHeadersAndTitle(w http.ResponseWriter, title string) {
	AddSafeHeaders(w)
	w.Write([]byte("<!DOCTYPE HTML PUBLIC \"-//W3C//DTD HTML 3.2 Final//EN\">"))
	w.Write([]byte("<TITLE>"+title+"</TITLE>"))
}





func main() {

	parseConfig()

	digestAll();

	if enableStartup {
		go startupCompress()
	}

	if enableCapture {
		go captureData()
	}

	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/get/", handleGet)
	http.HandleFunc("/stat/", handleStat)
	http.HandleFunc("/points/", handlePoints)
	err := http.ListenAndServe(":"+myPort, nil)
	if err != nil {
		fmt.Println("Error is ",err)
	}
}
