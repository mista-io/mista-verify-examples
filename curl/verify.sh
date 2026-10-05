#!/usr/bin/env bash
# Mista Verify over plain HTTP.
#
# Usage:
#   export MISTA_API_TOKEN="your_mista_api_token_here"
#   ./curl/verify.sh +15555550100 [channel]
#
# channel: auto (default), sms, whatsapp_sms, sms_whatsapp, whatsapp_only
set -euo pipefail

API="${MISTA_API_BASE_URL:-https://api.mista.io/api/v3}"
TO="${1:?Usage: $0 <phone in E.164, e.g. +15555550100> [channel]}"
CHANNEL="${2:-auto}"

if [[ -z "${MISTA_API_TOKEN:-}" ]]; then
  echo "Set MISTA_API_TOKEN first (Mista dashboard -> Settings -> API)." >&2
  exit 1
fi

echo "1) Starting verification for $TO via $CHANNEL..."
START=$(curl -sS -X POST "$API/verify" \
  -H "Authorization: Bearer $MISTA_API_TOKEN" \
  -H "Accept: application/json" \
  -H "Content-Type: application/json" \
  -d "{\"to\": \"$TO\", \"channel\": \"$CHANNEL\"}")
echo "$START"

SID=$(printf '%s' "$START" | sed -n 's/.*"sid"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p')
if [[ -z "$SID" ]]; then
  echo "No sid in the response — check the error above." >&2
  exit 1
fi

read -r -p "2) Enter the code you received: " CODE

echo "3) Checking code..."
curl -sS -X POST "$API/verify/check" \
  -H "Authorization: Bearer $MISTA_API_TOKEN" \
  -H "Accept: application/json" \
  -H "Content-Type: application/json" \
  -d "{\"sid\": \"$SID\", \"code\": \"$CODE\"}"
echo

echo "4) Current status:"
curl -sS "$API/verify/$SID" \
  -H "Authorization: Bearer $MISTA_API_TOKEN" \
  -H "Accept: application/json"
echo
