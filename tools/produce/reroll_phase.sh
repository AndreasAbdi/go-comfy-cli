#!/usr/bin/env bash
# Phase 2: re-render clips whose inputs changed after they were first rendered,
# then reassemble and rebuild the review packet.
# Waits for phase 1 (POST DONE) so only one model ever touches the GPU.
set -u
cd "$(dirname "$0")/../.."
PY=tools/refs/.venv/Scripts/python.exe
# stale_clips.py prints box-drawing-free ASCII, but Python still defaults to
# cp1252 on this shell and dies writing to a pipe.
export PYTHONIOENCODING=utf-8

until grep -qE "POST DONE" out/post_pass.log 2>/dev/null; do sleep 60; done
echo "=== phase 1 complete, starting re-roll phase ==="

# Keep the originals: 480p clips are small and a re-roll can come back worse.
BK=production/limbus-prologue/cut/_v1_clips
mkdir -p "$BK"
for f in production/limbus-prologue/shots/*/clips/*_CLIP_a1.mp4; do
  cp -n "$f" "$BK/$(basename "$f")" 2>/dev/null || true
done
echo "backed up $(ls "$BK" | wc -l) clips to $BK"

# grep -oP is unavailable in this Git Bash build ("-P supports only unibyte and
# UTF-8 locales"), so pull the list with sed.
LIST=$($PY tools/produce/stale_clips.py | sed -n 's/.*--shots \([0-9,]*\).*/\1/p' | tail -1)
if [ -z "$LIST" ]; then echo "nothing stale; skipping re-roll"; else
  echo "re-rolling: $LIST"
  $PY tools/produce/clips.py --megapixels 0.4 --force --shots "$LIST"
fi

echo "=== reassembling ==="
$PY tools/produce/assemble.py --attempt 1
echo "=== review packet ==="
$PY tools/produce/review_packet.py --attempt 1
echo "=== REROLL PHASE DONE ==="
