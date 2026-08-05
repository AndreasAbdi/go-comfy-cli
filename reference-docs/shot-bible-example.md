Yes. Here is a practical **v0.1 shot bible** designed for converting manga panels or generated keyframes into a coherent animated sequence.

# SHOT DIRECTION BIBLE

**Project:** [Project Name]
**Version:** 0.1
**Target format:** 16:9 animation
**Master frame rate:** 24 fps
**Primary use:** Manga-to-storyboard, storyboard-to-video, and generated-shot assembly

---

## 1. Creative Objective

The adaptation should preserve the source material’s narrative clarity, character emotion, and strongest compositions while translating static panels into intentional screen time.

The animation should not introduce movement merely to prevent stillness. Camera movement, cuts, character motion, and effects must serve one of the following purposes:

* Reveal information
* Clarify physical action
* Emphasize emotion
* Redirect attention
* Establish space
* Control tension
* Support dialogue or sound

When no movement improves the scene, use a strong held composition.

---

## 2. Default Directing Style

The default style is **manga-faithful cinematic limited animation**.

Its characteristics are:

* Preserve iconic manga compositions when practical.
* Use clear spatial continuity before stylized cutting.
* Favor strong key poses over constant low-value motion.
* Reserve close-ups for meaningful emotional or narrative changes.
* Use camera movement sparingly and deliberately.
* Allow environmental details and sound to carry quiet scenes.
* Use rapid editing only when the story’s intensity warrants it.
* Avoid excessive dissolves, zooms, camera shake, and decorative motion.
* Hold important reactions long enough to register.
* Prioritize silhouette, eye direction, and pose readability.

---

## 3. Technical Standards

### 3.1 Picture

* Resolution: 1920 × 1080
* Aspect ratio: 16:9
* Frame rate: 24 fps
* Delivery frame rate: constant
* Working color space: Rec. 709
* Pixel aspect ratio: square
* Default safe area: keep essential faces, text, and actions within the central 90% of frame

### 3.2 Animation Cadence

* Playback remains 24 fps.
* Ordinary limited animation may be animated on twos: one drawing every two frames.
* Slow or held movement may use threes.
* Fast action, camera movement, effects, and precise lip sync may use ones.
* Cadence changes should be intentional rather than accidental.

### 3.3 Shot Identification

Use the following naming structure:

`SEQ_SCENE_SHOT_VERSION`

Example:

`SQ01_SC03_SH020_v004`

Every shot record should include:

* Shot ID
* Source panel or source image
* Shot type
* Duration
* Start and end frame
* Dialogue
* Character action
* Camera action
* Transition
* Audio cue
* Continuity notes
* Generation or compositing notes

---

## 4. Shot Selection Rules

### 4.1 Establishing Shots

Use an establishing shot when:

* Entering a new location
* Spatial relationships are unclear
* A character’s position matters to the next action
* Time of day or environment affects the scene
* The audience needs a pause before an important event

Default duration:

* Simple location: 2–3 seconds
* Atmospheric location: 3–5 seconds
* Major reveal: 4–8 seconds

Do not automatically re-establish a location after every cut.

### 4.2 Wide Shots

Use wide shots for:

* Physical geography
* Multiple-character blocking
* Entrances and exits
* Full-body acting
* Action whose direction must remain clear
* Isolation or scale

A wide shot should communicate information that would be lost in closer framing.

### 4.3 Medium Shots

Use medium shots as the default conversational and behavioral framing.

They are preferred for:

* Dialogue with visible body language
* Character interactions
* Transitional actions
* Reactions that involve posture
* Two-character scenes

Avoid cutting between medium shots without a change in information, emotion, or action.

### 4.4 Close-Ups

Use close-ups for:

* Emotional changes
* Recognition
* Decisions
* Important dialogue
* Revealing a significant object
* Concealing spatial information intentionally

Do not use a close-up merely because a manga panel contains a large face. Determine what the panel is emphasizing.

### 4.5 Extreme Close-Ups and Inserts

Use extreme close-ups for brief, high-value details:

* Eyes shifting
* A finger tightening
* A weapon being prepared
* A message being read
* A drop of sweat
* A critical object

Default duration:

* Rapid insert: 6–12 frames
* Readable detail: 12–36 frames
* Dramatic detail: 1.5–3 seconds

The viewer must have enough time to identify the object.

---

## 5. Manga Panel Interpretation

### 5.1 Large or Full-Page Panels

Default interpretation:

* Major reveal
* Emotional climax
* Environment introduction
* Powerful action pose
* Moment of suspended time

Possible screen treatment:

* Extended hold
* Slow push-in
* Layered parallax
* Environmental animation
* Delayed reaction
* Sound-led reveal

Do not assume the entire panel must remain visible at once. A large panel may become multiple shots.

### 5.2 Narrow Sequential Panels

Default interpretation:

* Accelerated pacing
* Incremental movement
* Repeated observation
* Rapid reactions
* Fragmented action

Possible screen treatment:

* Short cuts
* Match cuts
* Inserts
* A continuous movement divided into beats

### 5.3 Repeated Similar Panels

Default interpretation:

* Deliberate pause
* Slow realization
* Subtle movement
* Emotional hesitation
* Time passing

Do not collapse repeated panels automatically. Their repetition may be the point.

### 5.4 Borderless Panels

Default interpretation:

* Memory
* Atmosphere
* Emotional expansion
* Dreamlike perception
* A moment outside ordinary time

Use softened transitions or altered sound only when consistent with the narrative.

### 5.5 Speed Lines

Determine whether the lines represent:

* Subject movement
* Camera movement
* Psychological intensity
* Impact
* Abstract emphasis

Do not convert all speed lines into literal camera motion.

### 5.6 Speech Bubbles

Speech bubble placement does not define the screen composition.

Remove speech bubbles before layout and reconstruct the obscured artwork when necessary.

Dialogue duration should be derived from recorded or estimated speech, not panel size alone.

---

## 6. Timing Standards

Timing values are starting points and must be tested in an animatic.

### 6.1 General Shot Durations

* Flash or impact image: 2–6 frames
* Very quick insert: 6–12 frames
* Quick reaction: 12–24 frames
* Ordinary reaction: 1–2 seconds
* Silent emotional reaction: 2–5 seconds
* Ordinary non-dialogue shot: 1.5–3 seconds
* Dialogue shot: spoken duration plus 6–18 frames of breathing room
* Establishing shot: 2–5 seconds
* Dramatic reveal: 3–8 seconds

### 6.2 Dialogue Timing

For each line:

1. Use the actual voice recording when available.
2. Add a short lead-in before speech begins.
3. Leave time for the listener’s reaction when important.
4. Do not cut immediately after every final syllable.
5. Allow overlap when natural.
6. Cut away from the speaker when the listener’s reaction carries more meaning.

Default dialogue padding:

* Before line: 4–12 frames
* After line: 6–18 frames
* Significant emotional aftermath: 18–72 frames

### 6.3 Reaction Timing

A reaction should begin after the triggering information becomes understandable.

Reaction shots may include:

* Recognition
* Processing
* Physical response
* Decision
* Recovery

Do not show every reaction as an immediate facial change.

### 6.4 Comedic Timing

Comedy depends on setup, recognition, and payoff.

Typical structure:

1. Establish normal expectation.
2. Hold long enough for the expectation to register.
3. Deliver the disruption.
4. Use a sharp cut or controlled pause.
5. Show the consequence or reaction.

Avoid smoothing intentionally abrupt comedic cuts.

### 6.5 Action Timing

Action should alternate between:

* Anticipation
* Execution
* Impact
* Reaction
* Recovery

Fast action does not require every phase to be shown equally.

Important impacts may use:

* One or two anticipation poses
* A rapid action frame
* A brief impact frame
* A stronger aftermath hold
* Sound that begins before or extends beyond the visible impact

---

## 7. Camera Rules

### 7.1 General Principle

The camera moves only when movement improves storytelling.

Every camera move must have a stated purpose.

Acceptable purposes include:

* Following action
* Revealing information
* Shifting emotional alignment
* Increasing or decreasing pressure
* Connecting subjects spatially
* Creating scale
* Simulating a character’s attention

### 7.2 Static Camera

Static framing is preferred when:

* Character acting carries the scene
* The composition is already strong
* Stillness supports tension
* Dialogue requires clarity
* Additional movement would distract

A static camera does not require a completely static image. Hair, breathing, lighting, atmosphere, or background activity may provide subtle life.

### 7.3 Push-In

Use a push-in for:

* Realization
* Escalating tension
* Emotional isolation
* Important information
* A narrowing of attention

Default scale change:

* Subtle: 100% to 105%
* Moderate: 100% to 112%
* Strong: 100% to 125%

Avoid repeating push-ins throughout the same conversation.

### 7.4 Pull-Back

Use a pull-back for:

* Revealing context
* Emotional distance
* Isolation
* Consequences
* Transition from personal to environmental scale

### 7.5 Pan and Tilt

Use a pan or tilt when:

* Following movement
* Revealing connected information
* Exploring a large source panel
* Preserving continuity across a location

Avoid panning across a still image solely because the image is wider than the screen.

### 7.6 Handheld Motion and Camera Shake

Use only for:

* Physical instability
* Subjective panic
* Extreme impact
* Documentary immediacy
* Environmental force

Camera shake should be brief, readable, and tied to a specific cause.

Do not apply continuous random shake.

---

## 8. Composition Rules

### 8.1 Focal Priority

Every shot should have one primary focal point.

Secondary information must not compete with:

* The speaking character
* The active object
* The emotional reaction
* The narrative reveal

Use framing, contrast, depth, motion, and eye direction to guide attention.

### 8.2 Headroom and Lead Room

Maintain appropriate space around the character unless discomfort is intentional.

Provide lead room in the direction of:

* Gaze
* Movement
* Anticipated entry
* Threat

Breaking this rule should create deliberate pressure.

### 8.3 Eyelines

Characters looking at one another must maintain compatible eyelines.

When generating shots:

* Record the target of each gaze.
* Preserve approximate eye height.
* Avoid unexplained eye-direction changes.
* Check eyelines across cuts.

### 8.4 Screen Direction

Characters and moving objects should maintain consistent direction across cuts.

If a character moves left-to-right, preserve that direction until:

* The axis is visibly crossed
* A neutral shot resets orientation
* A new establishing shot clarifies geography
* Disorientation is intentional

### 8.5 The 180-Degree Rule

Maintain the interaction axis during ordinary dialogue and action.

Crossing the axis is permitted when:

* The camera movement visibly crosses it
* The scene deliberately destabilizes
* A neutral frontal or overhead angle resets orientation
* A new establishing shot redefines space

---

## 9. Dialogue Coverage

Default dialogue coverage should not mechanically alternate between speaker close-ups.

Preferred options include:

* Two-shot with visible relationship
* Medium speaker shot
* Listener reaction
* Over-the-shoulder view
* Profile composition
* Wide shot held through several lines
* Detail insert connected to the dialogue
* Environmental cutaway when thematically relevant

Cut when:

* The emotional power shifts
* New information lands
* A character reacts
* Someone changes intention
* Physical behavior becomes important
* Silence changes the scene

Do not cut merely because the speaker changes.

---

## 10. Action Coverage

Before rapid action, establish:

* Character positions
* Threat direction
* Relevant environment
* Immediate objective
* Important objects

During action:

* Preserve readable silhouettes.
* Show direction before speed.
* Use close-ups for selected impacts, not every motion.
* Avoid changing camera angles so frequently that motion becomes unclear.
* Let sound bridge omitted actions.
* Use reaction shots to communicate force and consequence.

After action:

* Show the result.
* Allow enough time to understand changed positions or damage.
* Re-establish geography when necessary.

---

## 11. Transitions

### 11.1 Hard Cut

Use as the default transition.

A hard cut is appropriate for:

* Continuous action
* Dialogue
* Reactions
* Location changes with clear context
* Strong visual contrasts

### 11.2 Dissolve

Use for:

* Time passage
* Memory
* Emotional transition
* Related images separated by time or place

Avoid using dissolves to conceal weak cuts.

### 11.3 Fade

Use for:

* Major scene endings
* Significant time gaps
* Chapter or act boundaries
* Consciousness fading
* Deliberate closure

### 11.4 Match Cut

Use when two shots share:

* Shape
* Movement
* Composition
* Color
* Concept
* Sound

Match cuts should communicate an intentional relationship.

---

## 12. Limited Animation Rules

When movement resources are constrained, prioritize:

1. Eyes
2. Mouth
3. Head angle
4. Hands
5. Silhouette
6. Hair or clothing accents
7. Environmental motion
8. Camera movement

A held pose may be enhanced with:

* Blinking
* Breathing
* Small eye shifts
* Hair movement
* Cloth settling
* Light fluctuation
* Smoke, dust, rain, or particles
* Foreground parallax
* Controlled camera drift

Do not animate all secondary elements simultaneously.

---

## 13. Generated-Frame Rules

### 13.1 Reference Consistency

Each generated shot should use approved references for:

* Character identity
* Costume
* Props
* Location
* Lighting
* Time of day
* Art style
* Camera angle
* Emotional state

### 13.2 Generation Boundaries

Generate shots independently when they involve:

* A major camera-angle change
* A large pose change
* Significant occlusion
* New characters entering
* Rapid action
* A location transition

Do not rely on interpolation to bridge unrelated compositions.

### 13.3 Temporal Review

Flag a shot when any of the following changes unintentionally:

* Facial structure
* Eye color or shape
* Costume details
* Number of fingers
* Prop shape
* Background geometry
* Light direction
* Shadow position
* Line weight
* Character height
* Camera perspective

Minor texture variation may be corrected. Structural variation should usually be regenerated or repainted.

---

## 14. Interpolation Rules

Interpolation is permitted when:

* Camera perspective remains stable.
* Character identity remains consistent.
* Motion is modest and continuous.
* Occlusion is limited.
* Beginning and ending frames represent adjacent physical states.

Interpolation should not be trusted when:

* Hands cross faces.
* Limbs overlap heavily.
* The camera angle changes.
* An object appears or disappears.
* The character turns substantially.
* Hair or clothing changes topology.
* The two images are symbolic beats rather than physical endpoints.

Inspect all interpolated action frame by frame.

---

## 15. Compositing Standards

Whenever practical, separate:

* Background
* Midground
* Character
* Foreground
* Shadows
* Effects
* Atmosphere
* Lighting
* Text

Each layer should maintain:

* Clean edges
* Correct perspective
* Compatible sharpness
* Consistent grain
* Matching color temperature
* Appropriate depth and blur
* Stable masks across frames

Avoid excessive depth of field. Important faces and actions must remain readable.

---

## 16. Sound Direction

Sound should be introduced during the animatic stage.

Minimum temporary soundtrack:

* Dialogue
* Room tone
* Environmental ambience
* Key footsteps
* Object interaction
* Major impacts
* Music where required

Sound may begin before the corresponding image or continue after the cut.

Use silence intentionally before:

* Revelations
* Impacts
* Emotional decisions
* Punchlines
* Scene endings

Do not fill every moment with music or effects.

---

## 17. Editing Review Questions

For every shot, ask:

* What new information does this shot provide?
* Why does the shot begin here?
* Why does it end here?
* What should the viewer notice first?
* Is the emotional beat readable?
* Is the action direction clear?
* Does the cut preserve geography?
* Is the duration motivated?
* Does the camera movement have a purpose?
* Would the scene improve if this shot were removed?
* Is sound carrying information the image does not need to show?

If a shot has no clear purpose, combine it, shorten it, or remove it.

---

## 18. Shot Manifest Template

```yaml
shot_id: SQ01_SC01_SH010
source:
  panel: chapter_01_page_03_panel_02
  reference_images:
    - refs/character_a_front.png
    - refs/location_room_wide.png

story:
  purpose: Establish the character entering an unfamiliar room
  beat: Cautious arrival
  emotion: Controlled unease

picture:
  shot_size: wide
  camera_angle: eye_level
  composition: character enters from frame left
  focal_point: character silhouette
  aspect_ratio: "16:9"

timing:
  fps: 24
  duration_seconds: 3.0
  start_frame: 1
  end_frame: 72

character_action:
  - opens door
  - pauses
  - looks toward frame right

camera:
  movement: static
  reason: preserve tension and spatial clarity

animation:
  cadence: twos
  priorities:
    - door movement
    - head turn
    - eye direction
  secondary_motion:
    - coat settling
    - curtain movement

audio:
  dialogue: null
  ambience: interior_roomtone.wav
  effects:
    - door_open.wav
    - soft_footstep.wav
  music: none

transition:
  in: hard_cut
  out: hard_cut

continuity:
  entry_direction: left_to_right
  gaze_target: frame_right
  light_direction: window_from_right
  props:
    - shoulder_bag
    - closed_umbrella

generation:
  separate_layers:
    - character
    - door
    - room_background
    - foreground_shadow
  interpolation_allowed: true
  risk_notes:
    - hand on door handle
    - bag strap consistency

review_status: pending
```

---

## 19. Review Statuses

Every shot should use one status:

* `planned`
* `storyboarded`
* `generated`
* `assembled`
* `needs_regeneration`
* `needs_cleanup`
* `needs_timing_revision`
* `approved`
* `final`

A shot must not be marked final until picture, timing, continuity, sound, and technical output have been reviewed.

---

## 20. Project-Specific Overrides

Any scene may override the defaults in this bible when the choice is intentional and documented.

An override should record:

* Default rule being changed
* Reason for the change
* Scene or shot range affected
* Approval status

The bible exists to create consistency, not to prevent expressive decisions.

The next useful step is to tailor this around your intended genre and visual language—for example, atmospheric drama, high-energy action, horror, comedy, or motion-comic adaptation.
