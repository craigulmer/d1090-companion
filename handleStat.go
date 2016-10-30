package main

import (
	"net/http"
	"sort"
	"strings"
	"strconv"
	"bytes"
	"io/ioutil"

	"os"
)
func getCategories( icaos []string, flight_names map[string]string ) ([]string, []string) {

	icaos_mil := make([]string, 0, 20)
	icaos_pkg := make([]string, 0, 20)


	for _,k := range icaos {
		//Military planes seem to be in AE0000, nato may be in 500000
		if (k>="AE0000" && k<="AEFFFF") || (k>="500000" && k<="5FFFFF") { 
			icaos_mil = append(icaos_mil, k)
		}
		flt:=flight_names[k]
		if strings.HasPrefix(flt, "FDX") || strings.HasPrefix(flt, "UPS") {
			icaos_pkg = append(icaos_pkg, k)
		}
	}
	return icaos_mil, icaos_pkg
}

func generateStatPage(name string, plane_counts map[string]int, plane_ids map[string]string) bytes.Buffer {

	var obuf bytes.Buffer

	obuf.WriteString("<!DOCTYPE HTML PUBLIC \"-//W3C//DTD HTML 3.2 Final//EN\">")
	obuf.WriteString("<TITLE>Stats: "+name+"</TITLE>")
	obuf.WriteString("<h1>Stats for "+name+"</h1>\n")

	//Make a list of icaos
	icaos := make([]string, 0, len(plane_ids))
	for k := range plane_ids {
		icaos = append(icaos, k)
	}
	sort.Strings(icaos)

	//Get Lists
	icaos_mil, icaos_pkg := getCategories(icaos, plane_ids)
	
	obuf.WriteString("<table><tr><td valign=top>")

	//Dump Military Planes
	obuf.WriteString("<h2>Military Planes</h2>")
	obuf.WriteString("<table>")
	for _,k := range icaos_mil {
		obuf.WriteString("<tr><td>"+k+"</td><td>"+plane_ids[k]+"</td><td>"+strconv.Itoa(plane_counts[k])+"</td></tr>")
	}
	obuf.WriteString("</table>")
	
	//Next Column
	obuf.WriteString("</td><td width=100px></td><td valign=top>")

	//Dump Package Planes
	obuf.WriteString("<h2>Package Planes</h2>")
	obuf.WriteString("<table>")
	for _,k := range icaos_pkg {
		var link_id string = plane_ids[k]
		if link_id != "" {
			link_id = "<a href=https://flightaware.com/live/flight/"+link_id+">"+link_id+"</a>"
		}
		obuf.WriteString("<tr><td>"+k+"</td><td>"+link_id+"</td><td>"+strconv.Itoa(plane_counts[k])+"</td></tr>")
	}
	obuf.WriteString("</table>")
		
	//Next Column
	obuf.WriteString("</td><td width=100px></td><td valign=top>")

	obuf.WriteString("<h2>All Planes</h2>")
	obuf.WriteString("<table>")
	for _,k := range icaos {
		var link_id string = plane_ids[k]
		if link_id != "" {
			link_id = "<a href=https://flightaware.com/live/flight/"+link_id+">"+link_id+"</a>"
		}
		var parse_link string 
		parse_link = "<a href=/points/"+name+"/"+k+">"+strconv.Itoa(plane_counts[k])+"</a>"
		obuf.WriteString("<tr><td>"+k+"</td><td>"+link_id+"</td><td>"+parse_link+"</td></tr>")
	}
	obuf.WriteString("</table>")

	//End of columns
	obuf.WriteString("</td></tr></table>")

	stat_fname, is_today := convertSourceToStatFilename(name)

	//Check fname to see if it should be saved
	if (stat_fname!="") && (!is_today) {

		//Write to file
		fo, err := os.Create(stat_fname)
		if err != nil {
			return
		}
		defer fo.Close()
		fo.Write(obuf.Bytes())
	}

	//Either way, return the buffer
	return obuf
}

func convertSourceToStatFilename(s String)  (String, bool) {
	tokens := strings.Split(s,"/")
	name:=tokens[len(tokens)-1]
	
	spot:=strings.Index(name,".")

	if spot < 0 {
		return "", false
	}

	//Find out if it was for today
	sdate := getTodaysFilename()
	file_is_today := (sdate==name)

	stat_fname := dataDir + "/" + name[0:spot]+".stat"

	return stat_fname, file_is_today
}

func handleStat(w http.ResponseWriter, r *http.Request){

	var obuf bytes.Buffer

	stat_fname, is_today := convertSourceToStatFilename(r.URL.Path)
	if stat_fname=="" {
		return
	}

	AddSafeHeaders(w)
	
	if _, err := os.Stat(stat_fname); err==nil {
		//Use file if already exists
		b, _ := ioutil.ReadFile(stat_fname)
		obuf.Write(b)

	} else {
		//Generate stats file
		plane_counts, plane_ids:=parseFile(dataDir+"/"+name)
		obuf = generateStatPage(name, plane_counts, plane_ids)
	}
	w.Write(obuf.Bytes())
}
