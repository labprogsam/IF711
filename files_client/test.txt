package base

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

const clientFiles = "../files_client"

type Client interface {
	SendCommand(cmd string, params any) string
}

type BaseClient struct {
	client Client
}

func NewBaseClient(c Client) *BaseClient {
	return &BaseClient{client: c}
}

func (b *BaseClient) HandleCommand(args []string, params any) {
	command := args[0]

	switch command {
	case "LIST":
		start := time.Now()
		resp := b.client.SendCommand("LIST", params)
		log.Println(fmt.Sprintf("list - %v", time.Since(start)))
		fmt.Println("Files on server:\n" + resp)

	case "UPLOAD":
		if len(args) < 3 {
			fmt.Println("Usage: UPLOAD <ID> <filename>")
			return
		}
		filename := args[2]
		data, err := os.ReadFile(clientFiles + "/" + filename)
		if err != nil {
			fmt.Println("ERROR: File not found")
			return
		}
		encoded := base64.StdEncoding.EncodeToString(data)
		fullCmd := fmt.Sprintf("UPLOAD %s %s", filename, encoded)
		start := time.Now()
		resp := b.client.SendCommand(fullCmd, params)
		log.Println(fmt.Sprintf("upload - %v", time.Since(start)))
		fmt.Println(resp)

	case "DOWNLOAD":
		if len(args) < 3 {
			fmt.Println("Usage: DOWNLOAD <ID> <filename>")
			return
		}
		filename := args[2]
		start := time.Now()
		resp := b.client.SendCommand("DOWNLOAD "+filename, params)
		log.Println(fmt.Sprintf("download - %v", time.Since(start)))

		if strings.HasPrefix(resp, "ERROR") {
			fmt.Println(resp)
		} else {
			decoded, err := base64.StdEncoding.DecodeString(resp)
			if err != nil {
				fmt.Println("ERROR decoding file")
				return
			}
			os.WriteFile(clientFiles+"/"+filename, decoded, 0644)
			fmt.Println("File downloaded:", filename)
		}

	default:
		fmt.Println("Unknown command.")
	}
}
