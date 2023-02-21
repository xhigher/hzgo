package defines

import "strconv"

type AppLoc struct {
	Lat  float64 `json:"lat,omitempty"`
	Lng  float64 `json:"lng,omitempty"`
	City string  `json:"city,omitempty"`
	Addr string  `json:"addr,omitempty"`
}

type BaseParams struct {
	//应用包名
	Ap string `json:"app,omitempty"`
	//应用版本号
	Av string `json:"av,omitempty"`
	//dev_type 操作系统 1=android, 2=iOS, 3=windows
	Dt int32 `json:"dt,omitempty"`
	//brand 品牌
	Bd string `json:"bd,omitempty"`
	//model 机型
	Md string `json:"md,omitempty"`
	//os 安卓系统版本/iOS系统版本/UA
	Os string `json:"os,omitempty"`
	//dev_id 设备标识ID，唯一
	Did string `json:"did,omitempty"`
	//net_type 网络类型 1: Wi-Fi 2: 2G或3G 3: 4G 4: 其他
	Nt int32 `json:"nt,omitempty"`
	//channel 安装渠道
	Ch string `json:"ch,omitempty"`
	//IP地址
	Ip string `json:"ip,omitempty"`
	//location 地理位置
	Loc string `json:"loc,omitempty"`
	//imei 原值
	Imei string `json:"imei,omitempty"`
	//oaid 原值
	Oaid string `json:"oaid,omitempty"`
	//iOS广告标识符 原值
	Idfa string `json:"idfa,omitempty"`
	//数据摘要
	Ds string `json:"ds,omitempty"`
	//sign 签名
	Sign string `json:"sign,omitempty"`
	//timestamp 时间戳
	Ts int64 `json:"ts,omitempty"`
}

type UseridReq struct {
	Userid string `form:"userid" json:"userid" query:"userid"`
}

type RegisterReq struct {
	Username string `form:"username" json:"username" query:"username"`
	Password string `form:"password" json:"password" query:"password"`
}

type LoginReq struct {
	Username string `form:"username" json:"username" query:"username"`
	Password string `form:"password" json:"password" query:"password"`
}

type TokenUpdateReq struct {
	Audience string `form:"audience" json:"audience" query:"audience"`
	TokenId string `form:"token_id" json:"token_id" query:"token_id"`
	ExpiredAt int64 `form:"expired_at" json:"expired_at" query:"expired_at"`
	IssuedAt int64 `form:"issued_at" json:"issued_at" query:"issued_at"`
}

type TokenCheckReq struct {
	Userid string `form:"userid" json:"userid" query:"userid"`
	TokenId string `form:"token_id" json:"token_id" query:"token_id"`
}

type StatusPageReq struct {
	Status int32 `form:"status" json:"status" query:"status"`
	Offset int32 `form:"offset" json:"offset" query:"offset"`
	Limit int32  `form:"limit" json:"limit" query:"limit"`
}

type ChangeStatusReq struct {
	CommonIdReq
	Status int32 `form:"status" json:"status" query:"status"`
}

type BannerReq struct {
	Site string `json:"site"`
}

type ConfigReq struct {
	Sum string `json:"sum"`
	Id string `json:"id"`
}

type CommonIdReq struct {
	Id string `form:"id" json:"id" query:"id"`
}

func (r CommonIdReq) IntId() int32 {
	i, _ := strconv.Atoi(r.Id)
	return int32(i)
}

func (r CommonIdReq) Int64Id() int64 {
	i, _ := strconv.Atoi(r.Id)
	return int64(i)
}