package main

import (
	"fmt"
	"os"

	"github.com/jsirianni/dayz-query-go/dayz"
)

const (
	exitNewClientError       = 1
	exitServerInfoError      = 2
	exitParseServerInfoError = 3
)

func main() {
	// TODO(jsirianni): Should be a config option.
	serverAddr := "50.108.116.1:2324"

	dayzClient, err := dayz.NewClient(serverAddr, dayz.WithTimeoutSeconds(10))
	if err != nil {
		fmt.Println(err)
		os.Exit(exitNewClientError)
	}

	resp, err := dayzClient.ServerInfo()
	if err != nil {
		fmt.Println(err)
		os.Exit(exitServerInfoError)
	}

	info, err := dayz.ParseServerInfo(resp)
	if err != nil {
		fmt.Println(err)
		os.Exit(exitParseServerInfoError)
	}

	fmt.Println("Server Info:")
	fmt.Printf("  Protocol Version: %s\n", info.ProtocolVersion)
	fmt.Printf("  Server Name: %s\n", info.ServerName)
	fmt.Printf("  Map Name: %s\n", info.MapName)
	fmt.Printf("  Game Directory: %s\n", info.GameDirectory)
	fmt.Printf("  Game Description: %s\n", info.GameDescription)
	fmt.Printf("  App ID: %s\n", info.AppID)
	fmt.Printf("  Players: %s\n", info.Players)
	fmt.Printf("  Max Players: %s\n", info.MaxPlayers)
	fmt.Printf("  Bots: %s\n", info.Bots)
}
