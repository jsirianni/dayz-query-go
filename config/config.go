package config

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

const (
	// EnvServerList is the environment variable that
	// contains the list of DayZ server hostnames or IP
	// addresses to monitor. The list should be comma
	// separated, in the form of "host1:port1,host2:port2".
	// e.g: 50.108.116.1:2324,50.108.116.1:5700
	EnvServerList = "DAYZ_SERVER_LIST"
)

// ServerEndpoint is the endpoint of the DayZ server.
type ServerEndpoint struct {
	host string
	port string
}

// String returns the string representation of the server endpoint
// in the form of "host:port".
func (s ServerEndpoint) String() string {
	return net.JoinHostPort(s.host, s.port)
}

// New returns a new configuration
func New(logger *zap.Logger) (*Config, error) {
	c, err := ReadEnv()
	if err != nil {
		return nil, fmt.Errorf("read env: %w", err)
	}

	return c, nil
}

// Config is the configuration for the DayZ server monitor.
type Config struct {
	// Logger is a zap logger
	Logger *zap.Logger

	// ServerList is the list of DayZ server hostnames or IP
	ServerList []ServerEndpoint
}

// ReadEnv reads the configuration from the environment variables
// and returns the configuration. If there is an error reading
// the environment or the configuration is invalid, it returns
// an error with a nil configuration. The configuration will
// never be nil if there is no error.
func ReadEnv() (*Config, error) {
	c := &Config{}

	serverList, err := readServerList()
	if err != nil {
		return nil, fmt.Errorf("read server list: %w", err)
	}
	c.ServerList = serverList

	return c, nil
}

func readServerList() ([]ServerEndpoint, error) {
	v, ok := os.LookupEnv(EnvServerList)
	if !ok {
		return nil, fmt.Errorf("%s is a required option", EnvServerList)
	}

	if v == "" {
		return nil, fmt.Errorf("%s is empty", EnvServerList)
	}

	endpoints := strings.Split(v, ",")
	serverEndPoints := make([]ServerEndpoint, len(endpoints))

	for i, endpoint := range endpoints {
		host, port, err := net.SplitHostPort(endpoint)
		if err != nil {
			return nil, fmt.Errorf("invalid server endpoint %q while reading %s: %v", endpoint, EnvServerList, err)
		}

		if host == "" {
			return nil, fmt.Errorf("invalid server endpoint %q while reading %s: host is empty", endpoint, EnvServerList)
		}

		if port == "" {
			return nil, fmt.Errorf("invalid server endpoint %q while reading %s: port is empty", endpoint, EnvServerList)
		}

		if _, err := strconv.Atoi(port); err != nil {
			return nil, fmt.Errorf("invalid server endpoint %q while reading %s: port is not a number", endpoint, EnvServerList)
		}

		serverEndPoints[i] = ServerEndpoint{
			host: host,
			port: port,
		}
	}

	return serverEndPoints, nil
}
