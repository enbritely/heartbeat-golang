# Heartbeat library

## Usage

Include the following in your code to run the heartbeat service on your `10101` port:

```go
go RunHeartbeatService(":10101")
```

To have build number, build your go program with the following option:

```console
go build --ldflags="-X github.com/enbritely/heartbeat-golang.CommitHash `git rev-parse HEAD`"
```
