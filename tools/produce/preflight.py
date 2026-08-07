"""Check every shot's inputs resolve before committing to a long render pass.

A full clip pass is ~9 hours. Discovering a missing reference or a stale
workflow path at shot 30 wastes most of it, so verify all 45 up front.

Checks per shot: both reference images exist, the workflow named in the prompt
exists, and any R2V-A shot has both a driving recording on disk and a voiced
line in the EDL to position it with.

Run:  tools/refs/.venv/Scripts/python.exe tools/produce/preflight.py
"""

import csv
import json
import os
import re
import sys

REPO = os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
PROD = os.path.join(REPO, "production", "limbus-prologue")


def main():
    rows = list(csv.DictReader(
        open(os.path.join(PROD, "shots", "manifest.csv"), encoding="utf-8")))
    edl = json.load(open(os.path.join(PROD, "audio", "edl.json"), encoding="utf-8"))
    voiced = set()
    for l in edl["lines"]:
        if l["delivery"] == "vo" and l["wav"]:
            voiced.add(l["shot_id"])

    bad = []
    for r in rows:
        sid = r["shot_id"]
        p = os.path.join(PROD, "shots", sid, f"{sid}_MINIMAX_v001.prompt.md")
        if not os.path.exists(p):
            bad.append(f"{sid}: no prompt file")
            continue
        fm = re.match(r"^---\n(.*?)\n---\n", open(p, encoding="utf-8").read(),
                      re.S).group(1)

        def g(key):
            m = re.search(rf"^{key}:\s*(\S+)", fm, re.M)
            return m.group(1) if m else None

        for key in ("picture_1_subject", "picture_2_location", "workflow"):
            v = g(key)
            if not v or not os.path.exists(os.path.join(REPO, v)):
                bad.append(f"{sid}: {key} -> {v}")

        if g("route") == "R2V-A":
            auds = re.findall(r"^\s+-\s+(input/\S+\.wav)", fm, re.M)
            if not auds:
                bad.append(f"{sid}: route R2V-A but no reference_audio")
            elif not os.path.exists(os.path.join(REPO, auds[0])):
                bad.append(f"{sid}: audio missing {auds[0]}")
            if sid not in voiced:
                bad.append(f"{sid}: route R2V-A but no voiced line in the EDL")

    print(f"preflight: {len(rows)} shots, {len(bad)} problems")
    for b in bad:
        print("  " + b)
    return 1 if bad else 0


if __name__ == "__main__":
    sys.exit(main())
