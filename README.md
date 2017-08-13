# hologram UDP gateway

_UDP gateway for hologram cell service_

The hologw is a very simple gateway between UDP and the hologram IoT
service in order to submit data to Hologram.
The cellular device sends a UDP packet to port 9999 on which hologw runs,
the hologw forwards that packet via TCP to Hologram, reads the response, and
sends that back via UDP to the cellular device.

One added feature is that if the JSON message object contains a field
named "a" that is an integer then the GW delays the UDP response by
that many seconds. This provided to be able to test the UDP NAT timeout
of the cellular carrier.

There are no command-line options, so just run `./hologw` on an internet-accessible
server. Assuming your cellular device is a unix box, send it a packet
like so:
```
echo '{"a":10, "k":"XXXXXXXX","d":"Hello, World!","t":"test"}' | nc -q 20 -u your.host.or.ip 9999
```

This will submit the whole message to Hologram (which doesn't seem to mind the extra "a":10,
get a "[0,0]" response, then delay 10 seconds, and finally send it via UDP.
The `-q 20` causes nc to linger 20 seconds for the response.

The `hologw` binary is a Linux x64 executable that will run on any Linux distro.
