package base

import (
	"encoding/base64"
	"os"
	"strings"
)

const (
	requestTopic = "file/request"
	serverDir    = "../files_server"
)

type Server interface {
	HandleRequest()
	ListFiles()
	SendFile()
	SaveFile()
}

type BaseServer struct{}

func (b *BaseServer) SaveFile(filename, data string) string {
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "ERROR: Failed to decode file"
	}
	err = os.WriteFile(serverDir+"/"+filename, decoded, 0644)
	if err != nil {
		return "ERROR: Failed to save file"
	}
	return "SUCCESS: File uploaded"
}

func (b *BaseServer) SendFile(filename string) string {
	data, err := os.ReadFile(serverDir + "/" + filename)
	if err != nil {
		return "ERROR: File not found"
	}
	return base64.StdEncoding.EncodeToString(data)
}

func (b *BaseServer) ListFiles() string {
	files, err := os.ReadDir(serverDir)
	if err != nil {
		return "ERROR: Unable to list files"
	}
	var fileNames []string
	for _, f := range files {
		fileNames = append(fileNames, f.Name())
	}
	return strings.Join(fileNames, "\n")
}

func (b *BaseServer) HandleRequest(command, filename, data string) string {

	var response string

	switch command {
	case "LIST":
		response = b.ListFiles()
	case "DOWNLOAD":
		if filename == "" {
			response = "ERROR: Invalid DOWNLOAD command"
		} else {
			response = b.SendFile(filename)
		}
	case "UPLOAD":
		if filename == "" || data == "" {
			response = "ERROR: Invalid UPLOAD command"
		} else {
			response = b.SaveFile(filename, data)
		}
	default:
		response = "ERROR: Unknown command"
	}

	return response
}
