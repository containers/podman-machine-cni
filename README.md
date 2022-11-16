# Deprecated

This repository is no longer maintained. The core logic has been merged into [podman as of v4.0](https://github.com/containers/podman/pull/12283).

# podman-machine-cni

This plugin collects the port information of the container and sends information to a server on the host
operating system.  The information is used by the server to open and close port mappings on the host. It
is only meant to be used in a podman-machine virtual machine.  The plugin can
be enabled with the following stanza:

```
      {
         "type": "podman-machine",
         "capabilities": {
            "portMappings": true
         }
      },

```

The server in question is gvisor-tap-vsock.  The plugin connects to the server via RESTful API calls on
container start and stop (or die).  The plugin converts the port data information into a JSON payload
for the API endpoint.  On container start, ports on the host are opened and mapped; on stop, they are closed.
