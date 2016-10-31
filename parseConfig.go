package main

import (
	"flag"
	"fmt"
	"os"
	"log"
	"strings"
	"bufio"
	"path/filepath"
	"io"
)

func makeAbsolutePath(root_dir string, new_path string) string {

	if new_path=="." {
		return root_dir
	}

	if strings.HasPrefix(new_path,"/") {
		return new_path
	}
	return root_dir+"/"+new_path
}

func parseConfig() {

	var fname = flag.String("config","default.conf",
		                     "configuration file for d1090-companion")
	flag.Parse()

	//The default root dir is the place where the bin is installed
	tmp_rootDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	//Use either the absolute path handed to us, or our default config dir
	conf_file := *fname
	if !strings.HasPrefix(conf_file,"/") {
		conf_file = tmp_rootDir+"/configs/"+*fname
	}

	//Start parsing config or die
	f,err := os.Open(conf_file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	reader := bufio.NewReader(f)

	//Set the default vals so we always
	vals:=map[string]string{
		"enable_startup"    : "true",
		"enable_capture"    : "true",
		"enable_stdout"     : "false",
		"my_ip"             : "192.168.1.1",
		"my_port"           : "9090",
		"my_root_dir"       : ".",
		"my_data_dir"       : "data",
		"my_html_dir"       : "html",
		"my_goodle_api_key" : "YOUR_KEY_HERE",
		"d1090_ip"          : "192.168.1.1",
		"d1090_sport"       : "30003",
		"d1090_wport"       : "8080",
	}

	for {
		line,err := reader.ReadString('\n')
		if err!=nil{
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		line = strings.Split(line,"#")[0] //Remove comments
		fields:=strings.Fields(line)
		switch len(fields){
		case 0:
			continue //Empty line
		case 2:
			vals[fields[0]] = fields[1] //Remember Two entries
		default:
			//Too many/too few entries
			log.Fatal("Parse error for line:\n"+line+"\n")
		}
	}
	enableStartup = (vals["enable_startup"] == "true")
	enableCapture = (vals["enable_capture"] == "true")
	enableStdout  = (vals["enable_stdout"] == "true")

	rootDir = makeAbsolutePath(tmp_rootDir, vals["my_root_dir"])
	dataDir = makeAbsolutePath(rootDir,     vals["my_data_dir"])
	htmlDir = makeAbsolutePath(rootDir,     vals["my_html_dir"])

	myPort = vals["my_port"]

	d1090Link = "http://" + vals["d1090_ip"]+":"+ vals["d1090_wport"]
	d1090Raw  = vals["d1090_ip"]+":"+vals["d1090_sport"]
	
	googleApiKey = vals["my_google_api_key"]

	if enableStdout {
		fmt.Println(d1090Link+"\n"+d1090Raw+"\n"+rootDir+"\n"+dataDir+"\n"+htmlDir)
	}
}

