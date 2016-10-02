package main

import (
	"fmt"
	"net/http"
	"bufio"
	"io"
	"strings"
	"os"
)

func handlePoints(w http.ResponseWriter, r *http.Request){

	tokens:=strings.Split(r.URL.Path,"/")
	name := tokens[2]
	plane := tokens[3]

	points := parseFilePoints(dataDir+"/"+name, plane)

	AddSafeHeaders(w)

	f,err:=os.Open(htmlDir+"/pointplot_template.html")
	if err!=nil {
		return
	}
	defer f.Close()
	reader:=bufio.NewReader(f)
	for {
		line,err := reader.ReadString('\n')
		if err!=nil{
			if err == io.EOF {
				break
			}
		}
		if strings.Contains(line, "INSERT_DATA_HERE"){
			for i,v := range points {
				if i!=0 {
					fmt.Fprintf(w, ",")
				}
				fmt.Fprintf(w, "[%f,%f]",v[0],v[1])
			}
		} else {
			line = strings.Replace(line,"INSERT_TITLE_HERE", strings.TrimSuffix(name,".txt")+" "+plane,1)
			line = strings.Replace(line,"INSERT_GOOGLE_API_KEY_HERE", googleApiKey,1)
			w.Write([]byte(line))
		}
	}
	
}
