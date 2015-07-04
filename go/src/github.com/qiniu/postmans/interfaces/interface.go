package interfaces

type Msg struct {
	From string
	Msg  string
}

type GroupMsg struct {
	GroupId string
	From    string
	Msg     string
}

type Postman interface {
	SendMsg(to string, msg string) error
	RecvMsg() chan Msg
	SendGroupMsg(gid string, msg string) error
	RecvGroupMsg() chan GroupMsg
}
