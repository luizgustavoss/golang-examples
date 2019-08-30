package filewriter

import (
	"math/rand"
	"os"
	"strconv"
)

var maxSleepTime = 20

func WriteRandomValueInFile(file *os.File) {
	file.WriteString(strconv.Itoa(rand.Intn(maxSleepTime)))
}
