package github

import "strings"

type vPrefixSemanticVersionParser struct {
	parser semanticVersionParser
}

func newVPrefixSemanticVersionParser(
	parser semanticVersionParser,
) vPrefixSemanticVersionParser {
	return vPrefixSemanticVersionParser{
		parser: parser,
	}
}

func (p vPrefixSemanticVersionParser) Parse(tag string) (semanticVersion, bool) {
	version := strings.TrimPrefix(tag, "v")

	return p.parser.Parse(version)
}
