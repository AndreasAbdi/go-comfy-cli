# TV GIRL: EXIT THE SIGNAL — Montage Storyboard Plan

**Version:** 0.1  
**Status:** planned  
**Reference:** `reference-docs/shot-bible-example.md`  
**Source library:** `out/qwen-image-edit-outpainting-480p-corrected/`  
**Target:** 46-second, 16:9, 24 fps anime/illustration music-video montage  
**Working resolution:** 864 × 480 proxy; conform the final master to 1920 × 1080, Rec. 709

---

## 1. Creative Objective

Create a compact dream narrative rather than a reel of unrelated portraits.

**Logline:** A red-haired woman sits alone in a dark control room and channel-surfs through other people's moments—commutes, cafés, winter light, flowers, and water—until the images stop behaving like television and become an exit. When the signal returns to the red room, her chair is empty.

The different subjects and illustration styles are intentional. They are not presented as one character changing identity. Each image is a different “channel,” possible life, memory, or emotional frequency seen by the woman in the control room.

**Emotional arc:** controlled isolation → curiosity → sensory overload → release → absence.

**Primary visual motifs:**

- Eyes and faces as portals
- Circular match shapes: eyes, cups, flowers, ripples, chair backs
- Red as control or captivity
- Green and gold as waking life
- Blue as surrender and release
- TV interference only at the boundary between the control room and the other channels

---

## 2. Directing Rules

These rules tailor the shot bible to this montage:

1. Use hard cuts and shape matches as the default. Use only two dissolves: the water passage and the final empty-room release.
2. Reserve generated motion for six hero shots. The remaining images should use editorial motion, masks, parallax, light movement, and short holds.
3. Never animate a frame merely to avoid stillness. Strong portraits may remain held for 24–60 frames.
4. Do not imply that visually different women are the same literal person. Continuity belongs to color, gaze, and shape—not facial identity.
5. Preserve the source image's art style inside each shot. Do not normalize every frame into one aesthetic.
6. Use a maximum digital push of 112% on still images. Avoid continuous floating-camera motion.
7. Generated clips should contain one shot and one modest action. Do not ask interpolation to invent large pose or camera-angle changes.
8. Text is added during editing, never generated inside an image or video model.
9. Discard generated clip audio. Build one continuous soundtrack in the edit so the montage sounds like a single work.
10. Keep faces, hands, cups, flowers, and chair silhouettes inside the central 90% title/action safe area.

---

## 3. Source Asset Manifest

The aliases below are used throughout the storyboard. All paths are relative to the repository root.

| Alias | Corrected source | Visual role | Use |
|---|---|---|---|
| A01 | `out/qwen-image-edit-outpainting-480p-corrected/wallhaven-rr8pyq/ComfyUI_00053_.png` | red-haired woman in dark stair/control space | opening portal |
| A02 | `out/qwen-image-edit-outpainting-480p-corrected/__makima_chainsaw_man_and_1_more_drawn_by_k00s__sample-e2ce3cfc66a8beb08935c65118fe3b61/ComfyUI_00033_.png` | woman seated in saturated red room | control-room master |
| A03 | `out/qwen-image-edit-outpainting-480p-corrected/__makima_chainsaw_man_drawn_by_nyokki_dream666__sample-a75c25cafc79b56aef7524eac3a0c7f4/ComfyUI_00034_.png` | pale seated red-haired woman | drained-signal echo |
| A04 | `out/qwen-image-edit-outpainting-480p-corrected/images/ComfyUI_00051_.png` | graphic red-haired portrait with magenta/black background | title glitch texture |
| A05 | `out/qwen-image-edit-outpainting-480p-corrected/1017250634609314958/ComfyUI_00035_.png` | extreme close-up, flushed face | channel fragment |
| A06 | `out/qwen-image-edit-outpainting-480p-corrected/1117877938773911515/ComfyUI_00038_.png` | stern child portrait on white | channel fragment |
| A07 | `out/qwen-image-edit-outpainting-480p-corrected/1146940230209484233/ComfyUI_00039_.png` | butterfly and eye macro | recurring portal/match image |
| A08 | `out/qwen-image-edit-outpainting-480p-corrected/1150880879781597428/ComfyUI_00040_.png` | woman riding a green-lit bus | first sustained outside channel |
| A09 | `out/qwen-image-edit-outpainting-480p-corrected/665406913737236022/ComfyUI_00045_.png` | overhead café selfie with drinks | social-memory channel |
| A10 | `out/qwen-image-edit-outpainting-480p-corrected/843510205263987575/ComfyUI_00047_.png` | woman offering a drink on a street | invitation channel |
| A11 | `out/qwen-image-edit-outpainting-480p-corrected/31666003624972872/ComfyUI_00041_.png` | woman in pink against white | soft reaction beat |
| A12 | `out/qwen-image-edit-outpainting-480p-corrected/418764465373044331/ComfyUI_00042_.png` | winter close-up in knitted scarf | warmth channel |
| A13 | `out/qwen-image-edit-outpainting-480p-corrected/814518282648248352/ComfyUI_00046_.png` | blue eye behind snow/flowers | cold match insert |
| A14 | `out/qwen-image-edit-outpainting-480p-corrected/cool girlys are here/ComfyUI_00049_.png` | woman in low golden city light | dusk threshold |
| A15 | `out/qwen-image-edit-outpainting-480p-corrected/the girl/ComfyUI_00052_.png` | face framed by white flowers | natural-world portal |
| A16 | `out/qwen-image-edit-outpainting-480p-corrected/download/ComfyUI_00050_.png` | woman lying in a flower meadow | release hero image |
| A17 | `out/qwen-image-edit-outpainting-480p-corrected/591590101081951129/ComfyUI_00044_.png` | figure beneath green leaves | breath/counterpoint |
| A18 | `out/qwen-image-edit-outpainting-480p-corrected/1093319247064612531/ComfyUI_00036_.png` | woman floating in dark water and daisies | submerged dream hero |
| A19 | `out/qwen-image-edit-outpainting-480p-corrected/1116329826413159984/ComfyUI_00037_.png` | two pale figures in blue/coral abstract space | connection tableau |
| A20 | `out/qwen-image-edit-outpainting-480p-corrected/845550898835182253/ComfyUI_00048_.png` | upward-facing figure in bright blue water/light | emergence image |
| A21 | `out/qwen-image-edit-outpainting-480p-corrected/422494008813289154/ComfyUI_00043_.png` | crouched woman on white | utility silhouette; reserve |
| A22 | `out/qwen-image-edit-outpainting-480p-corrected/Makima/ComfyUI_00032_.png` | isolated red-haired profile on black | matte/glitch overlay; reserve |

**New required asset:**

- **N01 — empty red room:** create from A02 with Qwen Image Edit. Remove the seated woman completely while preserving the red wall, framed painting, black circular chair backs, perspective, shadows, and grain. Reconstruct the newly exposed chair upholstery and wall naturally. No new people, text, or props.

---

## 4. Overall Edit Structure

| Act | Time | Purpose | Palette | Editing character |
|---|---:|---|---|---|
| I — The Operator | 00:00–00:06.500 | Establish control room and the act of changing channels | black, red, magenta | stillness punctured by 2–4 frame signal tears |
| II — Other Lives | 00:06.500–00:22.000 | Let ordinary life become desirable and tactile | green, cream, pink, gold | longer readable shots and clean hard cuts |
| III — The Signal Opens | 00:22.000–00:39.000 | Images become nature, sensation, and release | green, white, deep blue, cyan | lyrical match cuts and selective generated motion |
| IV — Empty Chair | 00:39.000–00:46.000 | Return to the control room and reveal the escape | drained white → red → warm absence | deceleration, pull-back, final hold |

---

## 5. Master Shot Timeline

Frame numbers are 1-based and inclusive at 24 fps.

| Shot ID | Time / frames | Source | Picture and action | Camera / transition | Production route |
|---|---|---|---|---|---|
| `SQ01_SC01_SH010` | 00:00.000–00:01.500 / 1–36 | A01 | Dark close-up. The woman is nearly swallowed by the stairwell; her eyes catch a trace of red light. | Static for 24f, then a 12f push to 106%. Cut on a low CRT click. | **Hero video V01** or still hold if motion fails. |
| `SQ01_SC01_SH020` | 00:01.500–00:04.000 / 37–96 | A02 | Establish the red control room and centered seated figure. Let the symmetry feel coercive. | Static 36f, slow push 100→106%. Three horizontal tracking tears, each no longer than 3f. Hard cut. | Editorial still/parallax. |
| `SQ01_SC01_SH030` | 00:04.000–00:05.000 / 97–120 | A04 | Graphic interruption. Add the title `TV GIRL` for 16–18f; image tears sideways once. | No camera move. Add title in post. Smash cut. | Editorial still/title card. |
| `SQ01_SC01_SH040` | 00:05.000–00:06.500 / 121–156 | A05 → A06 → A07 | Three channel fragments: flushed face, stern child, butterfly eye. Each image provides a different emotional frequency. | 12f per image; cut on eye position and circular shapes. Final butterfly wing supplies the wipe edge. | Editorial rapid inserts. |
| `SQ01_SC02_SH050` | 00:06.500–00:10.000 / 157–240 | A08 | Bus passenger looks toward passing green light. Exterior foliage streams by; hair and shirt move slightly with the vehicle. | Locked medium shot with a very slow 100→104% push. Clean hard cut from butterfly to window reflection. | **Hero video V02.** |
| `SQ01_SC02_SH060` | 00:10.000–00:12.000 / 241–288 | A09 | Overhead café moment. Keep the raised hand, table, and drinks readable; allow a tiny shutter-flash lift near the end. | Reframe from 108→100% to reveal the table. Match cut from bus window oval to cup rim. | Editorial still with 2.5D separation. |
| `SQ01_SC02_SH070` | 00:12.000–00:14.500 / 289–348 | A10 | Street-side woman extends the drink toward camera. The cup crosses the foreground by only a few centimeters; background cars remain stable. | Static wide-angle feel. Cut as the cup nears lens. | Optional hero video; default editorial push. |
| `SQ01_SC02_SH080` | 00:14.500–00:16.000 / 349–384 | A11 | Soft pink portrait looks past camera. This is a calm held reaction, not a fashion pan. | Static 24f, then 12f drift left to create lead room. Hard cut. | Editorial still. |
| `SQ01_SC02_SH090` | 00:16.000–00:18.500 / 385–444 | A12 | Winter close-up. One blink, visible breath, and a small shift of knitted fabric. Preserve facial geometry. | Static close-up, no push. | Optional hero video; default layered still. |
| `SQ01_SC02_SH100` | 00:18.500–00:19.500 / 445–468 | A13 | Blue eye and white flakes/flowers. The eye is the focal point; particles cross at two depths. | 100→108% push. Cut on blink/black frame if available. | Editorial still with particles. |
| `SQ01_SC02_SH110` | 00:19.500–00:22.000 / 469–528 | A14 | Dusk portrait in amber shadow. Moving light, not camera motion, travels across the face. | Static. Let the last 8f fall nearly silent. | Editorial relight/masked grade. |
| `SQ01_SC03_SH120` | 00:22.000–00:25.000 / 529–600 | A15 | Flower-framed face. A breeze moves two or three petals, leaf shadows slide, and the subject blinks once near the end. | Slow pull-back 106→100% to reveal more flowers. Shape match from dusk eye highlight. | **Hero video V03.** |
| `SQ01_SC03_SH130` | 00:25.000–00:28.000 / 601–672 | A16 | Meadow release. Hair, shirt cuff, and flowers respond to one gentle gust; the hand opens toward the viewer. | Subtle rise/tilt toward the light, maximum 5% change. | **Hero video V04.** |
| `SQ01_SC03_SH140` | 00:28.000–00:29.500 / 673–708 | A17 | A second figure rests beneath leaves. This is a single breath between the meadow and water passages. | Static. Leaf shadows move; no character motion required. | Editorial still. |
| `SQ01_SC03_SH150` | 00:29.500–00:33.000 / 709–792 | A18 | Woman floats in dark blue water. Daisies drift, hair spreads slightly, one concentric ripple expands from her raised fingers. | Overhead static shot with a 100→104% push. Begin a 12f water-sound pre-lap. | **Hero video V05.** |
| `SQ01_SC03_SH160` | 00:33.000–00:35.500 / 793–852 | A19 | Two figures appear in a blue/coral suspended space. Preserve it as an iconic tableau; animate only drifting particulate and foreground parallax. | 108→100% pull-back. First and only soft 8f dissolve in. | Editorial still/composite. |
| `SQ01_SC03_SH170` | 00:35.500–00:38.000 / 853–912 | A20 | Upward-facing figure breaks into cyan light and spray. Treat the body as a silhouette emerging from water, not as a literal continuation of A19. | Gentle tilt upward; bright flare grows without clipping the face. Hard cut at peak cyan. | Optional hero video; editorial still is acceptable. |
| `SQ01_SC03_SH180` | 00:38.000–00:39.000 / 913–936 | A07 | Butterfly-eye image returns for exactly one second. The butterfly becomes a closing iris. | Fast push 100→112%. Cut through 2f black. | Editorial insert. |
| `SQ01_SC04_SH190` | 00:39.000–00:41.000 / 937–984 | A03 | Pale seated red-haired figure: the signal has been drained of red. Hold long enough for recognition. | Locked symmetrical frame. No glitch. Hard cut. | Editorial still. |
| `SQ01_SC04_SH200` | 00:41.000–00:43.000 / 985–1032 | A02 | Return to the saturated control room and seated figure. A final scan line crosses; the subject does not move. | Pull back 106→100%, reversing SH020. 4f white signal bloom at end. | Editorial still. |
| `SQ01_SC04_SH210` | 00:43.000–00:46.000 / 1033–1104 | N01 | Same red room, but the chair is empty. Warm natural light slowly replaces a narrow band of red. End title `EXIT THE SIGNAL` may appear for the final 36f. | Static 48f, then imperceptible pull-back for 24f. 12f dissolve from SH200. | **Hero video V06** from the edited N01 keyframe, or editorial relight if generation changes geometry. |

Total runtime: **46.000 seconds / 1,104 frames**.

---

## 6. Hero Video Generation Briefs

Each MiniMax prompt must be stored in its own Markdown file and follow `reference-docs/minimax-prompt-structure-guide.md`. Use all six sections in this exact order:

1. `subject_definitions`
2. `summary`
3. `retention_analysis`
4. `detailed_description`
5. `overall_soundscape`
6. `non_diegetic_music`

For image-driven shots, define the still as `<Picture 1>` and the visible person/environment as `<Subject 1>` or separate subjects when needed. Use `[keyframe completion]` in the summary. Each render should describe a **single continuous shot**. Do not ask MiniMax to cut between angles.

### V01 — Dark Signal Awakening (`SQ01_SC01_SH010`)

- **Reference:** A01 as `<Picture 1>` and opening keyframe.
- **Required motion:** faint red light flicker; one subtle eye lift or blink; minimal breathing.
- **Camera:** locked for most of the render; barely perceptible push during the final third.
- **Preserve:** red-haired identity, dark stair geometry, golden eyes, white shirt, heavy shadow.
- **Forbid:** walking, talking, hand gestures, new lights, anatomy changes, camera orbit.
- **Edit target:** select the strongest 1.5 seconds from the generated clip.

### V02 — The Bus Window (`SQ01_SC02_SH050`)

- **Reference:** A08 as `<Picture 1>` and opening keyframe.
- **Required motion:** green landscape and dappled light travel right-to-left outside; slight vehicle vibration; hair and loose fabric respond gently; one small eye shift toward the window.
- **Camera:** fixed medium framing, optional 4% push.
- **Preserve:** face, green shirt, seat, window frame, daylight direction, seated posture.
- **Forbid:** speaking, large head turn, new passengers, warped bus geometry, aggressive handheld shake.
- **Edit target:** 3.5 seconds.

### V03 — Flower Portal (`SQ01_SC03_SH120`)

- **Reference:** A15 as `<Picture 1>` and opening keyframe.
- **Required motion:** two or three flowers and strands of hair move in a soft breeze; leaf shadows drift across the face; one natural blink near the last third.
- **Camera:** slow 6% pull-back.
- **Preserve:** face shape, large eyes, white flowers, blue/green palette, upward gaze.
- **Forbid:** speaking, smile change, flying petals that cover the eyes, full-body invention.
- **Edit target:** 3 seconds.

### V04 — Meadow Release (`SQ01_SC03_SH130`)

- **Reference:** A16 as `<Picture 1>` and opening keyframe.
- **Required motion:** a single gentle gust travels through grass, flowers, long hair, and shirt; the visible hand relaxes and opens slightly; sunlight brightens by less than half a stop.
- **Camera:** tiny crane/tilt toward open sky, maintaining composition.
- **Preserve:** reclining pose, cream shirt, flower distribution, hand silhouette, warm late-afternoon light.
- **Forbid:** sitting up, reaching abruptly, speaking, new limbs, major camera rotation.
- **Edit target:** 3 seconds.

### V05 — Water Dream (`SQ01_SC03_SH150`)

- **Reference:** A18 as `<Picture 1>` and opening keyframe.
- **Required motion:** water ripples from raised fingertips; daisies drift at different speeds; hair floats subtly; chest movement is minimal.
- **Camera:** stable overhead composition with a 4% push.
- **Preserve:** closed eyes, pale dress, blue-black water, flower positions near the face, painterly surface texture.
- **Forbid:** opening eyes, submerging the face, large arm movement, added animals, camera roll.
- **Edit target:** 3.5 seconds.

### V06 — Empty Control Room (`SQ01_SC04_SH210`)

- **Reference:** N01 as `<Picture 1>` and opening keyframe.
- **Required motion:** dust in a narrow warm light beam; extremely slow light transition on the wall; no object movement except slight environmental flicker.
- **Camera:** locked, then a nearly invisible 2% pull-back.
- **Preserve:** empty chair, circular chair backs, red wall, painting, symmetry, room geometry.
- **Forbid:** any person, silhouette, face, text, new door, moving furniture, camera orbit.
- **Edit target:** 3 seconds.

**Fallback rule:** If a generated clip changes facial structure, hand anatomy, chair geometry, or illustration style, use the corrected still with masked environmental movement. Structural inconsistency is not worth extra motion.

---

## 7. Image Edit Plan for N01

Use `workflows/qwen-image-edit/` to create the empty-room ending.

**Input image:** A02  
**Reference image:** A02 again if the workflow requires a second image  
**Positive prompt:**

> Remove the seated red-haired woman completely and reconstruct the empty central black chair and the red wall behind her. Preserve the exact camera angle, symmetrical composition, circular black chair backs, framed painting, red palette, lighting direction, shadows, line work, texture, and 16:9 framing. The room must be empty. Do not add a person, silhouette, text, logo, doorway, or new object.

**Negative prompt:**

> person, woman, face, body, hands, silhouette, ghost, extra furniture, new doorway, text, logo, watermark, changed perspective, changed painting, warped chair, asymmetry

Generate at least four candidates. Approve only a candidate that preserves the chair and painting geometry closely enough to dissolve against A02 without a visible background jump.

---

## 8. Current Tool Mapping and Invocation

Run the CLI executable directly from PowerShell with `&`. Do not use `Start-Process`, because it can damage backslashes in `--set` values.

### Create N01 with Qwen Image Edit

Save the positive and negative prompts from Section 7 as Markdown or text files before running. If a Markdown path is supplied as a prompt value, the CLI reads the file and passes its contents to the workflow.

```powershell
$source = '.\out\qwen-image-edit-outpainting-480p-corrected\__makima_chainsaw_man_and_1_more_drawn_by_k00s__sample-e2ce3cfc66a8beb08935c65118fe3b61\ComfyUI_00033_.png'
& .\go-comfy-cli.exe run `
  --workflow .\workflows\qwen-image-edit\qwen-image-edit.json `
  --args-file .\workflows\qwen-image-edit\qwen-image-edit.args.yaml `
  --set "input_image=$source" `
  --set "reference_image=$source" `
  --set 'positive_prompt=.\prompts\tv-girl\n01-empty-room-positive.md' `
  --set 'negative_prompt=.\prompts\tv-girl\n01-empty-room-negative.md' `
  --output-folder .\out\tv-girl-montage\n01-empty-room
```

### Generate a hero clip with MiniMax Image to Video

Use the same command shape for V01–V06, changing the source, prompt, and output folder. The workflow's native result may be longer than the edit target; trim the approved range during assembly.

```powershell
$source = '.\out\qwen-image-edit-outpainting-480p-corrected\1150880879781597428\ComfyUI_00040_.png'
& .\go-comfy-cli.exe run `
  --workflow .\workflows\minimax-image-to-video\minimax-image-to-video.json `
  --args-file .\workflows\minimax-image-to-video\minimax-image-to-video.args.yaml `
  --set "input_image=$source" `
  --set 'positive_prompt=.\prompts\tv-girl\v02-bus-window.md' `
  --output-folder .\out\tv-girl-montage\v02-bus-window
```

### Workflow roles

| Need | Workflow | Notes |
|---|---|---|
| Remove the control-room subject for N01 | `workflows/qwen-image-edit/` | Generate several candidates; geometry continuity matters more than detail novelty. |
| Animate a corrected still | `workflows/minimax-image-to-video/` | Primary route for V01–V06. Use a six-section reference prompt even though the CLI alias is named `positive_prompt`. |
| Re-outpaint or repair framing | `workflows/qwen-image-edit-outpainting/` | Use only if a source edge fails after final 16:9 conform. |
| Create an entirely new insert | `workflows/anima-text-to-image/` | Not required for the master cut. Avoid adding imagery until the still animatic proves a gap. |
| Transfer motion from a selected reference clip | `workflows/minimax-reference-video/` | Optional later experiment. Its prompt must follow the six-section reference structure and use the correct `<Video N>` labels. |
| Drive a speaking shot from voice audio | `workflows/minimax-r2v-audio-image/` | Not used in this dialogue-free master plan. |

---

## 9. Transition and Compositing Plan

### Shape matches

- A07 butterfly/eye → A08 bus-window reflection
- A09 drink rim → A10 foreground cup
- A12 eye highlight → A13 blue eye
- A14 eye shadow → A15 flower-shadow pattern
- A15 flowers → A16 meadow flowers
- A16 open hand → A18 raised fingers/ripple
- A18 ripples → A19 circular blue/coral forms
- A20 cyan flare → A07 bright butterfly wing
- A02 occupied chair → N01 empty chair through a geometry-locked dissolve

### Layer treatment

For still-based shots, separate only what improves the beat:

- Subject mask
- Background
- One foreground occluder or particle layer
- Light/shadow pass
- Optional grain/scanline pass

Do not create deep artificial parallax on close faces. It will make ears, hair, and cheeks feel like cardboard planes.

### Signal effect limits

- Use signal tearing only in SH010–SH040, SH180, and SH200.
- Maximum tear duration: 3 frames.
- Maximum RGB offset: 3–5 pixels at proxy resolution.
- Keep scanlines subtle; they should be visible mainly in red-room blacks.
- The outside-world shots are clean, with film grain but no CRT overlay.

---

## 10. Sound and Music Direction

No dialogue.

**Score:** approximately 96 BPM, dream-pop/electronic, 4/4. Begin with filtered mono synth and CRT hum; open into stereo guitar/synth texture at 00:06.500; add soft drums by 00:12.000; remove percussion at 00:29.500; resolve to one sustained warm chord at 00:43.000.

**Sound cues:**

- 00:00.000 — low electrical room tone
- 00:01.500 — relay click entering the red room
- 00:04.000 — short tape-start chirp under title
- 00:05.000–00:06.500 — three tuned channel clicks
- 00:06.500 — bus ambience enters on the cut
- 00:10.000 — small camera shutter or glass clink
- 00:12.000 — ice/cup movement, kept quiet
- 00:16.000 — soft winter air and cloth texture
- 00:22.000 — breeze replaces city noise
- 00:29.000 — water pre-lap begins 12 frames before SH150
- 00:35.500 — low underwater swell
- 00:38.000 — single butterfly-wing flutter, stylized rather than literal
- 00:41.000 — CRT hum returns
- 00:43.000 — relay switches off; room tone gives way to distant morning birds

Generated clip audio should be muted. Avoid layering unrelated native audio from each generation.

---

## 11. Production Order

### Pass 1 — Animatic

1. Assemble all 21 shots from corrected stills at the exact durations above.
2. Add temporary score, hard cuts, and only the two specified dissolves.
3. Test the emotional arc before generating video.
4. Remove any shot that feels like a duplicate mood rather than new information.

### Pass 2 — Key Image

1. Generate N01 from A02 with Qwen Image Edit.
2. Confirm a geometry-locked A02→N01 dissolve.
3. If no candidate matches, create N01 by manual inpaint/composite rather than accepting a warped room.

### Pass 3 — Hero Motion

Generate in priority order:

1. V02 bus window
2. V05 water dream
3. V04 meadow release
4. V03 flower portal
5. V06 empty room
6. V01 dark awakening

Stop once the montage has enough motion. It is acceptable to keep V01 or V06 as still-based shots.

### Pass 4 — Editorial Motion and Compositing

Build restrained 2.5D and light passes for the remaining shots. Add particles only where listed. Check that no two consecutive shots use the same push direction.

### Pass 5 — Finish

1. Conform to 1920 × 1080, 24 fps, constant frame rate, Rec. 709.
2. Center-crop the 864 × 480 sources minimally to true 16:9 before upscale; protect faces near frame edges.
3. Normalize grain and sharpness without erasing differences in illustration style.
4. Mix the soundtrack as one continuous stereo piece.
5. Review every generated clip frame by frame for face, finger, background, and light-direction changes.

---

## 12. Approval Checklist

- [ ] The montage reads as a deliberate channel-surfing dream, not an asset showcase.
- [ ] The red-haired control-room figure is clearly the framing device.
- [ ] Different people are presented as different channels, not identity drift.
- [ ] Every generated prompt uses the required six-section MiniMax reference structure.
- [ ] All hero clips contain one modest action and one camera intention.
- [ ] A02 and N01 dissolve without background or chair geometry jumping.
- [ ] Glitches are confined to signal-boundary shots.
- [ ] Faces, hands, cups, and flowers remain structurally stable.
- [ ] Water and flower passages have enough screen time to feel like release.
- [ ] The final empty chair holds long enough to register before the end title.
- [ ] Final runtime is 46.000 seconds / 1,104 frames at 24 fps.
- [ ] Final deliverable is constant-frame-rate 1920 × 1080 Rec. 709.

---

## 13. Deferred Variants

These are optional after the primary cut works:

- **30-second cut:** remove SH070, SH080, SH090, SH140, and SH170; shorten SH020 and SH150.
- **Looping ending:** replace the final title with a 2-frame black gap followed by the first 12 frames of SH010.
- **More graphic cut:** use A22 as a white/red silhouette matte over SH030 and SH180.
- **More human cut:** insert A21 for 24–36 frames between SH110 and SH120 as a neutral white-space reset.

The 46-second version remains the master creative plan.
