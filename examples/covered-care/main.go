// SPDX-License-Identifier: Apache-2.0

// Command covered-care demonstrates the typed eligibility and Booking flow with
// caller-owned idempotency keys. It performs requests only when explicitly run.
package main

import (
	"context"
	"log"
	"os"

	heyrafiki "github.com/heyrafiki/rafiki-go"
)

func main() {
	client, err := heyrafiki.NewClient(os.Getenv("HEYRAFIKI_API_KEY"))
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	eligibility, err := client.EligibilityChecks.Create(ctx, heyrafiki.EligibilityCheckInput{
		MemberReference: "member_demo_jubilee_001",
		ServiceCode:     "psychotherapy-60",
		ScheduledAt:     "2026-08-12T07:00:00Z",
		Amount:          350000,
		Currency:        heyrafiki.CurrencyKES,
	}, heyrafiki.WriteOptions{IdempotencyKey: "eligibility-demo-20260812-001"})
	if err != nil {
		log.Fatal(err)
	}

	booking, err := client.Bookings.Create(ctx, heyrafiki.BookingInput{
		PractitionerID: "prc_2481",
		StartsAt:       "2026-08-12T07:00:00Z",
		EndsAt:         "2026-08-12T08:00:00Z",
		Format:         heyrafiki.CareFormatOnline,
		PaymentSource:  heyrafiki.PaymentSourceCovered,
	}, heyrafiki.WriteOptions{IdempotencyKey: "booking-demo-20260812-001"})
	if err != nil {
		log.Fatal(err)
	}

	_, err = client.Preauthorizations.Create(ctx, heyrafiki.PreauthorizationInput{
		EligibilityCheckID: eligibility.ID,
		BookingID:          booking.ID,
	}, heyrafiki.WriteOptions{IdempotencyKey: "preauthorization-demo-20260812-001"})
	if err != nil {
		log.Fatal(err)
	}
}
