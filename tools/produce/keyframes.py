"""Render shot keyframes with qwen-image-edit via go-comfy-cli.

Reads each shot's KEY prompt frontmatter for its plate, character references and
negative prompt, resolves the single reference slot (see below), and submits.

The reference slot is a bottleneck: qwen-image-edit has exactly two LoadImage
nodes. The plate takes one. For the other:
  - 1 character  -> that character's crop
  - 2+ characters -> the composited sheet from sheets.py
  - 0 characters  -> a prologue story CG, used purely as a style anchor, since
                     the slot cannot be left empty and the model needs all the
                     help it can get staying off photorealism.

Seeds are derived from the shot number and attempt so a rerun reproduces and a
re-roll is a deliberate `--attempt` bump rather than luck.

Usage:
  keyframes.py                      # every shot missing an approved keyframe
  keyframes.py --shots 010,020      # specific shots
  keyframes.py --attempt 2          # re-roll with different seeds
  keyframes.py --dry-run
"""

import argparse
import csv
import os
import re
import shutil
import subprocess
import sys
import time

REPO = os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
PROD = os.path.join(REPO, "production", "limbus-prologue")
CLI = os.path.join(REPO, "go-comfy-cli.exe")
WF = os.path.join(REPO, "workflows", "qwen-image-edit", "qwen-image-edit.json")
ARGS = os.path.join(REPO, "workflows", "qwen-image-edit", "qwen-image-edit.args.yaml")
STYLE_ANCHOR = "production/limbus-prologue/references/images/prologue/PRO-PAGE1-01.png"
LOG = os.path.join(REPO, "progress.txt")


def log(msg):
    line = f"[{time.strftime('%H:%M:%S')}] {msg}"
    print(line)
    with open(LOG, "a", encoding="utf-8") as f:
        f.write(line + "\n")


def frontmatter(path):
    txt = open(path, encoding="utf-8").read()
    m = re.match(r"^---\n(.*?)\n---\n", txt, re.S)
    return m.group(1), txt[m.end():].strip()


def shot_spec(sid):
    fm, body = frontmatter(
        os.path.join(PROD, "shots", sid, f"{sid}_KEY_v001.prompt.md"))
    plate = re.search(r"^input_image:\s*(\S+)", fm, re.M).group(1)
    chars = re.findall(r"^\s+-\s+(\S+\.png)", fm, re.M)
    neg = re.search(r"^negative_prompt: >-\n\s+(.*?)$", fm, re.M | re.S)
    neg = " ".join(neg.group(1).split()) if neg else ""
    return plate, chars, neg, body


# A full-body crop is the wrong reference for a tight insert: on an extreme
# close-up of Dante's dial the model gets a ~100px clock and invents the rest,
# which is how the first SH010 attempt grew a human face in the rim. Shots this
# tight on a single subject get a detail crop instead.
TIGHT = {"ECU", "INS", "CU"}
DETAIL = {
    "CHAR-DANTE-01":
        "production/limbus-prologue/references/images/characters/CHAR-DANTE-HEAD-01.png",
}


def reference_for(sid, chars, framing):
    if len(chars) == 1:
        rid = os.path.basename(chars[0]).replace(".png", "")
        if framing in TIGHT and rid in DETAIL:
            return DETAIL[rid]
        return chars[0]
    if len(chars) > 1:
        sheet = f"production/limbus-prologue/references/sheets/{sid}_SHEET.png"
        if not os.path.exists(os.path.join(REPO, sheet)):
            raise SystemExit(f"{sid}: sheet missing, run tools/produce/sheets.py")
        return sheet
    return STYLE_ANCHOR


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--shots", help="comma-separated shot numbers, e.g. 010,020")
    ap.add_argument("--attempt", type=int, default=1)
    ap.add_argument("--dry-run", action="store_true")
    ap.add_argument("--force", action="store_true", help="re-render even if present")
    a = ap.parse_args()

    rows = list(csv.DictReader(
        open(os.path.join(PROD, "shots", "manifest.csv"), encoding="utf-8")))
    if a.shots:
        want = {s.strip().zfill(3) for s in a.shots.split(",")}
        rows = [r for r in rows if r["shot_id"].split("SH")[1] in want]

    todo = []
    for r in rows:
        sid = r["shot_id"]
        dest = os.path.join(PROD, "shots", sid, "keyframes",
                            f"{sid}_KEY_v001_a{a.attempt}.png")
        if os.path.exists(dest) and not a.force:
            continue
        todo.append((r, dest))

    log(f"KEYFRAMES: {len(todo)} to render (attempt {a.attempt})")
    ok, fail = 0, []
    for i, (r, dest) in enumerate(todo, 1):
        sid = r["shot_id"]
        num = int(sid.split("SH")[1])
        plate, chars, neg, body = shot_spec(sid)
        ref = reference_for(sid, chars, r["framing"].strip())
        seed = num * 1000 + a.attempt

        work = os.path.join(PROD, "shots", sid, "keyframes")
        os.makedirs(work, exist_ok=True)
        ptxt = os.path.join(work, f"{sid}_KEY.submit.txt")
        with open(ptxt, "w", encoding="utf-8") as f:
            f.write(body)

        outdir = os.path.join(work, f"_raw_a{a.attempt}")
        cmd = [CLI, "run", "--workflow", WF, "--args-file", ARGS,
               "--set", f"input_image={os.path.join(REPO, plate)}",
               "--set", f"reference_image={os.path.join(REPO, ref)}",
               "--set", f"positive_prompt={ptxt}",
               "--set", f"negative_prompt={neg}",
               "--set", f"seed={seed}",
               "--output-folder", outdir]
        label = (f"[{i}/{len(todo)}] {sid} seed={seed} "
                 f"plate={os.path.basename(plate)} ref={os.path.basename(ref)}")
        if a.dry_run:
            print(label)
            continue

        t0 = time.time()
        p = subprocess.run(cmd, capture_output=True, text=True, encoding="utf-8",
                           errors="replace")
        if p.returncode != 0:
            log(f"{label} FAILED rc={p.returncode}: {(p.stderr or p.stdout)[-300:]}")
            fail.append(sid)
            continue
        pngs = sorted(f for f in os.listdir(outdir)) if os.path.isdir(outdir) else []
        pngs = [f for f in pngs if f.endswith(".png")]
        if not pngs:
            log(f"{label} FAILED: no output png")
            fail.append(sid)
            continue
        shutil.copyfile(os.path.join(outdir, pngs[-1]), dest)
        log(f"{label} ok {time.time() - t0:.0f}s -> {os.path.relpath(dest, REPO)}")
        ok += 1

    log(f"KEYFRAMES done: {ok} ok, {len(fail)} failed {fail if fail else ''}")
    return 1 if fail else 0


if __name__ == "__main__":
    sys.exit(main())
