"""Filmstrip + motion score for one rendered clip, for quick iteration.

Motion score is mean absolute frame-to-frame luma delta scaled 0-100. It exists
because "did the shot move" is otherwise invisible in a contact sheet: several
early clips looked fine as stills and were effectively frozen video.

Usage:  inspect_shot.py 050 [--attempt 1] [--out strip.png]
"""

import argparse
import os
import subprocess
import sys

import numpy as np
from PIL import Image, ImageDraw

REPO = os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
PROD = os.path.join(REPO, "production", "limbus-prologue")


def motion_score(path):
    p = subprocess.run(["ffmpeg", "-v", "error", "-i", path, "-vf",
                        "scale=128:72,format=gray", "-f", "rawvideo", "-"],
                       capture_output=True)
    a = np.frombuffer(p.stdout, dtype=np.uint8)
    n = a.size // (128 * 72)
    if n < 2:
        return 0.0, n
    a = a[: n * 128 * 72].reshape(n, -1).astype(np.int16)
    return round(float(np.abs(np.diff(a, axis=0)).mean()) * 100 / 255, 2), n


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("shot")
    ap.add_argument("--attempt", type=int, default=1)
    ap.add_argument("--out", default=None)
    ap.add_argument("--count", type=int, default=6)
    a = ap.parse_args()

    sid = f"SQ00_SC01_SH{a.shot.zfill(3)}"
    v = os.path.join(PROD, "shots", sid, "clips", f"{sid}_CLIP_a{a.attempt}.mp4")
    if not os.path.exists(v):
        raise SystemExit(f"no clip: {os.path.relpath(v, REPO)}")

    ms, n = motion_score(v)
    idxs = [round(i * (n - 1) / (a.count - 1)) for i in range(a.count)]
    tmp = os.path.join(PROD, "shots", sid, "clips", "_inspect")
    os.makedirs(tmp, exist_ok=True)
    for f in os.listdir(tmp):
        os.remove(os.path.join(tmp, f))
    sel = "+".join(f"eq(n\\,{i})" for i in idxs)
    subprocess.run(["ffmpeg", "-y", "-v", "error", "-i", v, "-vf",
                    f"select='{sel}'", "-vsync", "0",
                    os.path.join(tmp, "f%02d.png")], capture_output=True)

    fs = sorted(f for f in os.listdir(tmp) if f.endswith(".png"))
    ims = [Image.open(os.path.join(tmp, f)).convert("RGB") for f in fs]
    h = 340
    ts = [im.resize((round(im.width * h / im.height), h), Image.LANCZOS) for im in ims]
    gap = 6
    W = sum(t.width for t in ts) + gap * (len(ts) + 1)
    sh = Image.new("RGB", (W, h + 24 + gap * 2), (20, 20, 24))
    d = ImageDraw.Draw(sh)
    x = gap
    for t in ts:
        sh.paste(t, (x, gap))
        x += t.width + gap
    d.text((gap + 2, h + gap + 6), f"{sid}  {n} frames  motion={ms}",
           fill=(225, 225, 230))
    dest = a.out or os.path.join(PROD, "review", f"{sid}_inspect.png")
    os.makedirs(os.path.dirname(dest), exist_ok=True)
    sh.save(dest)
    print(f"{sid}: {n} frames, motion={ms} -> {os.path.relpath(dest, REPO)}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
