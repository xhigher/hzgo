package req

type UseridReq struct {
	Userid string `form:"userid" json:"userid" query:"userid"`
}

type RegisterReq struct {
	Username string `form:"username" json:"username" query:"username"`
	Password string `form:"password" json:"password" query:"password"`
}
