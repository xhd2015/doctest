#!/usr/bin/env bash
# Minimal repro: dependency package source change after warm go test cache.
#
# Verdict:
#   GO_BUG  — after mid.Version 1→2, go test is (cached) and still PASS
#   GO_OK   — after edit, test re-runs and FAILs (correct)
#   OTHER   — unexpected output
#
# Usage:
#   ./run.sh              # both variants
#   ./run.sh direct       # suite_direct only
#   ./run.sh blank        # suite_blank only
set -euo pipefail

ROOT="$(cd "$(dirname "$0")" && pwd)"
cd "$ROOT"

MODE="${1:-both}"
export GOCACHE="${GOCACHE:-$(mktemp -d)}"
echo "GOCACHE=$GOCACHE"
echo "go: $(go version)"
echo

MID="$ROOT/mid/mid.go"
VERDICTS_FILE="$(mktemp)"

restore_mid() {
  cat >"$MID" <<'EOF'
package mid

// Version is the only behavioral knob. run.sh rewrites the return value
// between warm runs. Keep the body tiny so inlining is likely (mirrors
// small intermediate Setup packages in doctest gen).
func Version() int {
	return 1
}
EOF
}

set_mid_v2() {
  cat >"$MID" <<'EOF'
package mid

// Version is the only behavioral knob. run.sh rewrites the return value
// between warm runs. Keep the body tiny so inlining is likely (mirrors
// small intermediate Setup packages in doctest gen).
func Version() int {
	return 2
}
EOF
}

run_variant() {
  local name="$1"
  local pkg="$2"
  local logdir
  logdir="$(mktemp -d)"

  echo "======== variant: $name ($pkg) ========"
  restore_mid

  echo "--- warm1 ---"
  go test -buildvcs=false "$pkg" | tee "$logdir/warm1.txt"
  echo "--- warm2 (expect cached) ---"
  go test -buildvcs=false "$pkg" | tee "$logdir/warm2.txt"

  set_mid_v2
  echo "--- after mid.Version 1→2 (GODEBUG=gocachetest=1) ---"
  set +e
  GODEBUG=gocachetest=1 go test -buildvcs=false "$pkg" >"$logdir/after.txt" 2>&1
  local rc=$?
  set -e
  cat "$logdir/after.txt"

  local cached=0 pass=0 fail=0
  if grep -q '(cached)' "$logdir/after.txt"; then
    cached=1
  fi
  if grep -qE '^ok[[:space:]]' "$logdir/after.txt"; then
    pass=1
  fi
  if grep -q 'FAIL' "$logdir/after.txt" || [[ "$rc" -ne 0 ]]; then
    fail=1
  fi

  if grep -q 'testcache:' "$logdir/after.txt"; then
    echo "--- testcache lines ---"
    grep 'testcache:' "$logdir/after.txt" || true
  fi

  local verdict
  if [[ "$cached" -eq 1 && "$pass" -eq 1 && "$fail" -eq 0 ]]; then
    verdict="GO_BUG"
  elif [[ "$fail" -eq 1 ]]; then
    verdict="GO_OK"
  else
    verdict="OTHER"
  fi

  echo
  echo "VERDICT[$name]=$verdict  (cached=$cached pass=$pass fail=$fail exit=$rc)"
  echo "logs: $logdir"
  echo
  echo "${name}=${verdict}" >>"$VERDICTS_FILE"

  restore_mid
}

if [[ "$MODE" == "both" || "$MODE" == "direct" ]]; then
  run_variant direct ./suite_direct
fi
if [[ "$MODE" == "both" || "$MODE" == "blank" ]]; then
  run_variant blank ./suite_blank
fi

echo "======== summary ========"
cat "$VERDICTS_FILE"
OVERALL=0
if grep -q '=GO_BUG$' "$VERDICTS_FILE"; then
  OVERALL=1
  echo "At least one variant looks like a Go testcache bug (stale PASS after dep source change)."
else
  echo "No GO_BUG verdict (Go re-ran after mid source change, or other outcome)."
fi
rm -f "$VERDICTS_FILE"
exit "$OVERALL"
