package logic

import (
	"fmt"
	"github.com/xhigher/hzgo/games/live_game/maps"
	"github.com/xhigher/hzgo/logger"
	"github.com/xhigher/hzgo/utils"
	"math/rand"
	"sort"
	"sync"
	"time"
)

var (
	robots = sync.Pool{
		New: func() interface{} {
			return newRobot(0)
		},
	}
)

func InitRobots(count int) {
	for i := 0; i < count; i++ {
		robots.Put(newRobot(i))
	}
}

func newRobot(i int) *Robot {
	return &Robot{
		&Player{
			id:          utils.IntToBase36(utils.NowTimeMillis() - 888999000000 + int64(i)),
			name:        utils.RandString(20),
			avatar:      fmt.Sprintf("https://hifun.yunwan.tech/res/img/avatar/%d.jpg", utils.RandInt64(100, 130)),
			role:        PlayerRobot,
			stepTime:    playerStepTime,
			skin:        randomPlayerSkin(),
			bubbleColor: randomBubbleColor(),
		},
		0,
		utils.RandInt32(10, 40),
		0,
		nil,
	}
}

func GeRobot() *Robot {
	robot := robots.Get().(*Robot)
	robot.bubbleCount = utils.RandInt(1, 2)
	robot.bubblePower = 1
	robot.pinCount = utils.RandInt(0, 1)

	robot.character = utils.RandInt(0, 3)
	return robot
}

func ReleaseRobot(r *Robot) {
	robots.Put(r)
}

type Robot struct {
	*Player
	character       int //性格  0为喜欢杀玩家   1为喜欢杀机器人  2喜欢炸箱子  3喜欢乱走乱放泡泡
	IQ              int32
	usedBubbleCount int
	timer           *time.Timer
}

func (r *Robot) Run() {
	if r.status == PlayerActive {
		logger.Infof("Robot.Run: id=%v, character=%v, IQ=%v", r.id, r.character, r.IQ)
		movableSites := &MovableSites{
			data: map[string]MovableSite{},
		}
		movableSites.Add(MovableSite{
			Base: &maps.Site{X: r.curSite.X, Y: r.curSite.Y}, Range: 0,
		})
		checkSites := []MovableSite{
			{Base: &maps.Site{X: r.curSite.X, Y: r.curSite.Y}},
		}
		r.getMovableSites(checkSites, movableSites, 1)

		logger.Infof("getMovableSites: id=%v, %v, %v", r.id, utils.JSONString(checkSites), utils.JSONString(movableSites))
		//获取预爆炸点
		bombSites := &BombSites{
			bombingBubbles: r.room.bubbles,
		}
		r.getBombSites(bombSites)
		logger.Infof("getBombSites: id=%v, %v", r.id, utils.JSONString(bombSites))
		//计算可以去到的点位的权重值，并附加到数组元素中
		siteList := movableSites.GetData()
		if len(siteList) > 0 {
			r.computeCanPosPower(siteList, bombSites.areas)
			//根据算出的点位权重再次排序,权重相同按照远近进行排序
			sort.SliceStable(siteList, func(i, j int) bool {
				return siteList[i].Power > siteList[j].Power
			})
			logger.Infof("getMovableSites: id=%v, siteList=%v", r.id, utils.JSONString(siteList))
			//拿到最高权重的目标位置
			powerSite := siteList[0]
			//前往该点位
			if powerSite.Start == nil {
				r.moveStop()
				if r.bubbleCount > 0 {
					if !bombSites.ExistBombArray(r.curSite) {
						r.usedBubbleCount++
						r.CreateBubble()
					}
				}
			} else {
				startSite := *powerSite.Start
				isDanger := false
				//智商大于0.7，且随机一般几率触发紧急避险
				if r.IQ > 30 && rand.Intn(10) < 5 {
					//紧急避险
					//获取即将爆炸的点位
					now := time.Now()
					delayTime := int64(r.stepTime * 2)
					var bombingBubbles []*Bubble
					var leftBubbles []*Bubble
					for _, b := range r.room.bubbles {
						if b.bombTime.Sub(now) < time.Duration(delayTime) {
							bombingBubbles = append(bombingBubbles, b)
						} else {
							leftBubbles = append(leftBubbles, b)
						}
					}
					bombSites2 := &BombSites{
						bombingBubbles: bombingBubbles,
						leftBubbles:    leftBubbles,
					}
					r.getBombSites(bombSites2)

					logger.Infof("id=%v, bombSites2.areas=%v, startSite=%v", r.id, utils.JSONString(bombSites2.areas), utils.JSONString(startSite))
					if utils.InArray(bombSites2.areas, startSite) {
						logger.Infof("InArray, id=%v, bombSites2.areas=%v, startSite=%v", r.id, utils.JSONString(bombSites2.areas), utils.JSONString(startSite))
						var toCheckSites []maps.Site
						for _, s := range siteList {
							if s.Range == 1 && !startSite.Equal(*s.Base) {
								toCheckSites = append(toCheckSites, *s.Base)
							}
						}
						for _, s := range toCheckSites {
							if utils.InArray(bombSites2.areas, s) {
								isDanger = true
								powerSite.Start = &maps.Site{
									X: s.X,
									Y: s.Y,
								}
								break
							}
						}
					}

				}

				//非紧急避险状态下，旁边如果有地方玩家，则放泡泡
				if !isDanger && r.bubbleCount > 0 {
					hasNearPlayer := false
					for _, s := range siteList {
						if s.Range <= 2 {
							for _, p := range r.room.players {
								if p.id != r.id && p.curSite.Equal(*s.Base) {
									hasNearPlayer = true
									break
								}
							}
						}
						if hasNearPlayer {
							break
						}
					}
					if hasNearPlayer {
						r.usedBubbleCount++
						r.CreateBubble()
					}

				}

				//如果当前点不是爆炸点，则随机概率不动
				if rand.Int31n(100) > r.IQ {
					if utils.InArray(bombSites.areas, r.curSite) {
						r.SetSite(*powerSite.Start)
					} else {
						r.moveStop()
					}
				} else {
					r.SetSite(*powerSite.Start)
				}
			}
		}

	}

	if r.IsTrapped() && r.pinCount > 0 {
		secs := utils.RandInt64(1000, 2000)
		r.timer = time.AfterFunc(time.Millisecond*time.Duration(secs), func() {
			r.usePin()
			r.timer = time.AfterFunc(time.Millisecond*time.Duration(r.stepTime), r.Run)
		})
	} else {
		r.timer = time.AfterFunc(time.Millisecond*time.Duration(r.stepTime), r.Run)
	}
}

func (r *Robot) getBoxSiteCount(s maps.Site) int {
	boxCount := 0
	for i := 1; i <= r.bubblePower; i++ {
		sites := r.room.mapData.GetValidSites([]maps.Site{
			{X: s.X - i, Y: s.Y},
			{X: s.X + i, Y: s.Y},
			{X: s.X, Y: s.Y - i},
			{X: s.X, Y: s.Y + i},
		})
		for _, ts := range sites {
			if r.room.mapData.ExistBox(ts) {
				boxCount++
			}
		}
	}
	return boxCount
}

func (r *Robot) getEmptySiteCount(s maps.Site) int {
	emptyCount := 0
	sites := r.room.mapData.GetValidSites([]maps.Site{
		{X: s.X - 1, Y: s.Y},
		{X: s.X + 1, Y: s.Y},
		{X: s.X, Y: s.Y - 1},
		{X: s.X, Y: s.Y + 1},
	})
	for _, ts := range sites {
		if r.room.mapData.ExistBox(ts) {
			emptyCount++
		}
	}
	return emptyCount
}
func (r *Robot) checkMovableSite(s maps.Site) bool {
	if r.curSite.Equal(s) {
		return true
	}
	if !r.room.mapData.IsEmptySite(s) {
		return false
	}
	if yes, _ := r.room.ExistBubble(s); yes {
		return false
	}
	return true
}

func (r *Robot) computeCanPosPower(movableSites []MovableSite, bombSites []maps.Site) {
	for i, ms := range movableSites {
		s := *ms.Base
		power := 0

		//如果该点位四周有箱子，则根据箱子个数加权重
		boxCount := r.getBoxSiteCount(s)
		//判断四周能移动的点位有几个，2个以上每一个减权重10
		emptyCount := r.getEmptySiteCount(s)

		//如果该点已经放了有泡泡，则减去权重
		isBubble, _ := r.room.ExistBubble(s)

		//如果该点位是预爆炸点，则减去权重
		isBombing := utils.InArray(bombSites, s)

		//如果该点位有道具，则增加权重
		isProp := false
		if ms.Range <= 6 {
			isProp, _ = r.room.ExistProp(s)
		}

		//这个点是否有玩家
		isNearUser, p := r.room.ExistPlayer(s)

		//如果附近有被泡住的玩家，则权重很高
		isTrappedUser := false
		if ms.Range <= 10 {
			if isNearUser && p.IsTrapped() {
				isTrappedUser = true
			}
		}

		//这个点是否有机器人
		isNearBot := false
		if p != nil && p.IsRobot() {
			isNearBot = true
		}

		//游戏内是否还剩余有玩家
		//hasHumanPlayer := r.room.HasHumanPlayer()

		switch r.character {
		case 0:
			{
				power += boxCount * 5
				power -= ms.Range * 2
				if r.bubbleCount != 1 || rand.Intn(10) < 5 {
					power -= (emptyCount - 1) * 5
				}
				if isBubble {
					power -= 100
				}
				if isBombing {
					power -= 290
				}
				if isProp && (r.bubbleCount < 4 || r.bubblePower < 4) {
					if yes, prop := r.room.ExistProp(s); yes {
						if prop.typ != 2 || r.stepTime > 200 {
							power += 50
						}
					}
				}
				if isTrappedUser {
					power += 500
				}
				if isNearUser {
					power += 300
				}
				break
			}
		case 1:
			{
				power += boxCount * 5
				power -= ms.Range * 2
				if r.bubbleCount != 1 || rand.Intn(10) < 5 {
					power -= (emptyCount - 1) * 5
				}
				if isBubble {
					power -= 100
				}
				if isBombing {
					power -= 290
				}
				if isProp && (r.bubbleCount < 4 || r.bubblePower < 4) {
					if yes, prop := r.room.ExistProp(s); yes {
						if prop.typ != 2 || r.stepTime > 200 {
							power += 50
						}
					}
				}
				if isTrappedUser {
					power += 50
				}
				if isNearBot {
					power += 300
				}
				break
			}
		case 2:
			{

				power += boxCount * 50
				power -= ms.Range * 20
				if r.bubbleCount != 1 || rand.Intn(10) < 5 {
					power -= (emptyCount - 1) * 50
				}
				if isBubble {
					power -= 100
				}
				if isBombing {
					power -= 300
				}
				if isProp {
					if yes, prop := r.room.ExistProp(s); yes {
						if prop.typ != 2 || r.stepTime > 200 {
							power += 50
						}
					}
				}
				if isTrappedUser {
					power += 500
				}
				break
			}
		case 3:
			{
				if boxCount < 3 {
					power += boxCount * 5
				}

				power -= ms.Range * 2
				if r.bubbleCount != 1 || rand.Intn(10) < 5 {
					power -= (emptyCount - 1) * 5
				}
				if isBubble {
					power -= 100
				}
				if isBombing {
					power -= 300
				}
				if isProp {
					if yes, prop := r.room.ExistProp(s); yes {
						if prop.typ != 2 || r.stepTime > 200 {
							power += 50
						}
					}
				}
				if isTrappedUser {
					power += 500
				}
				break
			}
		}

		movableSites[i].Power = power
		logger.Infof("computeCanPosPower: id=%v, %v", r.id, utils.JSONString(movableSites[i].Power))
	}
}

type MovableSite struct {
	Base  *maps.Site
	Start *maps.Site
	Range int
	Power int
}

type MovableSites struct {
	data map[string]MovableSite
}

func (m *MovableSites) Exists(site maps.Site) bool {
	if _, ok := m.data[fmt.Sprintf("%d-%d", site.X, site.Y)]; ok {
		return true
	}
	return false
}

func (m *MovableSites) Add(site MovableSite) {
	m.data[fmt.Sprintf("%d-%d", site.Base.X, site.Base.Y)] = site
}

func (m *MovableSites) GetData() []MovableSite {
	var data []MovableSite
	for _, s := range m.data {
		data = append(data, s)
	}
	return data
}

func (r *Robot) checkMoveSite(site maps.Site) bool {
	if r.curSite.Equal(site) {
		return true
	}
	if !r.room.mapData.IsEmptySite(site) {
		return false
	}
	//如果该坐标有泡泡，不能去
	if yes, _ := r.room.ExistBubble(site); yes {
		return false
	}
	logger.Infof("checkMoveSite: id=%v, %v, true", r.id, utils.JSONString(site))
	return true
}
func (r *Robot) getMovableSites(checkSites []MovableSite, movableSites *MovableSites, rn int) {
	var nextCheckSites []MovableSite //下一轮需要检测的点数组
	for i := 0; i < len(checkSites); i++ {
		s1 := checkSites[i]
		sites := r.room.mapData.GetValidSites([]maps.Site{
			{X: s1.Base.X, Y: s1.Base.Y + 1},
			{X: s1.Base.X, Y: s1.Base.Y - 1},
			{X: s1.Base.X + 1, Y: s1.Base.Y},
			{X: s1.Base.X - 1, Y: s1.Base.Y},
		})
		logger.Infof("GetValidSites: id=%v, %v", r.id, utils.JSONString(sites))
		for _, s2 := range sites {
			if r.checkMoveSite(s2) {
				if !movableSites.Exists(s2) {
					s3 := MovableSite{
						Base:  &s2,
						Range: rn,
					}
					if s1.Start != nil {
						s3.Start = s1.Start
					} else {
						s3.Start = &s2
					}
					movableSites.Add(s3)
					nextCheckSites = append(nextCheckSites, s3)
				}
			}
		}
	}
	logger.Infof("nextCheckSites: id=%v, %v, %v", r.id, utils.JSONString(nextCheckSites), utils.JSONString(movableSites))
	if len(nextCheckSites) > 0 {
		rn++
		r.getMovableSites(nextCheckSites, movableSites, rn)
	}
}

func (r *Robot) Stop() {
	r.timer.Stop()
}

type BombSites struct {
	areas          []maps.Site
	bombingBubbles []*Bubble //即将爆炸的泡泡
	leftBubbles    []*Bubble
}

func (bs *BombSites) AddBombArray(s maps.Site) {
	yes := false
	for _, b := range bs.areas {
		if b.Equal(s) {
			yes = true
			break
		}
	}
	if !yes {
		bs.areas = append(bs.areas, s)
	}
}

func (bs *BombSites) ExistBombArray(s maps.Site) bool {
	for _, b := range bs.areas {
		if b.Equal(s) {
			return true
		}
	}
	return false
}

func (r *Robot) getBombSites(data *BombSites) {
	for _, b := range data.bombingBubbles {
		data.AddBombArray(b.site)
		for i := 1; i <= b.power; i++ {
			sites := r.room.mapData.GetValidSites([]maps.Site{
				{X: b.site.X + i, Y: b.site.Y},
				{X: b.site.X - i, Y: b.site.Y},
				{X: b.site.X, Y: b.site.Y + i},
				{X: b.site.X, Y: b.site.Y - i},
			})
			for _, s := range sites {
				if r.room.mapData.IsEmptySite(s) {
					data.AddBombArray(s)
				}
			}
		}
	}
	//判断剩余的泡泡是否有在当前爆炸点内的，如果有，则计算引爆点
	var nextBombingBubbles []*Bubble
	var nextLeftBubbles []*Bubble
	for _, b := range data.leftBubbles {
		if data.ExistBombArray(b.site) {
			nextBombingBubbles = append(nextBombingBubbles, b)
		} else {
			nextLeftBubbles = append(nextLeftBubbles, b)
		}
	}
	if len(nextBombingBubbles) > 0 {
		data.bombingBubbles = nextBombingBubbles
		data.leftBubbles = nextLeftBubbles
		r.getBombSites(data)
	}
}
