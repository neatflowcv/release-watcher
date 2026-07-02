package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

func newTagsRequest(
	ctx context.Context,
	apiBaseURL *url.URL,
	repository repository,
	limit int,
) (*http.Request, error) {
	requestURL := *apiBaseURL
	requestURL.Path = path.Join(
		requestURL.Path,
		"repos",
		repository.owner,
		repository.name,
		"tags",
	)

	query := requestURL.Query()
	query.Set("per_page", strconv.Itoa(limit))
	requestURL.RawQuery = query.Encode()

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		requestURL.String(),
		http.NoBody,
	)
	if err != nil {
		return nil, fmt.Errorf("create github tags request: %w", err)
	}

	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("X-Github-Api-Version", "2022-11-28")

	return request, nil
}

func newTagRefRequest(
	ctx context.Context,
	apiBaseURL *url.URL,
	repository repository,
	tag string,
) (*http.Request, error) {
	return newGitHubRequest(
		ctx,
		apiBaseURL,
		"repos",
		repository.owner,
		repository.name,
		"git",
		"ref",
		"tags",
		tag,
	)
}

func newTagObjectRequest(
	ctx context.Context,
	apiBaseURL *url.URL,
	repository repository,
	sha string,
) (*http.Request, error) {
	return newGitHubRequest(
		ctx,
		apiBaseURL,
		"repos",
		repository.owner,
		repository.name,
		"git",
		"tags",
		sha,
	)
}

func newGitHubRequest(
	ctx context.Context,
	apiBaseURL *url.URL,
	pathParts ...string,
) (*http.Request, error) {
	requestURL := *apiBaseURL
	requestURL.Path = path.Join(append([]string{requestURL.Path}, pathParts...)...)

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		requestURL.String(),
		http.NoBody,
	)
	if err != nil {
		return nil, fmt.Errorf("create github request: %w", err)
	}

	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("X-Github-Api-Version", "2022-11-28")

	return request, nil
}
