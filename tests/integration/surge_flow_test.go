package integration

import "testing"

func TestSurgeFlowPublishesZoneSurgeChanged(t *testing.T) {
	requireIntegration(t)
	requireEnv(t, "KILAT_SURGE_FIXTURE_READY")

	resp := getJSON(t, "", "/api/v1/zones/at?lat=3.1478&lng=101.6953")
	assert2xx(t, resp)

	waitForKafkaEvent(t, envDefault("KILAT_ZONE_EVENTS_TOPIC", "zone.events"), "ZoneSurgeChanged")
	waitForKafkaEvent(t, envDefault("KILAT_NOTIFICATION_EVENTS_TOPIC", "notification.events"), "NotificationSent")
}
