package env

import (
	"log"
	"os"
	"path/filepath"

	"github.com/teapots/gzip"
	"github.com/teapots/request-logger"
	"github.com/teapots/static-serve"
	"github.com/teapots/teapot"

	"hd.qiniu.com/filters"
)

func ConfigFilters(tea *teapot.Teapot) {
	logOut := log.New(os.Stderr, "", log.LstdFlags|log.Lmicroseconds)

	staticOption := static.StaticOption{
		AnyMethod: tea.Config.RunMode.IsDev(),
	}

	loggerOption := reqlogger.LoggerOption{
		ColorMode:     tea.Config.RunMode.IsDev(),
		LineInfo:      true,
		ShortLine:     true,
		LogStackLevel: teapot.LevelCritical,
	}

	tea.Filter(
		// gzip
		gzip.All(),

		// 所有过滤器之前抓取 panic
		teapot.RecoveryFilter(),
	)

	// 静态文件
	tea.Filter(
		static.ServeFilter("public", filepath.Join(tea.Config.RunPath, "public"), staticOption),
	)

	if tea.Config.RunMode.IsDev() || tea.Config.RunMode.IsTest() {
		// 因为 chrome  的限制
		// 开发模式下删除 Postman Header 请求前缀
		tea.Filter(
			filters.HeaderRemovePrefixFilter("Postman-"),
		)
	}

	tea.Filter(
		// 在静态文件之后加入，跳过静态文件请求
		reqlogger.ReqLoggerFilter(logOut, loggerOption),

		// 在 action 里直接返回一般请求结果
		teapot.GenericOutFilter(),
	)
}
