// Package domain contains release watcher domain models.
package domain

// Project is a release target watched by the service.
type Project struct {
	name                string
	url                 string
	detectedVersion     string
	acknowledgedVersion string
}

// NewProject creates a project with its current release tracking state.
func NewProject(name, url, detectedVersion, acknowledgedVersion string) *Project {
	return &Project{
		name:                name,
		url:                 url,
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

// DetectedVersion returns the latest version detected by the service.
func (p *Project) DetectedVersion() string {
	return p.detectedVersion
}

// AcknowledgedVersion returns the version acknowledged by the user.
func (p *Project) AcknowledgedVersion() string {
	return p.acknowledgedVersion
}
