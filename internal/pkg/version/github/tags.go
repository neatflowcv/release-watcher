// Package github provides GitHub-backed version fetchers.
package github

import (
	"context"
	"errors"
	"fmt"

	"github.com/neatflowcv/release-watcher/internal/pkg/domain"
	"github.com/neatflowcv/release-watcher/internal/pkg/version"
	githubclient "github.com/neatflowcv/release-watcher/internal/pkg/version/github/client"
)

const latestTagLimit = 30

var (
	_ version.Fetcher = (*TagsFetcher)(nil)

	errNoSemanticTags = errors.New("github repository has no semantic version tags")
	errNoVerifiedTags = errors.New(
		"github repository has no verified semantic version tags",
	)
)

// TagsFetcher retrieves the latest version from a GitHub repository tag list.
type TagsFetcher struct {
	filter version.Filter
}

// NewTagsFetcher creates a GitHub tag version fetcher.
func NewTagsFetcher(filter version.Filter) *TagsFetcher {
	return &TagsFetcher{
		filter: filter,
	}
}

// GetLatestVersion returns the highest semantic version from recent GitHub tags.
func (f *TagsFetcher) GetLatestVersion(
	ctx context.Context,
	project *domain.Project,
) (string, error) {
	apiClient, err := githubclient.NewClient(project.URL())
	if err != nil {
		return "", fmt.Errorf("create github client: %w", err)
	}

	tag, err := f.getLatestTag(ctx, apiClient)
	if err != nil {
		return "", fmt.Errorf("get latest github tag: %w", err)
	}

	return tag, nil
}

func (f *TagsFetcher) getLatestTag(
	ctx context.Context,
	apiClient *githubclient.Client,
) (string, error) {
	// Thirty recent tags should be enough because release version tags are
	// expected to stay near the front of GitHub's tag list.
	tags, err := apiClient.ListTags(ctx, latestTagLimit)
	if err != nil {
		return "", fmt.Errorf("list github tags: %w", err)
	}

	tag, err := f.highestVerifiedSemanticTag(ctx, apiClient, tags)
	if err != nil {
		return "", err
	}

	return tag, nil
}

func (f *TagsFetcher) highestVerifiedSemanticTag(
	ctx context.Context,
	apiClient *githubclient.Client,
	tags []*githubclient.Tag,
) (string, error) {
	semanticTags := parseSemanticTags(tags, f.filter)
	if len(semanticTags) == 0 {
		return "", errNoSemanticTags
	}

	verifiedTags, err := f.verifiedSemanticTags(ctx, apiClient, semanticTags)
	if err != nil {
		return "", err
	}

	latestTag, ok := highestSemanticTag(verifiedTags)
	if !ok {
		return "", errNoVerifiedTags
	}

	return latestTag, nil
}

func (f *TagsFetcher) verifiedSemanticTags(
	ctx context.Context,
	apiClient *githubclient.Client,
	tags []semanticTag,
) ([]semanticTag, error) {
	verifiedTags := make([]semanticTag, 0, len(tags))

	for _, tag := range tags {
		verified, err := f.isVerifiedTag(ctx, apiClient, tag.tag)
		if err != nil {
			return nil, err
		}

		if verified {
			verifiedTags = append(verifiedTags, tag)
		}
	}

	return verifiedTags, nil
}

func (f *TagsFetcher) isVerifiedTag(
	ctx context.Context,
	apiClient *githubclient.Client,
	tag *githubclient.Tag,
) (bool, error) {
	tagRef, err := apiClient.GetAnnotatedTagRef(ctx, tag.Name)
	if errors.Is(err, githubclient.ErrTagIsNotAnnotated) {
		return true, nil
	}

	if err != nil {
		return false, fmt.Errorf("get github tag ref: %w", err)
	}

	detail, err := apiClient.GetTagDetail(ctx, tagRef)
	if err != nil {
		return false, fmt.Errorf("get github tag detail: %w", err)
	}

	return detail.Verification.Verified, nil
}
