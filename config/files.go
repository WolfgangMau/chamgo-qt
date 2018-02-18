package config

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func GetFilesInFolder(root string, ext string) []string {
	var files []string
	log.Printf("looking for files with extension %s in %s\n", root, ext)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		log.Printf("path: %s\n", path)
		if filepath.Ext(path) == ext {
			log.Printf("add %s\n", path)
			files = append(files, strings.Replace(path, root, "", 1))
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		log.Println(file)
	}
	return files
}

func ReadFileLines(path string) (res []string) {
	file, err := os.Open(path)
	if err != nil {
		log.Print(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		res = append(res, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Print(err)
	}
	return res
}
