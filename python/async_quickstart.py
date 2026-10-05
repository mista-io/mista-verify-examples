"""Async OTP quickstart with AsyncMista (for FastAPI, aiohttp, etc.).

    export MISTA_API_TOKEN="your_mista_api_token_here"
    python async_quickstart.py +15555550100
"""

import asyncio
import sys

from mista import AsyncMista, MistaError


async def main(phone: str) -> int:
    try:
        async with AsyncMista() as mista:
            verification = await mista.verify.start(to=phone, channel="auto")
            print(f"Code sent via {verification['channel']}. sid: {verification['sid']}")

            code = (await asyncio.to_thread(input, "Enter the code you received: ")).strip()
            result = await mista.verify.check(sid=verification["sid"], code=code)
    except MistaError as error:
        print(f"Mista error: {error}", file=sys.stderr)
        return 1

    if result["verified"]:
        print("Verified.")
        return 0
    print(f"Not verified ({result.get('reason')}).")
    return 1


if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: python async_quickstart.py <phone in E.164, e.g. +15555550100>", file=sys.stderr)
        sys.exit(1)
    sys.exit(asyncio.run(main(sys.argv[1])))
