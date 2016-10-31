package main

import (
	"io/ioutil"
	"compress/gzip"
	"strings"
	"os"
	"fmt"
)

func startupCompress(){

	//Walk through all input data files and compress (if not today's data)
	files,_ := ioutil.ReadDir(dataDir)
	for _,f := range files {

		//Bail if not a text file, or is today's data
		if ( (!strings.HasSuffix(f.Name(),"txt")) ||
			  ((getTodaysDate()+".txt"==f.Name()))    ){
			continue
		}

		//Read the file
		b, err := ioutil.ReadFile(dataDir+"/"+f.Name())
		if(err !=nil) {
			fmt.Println("No read\n");
			continue
		}

		//Open the file for writing
		fo, err := os.Create(dataDir+"/"+f.Name()+".gz")
		if (err != nil){
			fmt.Println("No create\n");
			continue
		}
		defer fo.Close()

		//Write buffer to gzip
		w := gzip.NewWriter(fo)
		w.Write(b)
		w.Close()

		//Remove original uncompressed file
		os.Remove(dataDir+"/"+f.Name())
	}

}
