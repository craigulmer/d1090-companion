package main

import (
	"fmt"
	"bufio"
	"net"
	"strings"
	"time"
	"os"
)


func captureData(){

	con, err:= net.Dial("tcp", d1090Raw)
	if err!=nil {
		fmt.Println("Error is",err)
		return
	}

	var sdate_prv string
	var fname string
	var f *os.File

	for {
		msg, err := bufio.NewReader(con).ReadString('\n')
		if err!=nil {
			fmt.Println("dumper error ",err)
			return
		}
		t := time.Now().UTC()
		sdate := t.Format("060102.txt")

		x := strings.Split(msg,",")
		if x[0]!="MSG" { 
			continue
		}

		if sdate != sdate_prv {
			if sdate_prv!="" {
				f.Close()
				//push fname to gzip
			}

			fname=dataDir+"/"+sdate
			flag := os.O_CREATE | os.O_WRONLY
			if _,err := os.Stat(fname); err==nil {
				flag = os.O_APPEND | os.O_WRONLY
			}

			f,err = os.OpenFile(fname, flag, 0600)	
			if err!=nil{
				continue
			}
			sdate_prv = sdate
		}
		
		id, d1, t1, d2, t2 := x[4],  x[6],  x[7],  x[8], x[9]
		flt, alt, lat, lon := x[10], x[11], x[14], x[15]

		switch x[1] {
		case "1":
			tmp := fmt.Sprintf("1\t%s\t%s\t%s\t%s\n", id, flt, d1, t1)
			if enableStdout {
				fmt.Printf(tmp)
			}
			f.Write([]byte(tmp))

		case "3":
			tmp := fmt.Sprintf("3\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n", id, lat, lon, alt, d1, t1, d2, t2)
			if enableStdout {
				fmt.Printf(tmp)
			}
			f.Write([]byte(tmp))
		}
	

		//fmt.Print("Message is :",string(msg))
	}
}
