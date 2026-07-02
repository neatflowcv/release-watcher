package client

// Tag is a GitHub repository tag response.
type Tag struct {
	Name       string    `json:"name"`
	ZipballURL string    `json:"zipball_url"`
	TarballURL string    `json:"tarball_url"`
	Commit     TagCommit `json:"commit"`
	NodeID     string    `json:"node_id"`
}

// TagCommit is the commit object embedded in a GitHub tag response.
type TagCommit struct {
	SHA string `json:"sha"`
	URL string `json:"url"`
}

// TagDetail is a GitHub annotated tag object response.
type TagDetail struct {
	NodeID       string                `json:"node_id"`
	SHA          string                `json:"sha"`
	URL          string                `json:"url"`
	Tagger       TagDetailTagger       `json:"tagger"`
	Object       TagDetailObject       `json:"object"`
	Tag          string                `json:"tag"`
	Message      string                `json:"message"`
	Verification TagDetailVerification `json:"verification"`
}

// TagDetailTagger is the tagger object embedded in a GitHub tag detail.
type TagDetailTagger struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Date  string `json:"date"`
}

// TagDetailObject is the tagged git object embedded in a GitHub tag detail.
type TagDetailObject struct {
	SHA  string `json:"sha"`
	Type string `json:"type"`
	URL  string `json:"url"`
}

// TagDetailVerification contains GitHub's tag signature verification result.
type TagDetailVerification struct {
	Verified   bool    `json:"verified"`
	Reason     string  `json:"reason"`
	Signature  *string `json:"signature"`
	Payload    *string `json:"payload"`
	VerifiedAt *string `json:"verified_at"`
}

// TagRef is a GitHub git ref response for a repository tag.
type TagRef struct {
	Ref    string       `json:"ref"`
	NodeID string       `json:"node_id"`
	URL    string       `json:"url"`
	Object TagRefObject `json:"object"`
}

// TagRefObject is the git object referenced by a GitHub tag ref.
type TagRefObject struct {
	SHA  string `json:"sha"`
	Type string `json:"type"`
	URL  string `json:"url"`
}
