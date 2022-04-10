package core

import "fmt"

// 定义一些AOI的边界值
const (
	AOI_MIN_X  int = 50
	AOI_MAX_X  int = 400
	AOI_CNTS_X int = 5
	AOI_MIN_Y  int = 50
	AOI_MAX_Y  int = 400
	AOI_CNTS_Y int = 5
)

/* AOI 区域管理模块 */

type AOIManager struct {
	MinX  int           // 区域的左边界坐标
	MaxX  int           // 区域的右边界坐标
	CntX  int           // X轴方向的格子数量
	MinY  int           // 区域的上边界坐标
	MaxY  int           // 区域的下边界坐标
	CntY  int           // Y方向格子的数量
	grids map[int]*Grid // 当前区域中有哪些格子 key = 格子ID，value = 格子对象
}

// NewAOIManager 初始化一个AOI区域管理模块
func NewAOIManager(minX, maxX, cntX, minY, maxY, cntY int) *AOIManager {
	aoiManager := &AOIManager{
		MinX:  minX,
		MaxX:  maxX,
		CntX:  cntX,
		MinY:  minY,
		MaxY:  maxY,
		CntY:  cntY,
		grids: make(map[int]*Grid),
	}

	// 给AOI初始化区域的所有格子进行编号和初始化
	for y := 0; y < cntY; y++ {
		for x := 0; x < cntX; x++ {
			// 格子编号
			gid := y*cntX + x
			aoiManager.grids[gid] = NewGrid(
				gid,
				minX+x*aoiManager.gridWidth(),
				minX+(x+1)*aoiManager.gridWidth(),
				minY+y*aoiManager.gridLength(),
				minY+(y+1)*aoiManager.gridLength(),
			)
		}
	}

	return aoiManager
}

// 得到每个格子在X轴方向的宽度
func (m *AOIManager) gridWidth() int {
	return (m.MaxX - m.MinX) / m.CntX
}

// 得到每个格子在Y轴方向的宽度
func (m *AOIManager) gridLength() int {
	return (m.MaxY - m.MinY) / m.CntY
}

// GetSurroundGridsByGid 根据格子GID得到周边九宫格格子的id集合
func (m *AOIManager) GetSurroundGridsByGid(gid int) (grids []*Grid) {
	// 判断gid是否在aoiManager中
	if _, ok := m.grids[gid]; !ok {
		return
	}

	// 初始化grids切片
	grids = append(grids, m.grids[gid])

	// 得到gid的x轴坐标  idx = gid % nx
	idx := gid % m.CntX

	// 判断gid左边是否有格子
	if idx > 0 {
		grids = append(grids, m.grids[gid-1])
	}

	// 判断gid右边是否有格子
	if idx < m.CntX-1 {
		grids = append(grids, m.grids[gid+1])
	}

	// 将X轴上的格子都取出，进行遍历，判断每个格子上下是否还有格子
	// 得到X轴格子的ID集合
	gidX := make([]int, 0, len(grids))
	for _, v := range grids {
		gidX = append(gidX, v.GID)
	}

	//遍历gidX
	for _, v := range gidX {
		idy := v / m.CntY
		// 上方是否有格子
		if idy > 0 {
			grids = append(grids, m.grids[v-m.CntX])
		}
		// 下方是否有格子
		if idy < m.CntY-1 {
			grids = append(grids, m.grids[v+m.CntX])
		}
	}
	return
}

// GetGidByPos 通过横纵坐标得到GID
func (m *AOIManager) GetGidByPos(x, y float32) int {
	idx := (int(x) - m.MinX) / m.gridWidth()
	idy := (int(y) - m.MinY) / m.gridLength()

	return idy*m.CntX + idx
}

// GetPlayerIdsByPos 通过横纵坐标得到周边九宫格内全部的playerID
func (m *AOIManager) GetPlayerIdsByPos(x, y float32) (playerIDs []int) {
	// 得到当前玩家的ID
	gid := m.GetGidByPos(x, y)
	// 得到当前九宫格的格子信息
	grids := m.GetSurroundGridsByGid(gid)

	// 将九宫格的信息全部的Player的ID 加到playerIDs
	for _, v := range grids {
		playerIDs = append(playerIDs, v.GetAllPlayerID()...)
	}
	return
}

// AddPidToGrid 添加一个PlayerID到一个格子中
func (m *AOIManager) AddPidToGrid(pid, gid int) {
	m.grids[gid].Add(pid)
}

// RemovePidFromGrid 移除一个格子中的PlayerID
func (m *AOIManager) RemovePidFromGrid(pid, gid int) {
	m.grids[gid].RemovePlayer(pid)
}

// GetPidsByGid 通过GID获取全部的PlayerID
func (m *AOIManager) GetPidsByGid(gid int) (playerIDs []int) {
	playerIDs = m.GetPidsByGid(gid)
	return
}

// AddPidByPos 通过坐标将PlayerID添加到一个格子中
func (m *AOIManager) AddPidByPos(pid int, x, y float32) {
	gid := m.GetGidByPos(x, y)
	m.grids[gid].Add(pid)
}

// RemovePidByPos 通过坐标把一个PlayerID	从一个格子中删除
func (m *AOIManager) RemovePidByPos(pid int, x, y float32) {
	gid := m.GetGidByPos(x, y)
	m.grids[gid].RemovePlayer(pid)
}

// 打印格子信息
func (m *AOIManager) String() string {
	s := fmt.Sprintf("AOIManager:\nMinX:%d, MaxX:%d, CntX:%d, MinY:%d, MaxY:%d, CntY:%d\n", m.MinX, m.MaxX, m.CntX, m.MinY, m.MaxY, m.CntY)

	for _, grid := range m.grids {
		s += fmt.Sprintln(grid.String())
	}
	return s
}
