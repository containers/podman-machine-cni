package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
)

type Expose struct {
	Local  string `json:"local"`
	Remote string `json:"remote"`
}

type Unexpose struct {
	Local string `json:"local"`
}

// getPrimaryIP extracts the host's IP address from an environment
// variable. It is an error if that IP is blank
func getPrimaryIP() (net.IP, error) {
	hostIP := os.Getenv("PODMAN_MACHINE_HOST")
	if len(hostIP) < 1 {
		return nil, errors.New("invalid PODMAN_MACHINE_HOST environment variable")
	}
	addr := net.ParseIP(hostIP)
	return addr, nil
}

func getAPIEndpoint() string {
	// read a envar this is required for testing
	endpoint := os.Getenv("GVPROXY_REMOTE_ADDR")
	if endpoint != "" {
		return endpoint
	}
	return apiEndpoint
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
		return annotateResponseError(resp.Body)
	}
	return nil
}

func annotateResponseError(r io.Reader) error {
	b, err := ioutil.ReadAll(r)
	if err == nil && len(b) > 0 {
		return fmt.Errorf("something went wrong with the request: %q", string(b))
	}
	return errors.New("something went wrong with the request, could not read response")
}
