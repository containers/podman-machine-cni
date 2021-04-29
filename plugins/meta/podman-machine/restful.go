package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"

	"github.com/sirupsen/logrus"
)

type Expose struct {
	Local  string `json:"local"`
	Remote string `json:"remote"`
}

type Unexpose struct {
	Local string `json:"local"`
}

func getPrimaryIP() (net.IP, error) {
	// no connection is actually made here
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := conn.Close(); err != nil {
			logrus.Error(err)
		}
	}()
	addr := conn.LocalAddr().(*net.UDPAddr)
	return addr.IP, nil
}

func postRequest(ctx context.Context, url *url.URL, body interface{}) error {
	var buf io.ReadWriter
	client := &http.Client{}
	buf = new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(body); err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url.String(), buf)
	if err != nil {
		return err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New("something went wrong with the request")
	}
	return nil
}
