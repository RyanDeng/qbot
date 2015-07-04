package env

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/teapots/config"
	"github.com/teapots/teapot"
	"hd.qiniu.com/env/adfilter"
	"hd.qiniu.com/env/global"
	"hd.qiniu.com/env/nrop2015"
	"hd.qiniu.com/env/supportopen"
)

func LoadConfig(tea *teapot.Teapot, file string, bind bool) (fileConf *teapot.Config, err error) {
	fileConf = &teapot.Config{}
	c, err := config.LoadIniFile(filepath.Join(tea.Config.RunPath, "conf/"+file+".ini"))
	if err != nil {
		return
	}
	fileConf.Configer = c
	if bind {
		tea.ImportConfig(c)
	}
	if !tea.Config.RunMode.IsProd() {
		file += ".dev.ini"
		if tea.Config.RunMode.IsTest() {
			file += ".test.ini"
		}
		if _, err := os.Stat("conf/" + file); err == nil {
			conf, err := config.LoadIniFile(filepath.Join(tea.Config.RunPath, "conf/"+file))
			if err != nil {
				return nil, err
			}
			conf.SetParent(c)
			c = conf
			tea.ImportConfig(c)
		}
	}
	return
}

func initAdminUids(activity string) {
	global.AdminUids[activity] = make(map[uint32]bool)
	uids := strings.Split(supportopen.Env.AdminUids, ",")
	for _, uid := range uids {
		uintUid, err := strconv.ParseUint(uid, 10, 32)
		if err != nil {
			continue
		}
		global.AdminUids[activity][uint32(uintUid)] = true
	}
}

func ConfigEnv(tea *teapot.Teapot) {
	defer func() {
		// 配置文件载入以后做一些其他配置
		ConfigDB(tea)
	}()

	// 加载app配置文件
	global.Env = &global.Setting{
		Config: tea.Config,
		Teapot: tea,
	}

	conf, err := LoadConfig(tea, "app", true)
	if err != nil {
		tea.Logger().Warnf("load app config error: %s", err)
		return
	}
	config.Decode(tea.Config, global.Env)

	// 加载nrop2015配置文件
	nrop2015.Env = &nrop2015.Setting{}
	conf, err = LoadConfig(tea, "nrop2015", false)
	if err != nil {
		tea.Logger().Warnf("load nrop2015 config error: %s", err)
	} else {
		config.Decode(conf, nrop2015.Env)
	}

	// 加载supportopen配置文件
	supportopen.Env = &supportopen.Setting{}
	conf, err = LoadConfig(tea, "supportopen", false)
	if err != nil {
		tea.Logger().Warnf("load supportopen config error: %s", err)
	} else {
		config.Decode(conf, supportopen.Env)
		initAdminUids("supportopen")
	}

	// 加载adfilter配置文件
	adfilter.Env = &adfilter.Setting{}
	conf, err = LoadConfig(tea, "adfilter", false)
	if err != nil {
		tea.Logger().Warnf("load adfilter config error: %s", err)
	} else {
		config.Decode(conf, adfilter.Env)
	}
}
