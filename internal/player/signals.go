package player

import (
	"codeberg.org/dergs/tonearm/internal/signals"
	v1 "codeberg.org/dergs/tonearm/pkg/tidalapi/models/v1"
)

var ControllableStateChanged = signals.NewStatefulSignal(ControllableState{
	HasTrack:    false,
	PlayerReady: true,
})

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

// Holds the relevant information about the currently playing or last played track.
// This can be nil if no track is currently playing. This is especially the case when
// the player has just been created.
//
// The signal fires after the new track information has been retrieved from the TIDAL API
// but always before the track starts playing.
var TrackChanged = signals.NewStatefulSignal[*Track](nil)

// Holds the current volume of the player as reported by playbin
//
// The signal fires whenever the volume changes. This could be either due to user input,
// external volume control or MPRIS.
var VolumeChanged = signals.NewStatefulSignal[float64](1.0)
