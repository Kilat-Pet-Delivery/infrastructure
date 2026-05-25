package integration

import (
	"net/http"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestChatFlowOverGatewayWebSocket(t *testing.T) {
	requireIntegration(t)
	customerToken := requireCustomerToken(t)
	runnerToken := requireRunnerToken(t)
	threadID := requireEnv(t, "KILAT_CHAT_THREAD_ID")

	customerConn := dialChat(t, customerToken)
	defer customerConn.Close()
	runnerConn := dialChat(t, runnerToken)
	defer runnerConn.Close()

	subscribe := map[string]any{"type": "subscribe", "thread_id": threadID}
	if err := customerConn.WriteJSON(subscribe); err != nil {
		t.Fatalf("customer subscribe: %v", err)
	}
	if err := runnerConn.WriteJSON(subscribe); err != nil {
		t.Fatalf("runner subscribe: %v", err)
	}

	message := map[string]any{"type": "send_message", "thread_id": threadID, "body": "integration chat ping"}
	if err := customerConn.WriteJSON(message); err != nil {
		t.Fatalf("customer send: %v", err)
	}

	_ = runnerConn.SetReadDeadline(time.Now().Add(20 * time.Second))
	var received map[string]any
	if err := runnerConn.ReadJSON(&received); err != nil {
		t.Fatalf("runner receive chat event: %v", err)
	}
	if received["type"] == "" {
		t.Fatalf("chat event missing type: %#v", received)
	}

	resp := postJSON(t, runnerToken, "/api/v1/threads/"+threadID+"/read", map[string]any{})
	assert2xx(t, resp)
}

func dialChat(t *testing.T, token string) *websocket.Conn {
	t.Helper()
	header := http.Header{}
	conn, resp, err := websocket.DefaultDialer.Dial(gatewayWSURL()+"/ws/chat?token="+token, header)
	if err != nil {
		if resp != nil {
			t.Fatalf("dial chat websocket: %v (%s)", err, resp.Status)
		}
		t.Fatalf("dial chat websocket: %v", err)
	}
	return conn
}
