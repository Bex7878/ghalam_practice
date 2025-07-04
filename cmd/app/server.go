package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	exeDir, err := os.Getwd() // Текущая директория исполняемого файла
	if err != nil {
		log.Fatal(err)
	}
	htmlPath := filepath.Join(exeDir, "frontend", "index.html")

	http.ServeFile(w, r, htmlPath)
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	files := []map[string]interface{}{
		{
			"Filename": "test1.txt",
			"Region":   "Almaty",
			"Path":     "/files/test1.txt",
			"Lat":      43.238949,
			"Lon":      76.889709,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(files)
	if err != nil {
		return
	}
}

func main() {

	http.HandleFunc("/api/search", enableCORS(SearchHandler))

	fs := http.FileServer(http.Dir("frontend"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path == "/" {
			HomeHandler(w, r)
			return
		}

		fs.ServeHTTP(w, r)
	})

	port := ":8080"
	fmt.Printf("Сервер запущен на порту %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			return
		}
		next(w, r)
	}
}
