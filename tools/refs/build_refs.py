"""Collect the Limbus prologue reference images into the production tree.

Character art ships in two shapes: the fixers (Lion/Wolf/Panther) are clean
full-body sprites, while the Sinners are dialogue *atlases* -- a grid of face
variants plus one full-body figure plus a column of spare heads. Feeding an
atlas to a generator as a character reference produces a collage, so the body
has to be cut out first. This finds it by labelling the alpha channel and
keeping the tallest connected component, which is the standing figure in every
atlas in this set.

Run:  tools/refs/.venv/Scripts/python.exe tools/refs/build_refs.py
"""

import os
import shutil
import sys

import numpy as np
from PIL import Image
from scipy import ndimage

Image.MAX_IMAGE_PIXELS = None

REPO = os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
SRC = os.path.join(REPO, "input", "limbuscompany")
CHARS = os.path.join(SRC, "art-assets", "characters")
PROLOGUE = os.path.join(SRC, "prologue", "art-assets")
OUT = os.path.join(REPO, "production", "limbus-prologue", "references", "images")

ALPHA_MIN = 24      # below this an atlas' antialiased fringe starts linking cells
MIN_FRAC = 0.010    # ignore components under 1% of the canvas height
PAD = 8             # transparent margin kept around a crop
MAX_EDGE = 1536     # downscale ceiling; references do not need atlas resolution


# key -> (source path relative to CHARS, extract mode)
#   "sprite" = already a single figure, just trim to the alpha bounding box
#   "atlas"  = cut the tallest connected component out of a sprite sheet
CHARACTERS = {
    "CHAR-DANTE-01":     ("Dante/Dante.png", "atlas"),
    "CHAR-FAUST-01":     ("Faust/Faust.png", "atlas"),
    "CHAR-LION-01":      ("Lion/Lion.png", "sprite"),
    "CHAR-WOLF-01":      ("Wolf/Wolf 108225.png", "sprite"),
    "CHAR-PANTHER-01":   ("Panther/Panther.png", "sprite"),
    "CHAR-YISANG-01":    ("Yi Sang/YiSang.png", "atlas"),
    "CHAR-GREGOR-01":    ("Gregor/Gregor.png", "atlas"),
    "CHAR-RODION-01":    ("Rodion/Rodion.png", "atlas"),
    "CHAR-ISHMAEL-01":   ("Ishmael/Ishmael.png", "atlas"),
    "CHAR-HEATHCLIFF-01": ("Heathcliff/Heathcliff.png", "atlas"),
    "CHAR-OUTIS-01":     ("Outis/Outis.png", "atlas"),
    "CHAR-DONQ-01":      ("Don Quixote/DonQuixote.png", "atlas"),
    "CHAR-MEURSAULT-01": ("Meursault/Meursault.png", "atlas"),
    "CHAR-RYOSHU-01":    ("Ryoshu/Ryoshu.png", "atlas"),
    "CHAR-HONGLU-01":    ("Hong Lu/HongLu.png", "atlas"),
    "CHAR-SINCLAIR-01":  ("Sinclair/Sinclair.png", "atlas"),
    "CHAR-VERGILIUS-01": ("Vergilius/Vergilius.png", "atlas"),
}

# Official LCB Sinner identity art -- the canonical uniform, kept as a costume
# authority alongside the sprite crop.
COSTUMES = {
    "CHAR-FAUST-01":      "Faust/Identities/LCB Sinner/10201_normal.png",
    "CHAR-YISANG-01":     "Yi Sang/Identities/LCB Sinner/10101_normal.png",
    "CHAR-GREGOR-01":     "Gregor/Identities/LCB Sinner/11201_normal.png",
    "CHAR-RODION-01":     "Rodion/Identities/LCB Sinner/10901_normal.png",
    "CHAR-ISHMAEL-01":    "Ishmael/Identities/LCB Sinner/10801_normal.png",
    "CHAR-HEATHCLIFF-01": "Heathcliff/Identities/LCB Sinner/10701_normal.png",
    "CHAR-OUTIS-01":      "Outis/Identities/LCB Sinner/11101_normal.png",
    "CHAR-DONQ-01":       "Don Quixote/Identities/LCB Sinner/10301_normal.png",
    "CHAR-MEURSAULT-01":  "Meursault/Identities/LCB Sinner/10501_normal.png",
    "CHAR-RYOSHU-01":     "Ryoshu/Identities/LCB Sinner/10401_normal.png",
    "CHAR-HONGLU-01":     "Hong Lu/Identities/LCB Sinner/10601_normal.png",
    "CHAR-SINCLAIR-01":   "Sinclair/Identities/LCB Sinner/11001_normal.png",
}

LOCATIONS = {
    "LOC-FOREST-01": "Dark_Forest.png",
    "LOC-BUSINT-01": "Mephitz_Inside.png",
}

# The story CGs are drawn as comic pages revealed panel by panel, so `S0_1_1`
# is panel one alone and `S0_1_5` is the finished page. Only the finished
# pages are useful as references; the partials are the same art with panels
# blanked out.
PROLOGUE_PAGES = {
    "PRO-PAGE1-01": "S0_1_5.png",
    "PRO-PAGE2-01": "S0_2_5.png",
    "PRO-PAGE3-01": "S0_3_4.png",
    "PRO-PAGE4-01": "S0_4.png",
    "PRO-PAGE5-01": "S0_5.png",
}


def load(path):
    img = Image.open(path)
    return img.convert("RGBA")


def alpha_bbox(img):
    a = np.array(img.getchannel("A"))
    ys, xs = np.where(a > ALPHA_MIN)
    if len(ys) == 0:
        raise ValueError("fully transparent image")
    return int(xs.min()), int(ys.min()), int(xs.max()) + 1, int(ys.max()) + 1


def tallest_component_bbox(img):
    """Bounding box of the tallest opaque blob -- the standing figure."""
    a = np.array(img.getchannel("A")) > ALPHA_MIN
    labels, n = ndimage.label(a)
    if n == 0:
        raise ValueError("no opaque pixels")
    slices = ndimage.find_objects(labels)
    min_h = img.height * MIN_FRAC
    best, best_h = None, -1.0
    for sl in slices:
        if sl is None:
            continue
        ys, xs = sl
        h = ys.stop - ys.start
        if h < min_h:
            continue
        if h > best_h:
            best, best_h = sl, h
    if best is None:
        raise ValueError("no component above the size floor")
    ys, xs = best
    return xs.start, ys.start, xs.stop, ys.stop


def crop(img, box):
    x0, y0, x1, y1 = box
    x0 = max(0, x0 - PAD)
    y0 = max(0, y0 - PAD)
    x1 = min(img.width, x1 + PAD)
    y1 = min(img.height, y1 + PAD)
    out = img.crop((x0, y0, x1, y1))
    scale = MAX_EDGE / max(out.width, out.height)
    if scale < 1.0:
        out = out.resize(
            (max(1, round(out.width * scale)), max(1, round(out.height * scale))),
            Image.LANCZOS,
        )
    return out


def emit(img, dest):
    os.makedirs(os.path.dirname(dest), exist_ok=True)
    img.save(dest, "PNG", optimize=True)
    return f"{img.width}x{img.height}"


def main():
    rows = []
    missing = []

    for key, (rel, mode) in CHARACTERS.items():
        src = os.path.join(CHARS, rel)
        if not os.path.exists(src):
            missing.append(src)
            continue
        img = load(src)
        box = alpha_bbox(img) if mode == "sprite" else tallest_component_bbox(img)
        dest = os.path.join(OUT, "characters", key + ".png")
        rows.append((key, mode, emit(crop(img, box), dest), rel))

    for key, rel in COSTUMES.items():
        src = os.path.join(CHARS, rel)
        if not os.path.exists(src):
            missing.append(src)
            continue
        img = load(src)
        dest = os.path.join(OUT, "costumes", key + "_LCB.png")
        rows.append((key + " (LCB)", "costume", emit(crop(img, alpha_bbox(img)), dest), rel))

    for key, name in LOCATIONS.items():
        src = os.path.join(PROLOGUE, name)
        if not os.path.exists(src):
            missing.append(src)
            continue
        dest = os.path.join(OUT, "locations", key + ".png")
        os.makedirs(os.path.dirname(dest), exist_ok=True)
        shutil.copyfile(src, dest)
        with Image.open(dest) as im:
            rows.append((key, "copy", f"{im.width}x{im.height}", name))

    for key, name in PROLOGUE_PAGES.items():
        src = os.path.join(PROLOGUE, name)
        if not os.path.exists(src):
            missing.append(src)
            continue
        dest = os.path.join(OUT, "prologue", key + ".png")
        os.makedirs(os.path.dirname(dest), exist_ok=True)
        shutil.copyfile(src, dest)
        with Image.open(dest) as im:
            rows.append((key, "copy", f"{im.width}x{im.height}", name))

    for key, mode, size, rel in rows:
        print(f"{key:<24} {mode:<8} {size:>11}  <- {rel}")
    print(f"\n{len(rows)} references written to {os.path.relpath(OUT, REPO)}")
    if missing:
        print("\nMISSING:", file=sys.stderr)
        for m in missing:
            print("  " + os.path.relpath(m, REPO), file=sys.stderr)
        return 1
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
