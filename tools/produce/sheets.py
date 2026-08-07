"""Composite multi-character reference sheets for keyframe generation.

qwen-image-edit has exactly two LoadImage nodes -- one source, one reference --
but shots in this cut cite up to four characters plus a location. The location
takes the source slot, so every character in the shot has to arrive through the
single reference slot. This lays them out side by side at a common height on a
neutral ground, which Qwen-Image-Edit-2511 reads as a character sheet.

Order matters: the sheet is built in the same order the prompt lists the
characters, so "the leftmost figure" in the prompt and in the sheet agree.

Run:  tools/refs/.venv/Scripts/python.exe tools/produce/sheets.py
"""

import csv
import os
import re

from PIL import Image

REPO = os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
PROD = os.path.join(REPO, "production", "limbus-prologue")
OUT = os.path.join(PROD, "references", "sheets")

HEIGHT = 1152        # tall enough that faces survive the model's own downscale
GAP = 48
MARGIN = 40
BG = (26, 28, 32)    # near-black, matches the scene key; avoids a white halo


def frontmatter(path):
    txt = open(path, encoding="utf-8").read()
    m = re.match(r"^---\n(.*?)\n---\n", txt, re.S)
    return m.group(1) if m else ""


FALLEN = ("production/limbus-prologue/references/images/characters/"
          "CHAR-DANTE-FALLEN-01.png")


def refs_of(shot_id):
    """(input_image, [character reference paths]) as the prompt declares them.

    Substitutes the fallen Dante where the shot has him on the ground. A group
    sheet feeds the wide shots, and a sheet built from the standing sprite put
    him upright in the establishing shot where the fixers are standing over his
    body.
    """
    p = os.path.join(PROD, "shots", shot_id, f"{shot_id}_KEY_v001.prompt.md")
    fm = frontmatter(p)
    src = re.search(r"^input_image:\s*(\S+)", fm, re.M)
    chars = re.findall(r"^\s+-\s+(\S+\.png)", fm, re.M)

    mm = os.path.join(PROD, "shots", shot_id, f"{shot_id}_MINIMAX_v001.prompt.md")
    if os.path.exists(mm) and "CHAR-DANTE-FALLEN-01.png" in open(mm, encoding="utf-8").read():
        chars = [FALLEN if c.endswith("CHAR-DANTE-01.png") else c for c in chars]
    return (src.group(1) if src else None), chars


def build_sheet(paths, dest):
    tiles = []
    for p in paths:
        im = Image.open(os.path.join(REPO, p)).convert("RGBA")
        w = max(1, round(im.width * HEIGHT / im.height))
        tiles.append(im.resize((w, HEIGHT), Image.LANCZOS))

    width = sum(t.width for t in tiles) + GAP * (len(tiles) - 1) + MARGIN * 2
    sheet = Image.new("RGB", (width, HEIGHT + MARGIN * 2), BG)
    x = MARGIN
    for t in tiles:
        sheet.paste(t, (x, MARGIN), t)
        x += t.width + GAP
    os.makedirs(os.path.dirname(dest), exist_ok=True)
    sheet.save(dest, "PNG")
    return sheet.size


def main():
    rows = list(csv.DictReader(
        open(os.path.join(PROD, "shots", "manifest.csv"), encoding="utf-8")))
    made, single, none = 0, 0, 0
    for row in rows:
        sid = row["shot_id"]
        _, chars = refs_of(sid)
        if len(chars) == 0:
            none += 1
        elif len(chars) == 1:
            single += 1                       # feed the character ref directly
        else:
            dest = os.path.join(OUT, f"{sid}_SHEET.png")
            size = build_sheet(chars, dest)
            names = [os.path.basename(c).replace(".png", "") for c in chars]
            print(f"{sid}  {size[0]:>5}x{size[1]}  {' + '.join(names)}")
            made += 1
    print(f"\n{made} sheets built, {single} single-character shots use their ref "
          f"directly, {none} shots have no character reference")


if __name__ == "__main__":
    main()
