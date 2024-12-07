package dayz

import (
	"bytes"
	"fmt"
)

// ServerInfo represents DayZ server information.
type ServerInfo struct {
	// TODO(jsirianni): Use correct types

	// ProtocolVersion is the protocol version of the server.
	ProtocolVersion string

	// ServerName is the name of the server.
	ServerName string

	// MapName is the name of the map.
	MapName string

	// GameDirectory is the directory of the game.
	GameDirectory string

	// AppID is the application ID of the game.
	AppID string

	// Players is the number of players connected to the server.
	Players string

	// MaxPlayers is the maximum number of players the server can hold.
	MaxPlayers string

	// Bots is the number of bots connected to the server.
	Bots string

	// ServerType is the type of server.
	ServerType string

	// OsType is the operating system of the server.
	OsType string

	// PasswordProtected is whether the server is password protected.
	PasswordProtected string

	// VacSecured is whether the server is VAC secured.
	VacSecured string

	// Version is the version of the server.
	Version string
}

func readNullTerminatedString(reader *bytes.Reader) (string, error) {
	var result []byte
	for {
		b, err := reader.ReadByte()
		if err != nil {
			return "", fmt.Errorf("reading null terminated string: %v", err)
		}
		if b == 0 {
			break
		}
		result = append(result, b)
	}
	return string(result), nil
}
