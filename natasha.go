package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"natasha/constants"
	"net/http"
	"os"
	"runtime"
	"strings"
)

func main() {
	fmt.Printf("Natasha welcomes you to a %v system\n", runtime.GOOS)
	fmt.Println()
	fmt.Println()

	//"/Users/DSharma/Downloads"
	path, outputPath := getRootPath()
	fileType := getFileTypeToSearch(path)
	searchFileSystem(path, outputPath, fileType)
}

func searchFileSystem(startingPoint, outputPath, filetype string) {
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
			searchFileSystem(startingPoint+"/"+file.Name(), outputPath, filetype)
		}
	} else {
		contentType, _ := detectContentType(startingPoint)
		switch filetype {
		case constants.FILE_TYPES[0]:
			{
				if strings.Contains(contentType, constants.FILE_TYPES_MAP[filetype]) {
					fmt.Println(startingPoint + " : " + contentType)
					copyFileContents(startingPoint, outputPath)
				}
			}
		case constants.FILE_TYPES[1]:
			{
				if strings.Contains(contentType, constants.FILE_TYPES_MAP[filetype]) {
					copyFileContents(startingPoint, outputPath)
				}
			}
		default:
			{
				fmt.Println("Invalid input")
			}
		}

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

func getRootPath() (string, string) {
	fmt.Printf("Please input root location to initiate search. Press enter to use %v as root location\n", constants.LINUX_ROOT)
	fmt.Println()
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	path := scanner.Text()

	if path == "" {
		path = constants.LINUX_ROOT
	}

	fmt.Printf("Please input output location to copy matching files. Press enter to use %v as default location\n", constants.DEFAULT_LINUX_OUTPUT_DIR)
	fmt.Println()
	scanner.Scan()
	outputPath := scanner.Text()

	if outputPath == "" {
		outputPath = constants.DEFAULT_LINUX_OUTPUT_DIR
	}

	return path, outputPath
}

func getFileTypeToSearch(searchPath string) string {
	fmt.Printf("Please input index for filetype to search. For ex - 1 for %v \n", constants.FILE_TYPES[0])
	fmt.Println()
	for i, fileType := range constants.FILE_TYPES {
		fmt.Printf("%d. %v\n", i+1, fileType)
	}
	fmt.Println()
	var i int
	fmt.Scanf("%d", &i)

	retryCounter := 3
	for i < 1 && retryCounter > 0 {
		fmt.Printf("No input detected. Please try again. %d tries remaining\n", retryCounter)
		retryCounter--
		fmt.Scanf("%d", &i)
	}
	if i < 1 {
		i = 1
	}

	fmt.Printf("Scanning for %v files in %v\n", constants.FILE_TYPES[i-1], searchPath)
	return constants.FILE_TYPES[i-1]

}

func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Open(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}
