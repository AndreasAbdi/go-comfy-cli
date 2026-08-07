"""Render shot clips with MiniMax H3 reference-to-video via go-comfy-cli.

Each shot is generated straight from two references -- the character (or the
composited multi-character sheet) and the location plate -- with no keyframe in
between. Shots carrying an original recording additionally drive on that audio.

Route comes from the prompt frontmatter, which the builder derives from whether
the shot has voiced dialogue:
  R2V-A  26 shots  minimax-r2v-2img-audio   two images + driving dialogue
  R2V    19 shots  minimax-r2v-2img         two images

Seeds derive from shot number and attempt, so a rerun reproduces and a re-roll
is an explicit `--attempt` bump.

Usage:
  clips.py                     # every shot without a clip yet
  clips.py --shots 030,040
  clips.py --attempt 2 --force
  clips.py --megapixels 0.6    # faster look-see pass
  clips.py --dry-run
"""

import argparse
import csv
import json
import os
import re
import shutil
import subprocess
import sys
import time

REPO = os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
PROD = os.path.join(REPO, "production", "limbus-prologue")
CLI = os.path.join(REPO, "go-comfy-cli.exe")
LOG = os.path.join(REPO, "progress.txt")

WORKFLOWS = {
    "R2V-A": "minimax-r2v-2img-audio",
    "R2V": "minimax-r2v-2img",
}

# MiniMaxH3ReferenceToVideo documents its trained range as ~124-362 frames, and
# 28 of the 45 shots are cut shorter than 124 -- some as short as 39. Asking for
# an out-of-range length is asking the model to work where it was never trained.
# So always GENERATE at least 124 frames and let assemble.py trim each clip back
# to its manifest length; it already conforms every clip to an exact frame count
# and keeps the head, which is the part the dialogue is positioned against.
MIN_FRAMES = 124
FPS = 24


def log(msg):
    line = f"[{time.strftime('%H:%M:%S')}] {msg}"
    print(line, flush=True)
    with open(LOG, "a", encoding="utf-8") as f:
        f.write(line + "\n")


def drive_audio(sid, wav, shot_start, shot_seconds, line_start, dest):
    """Build the audio the model should lip-sync to for THIS shot.

    Dialogue runs on a continuous timeline and deliberately overruns the cut in
    19 places: 9 of the 26 voiced shots are handed a recording longer than the
    clip. Feeding the whole recording makes the model fit the entire line into
    the shorter window, so mouths race. Instead give it exactly what is audible
    during the shot -- the line delayed to its in-shot offset, cut at the shot
    end. The full recording is still used at assembly, where it is free to
    continue across the cut.
    """
    offset = max(0.0, line_start - shot_start)
    ms = int(round(offset * 1000))
    run(["ffmpeg", "-y", "-v", "error", "-i", os.path.join(REPO, wav),
         "-af", f"adelay={ms}|{ms},atrim=0:{shot_seconds:.4f},asetpts=PTS-STARTPTS",
         "-c:a", "pcm_s16le", dest])
    return dest


def run(cmd):
    return subprocess.run(cmd, capture_output=True, text=True, encoding="utf-8",
                          errors="replace")


def spec(sid):
    path = os.path.join(PROD, "shots", sid, f"{sid}_MINIMAX_v001.prompt.md")
    txt = open(path, encoding="utf-8").read()
    m = re.match(r"^---\n(.*?)\n---\n", txt, re.S)
    fm, body = m.group(1), txt[m.end():].strip()

    def one(key):
        g = re.search(rf"^{key}:\s*(\S+)", fm, re.M)
        return g.group(1) if g else None

    auds = re.findall(r"^\s+-\s+(input/\S+\.wav)", fm, re.M)
    return dict(
        route=one("route"),
        subject=one("picture_1_subject"),
        location=one("picture_2_location"),
        seconds=float(one("duration_seconds")),
        audio=auds[0] if auds else None,
        body=body,
    )


def edl_index():
    p = os.path.join(PROD, "audio", "edl.json")
    e = json.load(open(p, encoding="utf-8"))
    shots = {s["shot_id"]: s for s in e["shots"]}
    first_vo = {}
    for l in e["lines"]:
        if l["delivery"] == "vo" and l["wav"] and l["shot_id"] not in first_vo:
            first_vo[l["shot_id"]] = l
    return shots, first_vo


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--shots")
    ap.add_argument("--attempt", type=int, default=1)
    ap.add_argument("--megapixels", type=float, default=1.0)
    ap.add_argument("--steps", type=int, default=20)
    ap.add_argument("--dry-run", action="store_true")
    ap.add_argument("--force", action="store_true")
    a = ap.parse_args()

    global SHOTS, FIRST_VO
    SHOTS, FIRST_VO = edl_index()
    rows = list(csv.DictReader(
        open(os.path.join(PROD, "shots", "manifest.csv"), encoding="utf-8")))
    if a.shots:
        want = {s.strip().zfill(3) for s in a.shots.split(",")}
        rows = [r for r in rows if r["shot_id"].split("SH")[1] in want]

    todo = []
    for r in rows:
        sid = r["shot_id"]
        dest = os.path.join(PROD, "shots", sid, "clips",
                            f"{sid}_CLIP_a{a.attempt}.mp4")
        if os.path.exists(dest) and not a.force:
            continue
        todo.append((sid, dest))

    log(f"CLIPS: {len(todo)} to render (attempt {a.attempt}, "
        f"{a.megapixels} MP, {a.steps} steps)")
    ok, fail = 0, []
    for i, (sid, dest) in enumerate(todo, 1):
        s = spec(sid)
        num = int(sid.split("SH")[1])
        gen_seconds = max(s["seconds"], MIN_FRAMES / FPS)
        seed = num * 1000 + a.attempt
        name = WORKFLOWS[s["route"]]
        wf = os.path.join(REPO, "workflows", name, f"{name}.json")
        args = os.path.join(REPO, "workflows", name, f"{name}.args.yaml")

        work = os.path.join(PROD, "shots", sid, "clips")
        os.makedirs(work, exist_ok=True)
        ptxt = os.path.join(work, f"{sid}_MINIMAX.submit.txt")
        with open(ptxt, "w", encoding="utf-8") as f:
            f.write(s["body"])

        # Clear the raw folder first. Otherwise a failed re-render leaves the
        # previous run's mp4 behind and the "newest file" pick silently
        # promotes stale output as if it were this attempt's.
        outdir = os.path.join(work, f"_raw_a{a.attempt}")
        shutil.rmtree(outdir, ignore_errors=True)
        cmd = [CLI, "run", "--workflow", wf, "--args-file", args,
               "--set", f"subject_image={os.path.join(REPO, s['subject'])}",
               "--set", f"location_image={os.path.join(REPO, s['location'])}",
               "--set", f"prompt={ptxt}",
               "--set", f"duration={gen_seconds:.4f}",
               "--set", "aspect=16:9 (Widescreen)",
               "--set", f"megapixels={a.megapixels}",
               "--set", f"seed={seed}",
               "--set", "seed_mode=fixed",
               "--set", f"steps={a.steps}",
               "--output-folder", outdir]
        if s["route"] == "R2V-A":
            shot = SHOTS[sid]
            line = FIRST_VO[sid]
            drv = os.path.join(work, f"{sid}_DRIVE.wav")
            drive_audio(sid, s["audio"], shot["start"], gen_seconds,
                        line["start"], drv)
            cmd[-2:-2] = ["--set", f"reference_audio={drv}"]

        trimmed = ("" if gen_seconds <= s["seconds"] + 1e-6
                   else f" (gen {gen_seconds:.2f}s -> trim to {s['seconds']:.2f}s)")
        label = (f"[{i}/{len(todo)}] {sid} {s['route']} {s['seconds']:.2f}s{trimmed} "
                 f"seed={seed} subj={os.path.basename(s['subject'])}"
                 + (f" aud={os.path.basename(s['audio'])}" if s["audio"] else ""))
        if a.dry_run:
            print(label)
            continue

        t0 = time.time()
        p = subprocess.run(cmd, capture_output=True, text=True, encoding="utf-8",
                           errors="replace")
        # SaveVideo's filename_prefix is "video/MiniMax_H3", and the CLI keeps
        # that subfolder, so the output is at _raw/video/*.mp4 rather than
        # directly in the output folder. Walk for it.
        vids = []
        for root, _, files in os.walk(outdir):
            vids += [os.path.join(root, f) for f in files if f.endswith(".mp4")]
        vids.sort(key=os.path.getmtime)
        if p.returncode != 0 or not vids:
            log(f"{label} FAILED rc={p.returncode}: {(p.stderr or p.stdout)[-400:]}")
            fail.append(sid)
            continue
        shutil.copyfile(vids[-1], dest)
        log(f"{label} ok {time.time() - t0:.0f}s -> {os.path.relpath(dest, REPO)}")
        ok += 1

    log(f"CLIPS done: {ok} ok, {len(fail)} failed {fail if fail else ''}")
    return 1 if fail else 0


if __name__ == "__main__":
    sys.exit(main())
