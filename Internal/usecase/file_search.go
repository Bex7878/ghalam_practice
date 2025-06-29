package usecases

import (
	"github.com/Bex7878/ghalam_practice/Internal/domain"
	"strings"
)

type FileSearchUseCase struct {
	fileRepo   domain.FileRepository
	extractors map[string]domain.DataExtractor
}

func NewFileSearchUseCase(repo domain.FileRepository, extractors map[string]domain.DataExtractor) *FileSearchUseCase {
	return &FileSearchUseCase{
		fileRepo:   repo,
		extractors: extractors,
	}
}

func (uc *FileSearchUseCase) SearchWithDynamicValues(req domain.SearchRequest) ([]domain.SearchResult, error) {
	// Если поисковый термин начинается с @, извлекаем значение
	if strings.HasPrefix(req.SearchTerm, "@") {
		content, err := uc.fileRepo.GetFileContent(req.Path)
		if err != nil {
			return nil, err
		}

		extractorKey := strings.TrimPrefix(req.SearchTerm, "@")
		if extractor, exists := uc.extractors[extractorKey]; exists {
			values, err := extractor.ExtractValues(content.Content)
			if err != nil {
				return nil, err
			}

			// Используем все найденные значения для поиска
			var allResults []domain.SearchResult
			for _, value := range values {
				newReq := req
				newReq.SearchTerm = value
				results, err := uc.Search(newReq)
				if err != nil {
					return nil, err
				}
				allResults = append(allResults, results...)
			}
			return allResults, nil
		}
	}

	return uc.Search(req)
}
