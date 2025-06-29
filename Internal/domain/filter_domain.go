package domain

type SearchRequest struct {
	Path          string
	SearchTerm    string
	Extensions    []string
	IgnoreCase    bool
	ShowLineNums  bool
	Recursive     bool
	OnlyShowMatch bool
}

type SearchResult struct {
	FilePath string
	Matches  []Match
	Error    string
}

type Match struct {
	LineNumber int
	Content    string
}

type FileContent struct {
	Name     string
	Content  string
	IsBinary bool
}
