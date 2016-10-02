package main

import (
	"os"
	"bufio"
	"strings"
	"fmt"
	"compress/gzip"
)

func parseFile(fn string) (plane_counts map[string]int, plane_ids map[string]string){

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


	plane_ids = make(map[string]string)
	plane_counts = make(map[string]int)

	for {
		line,err := reader.ReadString('\n')
		if err!=nil{
			return 
		}

		cols:=strings.Split(line,"\t")
		if cols[0]=="1" {
			plane_ids[cols[1]] = cols[2]
			count := plane_counts[cols[1]]
			plane_counts[cols[1]] = count+1
		} else if cols[0]=="3"{
			count := plane_counts[cols[1]]
			plane_counts[cols[1]] = count+1
		}
	}
	
	max_hits:=0
	for _,v := range plane_counts {
		if v>max_hits {
			max_hits = v
		}
	}
	fmt.Println(plane_ids)
	fmt.Println("Plane num: ",len(plane_ids))
	fmt.Println("Max planes: ",max_hits)

	return plane_counts, plane_ids
}
