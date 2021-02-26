package server

import (
	"os"
	"os/signal"
)

// MakeOSSignalChannel create a channel capturing OS Signals
func MakeOSSignalChannel() chan os.Signal {
	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)
	return c
}
