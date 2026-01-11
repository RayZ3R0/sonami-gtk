package mpris

import (
	"log/slog"

	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/prop"
)

var log = slog.With("module", "MPRIS")

const (
	mediaPlayerInterface = "org.mpris.MediaPlayer2"
	playerInterface      = "org.mpris.MediaPlayer2.Player"
)

func NewMprisServer(name string, desktopEntry string, identity string) *Server {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		panic(err)
	}

	object := NewMprisDBusObject(desktopEntry, identity)
	server := &Server{
		dbusConnection: conn,
		dbusName:       name,
		object:         object,
	}

	server.dbusConnection.Export(object, "/org/mpris/MediaPlayer2", mediaPlayerInterface)
	server.dbusConnection.Export(object, "/org/mpris/MediaPlayer2", playerInterface)
	server.properties, _ = prop.Export(server.dbusConnection, "/org/mpris/MediaPlayer2", object.Properties)

	return server
}
