package builder

import (
	"errors"
	"os/exec"
	"sync"

	"github.com/kabukky/klaus/filenames"
)

var (
	buildStats               = &BuildStats{}
	BuildAlreayUnderWayError = errors.New("A build is already under way.")
)

type BuildStats struct {
	IsBuilding      bool
	IsBuildingLock  sync.RWMutex
	ShouldBuild     bool
	ShouldBuildLock sync.RWMutex
}

func Build() (string, error) {
	buildStats.IsBuildingLock.RLock()
	if !buildStats.IsBuilding {
		buildStats.IsBuildingLock.RUnlock()
		buildStats.IsBuildingLock.Lock()
		buildStats.IsBuilding = true
		buildStats.IsBuildingLock.Unlock()
		cmd := exec.Command("go", "build", "-o", filenames.BinaryName)
		output, err := cmd.CombinedOutput()
		buildStats.IsBuildingLock.Lock()
		buildStats.IsBuilding = false
		buildStats.IsBuildingLock.Unlock()
		if err != nil {
			return string(output), err
		}
		return finishBuild()
	} else {
		buildStats.IsBuildingLock.RUnlock()
		buildStats.ShouldBuildLock.Lock()
		buildStats.ShouldBuild = true
		buildStats.ShouldBuildLock.Unlock()
		return "", BuildAlreayUnderWayError
	}
}

func finishBuild() (string, error) {
	// See if another build was scheduled while we built.
	buildStats.ShouldBuildLock.RLock()
	if buildStats.ShouldBuild {
		buildStats.ShouldBuildLock.RUnlock()
		buildStats.ShouldBuildLock.Lock()
		buildStats.ShouldBuild = false
		buildStats.ShouldBuildLock.Unlock()
		return Build()
	} else {
		buildStats.ShouldBuildLock.RUnlock()
	}
	return "", nil
}
