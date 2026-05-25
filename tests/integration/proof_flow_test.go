package integration

import "testing"

func TestProofFlowCompletesDeliveryAndUpdatesPaymentAndLoyalty(t *testing.T) {
	requireIntegration(t)
	runnerToken := requireRunnerToken(t)
	bookingID := requireEnv(t, "KILAT_ACTIVE_BOOKING_ID")

	assert2xx(t, postJSON(t, runnerToken, "/api/v1/bookings/"+bookingID+"/arrive-pickup", map[string]any{
		"qr_code": envDefault("KILAT_PICKUP_QR_CODE", "integration-qr"),
	}))
	assert2xx(t, postJSON(t, runnerToken, "/api/v1/bookings/"+bookingID+"/proof", map[string]any{
		"photo_url":      "https://example.test/proof.jpg",
		"signature_url":  "https://example.test/signature.png",
		"recipient_kind": "owner",
		"notes":          "integration proof",
	}))
	assert2xx(t, postJSON(t, runnerToken, "/api/v1/bookings/"+bookingID+"/complete", map[string]any{}))

	waitForKafkaEvent(t, envDefault("KILAT_PAYMENT_EVENTS_TOPIC", "payment.events"), "PaymentDisbursed")
	waitForKafkaEvent(t, envDefault("KILAT_LOYALTY_EVENTS_TOPIC", "loyalty.events"), "QuestProgressUpdated")
}
