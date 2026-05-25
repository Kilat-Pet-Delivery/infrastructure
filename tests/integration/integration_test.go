package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/segmentio/kafka-go"
)

const defaultTimeout = 20 * time.Second

func requireIntegration(t *testing.T) {
	t.Helper()
	if os.Getenv("KILAT_RUN_INTEGRATION") != "1" {
		t.Skip("set KILAT_RUN_INTEGRATION=1 to run docker-compose integration tests")
	}
}

func gatewayURL() string {
	if value := os.Getenv("KILAT_GATEWAY_URL"); value != "" {
		return strings.TrimRight(value, "/")
	}
	return "http://localhost:8080"
}

func gatewayWSURL() string {
	if value := os.Getenv("KILAT_GATEWAY_WS_URL"); value != "" {
		return strings.TrimRight(value, "/")
	}
	return "ws://localhost:8080"
}

func tokenFromEnvOrLogin(t *testing.T, key, email, password string) string {
	t.Helper()
	if token := os.Getenv(key); token != "" {
		return token
	}
	if email == "" || password == "" {
		t.Skipf("%s is not set and no seeded credentials were supplied", key)
	}

	body := map[string]string{"email": email, "password": password}
	resp := postJSON(t, "", "/api/v1/auth/login", body)
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		payload, _ := io.ReadAll(resp.Body)
		t.Skipf("login for %s failed with %s: %s", email, resp.Status, string(payload))
	}

	var parsed map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		t.Fatalf("decode login response: %v", err)
	}
	if token := findString(parsed, "access_token"); token != "" {
		return token
	}
	if token := findString(parsed, "token"); token != "" {
		return token
	}
	t.Fatalf("login response did not include access token")
	return ""
}

func requireCustomerToken(t *testing.T) string {
	return tokenFromEnvOrLogin(t, "KILAT_CUSTOMER_TOKEN", os.Getenv("KILAT_CUSTOMER_EMAIL"), os.Getenv("KILAT_CUSTOMER_PASSWORD"))
}

func requireRunnerToken(t *testing.T) string {
	return tokenFromEnvOrLogin(t, "KILAT_RUNNER_TOKEN", envDefault("KILAT_RUNNER_EMAIL", "runner.test@kilat.my"), envDefault("KILAT_RUNNER_PASSWORD", "TestRunner123!"))
}

func requireAgentToken(t *testing.T) string {
	return tokenFromEnvOrLogin(t, "KILAT_AGENT_TOKEN", os.Getenv("KILAT_AGENT_EMAIL"), os.Getenv("KILAT_AGENT_PASSWORD"))
}

func getJSON(t *testing.T, token, path string) *http.Response {
	t.Helper()
	return doJSON(t, http.MethodGet, token, path, nil)
}

func postJSON(t *testing.T, token, path string, body any) *http.Response {
	t.Helper()
	return doJSON(t, http.MethodPost, token, path, body)
}

func doJSON(t *testing.T, method, token, path string, body any) *http.Response {
	t.Helper()

	var reader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal request body: %v", err)
		}
		reader = bytes.NewReader(data)
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, method, gatewayURL()+path, reader)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("%s %s: %v", method, path, err)
	}
	return resp
}

func assert2xx(t *testing.T, resp *http.Response) {
	t.Helper()
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("status = %s, want 2xx: %s", resp.Status, string(body))
	}
}

func requireEnv(t *testing.T, key string) string {
	t.Helper()
	value := os.Getenv(key)
	if value == "" {
		t.Skipf("%s is required for this integration fixture", key)
	}
	return value
}

func envDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func findString(value any, key string) string {
	switch typed := value.(type) {
	case map[string]any:
		for k, v := range typed {
			if k == key {
				if str, ok := v.(string); ok {
					return str
				}
			}
			if found := findString(v, key); found != "" {
				return found
			}
		}
	case []any:
		for _, item := range typed {
			if found := findString(item, key); found != "" {
				return found
			}
		}
	}
	return ""
}

func waitForKafkaEvent(t *testing.T, topic, eventType string) []byte {
	t.Helper()
	brokers := strings.Split(envDefault("KILAT_KAFKA_BROKERS", "localhost:9092"), ",")
	groupID := fmt.Sprintf("kilat-integration-%d", time.Now().UnixNano())
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		GroupID: groupID,
		Topic:   topic,
	})
	defer reader.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			t.Fatalf("waiting for Kafka event %s on %s: %v", eventType, topic, err)
		}
		var event struct {
			Type string `json:"type"`
		}
		if err := json.Unmarshal(msg.Value, &event); err == nil && event.Type == eventType {
			return msg.Value
		}
	}
}
