package req

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