#!/usr/bin/env bash
# MCPByDojoGenesis/scripts/smoke-test.sh
#
# Smoke tests for the dojo MCP server.
# Tests tool registration, skill loading, and offline-mode handlers by
# sending JSON-RPC messages via stdin and checking stdout responses.
#
# Usage:
#   ./scripts/smoke-test.sh
#   DOJO_SKILLS_PATH=/path/to/plugins ./scripts/smoke-test.sh

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# Build the server binary
echo ""
echo "  MCP server smoke tests"
echo "  repo: $REPO_ROOT"
echo ""

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

# ─── Helper: send JSON-RPC request, get response ────────────────────────────

# Sends initialize + the actual request + a shutdown via stdin.
# Returns the JSON response for the actual request (line 2 of output).
mcp_call() {
  local method="$1"
  local params="$2"
  local id="${3:-1}"

  local init_req='{"jsonrpc":"2.0","id":0,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"smoke-test","version":"1.0.0"}}}'
  local actual_req="{\"jsonrpc\":\"2.0\",\"id\":$id,\"method\":\"$method\",\"params\":$params}"

  # Send both requests, capture output, timeout after 10s
  echo -e "${init_req}\n${actual_req}" | \
    timeout 10 /tmp/dojo-mcp-smoke 2>/dev/null | \
    tail -1
}

# macOS timeout fallback
if ! command -v timeout &>/dev/null; then
  if command -v gtimeout &>/dev/null; then
    timeout() { gtimeout "$@"; }
  else
    timeout() {
      local secs=$1; shift
      "$@" &
      local pid=$!
      ( sleep "$secs" && kill "$pid" 2>/dev/null ) &
      local watcher=$!
      wait "$pid" 2>/dev/null; local rc=$?
      kill "$watcher" 2>/dev/null; wait "$watcher" 2>/dev/null
      return $rc
    }
  fi
fi

printf "  %-44s  %s\n" "Test" "Result"
printf "  %s\n" "$(printf '─%.0s' $(seq 1 60))"

# ─── Test 1: Initialize response ────────────────────────────────────────────

init_resp=$(echo '{"jsonrpc":"2.0","id":0,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"smoke","version":"1.0"}}}' | \
  timeout 10 /tmp/dojo-mcp-smoke 2>/dev/null | head -1)

if echo "$init_resp" | grep -q '"serverInfo"'; then
  pass "initialize returns serverInfo"
else
  fail "initialize returns serverInfo" "missing serverInfo in response"
fi

if echo "$init_resp" | grep -q '"dojo-mcp-server"'; then
  pass "server name is dojo-mcp-server"
else
  fail "server name is dojo-mcp-server" "wrong server name"
fi

if echo "$init_resp" | grep -q '"3.1.0"'; then
  pass "server version is 3.1.0"
else
  fail "server version is 3.1.0" "wrong version"
fi

# ─── Test 2: tools/list returns 24 tools ─────────────────────────────────────

tools_resp=$(mcp_call "tools/list" '{}')

tool_count=$(echo "$tools_resp" | grep -o '"name"' | wc -l | tr -d ' ')
if [[ "$tool_count" -ge 20 ]]; then
  pass "tools/list returns $tool_count tools"
else
  fail "tools/list returns >=20 tools" "got $tool_count"
fi

# Check specific tool names exist
for tool in dojo_scout dojo_memory_list dojo_seed_list dojo_agent_list dojo_project_status dojo_converge dojo_health dojo_disposition_list; do
  if echo "$tools_resp" | grep -q "\"$tool\""; then
    pass "tool registered: $tool"
  else
    fail "tool registered: $tool" "not found in tools/list"
  fi
done

# ─── Test 3: dojo_scout (offline mode) ───────────────────────────────────────

scout_resp=$(mcp_call "tools/call" '{"name":"dojo_scout","arguments":{"situation":"should we use REST or gRPC"}}')

if echo "$scout_resp" | grep -qE "Tension|tension|Routes|routes|Scout|Framework"; then
  pass "dojo_scout returns scaffold"
else
  fail "dojo_scout returns scaffold" "missing expected sections"
fi

# ─── Test 4: dojo_search_skills ──────────────────────────────────────────────

search_resp=$(mcp_call "tools/call" '{"name":"dojo_search_skills","arguments":{"query":"debugging"}}')

if echo "$search_resp" | grep -q "debugging\|debug"; then
  pass "dojo_search_skills finds debugging"
else
  fail "dojo_search_skills finds debugging" "no results for 'debugging'"
fi

# ─── Test 5: dojo_invoke_skill ───────────────────────────────────────────────

invoke_resp=$(mcp_call "tools/call" '{"name":"dojo_invoke_skill","arguments":{"name":"strategic-scout"}}')

if echo "$invoke_resp" | grep -q "strategic-scout\|Strategic"; then
  pass "dojo_invoke_skill loads strategic-scout"
else
  fail "dojo_invoke_skill loads strategic-scout" "skill not found"
fi

# ─── Test 6: dojo_list_skills ────────────────────────────────────────────────

list_resp=$(mcp_call "tools/call" '{"name":"dojo_list_skills","arguments":{}}')

if echo "$list_resp" | grep -q "dojo-craft\|strategic-thinking"; then
  pass "dojo_list_skills shows plugins"
else
  fail "dojo_list_skills shows plugins" "no plugin listing"
fi

# ─── Test 7: dojo_apply_seed ─────────────────────────────────────────────────

seed_resp=$(mcp_call "tools/call" '{"name":"dojo_apply_seed","arguments":{"seed_name":"three_tiered_governance","situation":"managing plugin quality"}}')

if echo "$seed_resp" | grep -q "governance\|tier\|Governance"; then
  pass "dojo_apply_seed applies seed"
else
  fail "dojo_apply_seed applies seed" "seed content not found"
fi

# ─── Test 8: dojo_log_decision ───────────────────────────────────────────────

ADR_DIR="$TMPDIR/smoke-decisions"
mkdir -p "$ADR_DIR"

log_resp=$(echo "{\"jsonrpc\":\"2.0\",\"id\":0,\"method\":\"initialize\",\"params\":{\"protocolVersion\":\"2024-11-05\",\"capabilities\":{},\"clientInfo\":{\"name\":\"smoke\",\"version\":\"1.0\"}}}
{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"tools/call\",\"params\":{\"name\":\"dojo_log_decision\",\"arguments\":{\"title\":\"Smoke Test Decision\",\"context\":\"Testing ADR writer\",\"decision\":\"Use smoke tests\",\"consequences\":\"Better coverage\"}}}" | \
  DOJO_ADR_PATH="$ADR_DIR" timeout 10 /tmp/dojo-mcp-smoke 2>/dev/null | tail -1)

adr_count=$(ls "$ADR_DIR"/*.md 2>/dev/null | wc -l | tr -d ' ')
if [[ "$adr_count" -ge 1 ]]; then
  pass "dojo_log_decision writes ADR file"
else
  fail "dojo_log_decision writes ADR file" "no .md file in $ADR_DIR"
fi

# ─── Test 9: dojo_reflect ────────────────────────────────────────────────────

reflect_resp=$(mcp_call "tools/call" '{"name":"dojo_reflect","arguments":{"session_description":"debugging a production outage"}}')

if echo "$reflect_resp" | grep -qE "Reflect|Skills|Seeds|debug"; then
  pass "dojo_reflect returns frameworks"
else
  fail "dojo_reflect returns frameworks" "no reflection content"
fi

# ─── Test 10: dojo_disposition_list ──────────────────────────────────────────

disp_resp=$(mcp_call "tools/call" '{"name":"dojo_disposition_list","arguments":{}}')

if echo "$disp_resp" | grep -q "focused\|balanced\|exploratory\|deliberate"; then
  pass "dojo_disposition_list shows presets"
else
  fail "dojo_disposition_list shows presets" "missing disposition presets"
fi

# ─── Test 11: dojo_project_status ────────────────────────────────────────────

proj_resp=$(mcp_call "tools/call" '{"name":"dojo_project_status","arguments":{}}')

if echo "$proj_resp" | grep -qE "Project|project|No project|No active"; then
  pass "dojo_project_status responds"
else
  fail "dojo_project_status responds" "no project response"
fi

# ─── Test 12: dojo_converge (offline) ────────────────────────────────────────

conv_resp=$(mcp_call "tools/call" '{"name":"dojo_converge","arguments":{}}')

if echo "$conv_resp" | grep -qE "RED|YELLOW|GREEN|Convergence|dirty"; then
  pass "dojo_converge returns signal"
else
  fail "dojo_converge returns signal" "no convergence signal"
fi

# ─── Test 13: dojo_health (offline — should report unreachable) ──────────────

health_resp=$(mcp_call "tools/call" '{"name":"dojo_health","arguments":{}}')

if echo "$health_resp" | grep -qE "unreachable|offline|Gateway|health"; then
  pass "dojo_health reports offline status"
else
  fail "dojo_health reports offline status" "unexpected health response"
fi

# ─── Test 14: error handling — empty required args ───────────────────────────

err_resp=$(mcp_call "tools/call" '{"name":"dojo_scout","arguments":{"situation":""}}')

if echo "$err_resp" | grep -qE "error\|required\|empty\|isError"; then
  pass "dojo_scout rejects empty situation"
else
  fail "dojo_scout rejects empty situation" "no error for empty arg"
fi

err_resp=$(mcp_call "tools/call" '{"name":"dojo_invoke_skill","arguments":{"name":""}}')

if echo "$err_resp" | grep -qE "error\|required\|empty\|isError"; then
  pass "dojo_invoke_skill rejects empty name"
else
  fail "dojo_invoke_skill rejects empty name" "no error for empty arg"
fi

# ─── Test 15: dojo_memory_list (offline — should return error) ───────────────

mem_resp=$(mcp_call "tools/call" '{"name":"dojo_memory_list","arguments":{}}')

if echo "$mem_resp" | grep -qE "error\|Gateway\|unreachable\|offline\|isError"; then
  pass "dojo_memory_list reports offline"
else
  fail "dojo_memory_list reports offline" "should report Gateway required"
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

# Cleanup
rm -f /tmp/dojo-mcp-smoke
rm -rf "${ADR_DIR:-/tmp/smoke-decisions-nonexistent}"

exit $FAIL
