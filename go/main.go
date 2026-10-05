package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mista-io/mista-go"
)

func main() {
	ctx := context.Background()
	client := mista.NewClient("") // reads MISTA_API_TOKEN

	// 1. Send a code (WhatsApp first, SMS fallback)
	v, err := client.Verify.Start(ctx, &mista.StartVerificationParams{To: "+15555550100", Channel: "whatsapp_sms"})
	if err != nil {
		log.Fatal(err)
	}

	// 2. Ask the user for it
	var code string
	fmt.Print("Code: ")
	fmt.Scanln(&code)

	// 3. Check it
	result, err := client.Verify.Check(ctx, v.SID, code)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Verified:", result.Verified, result.Reason)
}
