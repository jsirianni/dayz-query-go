// Package dayz provides a simple API for fetching
// DayZ server information.
package dayz

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
	"time"
)

// ClientOption represents a DayZ server query client option.
type ClientOption func(*Client) error

// WithTimeoutSeconds sets the timeout (in seconds) for the
// query client. If a timeout less than 1 set, a 30 second
// timeout is used.
func WithTimeoutSeconds(timeout int) ClientOption {
	return func(c *Client) error {
		if timeout < 1 {
			timeout = 30
		}
		t := time.Duration(timeout) * time.Second
		return c.conn.SetDeadline(time.Now().Add(t))
	}
}

// NewClient creates a new DayZ server query client. If
// there is an error connecting to the server, a nil client
// is returned with the error.
func NewClient(serverAddr string, opts ...ClientOption) (*Client, error) {
	conn, err := net.Dial("udp", serverAddr)
	if err != nil {
		return nil, err
	}

	for _, opt := range opts {
		if err := opt(&Client{conn}); err != nil {
			return nil, fmt.Errorf("configuring option: %v", err)
		}
	}

	return &Client{conn}, nil
}

// Client represents a DayZ server query client.
type Client struct {
	conn net.Conn
}

// Query sends a query to the DayZ server and returns the
// response. If there is an error sending the query or reading
// the response, nil bytes are returned with the error.
func (c *Client) Query(query []byte) ([]byte, error) {
	_, err := c.conn.Write(query)
	if err != nil {
		return nil, fmt.Errorf("sending query: %v", err)
	}

	buffer := make([]byte, 2048)
	n, err := c.conn.Read(buffer)
	if err != nil {
		return nil, fmt.Errorf("reading response: %v", err)
	}

	return buffer[:n], nil
}

// ServerInfo queries the DayZ server for information and
// returns the parsed server information. If there is an
// error querying the server or parsing the response, an
// empty ServerInfo struct is returned with the error.
func (c *Client) ServerInfo() (ServerInfo, error) {
	raw, err := c.serverInfo()
	if err != nil {
		return ServerInfo{}, fmt.Errorf("server info: %v", err)
	}

	info, err := c.parseServerInfo(raw)
	if err != nil {
		return ServerInfo{}, fmt.Errorf("parse server info: %v", err)
	}

	return info, nil
}

func (c *Client) serverInfo() ([]byte, error) {
	// Initial query to retrieve the challenge number.
	query := []byte{
		0xFF, 0xFF, 0xFF, 0xFF,
		'T', 'S', 'o', 'u', 'r',
		'c', 'e', ' ', 'E', 'n',
		'g', 'i', 'n', 'e', ' ',
		'Q', 'u', 'e', 'r', 'y', 0x00,
	}

	resp, err := c.Query(query)
	if err != nil {
		return nil, fmt.Errorf("initial query: %v", err)
	}

	// Check if this is a challenge response
	// 0x41 indicates a challenge response
	if len(resp) >= 5 && resp[4] == 0x41 {
		// Decode the challenge number (last 4 bytes)
		challenge := binary.LittleEndian.Uint32(resp[5:9])

		// Append the challenge number to the query
		challengeBytes := make([]byte, 4)
		binary.LittleEndian.PutUint32(challengeBytes, challenge)
		query = append(query, challengeBytes...)

		// Resend the query with the challenge number
		resp, err = c.Query(query)
		if err != nil {
			return nil, fmt.Errorf("resending query: %v", err)
		}

		return resp, nil
	}

	return nil, fmt.Errorf("unexpected response")
}

func (c *Client) parseServerInfo(raw []byte) (ServerInfo, error) {
	info := ServerInfo{}

	// 0x49 indicates a server info response
	if len(raw) < 5 || raw[4] != 0x49 {
		return info, fmt.Errorf("invalid raw server info response with length %d", len(raw))
	}

	// Skip the header bytes
	reader := bytes.NewReader(raw[5:])

	protocol, err := reader.ReadByte()
	if err != nil {
		return info, fmt.Errorf("reading protocol version: %v", err)
	}
	info.ProtocolVersion = fmt.Sprintf("%d", protocol)

	serverName, err := readNullTerminatedString(reader)
	if err != nil {
		return info, fmt.Errorf("reading server name: %v", err)
	}
	info.ServerName = serverName

	mapName, err := readNullTerminatedString(reader)
	if err != nil {
		return info, fmt.Errorf("reading map name: %v", err)
	}
	info.MapName = mapName

	gameDirectory, err := readNullTerminatedString(reader)
	if err != nil {
		return info, fmt.Errorf("reading game directory: %v", err)
	}
	info.GameDirectory = gameDirectory

	// TODO(jsirianni): This value seems to be empty every time,
	// therefor we are reading it but not using it.
	_, err = readNullTerminatedString(reader)
	if err != nil {
		return info, fmt.Errorf("reading game description: %v", err)
	}

	var appID uint16
	if err := binary.Read(reader, binary.LittleEndian, &appID); err != nil {
		return info, fmt.Errorf("reading app id: %v", err)
	}
	info.AppID = fmt.Sprintf("%d", appID)

	players, err := reader.ReadByte()
	if err != nil {
		return info, fmt.Errorf("reading player count: %v", err)
	}
	info.Players = fmt.Sprintf("%d", players)

	maxPlayers, err := reader.ReadByte()
	if err != nil {
		return info, fmt.Errorf("reading max player count: %v", err)
	}
	info.MaxPlayers = fmt.Sprintf("%d", maxPlayers)

	bots, err := reader.ReadByte()
	if err != nil {
		return info, fmt.Errorf("reading bot count: %v", err)
	}
	info.Bots = fmt.Sprintf("%d", bots)

	serverType, err := reader.ReadByte()
	if err != nil {
		return info, fmt.Errorf("reading server type: %v", err)
	}
	// TODO(jsirianni): Parse server type:
	// d - dedicated, l - listen, p - proxy
	// It is probably always 'd' for DayZ servers.
	info.ServerType = string(serverType)

	// TOOD(jsirianni): Parse OS type: w - Windows, l - Linux
	osType, err := reader.ReadByte()
	if err != nil {
		return info, fmt.Errorf("reading os type: %v", err)
	}
	info.OsType = string(osType)

	passwordProtected, err := reader.ReadByte()
	if err != nil {
		return info, fmt.Errorf("reading password protected: %v", err)
	}
	info.PasswordProtected = strconv.FormatBool(passwordProtected == 1)

	vacSecured, err := reader.ReadByte()
	if err != nil {
		return info, fmt.Errorf("reading vac secured: %v", err)
	}
	info.VacSecured = strconv.FormatBool(vacSecured == 1)

	serverVersion, err := readNullTerminatedString(reader)
	if err != nil {
		return info, fmt.Errorf("reading server version: %v", err)
	}
	info.Version = serverVersion

	return info, nil
}

// TODO(jsirianni): Implement ModList.
// func (c *Client) ModList() (string, error) {
// 	query := []byte{0xFF, 0xFF, 0xFF, 0xFF, 'V'}

// 	resp, err := c.Query(query)
// 	if err != nil {
// 		return "", fmt.Errorf("initial query: %v", err)
// 	}

// 	// Check if this is a challenge response
// 	// 0x41 indicates a challenge response
// 	if len(resp) >= 5 && resp[4] == 0x41 {
// 		// Decode the challenge number (last 4 bytes)
// 		challenge := binary.LittleEndian.Uint32(resp[5:9])

// 		// Append the challenge number to the query
// 		challengeBytes := make([]byte, 4)
// 		binary.LittleEndian.PutUint32(challengeBytes, challenge)
// 		query = append(query, challengeBytes...)

// 		// Resend the query with the challenge number
// 		resp, err = c.Query(query)
// 		if err != nil {
// 			return "", fmt.Errorf("resending query: %v", err)
// 		}

// 		return string(resp), nil
// 	}

// 	return "", fmt.Errorf("unexpected response")
// }
