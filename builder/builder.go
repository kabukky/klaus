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
		setIsBuildung(true)
		cmd := exec.Command("go", "build", "-o", filenames.BinaryName)
		output, err := cmd.CombinedOutput()
		setIsBuildung(false)
		if err != nil {
			return string(output), err
		}
		return finishBuild()
	} else {
		buildStats.IsBuildingLock.RUnlock()
		setShouldBuild(true)
		return "", BuildAlreayUnderWayError
	}
}

func finishBuild() (string, error) {
	// See if another build was scheduled while we built.
	buildStats.ShouldBuildLock.RLock()
	if buildStats.ShouldBuild {
		buildStats.ShouldBuildLock.RUnlock()
		setShouldBuild(false)
		return Build()
	} else {
		buildStats.ShouldBuildLock.RUnlock()
	}
	return "", nil
}

func setIsBuildung(isBuilding bool) {
	buildStats.IsBuildingLock.Lock()
	buildStats.IsBuilding = isBuilding
	buildStats.IsBuildingLock.Unlock()
}

func setShouldBuild(shouldBuild bool) {
	buildStats.ShouldBuildLock.Lock()
	buildStats.ShouldBuild = shouldBuild
	buildStats.ShouldBuildLock.Unlock()
}
