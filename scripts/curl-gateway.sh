#!/usr/bin/env bash
set -euo pipefail

GATEWAY_BASE="${1:-http://localhost:8080}"

curl_json() {
  local label="$1"
  echo "==> ${label}"
  shift
  curl -sS -w "\nHTTP %{http_code}\n" "$@"
  echo
}

curl_json "GET /api/query" \
  "${GATEWAY_BASE}/api/query?name=Alice&age=30"

curl_json "GET /api/user/:id" \
  "${GATEWAY_BASE}/api/user/42"

curl_json "POST /api/create" \
  -H "Content-Type: application/json" \
  -d '{"title":"Example resource","status":"draft"}' \
  "${GATEWAY_BASE}/api/create"

curl_json "PUT /api/update/:id" \
  -X PUT \
  -H "Content-Type: application/json" \
  -d '{"status":"updated"}' \
  "${GATEWAY_BASE}/api/update/42"

curl_json "DELETE /api/delete/:id" \
  -X DELETE \
  "${GATEWAY_BASE}/api/delete/42"
