package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

const (
	LOG_DICTORY = "../log/test.log"
)

var logFile *os.File
var logger *log.Logger
var onlineConns = make(map[string]net.Conn)
var messageQueue = make(chan string, 1000)
var quitChan = make(chan bool)

func CheckError(err error) {
	if err != nil {
		panic(err)
	}

}
func ProcessInfo(conn net.Conn) {
	buf := make([]byte, 1024)
	// 清除ip
	defer func(conn net.Conn) {
		addr := fmt.Sprintf("%s", conn.RemoteAddr())
		delete(onlineConns, addr)
		conn.Close()
		for ip := range onlineConns {
			fmt.Println("now online cons:" + ip)
		}
	}(conn)

	for {
		numByte, err := conn.Read(buf)
		if err != nil {
			continue
		}
		if numByte != 0 {
			message := string(buf[0:numByte])
			//remoteaddr := conn.RemoteAddr()
			//fmt.Print(remoteaddr, "\n")
			//fmt.Printf("has receive message:%s\n ", string(buf[0:numByte]))
			messageQueue <- message
		}
	}
}

// 消费message
func ConsumeMessage() {
	for {
		select {
		case message := <-messageQueue:

			//处理消息
			doChanMessage(message)

		case <-quitChan:
			break
		}
	}
}

// 发送消息
func doChanMessage(message string) {
	contents := strings.Split(message, "#")
	if len(contents) > 1 {
		addr := contents[0]
		sendmessage := strings.Join(contents[1:], "#")
		addr = strings.Trim(addr, " ")
		if conn, ok := onlineConns[addr]; ok {
			_, err := conn.Write([]byte(sendmessage))
			if err != nil {
				fmt.Println("online conns send fail")
			}
		}
	} else {
		contents := strings.Split(message, "*")
		if strings.ToUpper(contents[1]) == "LIST" {
			var ips string = ""
			for i := range onlineConns {
				ips = ips + "|" + i
			}
			if conn, ok := onlineConns[contents[0]]; ok {
				_, err := conn.Write([]byte(ips))
				if err != nil {
					fmt.Println("online conns send fail")
				}
			}
		}

	}
}
func main() {
	logFile, err := os.OpenFile(LOG_DICTORY, os.O_RDWR|os.O_CREATE, 0)
	if err != nil {
		fmt.Println("log file create fail.....")
		os.Exit(-1)
	}
	defer logFile.Close()
	logger = log.New(logFile, "\r\n", log.Ldate|log.Ltime|log.Llongfile)
	listen_socker, err := net.Listen("tcp", "127.0.0.1:8080")
	CheckError(err)
	defer listen_socker.Close()
	fmt.Printf("server is starting.......\n")
	logger.Println("server is wating........")

	// 消费message
	go ConsumeMessage()
	for {
		conn, err := listen_socker.Accept()
		CheckError(err)
		// 将conn 储存到映射表
		addr := fmt.`("%s", conn.RemoteAddr())
		onlineConns[addr] = conn
		for ip := range onlineConns {
			fmt.Println(ip)
		}
		go ProcessInfo(conn)
	}

}
