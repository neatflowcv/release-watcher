package version

// Filter determines whether a version string should be considered.
type Filter interface {
	FilterVersion(version string) bool
}
