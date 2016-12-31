package filenames

import (
	"log"
	"os"
	"runtime"
)

var (
	WorkingDirectory string
	BinaryName       string
)

func init() {
	// WorkingDirectory
	var err error
	WorkingDirectory, err = os.Getwd()
	if err != nil {
		log.Fatalln("Could not get working directory:", err)
	}
	// BinaryName
	if runtime.GOOS == "windows" {
		BinaryName = "klaus-bin.exe"
	} else {
		BinaryName = "klaus-bin"
	}
}
