package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// /home/arseny/rbs_studying/rainbow
func main() {
	start := time.Now()
	defer func() {
		progTime := time.Since(start)
		fmt.Println("Время выполнения программы: ", progTime)
	}()

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

	sites, err := readSitesFromFile(*inputFilePath)
	if err != nil {
		fmt.Println("error by reading file: ", err)
		return
	}

	for _, site := range sites {
		err := processSite(site, *resultDir)
		if err != nil {
			fmt.Println("error by processing site: ", err)
		}
	}
}

func readSitesFromFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var sites []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, err
		}
		sites = append(sites, scanner.Text())
	}
	return sites, nil
}

func processSite(site string, resultDir string) error {
	resp, err := http.Get(site)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		filename := filepath.Join(resultDir, strings.ReplaceAll(site, "/", "_")+".html")
		return saveHTML(filename, resp.Body)
	} else {
		fmt.Println("Сайт", site, "вернул статус-код: ", resp.StatusCode)
		return nil
	}
}

func saveHTML(filename string, body io.Reader) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, body)
	if err != nil {
		return err
	}
	fmt.Println("Файл", filename, "сохранен.")
	return nil
}
