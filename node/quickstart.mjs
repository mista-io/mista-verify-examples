// OTP quickstart: send a code, read it from the terminal, check it.
//
//   export MISTA_API_TOKEN="your_mista_api_token_here"
//   node quickstart.mjs +15555550100 [channel]

import readline from "node:readline/promises";
import { stdin as input, stdout as output } from "node:process";
import { Mista, MistaError, RateLimitError } from "mista-sdk";

const [to, channel = "auto"] = process.argv.slice(2);
if (!to) {
  console.error("Usage: node quickstart.mjs <phone in E.164, e.g. +15555550100> [channel]");
  process.exit(1);
}

try {
  // Reads MISTA_API_TOKEN from the environment. Never hard-code your token.
  const mista = new Mista();

  const verification = await mista.verify.start({ to, channel });
  console.log(`Code sent to ${verification.to} via ${verification.channel}.`);
  console.log(`sid: ${verification.sid}  expires: ${verification.expires_at}`);

  const rl = readline.createInterface({ input, output });
  const code = (await rl.question("Enter the code you received: ")).trim();
  rl.close();

  const result = await mista.verify.check({ sid: verification.sid, code });
  if (result.verified) {
    console.log(`Verified at ${result.verified_at}.`);
  } else {
    // A wrong or expired code is not an exception — it comes back with a reason.
    console.log(`Not verified (${result.reason}).`);
    process.exitCode = 1;
  }
} catch (error) {
  if (error instanceof RateLimitError) {
    console.error(`Rate limited. Retry in ${error.retryAfter ?? "a few"} seconds.`);
  } else if (error instanceof MistaError) {
    console.error(`Mista error: ${error.message}`);
  } else {
    throw error;
  }
  process.exitCode = 1;
}
