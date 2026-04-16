#!/usr/bin/env bash
# MCPByDojoGenesis/scripts/smoke-test.sh
#
# Smoke tests for the dojo MCP server.
# Validates: build, tool registration, skill bundle, offline handlers.
# Uses JSON-RPC over stdin with background process + sleep for response collection.
#
# Usage:
#   ./scripts/smoke-test.sh
#   DOJO_SKILLS_PATH=/path/to/plugins ./scripts/smoke-test.sh

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

echo ""
echo "  MCP server smoke tests"
echo "  repo: $REPO_ROOT"
echo ""

# ─── Build ───────────────────────────────────────────────────────────────────

echo "  [build] compiling dojo-mcp-server..."
cd "$REPO_ROOT"
go build -o /tmp/dojo-mcp-smoke ./cmd/server/
echo "  [build] OK"
echo ""

PASS=0
FAIL=0
declare -a FAIL_DETAILS=()

pass() { PASS=$((PASS + 1)); printf "  %-44s  %s\n" "$1" "PASS"; }
fail() { FAIL=$((FAIL + 1)); FAIL_DETAILS+=("$1: $2"); printf "  %-44s  %s\n" "$1" "FAIL — $2"; }

# ─── Helper: send JSON-RPC, collect response ─────────────────────────────────

# Pipes one or more JSON-RPC lines to the server, waits briefly, kills, returns output.
mcp_send() {
  local input="$1"
  local tmpout
  tmpout="$(mktemp /tmp/mcp-out-XXXXXX)"

  echo "$input" | /tmp/dojo-mcp-smoke > "$tmpout" 2>/dev/null &
  local pid=$!
  sleep 1
  kill "$pid" 2>/dev/null
  wait "$pid" 2>/dev/null || true
  cat "$tmpout"
  rm -f "$tmpout"
}

INIT='{"jsonrpc":"2.0","id":0,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"smoke","version":"1.0"}}}'

printf "  %-44s  %s\n" "Test" "Result"
printf "  %s\n" "$(printf '─%.0s' $(seq 1 60))"

# ─── Test 1: Initialize ─────────────────────────────────────────────────────

resp=$(mcp_send "$INIT")
if echo "$resp" | grep -q '"dojo-mcp-server"'; then
  pass "initialize: server name"
else
  fail "initialize: server name" "$(echo "$resp" | head -c 100)"
fi

if echo "$resp" | grep -q '"3.1.0"'; then
  pass "initialize: version 3.1.0"
else
  fail "initialize: version 3.1.0" "wrong version"
fi

# ─── Test 2: tools/list ─────────────────────────────────────────────────────

resp=$(mcp_send "$(printf '%s\n%s' "$INIT" '{"jsonrpc":"2.0","id":1,"method":"tools/list","params":{}}')")
tool_count=$(echo "$resp" | grep -o '"name":"dojo_' | wc -l | tr -d ' ')

if [[ "$tool_count" -ge 20 ]]; then
  pass "tools/list: $tool_count dojo_ tools registered"
else
  fail "tools/list: >=20 tools" "got $tool_count"
fi

# Check key tool names
for tool in dojo_scout dojo_memory_list dojo_seed_list dojo_agent_list dojo_project_status dojo_converge dojo_health dojo_disposition_list; do
  if echo "$resp" | grep -q "\"$tool\""; then
    pass "tool: $tool"
  else
    fail "tool: $tool" "not in tools/list"
  fi
done

# ─── Test 3: dojo_scout ─────────────────────────────────────────────────────

resp=$(mcp_send "$(printf '%s\n%s' "$INIT" '{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"dojo_scout","arguments":{"situation":"REST vs gRPC"}}}')")
if echo "$resp" | grep -qiE "tension|route|scout|framework|synthesis"; then
  pass "dojo_scout: returns scaffold"
else
  fail "dojo_scout: returns scaffold" "$(echo "$resp" | tail -1 | head -c 120)"
fi

# ─── Test 4: dojo_search_skills ──────────────────────────────────────────────

resp=$(mcp_send "$(printf '%s\n%s' "$INIT" '{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"dojo_search_skills","arguments":{"query":"debugging"}}}')")
if echo "$resp" | grep -qi "debug"; then
  pass "dojo_search_skills: finds debugging"
else
  fail "dojo_search_skills: finds debugging" "no results"
fi

# ─── Test 5: dojo_invoke_skill ───────────────────────────────────────────────

resp=$(mcp_send "$(printf '%s\n%s' "$INIT" '{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"dojo_invoke_skill","arguments":{"name":"strategic-scout"}}}')")
if echo "$resp" | grep -qi "strategic"; then
  pass "dojo_invoke_skill: loads strategic-scout"
else
  fail "dojo_invoke_skill: loads strategic-scout" "not found"
fi

# ─── Test 6: dojo_list_skills ────────────────────────────────────────────────

resp=$(mcp_send "$(printf '%s\n%s' "$INIT" '{"jsonrpc":"2.0","id":5,"method":"tools/call","params":{"name":"dojo_list_skills","arguments":{}}}')")
if echo "$resp" | grep -qi "dojo-craft\|strategic-thinking"; then
  pass "dojo_list_skills: shows plugins"
else
  fail "dojo_list_skills: shows plugins" "no plugin listing"
fi

# ─── Test 7: dojo_disposition_list ───────────────────────────────────────────

resp=$(mcp_send "$(printf '%s\n%s' "$INIT" '{"jsonrpc":"2.0","id":6,"method":"tools/call","params":{"name":"dojo_disposition_list","arguments":{}}}')")
if echo "$resp" | grep -qi "focused\|balanced"; then
  pass "dojo_disposition_list: shows presets"
else
  fail "dojo_disposition_list: shows presets" "missing presets"
fi

# ─── Test 8: dojo_health (offline) ───────────────────────────────────────────

resp=$(mcp_send "$(printf '%s\n%s' "$INIT" '{"jsonrpc":"2.0","id":7,"method":"tools/call","params":{"name":"dojo_health","arguments":{}}}')")
if echo "$resp" | grep -q "UNREACHABLE\|Gateway Health\|unreachable"; then
  pass "dojo_health: reports offline"
else
  fail "dojo_health: reports offline" "$(echo "$resp" | tail -1 | head -c 120)"
fi

# ─── Test 9: dojo_converge ───────────────────────────────────────────────────

resp=$(mcp_send "$(printf '%s\n%s' "$INIT" '{"jsonrpc":"2.0","id":8,"method":"tools/call","params":{"name":"dojo_converge","arguments":{}}}')")
if echo "$resp" | grep -q "RED\|YELLOW\|GREEN\|Convergence"; then
  pass "dojo_converge: returns signal"
else
  fail "dojo_converge: returns signal" "$(echo "$resp" | tail -1 | head -c 120)"
fi

# ─── Test 10: error handling ────────────────────────────────────────────────

resp=$(mcp_send "$(printf '%s\n%s' "$INIT" '{"jsonrpc":"2.0","id":9,"method":"tools/call","params":{"name":"dojo_scout","arguments":{"situation":""}}}')")
if echo "$resp" | grep -q "isError\|required\|cannot be empty"; then
  pass "error: empty situation rejected"
else
  fail "error: empty situation rejected" "$(echo "$resp" | tail -1 | head -c 120)"
fi

# ─── Test 11: dojo_log_decision ──────────────────────────────────────────────

# dojo_log_decision uses DOJO_ADR_PATH env var — test via mcp_send
ADR_DIR="$(mktemp -d /tmp/mcp-adr-XXXXXX)"
# Override the mcp_send to inject env var for this test
adr_input="$(printf '%s\n%s' "$INIT" '{"jsonrpc":"2.0","id":10,"method":"tools/call","params":{"name":"dojo_log_decision","arguments":{"title":"Smoke Test","context":"Testing","decision":"Use tests","consequences":"Coverage"}}}')"
adr_tmpout="$(mktemp /tmp/mcp-out-XXXXXX)"
echo "$adr_input" | DOJO_ADR_PATH="$ADR_DIR" /tmp/dojo-mcp-smoke > "$adr_tmpout" 2>/dev/null &
adr_pid=$!
sleep 1
kill "$adr_pid" 2>/dev/null || true
wait "$adr_pid" 2>/dev/null || true
adr_resp=$(cat "$adr_tmpout")
rm -f "$adr_tmpout"
if echo "$adr_resp" | grep -q "decision\|ADR\|logged\|smoke"; then
  pass "dojo_log_decision: accepted"
else
  fail "dojo_log_decision: accepted" "$(echo "$adr_resp" | tail -1 | head -c 120)"
fi
rm -rf "$ADR_DIR"

# ─── Test 12: dojo_memory_list (offline error) ──────────────────────────────

resp=$(mcp_send "$(printf '%s\n%s' "$INIT" '{"jsonrpc":"2.0","id":11,"method":"tools/call","params":{"name":"dojo_memory_list","arguments":{}}}')")
if echo "$resp" | grep -q "Gateway\|gateway\|unavailable\|offline\|isError"; then
  pass "dojo_memory_list: offline error"
else
  fail "dojo_memory_list: offline error" "$(echo "$resp" | tail -1 | head -c 120)"
fi

# ─── Summary ─────────────────────────────────────────────────────────────────

TOTAL=$((PASS + FAIL))
printf "  %s\n" "$(printf '─%.0s' $(seq 1 60))"
echo ""
echo "  $PASS/$TOTAL passed — $FAIL failed"
echo ""

if [[ $FAIL -gt 0 ]]; then
  echo "  Failures:"
  for d in "${FAIL_DETAILS[@]}"; do
    echo "    - $d"
  done
  echo ""
fi

rm -f /tmp/dojo-mcp-smoke
exit $FAIL
