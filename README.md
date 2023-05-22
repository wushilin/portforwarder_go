# portforwarder_go
Golang version of port forwarder

# Build

```bash
$ go build .
```

# Running
```bash
$ ./portforwarder_go -b :443::www.google.com:443 -b :22::www.remote.com:22
```

# Log rotation
Please see github.com/wushilin/logd

