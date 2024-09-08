package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/pflag"
	"hippos/walrus"
)

func main() {
	var putData string
	var getBlobID string
	var apiURL string

	defaultURL := "http://localhost:31415"

	pflag.StringVar(&putData, "put", "", "Data to PUT to Walrus")
	pflag.StringVar(&getBlobID, "get", "", "Blob ID to GET from Walrus")
	pflag.StringVar(&apiURL, "api", defaultURL, "Walrus service URL")
	pflag.Parse()

	if pflag.CommandLine.HasFlags() && (putData == "" && getBlobID == "" && !pflag.CommandLine.Parsed()) {
		pflag.Usage()
		os.Exit(0)
	}

	client := walrus.NewWalrusClient(apiURL)

	if putData != "" {
		data := []byte(putData)
		blobID, err := client.PutData(data)
		if err != nil {
			log.Fatalf("Failed to put data: %v", err)
		}
		fmt.Println(blobID)
	}

	if getBlobID != "" {
		data, err := client.GetData(getBlobID)
		if err != nil {
			log.Fatalf("Failed to get data: %v", err)
		}
		fmt.Println(string(data))
	}

	if putData == "" && getBlobID == "" {
		fmt.Println("Please specify either --put or --get")
		pflag.Usage()
		os.Exit(1)
	}
}
