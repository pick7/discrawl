#!/usr/bin/env bash
set -euo pipefail

TARGET=${1:-}
BINARY=${2:-}
REQUIRED=${DISCRAWL_CODESIGN_REQUIRED:-0}
IDENTIFIER=${DISCRAWL_CODESIGN_IDENTIFIER:-org.openclaw.discrawl}
EXPECTED_TEAM_ID=${DISCRAWL_CODESIGN_TEAM_ID:-FWJYW4S8P8}
REQUIREMENT="identifier \"$IDENTIFIER\" and anchor apple generic and certificate 1[field.1.2.840.113635.100.6.2.6] exists and certificate leaf[field.1.2.840.113635.100.6.1.13] exists and certificate leaf[subject.OU] = \"$EXPECTED_TEAM_ID\""

case "$TARGET" in
  darwin_*) ;;
  *) exit 0 ;;
esac

if [[ "$REQUIRED" != 1 && -z "${CODESIGN_IDENTITY:-}" ]]; then
  echo "skipping Developer ID signing for snapshot target $TARGET"
  exit 0
fi

[[ "$(uname -s)" == Darwin ]] || {
  echo "official macOS release signing must run on macOS" >&2
  exit 1
}
[[ -n "${CODESIGN_IDENTITY:-}" ]] || {
  echo "CODESIGN_IDENTITY is required; run the release through mac-release codesign-run" >&2
  exit 1
}
[[ -f "$BINARY" ]] || {
  echo "release binary not found: $BINARY" >&2
  exit 1
}

codesign --force --options runtime --timestamp \
  --identifier "$IDENTIFIER" --sign "$CODESIGN_IDENTITY" "$BINARY"
codesign --verify --strict -R="$REQUIREMENT" --verbose=2 "$BINARY"

signature=$(codesign -dvvv "$BINARY" 2>&1)
grep -Fx "Identifier=$IDENTIFIER" <<<"$signature" >/dev/null
grep -Fx "TeamIdentifier=$EXPECTED_TEAM_ID" <<<"$signature" >/dev/null
grep -F "Authority=Developer ID Application:" <<<"$signature" >/dev/null
