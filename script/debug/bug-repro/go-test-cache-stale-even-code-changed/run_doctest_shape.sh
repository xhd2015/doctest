#!/usr/bin/env bash
# Doctest-shaped repro of stale go testcache after intermediate Setup edit.
#
# KEY FACTOR — not cwd / not testlog of gen files:
#
#   go test result cache uses two IDs:
#     id1 = link *action* ID  (hashes mid.a content ID as a packagefile)
#     id2 = linked *binary content* ID
#
#   Editing mid.Setup that only assigns an *unread* field:
#     • mid.a (and often leaf.a) rebuild
#     • optimizer may erase the write from the final binary (DCE/inline)
#     • binary content ID stays the SAME
#     • id1 MISSES, id2 may HIT → (cached) even though mid.go changed
#
#   Control (observe/): leaf reads WorkDir → content ID moves → re-run FAIL.
#
# Usage:
#   ./run_doctest_shape.sh           # dce + observe
#   ./run_doctest_shape.sh dce
#   ./run_doctest_shape.sh observe
set -euo pipefail

ROOT="$(cd "$(dirname "$0")" && pwd)"
MODE="${1:-both}"
VERDICTS=$(mktemp)
: >"$VERDICTS"

write_mid() {
  local dir=$1 tag=$2
  cat >"$dir/mid/mid.go" <<EOF
package mid

import (
	"testing"

	"dt.local/droot"
)

func Setup(t *testing.T, req *droot.Request) error {
	req.WorkDir = "${tag}"
	return nil
}
EOF
}

# Extract content-id *input* (second test ID's left-hand hex from gocachetest).
# Format: testcache: pkg: test ID <id> => <testID>
# First pair is usually action (id1), second is content (id2).
content_id_input() {
  local f=$1
  # second "test ID X => Y" line's X
  grep 'testcache:.*test ID' "$f" | sed -n '2s/.*test ID \([0-9a-f]*\) =>.*/\1/p' | head -1
}

run_one() {
  local name=$1
  local dir="$ROOT/doctest_shape/$name"
  local gocache
  gocache=$(mktemp -d)
  export GOCACHE="$gocache"
  cd "$dir"

  echo "========================================"
  echo "mode=$name  GOCACHE=$gocache"
  echo "========================================"
  write_mid "$dir" MID_V1

  echo "--- warm1 (GODEBUG=gocachetest=1) ---"
  GODEBUG=gocachetest=1 go test -buildvcs=false ./suite >"/tmp/ds-$name-w1.txt" 2>&1
  cat "/tmp/ds-$name-w1.txt"
  echo "--- warm2 ---"
  go test -buildvcs=false ./suite | tee "/tmp/ds-$name-w2.txt"

  local warm_cid
  warm_cid=$(content_id_input "/tmp/ds-$name-w1.txt")

  write_mid "$dir" MID_V2
  echo "--- after mid MID_V1→MID_V2 ---"
  set +e
  GODEBUG=gocachetest=1 go test -buildvcs=false ./suite >"/tmp/ds-$name-a1.txt" 2>&1
  local rc1=$?
  set -e
  cat "/tmp/ds-$name-a1.txt"

  echo "--- after2 (same mid V2, no edit) ---"
  go test -buildvcs=false ./suite | tee "/tmp/ds-$name-a2.txt"

  local after_cid
  after_cid=$(content_id_input "/tmp/ds-$name-a1.txt")

  local cached1=0 cached2=0 fail1=0
  grep -q '(cached)' "/tmp/ds-$name-a1.txt" && cached1=1
  grep -q '(cached)' "/tmp/ds-$name-a2.txt" && cached2=1
  grep -q 'FAIL' "/tmp/ds-$name-a1.txt" && fail1=1

  echo "--- analysis ---"
  echo "warm content-id input:  ${warm_cid:-?}"
  echo "after content-id input: ${after_cid:-?}"
  local cid_same=0
  if [[ -n "${warm_cid:-}" && -n "${after_cid:-}" && "$warm_cid" == "$after_cid" ]]; then
    cid_same=1
    echo "binary content ID: STABLE across mid.go edit"
  else
    echo "binary content ID: changed (or unavailable)"
  fi

  local verdict
  if [[ "$name" == "dce" ]]; then
    if [[ "$cached1" -eq 1 ]]; then
      verdict=BUG_REPRODUCED
      echo
      echo ">>> STALE (cached) on first run after mid.go edit (unread WorkDir / DCE path)."
    elif [[ "$cid_same" -eq 1 && "$cached2" -eq 1 && "$fail1" -eq 0 ]]; then
      # id2 same → Go may re-exec once (missing result blob) then cache;
      # content ID ignoring mid source is still the product-relevant bug.
      verdict=BUG_CONTENT_ID_STABLE
      echo
      echo ">>> Content ID stable after mid.go edit; after2 is (cached)."
      echo ">>> id1 missed (mid.a still link input); linked code effectively unchanged (DCE)."
    else
      verdict=NO_STALE_CACHE
    fi
    echo ">>> KEY FACTOR: not cwd — testlog never lists mid.go; binary content ID does."
  else
    if [[ "$fail1" -eq 1 && "$cached1" -eq 0 ]]; then
      verdict=GO_OK_RERUN_FAIL
      echo
      echo ">>> CONTROL: observing WorkDir forces content ID change → re-run FAIL."
    else
      verdict=UNEXPECTED
    fi
  fi

  echo
  echo "VERDICT[$name]=$verdict  a1_cached=$cached1 a2_cached=$cached2 fail=$fail1 cid_same=$cid_same"
  echo "$name=$verdict" >>"$VERDICTS"
  write_mid "$dir" MID_V1
  echo
}

if [[ "$MODE" == "both" || "$MODE" == "dce" ]]; then
  run_one dce
fi
if [[ "$MODE" == "both" || "$MODE" == "observe" ]]; then
  run_one observe
fi

echo "======== summary ========"
cat "$VERDICTS"
if grep -qE 'BUG_' "$VERDICTS"; then
  echo
  echo "Doctest-like stale-cache factor demonstrated (see VERDICT lines)."
  rm -f "$VERDICTS"
  exit 1
fi
rm -f "$VERDICTS"
exit 0
