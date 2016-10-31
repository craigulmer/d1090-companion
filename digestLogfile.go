package main

import (
	"os"
	"io"
	"io/ioutil"
	"bufio"
	"strings"
	"fmt"
	"compress/gzip"
)

func digestAll() {

	files,_ := ioutil.ReadDir(dataDir)
	for _,f := range files {
		if strings.HasSuffix(f.Name(),".txt") {
			digestLogFile(f.Name())
		}
	}
}


func digestLogFile(fn string) {
	
	fmt.Println("Digesting ",fn)


	//See if we need to regenerate (missing files, or newer source)


	tracks := make(map[string]string)
	plane_ids := make(map[string]string)
	plane_counts := make(map[string]int)

	f,err := os.Open(fn)
	if err != nil {
		return 
	}
	defer f.Close()

	var reader *bufio.Reader

	if strings.HasSuffix(fn,".gz") {
		gz, err := gzip.NewReader(f)
		if err != nil{
			return
		}
		defer gz.Close()
		
		reader = bufio.NewReader(gz)
	} else {
		reader = bufio.NewReader(f)

	}

	for {
		line,err := reader.ReadString('\n')
		if err!=nil{
			if err == io.EOF {
				break
			}
			fmt.Println("Err is "+err.Error())
			return
		}

		cols:=strings.Split(line,"\t")

		if cols[0]=="1" {
			//Update stats
			plane_ids[cols[1]] = cols[2]
			count := plane_counts[cols[1]]
			plane_counts[cols[1]] = count+1
		} else if cols[0]=="3" {
			//Update stats
			count := plane_counts[cols[1]]
			plane_counts[cols[1]] = count+1

			//Update tracks
			s := tracks[cols[1]]
			pt := cols[2]+" "+cols[3]
			if s!="" {
				s+=","+pt
			} else {
				s=pt
			}
			tracks[cols[1]] = s
		}
	}

	//Write out stat file
	generateStatPage(fn, plane_counts, plane_ids)

	//Write out track file



}
