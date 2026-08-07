#!/usr/bin/env python
"""Validate the prologue VO streams locally with faster-whisper.

Runs entirely offline of ComfyUI. For every S001B-*.wav this records:

  * detected language and confidence
  * the native transcript (word-level timings)
  * an English translation, used to score the WAV against candidate dialogue lines
  * true speech in/out points, which replace guessed lead-in/tail padding

and then tests the wav-number -> dialogue-id mapping at several offsets so an
off-by-one is detected rather than assumed.

    tools/asr/.venv/Scripts/python.exe tools/asr/transcribe_vo.py
    ... --model large-v3 --device cuda --out tools/asr/out
"""
from __future__ import annotations

import argparse
import csv
import difflib
import json
import os
import re
import sys
import unicodedata

sys.stdout.reconfigure(encoding="utf-8")

REPO = os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
DEF_DIALOGUE = os.path.join(
    REPO, "input", "limbuscompany", "dialogues", "EN",
    "00-prologue-selva-oscura", "tutorial", "EN_S001B.json")
DEF_AUDIO = os.path.join(REPO, "input", "limbuscompany", "prologue", "audio")

# current working hypothesis, from the duration/length regression
def mapped_id(wav_n: int) -> int:
    return wav_n - 1 if wav_n < 48 else wav_n


OFFSETS = (-2, -1, 0, 1, 2)


def norm(s: str) -> str:
    """Fold to comparable word tokens: strip markup, case, punctuation, redaction blocks."""
    s = unicodedata.normalize("NFKC", s or "")
    s = s.replace("■", " ")                 # ■ redaction blocks
    s = re.sub(r"[<>\[\]]", " ", s)
    s = s.lower()
    s = re.sub(r"[^\w\s]", " ", s)
    return re.sub(r"\s+", " ", s).strip()


def similarity(a: str, b: str) -> float:
    """Blend sequence ratio with token overlap; robust to loose localisation."""
    na, nb = norm(a), norm(b)
    if not na or not nb:
        return 0.0
    seq = difflib.SequenceMatcher(None, na, nb).ratio()
    ta, tb = set(na.split()), set(nb.split())
    jac = len(ta & tb) / len(ta | tb) if (ta | tb) else 0.0
    return 0.5 * seq + 0.5 * jac


def load_lines(path: str) -> dict[int, dict]:
    with open(path, encoding="utf-8-sig") as fh:
        return {e["id"]: e for e in json.load(fh)["dataList"]}


def wav_files(audio_dir: str, prefix: str) -> list[tuple[int, str]]:
    out = []
    for name in os.listdir(audio_dir):
        m = re.fullmatch(prefix + r"-(\d+)\.wav", name, re.I)
        if m:
            out.append((int(m.group(1)), os.path.join(audio_dir, name)))
    return sorted(out)


def main() -> int:
    ap = argparse.ArgumentParser(description=__doc__,
                                 formatter_class=argparse.RawDescriptionHelpFormatter)
    ap.add_argument("--audio-dir", default=DEF_AUDIO)
    ap.add_argument("--dialogue", default=DEF_DIALOGUE)
    ap.add_argument("--prefix", default="S001B")
    ap.add_argument("--model", default="large-v3")
    ap.add_argument("--device", default="cuda", choices=["cuda", "cpu", "auto"])
    ap.add_argument("--compute-type", default=None,
                    help="default float16 on cuda, int8 on cpu")
    ap.add_argument("--beam-size", type=int, default=5)
    ap.add_argument("--out", default=os.path.join(REPO, "tools", "asr", "out"))
    ap.add_argument("--limit", type=int, default=0, help="only the first N files")
    ap.add_argument("--force", action="store_true", help="ignore the cache")
    args = ap.parse_args()

    device = args.device
    if device == "auto":
        import ctranslate2
        device = "cuda" if ctranslate2.get_cuda_device_count() > 0 else "cpu"
    compute = args.compute_type or ("float16" if device == "cuda" else "int8")

    os.makedirs(args.out, exist_ok=True)
    cache_path = os.path.join(args.out, f"{args.prefix.lower()}-asr-cache.json")
    cache = {}
    if os.path.exists(cache_path) and not args.force:
        with open(cache_path, encoding="utf-8") as fh:
            cache = json.load(fh)

    lines = load_lines(args.dialogue)
    files = wav_files(args.audio_dir, args.prefix)
    if args.limit:
        files = files[: args.limit]
    if not files:
        print(f"no {args.prefix}-*.wav under {args.audio_dir}", file=sys.stderr)
        return 2

    from faster_whisper import WhisperModel
    print(f"loading {args.model} on {device}/{compute} ...", flush=True)
    model = WhisperModel(args.model, device=device, compute_type=compute)

    def run(path, task):
        segs, info = model.transcribe(path, task=task, beam_size=args.beam_size,
                                      word_timestamps=(task == "transcribe"),
                                      vad_filter=False)
        segs = list(segs)
        text = " ".join(s.text.strip() for s in segs).strip()
        starts = [w.start for s in segs for w in (s.words or [])]
        ends = [w.end for s in segs for w in (s.words or [])]
        return dict(text=text, language=info.language,
                    language_probability=round(info.language_probability, 4),
                    duration=round(info.duration, 3),
                    speech_start=round(min(starts), 3) if starts else None,
                    speech_end=round(max(ends), 3) if ends else None)

    todo = [(n, p) for n, p in files if str(n) not in cache]
    print(f"{len(files)} files, {len(todo)} to transcribe, {len(files)-len(todo)} cached\n", flush=True)
    for i, (n, path) in enumerate(files, 1):
        if str(n) in cache:
            continue
        native = run(path, "transcribe")
        english = run(path, "translate")
        cache[str(n)] = dict(wav=os.path.basename(path), **native, english=english["text"])
        print("[%2d/%2d] %s  %s p=%.2f  %.2fs  %s"
              % (i, len(files), os.path.basename(path), native["language"],
                 native["language_probability"], native["duration"],
                 native["text"][:60]), flush=True)
        with open(cache_path, "w", encoding="utf-8") as fh:
            json.dump(cache, fh, ensure_ascii=False, indent=1)

    # ---------------------------------------------------------------- analysis
    rows = []
    for n, path in files:
        rec = cache[str(n)]
        scores = {}
        for off in OFFSETS:
            e = lines.get(n + off)
            scores[off] = similarity(rec["english"], e["content"]) if e else 0.0
        best_off = max(scores, key=scores.get)
        cur = mapped_id(n)
        cur_line = lines.get(cur)
        rows.append(dict(
            wav=rec["wav"], wav_n=n, language=rec["language"],
            language_probability=rec["language_probability"], duration=rec["duration"],
            speech_start=rec["speech_start"], speech_end=rec["speech_end"],
            native_text=rec["text"], english_text=rec["english"],
            mapped_id=cur, mapped_speaker=(cur_line or {}).get("teller", ""),
            mapped_text=(cur_line or {}).get("content", ""),
            mapped_score=round(scores[mapped_id(n) - n], 4) if (mapped_id(n) - n) in scores else None,
            best_offset=best_off, best_id=n + best_off,
            best_score=round(scores[best_off], 4),
            agrees=(best_off == mapped_id(n) - n),
            scores={str(k): round(v, 4) for k, v in scores.items()}))

    with open(os.path.join(args.out, f"{args.prefix.lower()}-vo-manifest.csv"),
              "w", encoding="utf-8-sig", newline="") as fh:
        w = csv.DictWriter(fh, fieldnames=[k for k in rows[0] if k != "scores"])
        w.writeheader()
        for r in rows:
            w.writerow({k: v for k, v in r.items() if k != "scores"})
    with open(os.path.join(args.out, f"{args.prefix.lower()}-vo-manifest.json"),
              "w", encoding="utf-8") as fh:
        json.dump(rows, fh, ensure_ascii=False, indent=1)

    # ---------------------------------------------------------------- report
    from collections import Counter
    langs = Counter(r["language"] for r in rows)
    print("\n" + "=" * 72)
    print("LANGUAGE      :", ", ".join("%s x%d" % (k, v) for k, v in langs.most_common()))
    lp = [r["language_probability"] for r in rows]
    print("               mean confidence %.3f, min %.3f" % (sum(lp) / len(lp), min(lp)))

    print("\nMAPPING — mean translation similarity by offset")
    for off in OFFSETS:
        early = [r["scores"][str(off)] for r in rows if r["wav_n"] < 48]
        late = [r["scores"][str(off)] for r in rows if r["wav_n"] >= 48]
        print("  offset %+d   wav<48: %.3f (n=%d)   wav>=48: %.3f (n=%d)%s" % (
            off, sum(early) / len(early) if early else 0, len(early),
            sum(late) / len(late) if late else 0, len(late),
            "   <-- current rule" if off in (-1, 0) else ""))

    dis = [r for r in rows if not r["agrees"]]
    print("\nDISAGREEMENTS with the current rule: %d of %d" % (len(dis), len(rows)))
    for r in dis:
        print("  %s  mapped id %d (%.2f) -> best id %d (%.2f, offset %+d)"
              % (r["wav"], r["mapped_id"], r["mapped_score"] or 0,
                 r["best_id"], r["best_score"], r["best_offset"]))
        print("      heard    : %s" % r["english_text"][:90])
        print("      mapped to: %s" % (r["mapped_text"] or "")[:90])

    lead = [r["speech_start"] for r in rows if r["speech_start"] is not None]
    tail = [r["duration"] - r["speech_end"] for r in rows
            if r["speech_end"] is not None and r["duration"]]
    if lead:
        print("\nSPEECH BOUNDARIES (replaces guessed padding)")
        print("  lead-in  mean %.3fs  max %.3fs" % (sum(lead) / len(lead), max(lead)))
        print("  tail     mean %.3fs  max %.3fs" % (sum(tail) / len(tail), max(tail)))
    print("\nwrote %s" % os.path.join(args.out, f"{args.prefix.lower()}-vo-manifest.csv"))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
