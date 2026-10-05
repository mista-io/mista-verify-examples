// OTP quickstart: send a code, read it from the terminal, check it.
//
//	export MISTA_API_TOKEN="your_mista_api_token_here"
//	go run ./cmd/quickstart +15555550100 [channel]
package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mista-io/mista-go"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run ./cmd/quickstart <phone in E.164, e.g. +15555550100> [channel]")
		os.Exit(1)
	}
	to := os.Args[1]
	channel := "auto"
	if len(os.Args) > 2 {
		channel = os.Args[2]
	}

	// An empty token reads MISTA_API_TOKEN from the environment. Never hard-code your token.
	client := mista.NewClient("")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	verification, err := client.Verify.Start(ctx, &mista.StartVerificationParams{To: to, Channel: channel})
	if err != nil {
		exitWithError(err)
	}
	fmt.Printf("Code sent to %s via %s.\n", verification.To, verification.Channel)
	fmt.Printf("sid: %s  expires: %s\n", verification.SID, verification.ExpiresAt)

	fmt.Print("Enter the code you received: ")
	code, _ := bufio.NewReader(os.Stdin).ReadString('\n')

	result, err := client.Verify.Check(ctx, verification.SID, strings.TrimSpace(code))
	if err != nil {
		exitWithError(err)
	}
	if !result.Verified {
		// A wrong or expired code is not an error — it comes back with a reason.
		fmt.Printf("Not verified (%s).\n", result.Reason)
		os.Exit(1)
	}
	fmt.Printf("Verified at %s.\n", result.VerifiedAt)
}

func exitWithError(err error) {
	var apiErr *mista.APIError
	switch {
	case errors.Is(err, mista.ErrMissingToken):
		fmt.Fprintln(os.Stderr, "Set MISTA_API_TOKEN first (Mista dashboard -> Settings -> API).")
	case errors.Is(err, mista.ErrRateLimited) && errors.As(err, &apiErr):
		fmt.Fprintf(os.Stderr, "Rate limited. Retry in %d seconds.\n", apiErr.RetryAfter())
	default:
		fmt.Fprintln(os.Stderr, err)
	}
	os.Exit(1)
}
