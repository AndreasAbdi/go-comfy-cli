[CmdletBinding()]
param(
    [switch]$KeepIntermediates
)

$ErrorActionPreference = 'Stop'

$repo = Split-Path -Parent $PSScriptRoot
$ffmpeg = (Get-Command ffmpeg -ErrorAction Stop).Source
$ffprobe = (Get-Command ffprobe -ErrorAction Stop).Source
$outputRoot = Join-Path $repo 'out\tv-girl-movie'
$workRoot = Join-Path $outputRoot 'work'
$shotRoot = Join-Path $workRoot 'shots'
$assetRoot = Join-Path $outputRoot 'assets'
$visualMaster = Join-Path $workRoot 'tv-girl-picture-only.mp4'
$soundtrack = Join-Path $workRoot 'tv-girl-soundtrack.wav'
$movie = Join-Path $outputRoot 'tv-girl-exit-the-signal.mp4'
$reviewMovie = Join-Path $outputRoot 'tv-girl-exit-the-signal-review-540p.mp4'

New-Item -ItemType Directory -Force -Path $outputRoot, $workRoot, $shotRoot, $assetRoot | Out-Null

function Resolve-RepoPath {
    param([Parameter(Mandatory)][string]$RelativePath)

    $path = Join-Path $repo $RelativePath
    if (-not (Test-Path -LiteralPath $path)) {
        throw "Missing movie asset: $path"
    }
    return (Resolve-Path -LiteralPath $path).Path
}

$assets = @{
    A01 = Resolve-RepoPath 'out\qwen-image-edit-outpainting-480p-corrected\wallhaven-rr8pyq\ComfyUI_00053_.png'
    A02 = Resolve-RepoPath 'out\qwen-image-edit-outpainting-480p-corrected\__makima_chainsaw_man_and_1_more_drawn_by_k00s__sample-e2ce3cfc66a8beb08935c65118fe3b61\ComfyUI_00033_.png'
    A03 = Resolve-RepoPath 'out\qwen-image-edit-outpainting-480p-corrected\__makima_chainsaw_man_drawn_by_nyokki_dream666__sample-a75c25cafc79b56aef7524eac3a0c7f4\ComfyUI_00034_.png'
    A04 = Resolve-RepoPath 'out\qwen-image-edit-outpainting-480p-corrected\images\ComfyUI_00051_.png'
    A05 = Resolve-RepoPath 'out\qwen-image-edit-outpainting-480p-corrected\1017250634609314958\ComfyUI_00035_.png'
    A06 = Resolve-RepoPath 'out\qwen-image-edit-outpainting-480p-corrected\1117877938773911515\ComfyUI_00038_.png'
    A07 = Resolve-RepoPath 'out\qwen-image-edit-outpainting-480p-corrected\1146940230209484233\ComfyUI_00039_.png'
    A08 = Resolve-RepoPath 'out\qwen-image-edit-outpainting-480p-corrected\1150880879781597428\ComfyUI_00040_.png'
    A09 = Resolve-RepoPath 'out\qwen-image-edit-outpainting-480p-corrected\665406913737236022\ComfyUI_00045_.png'
    A10 = Resolve-RepoPath 'out\qwen-image-edit-outpainting-480p-corrected\843510205263987575\ComfyUI_00047_.png'
    A11 = Resolve-RepoPath 'out\qwen-image-edit-outpainting-480p-corrected\31666003624972872\ComfyUI_00041_.png'
    A12 = Resolve-RepoPath 'out\qwen-image-edit-outpainting-480p-corrected\418764465373044331\ComfyUI_00042_.png'
    A13 = Resolve-RepoPath 'out\qwen-image-edit-outpainting-480p-corrected\814518282648248352\ComfyUI_00046_.png'
    A14 = Resolve-RepoPath 'out\qwen-image-edit-outpainting-480p-corrected\cool girlys are here\ComfyUI_00049_.png'
    A15 = Resolve-RepoPath 'out\qwen-image-edit-outpainting-480p-corrected\the girl\ComfyUI_00052_.png'
    A16 = Resolve-RepoPath 'out\qwen-image-edit-outpainting-480p-corrected\download\ComfyUI_00050_.png'
    A17 = Resolve-RepoPath 'out\qwen-image-edit-outpainting-480p-corrected\591590101081951129\ComfyUI_00044_.png'
    A18 = Resolve-RepoPath 'out\qwen-image-edit-outpainting-480p-corrected\1093319247064612531\ComfyUI_00036_.png'
    A19 = Resolve-RepoPath 'out\qwen-image-edit-outpainting-480p-corrected\1116329826413159984\ComfyUI_00037_.png'
    A20 = Resolve-RepoPath 'out\qwen-image-edit-outpainting-480p-corrected\845550898835182253\ComfyUI_00048_.png'
    N01 = Resolve-RepoPath 'out\tv-girl-movie\assets\n01-empty-red-room.png'
}

# SH040 is represented as three twelve-frame editorial inserts. Together they
# retain the 36-frame duration specified by the storyboard.
$shots = @(
    [pscustomobject]@{ Id='010'; Asset='A01'; Frames=36; Motion='push6'; Grade='eq=contrast=1.08:saturation=0.88:brightness=-0.04'; Effect='control'; Title='' },
    [pscustomobject]@{ Id='020'; Asset='A02'; Frames=60; Motion='push6'; Grade='eq=contrast=1.09:saturation=1.10:brightness=-0.02'; Effect='control'; Title='' },
    [pscustomobject]@{ Id='030'; Asset='A04'; Frames=24; Motion='static'; Grade='eq=contrast=1.16:saturation=1.18'; Effect='glitch'; Title='TV GIRL' },
    [pscustomobject]@{ Id='040a'; Asset='A05'; Frames=12; Motion='push8'; Grade='eq=contrast=1.10:saturation=0.95'; Effect='fragment'; Title='' },
    [pscustomobject]@{ Id='040b'; Asset='A06'; Frames=12; Motion='pull6'; Grade='eq=contrast=1.15:saturation=0.76'; Effect='fragment'; Title='' },
    [pscustomobject]@{ Id='040c'; Asset='A07'; Frames=12; Motion='push12'; Grade='eq=contrast=1.12:saturation=1.12'; Effect='fragment'; Title='' },
    [pscustomobject]@{ Id='050'; Asset='A08'; Frames=84; Motion='push4'; Grade='eq=contrast=1.04:saturation=1.05:gamma=1.02'; Effect='clean'; Title='' },
    [pscustomobject]@{ Id='060'; Asset='A09'; Frames=48; Motion='pull8'; Grade='eq=contrast=1.03:saturation=1.06:brightness=0.01'; Effect='clean'; Title='' },
    [pscustomobject]@{ Id='070'; Asset='A10'; Frames=60; Motion='push7'; Grade='eq=contrast=1.07:saturation=1.07'; Effect='clean'; Title='' },
    [pscustomobject]@{ Id='080'; Asset='A11'; Frames=36; Motion='panleft'; Grade='eq=contrast=1.02:saturation=0.92:brightness=0.02'; Effect='clean'; Title='' },
    [pscustomobject]@{ Id='090'; Asset='A12'; Frames=60; Motion='static'; Grade='eq=contrast=1.03:saturation=0.90:brightness=0.015'; Effect='clean'; Title='' },
    [pscustomobject]@{ Id='100'; Asset='A13'; Frames=24; Motion='push8'; Grade='eq=contrast=1.11:saturation=0.88'; Effect='snow'; Title='' },
    [pscustomobject]@{ Id='110'; Asset='A14'; Frames=60; Motion='static'; Grade='eq=contrast=1.10:saturation=0.88:brightness=-0.015'; Effect='dusk'; Title='' },
    [pscustomobject]@{ Id='120'; Asset='A15'; Frames=72; Motion='pull6'; Grade='eq=contrast=1.04:saturation=1.08:brightness=0.01'; Effect='bloom'; Title='' },
    [pscustomobject]@{ Id='130'; Asset='A16'; Frames=72; Motion='pull5'; Grade='eq=contrast=1.03:saturation=1.10:brightness=0.015'; Effect='bloom'; Title='' },
    [pscustomobject]@{ Id='140'; Asset='A17'; Frames=36; Motion='static'; Grade='eq=contrast=1.05:saturation=1.08'; Effect='leaves'; Title='' },
    [pscustomobject]@{ Id='150'; Asset='A18'; Frames=84; Motion='push4'; Grade='eq=contrast=1.06:saturation=1.08:brightness=-0.01'; Effect='water'; Title='' },
    [pscustomobject]@{ Id='160'; Asset='A19'; Frames=60; Motion='pull8'; Grade='eq=contrast=1.02:saturation=1.09:brightness=0.01'; Effect='dream'; Title='' },
    [pscustomobject]@{ Id='170'; Asset='A20'; Frames=60; Motion='push5'; Grade='eq=contrast=1.08:saturation=1.12:brightness=0.02'; Effect='cyan'; Title='' },
    [pscustomobject]@{ Id='180'; Asset='A07'; Frames=24; Motion='push12'; Grade='eq=contrast=1.17:saturation=1.10'; Effect='glitch'; Title='' },
    [pscustomobject]@{ Id='190'; Asset='A03'; Frames=48; Motion='static'; Grade='eq=contrast=1.06:saturation=0.68:brightness=0.03'; Effect='pale'; Title='' },
    [pscustomobject]@{ Id='200'; Asset='A02'; Frames=48; Motion='pull6'; Grade='eq=contrast=1.11:saturation=1.12:brightness=-0.02'; Effect='control'; Title='' },
    [pscustomobject]@{ Id='210'; Asset='N01'; Frames=72; Motion='pull2'; Grade='eq=contrast=1.06:saturation=1.02:brightness=0.01'; Effect='release'; Title='EXIT THE SIGNAL' }
)

$totalFrames = ($shots | Measure-Object -Property Frames -Sum).Sum
if ($totalFrames -ne 1104) {
    throw "Storyboard frame total is $totalFrames; expected 1104."
}

function Get-ZoomExpression {
    param([string]$Motion, [int]$Frames)

    $denominator = [Math]::Max(1, $Frames - 1)
    switch ($Motion) {
        'push2'  { return "1+0.02*on/$denominator" }
        'push4'  { return "1+0.04*on/$denominator" }
        'push5'  { return "1+0.05*on/$denominator" }
        'push6'  { return "1+0.06*on/$denominator" }
        'push7'  { return "1+0.07*on/$denominator" }
        'push8'  { return "1+0.08*on/$denominator" }
        'push12' { return "1+0.12*on/$denominator" }
        'pull2'  { return "1.02-0.02*on/$denominator" }
        'pull5'  { return "1.05-0.05*on/$denominator" }
        'pull6'  { return "1.06-0.06*on/$denominator" }
        'pull8'  { return "1.08-0.08*on/$denominator" }
        'panleft' { return '1.06' }
        default  { return '1' }
    }
}

function Get-XExpression {
    param([string]$Motion, [int]$Frames)

    if ($Motion -eq 'panleft') {
        $denominator = [Math]::Max(1, $Frames - 1)
        return "(iw-iw/zoom)*(0.75-0.5*on/$denominator)"
    }
    return 'iw/2-(iw/zoom/2)'
}

function Get-EffectFilters {
    param([string]$Effect, [double]$Duration)

    switch ($Effect) {
        'control' {
            return @(
                'drawgrid=w=iw:h=6:t=1:c=black@0.055',
                "noise=alls=3:allf=t+u",
                "eq=brightness='0.006*sin(18*t)':eval=frame"
            )
        }
        'glitch' {
            return @(
                'rgbashift=rh=4:bh=-4:rv=1:bv=-1',
                'drawgrid=w=iw:h=5:t=1:c=black@0.08',
                'noise=alls=9:allf=t+u'
            )
        }
        'fragment' {
            return @('rgbashift=rh=2:bh=-2', 'noise=alls=5:allf=t+u')
        }
        'snow' {
            return @("noise=alls=3:allf=t+u", "eq=brightness='0.008*sin(3*t)':eval=frame")
        }
        'dusk' {
            return @("eq=brightness='-0.012+0.02*sin(1.3*t)':eval=frame")
        }
        'bloom' {
            return @("eq=brightness='0.006+0.012*sin(1.1*t)':eval=frame", 'colorbalance=gs=.015:bs=.01')
        }
        'leaves' {
            return @("eq=brightness='0.008*sin(2.1*t)':eval=frame")
        }
        'water' {
            return @("eq=brightness='0.006*sin(1.6*t)':eval=frame", 'colorbalance=bs=.035:gs=.01')
        }
        'dream' {
            return @("fade=t=in:st=0:d=0.333", 'colorbalance=bs=.025:rs=.01')
        }
        'cyan' {
            return @("eq=brightness='0.01+0.018*sin(1.2*t)':eval=frame", 'colorbalance=bs=.05:gs=.025')
        }
        'pale' {
            return @('curves=all=0/0.03 .55/.63 1/1')
        }
        'release' {
            return @(
                "fade=t=in:st=0:d=0.5",
                "eq=brightness='0.004+0.02*t/$Duration':eval=frame",
                'colorbalance=rs=.025:gs=.012'
            )
        }
        default {
            return @('noise=alls=2:allf=t+u')
        }
    }
}

function Get-TitleFilter {
    param([string]$Title, [double]$Duration)

    if ([string]::IsNullOrWhiteSpace($Title)) {
        return $null
    }

    $fontFile = 'C\:/Windows/Fonts/arialbd.ttf'
    if ($Title -eq 'TV GIRL') {
        return "drawtext=fontfile='$fontFile':text='TV GIRL':fontcolor=white:fontsize=128:x=(w-text_w)/2:y=(h-text_h)/2:borderw=2:bordercolor=black@0.55:alpha='if(lt(t,0.12),t/0.12,if(lt(t,$($Duration-0.16)),1,($Duration-t)/0.16))'"
    }

    return "drawtext=fontfile='$fontFile':text='EXIT THE SIGNAL':fontcolor=white:fontsize=58:x=(w-text_w)/2:y=h*0.79:borderw=2:bordercolor=black@0.52:alpha='if(lt(t,1.2),0,if(lt(t,1.7),(t-1.2)/0.5,1))'"
}

Write-Host "Rendering $($shots.Count) editorial segments ($totalFrames frames)..."
$shotFiles = New-Object System.Collections.Generic.List[string]

foreach ($shot in $shots) {
    $source = $assets[$shot.Asset]
    $duration = $shot.Frames / 24.0
    $zoom = Get-ZoomExpression -Motion $shot.Motion -Frames $shot.Frames
    $x = Get-XExpression -Motion $shot.Motion -Frames $shot.Frames
    $filters = New-Object System.Collections.Generic.List[string]
    $filters.Add('scale=1944:1080:flags=lanczos')
    $filters.Add('crop=1920:1080')
    $filters.Add("zoompan=z='$zoom':x='$x':y='ih/2-(ih/zoom/2)':d=1:s=1920x1080:fps=24")
    $filters.Add($shot.Grade)
    foreach ($effectFilter in (Get-EffectFilters -Effect $shot.Effect -Duration $duration)) {
        $filters.Add($effectFilter)
    }
    $filters.Add('vignette=PI/5.5')
    $titleFilter = Get-TitleFilter -Title $shot.Title -Duration $duration
    if ($null -ne $titleFilter) {
        $filters.Add($titleFilter)
    }
    $filters.Add('format=yuv420p')

    $shotFile = Join-Path $shotRoot ("shot-{0}.mp4" -f $shot.Id)
    $filterGraph = $filters -join ','

    if ((Test-Path -LiteralPath $shotFile) -and ((Get-Item -LiteralPath $shotFile).Length -gt 4096)) {
        $shotFiles.Add($shotFile)
        continue
    }

    & $ffmpeg -hide_banner -loglevel error -y `
        -framerate 24 -loop 1 -i $source `
        -vf $filterGraph -frames:v $shot.Frames -an `
        -c:v libx264 -preset medium -crf 16 -pix_fmt yuv420p -r 24 `
        $shotFile
    if ($LASTEXITCODE -ne 0) {
        throw "FFmpeg failed while rendering shot $($shot.Id)."
    }
    $shotFiles.Add($shotFile)
}

$concatFile = Join-Path $workRoot 'shots.concat.txt'
$concatLines = $shotFiles | ForEach-Object {
    $escaped = $_.Replace("'", "'\''")
    "file '$escaped'"
}
Set-Content -LiteralPath $concatFile -Value $concatLines -Encoding utf8

Write-Host 'Assembling picture master...'
& $ffmpeg -hide_banner -loglevel error -y `
    -f concat -safe 0 -i $concatFile `
    -vf 'fps=24,format=yuv420p' -frames:v 1104 -an `
    -c:v libx264 -preset slow -crf 16 -pix_fmt yuv420p -r 24 `
    $visualMaster
if ($LASTEXITCODE -ne 0) {
    throw 'FFmpeg failed while assembling the picture master.'
}

Write-Host 'Synthesizing original dream-pop soundtrack...'
$audioFilter = @"
[0:a]volume='0.022+0.013*between(t,6.5,39)':eval=frame,tremolo=f=0.12:d=0.35,lowpass=f=900,afade=t=in:st=0:d=2,afade=t=out:st=43:d=3[p1];
[1:a]volume='0.016+0.011*between(t,6.5,39)':eval=frame,tremolo=f=0.11:d=0.3,lowpass=f=1100,afade=t=in:st=1.5:d=2,afade=t=out:st=42:d=4[p2];
[2:a]volume='0.010+0.010*between(t,22,38)':eval=frame,tremolo=f=0.17:d=0.5,lowpass=f=1500,afade=t=in:st=5.5:d=2,afade=t=out:st=39:d=3[p3];
[3:a]lowpass=f=320,volume=0.018,afade=t=in:st=0:d=1.5,afade=t=out:st=43:d=3[room];
[4:a]volume='if(between(t,12,29.5),0.13*exp(-18*mod(t-12,0.625)),0)':eval=frame,lowpass=f=180[kick];
[5:a]highpass=f=2200,lowpass=f=7500,volume='if(between(t,12,29.5),0.016*exp(-50*mod(t-12,0.3125)),0)':eval=frame[hats];
[6:a]volume='0.20*(between(t,1.5,1.54)+between(t,4,4.05)+between(t,5,5.035)+between(t,5.5,5.535)+between(t,6,6.035)+between(t,41,41.04)+between(t,43,43.05))':eval=frame,highpass=f=700,lowpass=f=3200[clicks];
[7:a]highpass=f=350,lowpass=f=2600,volume='0.025*between(t,29.0,38.2)':eval=frame,afade=t=in:st=29:d=1.2,afade=t=out:st=37:d=1.2[water];
[p1][p2][p3][room][kick][hats][clicks][water]amix=inputs=8:normalize=0,acompressor=threshold=0.16:ratio=2.5:attack=20:release=250,alimiter=limit=0.86,loudnorm=I=-16:TP=-1.5:LRA=8,volume=1.3,aformat=sample_fmts=s16:sample_rates=48000:channel_layouts=stereo[out]
"@ -replace "`r?`n", ''

& $ffmpeg -hide_banner -loglevel error -y `
    -f lavfi -i 'sine=frequency=110:sample_rate=48000:duration=46' `
    -f lavfi -i 'sine=frequency=164.81:sample_rate=48000:duration=46' `
    -f lavfi -i 'sine=frequency=220:sample_rate=48000:duration=46' `
    -f lavfi -i 'anoisesrc=color=pink:sample_rate=48000:duration=46' `
    -f lavfi -i 'sine=frequency=58:sample_rate=48000:duration=46' `
    -f lavfi -i 'anoisesrc=color=white:sample_rate=48000:duration=46' `
    -f lavfi -i 'anoisesrc=color=white:sample_rate=48000:duration=46' `
    -f lavfi -i 'anoisesrc=color=brown:sample_rate=48000:duration=46' `
    -filter_complex $audioFilter -map '[out]' -t 46 `
    -c:a pcm_s16le $soundtrack
if ($LASTEXITCODE -ne 0) {
    throw 'FFmpeg failed while synthesizing the soundtrack.'
}

Write-Host 'Muxing final movie...'
& $ffmpeg -hide_banner -loglevel error -y `
    -i $visualMaster -i $soundtrack `
    -map 0:v:0 -map 1:a:0 -frames:v 1104 -t 46 `
    -c:v copy -c:a aac -b:a 256k -ar 48000 `
    -movflags +faststart -metadata title='TV GIRL: EXIT THE SIGNAL' `
    -metadata comment='Created from the TV Girl 0.1 storyboard at 24 fps.' `
    $movie
if ($LASTEXITCODE -ne 0) {
    throw 'FFmpeg failed while muxing the final movie.'
}

Write-Host 'Creating lightweight review encode...'
& $ffmpeg -hide_banner -loglevel error -y `
    -i $movie -vf 'scale=960:540:flags=lanczos' `
    -c:v libx264 -preset medium -crf 22 -pix_fmt yuv420p `
    -c:a aac -b:a 160k -movflags +faststart `
    $reviewMovie
if ($LASTEXITCODE -ne 0) {
    throw 'FFmpeg failed while creating the review encode.'
}

$probeJson = & $ffprobe -v error -show_entries 'format=duration:stream=index,codec_type,codec_name,width,height,r_frame_rate,avg_frame_rate,nb_frames,sample_rate,channels' -of json $movie
Set-Content -LiteralPath (Join-Path $outputRoot 'tv-girl-exit-the-signal.ffprobe.json') -Value $probeJson -Encoding utf8

if (-not $KeepIntermediates) {
    Remove-Item -LiteralPath $shotRoot -Recurse -Force
    Remove-Item -LiteralPath $concatFile -Force
}

Write-Host "Movie: $movie"
Write-Host "Review: $reviewMovie"
Write-Host "Frames: $totalFrames"
