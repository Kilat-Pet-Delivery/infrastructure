package integration

import "testing"

func TestReferralFlowPaysOutAfterEligibleBookings(t *testing.T) {
	requireIntegration(t)
	referrerToken := requireRunnerToken(t)
	refereeToken := requireEnv(t, "KILAT_REFEREE_TOKEN")
	referralID := requireEnv(t, "KILAT_REFERRAL_ID")

	assert2xx(t, postJSON(t, referrerToken, "/api/v1/referrals/code", map[string]any{}))
	assert2xx(t, postJSON(t, refereeToken, "/api/v1/referrals/"+referralID+"/redeem", map[string]any{}))

	waitForKafkaEvent(t, envDefault("KILAT_LOYALTY_EVENTS_TOPIC", "loyalty.events"), "ReferralPayoutDue")
	waitForKafkaEvent(t, envDefault("KILAT_PAYMENT_EVENTS_TOPIC", "payment.events"), "CreditDisbursed")
}
