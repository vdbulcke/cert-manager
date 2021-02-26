package server

import (
	"context"
	"os"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/vdbulcke/cert-manager/data"
	"github.com/vdbulcke/cert-manager/data/database"
	"github.com/vdbulcke/cert-manager/data/database/migration"
)

// StartServer start the server
// this function is blocking
func StartServer() {
	// TODO: Get this from config
	v := data.NewValidation()
	db := database.NewSqliteDB("./sqlite.db")
	// new logger
	logger := hclog.New(&hclog.LoggerOptions{
		Name:       "cert-manager",
		Level:      hclog.LevelFromString("DEBUG"),
		JSONFormat: true,
	})

	// create backend
	certBackend := data.NewCertBackend(logger, db, v)

	// run DB migration
	err := migration.DBMigration(db, logger)
	if err != nil {
		logger.Error("Error running db migration")
		os.Exit(1)
	}

	// create the server
	s := makeServer(logger, v, certBackend)

	// create channel for capturing OS signal
	c := MakeOSSignalChannel()

	// start the server
	go func() {
		logger.Info("Staring server")

		err := s.ListenAndServe()
		if err != nil {
			logger.Error("fail to start the server", "err", err)
			os.Exit(1)
		}
	}()

	// Block until a signal is received.
	sig := <-c

	logger.Info("Got signal", "sig", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)

}
