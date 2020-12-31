package main

import (
	"fmt"
	"github.com/changsongl/delay-queue/pkg/log"
	"github.com/changsongl/delay-queue/server"
	"github.com/changsongl/delay-queue/vars"
	"os"
)

const (
	env = vars.EnvDebug
)

func main() {
	os.Exit(run())
}

func run() int {
	l, err := log.New(env)
	if err != nil {
		fmt.Println(err)
		return 1
	}

	s := server.New(
		server.LoggerOption(l),
		server.EnvOption(env),
	)

	s.Init()

	err = s.Run(":8080")
	if err != nil {
		l.Error(err.Error())
		return 1
	}

	return 0
}
