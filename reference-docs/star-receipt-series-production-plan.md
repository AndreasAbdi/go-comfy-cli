# Star Receipt After School

## Ten-Minute Anime Production Plan

**Version:** 0.1  
**Format:** ten connected episodes, approximately 60 seconds each  
**Total runtime:** approximately 10 minutes  
**Picture:** 16:9, generated at 1344 x 768 where appropriate, delivered at 1920 x 1080  
**Frame rate:** constant 24 fps  
**Primary language:** English  
**Production approach:** audio-first limited animation assembled from short generated shots  
**Verification policy:** every material output requires independent subagent verification

This plan supplements the general directing rules in
[shot-bible-example.md](./shot-bible-example.md). Every MiniMax reference prompt must use the
ordered full-reference structure in
[minimax-prompt-structure-guide.md](./minimax-prompt-structure-guide.md).

## 1. Creative Objective

Create an original short-form anime with the warmth of an after-school ensemble comedy,
surreal visualizations of social anxiety, fast romantic-comedy reaction timing, and
hand-painted lunar fantasy. The named inspirations are creative touchstones only. Production
prompts should describe the original traits below instead of requesting a copy of any existing
series, character, or studio design.

The complete story must work as one ten-minute film while each one-minute episode also has its
own setup, turn, and payoff.

### Logline

Four students discover a celestial shopping arcade behind their abandoned music room. In the
arcade, every broken promise becomes a physical receipt. They form a band to settle their
promise debts before the lunar festival takes their mysterious keyboardist away from Earth.

### Themes

- Friendship is maintained through small promises, not grand declarations.
- Performance anxiety becomes manageable when it is shared.
- A perfect persona is less valuable than an imperfect, recognizable self.
- Music turns private feelings into something other people can help carry.

### Tone

- Sincere and cozy at its baseline.
- Abruptly graphic and exaggerated during jokes or anxiety spirals.
- Quiet, spacious, and painterly during emotional or celestial moments.
- Never cynical about a character's vulnerability.

## 2. Series Shape

There is no conventional opening sequence inside the one-minute episodes. Use a reusable
one-to-two-second title sting after the cold open in Episodes 2-9. Episodes 1 and 10 omit it to
preserve momentum.

| Episode | Title | One-minute story |
| --- | --- | --- |
| 1 | The Promise Receipt | Yuna says, "Hey, come on. You promised to take me shopping today. Wake up." Her tote prints a glowing receipt, and the music-room wall opens into a celestial storefront where Sayo is waiting. |
| 2 | Four Girls, Three Instruments | Yuna, Kiri, Sayo, and Emi attempt one rehearsal. A sour chord animates a swarm of tiny unpaid receipts, forcing them to play together. |
| 3 | An Audience of One | Kiri freezes before a single child in the shopping street. Her anxiety becomes an infinite empty theater until the band deliberately joins her imperfect first note. |
| 4 | The Dress With No Price | A celestial stage outfit makes Yuna effortlessly charismatic but erases her recognizable mannerisms. She rejects perfection and returns to her tangerine cardigan. |
| 5 | Rain on Side B | The club records rain for its song. In a quiet clubroom, Emi admits she records ordinary days because she is afraid of forgetting them. |
| 6 | Midnight Shopping Trip | The group explores the arcade. A playful spree turns serious when Kiri discovers that memories, not money, are accepted at the register. |
| 7 | The Princess Who Cannot Stay | Sayo reveals that she is the arcade's runaway heir and must return after the lunar festival. The others promise to finish their song before moonrise. |
| 8 | A Song Nobody Can Buy | The arcade offers the band instant fame in exchange for the memory of how they met. Yuna nearly accepts, then tears the contract receipt in half. |
| 9 | Returned Without Receipt | The clubroom disappears and the friends are separated into private fantasy spaces. Each fulfills one small promise, causing the room to rebuild itself around their music. |
| 10 | Promise, Paid in Full | The band performs on the school roof as Earth and the arcade overlap. Their song frees Sayo from the contract. At dawn, Yuna repeats the shopping line to the whole group as an affectionate joke. |

### Default episode rhythm

The exact timing may change with the final dialogue, but each episode should begin from this
shape:

- `00:00-00:07` — visual hook or cold-open line
- `00:07-00:09` — title sting when used
- `00:09-00:25` — ordinary-world setup
- `00:25-00:43` — promise complication or surreal escalation
- `00:43-00:55` — decision and emotional/comedic payoff
- `00:55-01:00` — button, reveal, or match cut into the next episode

Episodes should contain approximately 7-12 edited shots. A generated shot should normally be
3-9 seconds, even though MiniMax can generate longer clips. Held frames, reaction inserts,
sound bridges, and editorial reframing are part of the intended limited-animation language.

## 3. Character Continuity Bible

No production shot should be generated until the four lead reference packs pass the design-lock
gate. Stable identifiers, colors, proportions, accessories, and left/right placement matter more
than ornamental detail.

### Yuna Arai — vocalist and emotional center

- **Silhouette:** compact, energetic posture; chin-length warm-brown hair with outward tips.
- **Anchors:** two gold crescent hair clips over her left temple; tangerine cardigan; navy skirt;
  teal canvas shopping tote; amber-brown eyes.
- **Behavior:** occupies space without noticing; points with her whole arm; cheeks inflate before
  she complains; recovers quickly after embarrassment.
- **Voice:** clone from `shopping-promise.wav`; youthful, clear, impatient without hostility.
- **Never change:** hair clips, cardigan hue, tote color, eye color, or apparent age.

### Kiri Sato — guitarist

- **Silhouette:** narrow shoulders, slightly folded posture, oversized sleeves.
- **Anchors:** ink-blue bob covering part of the right eye; charcoal hoodie beneath the school
  blazer; blue star bandage on left index finger; cream electric guitar.
- **Behavior:** protects her centerline; glances before turning her head; fingers communicate
  emotion before her face does.
- **Voice:** original Qwen-designed voice; soft lower mezzo, dry close-mic quality, careful starts,
  occasional compressed rush when panic wins.
- **Never change:** covered-eye side, bandaged finger, guitar color, or hoodie shape.

### Sayo Mikage — keyboardist and celestial heir

- **Silhouette:** tall, vertical posture with restrained gestures.
- **Anchors:** long lavender braid over the right shoulder; silver moon brooch; ivory blouse;
  navy high-waisted skirt; cool gray-violet eyes.
- **Behavior:** holds eye contact too long; moves economically; formal composure breaks into very
  small, unmistakable smiles.
- **Voice:** original Qwen-designed voice; calm young contralto, precise consonants, gentle
  breath, faint formality, warmth emerging under pressure.
- **Never change:** braid side, moon brooch, cool palette, or measured physical rhythm.

### Emi Tachibana — drummer and field recordist

- **Silhouette:** grounded athletic stance; short auburn hair; sleeves commonly rolled.
- **Anchors:** green windbreaker; red over-ear headphones around neck; black field recorder;
  paired red drumsticks.
- **Behavior:** reacts with stillness; delivers jokes without signaling them; taps rhythms against
  available surfaces.
- **Voice:** original Qwen-designed voice; medium-low young voice, textured but clean, relaxed
  tempo, understated amusement.
- **Never change:** headphone color, windbreaker green, recorder shape, or drumstick pair.

### Tab — receipt familiar

- **Design:** a moth folded from cream receipt paper, blue ink markings, gold-lit edges, and a
  barcode-like tail.
- **Movement:** stop-motion-like folding on threes; brief one-frame snaps for comedy.
- **Sound:** paper flicks, tiny thermal-printer chirps, no spoken language.

## 4. Character Reference Plan

### Required master assets per lead

Generate and approve the following before episode work:

1. Neutral full-body front view.
2. Full-body three-quarter view.
3. Full-body profile facing left.
4. Full-body profile facing right.
5. Face close-up in neutral light.
6. Expression sheet: neutral, joy, irritation, embarrassment, fear, sadness, resolve.
7. Hands and signature-prop sheet.
8. Standard school outfit sheet with exact color swatches.
9. Performance outfit sheet when applicable.
10. Seated instrument pose and standing conversational pose.

Create one height lineup containing all four characters. Create separate scale sheets for the
clubroom, school roof, shopping street, and celestial arcade.

### Image workflow

1. Use **Anima** to explore the original hero design and clean anime illustration language.
2. Select one approved hero image per character; it becomes the identity source of truth.
3. Use **Qwen Image Edit** to produce new poses, expressions, outfits, and turnarounds while
   explicitly preserving identity, hair geometry, face, palette, and permanent accessories.
4. Use **Qwen Outpainting** to convert approved images to 16:9 compositions without changing
   the subject.
5. Reject any image with a changed anchor, ambiguous hands around an instrument, inconsistent
   apparent age, or a mismatched braid/part/accessory side.

### Reference selection per video shot

Do not overload MiniMax with the entire cast library. For a one-character shot, use:

- one approved face close-up;
- one pose or full-body reference;
- one outfit/prop reference if the shot makes it important;
- one location keyframe when needed.

For a two-character shot, use one face and one body reference for each character, plus at most
one location or composition image. Prefer reaction singles over crowded speaking shots.

### Naming convention

Use stable asset names:

```text
CHAR_YUNA_FACE_NEUTRAL_v001.png
CHAR_KIRI_BODY_3Q_v003.png
CHAR_SAYO_EXPR_RESOLVE_v002.png
LOC_CLUBROOM_WIDE_DAY_v001.png
EP03_SC02_SH050_KEYFRAME_v004.png
```

Do not overwrite an approved reference. Increment the version and record why it changed.

## 5. Voice Production and Transfer Plan

All production TTS is a command-line-only path. ComfyUI TTS nodes, workflows, APIs, and the
ComfyUI server are prohibited for voice design, cloning, auditions, and final dialogue. Final
dialogue cloning runs through the pre-existing native launcher at
`third_party/qwen3-tts.cpp/bin/qwen3-tts.ps1`. When a new identity must be designed, the
checked-in official Qwen VoiceDesign package may be invoked by the command-line runner at
`production/star-receipt/runtime/qwen_voice_design_cli.py`; its approved anchor is then cloned
through the native launcher for every production line. ComfyUI remains available only for
non-TTS image, video, and music workflows.

### Voice-source hierarchy

1. **Yuna:** use the provided `shopping-promise.wav` and its exact transcript as the clone
   reference. This is the only real supplied voice identity.
2. **Kiri, Sayo, and Emi:** create original anchor voices with
   `Qwen3-TTS-12Hz-1.7B-VoiceDesign`.
3. Approve one clean 8-15 second anchor passage for each designed voice.
4. Feed each approved anchor and its exact transcript into
   `Qwen3-TTS-12Hz-1.7B-Base` Voice Clone.
5. Reuse that same clone reference for every line belonging to that character.
6. Custom Voice may be used for temporary auditions, but it is not the identity master unless
   it wins a deliberate voice-lock review.
7. Every retained voice manifest records the executable/runner, model, source WAV, exact source
   transcript, command-line provenance, output hash, and an explicit `comfyui_tts_used: false`.
8. An independent subagent verifies transcript, speaker identity, technical integrity, timing,
   and the non-ComfyUI provenance before a voice can be locked.

Voice cloning or transfer must only use recordings whose use is authorized for this project.

### Voice anchor directions

These descriptions are starting points for Qwen Voice Design, not final prompts:

- **Kiri:** young adult female, soft lower mezzo, intimate close-mic tone, lightly breathy but
  intelligible, careful consonant attacks, a restrained pace that can suddenly accelerate under
  social pressure; nervousness should tighten vowels rather than turn into caricature.
- **Sayo:** young adult female contralto, smooth and composed, low dynamic range, precise diction,
  subtle formality, stable breath support; emotional warmth appears as softened consonants and a
  small rise in pace.
- **Emi:** young adult female, medium-low pitch, dry and natural, relaxed breath, steady pace,
  faint textured edge, understated amusement; never announcer-like or sleepy.

### Dialogue rules

- Write dialogue before generating video.
- Keep most lines between 1.5 and 5 seconds.
- Give each vocal event one speaker. Avoid overlapping generated dialogue.
- Generate 2-4 audio candidates for emotionally important lines and approve audio before picture.
- Preserve exact transcripts alongside every reference WAV.
- Do not alter pitch to create another character; each lead receives a distinct designed voice.
- Directly reuse the original shopping recording in Episodes 1 and 10. Other Yuna dialogue uses
  a clone derived from that authorized sample.
- Target approximately 15-25 seconds of dialogue per one-minute episode. Let acting, ambience,
  music, and reaction shots carry the remaining time.

### Audio-first animation pipeline

1. Finalize a short line and performance direction.
2. Generate and approve the Qwen TTS WAV.
3. Normalize conservatively without removing breath or performance texture.
4. Generate the character keyframe with the mouth naturally closed.
5. Run MiniMax R2V Audio + Image using the final WAV and approved character image.
6. In the MiniMax prompt, identify the audio reference and require the exact signal, words,
   timing, emotion, and speaker relationship.
7. In editorial, replace MiniMax's dialogue layer with the untouched approved Qwen WAV. This
   guarantees textual and vocal continuity even if the generative video soundtrack changes it.
8. Add room tone, effects, and music on separate layers.

### Voice file structure

```text
production/audio/voices/yuna/reference/
production/audio/voices/yuna/episode-01/
production/audio/voices/kiri/reference/
production/audio/voices/kiri/episode-03/
production/audio/voices/sayo/reference/
production/audio/voices/emi/reference/
```

Each reference folder should contain the WAV, exact transcript, voice-design instruction,
generation settings, model name, seed, approval status, and a note establishing permission or
provenance.

## 6. Facial Animation Rules

Facial performance is generated from locked image identity plus final audio, then judged as
acting rather than mere lip movement.

### Speaking-shot construction

- Begin and end with the mouth closed unless speech crosses the cut.
- Add 4-10 frames of visual anticipation before the first phoneme.
- Preserve the character's permanent eye shape, face width, hairline, and accessories.
- Use blinks at thought boundaries, not at regular mechanical intervals.
- Let eyebrows, cheeks, head angle, shoulders, and hands carry emotion alongside the mouth.
- Keep ordinary speaking shots between medium close-up and medium framing.
- Use profile dialogue only after identity is stable in profile references.
- Avoid hands crossing the mouth during important lip-sync phrases.
- Do not ask a single generated shot to perform large body travel, complex camera movement, and
  precise dialogue simultaneously.

### Dialogue coverage

- Use a speaker shot for the emotional start of a line.
- Cut to the listener when the information lands.
- Use inserts of hands, instruments, receipts, or the environment to bridge imperfect lip sync.
- Reserve extreme close-ups for a decision, realization, or deliberately heightened joke.
- Hold reactions for 12-48 frames depending on emotional weight.

### Facial acceptance gate

Reject or regenerate when:

- phonemes visibly lead or lag the final WAV;
- the jaw deforms beyond the approved face structure;
- eye color, eye count, hair clips, brooches, or braid side changes;
- teeth flicker or the mouth remains active after speech;
- the apparent age changes;
- a neutral line receives an unrelated exaggerated emotion;
- the character looks into camera without story motivation.

## 7. Visual Language

### Baseline illustration language

- Clean original 2D anime character design with confident graphite-like line variation.
- Soft painted backgrounds with simplified geometry and selective detail.
- Strong readable silhouettes and clear eye direction.
- Limited animation on twos, occasional threes for holds and Tab's folded movement.
- Ones reserved for fast effects, camera motion, instrument attacks, and impact accents.
- No photorealism, faux-3D plastic skin, logos, subtitles, watermarks, or imitation of an existing
  franchise's exact designs.

### Color scripts

**Ordinary school world**

- Warm cream walls, dusty teal shadows, tangerine accents, muted afternoon gold.
- Low-to-medium contrast and soft window light.
- Comfortable negative space; backgrounds feel occupied but not busy.

**Social-anxiety space**

- Paper-white or charcoal voids, harsh single-color accents, warped scale, repeated silhouettes.
- Abrupt graphic cuts, split frames, held poses, handwritten abstract marks without readable text.
- Sound may drop to breath, cloth, fingers, and one distant environmental detail.

**Celestial arcade**

- Indigo-violet sky interiors, silver counters, coral signage shapes, gold receipt light.
- Watercolor bloom, paper texture, impossible depth, slow parallax, and restrained sparkling
  particles.
- Wide compositions and slow reveals replace rapid comedy cutting.

### Camera grammar

- Static medium shots are the conversational default.
- Push in only for recognition, pressure, or a narrowing decision.
- Use a wide shot when geography, isolation, or ensemble blocking matters.
- Preserve the 180-degree axis and screen direction across generated shots.
- No decorative orbiting cameras, continuous random shake, or repeated push-ins.
- A large celestial reveal may use one deliberate pull-back or vertical rise.

### Comedy grammar

1. Establish an ordinary expectation.
2. Hold until it reads.
3. Deliver the disruption with a hard cut or one-frame graphic accent.
4. Hold the consequence longer than the disruption.
5. Return to ordinary scale without explaining the joke.

### Musical grammar

- Instrument motion favors a few correct, readable gestures over continuous fake virtuosity.
- Establish hand position before fast inserts.
- Cut on attacks but hold on sustained emotional notes.
- Reuse approved guitar, keyboard, and drum performance references to avoid technique drift.
- Give the final rooftop performance the most fluid animation and broadest camera scale.

## 8. Location and Prop Locks

Create the following reusable masters before Episode 1 production:

- **Clubroom:** wide north-facing master, door wall, window wall, instrument corner, seated
  dialogue axis, afternoon and rainy-night lighting.
- **School roof:** stairwell entrance, fence line, skyline direction, performance layout, dawn and
  moonlit versions.
- **Shopping street:** one continuous geography map, bus-stop end, music-store end, covered
  arcade, rain version.
- **Celestial arcade:** storefront threshold, central aisle, memory register, stage, Sayo's gate.

Lock Yuna's tote, Kiri's guitar, Sayo's portable keyboard, Emi's recorder/headphones/drumsticks,
Tab, and the glowing promise receipt as separate prop sheets.

## 9. Shot and Prompt Construction

Every shot record must include:

- shot ID and target duration;
- episode/scene purpose;
- approved picture references;
- approved audio reference and exact transcript when present;
- starting pose and ending pose;
- camera and screen direction;
- continuity anchors;
- MiniMax relationship markers;
- sound, music, and editorial notes;
- generation seed/version and acceptance status.

For MiniMax reference generation, write all required sections in this order:

1. `subject_definitions`
2. `summary`
3. `retention_analysis`
4. `detailed_description`
5. `overall_soundscape`
6. `non_diegetic_music`

When a Qwen WAV supplies final dialogue, define it as an audio reference bound to the visible
speaker. State whether it is copied directly or only supplies a voice identity. Put exact spoken
words inside `<d>[English] ...</d>`. Use stable subject and speaker IDs through the whole prompt.

## 10. Production Workflow

### Phase A — lock the bible

1. Approve final cast names, premise, episode arc, and wardrobe.
2. Generate and lock character, prop, height, and location references.
3. Generate Qwen voice auditions and approve one identity per lead.
4. Create reusable clone anchors and exact transcripts.
5. Generate one neutral and one emotional test line for every lead.

### Phase B — technical proof

Produce three tests before the full series:

1. Yuna saying the supplied shopping line in the clubroom.
2. A two-character Yuna/Sayo exchange with a listener reaction cut.
3. A silent six-to-eight-second celestial reveal with Tab and a controlled camera move.

The proof passes only when identity, lip sync, audio replacement, color, and editing cadence all
survive a final 1920 x 1080 render.

### Phase C — Episode 1 pilot

1. Write the 60-second script and timed audio play-through.
2. Build an animatic using approved stills and final voice lines.
3. Generate shots only after animatic timing is locked.
4. Edit, restore the original Qwen dialogue WAVs, mix, and review.
5. Update the bible with any continuity rule learned from the pilot.

### Phase D — batch production

Produce Episodes 2-9 in two batches, reusing approved assets and motion references. Produce
Episode 10 last so its rooftop performance can incorporate everything learned earlier.

### Phase E — finishing

- Conform the ten episodes into one ten-minute master.
- Preserve each episode as a separate 60-second deliverable.
- Normalize dialogue consistently and check intelligibility on speakers and headphones.
- Apply a single final color transform and constant 24 fps cadence.
- Review hard cuts, audio bridges, eye lines, screen direction, continuity anchors, and title cards.

## 11. Workload Estimate

A ten-minute limited-animation film at this density should contain approximately 80-120 unique
edited shots. Some are reusable stills, inserts, or alternate crops; the rest require generated
motion. Plan for 2-4 candidates per important shot and more for dialogue close-ups or the finale.

The efficient order is therefore:

1. 4 locked character packs.
2. 4 locked location packs.
3. 4 locked voice identities.
4. 3 technical proof shots.
5. 1 finished pilot episode.
6. 8 batch-produced middle episodes.
7. 1 higher-effort finale.

## 12. Approval Gates

### Design lock

- Every lead is recognizable in front, three-quarter, profile, seated, and expressive views.
- Permanent anchors survive Qwen edits and MiniMax motion tests.
- The ensemble reads as one original visual world.

### Voice lock

- Each voice is distinguishable without seeing the character.
- Each clone remains recognizable across neutral, excited, quiet, and distressed lines.
- The approved reference WAV and exact transcript are archived and unchanged.

### Motion lock

- Lip sync survives restoration of the original Qwen WAV.
- Eye, face, hair, hands, and accessories remain stable.
- The movement supports the shot's emotional purpose.

### Episode lock

- The one-minute episode has a readable beginning, turn, and ending.
- It advances either the promise plot or a character relationship.
- It connects cleanly to the following episode and still works as a standalone short.

### Series lock

- All ten episodes total approximately ten minutes.
- The repeated shopping line changes meaning between Episodes 1 and 10.
- Character appearance and voices remain consistent through the final rooftop scene.
- The finale resolves Sayo's contract, the band's promise, and Yuna's original complaint.

## 13. Independent Subagent Verification

Independent verification is mandatory throughout production. The agent that creates or directs
an output must always deploy a separate subagent to verify that the result matches the approved
plan, references, prompt, and technical requirements. Self-review may happen first, but it does
not replace the independent pass.

### Outputs requiring verification

Deploy a verifier after every material production outcome, including:

- character, expression, costume, prop, height, and location reference batches;
- Qwen voice auditions, approved voice anchors, clone tests, and episode dialogue batches;
- MiniMax facial-animation and motion tests;
- generated shot batches and any important regenerated shot;
- episode animatics, picture locks, audio mixes, and final episode exports;
- the combined ten-minute master and its delivery files;
- changes to this production bible that alter identity, continuity, voice, visual language, or
  acceptance criteria.

Routine bookkeeping, filenames, and other non-creative mechanical operations may be checked in
their surrounding batch rather than receiving a separate verifier invocation.

### Verifier responsibilities

The verifier receives the target output together with the relevant approved references, shot
record, exact dialogue transcript, generation prompt, and acceptance gate. It must independently
check, as applicable:

1. Character identity, apparent age, proportions, permanent anchors, costume, props, and screen
   direction.
2. Voice identity, exact words, pronunciation, emotion, timing, clipping, noise, and provenance.
3. Lip sync, mouth closure, facial stability, gaze, hands, instrument interaction, and intended
   acting beat.
4. Location geography, color script, visual language, camera purpose, and shot continuity.
5. Compliance with the six-section MiniMax prompt structure and declared reference-retention
   relationships.
6. Resolution, aspect ratio, duration, frame rate, audio format, synchronization, and file
   integrity.
7. Whether the output fulfills its narrative purpose inside the episode rather than merely
   looking attractive in isolation.

### Required verifier report

The verifier returns one of three explicit outcomes:

- **PASS:** matches the approved intent and may advance.
- **PASS WITH NOTES:** usable without regeneration; notes must be carried into later shots or
  finishing.
- **FAIL:** identifies the mismatch, affected frames or timestamps, likely cause, and a concrete
  regeneration or editing instruction.

Reports should cite the inspected filenames and, for time-based media, exact timestamps. A vague
approval such as "looks good" is not sufficient.

### Failure and escalation loop

1. A failed output does not advance to the next production gate.
2. Apply the verifier's correction through prompt revision, reference replacement, regeneration,
   or editing.
3. Deploy a subagent again to verify the revised output; the creator's own recheck is still not a
   substitute.
4. After two failed revisions for the same cause, compare the output directly with the locked
   reference and simplify the shot before another attempt.
5. Any proposed change to a locked character or voice identity must be recorded in this bible and
   independently verified before downstream work resumes.

### Production record

Maintain a verification ledger with:

```text
output_id | version | verifier | result | inspected_references | findings | required_action | date
```

No reference, shot, episode, or master is considered approved unless its latest version has a
corresponding independent subagent verification record.
