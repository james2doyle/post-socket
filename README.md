# POST Socket

> Turn a POST request (JSON) to Websocket event

### Config

```toml
# what port should the webhook web server run on?
WebhookPort = 8080
# what port should the websocket run on?
WebsocketPort = 8081
Logging = true
```

### Building with Docker

*Be sure to update the `EXPOSE` to follow the ports you setup in `config.toml`*

* `docker build -t post-socket .`
* `docker run post-socket:latest`

### Building with "go build"

* `go build`
* `./post-socket`

*With "go run":*

* `go run main.go`

### Deploying

You can use [now.sh](https://zeit.co/now) to deploy Docker apps. You can use their initial service for free. However, there is a limit on the number of deploys and also you are assigned random URL for each deploy. So this would not be a good option for long term usage.

Use `now` to deploy.

### Example

Open `example.html` on a webserver while the `post-socket` is running.

You can connect to the socket and send POST requests to the webhook server and see them come up on the socket connection.
