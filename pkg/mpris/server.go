package mpris

import (
	"log/slog"
	"reflect"
	"time"

	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/prop"
)

var capabilities = []string{
	"CanControl",
	"CanGoNext",
	"CanGoPrevious",
	"CanPause",
	"CanPlay",
	"CanSeek",
}

type Server struct {
	properties     *prop.Properties
	dbusConnection *dbus.Conn
	object         *MprisDBusObject
}

func (c *Server) EnableControl() {
	oldVar, _ := c.properties.Get(playerInterface, "CanControl")
	oldVal, ok := oldVar.Value().(bool)

	if !ok {
		log.Error("Unexpected non-boolean value in D-Bus PlaybackStatus value")
		for _, capability := range capabilities {
			c.properties.SetMust(playerInterface, capability, dbus.MakeVariant(true))
		}
		return
	}

	if oldVal != true {
		for _, capability := range capabilities {
			c.properties.SetMust(playerInterface, capability, dbus.MakeVariant(true))
		}
	}
}

func (c *Server) OnPause(cb func()) {
	c.object.OnPause = cb
}

func (c *Server) OnPlay(cb func()) {
	c.object.OnPlay = cb
}

func (c *Server) OnPlayPause(cb func()) {
	c.object.OnPlayPause = cb
}

func (c *Server) OnQuit(cb func()) {
	c.object.OnQuit = cb
}

func (c *Server) OnRaise(cb func()) {
	c.object.OnRaise = cb
}

func (c *Server) OnSeek(cb func(offset time.Duration)) {
	c.object.OnSeek = cb
}

func (c *Server) OnSetPosition(cb func(offset time.Duration)) {
	c.object.OnSetPosition = cb
}

func (c *Server) OnTrackNext(cb func()) {
	c.object.OnTrackNext = cb
}

func (c *Server) OnTrackPrevious(cb func()) {
	c.object.OnTrackPrevious = cb
}

func (c *Server) OnVolumeChanged(cb func(newVal float64)) {
	c.object.Properties[playerInterface]["Volume"].Callback = func(newVal *prop.Change) *dbus.Error {
		val, ok := newVal.Value.(float64)

		if !ok {
			slog.Error("Unexpected format from client for Volume value")
			return &dbus.ErrMsgInvalidArg
		}

		cb(val)
		return nil
	}
}

func (c *Server) SetPlaybackStatus(status PlaybackStatus) {
	oldVar, _ := c.properties.Get(playerInterface, "PlaybackStatus")
	oldVal, ok := oldVar.Value().(string)

	if !ok {
		log.Error("Unexpected non-string value in D-Bus PlaybackStatus value")
		c.properties.SetMust(playerInterface, "PlaybackStatus", dbus.MakeVariant(status))
		return
	}

	if PlaybackStatus(oldVal) != status {
		c.properties.SetMust(playerInterface, "PlaybackStatus", dbus.MakeVariant(status))
	}
}

func (c *Server) SetPosition(position time.Duration) {
	oldVar, _ := c.properties.Get(playerInterface, "Position")
	oldVal, ok := oldVar.Value().(int64)

	newVal := position.Microseconds()

	if !ok {
		log.Error("Unexpected non-integer value in D-Bus Position value")
		c.properties.SetMust(playerInterface, "Position", dbus.MakeVariant(newVal))
		return
	}

	if oldVal != newVal {
		c.properties.SetMust(playerInterface, "Position", dbus.MakeVariant(newVal))
	}
}

func (c *Server) SetTrackMetadata(metadata map[string]any) {
	oldVar, _ := c.properties.Get(playerInterface, "Metadata")
	oldVal, ok := oldVar.Value().(map[string]any)

	if !ok {
		log.Error("Unexpected non-(string keyed map) value in D-Bus Metadata value")
		c.properties.SetMust(playerInterface, "Metadata", dbus.MakeVariant(metadata))
		return
	}

	if !reflect.DeepEqual(oldVal, metadata) {
		c.properties.SetMust(playerInterface, "Metadata", dbus.MakeVariant(metadata))
	}
}

func (c *Server) SetVolume(volume float64) {
	oldVar, _ := c.properties.Get(playerInterface, "Volume")
	oldVal, ok := oldVar.Value().(float64)

	if !ok {
		log.Error("Unexpected non-float value in D-Bus Volume value")
		c.properties.SetMust(playerInterface, "Volume", dbus.MakeVariant(volume))
		return
	}

	if oldVal != volume {
		c.properties.SetMust(playerInterface, "Volume", dbus.MakeVariant(volume))
	}
}
