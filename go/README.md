# Mista Verify — Go

Uses the official [`mista-go`](https://github.com/mista-io/mista-go) module (Go 1.22+).

```bash
export MISTA_API_TOKEN="your_mista_api_token_here"
```

| Command | What it shows |
| --- | --- |
| [`cmd/quickstart`](cmd/quickstart/main.go) | Send a code, type it in, check it |
| [`cmd/channels`](cmd/channels/main.go) | Pick a channel — e.g. WhatsApp with SMS fallback — and look up status |
| [`cmd/server`](cmd/server/main.go) | net/http API for sign-up / login verification |

```bash
go run ./cmd/quickstart +15555550100
go run ./cmd/channels +15555550100 whatsapp_sms
go run ./cmd/server
```

`+15555550100` is a placeholder — use a real number you can receive messages on.

## The SDK in 10 lines

```go
client := mista.NewClient("") // reads MISTA_API_TOKEN

v, err := client.Verify.Start(ctx, &mista.StartVerificationParams{To: "+15555550100", Channel: "auto"})
if err != nil {
	return err
}
result, err := client.Verify.Check(ctx, v.SID, "123456")
if err != nil {
	return err
}
if !result.Verified {
	fmt.Println(result.Reason) // invalid_code | expired | max_attempts | not_open
}
```
