// Choosing a delivery channel, including WhatsApp with SMS fallback.
//
//   export MISTA_API_TOKEN="your_mista_api_token_here"
//   node channels.mjs +15555550100 whatsapp_sms
//
// Channels:
//   auto          follows your dashboard Verify settings (default)
//   sms           SMS only
//   whatsapp_sms  WhatsApp first, falls back to SMS
//   sms_whatsapp  SMS first, falls back to WhatsApp
//   whatsapp_only WhatsApp only, no fallback
//
// WhatsApp routing is available on the Growth and Pro plans.

import { Mista, MistaError } from "mista-sdk";

const CHANNELS = ["auto", "sms", "whatsapp_sms", "sms_whatsapp", "whatsapp_only"];

const [to, channel = "whatsapp_sms", senderId] = process.argv.slice(2);
if (!to || !CHANNELS.includes(channel)) {
  console.error(`Usage: node channels.mjs <phone> [${CHANNELS.join("|")}] [senderId]`);
  process.exit(1);
}

try {
  const mista = new Mista(); // reads MISTA_API_TOKEN

  const verification = await mista.verify.start({
    to,
    channel,
    // Optional: an approved sender ID for the SMS leg.
    ...(senderId ? { senderId } : {}),
  });

  console.log(`Requested channel: ${channel}`);
  console.log(`Mista is delivering via: ${verification.channel}`);
  console.log(`Status: ${verification.status}  sid: ${verification.sid}`);

  // Look the verification up again later, e.g. from a status page or a retry job.
  const current = await mista.verify.get(verification.sid);
  console.log(`Current status: ${current.status}, expires at ${current.expires_at}`);
} catch (error) {
  if (!(error instanceof MistaError)) throw error;
  console.error(`Mista error: ${error.message}`);
  process.exitCode = 1;
}
