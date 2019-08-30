package fileprocessor

import (
	"fmt"
	"sync"
	"time"

	"github.com/luizgustavoss/golang-examples/filereader"
)

var chunkFileSize = 100
var syncControl sync.WaitGroup

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

func processFile(syncControl *sync.WaitGroup, c chan<- string, fileName string) {
	defer syncControl.Done()
	seconds := filereader.ReadFileContentAsNumber(fileName)
	if seconds > 0 {
		time.Sleep(time.Duration(seconds) * time.Second)
	}
	c <- fmt.Sprintf("Process File: %s | timeToSleep: %d", fileName, seconds)
}

func processFilesAsynchronouslyInBatches(files []string) []chan string {

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

func printFileSetInfo(files []string) {
	fmt.Printf("Files to Process: \n\n")
	for _, name := range files {
		fmt.Printf("File: %s \n", name)
	}
}

func processFilesAsynchronously(files []string) []chan string {

	numberOfPenddingFiles := len(files)
	printFileSetInfo(files)

	channels := make([]chan string, numberOfPenddingFiles)

	for i := range channels {
		channels[i] = make(chan string, 1)
	}
	for i := 0; i < numberOfPenddingFiles; i++ {
		syncControl.Add(1)
		fmt.Printf("Process: %s \n", files[i])
		go processFile(&syncControl, channels[i], files[i])
	}
	syncControl.Wait()

	return channels
}

func ProcessFiles(baseDir string) []chan string {
	files := filereader.ReadRandomFiles(baseDir)
	if len(files) < 1 {
		fmt.Printf("There are no files to process! \n\n")
		return make([]chan string, 0)
	}
	return processFilesAsynchronously(files)
}
