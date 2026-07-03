// Package domain contains release watcher domain models.
package domain

// Project is a release target watched by the service.
type Project struct {
	name                string
	url                 string
	versionSource       VersionSource
	versionPattern      VersionPattern
	detectedVersion     string
	acknowledgedVersion string
}

// VersionSource identifies where project versions are discovered.
type VersionSource string

const (
	// VersionSourceGitHubReleases discovers versions from GitHub releases.
	VersionSourceGitHubReleases VersionSource = "github_releases"
	// VersionSourceGitHubTags discovers versions from GitHub repository tags.
	VersionSourceGitHubTags VersionSource = "github_tags"
)

// VersionPattern identifies how version strings are recognized.
type VersionPattern string

const (
	// VersionPatternSemver matches versions like 20.2.1.
	VersionPatternSemver VersionPattern = "semver"
	// VersionPatternVSemver matches versions like v20.2.1.
	VersionPatternVSemver VersionPattern = "v_semver"
	// VersionPatternPveVersion matches versions like 20.2.0-pve1 or 20.2.0-pve2.
	VersionPatternPveVersion VersionPattern = "pve_version"
	// VersionPatternGodotStable matches Godot stable versions like 4.7-stable.
	VersionPatternGodotStable VersionPattern = "godot_stable"
)

// NewProject creates a project with its current release tracking state.
func NewProject(
	name string,
	url string,
	versionSource VersionSource,
	versionPattern VersionPattern,
	detectedVersion string,
	acknowledgedVersion string,
) *Project {
	return &Project{
		name:                name,
		url:                 url,
		versionSource:       defaultVersionSource(versionSource),
		versionPattern:      defaultVersionPattern(versionPattern),
		detectedVersion:     detectedVersion,
		acknowledgedVersion: acknowledgedVersion,
	}
}

// Name returns the project name.
func (p *Project) Name() string {
	return p.name
}

// URL returns the release source URL.
func (p *Project) URL() string {
	return p.url
}

// VersionSource returns where project versions are discovered.
func (p *Project) VersionSource() VersionSource {
	return p.versionSource
}

// VersionPattern returns how version strings are recognized.
func (p *Project) VersionPattern() VersionPattern {
	return p.versionPattern
}

// DetectedVersion returns the latest version detected by the service.
func (p *Project) DetectedVersion() string {
	return p.detectedVersion
}

// AcknowledgedVersion returns the version acknowledged by the user.
func (p *Project) AcknowledgedVersion() string {
	return p.acknowledgedVersion
}

func defaultVersionSource(versionSource VersionSource) VersionSource {
	if versionSource == "" {
		return VersionSourceGitHubTags
	}

	return versionSource
}

func defaultVersionPattern(versionPattern VersionPattern) VersionPattern {
	if versionPattern == "" {
		return VersionPatternVSemver
	}

	return versionPattern
}
