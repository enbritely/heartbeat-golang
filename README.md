# Heartbeat library
[![Build Status](https://travis-ci.org/enbritely/heartbeat-golang.svg)](https://travis-ci.org/enbritely/heartbeat-golang)

## Install
```bash
go get github.com/enbritely/heartbeat-golang
```

## Usage
Include the following in your code to run the heartbeat service on your `10101` port:

```go
go RunHeartbeatService(":10101")
```

If you query your application on port `10101` you get the heartbeat signal as JSON message.
```bash
~ $ curl localhost:10101/heartbeat
{"status":"running","build":"767b1c413b3187be28d373e3c9d3f6be02451785","uptime":"10.245384286s"}
```

To have build number, build your go program with the following option:

```console
go build --ldflags="-X github.com/enbritely/heartbeat-golang.CommitHash `git rev-parse HEAD`"
```
