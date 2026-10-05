import readline from "node:readline/promises";
import { Mista } from "mista-sdk";

const mista = new Mista(); // reads MISTA_API_TOKEN

// 1. Send a code (WhatsApp first, SMS fallback)
const { sid } = await mista.verify.start({ to: "+15555550100", channel: "whatsapp_sms" });

// 2. Ask the user for it
const rl = readline.createInterface({ input: process.stdin, output: process.stdout });
const code = await rl.question("Code: ");
rl.close();

// 3. Check it
const result = await mista.verify.check({ sid, code });
console.log(result.verified ? "Verified" : `Not verified: ${result.reason}`);
