#!/usr/bin/env python
"""Derive the definitive WAV -> dialogue-line mapping by sequence alignment.

The recordings are an ordered subsequence of the script: every WAV corresponds to
exactly one line, lines appear in the same order as the WAVs, and unvoiced lines are
simply skipped. A fixed numeric offset cannot express that -- the observed offset
drifts from -1 to +4 across the scene -- so this does a global Needleman-Wunsch
alignment between the WAV sequence and the line sequence instead, scoring each pair
by the similarity of the WAV's English translation to the line text.

Reads the cache written by transcribe_vo.py; performs no transcription itself.

    tools/asr/.venv/Scripts/python.exe tools/asr/align_vo.py
"""
from __future__ import annotations

import argparse
import csv
import json
import os
import sys

sys.stdout.reconfigure(encoding="utf-8")
sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
from transcribe_vo import DEF_AUDIO, DEF_DIALOGUE, REPO, load_lines, similarity  # noqa: E402

SKIP_WAV = -2.00      # cost of leaving a recording unmatched (should never happen)


def align(wavs, lines, score, skip_line):
    n, m = len(wavs), len(lines)
    NEG = float("-inf")
    dp = [[NEG] * (m + 1) for _ in range(n + 1)]
    bt = [[None] * (m + 1) for _ in range(n + 1)]
    dp[0][0] = 0.0
    for j in range(1, m + 1):
        dp[0][j] = dp[0][j - 1] + skip_line
        bt[0][j] = "L"
    for i in range(1, n + 1):
        dp[i][0] = dp[i - 1][0] + SKIP_WAV
        bt[i][0] = "W"
    for i in range(1, n + 1):
        for j in range(1, m + 1):
            best, mv = dp[i - 1][j - 1] + score(i - 1, j - 1), "M"
            if dp[i][j - 1] + skip_line > best:
                best, mv = dp[i][j - 1] + skip_line, "L"
            if dp[i - 1][j] + SKIP_WAV > best:
                best, mv = dp[i - 1][j] + SKIP_WAV, "W"
            dp[i][j], bt[i][j] = best, mv
    i, j, out = n, m, []
    while i > 0 or j > 0:
        mv = bt[i][j]
        if mv == "M":
            out.append((i - 1, j - 1)); i -= 1; j -= 1
        elif mv == "L":
            j -= 1
        else:
            out.append((i - 1, None)); i -= 1
    return list(reversed(out)), dp[n][m]


def main() -> int:
    ap = argparse.ArgumentParser(description=__doc__,
                                 formatter_class=argparse.RawDescriptionHelpFormatter)
    ap.add_argument("--dialogue", default=DEF_DIALOGUE)
    ap.add_argument("--prefix", default="S001B")
    ap.add_argument("--asr-dir", default=os.path.join(REPO, "tools", "asr", "out"))
    ap.add_argument("--overrides", default=os.path.join(REPO, "tools", "asr", "overrides.json"),
                    help="verified manual corrections applied after alignment")
    ap.add_argument("--dur-weight", type=float, default=0.0,
                    help="weight of the duration prior; English length poorly predicts Korean audio")
    ap.add_argument("--skip-line", type=float, default=-0.10,
                    help="cost of leaving a script line unvoiced")
    ap.add_argument("--low", type=float, default=0.20,
                    help="flag pairs scoring below this for human review")
    args = ap.parse_args()

    cache_path = os.path.join(args.asr_dir, f"{args.prefix.lower()}-asr-cache.json")
    if not os.path.exists(cache_path):
        print("no ASR cache; run transcribe_vo.py first", file=sys.stderr)
        return 2
    with open(cache_path, encoding="utf-8") as fh:
        cache = json.load(fh)

    lines = load_lines(args.dialogue)
    ids = sorted(lines)
    wavs = sorted(cache, key=lambda k: int(k))

    def score(wi, li):
        rec = cache[wavs[wi]]
        s = similarity(rec["english"], lines[ids[li]]["content"])
        # secondary signal: longer lines take longer to say
        chars = len(lines[ids[li]]["content"])
        expect = 0.9 + chars / 13.0
        d = rec.get("duration") or 0.0
        agree = max(0.0, 1.0 - abs(d - expect) / max(expect, 1.0))
        w = args.dur_weight
        return (1.0 - w) * s + w * agree

    pairs, total = align(wavs, ids, score, args.skip_line)
    print("alignment score %.2f over %d recordings and %d lines" % (total, len(wavs), len(ids)))

    # Verified manual corrections. Whisper's translate pass drops clauses on short
    # utterances, which can pull the aligner off by a line; overrides record the
    # evidence for each correction rather than hiding it in a tuned heuristic.
    over = {}
    if os.path.exists(args.overrides):
        with open(args.overrides, encoding="utf-8") as fh:
            over = json.load(fh).get(args.prefix) or {}
    if over:
        pos = {lid: j for j, lid in enumerate(ids)}
        applied = []
        for k, (wi, li) in enumerate(pairs):
            n = str(int(wavs[wi]))
            if n in over:
                pairs[k] = (wi, pos[over[n]["line_id"]])
                applied.append("wav %s -> id %d" % (n, over[n]["line_id"]))
        print("applied %d verified override(s): %s" % (len(applied), ", ".join(applied)))
    print()

    rows, flagged, unmatched = [], [], []
    for wi, li in pairs:
        n = int(wavs[wi])
        rec = cache[wavs[wi]]
        if li is None:
            unmatched.append(n)
            continue
        lid = ids[li]
        sim = similarity(rec["english"], lines[lid]["content"])
        row = dict(wav=rec["wav"], wav_n=n, line_id=lid, offset=lid - n,
                   speaker=(lines[lid].get("teller") or "").strip(),
                   duration=rec["duration"],
                   speech_start=rec["speech_start"], speech_end=rec["speech_end"],
                   language=rec["language"], language_probability=rec["language_probability"],
                   similarity=round(sim, 4), native_text=rec["text"],
                   english_text=rec["english"], line_text=lines[lid]["content"])
        rows.append(row)
        if sim < args.low:
            flagged.append(row)

    out_csv = os.path.join(args.asr_dir, f"{args.prefix.lower()}-vo-mapping.csv")
    with open(out_csv, "w", encoding="utf-8-sig", newline="") as fh:
        w = csv.DictWriter(fh, fieldnames=list(rows[0]))
        w.writeheader()
        w.writerows(rows)
    with open(out_csv.replace(".csv", ".json"), "w", encoding="utf-8") as fh:
        json.dump(rows, fh, ensure_ascii=False, indent=1)

    print("OFFSET (line_id - wav_number) BY RANGE")
    runs, cur = [], None
    for r in rows:
        if cur and cur[2] == r["offset"]:
            cur[1] = r["wav_n"]
        else:
            cur = [r["wav_n"], r["wav_n"], r["offset"]]
            runs.append(cur)
    for lo, hi, off in runs:
        print("  wav %2d-%2d  ->  offset %+d" % (lo, hi, off))

    print("\nmean similarity %.3f; %d of %d pairs below %.2f"
          % (sum(r["similarity"] for r in rows) / len(rows), len(flagged), len(rows), args.low))
    if unmatched:
        print("UNMATCHED recordings: %s" % unmatched)
    if flagged:
        print("\nFLAGGED FOR HUMAN REVIEW (short or non-verbal utterances score low by nature)")
        for r in flagged:
            print("  %s -> id %-3d (%.2f)  heard %-42s | line %s"
                  % (r["wav"], r["line_id"], r["similarity"],
                     (r["english_text"] or "")[:42], r["line_text"][:42]))
    print("\nwrote %s" % out_csv)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
