package platform

import (
	model "github.com/xhigher/hzgo/demos/admin/model/platform"
)

type TraceLog struct {
	Module string      `json:"module"`
	Action string      `json:"action"`
	Params interface{} `json:"params"`
	Result interface{} `json:"result"`
	Roles  []string    `json:"roles"`
	Uid    string      `json:"uid"`
}

func AddLog(module, path string, params, result interface{}, roles []string, uid string) {
	err := model.AddTraceLog(module, path, params, result, roles, uid)
	if err != nil {
		return
	}
}
