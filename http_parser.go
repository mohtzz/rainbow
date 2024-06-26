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
	inputFilePath := flag.String("src", "sites.txt", "input")
	resultDir := flag.String("dst", "result", "output")
	flag.Parse()

	file, err := os.Open(*inputFilePath)
	if err != nil {
		fmt.Println("error by opening a file:", err)
		return
	}
	defer file.Close()

	var sites []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		sites = append(sites, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("error by reading file: ", err)
		return
	}

	for _, site := range sites {
		resp, err := http.Get(site)
		if err != nil {
			fmt.Println("error by sending response: ", err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			filename := filepath.Join(*resultDir, strings.ReplaceAll(site, "/", "_")+".html")
			file, err := os.Create(filename)
			if err != nil {
				fmt.Println("error by creating file: ", err)
				continue
			}
			defer file.Close()

			fmt.Println("Файл", filename, "сохранен.")

		} else {
			fmt.Println("Сайт", site, "вернул статус-код: ", resp.StatusCode)
		}
	}
}
