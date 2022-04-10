package core

import (
	"errors"
	"fmt"
	"sync"
)

/* 一个AOI的格子类型 */

type Grid struct {
	GID       int          // 格子ID
	MinX      int          // 格子左边边界坐标
	MaxX      int          // 格子右边边界坐标
	MinY      int          // 格子上边边界坐标
	MaxY      int          // 格子下边边界坐标
	playerIDs map[int]bool // 当前格子内玩家或者物体成员的ID集合
	pIDLock   sync.RWMutex // 保护当前集合的锁
}

// NewGrid 初始化当前格子的方法
func NewGrid(gID, minX, maxX, minY, maxY int) *Grid {
	return &Grid{
		GID:       gID,
		MinX:      minX,
		MaxX:      maxX,
		MinY:      minY,
		MaxY:      maxY,
		playerIDs: make(map[int]bool),
	}
}

// Add 给格子添加一个玩家
func (g *Grid) Add(playerID int) {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()
	g.playerIDs[playerID] = true
}

// RemovePlayer 从格子中删除一个玩家
func (g *Grid) RemovePlayer(playerID int) error {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()

	if _, ok := g.playerIDs[playerID]; ok {
		delete(g.playerIDs, playerID)
		return nil
	}
	return errors.New("remove player failed")
}

// GetAllPlayerID 得到当前格子中的所有玩家
func (g *Grid) GetAllPlayerID() (players []int) {
	g.pIDLock.RLock()
	defer g.pIDLock.RUnlock()

	for k, _ := range g.playerIDs {
		players = append(players, k)
	}
	return players
}

// 调试使用-打印出格子的基本信息
func (g *Grid) String() string {
	return fmt.Sprintf("Grid ID:%d, minX:%d, maxX:%d, minY:%d, maxY:%d, playerIDs:%v",
		g.GID, g.MinX, g.MaxX, g.MinY, g.MaxY, g.playerIDs)
}
