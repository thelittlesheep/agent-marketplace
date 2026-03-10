package main

// Platform defines the interface for a social media platform parser.
// To add a new platform: create xxx.go, implement Platform, call RegisterPlatform() in init().
type Platform interface {
	Name() string
	MatchURL(rawURL string) bool
	NormalizeURL(rawURL string) (string, error)
	ParsePost(markdown string, sourceURL string) (*Post, error)
	ParseComments(markdown string) ([]*Comment, error)
}

// Post represents a social media post.
type Post struct {
	Author     string      `json:"author"`
	Content    string      `json:"content"`
	Images     []string    `json:"images,omitempty"`
	QuotedPost *Post       `json:"quoted_post,omitempty"`
	Engagement *Engagement `json:"engagement,omitempty"`
	SourceURL  string      `json:"source_url"`
	Platform   string      `json:"platform"`
}

// Comment represents a single comment/reply.
type Comment struct {
	Author     string      `json:"author"`
	Content    string      `json:"content"`
	Engagement *Engagement `json:"engagement,omitempty"`
	IsAuthor   bool        `json:"is_author,omitempty"`
	Children   []*Comment  `json:"children,omitempty"`

	// parentIndex is used internally for tree building.
	// -1 means root-level comment, >= 0 means index of parent in flat list.
	parentIndex int
}

// Engagement holds interaction metrics. Pointer fields: nil = not available.
type Engagement struct {
	Likes   *int `json:"likes,omitempty"`
	Replies *int `json:"replies,omitempty"`
	Reposts *int `json:"reposts,omitempty"`
	Quotes  *int `json:"quotes,omitempty"`
}

// Result wraps the full output for JSON serialization.
type Result struct {
	Post     *Post      `json:"post"`
	Comments []*Comment `json:"comments,omitempty"`
}

// intPtr is a helper to create a pointer to an int.
func intPtr(n int) *int {
	return &n
}
