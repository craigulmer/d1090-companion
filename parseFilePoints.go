package main

import (
	"os"
	"io"
	"bufio"
	"strings"
	"fmt"
	"compress/gzip"
	"strconv"
)

func parseFilePoints(fn string, aid string) ([][2]float64){

	values := [][2]float64{}

	f,err := os.Open(fn)
	if err != nil {
		return values
	}
	defer f.Close()

	var reader *bufio.Reader

	if strings.HasSuffix(fn,".gz") {
		gz, err := gzip.NewReader(f)
		if err != nil{
			return values
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
			return values
		}

		cols:=strings.Split(line,"\t")
		if (cols[0]=="3") && (cols[1]==aid) {

			xy:=[2]float64{}
			xy[0],err= strconv.ParseFloat(cols[2],64)
			if err!=nil{
				continue
			}
			xy[1],err= strconv.ParseFloat(cols[3],64)
			if err!=nil{
				continue
			}
			values = append(values, xy)
		}
	}

	return values
}
