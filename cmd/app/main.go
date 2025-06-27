package main

import (
	"github.com/Bex7878/ghalam_practice/Internal/domain"
	usecases "github.com/Bex7878/ghalam_practice/Internal/usecase"
	"log"
	"net/http"
	"strings"
)

func main() {
	// Инициализация репозиториев
	fileRepo := file_repository.NewLocalFileRepository()
	zipRepo := file_repository.NewZipFileRepository()

	// Агрегированный репозиторий
	repo := NewAggregateFileRepository(fileRepo, zipRepo)

	// Инициализация use case
	searchUC := usecases.NewFileSearchUseCase(repo)

	// Инициализация обработчиков
	handler := api.NewFileSearchHandler(searchUC)

	// Настройка маршрутов
	http.HandleFunc("/api/search", handler.Search)
	http.HandleFunc("/api/content", handler.GetContent)

	// Запуск сервера
	log.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// AggregateFileRepository объединяет несколько репозиториев
type AggregateFileRepository struct {
	repos map[string]domain.FileRepository
}

func NewAggregateFileRepository(repos ...domain.FileRepository) *AggregateFileRepository {
	a := &AggregateFileRepository{
		repos: make(map[string]domain.FileRepository),
	}

	for _, repo := range repos {
		switch r := repo.(type) {
		case *file_repository.LocalFileRepository:
			a.repos["local"] = r
		case *file_repository.ZipFileRepository:
			a.repos["zip"] = r
		}
	}

	return a
}

func (a *AggregateFileRepository) SearchInFile(path string, req domain.SearchRequest) (domain.SearchResult, error) {
	// Выбираем нужный репозиторий в зависимости от расширения файла
	if strings.HasSuffix(strings.ToLower(path), ".zip") {
		return a.repos["zip"].SearchInFile(path, req)
	}
	return a.repos["local"].SearchInFile(path, req)
}

// Остальные методы агрегированного репозитория...
