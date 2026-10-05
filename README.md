# Mista Verify Examples

Examples for integrating Mista Verify into your application.

Mista Verify provides a single API for verification with intelligent delivery, verification attempts, and channel fallback.

## Quick Start

See the documentation: https://docs.mista.io/

Set your API token (Mista dashboard → Settings → API), then replace `+15555550100` in the example with your phone number.

```bash
export MISTA_API_TOKEN="your_api_token_here"
```

| Language | Run |
| --- | --- |
| [Node.js](node/verify.js) | `cd node && npm install && node verify.js` |
| [Python](python/verify.py) | `cd python && pip install -r requirements.txt && python verify.py` |
| [Go](go/main.go) | `cd go && go run .` |

Each example does the same three steps: send a code, ask for it, check it.

## Channels

| `channel` | Delivery |
| --- | --- |
| `auto` | Uses your dashboard settings |
| `sms` | SMS only |
| `whatsapp_sms` | WhatsApp, falls back to SMS |
| `sms_whatsapp` | SMS, falls back to WhatsApp |
| `whatsapp_only` | WhatsApp only |

## What you can build

- OTP verification
- SMS verification
- WhatsApp fallback
- Authentication flows
- Onboarding verification
