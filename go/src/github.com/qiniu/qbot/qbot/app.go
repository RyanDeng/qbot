package main

import (
	"encoding/json"
	"github.com/qiniu/postmans/interfaces"
	"github.com/qiniu/qbot"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type Config struct {
	DBSettings qbot.DBSettings `json:"db"`
	ManagePost uint32          `json:"admin_post"`
}

func main() {

	confFile, err := os.Open("service.conf")
	if err != nil {
		log.Fatal("no config file")
		os.Exit(-1)
	}

	confData, err := ioutil.ReadAll(confFile)
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}

	var conf Config
	err = json.Unmarshal(confData, &conf)
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}
	cfg := &qbot.Config{
		DBSettings: conf.DBSettings,
	}
	service, err := qbot.NewService(cfg)
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}

	go func() {
		time.Sleep(5 * time.Second)
		r, ok, err := service.ReminderTbl.GetAndDelete()
		if err != nil {
			log.Error("GetAndDelete:", err)
			continue
		}
		if ok {
			for _, to := range r.Tos {
				postman.SendMsg(to, r.Event)
			}
		}
	}()

	var postman interfaces.Postman
	for {
		select {
		case msg := <-postman.RecvMsg():
			go func() {
				service.AI(msg)
			}()
		case grpmsg := <-postman.RecvGroupMsg():
			go func() {
				service.GroupAI(grpmsg)
			}()
		default:
		}
	}

}
