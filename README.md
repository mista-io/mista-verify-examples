# Mista Verify Examples

Examples for integrating Mista Verify into your application.

Mista Verify provides a single API for verification with intelligent delivery, verification attempts, and channel fallback.

## Quick Start

See the documentation: https://docs.mista.io/

1. Create an API token in the Mista dashboard (**Settings → API**).
2. Export it — every example reads the token from `MISTA_API_TOKEN`:

   ```bash
   export MISTA_API_TOKEN="your_mista_api_token_here"
   ```

3. Pick a language and run the quickstart with your own phone number (E.164 format):

   | Language | Install | Run |
   | --- | --- | --- |
   | Node.js 18+ | `cd node && npm install` | `node quickstart.mjs +15555550100` |
   | Python 3.9+ | `cd python && pip install -r requirements.txt` | `python quickstart.py +15555550100` |
   | Go 1.22+ | `cd go` | `go run ./cmd/quickstart +15555550100` |
   | cURL | — | `./curl/verify.sh +15555550100` |

`+15555550100` is a placeholder — replace it with a real number you can receive messages on.

## What you can build

- OTP verification
- SMS verification
- WhatsApp fallback
- Authentication flows
- Onboarding verification

## Examples

| Example | Node.js | Python | Go |
| --- | --- | --- | --- |
| OTP quickstart (send a code, check it) | [`quickstart.mjs`](node/quickstart.mjs) | [`quickstart.py`](python/quickstart.py) | [`cmd/quickstart`](go/cmd/quickstart/main.go) |
| Channels and WhatsApp → SMS fallback | [`channels.mjs`](node/channels.mjs) | [`channels.py`](python/channels.py) | [`cmd/channels`](go/cmd/channels/main.go) |
| Sign-up / login API server | [`server.mjs`](node/server.mjs) (Express) | [`server.py`](python/server.py) (Flask) | [`cmd/server`](go/cmd/server/main.go) (net/http) |
| Async client | — | [`async_quickstart.py`](python/async_quickstart.py) | — |
| Raw HTTP | [`curl/verify.sh`](curl/verify.sh) | | |

## How it works

```text
Your app                         Mista Verify                    User
   │  verify.start(to, channel)      │                             │
   ├────────────────────────────────►│  code via SMS / WhatsApp    │
   │◄──────────── { sid } ───────────┤────────────────────────────►│
   │                                 │                             │
   │  verify.check(sid, code)        │         user types code     │
   ├────────────────────────────────►│◄────────────────────────────┤
   │◄──── { verified: true } ────────┤                             │
```

1. **Start** — `POST /api/v3/verify` sends a one-time code and returns a `sid`. Keep the `sid` on your server (session, database), not just in the browser.
2. **Check** — `POST /api/v3/verify/check` with the `sid` and the code the user typed. A wrong code is **not** an error: you get `verified: false` with a `reason`.
3. **Get** (optional) — `GET /api/v3/verify/{sid}` returns the current status.

### Channels

| `channel` | Delivery |
| --- | --- |
| `auto` (default) | Follows your dashboard Verify settings |
| `sms` | SMS only |
| `whatsapp_sms` | WhatsApp first, falls back to SMS |
| `sms_whatsapp` | SMS first, falls back to WhatsApp |
| `whatsapp_only` | WhatsApp only, no fallback |

WhatsApp routing is available on the Growth and Pro plans.

### Check results

| `reason` | Meaning | What to tell the user |
| --- | --- | --- |
| `invalid_code` | The code didn't match | "That code is incorrect. Try again." |
| `expired` | The code timed out | "This code has expired. Request a new one." |
| `max_attempts` | Too many wrong codes | "Too many attempts. Request a new code." |
| `not_open` | Verification was already closed | "Request a new code." |

## Security notes

- **Never commit your API token.** Use environment variables or a secrets manager. `.env` is git-ignored; `.env.example` only holds placeholders.
- Call Mista Verify **from your backend only** — never ship your token in a browser or mobile app.
- Bind the `sid` to the user/phone on your server, and trust the phone number only after `verified: true`.
- Rate-limit your own `start` endpoint per IP and per phone number to prevent abuse.

## SDKs

- Node.js — [`mista-sdk`](https://github.com/mista-io/mista-node) · `npm install mista-sdk`
- Python — [`mista`](https://github.com/mista-io/mista-python) · `pip install mista`
- Go — [`mista-go`](https://github.com/mista-io/mista-go) · `go get github.com/mista-io/mista-go`

## License

[MIT](LICENSE)
