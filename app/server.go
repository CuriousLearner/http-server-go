package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func generateHeadersMap(request []string) map[string]string {
	headersMap := make(map[string]string)
	for i := 1; i < len(request)-2; i++ {
		header := strings.Split(request[i], ": ")
		if len(header) != 2 {
			fmt.Println("Error parsing header: ", request[i])
			continue
		}
		headersMap[header[0]] = header[1]
	}
	return headersMap
}

func generateResponse(responseType string, content string) string {
	response := "HTTP/1.1 %s\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s"
	return fmt.Sprintf(response, responseType, len(content), content)
}

func route(conn net.Conn, uriPath string, headersMap map[string]string) {
	var err error
	var response, content, responseType string
	OkResponse := "200 OK"
	NotFoundResponseType := "404 Not Found"

	switch {
	case uriPath == "/":
		responseType = OkResponse
	case strings.Contains(uriPath, "/echo/"):
		responseType = OkResponse
		content = strings.Split(uriPath, "/echo/")[1]
	case strings.Contains(uriPath, "/user-agent"):
		responseType = OkResponse
		content = headersMap["User-Agent"]
	default:
		responseType = NotFoundResponseType
		content = ""
	}
	response = generateResponse(responseType, content)
	_, err = conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Error writing: ", err.Error())
	}
}

func handleConnection(conn net.Conn) {

	buffer := make([]byte, 1024)
	_, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading: ", err.Error())
	}

	request := strings.Split(string(buffer), "\r\n")
	headersMap := generateHeadersMap(request)
	requestStartLine := strings.Split(request[0], " ")
	uriPath := requestStartLine[1]

	route(conn, uriPath, headersMap)
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
		}

		go handleConnection(conn)
	}

}
