"""Choosing a delivery channel, including WhatsApp with SMS fallback.

    export MISTA_API_TOKEN="your_mista_api_token_here"
    python channels.py +15555550100 --channel whatsapp_sms

Channels:
    auto           follows your dashboard Verify settings (default)
    sms            SMS only
    whatsapp_sms   WhatsApp first, falls back to SMS
    sms_whatsapp   SMS first, falls back to WhatsApp
    whatsapp_only  WhatsApp only, no fallback

WhatsApp routing is available on the Growth and Pro plans.
"""

import argparse
import sys

from mista import Mista, MistaError

CHANNELS = ["auto", "sms", "whatsapp_sms", "sms_whatsapp", "whatsapp_only"]


def main() -> int:
    parser = argparse.ArgumentParser(description="Mista Verify channels")
    parser.add_argument("phone", help="Phone number in E.164 format, e.g. +15555550100")
    parser.add_argument("--channel", default="whatsapp_sms", choices=CHANNELS)
    parser.add_argument("--sender-id", help="Optional approved sender ID for the SMS leg")
    args = parser.parse_args()

    try:
        mista = Mista()
        verification = mista.verify.start(
            to=args.phone,
            channel=args.channel,
            sender_id=args.sender_id,
        )
        print(f"Requested channel: {args.channel}")
        print(f"Mista is delivering via: {verification['channel']}")
        print(f"Status: {verification['status']}  sid: {verification['sid']}")

        # Look the verification up again later, e.g. from a status page or a retry job.
        current = mista.verify.get(verification["sid"])
        print(f"Current status: {current['status']}, expires at {current['expires_at']}")
    except MistaError as error:
        print(f"Mista error: {error}", file=sys.stderr)
        return 1
    return 0


if __name__ == "__main__":
    sys.exit(main())
