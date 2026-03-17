package adapters

import (
	"agent-in-go/pkg/session"
	"fmt"
	"net/http"

	"golang.org/x/net/websocket"
)

type WSAdapter struct {
	port  string
	store *session.SessionStore
}

func NewWSAdapter(port string, store *session.SessionStore) *WSAdapter {
	return &WSAdapter{port: port, store: store}
}

func (a *WSAdapter) Name() string { return "WebSocket" }

func (a *WSAdapter) Start() error {
	http.Handle("/ws", websocket.Handler(func(conn *websocket.Conn) {
		sessionID := conn.Request().URL.Query().Get("session_id")
		if sessionID == "" {
			sessionID = conn.RemoteAddr().String()
		}
		for {
			var question string
			if err := websocket.Message.Receive(conn, &question); err != nil {
				break
			}
			answer := a.store.Ask(sessionID, question) // ← store.Ask, not sess.Ask
			websocket.Message.Send(conn, answer)
		}
	}))
	fmt.Printf("[WS] listening on :%s\n", a.port)
	return http.ListenAndServe(":"+a.port, nil)
}
