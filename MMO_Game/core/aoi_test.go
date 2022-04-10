package core

import (
	"fmt"
	"testing"
)

//func TestNewAOIManager(t *testing.T) {
//	aoiManager := NewAOIManager(0, 250, 5, 0, 250, 5)
//	fmt.Println(aoiManager.String())
//}

func TestAOIManager_GetSurroundGridsByGidGetSurroundGridsByGid(t *testing.T) {
	aoiManager := NewAOIManager(0, 250, 5, 0, 250, 5)
	fmt.Println("test begin")

	for gid, _ := range aoiManager.grids {
		grids := aoiManager.GetSurroundGridsByGid(gid)
		fmt.Println("gid: ", gid, " grids len = ", len(grids))

		gIDs := make([]int, 0, len(grids))
		for _, grid := range grids {
			gIDs = append(gIDs, grid.GID)
		}
		fmt.Printf("surrounding grid IDs are %v\n", gIDs)
	}
}
