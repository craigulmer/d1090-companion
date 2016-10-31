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

	return obuf
}



func getOrGenDigest(day_name string, digest_name string)  (obuf bytes.Buffer) {

	//Get last item in path, strip off all suffixes 
	day_name = cleanDayName(day_name)

	//See if the desired digest file already exists
	digest_fname := dataDir + "/" + day_name + "." + digest_name
	if _, err := os.Stat(digest_fname); err==nil {
		//Use file if already exists
		b, _ := ioutil.ReadFile(digest_fname)
		obuf.Write(b)	
		return obuf
	}

	//Not available, we need to parse 
	plane_counts, plane_ids := parseFile(day_name)
	is_today := (getTodaysDate() == day_name)

	//Work through all the digests
	digests := [2]string{ "stat", "track" }
	for _, d := range digests {
		//When today, only do the digest we need
		if (is_today && (d!=digest_name)) {
			continue
		}
		var tbuf bytes.Buffer
		switch d {
		case "stat":
			tbuf = generateStatPage(day_name, plane_counts, plane_ids)
		case "track":
			tbuf.Reset()
		default:
			//fmt.Printf("Unknown: %s\n", d)
			tbuf.Reset()
		}
		
		//Write valid results to file, if not today
		if((!is_today) && (tbuf.Len()>0)){
			fo, err := os.Create(dataDir+"/"+day_name+"."+d)
			if err != nil {
				return
			}
			defer fo.Close()
			fo.Write(tbuf.Bytes())
		}

		//Pass back the buffer
		if(d==digest_name){
			obuf = tbuf
		}
	}

	return obuf
}

func handleStat(w http.ResponseWriter, r *http.Request){

	AddSafeHeaders(w)
	obuf := getOrGenDigest(r.URL.Path, "stat")	
	w.Write(obuf.Bytes())
}
