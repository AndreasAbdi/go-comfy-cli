"""Build per-shot Picture 1 references cropped to the shot's own framing.

MiniMax H3 R2V reproduces the framing and pose of its reference image far more
strongly than it follows the framing described in the prompt. Two SH030 attempts
proved it: handed full-figure standing references, it returned full-figure
standing characters side by side facing camera, both times, regardless of the
prompt asking for a low medium shot of Panther over a body -- and the second
attempt had only ONE character as Picture 1 and still composed two.

Rather than keep fighting that, use it. If the reference is cropped to a medium
shot, the model returns a medium shot. So each shot gets a Picture 1 cropped to
its own framing tier, taken from the character's full-figure art.

Crop windows are fractions of the figure's height measured from the top of the
alpha bounding box, so they track the actual figure rather than the canvas:

  ECU / INS   head and immediate detail
  CU          head and shoulders
  MED / OTS   waist up
  ACT / POV   three-quarter figure
  WIDE / EST  whole figure, untouched

Run:  tools/refs/.venv/Scripts/python.exe tools/produce/framing_crops.py
"""

import csv
import os
import re

import numpy as np
from PIL import Image

REPO = os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
PROD = os.path.join(REPO, "production", "limbus-prologue")
OUT = os.path.join(PROD, "references", "framing")

# framing -> (top, bottom) as a fraction of figure height
WINDOWS = {
    "ECU": (0.00, 0.20),
    "INS": (0.00, 0.24),
    "CU":  (0.00, 0.34),
    "MED": (0.00, 0.58),
    "OTS": (0.00, 0.58),
    "ACT": (0.00, 0.80),
    "POV": (0.00, 0.80),
}
MIN_EDGE = 768      # upscale small crops so the model gets real detail
ALPHA_MIN = 24
BG = (18, 18, 22)


def figure_box(img):
    a = np.array(img.getchannel("A")) > ALPHA_MIN
    ys, xs = np.where(a)
    return int(xs.min()), int(ys.min()), int(xs.max()) + 1, int(ys.max()) + 1


def head_top(img, x0, x1, y0, y1):
    """Row where the body actually starts, ignoring raised weapons.

    Panther's polearm and Wolf's shoulder blade rise above the head, so the
    alpha bounding box top is the weapon tip. Measuring the crop windows from
    there put Panther's face at the bottom of her own close-up. A weapon haft
    is a thin vertical line; a head is not. So scan down for the first row wide
    enough to be a head.
    """
    a = np.array(img.getchannel("A"))[y0:y1, x0:x1] > ALPHA_MIN
    if a.size == 0:
        return y0
    counts = a.sum(axis=1)
    need = max(4, int((x1 - x0) * 0.13))
    rows = np.where(counts >= need)[0]
    if rows.size == 0:
        return y0
    # require the width to persist, so a crossguard does not read as a head
    for r in rows:
        if counts[r:r + 12].min() >= need * 0.7:
            return y0 + int(r)
    return y0 + int(rows[0])


def crop_to(path, framing, dest):
    img = Image.open(os.path.join(REPO, path)).convert("RGBA")
    win = WINDOWS.get(framing)
    if win is None:
        out = img
    else:
        x0, y0, x1, y1 = figure_box(img)
        hy = head_top(img, x0, x1, y0, y1)
        h = y1 - hy                      # head-to-foot, not weapon-tip-to-foot
        top = hy + int(h * win[0])
        bot = hy + int(h * win[1])
        # Crop horizontally to the figure plus a small margin. Forcing a 16:9
        # window instead makes the width a multiple of the crop HEIGHT, which
        # strands a narrow figure in a sea of empty background on medium shots
        # and buries the face. The reference is a framing cue, not a plate; its
        # aspect does not have to match the output.
        mx = int((x1 - x0) * 0.15)
        box = (max(0, x0 - mx), max(0, top),
               min(img.width, x1 + mx), min(img.height, bot))
        out = img.crop(box)

    if max(out.size) < MIN_EDGE:
        s = MIN_EDGE / max(out.size)
        out = out.resize((max(1, round(out.width * s)), max(1, round(out.height * s))),
                         Image.LANCZOS)
    flat = Image.new("RGB", out.size, BG)
    flat.paste(out, (0, 0), out)
    os.makedirs(os.path.dirname(dest), exist_ok=True)
    flat.save(dest, "PNG")
    return flat.size


# How much of the location plate a shot should see. Picture 1 was cropped to the
# shot's framing while Picture 2 stayed a full wide establishing landscape in
# every shot -- two references disagreeing about scale, and the wide one is the
# whole frame. The 21-shot storyboard came back as medium-to-full figures almost
# throughout, including two extreme close-ups. So zoom the plate to match.
PLATE_ZOOM = {
    "ECU": 0.22,
    "INS": 0.26,
    "CU":  0.38,
    "MED": 0.62,
    "OTS": 0.62,
    "ACT": 0.80,
    "POV": 0.80,
    "WIDE": 1.0,
    "EST": 1.0,
}


def crop_plate(path, framing, dest):
    """Centre-crop the location plate to the fraction the framing implies."""
    z = PLATE_ZOOM.get(framing, 1.0)
    img = Image.open(os.path.join(REPO, path)).convert("RGB")
    if z >= 1.0:
        out = img
    else:
        w, h = img.size
        cw, ch = int(w * z), int(h * z)
        # bias slightly below centre: the action sits on the forest floor,
        # not in the canopy
        cx, cy = w // 2, int(h * 0.58)
        x0 = max(0, min(w - cw, cx - cw // 2))
        y0 = max(0, min(h - ch, cy - ch // 2))
        out = img.crop((x0, y0, x0 + cw, y0 + ch))
    if max(out.size) < MIN_EDGE:
        s = MIN_EDGE / max(out.size)
        out = out.resize((round(out.width * s), round(out.height * s)), Image.LANCZOS)
    os.makedirs(os.path.dirname(dest), exist_ok=True)
    out.save(dest, "PNG")
    return out.size


def main():
    rows = list(csv.DictReader(
        open(os.path.join(PROD, "shots", "manifest.csv"), encoding="utf-8")))
    made = 0
    for r in rows:
        sid = r["shot_id"]
        p = os.path.join(PROD, "shots", sid, f"{sid}_MINIMAX_v001.prompt.md")
        fm = re.match(r"^---\n(.*?)\n---\n", open(p, encoding="utf-8").read(), re.S).group(1)
        # Read picture_1_SOURCE, not picture_1_subject: once a crop exists the
        # builder points `subject` at the crop, so cropping from `subject`
        # would re-crop this script's own output every run.
        m = re.search(r"^picture_1_source:\s*(\S+)", fm, re.M)
        pic1 = m.group(1) if m else re.search(
            r"^picture_1_subject:\s*(\S+)", fm, re.M).group(1)
        framing = r["framing"].strip()

        # Location plate, zoomed to the same framing tier as the subject.
        plate = re.search(r"^picture_2_location:\s*(\S+)", fm, re.M)
        if plate:
            crop_plate(plate.group(1), framing,
                       os.path.join(OUT, f"{sid}_PLATE.png"))

        # Sheets and the style anchor are already the right thing; do not crop.
        # Neither are the composed detail references (CHAR-DANTE-HEAD-01,
        # -HEADSIDE-, -FALLEN-): they are already framed for their purpose, and
        # they are flattened RGB with no alpha for figure_box to measure.
        if any(t in pic1 for t in ("_SHEET", "PRO-PAGE", "-HEAD-",
                                   "-HEADSIDE-", "-FALLEN-")):
            continue
        dest = os.path.join(OUT, f"{sid}_PIC1.png")
        size = crop_to(pic1, framing, dest)
        print(f"{sid}  {framing:<4} {size[0]:>5}x{size[1]:<5} <- {os.path.basename(pic1)}")
        made += 1
    print(f"\n{made} framing-matched references written to "
          f"{os.path.relpath(OUT, REPO)}")


if __name__ == "__main__":
    main()
