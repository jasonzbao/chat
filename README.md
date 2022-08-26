### SETUP

1. Install golang 1.18
2. run `docker-compose up`
3. `cd go_client; go run main.go`


### About
A distributed, linearly scalable chat agent with support for /exit, /name, /color, and /private.

Client connects to server via a websocket and listens to events in the "all" channel and their personal channel.

Events are passed through a Redis pub/sub bus that routes requests to the right listeners.

Non-private messages are persisted into postgres. Ideally, this would be replaced by a No-SQL DB and for maximum throughput, would be batched per server.
