package filereader

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

var baseDir string

func ReadRandomFiles(dir string) []string {
	baseDir = dir
	var files []string
	err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		panic(err)
	}
	return append(files[1:]) //ignore first item (dir)
}

func ReadFileContentAsNumber(fileName string) int {
	fileData, _ := ioutil.ReadFile(fileName)
	data, _ := strconv.Atoi(string(fileData))
	fmt.Printf("File: %s | Content: %d\n", fileName, data)
	return data
}
