# Runner App Smoke Collection

Open this folder in Bruno and select the `local` environment. Start the local stack first:

```bash
make integration-up
make seed
```

Seeded runner credentials:

- Email: `runner.test@kilat.my`
- Password: `TestRunner123!`

Fill the environment values after login:

- `runnerToken`, `customerToken`, `agentToken`
- `threadId`, `messageId`, `bookingId`, `incidentId`
- `assigneeUserId`, `questId`, `referralId`, `bankAccountId`

## Manual Smoke Walkthrough

1. Sign up or log in as a customer and the seeded runner.
2. Create or reuse an accepted booking, then copy its `bookingId`.
3. Open chat, list threads, send a message, fetch messages, and mark the message read.
4. Arrive at pickup, submit proof, and complete the booking.
5. Confirm payment disbursement and loyalty progress in the relevant responses or Kafka logs.
6. Create a bank account, set it as default, and submit a cash-out from the app flow.
7. Create an SOS incident, assign it to an agent, transition it to in progress, then resolve it.
8. Check zones at the KLCC point and verify active zone metadata is returned.
9. Create a referral code, redeem an eligible referral, and confirm payout behavior.

When finished:

```bash
make integration-down
```
