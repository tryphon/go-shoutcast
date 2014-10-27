# Go ShoutCast

Connect a ShoutCast server to send stream

# Usage

    client := shoutcast.Client{
	  Host:     "example.com:8000",
	  Password: "secret",
	  Headers: map[string]string{
		"icy-bt":  "96",
		"icy-pub": "1",
        "content-type": "audio/mpeg",
	  },
    }

    connection, err := client.Connect()
    if err != nil {
      // ...
    }

    // You can send stream data
    connection.Write(...)
