package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)


func handleRoot( w http.ResponseWriter, r *http.Request){
	
	AddSafeHeadersAndTitle(w, "Dump1090 Helpers")

	w.Write([]byte("<h1>Local Links</h1>"))
	w.Write([]byte("<ul><a href="+d1090Link+">Dump1090 Companion Web UI</a></ul>"))
		

	w.Write([]byte("<h1>Available Files</h1>\n"))

	fmt.Fprintf(w, "<table><tr><td>Name</td></tr>\n");
	files,_ := ioutil.ReadDir(dataDir)
	for _,f := range files {
		if !(strings.HasSuffix(f.Name(),"stat")||strings.HasPrefix(f.Name(),".")) {
			fmt.Fprintf(w, "<tr>");
			fmt.Fprintf(w, "<td><a href=/get/"+f.Name()+">"+f.Name()+"</a></td>")
			fmt.Fprintf(w, "<td><a href=/stat/"+f.Name()+">Stats</a></td>")
			fmt.Fprintf(w, "</tr>")
		}
	}
	fmt.Fprintf(w,"</table>")

}
