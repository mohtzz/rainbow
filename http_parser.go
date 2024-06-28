package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// /home/arseny/rbs_studying
func main() {
	inputFilePath := flag.String("src", "", "input")
	resultDir := flag.String("dst", "result", "output")
	flag.Parse()

	// проверка на корректность флагов

	if *inputFilePath == "" || *resultDir == "" {
		fmt.Println("wrong path/path doesn't exist")
		return
	}
	// создание директории при ее отсутствии
	if _, err := os.Stat(*resultDir); os.IsNotExist(err) {
		os.MkdirAll(*resultDir, 0755)
	}

	file, err := os.Open(*inputFilePath)
	if err != nil {
		fmt.Println("error by opening a file:", err)
		return
	}
	defer file.Close()

	var sites []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			fmt.Println("error by reading file: ", err)
			return
		}
		sites = append(sites, scanner.Text())
	}

	for _, site := range sites {
		resp, err := http.Get(site)
		if err != nil {
			fmt.Println("error by sending response: ", err)
			continue
		}
		defer resp.Body.Close() // проблемный дефер

		if resp.StatusCode == http.StatusOK {
			filename := filepath.Join(*resultDir, strings.ReplaceAll(site, "/", "_")+".html")
			file, err := os.Create(filename)
			if err != nil {
				fmt.Println("error by creating file: ", err)
				continue
			}
			defer file.Close() // проблемный дефер

			fmt.Println("Файл", filename, "сохранен.")

		} else {
			fmt.Println("Сайт", site, "вернул статус-код: ", resp.StatusCode)
		}
	}
}

// func openFile(inputFilePath *string){
// 	file, err := os.Open(*inputFilePath)
// 	if err != nil {
// 		fmt.Println("error by opening a file:", err)
// 		return
// 	}
// 	defer file.Close()
// }
