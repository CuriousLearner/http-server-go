package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

func handleConnection(conn net.Conn) {

	buffer := make([]byte, 1024)
	_, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading: ", err.Error())
	}

	request := strings.Split(string(buffer), "\r\n")
	headersMap := make(map[string]string)
	for i := 1; i < len(request)-2; i++ {
		header := strings.Split(request[i], ": ")
		headersMap[header[0]] = header[1]
	}
	requestStartLine := strings.Split(request[0], " ")
	uriPath := requestStartLine[1]

	if uriPath == "/" {
		_, err = conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\n"))
		if err != nil {
			fmt.Println("Error writing: ", err.Error())
			os.Exit(1)
		}
	} else if strings.Contains(uriPath, "/echo/") {
		content := strings.Split(uriPath, "/echo/")[1]
		_, err = conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n" + "Content-Length: " + strconv.Itoa(len(content)) + "\r\n\r\n" + content))
		if err != nil {
			fmt.Println("Error writing: ", err.Error())
			os.Exit(1)
		}
	} else if strings.Contains(uriPath, "/user-agent") {
		_, err = conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n" + "Content-Length: " + strconv.Itoa(len(headersMap["User-Agent"])) + "\r\n\r\n" + headersMap["User-Agent"]))
		if err != nil {
			fmt.Println("Error writing: ", err.Error())
			os.Exit(1)
		}
	} else {
		_, err = conn.Write([]byte("HTTP/1.1 404 Not Found\r\nContent-Type: text/plain\r\n\r\n"))
		if err != nil {
			fmt.Println("Error writing: ", err.Error())
			os.Exit(1)
		}
	}
}

func main() {
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleConnection(conn)
	}

}
