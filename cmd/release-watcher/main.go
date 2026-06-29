// Package main starts the release-watcher service.
package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/neatflowcv/release-watcher/internal/app/server"
	pgxrepository "github.com/neatflowcv/release-watcher/internal/pkg/repository/pgx"
)

const (
	databaseURLEnv = "RELEASE_WATCHER_DATABASE_URL"
	addressEnv     = "RELEASE_WATCHER_ADDRESS"
)

func main() {
	err := run()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}
}

func run() error {
	dsn := os.Getenv(databaseURLEnv)
	if dsn == "" {
		return fmt.Errorf("%w: %s", errMissingEnvironment, databaseURLEnv)
	}

	projectRepository, err := pgxrepository.NewRepository(dsn)
	if err != nil {
		return fmt.Errorf("create project repository: %w", err)
	}
	defer projectRepository.Close()

	err = server.Run(projectRepository, os.Getenv(addressEnv))
	if err != nil {
		return fmt.Errorf("run server: %w", err)
	}

	return nil
}

var errMissingEnvironment = errors.New("missing environment variable")
