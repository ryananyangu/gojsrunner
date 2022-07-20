package utils

import (
	"io/ioutil"
	"os"
	"time"
)

func ReadFile(file string) (string, error) {
	openedFile, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer openedFile.Close()

	byteValue, _ := ioutil.ReadAll(openedFile)

	return string(byteValue[:]), nil

}

func TimeTrack(start time.Time, name string) (string, time.Duration) {
	elapsed := time.Since(start)
	return name, elapsed
}
