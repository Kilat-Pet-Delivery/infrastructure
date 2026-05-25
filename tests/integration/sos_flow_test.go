package integration

import "testing"

func TestSOSFlowAssignsIncidentAndPushesNotification(t *testing.T) {
	requireIntegration(t)
	runnerToken := requireRunnerToken(t)
	_ = requireAgentToken(t)

	resp := postJSON(t, runnerToken, "/api/v1/incidents", map[string]any{
		"type":     "sos",
		"severity": "high",
		"notes":    "integration SOS",
	})
	assert2xx(t, resp)

	waitForKafkaEvent(t, envDefault("KILAT_INCIDENT_EVENTS_TOPIC", "incident.events"), "IncidentAssigned")
	waitForKafkaEvent(t, envDefault("KILAT_NOTIFICATION_EVENTS_TOPIC", "notification.events"), "NotificationSent")
}
