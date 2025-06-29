package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Регистрируем обработчики
	http.HandleFunc("/api/search", searchHandler)
	http.HandleFunc("/", homeHandler)

	// Запускаем сервер
	port := ":8080"
	fmt.Printf("Сервер запущен на порту %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

// Пустые заглушки обработчиков, реализация будет в main.go
func homeHandler(w http.ResponseWriter, r *http.Request) {
	HomeHandler(w, r)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	SearchHandler(w, r)
}
