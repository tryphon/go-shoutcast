package shoutcast

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"sort"
	"strings"
	"time"
)

type Client struct {
	Host     string
	Password string
	Headers  map[string]string
	Timeout  time.Duration

	connection net.Conn
	readWriter *bufio.ReadWriter
}

func (client *Client) Connect() (connection net.Conn, err error) {
	err = client.Dial()
	if err != nil {
		return
	}

	err = client.SendPassword()
	if err != nil {
		return
	}

	err = client.SendHeaders()
	if err != nil {
		return
	}

	connection = client.connection
	return
}

func (client *Client) SendPassword() error {
	err := client.Write(client.Password)
	if err != nil {
		return err
	}
	client.Flush()

	var response string
	response, err = client.Read()
	if err != nil {
		return err
	}

	if !strings.HasPrefix(response, "OK") {
		return errors.New(fmt.Sprintf("Shoutcast server rejects request: %s", response))
	}

	return nil
}

func (client *Client) sortedHeaderNames() []string {
	var keys []string
	for key := range client.Headers {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func (client *Client) SendHeaders() error {
	for _, attribute := range client.sortedHeaderNames() {
		value := client.Headers[attribute]
		if err := client.SendHeader(attribute, value); err != nil {
			return err
		}
	}

	if err := client.Write(""); err != nil {
		return nil
	}
	return client.Flush()
}

func (client *Client) SendHeader(name, value string) error {
	return client.Write(fmt.Sprintf("%s:%s", name, value))
}

func (client *Client) Write(message string) error {
	_, err := client.readWriter.WriteString(fmt.Sprintf("%s\n", message))
	return err
}

func (client *Client) Read() (response string, err error) {
	response, err = client.readWriter.ReadString('\n')
	if err != nil {
		return
	}
	response = strings.TrimSpace(response)
	return
}

func (client *Client) Flush() error {
	return client.readWriter.Flush()
}

func (client *Client) Dial() error {
	if client.Timeout == 0 {
		client.Timeout = 10 * time.Second
	}
	connection, err := net.DialTimeout("tcp", client.Host, client.Timeout)
	if err != nil {
		return err
	}
	client.connection = connection
	client.readWriter = bufio.NewReadWriter(bufio.NewReader(connection), bufio.NewWriter(connection))
	return nil
}

func (client *Client) Close() {
	if client.connection != nil {
		client.connection.Close()
		client.connection = nil
	}
	client.readWriter = nil
}
