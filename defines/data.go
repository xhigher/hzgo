package defines

type UseridData struct {
	Userid string `json:"userid"`
}

type TokenData struct {
	Token string `json:"token"`
	Et    int64  `json:"et"`
}
