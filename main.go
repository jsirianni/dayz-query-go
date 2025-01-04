package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/jsirianni/dayz-query-go/config"
	"github.com/jsirianni/dayz-query-go/dayz"
	"go.uber.org/zap"
)

const (
	exitNewClientError = 1
	exitNewConfigError = 5
	exitNewLoggerError = 6
)

func main() {
	deerIsle := "50.108.116.1:2324"
	namalsk := "50.108.116.1:2315"
	s := []string{deerIsle, namalsk}
	os.Setenv("DAYZ_SERVER_LIST", strings.Join(s, ","))

	logger, err := zap.NewProduction()
	if err != nil {
		err := fmt.Errorf("new logger: %v", err)
		fmt.Println(err)
		os.Exit(exitNewLoggerError)
	}

	config, err := config.New(logger)
	if err != nil {
		err := fmt.Errorf("new config: %v", err)
		fmt.Println(err)
		os.Exit(exitNewConfigError)
	}

	clients := make([]*dayz.Client, 0, len(config.ServerList))

	for _, server := range config.ServerList {
		dayzClient, err := dayz.NewClient(server.String(), dayz.WithTimeoutSeconds(10))
		if err != nil {
			logger.Error("new client", zap.Error(err))
			os.Exit(exitNewClientError)
		}
		clients = append(clients, dayzClient)
	}

	for _, dayzClient := range clients {
		info, err := dayzClient.ServerInfo()
		if err != nil {
			logger.Error("server info", zap.Error(err))
			continue
		}

		logger.Info(
			"server info",
			zap.String("protocol_version", info.ProtocolVersion),
			zap.String("server_name", info.ServerName),
			zap.String("map_name", info.MapName),
			zap.String("game_directory", info.GameDirectory),
			zap.String("app_id", info.AppID),
			zap.String("players", info.Players),
			zap.String("max_players", info.MaxPlayers),
			zap.String("bots", info.Bots),
			zap.String("server_type", info.ServerType),
			zap.String("os_type", info.OsType),
			zap.String("password_protected", info.PasswordProtected),
			zap.String("vac_secured", info.VacSecured),
			zap.String("version", info.Version),
		)
	}
}
