"""Sign-up / login verification API with Flask.

    export MISTA_API_TOKEN="your_mista_api_token_here"
    python server.py

    curl -X POST localhost:3000/api/verify/start -H 'Content-Type: application/json' \\
         -d '{"phone": "+15555550100", "channel": "whatsapp_sms"}'
    curl -X POST localhost:3000/api/verify/check -H 'Content-Type: application/json' \\
         -d '{"sid": "SID_FROM_START", "code": "123456"}'
"""

import logging
import os

from flask import Flask, jsonify, request
from mista import (
    AuthenticationError,
    BadRequestError,
    Mista,
    MistaError,
    RateLimitError,
    ValidationError,
)

app = Flask(__name__)
mista = Mista()  # reads MISTA_API_TOKEN

# In production, keep this in your session store or database.
pending: dict = {}  # sid -> {"phone": ...}

REASON_MESSAGES = {
    "invalid_code": "That code is incorrect. Try again.",
    "expired": "This code has expired. Request a new one.",
    "max_attempts": "Too many attempts. Request a new code.",
    "not_open": "This verification is no longer active. Request a new code.",
}


@app.post("/api/verify/start")
def start_verification():
    body = request.get_json(silent=True) or {}
    phone = body.get("phone")
    channel = body.get("channel", "auto")
    if not isinstance(phone, str) or not phone.startswith("+"):
        return jsonify(error="phone must be in E.164 format, e.g. +15555550100"), 400

    verification = mista.verify.start(to=phone, channel=channel)
    pending[verification["sid"]] = {"phone": phone}
    return jsonify(
        sid=verification["sid"],
        channel=verification["channel"],
        expires_at=verification["expires_at"],
    )


@app.post("/api/verify/check")
def check_verification():
    body = request.get_json(silent=True) or {}
    sid = body.get("sid")
    code = body.get("code")
    entry = pending.get(sid) if isinstance(sid, str) else None
    if entry is None or not isinstance(code, str):
        return jsonify(error="Unknown verification. Request a new code."), 400

    result = mista.verify.check(sid=sid, code=code.strip())
    if not result["verified"]:
        reason = result.get("reason")
        message = REASON_MESSAGES.get(reason, "Verification failed. Request a new code.")
        return jsonify(verified=False, reason=reason, error=message), 422

    pending.pop(sid, None)
    # The phone number is now proven. Create the account / start the session here.
    return jsonify(verified=True, phone=entry["phone"])


@app.errorhandler(MistaError)
def handle_mista_error(error: MistaError):
    if isinstance(error, RateLimitError):
        return jsonify(error="Too many requests. Try again shortly.", retry_after=error.retry_after), 429
    if isinstance(error, (ValidationError, BadRequestError)):
        return jsonify(error=str(error)), 400
    if isinstance(error, AuthenticationError):
        logging.error("Mista rejected the API token. Check MISTA_API_TOKEN.")
    else:
        logging.exception("Mista request failed")
    return jsonify(error="Verification is temporarily unavailable."), 502


if __name__ == "__main__":
    app.run(port=int(os.environ.get("PORT", "3000")))
