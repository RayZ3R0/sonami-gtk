package player

import (
	"github.com/RayZ3R0/sonami-gtk/internal/settings"
	"github.com/RayZ3R0/sonami-gtk/internal/signals"
	v1 "github.com/RayZ3R0/sonami-gtk/pkg/tidalapi/models/v1"
	"github.com/RayZ3R0/sonami-gtk/pkg/sonami"
)

// Holds the current audio stream quality.
var AudioStreamQuality = signals.NewStatefulSignal[*StreamQuality](nil)

// Holds the current playback state including the expected duration and playing position of the
// currently playing track, if any.
//
// This signal fires when the player changes between playback states such as buffering, playing, paused, or stopped.
// This signal also fires whenever the expected track duration changes or at least every 250 milliseconds during playback.
var PlaybackStateChanged = signals.NewStatefulSignal[*PlaybackState](&PlaybackState{Status: PlaybackStatusStopped})

// Holds the definitive playback quality of the currently playing or last played track
// This can be lower than the chosen quality depending on the tracks original quality
// or any other reason TIDAL's API may have for offering a lower quality.
//
// The signal fires shortly after the new track information has been broadcasted on
// the TrackChanged signal. The signal may fire either before the track starts playing
// in case of a cold start or extremely shortly after a track has started playing in
// case of gapless playback.
var PlaybackQualityChanged = signals.NewStatefulSignal[v1.AudioQuality](v1.AudioQualityHighResLossless)

// Holds the user-selected repeat mode.
//
// The signal fires whenever the user changes the repeat mode.
var RepeatModeChanged = signals.NewStatefulSignal[RepeatMode](RepeatModeNone)

// Holds the seed chosen when the user selects shuffle mode.
//
// The signal fires whenever the user toggles the shuffle mode.
// The value is 0 if shuffle mode is disabled.
var ShuffleStateChanged = signals.NewStatefulSignal[bool](false)

var SourceChanged = signals.NewStatefulSignal[sonami.PlaybackSource](nil)

// Holds the relevant information about the currently playing or last played track.
// This can be nil if no track is currently playing. This is especially the case when
// the player has just been created.
//
// The signal fires after the new track information has been retrieved from the TIDAL API
// but always before the track starts playing.
var TrackChanged = signals.NewStatefulSignal[sonami.Track](nil)

// Holds the current volume of the player as reported by playbin
//
// The signal fires whenever the volume changes. This could be either due to user input,
// external volume control or MPRIS.
var VolumeChanged = signals.NewStatefulSignal[float64](settings.Player().GetVolume())
