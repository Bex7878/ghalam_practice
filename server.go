package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	searchResults []FileSearchResult
	mu            sync.Mutex
)

// Координаты регионов (примерные)
var regionCoordinates = map[string]struct {
	MinLat, MaxLat float64
	MinLon, MaxLon float64
}{
	"abay":   {48.0, 50.0, 72.0, 76.0},
	"akmola": {50.0, 53.0, 66.0, 72.0},
	// Добавьте координаты для остальных регионов...
}

type FileSearchResult struct {
	Path     string  `json:"path"`
	Filename string  `json:"filename"`
	Content  string  `json:"content,omitempty"`
	Region   string  `json:"region"`
	Lat      float64 `json:"lat"`
	Lon      float64 `json:"lon"`
}

func main() {
	// Настройка обработчиков
	http.HandleFunc("/api/search", searchHandler)
	http.HandleFunc("/api/files", filesHandler)
	http.HandleFunc("/", homeHandler)

	// Запуск сервера
	port := ":8080"
	fmt.Printf("Сервер запущен на порту %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path == "/" {
		path = "/home.html"
	}

	absPath, err := filepath.Abs("." + path)
	if err != nil {
		http.Error(w, "Ошибка пути", http.StatusInternalServerError)
		return
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}

	http.ServeFile(w, r, absPath)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	searchTerm := r.URL.Query().Get("q")
	region := r.URL.Query().Get("region")
	ext := r.URL.Query().Get("ext")

	if searchTerm == "" {
		http.Error(w, "Не указан поисковый запрос", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	// Очищаем предыдущие результаты
	searchResults = nil

	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			// Проверяем расширение файла
			if ext != "" && filepath.Ext(path) != ext {
				return nil
			}

			// Проверяем содержимое файла
			content, err := os.ReadFile(path)
			if err != nil {
				return nil
			}

			if strings.Contains(string(content), searchTerm) {
				// Определяем регион файла (в реальной системе это должно быть из метаданных файла)
				fileRegion := determineFileRegion(path)

				// Если указан регион и файл не принадлежит ему - пропускаем
				if region != "" && region != "all" && fileRegion != region {
					return nil
				}

				// Генерируем случайные координаты в пределах региона (для демонстрации)
				var lat, lon float64
				if coords, ok := regionCoordinates[fileRegion]; ok {
					lat = coords.MinLat + (coords.MaxLat-coords.MinLat)*0.5
					lon = coords.MinLon + (coords.MaxLon-coords.MinLon)*0.5
				}

				searchResults = append(searchResults, FileSearchResult{
					Path:     path,
					Filename: filepath.Base(path),
					Content:  string(content),
					Region:   fileRegion,
					Lat:      lat,
					Lon:      lon,
				})
			}
		}
		return nil
	})

	if err != nil {
		http.Error(w, "Ошибка поиска", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(searchResults)
}

func filesHandler(w http.ResponseWriter, r *http.Request) {
	region := r.URL.Query().Get("region")

	mu.Lock()
	defer mu.Unlock()

	var results []FileSearchResult
	for _, res := range searchResults {
		if region == "" || region == "all" || res.Region == region {
			results = append(results, res)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

// Функция для определения региона файла (заглушка)
func determineFileRegion(path string) string {
	// В реальной системе это должно определяться из метаданных файла
	// или по его расположению в файловой системе
	// Здесь просто возвращаем случайный регион для демонстрации
	regions := []string{
		"abay", "akmola", "aktobe", "almaty", "atyrau", "east-kazakhstan",
		"zhambyl", "zhetysu", "west-kazakhstan", "karaganda", "kostanay",
		"kzylorda", "mangystau", "pavlodar", "north-kazakhstan", "turkistan",
		"ulytau", "almaty-city", "astana", "shymkent",
	}
	return regions[len(path)%len(regions)]
}
