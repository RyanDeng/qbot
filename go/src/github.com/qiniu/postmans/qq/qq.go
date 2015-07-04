package qq

import (
	"encoding/json"
	"fmt"
	. "github.com/qiniu/postmans/interfaces"
	"net/http"
	"net/url"
)

type config struct {
	Host   string `json:"host"`
	JsHost string `json:"jsHost"`
}

type QQ struct {
	config
	buddyMsgs chan Msg
	groupMsgs chan GroupMsg
}

func NewQQ(confStr string) (man Postman, err error) {

	var conf config
	err = json.Unmarshal([]byte(confStr), &conf)
	if err != nil {
		return
	}
	qq := &QQ{
		config:    conf,
		buddyMsgs: make(chan Msg),
		groupMsgs: make(chan GroupMsg),
	}
	go func() {
		panic(qq.Run())
	}()

	man = qq
	return
}

func (qq *QQ) Type() string {
	return "QQ"
}

func (qq *QQ) SendMsg(to, msg string) (err error) {

	resp, err := http.PostForm(qq.JsHost+"/msg", url.Values{"to": {to}, "msg": {msg}})
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("error: %d", resp.StatusCode)
	}

	return
}

func (qq *QQ) SendGroupMsg(gid, msg string) (err error) {

	resp, err := http.PostForm(qq.JsHost+"/grpMsg", url.Values{"group": {gid}, "msg": {msg}})
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("error: %d", resp.StatusCode)
	}
	return
}

func (qq *QQ) RecvMsg() chan Msg {
	return qq.buddyMsgs
}

func (qq *QQ) RecvGroupMsg() chan GroupMsg {
	return qq.groupMsgs
}

func (qq *QQ) msg(w http.ResponseWriter, r *http.Request) {

	from := r.FormValue("from")
	msg := r.FormValue("msg")
	if from == "" {
		w.WriteHeader(400)
		return
	}
	qq.buddyMsgs <- Msg{from, msg}
	w.WriteHeader(200)
}

func (qq *QQ) grpMsg(w http.ResponseWriter, r *http.Request) {

	from := r.FormValue("from")
	groupid := r.FormValue("group")
	msg := r.FormValue("msg")
	if from == "" || groupid == "" {
		w.WriteHeader(400)
		return
	}
	qq.groupMsgs <- GroupMsg{groupid, from, msg}
	w.WriteHeader(200)
}

func (qq *QQ) Run() (err error) {

	http.HandleFunc("/msg", qq.msg)
	http.HandleFunc("/grpMsg", qq.grpMsg)
	return http.ListenAndServe(qq.Host, nil)
}

func init() {
	Register("QQ", NewQQ)
}
