// Package dayz provides a simple API for fetching
// DayZ server information.
package dayz

import (
	"encoding/binary"
	"fmt"
	"net"
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

func (c *Client) ServerInfo() ([]byte, error) {
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
