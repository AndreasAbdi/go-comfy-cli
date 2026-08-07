"""Derive two-image R2V workflows from the shipped minimax-r2v-audio-image.

The stock workflow wires one reference image. MiniMaxH3ReferenceToVideo already
declares a second slot (`ref_images.ref_image_1`) with a null link, so a second
image costs one LoadImage node and one link -- no change to the node's dynamic
input group, which is the part that would risk breaking UI->API conversion.

Two references is exactly what a shot needs: the character (or the composited
sheet, for multi-character shots) plus the location plate. That is enough to
skip keyframes entirely and generate each shot straight from references, which
is both fewer generations and one less lossy hop than
reference -> qwen keyframe -> video.

Produces:
  minimax-r2v-2img-audio/  two images + driving dialogue audio  (24 VO shots)
  minimax-r2v-2img/        two images, no audio reference       (21 silent shots)

Run:  tools/refs/.venv/Scripts/python.exe tools/produce/derive_workflows.py
"""

import json
import os
import shutil

REPO = os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
WF = os.path.join(REPO, "workflows")
SRC = os.path.join(WF, "minimax-r2v-audio-image", "minimax-r2v-audio-image.json")

R2V = 136          # MiniMaxH3ReferenceToVideo
SLOT_IMG1 = 4      # ref_images.ref_image_1
SLOT_AUDIO0 = 7    # ref_audios.ref_audio_0
LOADIMAGE_0 = 149
LOADAUDIO = 148

COMMON_ALIASES = """
  shift_video:
    selector: '.nodes[] | select(.type == "MiniMaxH3SigmaShift") | .widgets_values[0]'
    type: number
    cardinality: one
  shift_audio:
    selector: '.nodes[] | select(.type == "MiniMaxH3SigmaShift") | .widgets_values[1]'
    type: number
    cardinality: one
  prompt:
    selector: '.nodes[] | select(.title == "Input Text (Prompt)") | .widgets_values[0]'
    type: string
    cardinality: one
  # Seconds. Quantised downstream to a legal MiniMax length by
  # max(5,round(a*24)) + (5 - (max(5,round(a*24)) % 17)) % 17
  duration:
    selector: '.nodes[] | select(.title == "Float (Duration)") | .widgets_values[0]'
    type: number
    cardinality: one
  # width/height on the R2V node are dead -- they are driven by these.
  aspect:
    selector: '.nodes[] | select(.id == 115) | .widgets_values[0]'
    type: string
    cardinality: one
  megapixels:
    selector: '.nodes[] | select(.id == 115) | .widgets_values[1]'
    type: number
    cardinality: one
  seed:
    selector: '.nodes[] | select(.id == 129) | .widgets_values[0]'
    type: number
    cardinality: one
  seed_mode:
    selector: '.nodes[] | select(.id == 129) | .widgets_values[1]'
    type: string
    cardinality: one
  steps:
    selector: '.nodes[] | select(.id == 124) | .widgets_values[1]'
    type: number
    cardinality: one
"""

IMG_ALIASES = """  subject_image:
    selector: '.nodes[] | select(.id == 149) | .widgets_values[0]'
    type: string
    cardinality: one
  location_image:
    selector: '.nodes[] | select(.id == 151) | .widgets_values[0]'
    type: string
    cardinality: one
"""

AUDIO_ALIAS = """  reference_audio:
    selector: '.nodes[] | select(.id == 148) | .widgets_values[0]'
    type: string
    cardinality: one
"""


def node(d, nid):
    return next(n for n in d["nodes"] if n["id"] == nid)


def build(with_audio):
    d = json.loads(open(SRC, encoding="utf-8").read())

    # -- second LoadImage, cloned from the first so its shape stays valid
    src_li = node(d, LOADIMAGE_0)
    new_id = 151
    new_link = d["last_link_id"] + 1
    li = json.loads(json.dumps(src_li))
    li["id"] = new_id
    li["pos"] = [src_li.get("pos", [0, 0])[0], src_li.get("pos", [0, 0])[1] + 360]
    li["widgets_values"] = ["ComfyUI_00046_.png", "image"]
    li["outputs"][0]["links"] = [new_link]
    if len(li["outputs"]) > 1:
        li["outputs"][1]["links"] = []
    d["nodes"].append(li)

    r2v = node(d, R2V)
    r2v["inputs"][SLOT_IMG1]["link"] = new_link
    d["links"].append([new_link, new_id, 0, R2V, SLOT_IMG1, "IMAGE"])
    d["last_node_id"] = max(d["last_node_id"], new_id)
    d["last_link_id"] = new_link

    # -- MiniMaxH3SigmaShift between the UNet and everything that consumes it.
    # The stock workflow omits this node, so the model samples unshifted. Its
    # documented defaults are shift_video=12.0 / shift_audio=3.0, i.e. MiniMax
    # expects a substantial shift, and running without one is a prime suspect
    # for the near-frozen output (motion scores 0.17-0.47 across three clips).
    ss_id = 152
    ss_link = d["last_link_id"] + 1
    consumers = [(L[3], L[4]) for L in d["links"] if L[1] == 127]
    d["links"] = [L for L in d["links"] if L[1] != 127]
    d["nodes"].append({
        "id": ss_id, "type": "MiniMaxH3SigmaShift", "pos": [200, 900],
        "size": [280, 82], "flags": {}, "order": 5, "mode": 0,
        # Widget-backed inputs must be declared here too, each with a `widget`
        # ref and a null link, or UI->API conversion drops the widget values
        # and ComfyUI rejects the prompt with "Required input is missing".
        "inputs": [{"localized_name": "model", "name": "model", "type": "MODEL",
                    "link": ss_link},
                   {"localized_name": "shift_video", "name": "shift_video",
                    "type": "FLOAT", "widget": {"name": "shift_video"},
                    "link": None},
                   {"localized_name": "shift_audio", "name": "shift_audio",
                    "type": "FLOAT", "widget": {"name": "shift_audio"},
                    "link": None}],
        "outputs": [{"localized_name": "MODEL", "name": "MODEL", "type": "MODEL",
                     "links": []}],
        "properties": {"Node name for S&R": "MiniMaxH3SigmaShift"},
        "widgets_values": [12.0, 3.0],
    })
    d["links"].append([ss_link, 127, 0, ss_id, 0, "MODEL"])
    nxt = ss_link
    for tgt, slot in consumers:
        nxt += 1
        d["links"].append([nxt, ss_id, 0, tgt, slot, "MODEL"])
        for n in d["nodes"]:
            if n["id"] == ss_id:
                n["outputs"][0]["links"].append(nxt)
            if n["id"] == tgt:
                n["inputs"][slot]["link"] = nxt
    d["last_node_id"] = max(d["last_node_id"], ss_id)
    d["last_link_id"] = nxt

    if not with_audio:
        # Detach the driving audio. audio_vae stays wired -- it is a decoder,
        # not a reference; the model still emits an audio track, which the
        # assembly stage discards for silent shots.
        lk = r2v["inputs"][SLOT_AUDIO0]["link"]
        r2v["inputs"][SLOT_AUDIO0]["link"] = None
        d["links"] = [L for L in d["links"] if L[0] != lk]
        d["nodes"] = [n for n in d["nodes"] if n["id"] != LOADAUDIO]

    return d


def emit(name, doc, args_body):
    out = os.path.join(WF, name)
    os.makedirs(out, exist_ok=True)
    with open(os.path.join(out, f"{name}.json"), "w", encoding="utf-8") as f:
        json.dump(doc, f, ensure_ascii=False, indent=2)
    with open(os.path.join(out, f"{name}.args.yaml"), "w", encoding="utf-8") as f:
        f.write("version: 1\naliases:\n" + args_body)
    print(f"wrote workflows/{name}/ ({len(doc['nodes'])} nodes)")


def main():
    emit("minimax-r2v-2img-audio", build(True),
         IMG_ALIASES + AUDIO_ALIAS + COMMON_ALIASES)
    emit("minimax-r2v-2img", build(False),
         IMG_ALIASES + COMMON_ALIASES)


if __name__ == "__main__":
    main()
