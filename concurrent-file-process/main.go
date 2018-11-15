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

var syncControl sync.WaitGroup
var randomFilesQuantity = 100
var maxSleepTime = 10
var chunkFileSize = 5
var baseDir = os.TempDir() + "/poc"

func configureRandomSeed() {
	now := time.Now()
	rand.Seed(now.UnixNano())
}

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

func writeRandomValueInFile(file *os.File) {
	file.WriteString(strconv.Itoa(rand.Intn(maxSleepTime)))
}

func createRandomFiles() {
	configureRandomSeed()
	numberOfFiles := rand.Intn(randomFilesQuantity)
	for fileIndex := 0; fileIndex < numberOfFiles; fileIndex++ {
		filePath := generateFileName(fileIndex)
		file := createFile(filePath)
		writeRandomValueInFile(file)
		defer file.Close()
	}
}

func processFiles() []chan string {
	files := readRandomFiles()
	if len(files) < 1 {
		fmt.Printf("There are no files to process! \n\n")
		return make([]chan string, 0)
	}
	return processFilesAsynchronously(files)
}

func printFileSetInfo(files []string) {
	fmt.Printf("Files to Process: \n\n")
	for _, name := range files {
		fmt.Printf("File: %s \n", name)
	}
}

func processFilesAsynchronously(files []string) []chan string {

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
		copy(files, append(files[chunkSize:]))

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
		syncControl.Add(1)
		fmt.Printf("Process: %s \n", filesChunk[i])
		go processFile(&syncControl, chans[i], filesChunk[i])
	}
	syncControl.Wait()
	return chans
}

func readFileContentAsNumber(fileName string) int {
	fileData, _ := ioutil.ReadFile(fileName)
	data, _ := strconv.Atoi(string(fileData))
	fmt.Printf("File: %s | Content: %d\n", fileName, data)
	return data
}

func processFile(syncControl *sync.WaitGroup, c chan<- string, fileName string) {
	defer syncControl.Done()
	seconds := readFileContentAsNumber(fileName)
	if seconds > 0 {
		time.Sleep(time.Duration(seconds) * time.Second)
	}
	c <- fmt.Sprintf("Process File: %s | timeToSleep: %d", fileName, seconds)
}

func readRandomFiles() []string {
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

func main() {

	createRandomFiles()
	channels := processFiles()

	channlesLength := len(channels)
	for i := 0; i < channlesLength; i++ {
		c := channels[i]
		msg := <-c
		fmt.Printf("Message: [%s]\n", msg)
	}
}
