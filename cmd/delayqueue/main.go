package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/changsongl/delay-queue/api"
	"github.com/changsongl/delay-queue/config"
	"github.com/changsongl/delay-queue/pkg/log"
	"github.com/changsongl/delay-queue/server"
	"github.com/changsongl/delay-queue/vars"
	"os"
	"strings"
)

var (
	configFile = flag.String("config.file", "../../config/config.yaml", "config file")
	configType = flag.String("config.type", "", "config type")
	env        = flag.String("env", "release", "delay queue env")

	ErrorInvalidFileType = errors.New("invalid config file type")
)

// load config file and type
func loadConfigFlags() (file string, fileType config.FileType, err error) {
	t := *configType
	f := *configFile

	// if file type is not provided, load file type from file
	if t == "" {
		extSlice := strings.Split(f, ".")
		lenExt := len(extSlice)
		if lenExt == 0 {
			return "", "", ErrorInvalidFileType
		}

		t = extSlice[lenExt-1]
	}

	t = strings.ToLower(t)

	switch t {
	case "yaml":
		return f, config.FileTypeYaml, nil
	case "json":
		return f, config.FileTypeJson, nil
	default:
		return "", "", ErrorInvalidFileType
	}
}

// load env
func loadEnv() (vars.Env, error) {
	envType := vars.Env(*env)
	if envType != vars.EnvDebug && envType != vars.EnvRelease {
		return "", errors.New(fmt.Sprintf("invalid env (%s)", envType))
	}

	return envType, nil
}

func main() {
	os.Exit(run())
}

func run() int {

	flag.Parse()
	file, fileType, err := loadConfigFlags()
	if err != nil {
		fmt.Printf("Load conifuration failed: %v\n", err)
		return 1
	}
	dqEnv, err := loadEnv()
	if err != nil {
		fmt.Printf("Load env failed: %v\n", err)
		return 1
	}

	l, err := log.New(dqEnv)
	if err != nil {
		fmt.Printf("Init log failed: %v\n", err)
		return 1
	}

	l.Info("Init configuration",
		log.String("file", file), log.String("file.type", string(fileType)))
	conf := config.New()
	err = conf.Load(file, fileType)
	if err != nil {
		l.Error(err.Error())
		return 1
	}

	dqApi := api.NewApi(l.WithModule("api"))

	l.Info("Init server",
		log.String("env", string(dqEnv)))
	s := server.New(
		server.LoggerOption(l),
		server.EnvOption(dqEnv),
		server.BeforeStartEventOption(),
		server.AfterStopEventOption(),
	)
	s.Init()
	s.RegisterRouters(dqApi.RouterFunc())

	l.Info("Run server", log.String("address", ":8080"))
	err = s.Run(":8080")
	if err != nil {
		l.Error(err.Error())
		return 1
	}

	return 0
}
