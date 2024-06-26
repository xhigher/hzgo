package events

const (

	//客户端事件 begin
	Join  = "matching"
	Ready = "ready"

	CreateBubble = "createBubble"
	SetPos       = "setPos"
	MoveStop     = "moveStop"
	UseNeedle    = "useNeedle"

	//客户端事件 end

	DisappearProp = "disappearProp"

	Login        = "login"
	LoginError   = "loginError"
	LoginSuccess = "loginSuccess"

	JoinError   = "joinError"
	JoinSuccess = "matching"

	LoadProcess = "loadProcess"
	GameReady   = "gameReady"
	GameStart   = "gameStart"
	GameOver    = "gameOver"
	exitRoom    = "exitRoom"

	BubbleBomb       = "bomb"
	Move             = "move"
	ChangeUserStatus = "changeUserStatus"
)
