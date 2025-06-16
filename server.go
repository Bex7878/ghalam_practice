package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type DataRecord struct {
	Region string    `json:"region"`
	Data   float64   `json:"data"`
	Date   time.Time `json:"date"`
}

func main() {
	// Настройка обработчиков
	http.HandleFunc("/api/data", dataHandler)
	http.HandleFunc("/", homeHandler)

	// Запуск сервера
	port := ":8080"
	fmt.Printf("Сервер запущен на порту %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Отдаем статические файлы
	path := r.URL.Path
	if path == "/" {
		path = "/home.html"
	}

	// Получаем абсолютный путь к файлу
	absPath, err := filepath.Abs("." + path)
	if err != nil {
		http.Error(w, "Ошибка пути", http.StatusInternalServerError)
		return
	}

	// Проверяем существование файла
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}

	http.ServeFile(w, r, absPath)
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
	// Парсим параметры запроса
	region := r.URL.Query().Get("region")
	dateFrom := r.URL.Query().Get("date_from")
	dateTo := r.URL.Query().Get("date_to")

	// Здесь должна быть логика получения данных из БД или другого источника
	// Для примера используем фиктивные данные
	records := generateSampleData(region, dateFrom, dateTo)

	// Устанавливаем заголовок Content-Type
	w.Header().Set("Content-Type", "application/json")

	// Кодируем данные в JSON и отправляем
	if err := json.NewEncoder(w).Encode(records); err != nil {
		http.Error(w, "Ошибка кодирования JSON", http.StatusInternalServerError)
		return
	}
}

func generateSampleData(region, dateFrom, dateTo string) []DataRecord {
	// Генерируем фиктивные данные для демонстрации
	var records []DataRecord
	regions := []string{
		"abay", "akmola", "aktobe", "almaty", "atyrau", "east-kazakhstan",
		"zhambyl", "zhetysu", "west-kazakhstan", "karaganda", "kostanay",
		"kzylorda", "mangystau", "pavlodar", "north-kazakhstan", "turkistan",
		"ulytau", "almaty-city", "astana", "shymkent",
	}

	now := time.Now()
	for i, r := range regions {
		// Фильтрация по региону, если указан
		if region != "" && region != "all" && region != r {
			continue
		}

		// Добавляем несколько записей для каждого региона
		for j := 0; j < 3; j++ {
			date := now.AddDate(0, 0, -j*10)
			records = append(records, DataRecord{
				Region: r,
				Data:   float64((i+1)*100 + j*10),
				Date:   date,
			})
		}
	}

	return records
}
