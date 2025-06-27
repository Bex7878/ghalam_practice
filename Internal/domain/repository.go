package domain

type FileRepository interface {
	SearchInFile(path string, req SearchRequest) (SearchResult, error)
	GetFileContent(path string) (FileContent, error)
	ListFiles(dirPath string, recursive bool, extensions []string) ([]string, error)
}
