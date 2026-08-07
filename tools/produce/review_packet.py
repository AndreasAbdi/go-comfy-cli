"""Build a review packet a second model can actually judge.

A reviewer cannot watch an mp4, so this renders what matters into stills and
numbers: a filmstrip per shot (start / quarter / middle / three-quarter / end),
a motion score so "it rendered a still image for five seconds" is visible
without watching, and a manifest tying each shot to its intent, references and
dialogue.

Motion is reported three ways because a single number misleads: `motion` (mean
luma delta) under-reports a talking close-up, `changing` (share of pixels moving
more than 8 levels between frames) shows whether anything moves at all, and
`endpoints` (share differing between first and last frame) catches a clip that
goes nowhere. Early clips looked fine as stills and were effectively frozen
video; the reviewer must be able to see that without taking it on trust.

Usage:  review_packet.py [--attempt 1]
"""

import argparse
import csv
import json
import os
import re
import subprocess
import sys

import numpy as np
from PIL import Image, ImageDraw

REPO = os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
PROD = os.path.join(REPO, "production", "limbus-prologue")
OUT = os.path.join(PROD, "review")
STRIP_H = 260


def run(cmd):
    return subprocess.run(cmd, capture_output=True, text=True, encoding="utf-8",
                          errors="replace")


def frames_at(video, idxs, dest_dir):
    os.makedirs(dest_dir, exist_ok=True)
    sel = "+".join(f"eq(n\\,{i})" for i in idxs)
    run(["ffmpeg", "-y", "-v", "error", "-i", video, "-vf", f"select='{sel}'",
         "-vsync", "0", os.path.join(dest_dir, "f%02d.png")])
    return sorted(os.path.join(dest_dir, f) for f in os.listdir(dest_dir)
                  if f.endswith(".png"))


def motion_stats(video):
    """How much the shot actually moves.

    Mean luma delta alone under-reports a talking close-up: a moving mouth is a
    small area and barely shifts the mean. `changing` -- the share of pixels
    that shift by more than 8 levels between consecutive frames -- separates
    "small thing moving" from "nothing moving", and `endpoints` catches a clip
    that drifts nowhere across its whole duration.
    """
    p = subprocess.run(["ffmpeg", "-v", "error", "-i", video, "-vf",
                        "scale=192:108,format=gray", "-f", "rawvideo", "-"],
                       capture_output=True)
    a = np.frombuffer(p.stdout, dtype=np.uint8)
    n = a.size // (192 * 108)
    if n < 2:
        return dict(motion=0.0, changing=0.0, endpoints=0.0, frames=n)
    a = a[: n * 192 * 108].reshape(n, -1).astype(np.int16)
    d = np.abs(np.diff(a, axis=0))
    return dict(
        motion=round(float(d.mean()) * 100 / 255, 2),
        changing=round(float((d > 8).mean()) * 100, 2),
        endpoints=round(float((np.abs(a[0] - a[-1]) > 8).mean()) * 100, 2),
        frames=n,
    )


def strip(paths, dest, label):
    ims = [Image.open(p).convert("RGB") for p in paths]
    ts = [im.resize((round(im.width * STRIP_H / im.height), STRIP_H), Image.LANCZOS)
          for im in ims]
    gap = 6
    w = sum(t.width for t in ts) + gap * (len(ts) + 1)
    sh = Image.new("RGB", (w, STRIP_H + 26 + gap * 2), (18, 18, 22))
    d = ImageDraw.Draw(sh)
    x = gap
    for t in ts:
        sh.paste(t, (x, gap))
        x += t.width + gap
    d.text((gap + 2, STRIP_H + gap + 6), label, fill=(220, 220, 226))
    sh.save(dest)
    return sh.size


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--attempt", type=int, default=1)
    a = ap.parse_args()

    edl = json.load(open(os.path.join(PROD, "audio", "edl.json"), encoding="utf-8"))
    rows = {r["shot_id"]: r for r in csv.DictReader(
        open(os.path.join(PROD, "shots", "manifest.csv"), encoding="utf-8"))}
    lines_by_shot = {}
    for l in edl["lines"]:
        lines_by_shot.setdefault(l["shot_id"], []).append(l)

    os.makedirs(OUT, exist_ok=True)
    tmp = os.path.join(OUT, "_frames")
    report, missing = [], []

    for s in edl["shots"]:
        sid = s["shot_id"]
        v = os.path.join(PROD, "shots", sid, "clips", f"{sid}_CLIP_a{a.attempt}.mp4")
        if not os.path.exists(v):
            missing.append(sid)
            continue
        n = s["frames"]
        idxs = [0, n // 4, n // 2, (3 * n) // 4, n - 1]
        fs = frames_at(v, idxs, os.path.join(tmp, sid))
        st = motion_stats(v)
        r = rows[sid]
        label = (f"{sid}  {s['seconds']:.2f}s  {r['framing']}  {r['camera']}  "
                 f"motion={st['motion']} changing={st['changing']}% "
                 f"endpoints={st['endpoints']}%")
        strip(fs, os.path.join(OUT, f"{sid}_strip.png"), label)
        st.pop("frames", None)          # the edit's frame count wins, not the raw clip's
        report.append(dict(
            shot_id=sid, seconds=s["seconds"], frames=n, route=s["route"],
            framing=r["framing"], camera=r["camera"], **st,
            transition_out=r["transition_out"],
            references=r["reference_images"].split(),
            lines=[dict(speaker=l["speaker"], delivery=l["delivery"],
                        text=l["text"]) for l in lines_by_shot.get(sid, [])],
        ))

    # A single grid of one frame per shot, in cut order. Style consistency and
    # palette drift across the whole cut are visible here and not in 45
    # separate filmstrips.
    if report:
        tiles = []
        for r in report:
            mid = os.path.join(tmp, r["shot_id"], "f03.png")
            if not os.path.exists(mid):
                cand = sorted(f for f in os.listdir(os.path.join(tmp, r["shot_id"]))
                              if f.endswith(".png"))
                mid = os.path.join(tmp, r["shot_id"], cand[len(cand) // 2])
            tiles.append((r["shot_id"].split("SH")[1], Image.open(mid).convert("RGB")))
        tw, th, cols = 320, 180, 5
        rows_n = (len(tiles) + cols - 1) // cols
        board = Image.new("RGB", (cols * (tw + 8) + 8, rows_n * (th + 24) + 8),
                          (16, 16, 20))
        d = ImageDraw.Draw(board)
        for i, (num, im) in enumerate(tiles):
            x = 8 + (i % cols) * (tw + 8)
            y = 8 + (i // cols) * (th + 24)
            board.paste(im.resize((tw, th), Image.LANCZOS), (x, y))
            d.text((x + 2, y + th + 5), f"SH{num}", fill=(210, 210, 216))
        board.save(os.path.join(OUT, "storyboard.png"))
        print(f"storyboard: {len(tiles)} shots -> review/storyboard.png")

    with open(os.path.join(OUT, "report.json"), "w", encoding="utf-8") as f:
        json.dump(dict(attempt=a.attempt, missing=missing, shots=report),
                  f, ensure_ascii=False, indent=2)

    if report:
        ch = sorted(r["changing"] for r in report)
        frozen = [r["shot_id"] for r in report if r["changing"] < 0.5]
        print(f"packet: {len(report)} shots, {len(missing)} missing")
        print(f"changing pixels %: min={ch[0]} median={ch[len(ch)//2]} max={ch[-1]}")
        if frozen:
            print(f"NEARLY FROZEN (<0.5% changing): {', '.join(frozen)}")
    else:
        print(f"no clips found for attempt {a.attempt}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
