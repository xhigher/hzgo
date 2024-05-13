package notice

type RegisterReq struct {
	Uid string `json:"uid"`
	Did string `json:"did"`
	Nid string `json:"nid"`
}

type MessageBroadcastReq struct {
	From string `json:"from"`
	To   string `json:"to"`
	Data []byte `json:"data"`
}
