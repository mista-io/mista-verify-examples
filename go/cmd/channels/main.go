// Choosing a delivery channel, including WhatsApp with SMS fallback.
//
//	export MISTA_API_TOKEN="your_mista_api_token_here"
//	go run ./cmd/channels +15555550100 whatsapp_sms [senderID]
//
// Channels:
//
//	auto           follows your dashboard Verify settings (default)
//	sms            SMS only
//	whatsapp_sms   WhatsApp first, falls back to SMS
//	sms_whatsapp   SMS first, falls back to WhatsApp
//	whatsapp_only  WhatsApp only, no fallback
//
// WhatsApp routing is available on the Growth and Pro plans.
package main

import (
	"context"
	"fmt"
	"os"
	"slices"
	"time"

	"github.com/mista-io/mista-go"
)

var channels = []string{"auto", "sms", "whatsapp_sms", "sms_whatsapp", "whatsapp_only"}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: go run ./cmd/channels <phone> [%v] [senderID]\n", channels)
		os.Exit(1)
	}
	params := &mista.StartVerificationParams{To: os.Args[1], Channel: "whatsapp_sms"}
	if len(os.Args) > 2 {
		params.Channel = os.Args[2]
	}
	if len(os.Args) > 3 {
		// Optional: an approved sender ID for the SMS leg.
		params.SenderID = os.Args[3]
	}
	if !slices.Contains(channels, params.Channel) {
		fmt.Fprintf(os.Stderr, "Unknown channel %q. Use one of %v.\n", params.Channel, channels)
		os.Exit(1)
	}

	client := mista.NewClient("")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	verification, err := client.Verify.Start(ctx, params)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("Requested channel: %s\n", params.Channel)
	fmt.Printf("Mista is delivering via: %s\n", verification.Channel)
	fmt.Printf("Status: %s  sid: %s\n", verification.Status, verification.SID)

	// Look the verification up again later, e.g. from a status page or a retry job.
	current, err := client.Verify.Get(ctx, verification.SID)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("Current status: %s, expires at %s\n", current.Status, current.ExpiresAt)
}
