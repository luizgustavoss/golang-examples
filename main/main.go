package main

import (
	"fmt"
	"os"

	"github.com/luizgustavoss/golang-examples/filecreator"
	"github.com/luizgustavoss/golang-examples/fileprocessor"
)

var baseDir = os.TempDir() + "/poc"

func main() {

	filecreator.CreateRandomFiles(baseDir)
	channels := fileprocessor.ProcessFiles(baseDir)

	channlesLength := len(channels)
	for i := 0; i < channlesLength; i++ {
		c := channels[i]
		msg := <-c
		fmt.Printf("Message: [%s]\n", msg)
	}
}
