package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/plugins/pkg/testutils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("machine tests", func() {
	// this plugin needs no netns so we just set a path
	// this allows us to run the test rootless
	bogusNetns := "/some/path"
	const IFNAME string = "dummy0"

	var gvproxyCommand *exec.Cmd

	BeforeEach(func() {
		// start gvproxy in NETNS
		gvproxyCommand = exec.Command(os.Getenv("GVPROXY_PATH"), "--listen", "tcp://127.0.0.1:7777")
		gvproxyCommand.Stderr = os.Stderr
		gvproxyCommand.Stdout = os.Stdout
		err := gvproxyCommand.Start()
		Expect(err).NotTo(HaveOccurred())
		// wait 1 second to give gvproxy some time to come up
		time.Sleep(1 * time.Second)

		// set required vars for testing
		os.Setenv("GVPROXY_REMOTE_ADDR", "127.0.0.1")
		os.Setenv("PODMAN_MACHINE_HOST", "127.0.0.1")
	})

	AfterEach(func() {
		err := gvproxyCommand.Process.Kill()
		Expect(err).NotTo(HaveOccurred())
	})

	for _, ver := range []string{"0.3.0", "0.3.1", "0.4.0", "1.0.0"} {
		ver := ver
		It(fmt.Sprintf("podman machine with ports v[%s]", ver), func() {
			fullConf := []byte(fmt.Sprintf(`{
				"cniVersion": "%s",
				"name": "test",
				"type": "podman-machine",
				"runtimeConfig": {
					"portMappings": [
				  		{ "hostPort": 9999, "containerPort": 80, "protocol": "tcp"}
					]
				},
				"prevResult": {
				  "interfaces": [
					{
					  "name": "dummy0",
					  "mac": "a6:a7:ca:6b:34:2e"
					},
					{
					  "name": "vetha0a83b38",
					  "mac": "9a:45:bd:b0:2d:dd"
					},
					{
					  "name": "eth0",
					  "mac": "ea:63:0e:63:3e:86",
					  "sandbox": "/var/run/netns/baude"
					}
				  ],
				  "ips": [
					{
					  "version": "4",
					  "interface": 2,
					  "address": "10.88.8.5/24",
					  "gateway": "10.88.8.1"
					}
				  ],
				  "routes": [
					{
					  "dst": "0.0.0.0/0"
					}
				  ]
			  }
				  }`, ver))

			args := &skel.CmdArgs{
				ContainerID: "dummy",
				Netns:       bogusNetns,
				IfName:      IFNAME,
				StdinData:   fullConf,
			}

			r, _, err := testutils.CmdAdd(bogusNetns, args.ContainerID, IFNAME, fullConf, func() error {
				return cmdAdd(args)
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(r.Version()).To(Equal(ver))

			err = testutils.CmdCheck(bogusNetns, args.ContainerID, IFNAME, fullConf, func() error {
				return cmdCheck(args)
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(r.Version()).To(Equal(ver))

			err = testutils.CmdDel(bogusNetns, args.ContainerID, IFNAME, func() error {
				return cmdDel(args)
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(r.Version()).To(Equal(ver))
		})

		It(fmt.Sprintf("podman machine with no ports v[%s]", ver), func() {
			fullConf := []byte(fmt.Sprintf(`{
				"cniVersion": "%s",
				"name": "test",
				"type": "podman-machine",
				"prevResult": {
				  "interfaces": [
					{
					  "name": "dummy0",
					  "mac": "a6:a7:ca:6b:34:2e"
					},
					{
					  "name": "vetha0a83b38",
					  "mac": "9a:45:bd:b0:2d:dd"
					},
					{
					  "name": "eth0",
					  "mac": "ea:63:0e:63:3e:86",
					  "sandbox": "/var/run/netns/baude"
					}
				  ],
				  "ips": [
					{
					  "version": "4",
					  "interface": 2,
					  "address": "10.88.8.5/24",
					  "gateway": "10.88.8.1"
					}
				  ],
				  "routes": [
					{
					  "dst": "0.0.0.0/0"
					}
				  ]
			  }
				  }`, ver))

			args := &skel.CmdArgs{
				ContainerID: "dummy",
				Netns:       bogusNetns,
				IfName:      IFNAME,
				StdinData:   fullConf,
			}

			r, _, err := testutils.CmdAdd(bogusNetns, args.ContainerID, IFNAME, fullConf, func() error {
				return cmdAdd(args)
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(r.Version()).To(Equal(ver))

			err = testutils.CmdCheck(bogusNetns, args.ContainerID, IFNAME, fullConf, func() error {
				return cmdCheck(args)
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(r.Version()).To(Equal(ver))

			err = testutils.CmdDel(bogusNetns, args.ContainerID, IFNAME, func() error {
				return cmdDel(args)
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(r.Version()).To(Equal(ver))
		})
	}
})
