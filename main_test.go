package main

import (
  "fmt"
  "net/http"
  "os"
  "strings"
  "testing"

  "github.com/gorilla/websocket"
  . "github.com/smartystreets/goconvey/convey"
)

func TestMain(m *testing.M) {
  // setup the config
  conf = appConfig{
    WebhookPort:   8080,
    WebsocketPort: 8081,
    Logging:       false,
  }

  RunWebhookServer()
  RunWebsocketServer()

  // call flag.Parse() here if TestMain uses flags
  os.Exit(m.Run())
}

func Test_Preflight_Request(t *testing.T) {
  Convey("It can accept a post request and convert it to a websocket message", t, func() {
    defer func() {
      So(recover(), ShouldBeNil)
    }()

    webhookURL := fmt.Sprintf("%s://%s:%d", "http", "localhost", conf.WebhookPort)

    // send the post request
    req, err := http.NewRequest(http.MethodOptions, webhookURL, nil)
    So(err, ShouldBeNil)
    client := &http.Client{}
    resp, err := client.Do(req)
    So(err, ShouldBeNil)
    So(resp.StatusCode, ShouldEqual, http.StatusOK)
    So(resp.Header.Get("Access-Control-Allow-Origin"), ShouldEqual, "*")
    So(resp.Header.Get("Access-Control-Allow-Headers"), ShouldEqual, "Content-Type, Origin, Accept, token")
    So(resp.Header.Get("Access-Control-Allow-Methods"), ShouldEqual, "GET,POST,OPTIONS")
  })
}

func Test_Post_Request_And_Websocket_Event(t *testing.T) {
  Convey("It can accept a post request and convert it to a websocket message", t, func() {
    defer func() {
      So(recover(), ShouldBeNil)
    }()

    websocketURL := fmt.Sprintf("%s://%s:%d", "ws", "localhost", conf.WebsocketPort)
    // Connect to the server over the websocket
    ws, _, err := websocket.DefaultDialer.Dial(websocketURL, nil)
    So(err, ShouldBeNil)
    defer ws.Close()

    webhookURL := fmt.Sprintf("%s://%s:%d", "http", "localhost", conf.WebhookPort)
    message := `{"type": "message", "message": "1536452269000"}`

    // send the post request
    resp, err := http.Post(webhookURL, "application/json", strings.NewReader(message))
    So(err, ShouldBeNil)
    So(resp.StatusCode, ShouldEqual, http.StatusAccepted)

    // Read response and check to see if it's what we expect.
    _, p, err := ws.ReadMessage()
    So(err, ShouldBeNil)
    So(string(p), ShouldEqual, message)
  })
}
