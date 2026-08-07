# VO validation (local, no ComfyUI)

Validates the prologue voice streams and derives the authoritative WAV → dialogue-line mapping
using [faster-whisper](https://github.com/SYSTRAN/faster-whisper) in an isolated venv. Nothing
here talks to ComfyUI, so it runs whether or not the server is up.

## Setup

```powershell
uv venv tools\asr\.venv --python 3.12
uv pip install --python tools\asr\.venv\Scripts\python.exe faster-whisper
```

`.venv/` is gitignored. CUDA is used automatically when available (`ctranslate2` reports the
device); pass `--device cpu` to force CPU, which works but is slower.

## Use

```powershell
# transcribe + translate every S001B-*.wav (caches per file; re-runs are free)
tools\asr\.venv\Scripts\python.exe tools\asr\transcribe_vo.py --model large-v3

# derive the mapping from the cache (no transcription)
tools\asr\.venv\Scripts\python.exe tools\asr\align_vo.py
```

Useful flags: `--limit N` and `--model tiny` for a fast smoke test, `--force` to ignore the
cache, `--prefix S002B` for another scene, `--device cpu`.

## What each script does

**`transcribe_vo.py`** — two passes per recording: `transcribe` (native language, word
timestamps) and `translate` (English, used for matching). Records detected language and
confidence, the native transcript, the English translation, and true speech in/out points.
Caches to `out/s001b-asr-cache.json` after every file, so an interrupted run resumes.

**`align_vo.py`** — recovers the mapping. The recordings are an ordered *subsequence* of the
script: every WAV is one line, order is preserved, unvoiced lines are skipped. A fixed numeric
offset cannot express that — the real offset drifts from −1 to +4 across `S001B` — so this runs
a global Needleman-Wunsch alignment scoring each WAV's English translation against each
candidate line. Applies verified corrections from `overrides.json` afterward.

**`overrides.json`** — manual corrections, each carrying the evidence that justifies it. Needed
where Whisper's English pass drops a clause on a short utterance, or where two lines share
identical English and similarity cannot separate them. Prefer an override with stated evidence
over tuning the scoring until it happens to produce the right answer.

## Outputs (`out/`)

| File | Contents |
|---|---|
| `s001b-asr-cache.json` | Raw per-recording ASR results |
| `s001b-vo-mapping.csv` / `.json` | **Authoritative mapping.** WAV, line id, offset, speaker, similarity, speech boundaries, native + English text |
| `s001b-vo-manifest.csv` / `.json` | Per-WAV diagnostics including the fixed-offset scores that the alignment replaced |

`s001b-vo-mapping.json` is consumed directly by the shot-plan build scripts. Regenerate it
before rebuilding any plan.

## Findings for `S001B`

- **Korean**, 56 of 59 at mean confidence 0.974. The three outliers (`ja` ×2, `en` ×1) are
  misdetections on non-verbal grunts.
- Offsets: `wav 01–35 → −1`, `36 → 0`, `37–47 → +1`, `52–57 → +2`, `58–70 → +3`, `72–86 → +4`.
- **Dante has 0 voiced lines of 23** — he is the silent protagonist. The alignment never
  optimised for this, so it independently corroborates the mapping.
- Eight pairs score below 0.20 and are flagged for confirmation by ear. All are short or
  non-verbal (`Tch!` for `*Scoff*`, `N-n-n-n…` for `……`) and bracketed by confident neighbours.
