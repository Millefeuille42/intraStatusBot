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
	case resp := <-respCh:
		return resp, nil
	case err := <-errCh:
		return nil, err
	case <-time.After(3 * time.Second):
		return nil, errors.New("timeout")
	}
}

func commandSendIntraStatus(agent discordAgent) {
	t := time.Now()
	response, err := makeReqWithTimeout("https://intra.42.fr/")
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "ERROR:", err)
		return
	}
	responseTime := fmt.Sprintf("%d", time.Since(t).Milliseconds())
	agent.channel = agent.message.ChannelID
	if response.StatusCode >= 200 && response.StatusCode < 400 {
		sendMessageWithMention("Intra seems up! "+
			"("+responseTime+"ms)", "", agent)
		return
	}
	sendMessageWithMention("Intra seems down... As usual!", "", agent)
}
