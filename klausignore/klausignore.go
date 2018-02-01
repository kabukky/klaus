package klausignore

import (
	"io/ioutil"
	"log"
	"strings"

	"github.com/kabukky/klaus/filenames"
)

const filename = ".klausignore"

//Ignore is a list of directories ignored by klaus
var Ignore []string

func init() {
	path := filenames.WorkingDirectory + "/" + filename
	contents, err := readLines(path)

	if err != nil {
		log.Printf("Error reading %s", path)

		Ignore = []string{}
	}

	Ignore = contents
}

func readLines(path string) ([]string, error) {
	content, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}
	return strings.Split(string(content), "\n"), nil
}
