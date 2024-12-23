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
	exitModeListError        = 4
)

func main() {
	// TODO(jsirianni): Should be a config option.
	serverAddr := "50.108.116.1:2324"

	dayzClient, err := dayz.NewClient(serverAddr, dayz.WithTimeoutSeconds(10))
	if err != nil {
		err := fmt.Errorf("new client: %v", err)
		fmt.Println(err)
		os.Exit(exitNewClientError)
	}

	info, err := dayzClient.ServerInfo()
	if err != nil {
		err := fmt.Errorf("server info: %v", err)
		fmt.Println(err)
		os.Exit(exitServerInfoError)
	}

	fmt.Println("Server Info:")
	fmt.Printf("  Protocol Version: %s\n", info.ProtocolVersion)
	fmt.Printf("  Server Name: %s\n", info.ServerName)
	fmt.Printf("  Map Name: %s\n", info.MapName)
	fmt.Printf("  Game Directory: %s\n", info.GameDirectory)
	fmt.Printf("  App ID: %s\n", info.AppID)
	fmt.Printf("  Players: %s\n", info.Players)
	fmt.Printf("  Max Players: %s\n", info.MaxPlayers)
	fmt.Printf("  Bots: %s\n", info.Bots)
	fmt.Printf("  Server Type: %s\n", info.ServerType)
	fmt.Printf("  OS Type: %s\n", info.OsType)
	fmt.Printf("  Password Protected: %s\n", info.PasswordProtected)
	fmt.Printf("  VAC Secured: %s\n", info.VacSecured)
	fmt.Printf("  Version: %s\n", info.Version)
}
