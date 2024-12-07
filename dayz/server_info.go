package dayz

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
)

func ParseServerInfo(raw []byte) (ServerInfo, error) {
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

type ServerInfo struct {
	// TODO(jsirianni): Use correct types

	ProtocolVersion   string
	ServerName        string
	MapName           string
	GameDirectory     string
	AppID             string
	Players           string
	MaxPlayers        string
	Bots              string
	ServerType        string
	OsType            string
	PasswordProtected string
	VacSecured        string
	Version           string
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
