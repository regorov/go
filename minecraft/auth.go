package minecraft

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func Login(username string, password string, debugWriter io.Writer) (client *Client, err error) {
	params := url.Values{
		"user":     {username},
		"password": {password},
		"version":  {"13"},
	}

	if debugWriter != nil {
		fmt.Fprintf(debugWriter, "Logging in to minecraft.net as '%s'\n", username)
		fmt.Fprintf(debugWriter, "POST https://login.minecraft.net username=%s&password=...&version=13\n", username)
	}

	resp, err := http.PostForm("https://login.minecraft.net", params)
	if err != nil {
		return nil, err
	}

	respdata, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	respparts := strings.Split(string(respdata), ":")

	client = newClient(respparts[2], respparts[3], debugWriter)

	if debugWriter != nil {
		fmt.Fprintf(debugWriter, "Session ID: %s\n\n", client.sessionId)
	}

	return client, nil
}

func LoginOffline(username string, debugWriter io.Writer) (client *Client) {
	return newClient(username, "", debugWriter)
}

// Runs in the background, sending a keep-alive request to login.minecraft.net every 5 minutes.
func (client *Client) HTTPKeepAlive() {
	params := url.Values{
		"name":    {client.username},
		"session": {client.sessionId},
	}

	sessionURL := "https://login.minecraft.net/session?" + params.Encode()
	ticker := time.NewTicker(Tick * 6000)

	for {
		select {
		case <-ticker.C:
			if client.DebugWriter != nil {
				fmt.Fprintf(client.DebugWriter, "(HTTP keep-alive) GET %s\n", sessionURL)
			}

			resp, err := http.Get(sessionURL)
			if err != nil {
				client.ErrChan <- err
			}

			resp.Body.Close()

		case <-client.stopHTTPKeepAlive:
			client.stopHTTPKeepAlive <- struct{}{}
			return
		}
	}
}

func (client *Client) Logout() {
	client.stopHTTPKeepAlive <- struct{}{}
	<-client.stopHTTPKeepAlive
}
