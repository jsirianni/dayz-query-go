package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

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
		dayzClient, err := dayz.NewClient(logger, server.String(), dayz.WithTimeoutSeconds(10))
		if err != nil {
			logger.Error("new client", zap.Error(err))
			os.Exit(exitNewClientError)
		}
		clients = append(clients, dayzClient)
	}

	signalCtx, signalCancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer signalCancel()

	wg := sync.WaitGroup{}
	for _, dayzClient := range clients {
		wg.Add(1)
		go func(c *dayz.Client) {
			if err := dayzClient.Run(signalCtx); err != nil {
				logger.Error("client run", zap.Error(err))
			}
			wg.Done()
		}(dayzClient)
		logger.Info("client started")
	}

	<-signalCtx.Done()
	logger.Info("signal received, shutting down")
	wg.Wait()
	logger.Info("all clients stopped")
}
