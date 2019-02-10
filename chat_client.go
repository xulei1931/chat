package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func CheckError(err error) {
	if err != nil {
		fmt.Println("Error：%s", err.Error())
		//os.Exit(1)
	}

}
func MessageSend(conn net.Conn) {
	var input string
	for {
		reader := bufio.NewReader(os.Stdin)
		data, _, _ := reader.ReadLine()
		input = string(data)
		if strings.ToUpper(input) == "EXIT" {
			conn.Close()
			os.Exit(1)
		}
		_, err := conn.Write([]byte(input))
		if err != nil {
			fmt.Println("client connect error: " + err.Error())
			break

		}

	}

}

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	CheckError(err)
	defer conn.Close()
	go MessageSend(conn)
	buf := make([]byte, 1024)
	for {
		_, err := conn.Read(buf)
		if err != nil {
			fmt.Println("您已退出，欢迎下次光临！")
			os.Exit(0)
		}
		fmt.Println("receive server message contents:" + string(buf))

	}
	fmt.Println("end.....")

}
