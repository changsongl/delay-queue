package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/changsongl/delay-queue/api"
	"github.com/changsongl/delay-queue/bucket"
	"github.com/changsongl/delay-queue/config"
	"github.com/changsongl/delay-queue/dispatch"
	"github.com/changsongl/delay-queue/pkg/log"
	client "github.com/changsongl/delay-queue/pkg/redis"
	"github.com/changsongl/delay-queue/pool"
	"github.com/changsongl/delay-queue/queue"
	"github.com/changsongl/delay-queue/server"
	"github.com/changsongl/delay-queue/store/redis"
	"github.com/changsongl/delay-queue/timer"
	"github.com/changsongl/delay-queue/vars"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	// configuration and environment
	configFile = flag.String("config.file", "../../config/config.yaml", "config file")
	configType = flag.String("config.type", "", "config type: yaml, json")
	env        = flag.String("env", "release", "delay queue env: debug, release")
	version    = flag.Bool("version", false, "display build info")

	// ErrorInvalidFileType configuration file type is invalid
	ErrorInvalidFileType = errors.New("invalid config file type")
)

// loadConfigFlags load config file and type
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
		return f, config.FileTypeYAML, nil
	case "json":
		return f, config.FileTypeJSON, nil
	default:
		return "", "", ErrorInvalidFileType
	}
}

// load env
func loadEnv() (vars.Env, error) {
	envType := vars.Env(*env)
	if envType != vars.EnvDebug && envType != vars.EnvRelease {
		return "", fmt.Errorf("invalid env (%s)", envType)
	}

	return envType, nil
}

// main function
func main() {
	os.Exit(run())
}

// run function
func run() int {

	// parse flags
	flag.Parse()
	if *version {
		fmt.Printf(vars.BuildInfo())
		return 0
	}

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

	// get logger
	l, err := createMainLog()
	if err != nil {
		fmt.Printf("Init log failed: %v\n", err)
		return 1
	}

	fmt.Println(logo())

	// load config file
	l.Info("Init configuration",
		log.String("file", file), log.String("file.type", string(fileType)))
	conf := config.New()
	err = conf.Load(file, fileType)
	if err != nil {
		l.Error("conf.Load failed", log.Error(err))
		return 1
	}

	// print config
	l.Info("Loaded Configuration", log.String("Configuration", conf.String()))

	wg := sync.WaitGroup{}
	wg.Add(1)

	// init dispatcher of delay queue, with timer, bucket, queue, job pool components.
	disp := dispatch.NewDispatch(l,
		func() (bucket.Bucket, pool.Pool, queue.Queue, timer.Timer) {
			cli := client.New(conf.Redis)
			s := redis.NewStore(cli)

			b := bucket.New(s, l, conf.DelayQueue.BucketSize, conf.DelayQueue.BucketName)
			if maxFetchNum := conf.DelayQueue.BucketMaxFetchNum; maxFetchNum != 0 {
				b.SetMaxFetchNum(maxFetchNum)
			}

			p := pool.New(s, l)
			q := queue.New(s, l, conf.DelayQueue.QueueName)
			t := timer.New(
				l, time.Duration(conf.DelayQueue.TimerFetchInterval)*time.Millisecond,
				time.Duration(conf.DelayQueue.TimerFetchDelay)*time.Millisecond,
			)
			return b, p, q, t
		},
	)
	go func() {
		disp.Run()
		wg.Done()
	}()

	// run http server to receive requests from user
	dqAPI := api.NewAPI(l, disp)
	l.Info("Init server",
		log.String("env", string(dqEnv)))
	s := server.New(
		server.LoggerOption(l),
		server.EnvOption(dqEnv),
		server.BeforeStartEventOption(),
		server.AfterStopEventOption(),
	)
	s.Init()
	s.RegisterRouters(dqAPI.RouterFunc())
	err = s.Run(conf.DelayQueue.BindAddress)
	if err != nil {
		l.Error("s.Run failed", log.Error(err))
		return 1
	}

	wg.Wait()

	return 0
}

// createMainLog create logger with name "main"
func createMainLog() (log.Logger, error) {
	l, err := log.New()
	if err != nil {
		return nil, err
	}

	return l.WithModule("main"), nil
}

// logo of delay queue
func logo() string {
	return "" +
		"  ____       _                ___                        \n" +
		" |  _ \\  ___| | __ _ _   _   / _ \\ _   _  ___ _   _  ___ \n" +
		" | | | |/ _ \\ |/ _` | | | | | | | | | | |/ _ \\ | | |/ _ \\\n" +
		" | |_| |  __/ | (_| | |_| | | |_| | |_| |  __/ |_| |  __/\n" +
		" |____/ \\___|_|\\__,_|\\__, |  \\__\\_\\\\__,_|\\___|\\__,_|\\___|\n" +
		"                     |___/  "
}
