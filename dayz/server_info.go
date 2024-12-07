package dayz

import (
	"bytes"
	"encoding/binary"
	"fmt"
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

	// Read strings (null-terminated)
	info.ServerName = readNullTerminatedString(reader)
	info.MapName = readNullTerminatedString(reader)
	info.GameDirectory = readNullTerminatedString(reader)
	info.GameDescription = readNullTerminatedString(reader)

	binary.Read(reader, binary.LittleEndian, &info.AppID)

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

	return info, nil
}

type ServerInfo struct {
	// TODO(jsirianni): Use correct types

	ProtocolVersion string
	ServerName      string
	MapName         string
	GameDirectory   string
	GameDescription string
	AppID           string
	Players         string
	MaxPlayers      string
	Bots            string
}

func readNullTerminatedString(reader *bytes.Reader) string {
	var result []byte
	for {
		b, err := reader.ReadByte()
		if err != nil || b == 0 {
			break
		}
		result = append(result, b)
	}
	return string(result)
}
