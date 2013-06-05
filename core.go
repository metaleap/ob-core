package core

import (
	"fmt"
	"os"
	"path/filepath"

	uio "github.com/metaleap/go-util/io"
)

var (
	DirPath string
)

func isObDir(dirPath string) bool {
	return uio.DirsFilesExist(dirPath, "client")
}

func GuessDirPath(userSpecified string) (guess string) {
	if guess = userSpecified; !isObDir(guess) {
		if guess = os.Getenv("OBPATH"); !isObDir(guess) {
			guess = "."
		}
	}
	return
}

func Init(dirPath string) (err error) {
	if DirPath, err = filepath.Abs(dirPath); (err == nil) && !isObDir(DirPath) {
		err = fmt.Errorf("The specified directory '%s' does not contain a valid OpenBase installation.", DirPath)
	}
	return
}
