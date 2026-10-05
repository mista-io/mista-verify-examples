"""OTP quickstart: send a code, read it from the terminal, check it.

    export MISTA_API_TOKEN="your_mista_api_token_here"
    python quickstart.py +15555550100 [--channel whatsapp_sms]
"""

import argparse
import sys

from mista import Mista, MistaError, RateLimitError

CHANNELS = ["auto", "sms", "whatsapp_sms", "sms_whatsapp", "whatsapp_only"]


def main() -> int:
    parser = argparse.ArgumentParser(description="Mista Verify quickstart")
    parser.add_argument("phone", help="Phone number in E.164 format, e.g. +15555550100")
    parser.add_argument("--channel", default="auto", choices=CHANNELS)
    args = parser.parse_args()

    try:
        # Reads MISTA_API_TOKEN from the environment. Never hard-code your token.
        mista = Mista()

        verification = mista.verify.start(to=args.phone, channel=args.channel)
        print(f"Code sent to {verification['to']} via {verification['channel']}.")
        print(f"sid: {verification['sid']}  expires: {verification['expires_at']}")

        code = input("Enter the code you received: ").strip()
        result = mista.verify.check(sid=verification["sid"], code=code)
    except RateLimitError as error:
        print(f"Rate limited. Retry in {error.retry_after or 'a few'} seconds.", file=sys.stderr)
        return 1
    except MistaError as error:
        print(f"Mista error: {error}", file=sys.stderr)
        return 1

    if result["verified"]:
        print(f"Verified at {result.get('verified_at')}.")
        return 0

    # A wrong or expired code is not an exception — it comes back with a reason.
    print(f"Not verified ({result.get('reason')}).")
    return 1


if __name__ == "__main__":
    sys.exit(main())
