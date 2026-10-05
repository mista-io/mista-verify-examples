// Sign-up / login verification API with Express.
//
//   export MISTA_API_TOKEN="your_mista_api_token_here"
//   node server.mjs
//
//   curl -X POST localhost:3000/api/verify/start -H 'Content-Type: application/json' \
//        -d '{"phone": "+15555550100", "channel": "whatsapp_sms"}'
//   curl -X POST localhost:3000/api/verify/check -H 'Content-Type: application/json' \
//        -d '{"sid": "SID_FROM_START", "code": "123456"}'

import express from "express";
import {
  AuthenticationError,
  BadRequestError,
  Mista,
  MistaError,
  RateLimitError,
  ValidationError,
} from "mista-sdk";

const mista = new Mista();
const app = express();
app.use(express.json());

// In production, keep this in your session store or database.
const pending = new Map(); // sid -> { phone }

const REASON_MESSAGES = {
  invalid_code: "That code is incorrect. Try again.",
  expired: "This code has expired. Request a new one.",
  max_attempts: "Too many attempts. Request a new code.",
  not_open: "This verification is no longer active. Request a new code.",
};

app.post("/api/verify/start", async (req, res, next) => {
  const { phone, channel = "auto" } = req.body ?? {};
  if (typeof phone !== "string" || !phone.startsWith("+")) {
    return res.status(400).json({ error: "phone must be in E.164 format, e.g. +15555550100" });
  }

  try {
    const verification = await mista.verify.start({ to: phone, channel });
    pending.set(verification.sid, { phone });
    res.json({
      sid: verification.sid,
      channel: verification.channel,
      expires_at: verification.expires_at,
    });
  } catch (error) {
    next(error);
  }
});

app.post("/api/verify/check", async (req, res, next) => {
  const { sid, code } = req.body ?? {};
  const entry = typeof sid === "string" ? pending.get(sid) : undefined;
  if (!entry || typeof code !== "string") {
    return res.status(400).json({ error: "Unknown verification. Request a new code." });
  }

  try {
    const result = await mista.verify.check({ sid, code: code.trim() });
    if (!result.verified) {
      return res.status(422).json({
        verified: false,
        reason: result.reason,
        error: REASON_MESSAGES[result.reason] ?? "Verification failed. Request a new code.",
      });
    }

    pending.delete(sid);
    // The phone number is now proven. Create the account / start the session here.
    res.json({ verified: true, phone: entry.phone });
  } catch (error) {
    next(error);
  }
});

app.use((error, _req, res, _next) => {
  if (error instanceof RateLimitError) {
    return res.status(429).json({ error: "Too many requests. Try again shortly.", retry_after: error.retryAfter });
  }
  if (error instanceof ValidationError || error instanceof BadRequestError) {
    return res.status(400).json({ error: error.message });
  }
  if (error instanceof AuthenticationError) {
    console.error("Mista rejected the API token. Check MISTA_API_TOKEN.");
    return res.status(500).json({ error: "Verification is temporarily unavailable." });
  }
  if (error instanceof MistaError) {
    console.error(error);
    return res.status(502).json({ error: "Verification is temporarily unavailable." });
  }
  console.error(error);
  res.status(500).json({ error: "Internal server error" });
});

const port = Number(process.env.PORT) || 3000;
app.listen(port, () => console.log(`Listening on http://localhost:${port}`));
