// Copyright 2021 authors
// Copyright 2017 CNI authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/types"
	"github.com/containernetworking/cni/pkg/version"
	bv "github.com/containernetworking/plugins/pkg/utils/buildversion"
	"github.com/pkg/errors"
)

func cmdAdd(args *skel.CmdArgs) error {
	// Get the port information from the chained plugin
	portMaps, err := parseConfig(args.StdinData, args.Args)
	if err != nil {
		return errors.Wrap(err, "failed to parse config")
	}
	hostAddr, err := getPrimaryIP()
	if err != nil {
		return err
	}
	// No portmappings, do nothing
	if len(portMaps.RuntimeConfig.PortMaps) < 1 {
		return nil
	}
	// Iterate and send requests to the server
	for _, pm := range portMaps.RuntimeConfig.PortMaps {
		hostPort := strconv.Itoa(pm.HostPort)
		u, err := url.Parse(fmt.Sprintf("http://%s:%s/services/forwarder/expose", apiEndpoint, apiEndpointPort))
		if err != nil {
			return err
		}
		expose := Expose{
			Local:  fmt.Sprintf("%s:%s", "0.0.0.0", hostPort),
			Remote: fmt.Sprintf("%s:%s", hostAddr.String(), hostPort),
		}
		if err := postRequest(context.Background(), u, expose); err != nil {
			return err
		}
	}

	// Have to do this for chained plugins, which this is
	return types.PrintResult(portMaps.PrevResult, portMaps.CNIVersion)
}

func cmdDel(args *skel.CmdArgs) error {
	portMaps, err := parseConfig(args.StdinData, args.Args)
	if err != nil {
		return errors.Wrap(err, "failed to parse config")
	}
	// No portmappings, do nothing
	if len(portMaps.RuntimeConfig.PortMaps) < 1 {
		return nil
	}
	for _, pm := range portMaps.RuntimeConfig.PortMaps {
		hostPort := strconv.Itoa(pm.HostPort)
		u, err := url.Parse(fmt.Sprintf("http://%s:%s/services/forwarder/unexpose", apiEndpoint, apiEndpointPort))
		if err != nil {
			return err
		}
		unexpose := Unexpose{Local: fmt.Sprintf("0.0.0.0:%s", hostPort)}
		if err := postRequest(context.Background(), u, unexpose); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	skel.PluginMain(cmdAdd, cmdCheck, cmdDel, version.All, bv.BuildString("machine"))
}

func cmdCheck(args *skel.CmdArgs) error {
	client := &http.Client{}
	u, err := url.Parse(fmt.Sprintf("http://%s:%s/status", apiEndpoint, apiEndpointPort))
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, u.String(), nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New("something went wrong with the request")
	}
	return nil
}

// parseConfig parses the supplied configuration (and prevResult) from stdin.
func parseConfig(stdin []byte, args string) (*PortMapConf, error) {
	conf := PortMapConf{}
	if err := json.Unmarshal(stdin, &conf); err != nil {
		return nil, fmt.Errorf("failed to parse network configuration: %v", err)
	}
	// Parse previous result.
	if conf.RawPrevResult != nil {
		if err := version.ParsePrevResult(&conf.NetConf); err != nil {
			return nil, errors.Wrap(err, "could not parse prevResult")
		}
	}
	return &conf, nil
}
