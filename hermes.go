//go:build linux || darwin

/*
Create: 2022/8/11
Project: Hermes
Github: https://github.com/landers1037
Copyright Renj
*/

// Package main
package main

import (
	"os"

	"github.com/JJApplication/fushin/log"
	"github.com/JJApplication/fushin/server/uds"
	"github.com/JJApplication/fushin/utils/env"
)

var logger log.Logger

func init() {
	logger = log.Logger{
		Name:   Hermes,
		Option: log.DefaultOption,
		Sync:   true,
	}
	_ = logger.Init()
}

func main() {
	hermes := HermesCore{
		AppName:   Hermes,
		Mail:      new(mail),
		EnvLoader: new(env.EnvLoader),
		UdsServer: uds.New(Hermes),
	}
	hermes.Init()
	if err := hermes.Run(); err != nil {
		logger.ErrorF("%s exit with error: %s", Hermes, err.Error())
		os.Exit(1)
	}
}
