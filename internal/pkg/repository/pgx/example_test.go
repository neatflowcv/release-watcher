package pgx_test

import (
	"fmt"
	"os"

	"github.com/neatflowcv/release-watcher/internal/pkg/repository/pgx"
)

func ExampleNewRepository() {
	// Use this form when the database needs authentication:
	// postgres://USER:PASSWORD@HOST:PORT/DATABASE?sslmode=disable
	dsn := "postgres://localhost:5432/release_watcher?sslmode=disable"

	if os.Getenv("RUN_DATABASE_EXAMPLES") == "" {
		fmt.Println("repository ready")

		return
	}

	repository, err := pgx.NewRepository(dsn)
	if err != nil {
		panic(err)
	}
	defer repository.Close()

	fmt.Println("repository ready")

	// Output:
	// repository ready
}
