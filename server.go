package main

import "net"
import "fmt"
import (
	"bufio"
	"time"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"io/ioutil"
)

//import "strings" // only needed below for sample processing

type Message struct {
	Ip   string
	Time time.Time
}

type IpName struct {
	Ip   string `json:"ip"`
	Name string `json:"name"`
}

var dbWhoOnLine map[string]time.Time
var dbIpName map[string]string

func init() {
	dbWhoOnLine = make(map[string]time.Time)
	dbIpName = make(map[string]string)
	var tempIpName []IpName
	file, e := ioutil.ReadFile("./dbIpName.json")
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}
	e = json.Unmarshal(file, &tempIpName)
	if e != nil {
		log.Println(e)
	}
	for _, val := range tempIpName {
		dbIpName[val.Ip] = val.Name
	}
}

func main() {
	fmt.Println("Launching server...")

	// listen on all interfaces
	ln, _ := net.Listen("tcp", ":10543")

	// accept connection on port
	// run loop forever (or until ctrl-c)
	go func() {
		f, err := os.OpenFile("log.log", os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Fatal(err)
			}
			// will listen for message to process ending in newline (\n)
			message, _ := bufio.NewReader(conn).ReadString('\n')
			// output message received
			mes := Message{}
			err = json.Unmarshal([]byte(message), &mes)
			if err != nil {
				log.Println(err)
			}
			// send new string back to client
			conn.Write([]byte("Done!\n"))

			_, ok := dbIpName[mes.Ip]
			fmt.Println(mes.Ip)
			if !ok {
				dbIpName[mes.Ip] = mes.Ip
			}

			s := dbIpName[mes.Ip] + " включился в " + mes.Time.String() + "\n"

			dbWhoOnLine[mes.Ip] = mes.Time

			if _, err = f.WriteString(s); err != nil {
				panic(err)
			}

		}
	}()

	time.Sleep(time.Second * 5)

	http.HandleFunc("/", IndexHandle)
	http.HandleFunc("/file", FileHandle)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func IndexHandle(w http.ResponseWriter, r *http.Request) {
	body := ""
	for key, val := range dbWhoOnLine {
		body += dbIpName[key] + " включился в " + val.String() + "\n"
	}
	fmt.Fprint(w, body)
}

func FileHandle(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("log.log")
	if err != nil {
		// handle the error here
		return
	}
	defer file.Close()

	// get the file size
	stat, err := file.Stat()
	if err != nil {
		return
	}

	bs := make([]byte, stat.Size())
	_, err = file.Read(bs)
	if err != nil {
		return
	}
	str := string(bs)
	fmt.Fprint(w, str)
}
