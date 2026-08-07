"""Composite extracted references into one sheet so the crops can be eyeballed.

Usage: contact_sheet.py <subdir> <out.png> [row_height]
"""

import os
import sys

from PIL import Image, ImageDraw

REPO = os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
OUT = os.path.join(REPO, "production", "limbus-prologue", "references", "images")

BG = (30, 30, 34)
LABEL_H = 26


def main():
    sub = sys.argv[1]
    dest = sys.argv[2]
    h = int(sys.argv[3]) if len(sys.argv) > 3 else 420

    d = os.path.join(OUT, sub)
    names = sorted(f for f in os.listdir(d) if f.endswith(".png"))
    tiles = []
    for n in names:
        im = Image.open(os.path.join(d, n)).convert("RGBA")
        w = max(1, round(im.width * h / im.height))
        tiles.append((n, im.resize((w, h), Image.LANCZOS)))

    gap = 8
    total = sum(t.width for _, t in tiles) + gap * (len(tiles) + 1)
    sheet = Image.new("RGB", (total, h + LABEL_H + gap * 2), BG)
    draw = ImageDraw.Draw(sheet)
    x = gap
    for n, t in tiles:
        sheet.paste(t, (x, gap), t)
        draw.text((x + 2, gap + h + 4), n.replace(".png", ""), fill=(210, 210, 215))
        x += t.width + gap

    os.makedirs(os.path.dirname(dest), exist_ok=True)
    sheet.save(dest, "PNG")
    print(f"{dest}  {sheet.width}x{sheet.height}  ({len(tiles)} tiles)")


if __name__ == "__main__":
    main()
