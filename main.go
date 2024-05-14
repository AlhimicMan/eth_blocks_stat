package main

import (
	"eth_blocks_stat/infrastructure/config"
	"eth_blocks_stat/infrastructure/http_server"
	"fmt"
	"os"
	"strconv"
)

func main() {
	apiKey := os.Getenv("GET_BLOCK_API_KEY")
	if len(apiKey) == 0 {
		fmt.Println("Please set GET_BLOCK_API_KEY environment variable")
		return
	}
	serverHost := os.Getenv("HTTP_SERVER_HOST")
	if len(apiKey) == 0 {
		fmt.Println("Please set HTTP_SERVER_HOST environment variable")
		return
	}
	serverPortEnv := os.Getenv("HTTP_SERVER_PORT")
	if len(apiKey) == 0 {
		fmt.Println("Please set HTTP_SERVER_PORT environment variable")
		return
	}
	serverPort, err := strconv.Atoi(serverPortEnv)
	if err != nil {
		fmt.Println("Cannot convert HTTP_SERVER_PORT environment variable to int")
	}

	cfg := config.ServiceConfig{
		APIKey:         apiKey,
		HTTPServerHost: serverHost,
		HTTPServerPort: serverPort,
	}

	err = http_server.RunServer(cfg)
	if err != nil {
		fmt.Printf("Error starting HTTP server: %s", err)
	}
	return
}
