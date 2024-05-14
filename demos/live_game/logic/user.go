package logic

import (
	"github.com/xhigher/hzgo/demos/live_game/model/store"
	"github.com/xhigher/hzgo/utils"
)

var testUsers = map[string]store.UserInfo{
	"fsbi5y": {
		Id:       "fsbi5y",
		Nickname: "Tom",
		Avatar:   "https://hifun.yunwan.tech/res/img/avatar/100.jpg",
		Token:    "CSQFEOXZOKAPTENXPOYW",
		Sex:      1,
		Level:    1,
		Openid:   "CSQFEOXZOKAPTENXPOYW",
		Skin:     randomPlayerSkin(),
	},
	"fsbi62": {
		Id:       "fsbi62",
		Nickname: "Tom",
		Avatar:   "https://hifun.yunwan.tech/res/img/avatar/101.jpg",
		Token:    "GFDUVSYXPHFFAZMUGCAI",
		Sex:      1,
		Level:    1,
		Openid:   "GFDUVSYXPHFFAZMUGCAI",
		Skin:     randomPlayerSkin(),
	},
	"fsbi66": {
		Id:       "fsbi66",
		Nickname: "Tom",
		Avatar:   "https://hifun.yunwan.tech/res/img/avatar/102.jpg",
		Token:    "BSZMDGOFPSGORNYQYQGU",
		Sex:      1,
		Level:    1,
		Openid:   "BSZMDGOFPSGORNYQYQGU",
		Skin:     randomPlayerSkin(),
	},
	"fsbi6a": {
		Id:       "fsbi6a",
		Nickname: "Tom",
		Avatar:   "https://hifun.yunwan.tech/res/img/avatar/103.jpg",
		Token:    "ZGLQDITYOBJCMCNRIODT",
		Sex:      1,
		Level:    1,
		Openid:   "ZGLQDITYOBJCMCNRIODT",
		Skin:     randomPlayerSkin(),
	},
	"fsbi6e": {
		Id:       "fsbi6e",
		Nickname: "Tom",
		Avatar:   "https://hifun.yunwan.tech/res/img/avatar/104.jpg",
		Token:    "JTLSGEBRQLWVPFCQEUYL",
		Sex:      1,
		Level:    1,
		Openid:   "JTLSGEBRQLWVPFCQEUYL",
		Skin:     randomPlayerSkin(),
	},
	"fsbi6i": {
		Id:       "fsbi6i",
		Nickname: "Tom",
		Avatar:   "https://hifun.yunwan.tech/res/img/avatar/105.jpg",
		Token:    "CRLZESWVNUQQLYUXKGOB",
		Sex:      1,
		Level:    1,
		Openid:   "CRLZESWVNUQQLYUXKGOB",
		Skin:     randomPlayerSkin(),
	},
	"fsbi6m": {
		Id:       "fsbi6m",
		Nickname: "Tom",
		Avatar:   "https://hifun.yunwan.tech/res/img/avatar/106.jpg",
		Token:    "QDZOVSRIJSAPOKRQUEKY",
		Sex:      1,
		Level:    1,
		Openid:   "QDZOVSRIJSAPOKRQUEKY",
		Skin:     randomPlayerSkin(),
	},
	"fsbi6q": {
		Id:       "fsbi6q",
		Nickname: "Tom",
		Avatar:   "https://hifun.yunwan.tech/res/img/avatar/107.jpg",
		Token:    "CWDMXMAJHZRGKNJTYJJH",
		Sex:      1,
		Level:    1,
		Openid:   "CWDMXMAJHZRGKNJTYJJH",
		Skin:     randomPlayerSkin(),
	},
	"fsbi6u": {
		Id:       "fsbi6u",
		Nickname: "Tom",
		Avatar:   "https://hifun.yunwan.tech/res/img/avatar/108.jpg",
		Token:    "WYNNGDYAARFLQTVDENCX",
		Sex:      1,
		Level:    1,
		Openid:   "WYNNGDYAARFLQTVDENCX",
		Skin:     randomPlayerSkin(),
	},
	"fsbi6y": {
		Id:       "fsbi6y",
		Nickname: "Tom",
		Avatar:   "https://hifun.yunwan.tech/res/img/avatar/109.jpg",
		Token:    "IKFQDBKFBGZLLNVRNMFP",
		Sex:      1,
		Level:    1,
		Openid:   "IKFQDBKFBGZLLNVRNMFP",
		Skin:     randomPlayerSkin(),
	},
}

func GetUserId() string {
	return utils.IntToBase36(utils.NowTimeMillis() - 888999000000)
}

func GetUser(id string) (yes bool, user store.UserInfo) {
	user, yes = testUsers[id]
	return
}
