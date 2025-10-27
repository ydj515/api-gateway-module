#!/usr/bin/env bash
set -euo pipefail

# Repeatedly hits a gateway endpoint to surface the circuit breaker OPEN state.
# If the downstream service is healthy the breaker will not trip; stop it first.

GATEWAY_BASE="${1:-http://localhost:8080}"
ENDPOINT="${2:-/api/query?name=breaker&age=0}"
ATTEMPTS="${3:-10}"
SLEEP_BETWEEN="${4:-0.5}"

if ! [[ "${ATTEMPTS}" =~ ^[0-9]+$ ]] || [ "${ATTEMPTS}" -lt 1 ]; then
  echo "ATTEMPTS must be a positive integer (got: ${ATTEMPTS})" >&2
  exit 1
fi

FULL_URL="${GATEWAY_BASE%/}${ENDPOINT}"

BODY_FILE="$(mktemp)"
ERR_FILE="$(mktemp)"

cleanup() {
  rm -f "${BODY_FILE}" "${ERR_FILE}"
}
trap cleanup EXIT

echo "Circuit breaker test against ${FULL_URL}"
echo "Attempts: ${ATTEMPTS}, delay: ${SLEEP_BETWEEN}s"
echo

breaker_open=0
last_attempt=0

for attempt in $(seq 1 "${ATTEMPTS}"); do
  last_attempt="${attempt}"
  : >"${BODY_FILE}"
  : >"${ERR_FILE}"

  echo "---- Attempt ${attempt} ----"

  http_code="$(curl -sS -o "${BODY_FILE}" -w "%{http_code}" "${FULL_URL}" 2>"${ERR_FILE}" || echo "000")"

  if [ "${http_code}" = "000" ]; then
    echo "No HTTP response (curl error):"
    sed 's/^/  /' "${ERR_FILE}"
  else
    echo "HTTP ${http_code}"
    if [ -s "${BODY_FILE}" ]; then
      sed 's/^/  /' "${BODY_FILE}"
    fi
  fi

  if grep -qiE "circuit breaker is open|too many requests" "${BODY_FILE}" "${ERR_FILE}" 2>/dev/null; then
    breaker_open=1
    echo "** Circuit breaker OPEN detected on attempt ${attempt} **"
    break
  fi

  if [ "${attempt}" -lt "${ATTEMPTS}" ]; then
    sleep "${SLEEP_BETWEEN}"
    echo
  fi
done

echo
if [ "${breaker_open}" -eq 1 ]; then
  echo "Result: breaker entered OPEN state after ${last_attempt} attempts."
else
  echo "Result: breaker OPEN state not observed. Ensure the downstream target is failing (e.g., stop the backend) and retry."
fi
