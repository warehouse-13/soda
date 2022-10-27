## soda

This is an educational repo to explain and reproduce [this bug][bug] in flintlock.

A blog write up exists [here][blog], but this is a repo with a little more explanation
about why we saw this error in Liquid Metal.

### Running the repro

Open a terminal with at least 3 panes. I like to use [tmux](https://github.com/tmux/tmux/wiki) ([my config][tmux]).

_nb: I am running on Linux, some commands may not work for you if you use macOS._

Clone this repo, build the binaries:
```bash
git clone https://github.com/warehouse-13/soda
cd soda
make build
```

In one window, start the service and the [pprof][service]:
```bash
./srv
```

In another, start watching for established connections on the server port `1430`:
```bash
watch -n 0.1 "netstat -a | awk '/:1430/ && /ESTABLISHED/' | wc -l"
```

Open your browser to http://localhost:1431/debug/pprof/.

In the last terminal window, start the client pointing it at the server:
```bash
./cli --address localhost:1430
```

In the client window you will see lots of random numbers being generated. These
are not important, they just show the server is responding to requests.

In the `netstat` window, you will see the number ticking up. How high it gets depends
on whatever open connection limit you have on your machine.

In the browser if you refresh it, you will see that the number of `goroutines` has
jumped by a couple of thousand.

Eventually you will see the client get stuck for a bit and then return this error:
```
could not make call rpc error: code = Unavailable desc = failed to receive server preface within timeout
```

### What does this mean?

This is what happens when you do not close client connections.

If you stop either the client or the server the connections will be closed. This is
fine for a test like this, or for short-lived programs which exit immediately after call,
but not ideal for long running services.

### Fix it

Stop the server and the client.

Open `client/main.go`.

Uncomment the following line:
```go
	// defer conn.Close()
```

Restart everything. This time you will see that connections and goroutines
stay at a reasonable level.

There are more notes in the code itself as well as some branches showing progression
of the solve.

[bug]: https://github.com/weaveworks-liquidmetal/flintlock/issues/503
[blog]: https://cbctl.dev/blog/close-grpc-connections
[tmux]: https://gist.github.com/Callisto13/b4cc217ca4f1c2f7f51405d62b941adb
[pprof]: https://pkg.go.dev/net/http/pprof
