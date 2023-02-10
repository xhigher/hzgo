package defines

type UseridData struct {
	Userid string `json:"userid"`
}

type TokenData struct {
	Token string `json:"token"`
	Et    int64  `json:"et"`
}

type PageData struct {
	Total int32 `json:"total"`
	Data interface{} `json:"data"`
	Offset int32 `json:"offset"`
	Limit int32 `json:"limit"`
}