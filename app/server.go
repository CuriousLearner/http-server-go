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

func route(conn net.Conn, request Request) {
	var err error
	var response, content, responseType string
	contentType := TEXT_PLAIN

	switch {
	case request.uriPath == "/":
		responseType = OkResponse
	case strings.Contains(request.uriPath, "/echo/"):
		responseType = OkResponse
		content = strings.Split(request.uriPath, "/echo/")[1]
	case strings.Contains(request.uriPath, "/user-agent"):
		responseType = OkResponse
		content = request.headersMap["User-Agent"]
	case strings.Contains(request.uriPath, "/files/"):
		content = ""
		filename := strings.TrimPrefix(request.uriPath, "/files/")
		filePath := filepath.Join(directory, filename)
		if request.httpMethod == "POST" {
			responseType = CreatedresponseType
			err := os.WriteFile(filePath, []byte(request.requestBody), 0644)
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
			contentType = OCTET_STREAM
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

	request := parseRequest(buffer)

	route(conn, request)
}

func parseRequest(buffer []byte) (request Request) {
	fmt.Println("Request received:\n", string(buffer))
	requestString := strings.Split(string(buffer), "\r\n")
	headersMap := generateHeadersMap(requestString)
	requestStartLine := strings.Split(requestString[0], " ")
	httpMethod := requestStartLine[0]
	uriPath := requestStartLine[1]
	requestBody := strings.ReplaceAll(requestString[len(requestString)-1], "\x00", "")
	return Request{httpMethod, uriPath, headersMap, requestBody}
}

type Request struct {
	httpMethod  string
	uriPath     string
	headersMap  map[string]string
	requestBody string
}

const (
	OkResponse             = "200 OK"
	NotFoundResponseType   = "404 Not Found"
	CreatedresponseType    = "201 Created"
	BadRequestResponseType = "400 Bad Request"
)

const (
	TEXT_PLAIN   = "text/plain"
	OCTET_STREAM = "application/octet-stream"
)

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
