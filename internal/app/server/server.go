// Package server provides the release-watcher REST API.
package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/neatflowcv/release-watcher/internal/pkg/domain"
	"github.com/neatflowcv/release-watcher/internal/pkg/repository"
)

const (
	defaultAddress    = ":8080"
	readHeaderTimeout = 5 * time.Second
)

type createProjectInput struct {
	Body createProjectBody
}

type projectOutput struct {
	Body projectResponseBody
}

type listProjectsOutput struct {
	Body []projectResponseBody
}

type createProjectBody struct {
	Name string `doc:"Project name"               json:"name"`
	URL  string `doc:"Project release source URL" json:"url"`
}

type projectResponseBody struct {
	Name                string `doc:"Project name"                           json:"name"`
	URL                 string `doc:"Project release source URL"             json:"url"`
	DetectedVersion     string `doc:"Latest version detected by the service" json:"detectedVersion"`
	AcknowledgedVersion string `doc:"Version acknowledged by the user"       json:"acknowledgedVersion"`
}

// NewHandler creates an HTTP handler with all API routes registered.
func NewHandler(projectRepository repository.Repository) http.Handler {
	mux := http.NewServeMux()
	api := humago.New(mux, huma.DefaultConfig("release-watcher", "0.1.0"))

	registerProjectRoutes(api, projectRepository)

	return mux
}

// Run starts the REST API server.
func Run(projectRepository repository.Repository, address string) error {
	if address == "" {
		address = defaultAddress
	}

	server := http.Server{ //nolint:exhaustruct // Only the server settings used by this app are set.
		Addr:              address,
		Handler:           NewHandler(projectRepository),
		ReadHeaderTimeout: readHeaderTimeout,
	}

	err := server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("run server: %w", err)
	}

	return nil
}

func registerProjectRoutes(
	api huma.API,
	projectRepository repository.Repository,
) {
	huma.Register(api, huma.Operation{ //nolint:exhaustruct // Optional OpenAPI fields are intentionally omitted.
		OperationID:   "create-project",
		Method:        http.MethodPost,
		Path:          "/projects",
		Summary:       "Create project",
		DefaultStatus: http.StatusCreated,
	}, func(ctx context.Context, input *createProjectInput) (*projectOutput, error) {
		project := domain.NewProject(
			input.Body.Name,
			input.Body.URL,
			"",
			"",
		)

		err := projectRepository.CreateProject(ctx, project)
		if err != nil {
			return nil, huma.Error500InternalServerError("create project", err)
		}

		return &projectOutput{Body: projectBodyFromDomain(project)}, nil
	})

	huma.Register(api, huma.Operation{ //nolint:exhaustruct // Optional OpenAPI fields are intentionally omitted.
		OperationID: "list-projects",
		Method:      http.MethodGet,
		Path:        "/projects",
		Summary:     "List projects",
	}, func(ctx context.Context, _ *struct{}) (*listProjectsOutput, error) {
		projects, err := projectRepository.ListProjects(ctx)
		if err != nil {
			return nil, huma.Error500InternalServerError("list projects", err)
		}

		return &listProjectsOutput{Body: projectBodiesFromDomain(projects)}, nil
	})
}

func projectBodiesFromDomain(projects []*domain.Project) []projectResponseBody {
	bodies := make([]projectResponseBody, 0, len(projects))

	for _, project := range projects {
		bodies = append(bodies, projectBodyFromDomain(project))
	}

	return bodies
}

func projectBodyFromDomain(project *domain.Project) projectResponseBody {
	return projectResponseBody{
		Name:                project.Name(),
		URL:                 project.URL(),
		DetectedVersion:     project.DetectedVersion(),
		AcknowledgedVersion: project.AcknowledgedVersion(),
	}
}
