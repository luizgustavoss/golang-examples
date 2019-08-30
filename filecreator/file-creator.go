package filecreator

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/luizgustavoss/golang-examples/filewriter"
	"github.com/luizgustavoss/golang-examples/seedconfigurator"
)

var randomFilesQuantity = 1000
var baseDir string

func generateFileName(fileIndex int) string {
	filePath := fmt.Sprintf("%s/file%d.%s", baseDir, fileIndex, "txt")
	return filePath
}

func createFile(filePath string) *os.File {
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Error creating file %s: %v\n", filePath, err)
		os.Exit(1)
	}
	return file
}

func CreateRandomFiles(dir string) {
	baseDir = dir
	seedconfigurator.ConfigureRandomSeed()
	numberOfFiles := rand.Intn(randomFilesQuantity)
	for fileIndex := 0; fileIndex < numberOfFiles; fileIndex++ {
		filePath := generateFileName(fileIndex)
		file := createFile(filePath)
		filewriter.WriteRandomValueInFile(file)
		defer file.Close()
	}
}
