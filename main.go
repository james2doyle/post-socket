package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "net"
    "net/http"
    "runtime"

    "github.com/BurntSushi/toml"
    "github.com/gobwas/ws"
    "github.com/gobwas/ws/wsutil"
)

type (
    appConfig struct {
        WebhookPort   int
        WebsocketPort int
        Logging       bool
    }
)

var conf appConfig
var clients []*net.Conn

func logger(format string, args ...interface{}) {
    if conf.Logging {
        log.Printf(format, args...)
    }
}

func addCorsHeader(w http.ResponseWriter) {
    headers := w.Header()
    headers.Add("Access-Control-Allow-Origin", "*")
    headers.Add("Vary", "Origin")
    headers.Add("Vary", "Access-Control-Request-Method")
    headers.Add("Vary", "Access-Control-Request-Headers")
    headers.Add("Access-Control-Allow-Headers", "Content-Type, Origin, Accept, token")
    headers.Add("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
}

func handleWebhook() http.HandlerFunc {
    logger("webhook server started on %d", conf.WebhookPort)
    return func(w http.ResponseWriter, r *http.Request) {
        // handle preflight requests
        if r.Method == "OPTIONS" {
            addCorsHeader(w)
            w.WriteHeader(http.StatusOK)
            return
        }

        if r.Method != "POST" {
            logger("%s", "method not post")
            http.Error(w, "method must be POST", http.StatusMethodNotAllowed)
            return
        }

        // Read body
        msg, err := ioutil.ReadAll(r.Body)
        if msg == nil {
            return
        }
        logger("post body: %s", msg)
        defer r.Body.Close()
        if err != nil {
            http.Error(w, err.Error(), 500)
            return
        }

        for _, client := range clients {
            err = wsutil.WriteServerMessage(*client, ws.OpText, msg)
            if err != nil {
                // TODO: handle closed clients
                logger("webhook write socket error: %s", err)
                // clients[index] = clients[len(clients)-1]
                // w.WriteHeader(http.StatusInternalServerError)
                // fmt.Fprintf(w, "error sending websocket message for client: %v", err)
                // return
            }
        }

        w.WriteHeader(http.StatusOK)
        fmt.Fprint(w, "sent websocket message for client")
    }
}

func handleWebsocket() http.HandlerFunc {
    logger("websocket server started on %d", conf.WebsocketPort)
    return func(w http.ResponseWriter, r *http.Request) {
        conn, _, _, err := ws.UpgradeHTTP(r, w)
        if err != nil {
            logger("%s", err)
        }

        // queue up the connections
        clients = append(clients, &conn)

        go func() {
            defer conn.Close()

            for {
                msg, op, err := wsutil.ReadClientData(conn)
                if msg == nil {
                    break
                }
                if err != nil {
                    logger("%s", err)
                }

                for _, client := range clients {
                    err = wsutil.WriteServerMessage(*client, op, msg)
                    if err != nil {
                        logger("write socket error: %s", err)
                    }
                }
            }
        }()
    }
}

// RunWebhookServer will run the Webhook server in a goroutine
func RunWebhookServer() {
    // run the server that accepts the webhooks
    go func() {
        http.ListenAndServe(fmt.Sprintf(":%d", conf.WebhookPort), handleWebhook())
    }()
}

// RunWebsocketServer will run the Websocket server in a goroutine
func RunWebsocketServer() {
    // run the server that providers the websocket
    go func() {
        http.ListenAndServe(fmt.Sprintf(":%d", conf.WebsocketPort), handleWebsocket())
    }()
}

func main() {
    if _, err := toml.DecodeFile("./config.toml", &conf); err != nil {
        log.Fatal(err)
    }

    RunWebhookServer()
    RunWebsocketServer()

    runtime.Goexit()

    fmt.Println("Exit")
}
