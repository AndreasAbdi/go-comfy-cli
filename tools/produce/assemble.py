"""Assemble the rendered clips, dialogue and subtitles into the finished cut.

Reads production/limbus-prologue/audio/edl.json, which carries per-LINE global
timings. Those cannot be derived from the shot rows: dialogue runs on one
continuous timeline and deliberately overruns the shot boundary in 19 places,
so a line's position is not its shot's start plus an offset.

Video: every clip is conformed to exactly its manifest frame count at 24 fps and
a uniform size, then concatenated. A clip that came back short is padded by
freezing its last frame rather than letting concat drift the timeline.

Audio: each original recording is delayed to its global start and mixed. Dante's
lines are silent by design -- he has no voice -- and contribute subtitles only.

Usage:
  assemble.py                       # full cut from attempt-1 clips
  assemble.py --attempt 2
  assemble.py --burn-subs           # hard-burn subtitles instead of soft
"""

import argparse
import json
import os
import subprocess
import sys

REPO = os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
PROD = os.path.join(REPO, "production", "limbus-prologue")
EDL = os.path.join(PROD, "audio", "edl.json")
OUT = os.path.join(PROD, "cut")
FPS = 24
# ResolutionSelector rounds to multiples of 32, so 1.0 MP at 16:9 comes back as
# 1376x768 -- aspect 1.792, not 1.778. Scaling that into a 1344x768 box with
# `pad` would letterbox every single shot with 9px black bars. Fill the target
# and centre-crop the 1% overshoot instead.
W, H = 1920, 1080


def run(cmd, **kw):
    p = subprocess.run(cmd, capture_output=True, text=True, encoding="utf-8",
                       errors="replace", **kw)
    if p.returncode != 0:
        raise SystemExit(f"ffmpeg failed:\n{' '.join(cmd[:6])}...\n{p.stderr[-1500:]}")
    return p


def probe_frames(path):
    """Decoded frame count. -count_frames, not -count_packets: packets can
    disagree with frames and the conform is checked against this number."""
    p = run(["ffprobe", "-v", "error", "-select_streams", "v:0",
             "-count_frames", "-show_entries", "stream=nb_read_frames",
             "-of", "csv=p=0", path])
    return int(p.stdout.strip() or 0)


def default_score():
    """Newest generated score bed under audio/bed, if one exists."""
    d = os.path.join(PROD, "audio", "bed")
    cands = []
    for root, _, files in os.walk(d):
        cands += [os.path.join(root, f) for f in files
                  if f.endswith((".flac", ".wav", ".mp3"))]
    return max(cands, key=os.path.getmtime) if cands else None


def srt_time(t):
    h, r = divmod(t, 3600)
    m, s = divmod(r, 60)
    return f"{int(h):02d}:{int(m):02d}:{s:06.3f}".replace(".", ",")


CPS = 17.0        # characters per second a viewer can comfortably read
MIN_SUB = 1.2
GAP = 0.08        # never let two subtitles touch


def write_srt(lines, dest):
    """Subtitles timed for reading, not just for speech length.

    A line's spoken duration is the floor, not the answer: an 89-character line
    delivered in 3 seconds is unreadable, so each subtitle is also given at
    least text-length/CPS. Ends are then clamped so consecutive subtitles never
    overlap -- dialogue here deliberately overruns shot boundaries, so lines sit
    close together and naive end times collide.
    """
    out = []
    for i, l in enumerate(lines):
        need = max(MIN_SUB, l["seconds"] + 0.3, len(l["text"]) / CPS)
        end = l["start"] + need
        if i + 1 < len(lines):
            end = min(end, lines[i + 1]["start"] - GAP)
        end = max(end, l["start"] + MIN_SUB * 0.6)   # never collapse to nothing
        out.append(f"{i + 1}\n{srt_time(l['start'])} --> {srt_time(end)}\n"
                   f"{l['text']}\n")
    with open(dest, "w", encoding="utf-8") as f:
        f.write("\n".join(out))
    return dest


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--attempt", type=int, default=1)
    ap.add_argument("--burn-subs", action="store_true")
    ap.add_argument("--out", default=None)
    ap.add_argument("--score", default=None, help="score bed; defaults to the generated one")
    a = ap.parse_args()

    edl = json.load(open(EDL, encoding="utf-8"))
    os.makedirs(OUT, exist_ok=True)
    work = os.path.join(OUT, f"_work_a{a.attempt}")
    os.makedirs(work, exist_ok=True)

    # ---- conform each clip to its exact frame count and size
    conformed, missing = [], []
    for s in edl["shots"]:
        sid = s["shot_id"]
        src = os.path.join(PROD, "shots", sid, "clips",
                           f"{sid}_CLIP_a{a.attempt}.mp4")
        if not os.path.exists(src):
            missing.append(sid)
            continue
        dst = os.path.join(work, f"{sid}.mp4")
        have = probe_frames(src)
        want = s["frames"]
        vf = (f"scale={W}:{H}:force_original_aspect_ratio=increase,"
              f"crop={W}:{H},setsar=1")
        if have < want:
            # freeze the tail rather than let the timeline drift
            vf += f",tpad=stop_mode=clone:stop={want - have}"
        # `trim=end_frame=N` does NOT give N frames on these clips -- MiniMax
        # writes a 1/12288 timebase and trim came back 3 frames long every time,
        # which would have drifted the picture seconds ahead of the dialogue
        # across 45 shots. select on the frame index and restamp PTS instead.
        vf += f",select='lt(n\\,{want})',setpts=N/{FPS}/TB"
        run(["ffmpeg", "-y", "-i", src, "-an", "-vf", vf,
             "-r", str(FPS), "-vsync", "cfr",
             "-c:v", "libx264", "-crf", "16", "-preset", "medium",
             "-pix_fmt", "yuv420p", dst])
        got = probe_frames(dst)
        if got != want:
            raise SystemExit(f"{sid}: conform produced {got} frames, wanted {want}")
        conformed.append((sid, dst, have, want))

    if missing:
        print(f"MISSING {len(missing)} clips: {', '.join(missing[:8])}"
              f"{' ...' if len(missing) > 8 else ''}")
        if not conformed:
            return 1
    for sid, _, have, want in conformed:
        if have != want:
            print(f"  {sid}: {have} frames -> conformed to {want}")

    # ---- concat video
    listf = os.path.join(work, "concat.txt")
    with open(listf, "w", encoding="utf-8") as f:
        for _, p, _, _ in conformed:
            f.write(f"file '{p.replace(os.sep, '/')}'\n")
    silent = os.path.join(work, "video.mp4")
    run(["ffmpeg", "-y", "-f", "concat", "-safe", "0", "-i", listf,
         "-c", "copy", silent])

    # ---- dialogue bed
    vo = [l for l in edl["lines"] if l["delivery"] == "vo" and l["wav"]]
    tts = [l for l in edl["lines"] if l["delivery"] == "tts"]
    for l in tts:
        cand = os.path.join(PROD, "audio", "tts", f"line_{l['line_id']:03d}.wav")
        if os.path.exists(cand):
            l["wav"] = os.path.relpath(cand, REPO)
            vo.append(l)
    vo.sort(key=lambda l: l["start"])

    dialogue = os.path.join(work, "dialogue.wav")
    if vo:
        cmd = ["ffmpeg", "-y"]
        for l in vo:
            cmd += ["-i", os.path.join(REPO, l["wav"])]
        parts, labels = [], []
        for i, l in enumerate(vo):
            ms = int(round(l["start"] * 1000))
            parts.append(f"[{i}:a]aresample=48000,aformat=sample_fmts=fltp:"
                         f"channel_layouts=stereo,adelay={ms}|{ms}[d{i}]")
            labels.append(f"[d{i}]")
        total = edl["total_seconds"]
        fc = (";".join(parts) + ";" + "".join(labels) +
              f"amix=inputs={len(vo)}:normalize=0:dropout_transition=0[mix];"
              f"[mix]apad,atrim=0:{total},alimiter=limit=0.95[a]")
        cmd += ["-filter_complex", fc, "-map", "[a]", "-c:a", "pcm_s16le", dialogue]
        run(cmd)
    else:
        run(["ffmpeg", "-y", "-f", "lavfi", "-t", str(edl["total_seconds"]),
             "-i", "anullsrc=r=48000:cl=stereo", "-c:a", "pcm_s16le", dialogue])

    # ---- score bed under the dialogue
    # The plan builds ambience and score once in the edit and discards all
    # generated clip audio, which is why every clip is conformed with -an.
    # ACE-Step came back at 0.0 dBFS true peak with 19.4 LU range, so it is
    # loudness-normalised and sidechain-ducked rather than mixed raw.
    score = a.score or default_score()
    mixed = os.path.join(work, "mix.wav")
    if score and os.path.exists(score):
        run(["ffmpeg", "-y", "-i", dialogue, "-i", score, "-filter_complex",
             f"[0:a]aresample=48000,loudnorm=I=-15:TP=-1.5:LRA=11[dx];"
             f"[1:a]aresample=48000,atrim=0:{edl['total_seconds']},"
             f"loudnorm=I=-30:TP=-6:LRA=7[sc];"
             f"[dx]asplit=2[d1][dkey];"
             f"[sc][dkey]sidechaincompress=threshold=0.05:ratio=6:attack=20:"
             f"release=600[scd];"
             f"[d1][scd]amix=inputs=2:normalize=0:duration=first,"
             f"loudnorm=I=-16:TP=-1.5:LRA=11,alimiter=limit=0.9[a]",
             "-map", "[a]", "-c:a", "pcm_s16le", mixed])
    else:
        print("no score bed found -- dialogue only")
        mixed = dialogue

    # ---- subtitles
    srt = write_srt(sorted(edl["lines"], key=lambda l: l["start"]),
                    os.path.join(OUT, "subtitles.srt"))

    # ---- mux
    dest = a.out or os.path.join(OUT, f"limbus-prologue-3min_a{a.attempt}.mp4")
    cmd = ["ffmpeg", "-y", "-i", silent, "-i", mixed]
    tail = ["-map", "0:v", "-map", "1:a", "-c:a", "aac", "-b:a", "192k", "-ar", "48000",
            "-shortest", dest]
    if a.burn_subs:
        esc = srt.replace("\\", "/").replace(":", "\\:")
        cmd += ["-vf", f"subtitles='{esc}':force_style="
                       "'FontName=Arial,FontSize=20,PrimaryColour=&H00FFFFFF,"
                       "OutlineColour=&H00000000,Outline=2,MarginV=40'",
                "-c:v", "libx264", "-crf", "16", "-preset", "medium"]
    else:
        # soft subtitle track, so the frame stays clean as the plan requires
        cmd += ["-i", srt, "-c:v", "copy", "-c:s", "mov_text"]
        tail = ["-map", "2:s"] + tail
    run(cmd + tail)

    n = len(conformed)
    print(f"\nassembled {n}/{len(edl['shots'])} shots -> {os.path.relpath(dest, REPO)}")
    print(f"  {edl['total_frames']} frames, {edl['total_seconds']:.2f}s, "
          f"{len(vo)} spoken lines, {len(edl['lines'])} subtitles")
    return 0


if __name__ == "__main__":
    sys.exit(main())
