package main

import "net"
import "fmt"
import "bufio"
import (
	"time"
	"encoding/json"
	"io/ioutil"
	"os"
	"log"
)

type Message struct {
	Ip string
	Time time.Time
}

type Conf struct {
	Ip string `json:"ip"`
	IpServer string `json:"ipserver"`
	IpServerPort string `json:"ipserverport"`
}

func main() {
	file, e := ioutil.ReadFile("./config.json")
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}
	var i Conf
	e = json.Unmarshal(file, &i)
	if e != nil {
		log.Println(e)
	}

	// connect to this socket
	var conn net.Conn
	var err error
	mes := Message{Ip:i.Ip, Time: time.Now()}
	address := i.IpServer + ":" +i.IpServerPort
	for ;; {
		conn, err = net.Dial("tcp", address)
		if err == nil {
			break
		}
		time.Sleep(time.Second * 10)
	}
	for {
		// read in input from stdin
		// send to socket
		text, _ := json.Marshal(mes)
		fmt.Fprintf(conn, string(text) + "\n")
		// listen for reply
		message, _ := bufio.NewReader(conn).ReadString('\n')
		if message == "Done!\n" {
			fmt.Println("Done!")
			return
		}
		fmt.Print("Message from server: "+message)
	}
}