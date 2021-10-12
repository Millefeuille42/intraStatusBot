package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
)

func makeReqOnChannel(url string, respCh chan *http.Response, errCh chan error) {
	resp, err := http.Get(url)
	if err != nil {
		errCh <- err
		return
	}
	respCh <- resp
}

func makeReqWithTimeout(url string) (*http.Response, error) {
	respCh := make(chan *http.Response, 1)
	errCh := make(chan error, 1)

	go makeReqOnChannel(url, respCh, errCh)

	select {
	case resp := <- respCh:
		return resp, nil
	case err := <- errCh:
		return nil, err
	case <- time.After(3 * time.Second):
		return nil, errors.New("timeout")
	}
}

func commandSendIntraStatus(agent discordAgent) {
	response, err := makeReqWithTimeout("https://intra.42.fr/pommedeterre")
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "ERROR:", err)
		return
	}
	if response.StatusCode >= 200 || response.StatusCode < 400 {
		_, _ = agent.session.ChannelMessageSend(agent.)
		return
	}
	fmt.Println("As usual...")
}
