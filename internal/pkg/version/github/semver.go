package github

import (
	"strconv"
	"strings"

	"github.com/neatflowcv/release-watcher/internal/pkg/version"
	githubclient "github.com/neatflowcv/release-watcher/internal/pkg/version/github/client"
)

const semverPartCount = 3

type semanticVersionParser interface {
	Parse(tag string) (semanticVersion, bool)
}

type fixedSemanticVersionParser struct{}

type semanticVersion struct {
	major int
	minor int
	patch int
}

type semanticTag struct {
	name    string
	tag     *githubclient.Tag
	version semanticVersion
}

func parseSemanticTags(
	tags []*githubclient.Tag,
	filter version.Filter,
) []semanticTag {
	parser := newVPrefixSemanticVersionParser(fixedSemanticVersionParser{})
	semanticTags := make([]semanticTag, 0, len(tags))

	for _, tag := range tags {
		if !filter.FilterVersion(tag.Name) {
			continue
		}

		version, ok := parser.Parse(tag.Name)
		if !ok {
			continue
		}

		semanticTags = append(semanticTags, semanticTag{
			name:    tag.Name,
			tag:     tag,
			version: version,
		})
	}

	return semanticTags
}

func highestSemanticTag(semanticTags []semanticTag) (string, bool) {
	if len(semanticTags) == 0 {
		return "", false
	}

	highestTag := semanticTags[0]

	for _, tag := range semanticTags[1:] {
		if compareSemanticVersion(tag.version, highestTag.version) > 0 {
			highestTag = tag
		}
	}

	return highestTag.name, true
}

func (p fixedSemanticVersionParser) Parse(tag string) (semanticVersion, bool) {
	parts := strings.Split(tag, ".")
	if len(parts) != semverPartCount {
		return emptySemanticVersion(), false
	}

	major, parsed := parseSemanticPart(parts[0])
	if !parsed {
		return emptySemanticVersion(), false
	}

	minor, parsed := parseSemanticPart(parts[1])
	if !parsed {
		return emptySemanticVersion(), false
	}

	patch, parsed := parseSemanticPart(parts[2])
	if !parsed {
		return emptySemanticVersion(), false
	}

	return semanticVersion{
		major: major,
		minor: minor,
		patch: patch,
	}, true
}

func parseSemanticPart(value string) (int, bool) {
	if value == "" || hasInvalidLeadingZero(value) {
		return 0, false
	}

	number, err := strconv.Atoi(value)
	if err != nil {
		return 0, false
	}

	return number, true
}

func emptySemanticVersion() semanticVersion {
	return semanticVersion{
		major: 0,
		minor: 0,
		patch: 0,
	}
}

func hasInvalidLeadingZero(value string) bool {
	return len(value) > 1 && strings.HasPrefix(value, "0")
}

func compareSemanticVersion(left, right semanticVersion) int {
	if compared := compareInt(left.major, right.major); compared != 0 {
		return compared
	}

	if compared := compareInt(left.minor, right.minor); compared != 0 {
		return compared
	}

	return compareInt(left.patch, right.patch)
}

func compareInt(left, right int) int {
	switch {
	case left > right:
		return 1
	case left < right:
		return -1
	default:
		return 0
	}
}
