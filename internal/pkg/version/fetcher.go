// Package version defines release version lookup boundaries.
package version

import (
	"context"

	"github.com/neatflowcv/release-watcher/internal/pkg/domain"
)

// Fetcher retrieves the latest release version for a project.
type Fetcher interface {
	GetLatestVersion(ctx context.Context, project *domain.Project) (string, error)
}
