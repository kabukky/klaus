package runner

import (
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/kabukky/klaus/filenames"
)

var (
	binary Binary
)

type Binary struct {
	sync.Mutex
	cmd *exec.Cmd
}

func Run() error {
	binary.Lock()
	defer binary.Unlock()
	if binary.cmd != nil && binary.cmd.Process != nil {
		log.Println("Killing running process...")
		err := binary.cmd.Process.Kill()
		if err != nil {
			log.Println("Couldn't kill process:", err)
		}
	}
	binary.cmd = exec.Command(filepath.Join(filenames.WorkingDirectory, filenames.BinaryName), os.Args[1:]...)
	stdout, err := binary.cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := binary.cmd.StderrPipe()
	if err != nil {
		return err
	}
	go io.Copy(os.Stdout, stdout)
	go io.Copy(os.Stdout, stderr)
	err = binary.cmd.Start()
	if err != nil {
		return err
	}
	go binary.cmd.Wait()
	log.Println("Starting process.")
	return nil
}
