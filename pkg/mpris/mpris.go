package mpris

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/prop"
)

var log = slog.With("module", "MPRIS")

const (
	mediaPlayerInterface = "org.mpris.MediaPlayer2"
	playerInterface      = "org.mpris.MediaPlayer2.Player"
)

func NewMprisServer(name string) *Server {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		panic(err)
	}

	object := NewMprisDBusObject()
	server := &Server{
		dbusConnection: conn,
		object:         object,
	}

	server.dbusConnection.Export(object, "/org/mpris/MediaPlayer2", mediaPlayerInterface)
	server.dbusConnection.Export(object, "/org/mpris/MediaPlayer2", playerInterface)
	server.properties, _ = prop.Export(server.dbusConnection, "/org/mpris/MediaPlayer2", object.Properties)

	reply, err := server.dbusConnection.RequestName(name, dbus.NameFlagDoNotQueue)
	if reply != dbus.RequestNameReplyPrimaryOwner {
		fmt.Fprintln(os.Stderr, "name already taken")
		os.Exit(1)
	}

	return server
}
