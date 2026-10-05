# Mista Verify — Node.js

Uses the official [`mista-sdk`](https://github.com/mista-io/mista-node) package (Node.js 18+).

```bash
npm install
export MISTA_API_TOKEN="your_mista_api_token_here"
```

| File | What it shows |
| --- | --- |
| [`quickstart.mjs`](quickstart.mjs) | Send a code, type it in, check it |
| [`channels.mjs`](channels.mjs) | Pick a channel — e.g. WhatsApp with SMS fallback — and look up status |
| [`server.mjs`](server.mjs) | Express API for sign-up / login verification |

```bash
node quickstart.mjs +15555550100
node channels.mjs +15555550100 whatsapp_sms
node server.mjs
```

`+15555550100` is a placeholder — use a real number you can receive messages on.

On Node.js 20.6+ you can load the token from the repo's `.env` instead of exporting it:

```bash
cp ../.env.example ../.env
node --env-file=../.env quickstart.mjs +15555550100
```

## The SDK in 10 lines

```js
import { Mista } from "mista-sdk";

const mista = new Mista(); // reads MISTA_API_TOKEN

const verification = await mista.verify.start({ to: "+15555550100", channel: "auto" });
const result = await mista.verify.check({ sid: verification.sid, code: "123456" });

if (result.verified) {
  // phone number confirmed
} else {
  console.log(result.reason); // invalid_code | expired | max_attempts | not_open
}
```
