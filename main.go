package main

import (
	"archive/zip"
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// Парсинг аргументов командной строки
	searchTerm := flag.String("s", "", "Строка для поиска")
	ignoreCase := flag.Bool("i", false, "Игнорировать регистр при поиске")
	showLineNumbers := flag.Bool("n", false, "Показывать номера строк")
	recursive := flag.Bool("r", false, "Рекурсивный поиск в поддиректориях")
	showMatchesOnly := flag.Bool("o", false, "Показывать только файлы с совпадениями")
	fileExtensions := flag.String("ext", "", "Фильтр по расширениям файлов (через запятую, например '.html,.txt,.zip')")
	flag.Parse()

	if len(flag.Args()) < 1 {
		fmt.Println("Использование: filesearch [опции] <путь_к_файлу_или_директории>")
		fmt.Println("Опции:")
		flag.PrintDefaults()
		return
	}

	path := flag.Args()[0]

	// Получаем список расширений для фильтрации
	extList := strings.Split(strings.ToLower(*fileExtensions), ",")

	// Проверка, является ли путь директорией
	fileInfo, err := os.Stat(path)
	if err != nil {
		log.Fatalf("Ошибка: %v", err)
	}

	if fileInfo.IsDir() {
		// Поиск в директории
		searchInDirectory(path, *searchTerm, *ignoreCase, *showLineNumbers, *recursive, *showMatchesOnly, extList)
	} else {
		// Проверка расширения файла
		if !hasValidExtension(path, extList) {
			fmt.Printf("Файл %s не соответствует указанным расширениям\n", path)
			return
		}

		// Поиск в одиночном файле
		if *searchTerm != "" {
			fmt.Printf("=== Поиск в файле: %s ===\n", path)
			found := searchInFile(path, *searchTerm, *ignoreCase, *showLineNumbers, *showMatchesOnly)
			if !found && !*showMatchesOnly {
				fmt.Println("Совпадений не найдено.")
			}
		} else {
			// Просто показать содержимое файла
			printFileContent(path)
		}
	}
}

func hasValidExtension(filePath string, extList []string) bool {
	if len(extList) == 0 || (len(extList) == 1 && extList[0] == "") {
		return true
	}

	ext := strings.ToLower(filepath.Ext(filePath))
	for _, e := range extList {
		if ext == strings.TrimSpace(e) {
			return true
		}
	}
	return false
}

func searchInDirectory(dirPath, searchTerm string, ignoreCase, showLineNumbers, recursive, showMatchesOnly bool, extList []string) {
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Пропускаем директории, если не рекурсивный поиск
		if info.IsDir() && path != dirPath && !recursive {
			return filepath.SkipDir
		}

		// Обрабатываем только обычные файлы с нужным расширением
		if !info.IsDir() && info.Mode().IsRegular() && hasValidExtension(path, extList) {
			if searchTerm != "" {
				if showMatchesOnly {
					// Быстрая проверка на наличие строки в файле
					if fileContains(path, searchTerm, ignoreCase) {
						fmt.Printf("\n=== Найдено в файле: %s ===\n", path)
						searchInFile(path, searchTerm, ignoreCase, showLineNumbers, showMatchesOnly)
					}
				} else {
					fmt.Printf("\n=== Поиск в файле: %s ===\n", path)
					searchInFile(path, searchTerm, ignoreCase, showLineNumbers, showMatchesOnly)
				}
			} else {
				// Просто показать содержимое файла
				fmt.Printf("\n=== Содержимое файла: %s ===\n", path)
				printFileContent(path)
			}
		}

		return nil
	})

	if err != nil {
		log.Fatalf("Ошибка при обходе директории: %v", err)
	}
}

func fileContains(filePath, searchTerm string, ignoreCase bool) bool {
	// Для ZIP-файлов проверяем все файлы внутри
	if strings.HasSuffix(strings.ToLower(filePath), ".zip") {
		return zipContains(filePath, searchTerm, ignoreCase)
	}

	// Для обычных файлов
	file, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	searchTermLower := searchTerm
	if ignoreCase {
		searchTermLower = strings.ToLower(searchTerm)
	}

	for scanner.Scan() {
		line := scanner.Text()
		if ignoreCase {
			line = strings.ToLower(line)
		}
		if strings.Contains(line, searchTermLower) {
			return true
		}
	}

	return false
}

func zipContains(zipPath, searchTerm string, ignoreCase bool) bool {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return false
	}
	defer r.Close()

	searchTermLower := searchTerm
	if ignoreCase {
		searchTermLower = strings.ToLower(searchTerm)
	}

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			continue
		}

		scanner := bufio.NewScanner(rc)
		for scanner.Scan() {
			line := scanner.Text()
			if ignoreCase {
				line = strings.ToLower(line)
			}
			if strings.Contains(line, searchTermLower) {
				rc.Close()
				return true
			}
		}
		rc.Close()
	}

	return false
}

func searchInFile(filePath, searchTerm string, ignoreCase, showLineNumbers, showMatchesOnly bool) bool {
	// Обработка ZIP-архивов
	if strings.HasSuffix(strings.ToLower(filePath), ".zip") {
		return searchInZipFile(filePath, searchTerm, ignoreCase, showLineNumbers, showMatchesOnly)
	}

	// Обработка обычных файлов
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Ошибка открытия файла %s: %v", filePath, err)
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 1
	found := false

	searchTermLower := searchTerm
	if ignoreCase {
		searchTermLower = strings.ToLower(searchTerm)
	}

	for scanner.Scan() {
		line := scanner.Text()
		lineToCompare := line
		if ignoreCase {
			lineToCompare = strings.ToLower(line)
		}

		if strings.Contains(lineToCompare, searchTermLower) {
			found = true
			if showLineNumbers {
				fmt.Printf("%d: ", lineNumber)
			}
			fmt.Println(line)
		}
		lineNumber++
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Ошибка чтения файла %s: %v", filePath, err)
	}

	return found
}

func searchInZipFile(zipPath, searchTerm string, ignoreCase, showLineNumbers, showMatchesOnly bool) bool {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		log.Printf("Ошибка открытия ZIP-архива %s: %v", zipPath, err)
		return false
	}
	defer r.Close()

	searchTermLower := searchTerm
	if ignoreCase {
		searchTermLower = strings.ToLower(searchTerm)
	}

	anyFound := false

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			log.Printf("Ошибка открытия файла %s в архиве %s: %v", f.Name, zipPath, err)
			continue
		}

		scanner := bufio.NewScanner(rc)
		lineNumber := 1
		fileFound := false

		for scanner.Scan() {
			line := scanner.Text()
			lineToCompare := line
			if ignoreCase {
				lineToCompare = strings.ToLower(line)
			}

			if strings.Contains(lineToCompare, searchTermLower) {
				if !fileFound {
					fmt.Printf("\n=== Найдено в файле внутри архива: %s/%s ===\n", zipPath, f.Name)
					fileFound = true
					anyFound = true
				}
				if showLineNumbers {
					fmt.Printf("%d: ", lineNumber)
				}
				fmt.Println(line)
			}
			lineNumber++
		}

		rc.Close()
	}

	return anyFound
}

func printFileContent(filePath string) {
	// Обработка ZIP-архивов
	if strings.HasSuffix(strings.ToLower(filePath), ".zip") {
		printZipContent(filePath)
		return
	}

	// Обработка обычных файлов
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("Ошибка чтения файла %s: %v", filePath, err)
		return
	}

	if isText(content) {
		fmt.Println(string(content))
	} else {
		fmt.Println("Файл не является текстовым (бинарный файл).")
	}
}

func printZipContent(zipPath string) {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		log.Printf("Ошибка открытия ZIP-архива %s: %v", zipPath, err)
		return
	}
	defer r.Close()

	fmt.Printf("Содержимое ZIP-архива: %s\n", zipPath)
	for _, f := range r.File {
		fmt.Printf("  %s (размер: %d байт)\n", f.Name, f.UncompressedSize64)
	}
}

func isText(content []byte) bool {
	if len(content) == 0 {
		return true
	}

	nonTextChars := 0
	for _, b := range content {
		if b < 32 && b != '\t' && b != '\n' && b != '\r' && b != '\f' {
			nonTextChars++
			if float64(nonTextChars)/float64(len(content)) > 0.05 {
				return false
			}
		}
	}
	return true
}
