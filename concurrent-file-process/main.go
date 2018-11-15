package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

var controle sync.WaitGroup

func configureRandomSeed() {
	now := time.Now()
	rand.Seed(now.UnixNano())
}

func generateFileName(baseDir string, fileIndex int) string {
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

func writeRandomValueInFile(file *os.File) {
	file.WriteString(strconv.Itoa(rand.Intn(5)))
}

func createRandomFiles(baseDir string) {
	configureRandomSeed()
	numberOfFiles := rand.Intn(50)
	for fileIndex := 0; fileIndex < numberOfFiles; fileIndex++ {
		filePath := generateFileName(baseDir, fileIndex)
		file := createFile(filePath)
		writeRandomValueInFile(file)
		defer file.Close()
	}
}

func processFiles(baseDir string) []chan string {
	files := readRandomFiles(baseDir)
	return processFilesAsynchronously(append(files[1:])) // remover o primeiro item que é o próprio diretório
}

func printFileSetInfo(files []string) {
	fmt.Printf("Files to Process: \n\n")
	for _, name := range files {
		fmt.Printf("File: %s \n", name)
	}
}

func processFilesAsynchronously(files []string) []chan string {

	var chunkFileSize = 5
	numberOfPenddingFiles := len(files)
	printFileSetInfo(files)

	channels := make([]chan string, 0)

	var chunkSize int
	var lastInteraction = false

	for !lastInteraction {
		fmt.Printf("numberOfPenddingFiles: %d \n", numberOfPenddingFiles)

		if numberOfPenddingFiles > chunkFileSize {
			numberOfPenddingFiles -= chunkFileSize
			chunkSize = chunkFileSize
		} else {
			chunkSize = numberOfPenddingFiles
			numberOfPenddingFiles = 0
			lastInteraction = true
		}

		fmt.Printf("chunkSize: %d \n", chunkSize)

		filesChunk := make([]string, chunkSize)
		copy(filesChunk, files[0:chunkSize])
		copy(files, append(files[chunkSize+1:]))

		chans := processChunk(filesChunk)
		channels = append(channels, chans...)
	}
	return channels
}

func processChunk(filesChunk []string) []chan string {

	var chunkSize = len(filesChunk)

	chans := make([]chan string, chunkSize)
	for i := range chans {
		chans[i] = make(chan string, 1)
	}
	for i := 0; i < chunkSize; i++ {
		controle.Add(1)
		fmt.Printf("Process: %s \n", filesChunk[i])
		go processFile(&controle, chans[i], filesChunk[i])
	}
	controle.Wait()
	return chans
}

func readFileContentAsNumber(fileName string) int {
	fileData, _ := ioutil.ReadFile(fileName)
	data, _ := strconv.Atoi(string(fileData))
	fmt.Printf("File: %s | Content: %d\n", fileName, data)
	return data
}

func processFile(controle *sync.WaitGroup, c chan<- string, fileName string) {
	defer controle.Done()
	seconds := readFileContentAsNumber(fileName)
	if seconds > 0 {
		time.Sleep(time.Duration(seconds) * time.Second)
	}
	c <- fmt.Sprintf("Process File: %s | timeToSleep: %d", fileName, seconds)
}

func readRandomFiles(baseDir string) []string {
	var files []string
	err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		panic(err)
	}
	return files
}

func main() {

	baseDir := os.TempDir() + "/poc"
	createRandomFiles(baseDir)
	channels := processFiles(baseDir)

	channlesLength := len(channels)
	for i := 0; i < channlesLength; i++ {
		c := channels[i]
		msg := <-c
		fmt.Printf("Message: [%s]\n", msg)
	}
}
