package main

import (
	"io/ioutil"
	"net/http"
	"strings"
)

func handleGet( w http.ResponseWriter, r *http.Request){
	
	w.Header().Set("Content-Type", "text/plain; charset=utf-8") // normal header

	tokens:=strings.Split(r.URL.Path,"/")
	name:=tokens[len(tokens)-1]
	
	fbytes, err:= ioutil.ReadFile(dataDir+"/"+name)
	if err!=nil {
		return
	}
	w.Write(fbytes)
}
