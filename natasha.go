package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
)

func main() {
	fmt.Printf("Natasha welcomes you to a %v system\n", runtime.GOOS)
	fmt.Println("Using / as root")

	searchFileSystem("/Users/DSharma/Downloads")
}

func searchFileSystem(startingPoint string) {
	f, err := os.Stat(startingPoint)
	if err != nil {
		log.Fatal(err)
	}

	if f.IsDir() {
		files, err := ioutil.ReadDir(startingPoint)
		if err != nil {
			log.Fatal(err)
		}

		for _, file := range files {
			searchFileSystem(startingPoint + "/" + file.Name())
		}
	} else {

		contentType, _ := detectContentType(startingPoint)

		fmt.Println(startingPoint + " : " + contentType)
	}

}

func detectContentType(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}

	buffer := make([]byte, 512)
	_, err2 := file.Read(buffer)
	if err2 != nil {
		return "", err2
	}

	contentType := http.DetectContentType(buffer)
	return contentType, nil
}
