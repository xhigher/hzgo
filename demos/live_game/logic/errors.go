package logic

import "errors"

var (
	errMatchWait    = errors.New("游戏报名还未开始")
	errMatchOngoing = errors.New("游戏进行中")
	errMatchEnd     = errors.New("游戏已结束")
	errPlayerJoined = errors.New("您已报名")
)
