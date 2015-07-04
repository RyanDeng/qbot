package main

import (
	"runtime"

	"github.com/teapots/teapot"
	"hd.qiniu.com/env"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	tea := teapot.New()

	env.ConfigEnv(tea)
	env.ConfigRoutes(tea)
	env.ConfigProviders(tea)
	env.ConfigFilters(tea)
	env.ConfigJobs(tea)
	tea.Run()
}
