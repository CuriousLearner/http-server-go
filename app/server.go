package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"path/filepath"
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

func generateResponse(responseType string, content string, contentType string) string {
	response := "HTTP/1.1 %s\r\nContent-Type: %s\r\nContent-Length: %d\r\n\r\n%s"
	return fmt.Sprintf(response, responseType, contentType, len(content), content)
}

func route(conn net.Conn, httpMethod string, uriPath string, headersMap map[string]string, requestBody string) {
	var err error
	var response, content, responseType string
	OkResponse := "200 OK"
	NotFoundResponseType := "404 Not Found"
	CreatedresponseType := "201 Created"
	BadRequestResponseType := "400 Bad Request"
	contentType := "text/plain"

	switch {
	case uriPath == "/":
		responseType = OkResponse
	case strings.Contains(uriPath, "/echo/"):
		responseType = OkResponse
		content = strings.Split(uriPath, "/echo/")[1]
	case strings.Contains(uriPath, "/user-agent"):
		responseType = OkResponse
		content = headersMap["User-Agent"]
	case strings.Contains(uriPath, "/files/"):
		content = ""
		filename := strings.TrimPrefix(uriPath, "/files/")
		filePath := filepath.Join(directory, filename)
		if httpMethod == "POST" {
			responseType = CreatedresponseType
			err := os.WriteFile(filePath, []byte(requestBody), 0644)
			if err != nil {
				fmt.Println("Error writing file:", err)
				responseType = BadRequestResponseType
				break
			}
		} else {
			responseType = NotFoundResponseType
			data, err := os.ReadFile(filePath)
			if err != nil {
				fmt.Println("Error reading file:", err)
				break
			}
			responseType = OkResponse
			content = string(data)
			contentType = "application/octet-stream"
		}
	default:
		responseType = NotFoundResponseType
		content = ""
	}
	response = generateResponse(responseType, content, contentType)
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
	httpMethod := requestStartLine[0]
	uriPath := requestStartLine[1]
	requestBody := strings.ReplaceAll(request[len(request)-1], "\x00", "")
	route(conn, httpMethod, uriPath, headersMap, requestBody)
}

var directory string

func main() {
	fmt.Println("Logs from your program will appear here!")
	flag.StringVar(&directory, "directory", "", "the directory to serve files from")
	flag.Parse()

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
