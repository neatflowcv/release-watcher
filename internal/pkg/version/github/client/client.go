// Package client provides GitHub repository API clients.
package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	defaultAPIBaseURL       = "https://api.github.com"
	defaultTimeout          = 10 * time.Second
	ownerPathIndex          = 0
	repositoryNamePathIndex = 1
	repositoryPathPartCount = 2
)

var (
	errInvalidGitHubURL = errors.New("invalid github repository url")
	errTagsRequest      = errors.New("github tags request failed")

	// ErrTagIsNotAnnotated indicates that the tag ref points directly to an
	// object instead of an annotated tag object.
	ErrTagIsNotAnnotated = errors.New("github tag is not annotated")
)

// Client calls GitHub repository APIs.
type Client struct {
	apiBaseURL *url.URL
	httpClient *http.Client
	repository repository
}

type repository struct {
	owner string
	name  string
}

// NewClient creates a GitHub API client for a repository URL.
func NewClient(rawURL string) (*Client, error) {
	repository, err := parseRepositoryURL(rawURL)
	if err != nil {
		return nil, fmt.Errorf("parse github repository url: %w", err)
	}

	apiBaseURL, err := url.Parse(defaultAPIBaseURL)
	if err != nil {
		panic(fmt.Sprintf("parse default github api base url: %v", err))
	}

	return &Client{
		apiBaseURL: apiBaseURL,
		httpClient: &http.Client{
			Transport:     nil,
			CheckRedirect: nil,
			Jar:           nil,
			Timeout:       defaultTimeout,
		},
		repository: repository,
	}, nil
}

func parseRepositoryURL(rawURL string) (repository, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return repository{}, fmt.Errorf("%w: %w", errInvalidGitHubURL, err)
	}

	if !isGitHubHost(parsedURL.Hostname()) {
		return repository{}, errInvalidGitHubURL
	}

	parts := strings.Split(strings.Trim(parsedURL.Path, "/"), "/")
	if len(parts) < repositoryPathPartCount {
		return repository{}, errInvalidGitHubURL
	}

	name := strings.TrimSuffix(parts[repositoryNamePathIndex], ".git")
	if parts[ownerPathIndex] == "" || name == "" {
		return repository{}, errInvalidGitHubURL
	}

	return repository{
		owner: parts[ownerPathIndex],
		name:  name,
	}, nil
}

func isGitHubHost(host string) bool {
	return strings.EqualFold(host, "github.com") ||
		strings.EqualFold(host, "www.github.com")
}

// ListTags returns recent repository tags up to the requested limit.
func (c *Client) ListTags(ctx context.Context, limit int) ([]*Tag, error) {
	request, err := newTagsRequest(ctx, c.apiBaseURL, c.repository, limit)
	if err != nil {
		return nil, err
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("send github tags request: %w", err)
	}
	defer func() {
		_ = response.Body.Close()
	}()

	tags, err := readTags(response)
	if err != nil {
		return nil, err
	}

	return tags, nil
}

// GetTagDetail returns the annotated tag object detail for a repository tag.
func (c *Client) GetTagDetail(
	ctx context.Context,
	tagRef *TagRef,
) (*TagDetail, error) {
	request, err := newTagObjectRequest(ctx, c.apiBaseURL, c.repository, tagRef.Object.SHA)
	if err != nil {
		return nil, err
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("send github tag object request: %w", err)
	}
	defer func() {
		_ = response.Body.Close()
	}()

	tagDetail, err := readTagDetail(response)
	if err != nil {
		return nil, err
	}

	return tagDetail, nil
}

// GetAnnotatedTagRef returns the annotated git ref for a repository tag.
func (c *Client) GetAnnotatedTagRef(
	ctx context.Context,
	tag string,
) (*TagRef, error) {
	request, err := newTagRefRequest(ctx, c.apiBaseURL, c.repository, tag)
	if err != nil {
		return nil, err
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("send github tag ref request: %w", err)
	}
	defer func() {
		_ = response.Body.Close()
	}()

	tagRef, err := readTagRef(response)
	if err != nil {
		return nil, err
	}

	// Non-tag refs are lightweight tags. They do not have annotated tag
	// details, so exclude them from this annotated-tag-specific method.
	if tagRef.Object.Type != "tag" {
		return nil, ErrTagIsNotAnnotated
	}

	return tagRef, nil
}

func readTags(response *http.Response) ([]*Tag, error) {
	if response.StatusCode != http.StatusOK {
		return nil, readGitHubErrorResponse(response)
	}

	var tags []*Tag

	decoder := json.NewDecoder(response.Body)

	err := decoder.Decode(&tags)
	if err != nil {
		return nil, fmt.Errorf("decode github tags response: %w", err)
	}

	return tags, nil
}

func readTagRef(response *http.Response) (*TagRef, error) {
	if response.StatusCode != http.StatusOK {
		return nil, readGitHubErrorResponse(response)
	}

	var tagRef TagRef

	decoder := json.NewDecoder(response.Body)

	err := decoder.Decode(&tagRef)
	if err != nil {
		return nil, fmt.Errorf("decode github tag ref response: %w", err)
	}

	return &tagRef, nil
}

func readTagDetail(response *http.Response) (*TagDetail, error) {
	if response.StatusCode != http.StatusOK {
		return nil, readGitHubErrorResponse(response)
	}

	var tagDetail TagDetail

	decoder := json.NewDecoder(response.Body)

	err := decoder.Decode(&tagDetail)
	if err != nil {
		return nil, fmt.Errorf("decode github tag object response: %w", err)
	}

	return &tagDetail, nil
}

func readGitHubErrorResponse(response *http.Response) error {
	body, readErr := io.ReadAll(response.Body)
	if readErr != nil {
		return fmt.Errorf(
			"read github error response: %w",
			readErr,
		)
	}

	return fmt.Errorf(
		"%w: status %d: %s",
		errTagsRequest,
		response.StatusCode,
		strings.TrimSpace(string(body)),
	)
}
