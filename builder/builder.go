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
	if !buildStats.getIsBuilding() {
		buildStats.setIsBuilding(true)
		cmd := exec.Command("go", "build", "-i", "-o", filenames.BinaryName)
		output, err := cmd.CombinedOutput()
		buildStats.setIsBuilding(false)
		if err != nil {
			return string(output), err
		}
		return finishBuild()
	} else {
		buildStats.setShouldBuild(true)
		return "", BuildAlreayUnderWayError
	}
}

func finishBuild() (string, error) {
	// See if another build was scheduled while we built.
	if buildStats.getShouldBuild() {
		buildStats.setShouldBuild(false)
		return Build()
	}
	return "", nil
}

func (bs *BuildStats) setIsBuilding(isBuilding bool) {
	bs.IsBuildingLock.Lock()
	bs.IsBuilding = isBuilding
	bs.IsBuildingLock.Unlock()
}

func (bs *BuildStats) getIsBuilding() bool {
	bs.IsBuildingLock.RLock()
	defer bs.IsBuildingLock.RUnlock()
	return bs.IsBuilding
}

func (bs *BuildStats) setShouldBuild(shouldBuild bool) {
	bs.ShouldBuildLock.Lock()
	bs.ShouldBuild = shouldBuild
	bs.ShouldBuildLock.Unlock()
}

func (bs *BuildStats) getShouldBuild() bool {
	bs.ShouldBuildLock.RLock()
	defer bs.ShouldBuildLock.RUnlock()
	return bs.ShouldBuild
}
