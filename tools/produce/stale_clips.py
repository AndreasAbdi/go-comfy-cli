"""List clips whose inputs changed after they were rendered.

Fixes landed mid-pass, so some clips were generated against prompts that no
longer exist. Blanket-re-rolling everything wastes ~5 min per shot; this works
out which shots actually changed, by comparing each clip's render time against
when each fix went in and whether that fix touches that shot at all.

A WIDE shot, for instance, is untouched by the plate-zoom fix because its zoom
factor is 1.0 -- the plate it gets is the same image it always got.

Run:  tools/refs/.venv/Scripts/python.exe tools/produce/stale_clips.py
Then: clips.py --attempt 2 --shots $(that list)
"""

import csv
import datetime
import os
import re
import sys

REPO = os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
PROD = os.path.join(REPO, "production", "limbus-prologue")

# When each fix landed, and which shots it can possibly affect.
FIXES = [
    # (name, time applied, predicate(framing, prompt_text) -> affected?)
    ("fallen", datetime.time(14, 25),
     lambda fr, mm: "CHAR-DANTE-FALLEN-01.png" in mm),
    ("plate", datetime.time(14, 59),
     lambda fr, mm: fr not in ("WIDE", "EST")),      # zoom 1.0 = unchanged
    ("ecu-trim", datetime.time(15, 20),
     lambda fr, mm: fr in ("ECU", "INS")),
]

LOGS = ["out/full_pass_480.log", "out/full_pass.log"]


def render_times():
    """shot_id -> time it was rendered, across every pass log."""
    out = {}
    for lg in LOGS:
        p = os.path.join(REPO, lg)
        if not os.path.exists(p):
            continue
        for line in open(p, encoding="utf-8"):
            m = re.match(r"\[(\d\d):(\d\d):\d\d\].*?(SQ00_SC01_SH\d+).*\bok\b", line)
            if m:
                t = datetime.time(int(m.group(1)), int(m.group(2)))
                # earliest render wins: a later entry is a re-roll of the same
                # shot and would mask the original staleness
                out.setdefault(m.group(3), t)
    return out


def main():
    rows = {r["shot_id"]: r for r in csv.DictReader(
        open(os.path.join(PROD, "shots", "manifest.csv"), encoding="utf-8"))}
    times = render_times()

    stale = []
    for sid, t in sorted(times.items()):
        if sid not in rows:
            continue
        fr = rows[sid]["framing"].strip()
        mm = open(os.path.join(PROD, "shots", sid,
                               f"{sid}_MINIMAX_v001.prompt.md"), encoding="utf-8").read()
        why = [name for name, applied, pred in FIXES if t < applied and pred(fr, mm)]
        if why:
            stale.append((sid, fr, ",".join(why)))

    for sid, fr, why in stale:
        print(f"  {sid} {fr:<4} {why}")
    nums = ",".join(s.split("SH")[1] for s, _, _ in stale)
    print(f"\n{len(stale)} stale of {len(times)} rendered")
    if nums:
        print(f"\nclips.py --attempt 2 --shots {nums}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
