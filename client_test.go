package shoutcast

import (
	"bufio"
	"bytes"
	"testing"
)

func testClient() (*Client, *bytes.Buffer, *bytes.Buffer) {
	var input, output bytes.Buffer

	client := &Client{
		Host:     "example.com:8000",
		Password: "secret",
		Headers: map[string]string{
			"icy-bt":  "96",
			"icy-pub": "1",
		},
		readWriter: bufio.NewReadWriter(bufio.NewReader(&input), bufio.NewWriter(&output)),
	}
	return client, &output, &input
}

func TestClient_Write(t *testing.T) {
	client, output, _ := testClient()
	client.Write("dummy")
	client.Flush()
	if output.String() != "dummy\n" {
		t.Errorf("Wrong result :\n got: %v\nwant: %v", output.String(), "dummy\n")
	}
}

func TestClient_SendHeader(t *testing.T) {
	client, output, _ := testClient()
	client.SendHeader("name", "value")
	client.Flush()
	if output.String() != "name:value\n" {
		t.Errorf("Wrong result :\n got: %v\nwant: %v", output.String(), "name:value\n")
	}
}

func TestClient_Read(t *testing.T) {
	client, _, input := testClient()
	input.WriteString("OK2\n")

	response, _ := client.Read()
	if response != "OK2" {
		t.Errorf("Wrong response :\n got: %v\nwant: %v", response, "OK2")
	}
}

func TestClient_SendPassword(t *testing.T) {
	client, output, input := testClient()
	input.WriteString("OK2\n")

	err := client.SendPassword()
	if err != nil {
		t.Fatal(err)
	}

	if output.String() != "secret\n" {
		t.Errorf("Wrong output :\n got: %v\nwant: %v", output.String(), "secret\n")
	}
}

func TestClient_SendHeaders(t *testing.T) {
	client, output, _ := testClient()
	client.SendHeaders()
	client.Flush()

	if output.String() != "icy-bt:96\nicy-pub:1\n\n" {
		t.Errorf("Wrong output :\n got: '%v'\nwant: '%v'", output.String(), "icy-bt:96\nicy-pub:1\n\n")
	}
}
