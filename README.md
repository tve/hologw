# hologram UDP gateway

_UDP gateway for hologram cell service_

The hologw will someday implement a simple gateway between UDP and the hologram IoT
service. Right now all it does is listen on UDP port 9999, print the contents of
incoming packets, and respond with a single packet containing `OK\n`.

A simple way to send it packets from a Linux device is to use something like:
```
echo '{"k":"=XXXXXXX","d":"Hello, World!","t":"test"}' | nc -q 10 -u your.host.or.ip 9999
```

The `-q 10` causes nc to linger 10 seconds for the OK response.
