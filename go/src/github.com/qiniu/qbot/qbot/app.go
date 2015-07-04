package main

import (
	"encoding/json"
	"github.com/qiniu/postmans"

	"github.com/qiniu/qbot"
	"io/ioutil"
	"log"
	"os"
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

	confStr := `{"host": ":8870", "jsHost": "http://192.168.200.244:8890"}`
	gen, ok := postmans.Get("QQ")
	if !ok {
		panic("QQ not found")
	}

	postman, err := gen(confStr)
	if err != nil {
		panic(err)
	}

	service, err := qbot.NewService(cfg)
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}

	for {
		select {
		case msg := <-postman.RecvMsg():
			go func() {
				service.AI(&msg, postman)
			}()

		case grpmsg := <-postman.RecvGroupMsg():
			go func() {
				service.GroupAI(&grpmsg, postman)
			}()
		default:
		}
	}

}
