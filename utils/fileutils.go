package utils

import (
	"io"
	"net/http"
	"os"
	"strings"
)

func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func DetectContentType(filePath string) (string, error) {
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

func CopyFileContents(src, dst string) (err error) {
	if dirExists, err := Exists(dst); !dirExists && err == nil {
		os.Mkdir(dst, os.ModePerm)
	}

	filename := src[strings.LastIndex(src, "/")+1 : len(src)]

	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	destination := dst + "/" + filename

	out, err := os.OpenFile(destination, os.O_RDWR|os.O_CREATE, 0666)
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
