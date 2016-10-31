package main

import (
	"os"
	"bufio"
	"strings"
	"fmt"
	"compress/gzip"
)

func parseFile(day_name string) (plane_counts map[string]int, plane_ids map[string]string){

	day_name = cleanDayName(day_name)

	var reader *bufio.Reader

	//Try finding the file using different extensions
	suffixes:=[2]string{ ".txt.gz", ".txt"}
	for _,s := range suffixes {
		f,err := os.Open(dataDir+"/"+day_name+s)
		if err!= nil {
			continue
		}
		defer f.Close()

		if s==".txt.gz" {
			gz, err := gzip.NewReader(f)
			if err != nil {
				return
			}
			defer gz.Close()
			reader = bufio.NewReader(gz)	
		} else {
			reader = bufio.NewReader(f)
		}
		break
	}
	//Check for not found
	if reader==nil {
		return
	}

	//File valid, we can create the maps now
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
