package core

import (
	"sync"
)

type WorldManager struct {
	AoiManager *AOIManager       // 当前世界地图的aoi管理模块
	Players    map[int32]*Player // 当前全部在线的player集合
	pLock      sync.RWMutex      // 保护player的锁
}

// WorldManagerObj 全局地图
var WorldManagerObj *WorldManager

// 初始化，创建一个全局的世界管理模块
func init() {
	WorldManagerObj = &WorldManager{
		// 创建地图
		AoiManager: NewAOIManager(AOI_MIN_X, AOI_MAX_X, AOI_CNTS_X, AOI_MIN_Y, AOI_MAX_Y, AOI_CNTS_Y),
		Players:    make(map[int32]*Player),
	}
}

// AddPlayer 添加一个玩家
func (wm *WorldManager) AddPlayer(player *Player) {
	wm.pLock.Lock()
	wm.Players[player.Pid] = player
	wm.pLock.Unlock()

	// 将player添加到aoiManager中
	wm.AoiManager.AddPidByPos(int(player.Pid), player.X, player.Z)
}

// RemovePlayerByPid 删除一个玩家
func (wm *WorldManager) RemovePlayerByPid(pid int32) {
	player := wm.Players[pid]
	wm.AoiManager.RemovePidByPos(int(pid), player.X, player.Z)

	wm.pLock.Lock()
	delete(wm.Players, pid)
	wm.pLock.Unlock()
}

// GetPlayerByPid 通过PID查询player对象
func (wm *WorldManager) GetPlayerByPid(pid int32) *Player {
	wm.pLock.RLock()
	defer wm.pLock.RUnlock()

	return wm.Players[pid]
}

// GetAllPlayers 获取所有的在线玩家
func (wm *WorldManager) GetAllPlayers() []*Player {
	wm.pLock.RLock()
	defer wm.pLock.RUnlock()

	players := make([]*Player, 0)

	for _, v := range wm.Players {
		players = append(players, v)
	}
	return players
}
